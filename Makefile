.PHONY: build up down logs clean restart gateway-logs go-logs db-logs test health

# Cores para output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m

build: ## Build de todos os serviços
	@echo -e "$(GREEN)Building all services...$(NC)"
	docker compose build

up: ## Inicia todos os serviços
	@echo -e "$(GREEN)Starting all services...$(NC)"
	docker compose up -d
	@echo -e "$(GREEN)Services started!$(NC)"
	@echo -e "$(YELLOW)Gateway: http://localhost:8080$(NC)"
	@echo -e "$(YELLOW)Go API:  http://localhost:8001/go$(NC)"
	@echo -e "$(YELLOW)DB:      localhost:5432$(NC)"

down: ## Para todos os serviços
	@echo -e "$(RED)Stopping all services...$(NC)"
	docker compose down

restart: down up ## Reinicia todos os serviços

clean: ## Para serviços e remove volumes
	@echo -e "$(RED)Stopping and removing all data...$(NC)"
	docker compose down -v
	docker system prune -f

logs: ## Mostra logs de todos os serviços
	docker compose logs -f

gateway-logs: ## Mostra logs do gateway
	docker compose logs -f rust_gateway

go-logs: ## Mostra logs da Go API
	docker compose logs -f go_app

db-logs: ## Mostra logs do banco
	docker compose logs -f go_db

ps: ## Lista status dos containers
	@docker compose ps

health: ## Verifica health dos serviços
	@echo -e "$(GREEN)Checking services health...$(NC)"
	@echo -e "\n$(YELLOW)Gateway Health:$(NC)"
	@curl -s http://localhost:8080/ | jq . || echo "Gateway not responding"
	@echo -e "\n$(YELLOW)Go API Health:$(NC)"
	@curl -s http://localhost:8001/go/healthcheck | jq . || echo "Go API not responding"
	@echo -e "\n$(YELLOW)Database:$(NC)"
	@docker exec go_db pg_isready -U postgres && echo "Database is ready" || echo "Database not ready"

test: ## Testa o gateway com requisições
	@echo -e "$(GREEN)Testing gateway...$(NC)"
	@echo -e "\n$(YELLOW)1. Testing health endpoint:$(NC)"
	@curl -s http://localhost:8080/ | jq .
	@echo -e "\n$(YELLOW)2. Testing proxy to Go API:$(NC)"
	@curl -s http://localhost:8080/go/users | jq .
	@echo -e "\n$(YELLOW)3. Testing rate limit (10 requests):$(NC)"
	@for i in {1..10}; do \
		curl -s -o /dev/null -w "Request $i: %{http_code}\n" http://localhost:8080/go/users; \
	done

dev: ## Modo desenvolvimento (build + up + logs)
	@make build
	@make up
	@make logs

rebuild: ## Rebuild completo (limpa e reconstrói)
	@make clean
	@make build
	@make up