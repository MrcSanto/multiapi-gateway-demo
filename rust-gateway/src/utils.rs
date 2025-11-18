use axum::http::HeaderMap;

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

pub fn filter_headers(headers: &HeaderMap) -> HeaderMap {
    let mut filtered = HeaderMap::new();
    
    for (name, value) in headers.iter() {
        let name_lower = name.as_str().to_lowercase();
        if !HOP_HEADERS.contains(&name_lower.as_str()) {
            filtered.insert(name.clone(), value.clone());
        }
    }
    
    filtered
}