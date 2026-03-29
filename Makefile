SHELL := /bin/bash
.ONESHELL:
.DEFAULT_GOAL := help

-include .env
export

# ==========================================
# VARIÁVEIS DA APLICAÇÃO (GO)
# ==========================================
APP_NAME ?= go-lab
BIN_DIR ?= ./bin
BIN_PATH ?= $(BIN_DIR)/api

# ==========================================
# VARIÁVEIS DOS CONTAINERS (DOCKER COMPOSE)
# ==========================================
BIND_HOST ?= 127.0.0.1
PG_CONTAINER ?= postgres-abilico
RABBITMQ_CONTAINER ?= app_rabbitmq

# Variáveis do Postgres
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_PATH=./internal/migrations


# ==========================================
# COMANDOS DISPONÍVEIS
# ==========================================
.PHONY: help deps down reset ps logs \
        createdb dropdb psql \
        status createmigration migrateup migratedown migrate \
        build run

help: ## Mostra os comandos disponíveis
	@echo ""
	@echo "Uso: make <comando>"
	@echo ""
	@awk 'BEGIN {FS=":.*## "}; \
	/^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0,5)} \
	/^[a-zA-Z0-9_.-]+:.*## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' \
	$(MAKEFILE_LIST)
	@echo ""

# ------------------------------------------
# INFRAESTRUTURA (DOCKER COMPOSE)
# ------------------------------------------
deps: ## Sobe Postgres e RabbitMQ via Docker Compose (em background)
	@docker compose up -d
	@echo "✅ Dependências iniciadas (Postgres e RabbitMQ)"

down: ## Para os containers (sem remover dados)
	@docker compose stop
	@echo "🛑 Containers parados"

reset: ## Remove containers e volumes (CUIDADO: perde dados do banco e filas)
	@docker compose down -v
	@echo "🧨 Containers e volumes removidos"

ps: ## Lista o status dos containers
	@docker compose ps

logs: ## Mostra os logs dos containers em tempo real
	@docker compose logs -f

# ------------------------------------------
# BANCO DE DADOS (POSTGRES)
# ------------------------------------------
createdb: ## Cria o banco de dados dentro do container
	@docker exec -it $(PG_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb: ## Deleta o banco de dados dentro do container
	@docker exec -it $(PG_CONTAINER) dropdb $(DB_NAME)

psql: ## Acessa o terminal interativo do Postgres
	@docker exec -it $(PG_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

# ------------------------------------------
# MIGRATIONS & SQLC
# ------------------------------------------
status: ## Mostra URL do banco e caminho das migrations (debug)
	@echo "MIGRATIONS_PATH=$(MIGRATIONS_PATH)"
	@echo "DB_URL=$(DB_URL)"

createmigration: ## Cria uma nova migration (ex: make createmigration name=init)
	@test -n "$(name)" || (echo "❌ use: make createmigration name=nome_da_migration" && exit 1)
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) $(name)

migrateup: ## Aplica as migrations no banco (up)
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migratedown: ## Reverte as migrations no banco (down)
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate: migrateup ## Alias rápido para aplicar as migrations

# ------------------------------------------
# APLICAÇÃO GO
# ------------------------------------------
build: ## Compila a API em um binário
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_PATH) ./cmd/api
	@echo "✅ Binário gerado em: $(BIN_PATH)"

run: ## Roda a API em modo de desenvolvimento
	@go run ./cmd/api