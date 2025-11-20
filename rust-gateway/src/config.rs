use std::env;

#[derive(Clone, Debug)]
pub struct Config {
    pub upstreams: Vec<String>,
    pub rate_limit: usize,
    pub window_sec: u64,
    pub cb_fail_threshold: usize,
    pub cb_cooldown_sec: u64,
    pub port: String,
}

impl Config {
    pub fn from_env() -> Self {
        let upstreams: Vec<String> = env::var("UPSTREAMS")
            .unwrap_or_else(|_| {
                "http://go-api:8000,http://node-api:8000,http://python-api:8000".to_string()
            })
            .split(',')
            .filter_map(|s| {
                let trimmed = s.trim().trim_end_matches('/');
                if trimmed.is_empty() {
                    None
                } else {
                    Some(trimmed.to_string())
                }
            })
            .collect();
        // printando os upstreams para visualização
        for (i, up) in upstreams.iter().enumerate() {
            println!("Upstream {} -> {}", i, up);
        }

        Self {
            upstreams,
            rate_limit: Self::parse_env("RATE_LIMIT", 30),
            window_sec: Self::parse_env("WINDOW_SEC", 60),
            cb_fail_threshold: Self::parse_env("CB_FAIL_THRESHOLD", 3),
            cb_cooldown_sec: Self::parse_env("CB_COOLDOWN_SEC", 10),
            port: env::var("PORT").unwrap_or_else(|_| "8080".to_string()),
        }
    }

    // método para ler uma variavel de ambiente, tentar converter se possível,
    // caso de erro, utilizamos o valor default definido
    fn parse_env<T: std::str::FromStr>(key: &str, default: T) -> T {
        env::var(key)
            .ok()
            .and_then(|v| v.parse().ok())
            .unwrap_or(default)
    }
}