.PHONY: help install install-api install-web api web db-up db-down db-logs db-ps db-init

COMPOSE := docker compose -f database/docker-compose.yml
API_DIR := apps/api
WEB_DIR := apps/web

help:
	@echo "Available targets:"
	@echo "  make install      Install backend and web dependencies"
	@echo "  make install-api  Download Go dependencies"
	@echo "  make install-web  Install web dependencies with pnpm"
	@echo "  make api          Run the Go API against the local dev database"
	@echo "  make web          Run the web UI in dev mode"
	@echo "  make db-up        Start the Postgres container"
	@echo "  make db-down      Stop the Postgres container"
	@echo "  make db-logs      Tail Postgres logs"
	@echo "  make db-ps        Show Postgres container status"
	@echo "  make db-init      Apply the SQL schema to the local dev database"

install: install-api install-web

install-api:
	cd $(API_DIR) && go mod download

install-web:
	cd $(WEB_DIR) && pnpm install

api:
	cd $(API_DIR) && go run ./cmd/api

web:
	cd $(WEB_DIR) && pnpm dev

db-up:
	$(COMPOSE) up -d

db-down:
	$(COMPOSE) down

db-logs:
	$(COMPOSE) logs -f db

db-ps:
	$(COMPOSE) ps

db-init:
	psql "$(DATABASE_URL)" -f database/gated-community-schema.sql
