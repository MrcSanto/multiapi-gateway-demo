use crate::{state::AppState, state::CircuitState, utils::filter_headers};
use axum::{
    body::Body,
    extract::{Request, State},
    http::{Method, StatusCode},
    response::{IntoResponse, Response},
};
use std::sync::Arc;

pub async fn proxy_handler(
    State(state): State<Arc<AppState>>,
    method: Method,
    req: Request,
) -> Result<Response, (StatusCode, &'static str)> {
    // Escolher upstream
    let idx = state
        .pick_upstream()
        .await
        .map_err(|_| (StatusCode::SERVICE_UNAVAILABLE, "No healthy upstreams"))?;

    let upstream = &state.config.upstreams[idx];
    let path = req.uri().path();
    let query = req.uri().query().unwrap_or("");

    println!("→ Gateway encaminhando para instância: {}", upstream);

    // Construir URL
    let url = if query.is_empty() {
        format!("{}{}", upstream, path)
    } else {
        format!("{}{}?{}", upstream, path, query)
    };

    // Filtrar headers
    let headers = filter_headers(req.headers());

    // Ler body
    let body_bytes = axum::body::to_bytes(req.into_body(), usize::MAX)
        .await
        .map_err(|_| (StatusCode::BAD_REQUEST, "Failed to read body"))?;

    // Fazer request para upstream
    let result = state
        .client
        .request(method, &url)
        .headers(headers)
        .body(body_bytes.to_vec())
        .send()
        .await;

    let resp = match result {
        Ok(r) => r,
        Err(_) => {
            state.on_failure(idx).await;
            return Err((StatusCode::BAD_GATEWAY, "Upstream error"));
        }
    };

    // Circuit breaker logic
    if resp.status().as_u16() >= 500 {
        state.on_failure(idx).await;
    } else {
        state.on_success(idx).await;
    }

    // Construir resposta
    let status = StatusCode::from_u16(resp.status().as_u16()).unwrap();
    let resp_headers = filter_headers(resp.headers());
    let body = resp.bytes().await.unwrap_or_default();

    let mut response = Response::new(Body::from(body));
    *response.status_mut() = status;
    *response.headers_mut() = resp_headers;

    Ok(response)
}

pub async fn health_handler(State(state): State<Arc<AppState>>) -> impl IntoResponse {
    let cbs = state.circuit_breakers.read().await;
    let cb_status: Vec<String> = cbs
        .iter()
        .map(|cb| {
            format!(
                "{{state: {:?}, fails: {}}}",
                match cb.state {
                    CircuitState::Closed => "closed",
                    CircuitState::Open => "open",
                    CircuitState::Half => "half",
                },
                cb.fail_count
            )
        })
        .collect();

    axum::Json(serde_json::json!({
        "upstreams": state.config.upstreams,
        "circuit_breakers": cb_status,
        "rate_limit": state.config.rate_limit,
        "window_sec": state.config.window_sec,
    }))
}