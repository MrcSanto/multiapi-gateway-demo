mod config;
mod handlers;
mod middleware;
mod state;
mod utils;

use axum::{middleware as axum_middleware, routing::{any, get}, Router};
use config::Config;
use handlers::{health_handler, proxy_handler};
use middleware::rate_limit_middleware;
use state::AppState;
use std::sync::Arc;

#[tokio::main]
async fn main() {
    // Carregar configuração
    let config = Config::from_env();
    println!("  Configuração carregada:");
    println!("   Upstreams: {:?}", config.upstreams);
    println!("   Rate Limit: {} req/{} sec", config.rate_limit, config.window_sec);
    println!("   Circuit Breaker: {} falhas, {} sec cooldown", 
             config.cb_fail_threshold, config.cb_cooldown_sec);

    // Criar estado compartilhado
    let state = Arc::new(AppState::new(config.clone()));

    // Criar router
    let app = Router::new()
        .route("/_health", get(health_handler))
        .fallback(any(proxy_handler))
        .layer(axum_middleware::from_fn_with_state(
            state.clone(),
            rate_limit_middleware,
        ))
        .with_state(state);

    // Bind no endereço
    let addr = format!("0.0.0.0:{}", config.port);
    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .expect("Falha ao fazer bind no endereço");

    println!("  API Gateway listening on http://{}", addr);

    // Iniciar servidor
    axum::serve(listener, app)
        .await
        .expect("Falha ao iniciar servidor");
}