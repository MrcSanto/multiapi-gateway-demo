use axum::{
    body::Body,
    extract::{Request, State},
    http::{HeaderMap, Method, StatusCode},
    middleware::{self, Next},
    response::{IntoResponse, Response},
    routing::{any, get},
    Router,
};
use reqwest::Client;
use std::{
    collections::{HashMap, VecDeque},
    sync::Arc,
    time::{Duration, Instant},
};
use tokio::sync::RwLock;

// config
#[derive(Clone)]
struct Config {
    upstreams: Vec<String>,
    rate_limit: usize,
    window_sec: u64,
    cb_fail_threshold: usize,
    cb_cooldown_sec:  u64,
}

impl Config {
    fn from_env() -> Self {
        let upstreams: Vec<String> = std::env::var("UPSTREAMS")
            .unwrap_or_else(|_| "http://go-api:8000,http://node-api:8000,http://python-api:8000".to_string())
            .split(",")
            .filter_map(|s| {
                let trimmed = s.trim().trim_end_matches("/");
                if trimmed.is_empty() {
                    None
                } else {
                    Some(trimmed.to_string())
                }
            })
            .collect();

        Self {
            upstreams,
            rate_limit: std::env::var("RATE_LIMIT")
                .unwrap_or_else(|_|"30".to_string())
                .parse()
                .unwrap_or(30),
            window_sec: std::env::var("WINDOW_SEC")
                .unwrap_or_else(|_|"60".to_string())
                .parse()
                .unwrap_or(60),
            cb_fail_threshold: std::env::var("CB_FAIL_THRESHOLD")
                .unwrap_or_else(|_|"3".to_string())
                .parse()
                .unwrap_or(3),
            cb_cooldown_sec: std::env::var("CB_COOLDOWN_SEC")
                .unwrap_or_else(|_|"10".to_string())
                .parse()
                .unwrap_or(10),
        }

    }   
}

// estados
#[derive(Clone)]
enum CircuitState {
    Closed, 
    Open, 
    Half
}

struct CircuitBreaker {
    state: CircuitState,
    fail_count: usize,
    opened_at: Option<Instant>
}

struct AppState {
    config: Config,
    client: Client,
    rate_limits: RwLock<HashMap<String, VecDeque<Instant>>>,
    round_robin: RwLock<usize>,
    circuit_breakers: RwLock<Vec<CircuitBreaker>>
}

