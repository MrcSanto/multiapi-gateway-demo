use crate::config::Config;
use axum::http::StatusCode;
use reqwest::Client;
use std::{
    collections::{HashMap, VecDeque},
    time::{Duration, Instant},
};
use tokio::sync::RwLock;

#[derive(Clone, Debug)]
pub enum CircuitState {
    Closed,
    Open,
    Half,
}

pub struct CircuitBreaker {
    pub state: CircuitState,
    pub fail_count: usize,
    pub opened_at: Option<Instant>,
}

pub struct AppState {
    pub config: Config,
    pub client: Client,
    pub rate_limits: RwLock<HashMap<String, VecDeque<Instant>>>,
    pub round_robin: RwLock<usize>,
    pub circuit_breakers: RwLock<Vec<CircuitBreaker>>,
}

impl AppState {
    pub fn new(config: Config) -> Self {
        let num_upstreams = config.upstreams.len();
        let circuit_breakers = (0..num_upstreams)
            .map(|_| CircuitBreaker {
                state: CircuitState::Closed,
                fail_count: 0,
                opened_at: None,
            })
            .collect();

        Self {
            config,
            client: Client::builder()
                .timeout(Duration::from_secs(10))
                .redirect(reqwest::redirect::Policy::none())
                .build()
                .unwrap(),
            rate_limits: RwLock::new(HashMap::new()),
            round_robin: RwLock::new(0),
            circuit_breakers: RwLock::new(circuit_breakers),
        }
    }

    pub async fn check_rate_limit(&self, ip: String) -> bool {
        let mut limits = self.rate_limits.write().await;
        let now = Instant::now();
        let window = Duration::from_secs(self.config.window_sec);

        let queue = limits.entry(ip).or_insert_with(VecDeque::new);

        // Remove timestamps fora da janela
        while let Some(&front) = queue.front() {
            if now.duration_since(front) > window {
                queue.pop_front();
            } else {
                break;
            }
        }

        if queue.len() >= self.config.rate_limit {
            return true; // Rate limit exceeded
        }

        queue.push_back(now);
        false
    }

    pub async fn pick_upstream(&self) -> Result<usize, StatusCode> {
        let mut cbs = self.circuit_breakers.write().await;
        let mut rr = self.round_robin.write().await;

        for _ in 0..self.config.upstreams.len() {
            let idx = *rr;
            *rr = (*rr + 1) % self.config.upstreams.len();

            let cb = &mut cbs[idx];

            match cb.state {
                CircuitState::Open => {
                    if let Some(opened_at) = cb.opened_at {
                        if opened_at.elapsed().as_secs() >= self.config.cb_cooldown_sec {
                            cb.state = CircuitState::Half;
                            return Ok(idx);
                        }
                    }
                }
                _ => return Ok(idx),
            }
        }

        Err(StatusCode::SERVICE_UNAVAILABLE)
    }

    pub async fn on_success(&self, idx: usize) {
        let mut cbs = self.circuit_breakers.write().await;
        let cb = &mut cbs[idx];
        cb.fail_count = 0;
        cb.state = CircuitState::Closed;
    }

    pub async fn on_failure(&self, idx: usize) {
        let mut cbs = self.circuit_breakers.write().await;
        let cb = &mut cbs[idx];
        cb.fail_count += 1;

        let should_open = matches!(cb.state, CircuitState::Half)
            || cb.fail_count >= self.config.cb_fail_threshold;

        if should_open {
            cb.state = CircuitState::Open;
            cb.opened_at = Some(Instant::now());
        }
    }
}