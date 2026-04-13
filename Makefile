.PHONY: help install install-api install-web api web db-up db-down db-logs db-ps db-init db-seed db-init-docker db-seed-docker db-docker-run

COMPOSE := docker compose -f database/docker-compose.yml
API_DIR := apps/api
WEB_DIR := apps/web
DB_CONTAINER := db
DB_USER := user
DB_NAME := comune_dev
SCRIPT ?=

ifneq ($(filter db-docker-run,$(MAKECMDGOALS)),)
SCRIPT := $(or $(SCRIPT),$(word 2,$(MAKECMDGOALS)))
ifneq ($(strip $(word 2,$(MAKECMDGOALS))),)
$(eval $(word 2,$(MAKECMDGOALS)):;@:)
endif
endif

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
	@echo "  make db-init      Apply the SQL schema using local psql and DATABASE_URL"
	@echo "  make db-seed      Apply the local development seed data using local psql and DATABASE_URL"
	@echo "  make db-init-docker Apply the SQL schema using psql inside the Docker db container"
	@echo "  make db-seed-docker Apply the local development seed using psql inside the Docker db container"
	@echo "  make db-docker-run SCRIPT=database/file.sql"
	@echo "  make db-docker-run database/file.sql"

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

db-seed:
	psql "$(DATABASE_URL)" -f database/seed.sql

db-init-docker:
	$(COMPOSE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -f /dev/stdin < database/gated-community-schema.sql

db-seed-docker:
	$(COMPOSE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -f /dev/stdin < database/seed.sql

db-docker-run:
	@if [ -z "$(SCRIPT)" ]; then \
		echo "Usage: make db-docker-run SCRIPT=database/file.sql"; \
		echo "   or: make db-docker-run database/file.sql"; \
		exit 1; \
	fi
	@if [ ! -f "$(SCRIPT)" ]; then \
		echo "SQL script not found: $(SCRIPT)"; \
		exit 1; \
	fi
	$(COMPOSE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -f /dev/stdin < "$(SCRIPT)"
