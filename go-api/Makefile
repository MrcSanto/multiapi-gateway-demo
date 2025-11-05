COMPOSE := docker compose
IMAGE := docker image
APP := go_app
DB := go_db
HOST ?= http://localhost:8000

up:
	$(COMPOSE) up -d --build
	$(IMAGE) prune -f
	@echo "API em: $(HOST)"

logs:
	$(COMPOSE) logs -f $(APP)

down:
	$(COMPOSE) down -v

rebuild:
	$(COMPOSE) build --no-cache
	$(IMAGE) prune -f
	$(COMPOSE) up -d

ps:
	$(COMPOSE) ps

psql:
	$(COMPOSE) exec $(DB) psql -U $$POSTGRES_USER -d $$POSTGRES_DB
	
	
