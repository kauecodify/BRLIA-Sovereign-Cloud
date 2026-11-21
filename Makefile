# BRLIA Sovereign Cloud - Makefile
# Comandos utilitarios para desenvolvimento e deploy

.PHONY: help build run stop clean test lint docker docker-up docker-down docker-build docker-push migrate seed backup restore

# Cores para output
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

# Variaveis de ambiente
ENV ?= development
BIN_DIR ?= ./bin
DOCKER_REGISTRY ?= ghcr.io/kauecodify
DOCKER_TAG ?= latest

help: ## Mostra esta ajuda
	@echo "BRLIA Sovereign Cloud - Comandos Disponiveis"
	@echo ""
	@echo "Comandos de Build:"
	@echo "  make build              - Builda todos os servicos"
	@echo "  make build-control      - Builda apenas o Control API"
	@echo "  make build-waterfall    - Builda apenas o Waterfall Engine"
	@echo ""
	@echo "Comandos de Execucao:"
	@echo "  make run                - Executa todos os servicos localmente"
	@echo "  make run-control        - Executa apenas o Control API"
	@echo "  make run-waterfall      - Executa apenas o Waterfall Engine"
	@echo ""
	@echo "Comandos Docker:"
	@echo "  make docker-up          - Sobe todos os containers"
	@echo "  make docker-down        - Para todos os containers"
	@echo "  make docker-build       - Builda imagens Docker"
	@echo "  make docker-push        - Publica imagens Docker"
	@echo ""
	@echo "Comandos de Database:"
	@echo "  make migrate            - Aplica migracoes do banco"
	@echo "  make seed               - Popula banco com dados iniciais"
	@echo ""
	@echo "Comandos de Manutencao:"
	@echo "  make clean              - Remove arquivos temporarios"
	@echo "  make test               - Executa todos os testes"
	@echo "  make lint               - Executa linting do codigo"
	@echo "  make backup             - Cria backup do sistema"
	@echo "  make restore            - Restaura de backup"
	@echo "  make stop               - Para todos os servicos"

# Build
build: build-control build-waterfall ## Builda todos os servicos

build-control: ## Builda o Control API
	@echo "$(GREEN)Building Control API...$(NC)"
	cd services/control-api && go build -o $(BIN_DIR)/control-api main.go
	@echo "$(GREEN)Control API built successfully$(NC)"

build-waterfall: ## Builda o Waterfall Engine
	@echo "$(GREEN)Building Waterfall Engine...$(NC)"
	cd services/waterfall && go build -o $(BIN_DIR)/waterfall main.go
	@echo "$(GREEN)Waterfall Engine built successfully$(NC)"

# Run
run: run-control run-waterfall ## Executa todos os servicos

run-control: ## Executa o Control API
	@echo "$(GREEN)Starting Control API...$(NC)"
	$(BIN_DIR)/control-api

run-waterfall: ## Executa o Waterfall Engine
	@echo "$(GREEN)Starting Waterfall Engine...$(NC)"
	$(BIN_DIR)/waterfall

# Docker
docker-up: ## Sobe todos os containers
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Containers started. Access:"
	@echo "  - Control API: http://localhost:8080"
	@echo "  - Waterfall:   http://localhost:8081"
	@echo "  - Admin Console: http://localhost:3001"
	@echo "  - Prometheus:  http://localhost:9090"
	@echo "  - Grafana:     http://localhost:3000"
	@echo "  - Keycloak:    http://localhost:8082"
	@echo "  - Vault:       http://localhost:8200$(NC)"

docker-down: ## Para todos os containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker-compose down

docker-build: ## Builda imagens Docker
	@echo "$(GREEN)Building Docker images...$(NC)"
	docker-compose build

docker-push: ## Publica imagens Docker
	@echo "$(GREEN)Pushing Docker images to $(DOCKER_REGISTRY)...$(NC)"
	docker tag brlia-control-api:latest $(DOCKER_REGISTRY)/brlia-control-api:$(DOCKER_TAG)
	docker tag brlia-waterfall:latest $(DOCKER_REGISTRY)/brlia-waterfall:$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/brlia-control-api:$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/brlia-waterfall:$(DOCKER_TAG)
	@echo "$(GREEN)Images pushed successfully$(NC)"

# Database
migrate: ## Aplica migracoes do banco
	@echo "$(GREEN)Running database migrations...$(NC)"
	# Implementar com goose, flyway ou similar
	@echo "$(YELLOW)Database migrations not yet implemented$(NC)"

seed: ## Popula banco com dados iniciais
	@echo "$(GREEN)Seeding database...$(NC)"
	# Implementar seeding
	@echo "$(YELLOW)Database seeding not yet implemented$(NC)"

# Testes
test: test-unit test-integration test-e2e ## Executa todos os testes

test-unit: ## Executa testes unitarios
	@echo "$(GREEN)Running unit tests...$(NC)"
	cd services/control-api && go test ./... -v -short
	cd services/waterfall && go test ./... -v -short

# Manutencao
clean: ## Remove arquivos temporarios
	@echo "$(YELLOW)Cleaning up...$(NC)"
	rm -rf $(BIN_DIR)/*
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(NC)"

stop: ## Para todos os servicos
	@echo "$(YELLOW)Stopping all services...$(NC)"
	pkill -f "control-api\|waterfall" || true
	docker-compose down || true
	@echo "$(GREEN)All services stopped$(NC)"

lint: ## Executa linting do codigo
	@echo "$(GREEN)Running linter...$(NC)"
	cd services/control-api && go vet ./...
	cd services/waterfall && go vet ./...
	@echo "$(GREEN)Linting complete$(NC)"

backup: ## Cria backup do sistema
	@echo "$(GREEN)Creating backup...$(NC)"
	mkdir -p backups
	BACKUP_FILE="backups/brlia-backup-$(date +%Y%m%d-%H%M%S).tar.gz"
	tar -czf $${BACKUP_FILE} bin/ services/control-api/ services/waterfall/ apps/ database/
	@echo "$(GREEN)Backup created: $${BACKUP_FILE}$(NC)"

restore: ## Restaura de backup
	@echo "$(YELLOW)Restoring from backup...$(NC)"
	LATEST_BACKUP=$(ls -t backups/*.tar.gz | head -1)
	if [ -z "$${LATEST_BACKUP}" ]; then
		@echo "$(YELLOW)No backup found$(NC)"
		return 1
	fi
	tar -xzf $${LATEST_BACKUP} -C .
	@echo "$(GREEN)Restored from: $${LATEST_BACKUP}$(NC)"
