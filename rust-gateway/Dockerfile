# Build stage
FROM rust:1.84-slim AS builder

# Instalar dependências de build
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar arquivos do projeto
COPY Cargo.toml Cargo.lock ./
COPY src ./src

# Build release
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Instalar apenas dependências de runtime
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar binário do build stage
COPY --from=builder /app/target/release/rust-gateway /app/rust-gateway

# Variáveis de ambiente
ENV UPSTREAMS="http://mock1:8000,http://mock2:8000,http://mock3:8000"
ENV RATE_LIMIT=30
ENV WINDOW_SEC=60
ENV CB_FAIL_THRESHOLD=3
ENV CB_COOLDOWN_SEC=10
ENV PORT=8000

EXPOSE 8000

CMD ["/app/rust-gateway"]