impl AppState {
    fn new(config: Config) -> Self {
        let num_upstreams: usize = config.upstreams.len();
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

    async fn check_limit_rate(&self, ip: String) -> bool {
        let mut limits: tokio::sync::RwLockWriteGuard<'_, HashMap<String, VecDeque<Instant>>> = self.rate_limits.write().await;
        let now: Instant = Instant::now();
        let window: Duration = Duration::from_secs(self.config.window_sec);

        let queue: &mut VecDeque<Instant> = limits.entry(ip).or_insert_with(VecDeque::new);
        

        // remove os timestamps de fora da janela
        while let Some(&front) = queue.front() {
            if now.duration_since(front) > window {
                queue.pop_front();
            } else {
                break;
            }
        }

        if queue.len() >= self.config.rate_limit {
            return true; // excedemos o rate limit
        }

        queue.push_back(now);
        false
    }

    async fn pick_upstream(&self) -> Result<usize, StatusCode> {
        let mut cbs: tokio::sync::RwLockWriteGuard<'_, Vec<CircuitBreaker>> = self.circuit_breakers.write().await;
        let mut rr: tokio::sync::RwLockWriteGuard<'_, usize> = self.round_robin.write().await;

        for _ in 0..self.config.upstreams.len() {
            let idx: usize = *rr;
            *rr = (*rr + 1) % self.config.upstreams.len();

            let cb: &mut CircuitBreaker = &mut cbs[idx];

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

    async fn on_sucess(&self, idx: usize) {
        let mut cbs: tokio::sync::RwLockWriteGuard<'_, Vec<CircuitBreaker>> = self.circuit_breakers.write().await;
        let cb: &mut CircuitBreaker = &mut cbs[idx];

        cb.fail_count = 0;
        cb.state = CircuitState::Closed;
    }

    async fn on_failure(&self, idx: usize) {
        let mut cbs: tokio::sync::RwLockWriteGuard<'_, Vec<CircuitBreaker>> = self.circuit_breakers.write().await;
        let cb: &mut CircuitBreaker = &mut cbs[idx];

        cb.fail_count += 1;

        let should_open: bool = matches!(cb.state, CircuitState::Half)
            || cb.fail_count >= self.config.cb_fail_threshold;

        if should_open {
            cb.state = CircuitState::Open;
            cb.opened_at = Some(Instant::now())
        }
    }

}

// middleware
async fn rate_limit_middleware(
    State(state): State<Arc<AppState>>,
    req: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    let ip = extract_client_ip(&req);

    if state.check_limit_rate(ip).await {
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    Ok(next.run(req).await)
}

fn extract_client_ip(req: &Request) -> String {
    req.headers()
        .get("x-forwarded-for")
        .and_then(|v| v.to_str().ok())
        .and_then(|s| s.split(',').next())
        .map(|s| s.trim().to_string())
        .unwrap_or_else(|| "unknown".to_string())
}

//handlers
async fn proxy_handler(
    State(state): State<Arc<AppState>>,
    method: Method,
    req: Request
) -> Result<Response, (StatusCode, &'static str)> {
    let idx: usize = state
        .pick_upstream()
        .await
        .map_err(|_| (StatusCode::SERVICE_UNAVAILABLE, "No healthy upstreams"))?;

    let upstream: &String = &state.config.upstreams[idx];
    let path: &str = req.uri().path();
    let query: &str = req.uri().query().unwrap_or("");

    let url: String = if query.is_empty() {
        format!("{}{}", upstream, path)
    } else {
        format!("{}{}?{}", upstream, path, query)
    };

    // filtrando headers hop-by-hop
    let headers: HeaderMap = filter_headers(req.headers());

    // fazer request upstream
    let body_bytes = axum::body::to_bytes(req.into_body(), usize::MAX)
        .await
        .map_err(|_| (StatusCode::BAD_REQUEST, "Failed to read body"))?;

    // fazer request upstream
    let result: Result<reqwest::Response, reqwest::Error> = state
        .client
        .request(method, &url)
        .headers(headers)
        .body(body_bytes.to_vec())
        .send()
        .await;

    let resp: reqwest::Response = match result {
        Ok(r) => r,
        Err(_) => {
            state.on_failure(idx).await;
            return Err((StatusCode::BAD_GATEWAY, "Upstream error"))
        }
    };

    // logica do circuit breaker
    if resp.status().as_u16() >= 500 {
        state.on_failure(idx).await;
    } else {
        state.on_sucess(idx).await;
    }

    // construindo resposta
    let status: StatusCode = StatusCode::from_u16(resp.status().as_u16()).unwrap();
    let resp_headers: HeaderMap = filter_headers(resp.headers());
    let body: axum::body::Bytes = resp.bytes().await.unwrap_or_default();

    let mut response: Response<Body> = Response::new(Body::from(body));
    *response.status_mut() = status;
    *response.headers_mut() = resp_headers;

    Ok(response)
}

fn filter_headers(headers: &HeaderMap) -> HeaderMap {
    const HOP_HEADERS: &[&str] = &[
        "connection",
        "keep-alive",
        "proxy-authenticate",
        "proxy-authorization",
        "te",
        "trailers",
        "transfer-encoding",
        "upgrade",
    ];

    let mut filtered = HeaderMap::new();
    for (name, value) in headers.iter() {
        if !HOP_HEADERS.contains(&name.as_str().to_lowercase().as_str()) {
            filtered.insert(name.clone(), value.clone());
        }
    }

    filtered
}

async fn health_handler(State(state): State<Arc<AppState>>) -> impl IntoResponse {
    let cbs: tokio::sync::RwLockReadGuard<'_, Vec<CircuitBreaker>> = state.circuit_breakers.read().await;
    let cb_status: Vec<String> = cbs
        .iter()
        .map(|cb| {
            format!(
                "{{state: {:?}, fails: {}}}",
                match cb.state {
                    CircuitState::Closed => "closed",
                    CircuitState::Half => "half",
                    CircuitState::Open => "open"
                },
                cb.fail_count
            )
        })
        .collect();

    axum::Json(serde_json::json!({
        "upstreams": state.config.upstreams,
        "circuit_breakers": cb_status,
        "rate_limit": state.config.rate_limit,
        "window_sec": state.config.window_sec
    }))
}

#[tokio::main]
async fn main() {
    let config: Config = Config::from_env();
    let state: Arc<AppState> = Arc::new(AppState::new(config));

    let app: Router = Router::new()
        .route("/_health", get(health_handler))
        .fallback(any(proxy_handler))
        .layer(middleware::from_fn_with_state(
            state.clone(), 
            rate_limit_middleware
        ))
        .with_state(state);

    let port: String = std::env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    let addr: String = format!("0.0.0.0:{}", port);

    let listener: tokio::net::TcpListener = tokio::net::TcpListener::bind(&addr)
        .await
        .unwrap();

    println!("API Gateway listening on http://{}", addr);

    axum::serve(listener, app).await.unwrap();
}
