# BRLIA Sovereign Cloud - Makefile
# Comandos utilitarios para desenvolvimento e deploy

.PHONY: help build build-services build-adapters run stop clean test test-unit test-integration test-e2e lint docker docker-up docker-down docker-build docker-push migrate seed backup restore

# Cores para output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

# Variaveis de ambiente
ENV ?= development
BIN_DIR ?= ./bin
DOCKER_REGISTRY ?= ghcr.io/kauecodify
DOCKER_TAG ?= latest

# Componentes
SERVICES := control-api waterfall identity compute storage network security billing audit ai-platform hpc quantum-broker government
ADAPTERS := kubernetes kvm ceph slurm hsm datacenter
DOCKER_SERVICES := identity compute storage network security billing audit ai-platform hpc quantum-broker government

help: ## Mostra esta ajuda
	@echo "BRLIA Sovereign Cloud - Comandos Disponiveis"
	@echo ""
	@echo " Comandos de Build:"
	@echo "   make build            - Builda todos os servicos"
	@echo "   make build-services   - Builda apenas os servicos Go"
	@echo "   make build-adapters   - Builda os adaptadores de infraestrutura"
	@echo ""
	@echo " Comandos de Execucao:"
	@echo "   make run              - Executa os servicos compilados localmente"
	@echo ""
	@echo " Comandos Docker:"
	@echo "   make docker-up        - Sobe todos os containers"
	@echo "   make docker-down      - Para todos os containers"
	@echo "   make docker-build     - Builda imagens Docker"
	@echo "   make docker-push      - Publica imagens Docker"
	@echo ""
	@echo " Comandos de Database:"
	@echo "   make migrate          - Aplica migracoes do banco"
	@echo "   make seed             - Popula banco com dados iniciais"
	@echo ""
	@echo " Comandos de Qualidade:"
	@echo "   make test             - Testes unitarios + lint"
	@echo "   make test-integration - Testes de integracao (requer Docker)"
	@echo "   make test-e2e         - Testes ponta a ponta (requer stack no ar)"
	@echo "   make lint             - Executa go vet em todos os modulos"
	@echo ""
	@echo " Comandos de Manutencao:"
	@echo "   make clean            - Remove arquivos temporarios"
	@echo "   make backup           - Cria backup do sistema"
	@echo "   make restore          - Restaura de backup"
	@echo "   make stop             - Para todos os servicos"

# Build
build: build-services build-adapters ## Builda todos os componentes

build-services: ## Builda os servicos Go
	@mkdir -p $(BIN_DIR)
	@for svc in $(SERVICES); do \
		echo "$(GREEN)Building $$svc...$(NC)"; \
		(cd services/$$svc && go build -o ../../$(BIN_DIR)/$$svc main.go) || exit 1; \
	done
	@echo "$(GREEN)All services built$(NC)"

build-adapters: ## Builda os adaptadores de infraestrutura
	@mkdir -p $(BIN_DIR)/adapters
	@for adp in $(ADAPTERS); do \
		echo "$(GREEN)Building adapter $$adp...$(NC)"; \
		(cd adapters/$$adp && go build -o ../../$(BIN_DIR)/adapters/$$adp main.go) || exit 1; \
	done
	@echo "$(GREEN)All adapters built$(NC)"

# Run
run: build-services ## Executa os servicos compilados localmente
	@for svc in $(SERVICES); do \
		echo "$(GREEN)Starting $$svc...$(NC)"; \
		$(BIN_DIR)/$$svc & \
	done
	@wait

# Docker
docker-up: ## Sobe todos os containers
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker compose up -d
	@echo "$(GREEN)Containers started. Access:"
	@echo "  - Control API:       http://localhost:8080"
	@echo "  - Waterfall:         http://localhost:8081"
	@echo "  - Keycloak:          http://localhost:8082"
	@echo "  - Identity:          http://localhost:8083"
	@echo "  - Compute:           http://localhost:8084"
	@echo "  - Storage:           http://localhost:8085"
	@echo "  - Network:           http://localhost:8086"
	@echo "  - Security:          http://localhost:8087"
	@echo "  - Billing:           http://localhost:8088"
	@echo "  - Audit:             http://localhost:8089"
	@echo "  - AI Platform:       http://localhost:8090"
	@echo "  - HPC:               http://localhost:8091"
	@echo "  - Quantum Broker:    http://localhost:8092"
	@echo "  - Government:        http://localhost:8093"
	@echo "  - Admin Console:     http://localhost:3001"
	@echo "  - Customer Portal:   http://localhost:3002"
	@echo "  - Developer Portal:  http://localhost:3003"
	@echo "  - Operations Center: http://localhost:3004"
	@echo "  - Prometheus:        http://localhost:9090"
	@echo "  - Grafana:           http://localhost:3000"
	@echo "  - Vault:             http://localhost:8200$(NC)"

docker-down: ## Para todos os containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker compose down

docker-build: ## Builda imagens Docker
	@echo "$(GREEN)Building Docker images...$(NC)"
	docker compose build

docker-push: ## Publica imagens Docker
	@echo "$(GREEN)Pushing Docker images to $(DOCKER_REGISTRY)...$(NC)"
	docker tag brlia-control-api:latest $(DOCKER_REGISTRY)/brlia-control-api:$(DOCKER_TAG) && docker push $(DOCKER_REGISTRY)/brlia-control-api:$(DOCKER_TAG)
	docker tag brlia-waterfall:latest $(DOCKER_REGISTRY)/brlia-waterfall:$(DOCKER_TAG) && docker push $(DOCKER_REGISTRY)/brlia-waterfall:$(DOCKER_TAG)
	@for svc in $(DOCKER_SERVICES); do \
		docker tag brlia-$$svc:latest $(DOCKER_REGISTRY)/brlia-$$svc:$(DOCKER_TAG) && \
		docker push $(DOCKER_REGISTRY)/brlia-$$svc:$(DOCKER_TAG) || exit 1; \
	done
	@echo "$(GREEN)Images pushed successfully$(NC)"

# Database
migrate: ## Aplica migracoes do banco
	@echo "$(GREEN)Running database migrations...$(NC)"
	@for migration in database/migrations/*.sql; do \
		echo "Applying $$migration"; \
		docker compose exec -T postgres psql -U brlia -d brlia_control < "$$migration" || exit 1; \
	done
	@echo "$(GREEN)Migrations applied$(NC)"

seed: ## Popula banco com dados iniciais
	@echo "$(GREEN)Seeding database...$(NC)"
	@for seed in database/seeds/*.sql; do \
		echo "Applying $$seed"; \
		docker compose exec -T postgres psql -U brlia -d brlia_control < "$$seed" || exit 1; \
	done
	@echo "$(GREEN)Database seeded$(NC)"

# Testes
test: test-unit lint ## Executa quality gates locais

test-unit: ## Executa testes unitarios
	@echo "$(GREEN)Running unit tests...$(NC)"
	@for svc in $(SERVICES); do \
		echo "$(GREEN)Testing $$svc...$(NC)"; \
		(cd services/$$svc && go test ./... -count=1) || exit 1; \
	done
	@echo "$(GREEN)Unit tests passed$(NC)"

test-integration: ## Executa testes de integracao
	@echo "$(GREEN)Running integration tests...$(NC)"
	bash tests/integration/run.sh || { echo "$(RED)Integration tests failed$(NC)"; exit 1; }

test-e2e: ## Executa testes ponta a ponta
	@echo "$(GREEN)Running e2e tests...$(NC)"
	bash tests/e2e/run.sh || { echo "$(RED)E2E tests failed$(NC)"; exit 1; }

# Manutencao
clean: ## Remove arquivos temporarios
	@echo "$(YELLOW)Cleaning up...$(NC)"
	rm -rf $(BIN_DIR)
	docker system prune -f
	@echo "$(GREEN)Cleanup complete$(NC)"

stop: ## Para todos os servicos
	@echo "$(YELLOW)Stopping all services...$(NC)"
	pkill -f "bin/control-api|bin/waterfall|bin/identity|bin/compute|bin/storage|bin/network|bin/security|bin/billing|bin/audit|bin/ai-platform|bin/hpc|bin/quantum-broker|bin/government" || true
	docker compose down || true
	@echo "$(GREEN)All services stopped$(NC)"

lint: ## Executa go vet em todos os modulos
	@echo "$(GREEN)Running go vet...$(NC)"
	@for svc in $(SERVICES); do \
		(cd services/$$svc && go vet ./...) || exit 1; \
	done
	@for adp in $(ADAPTERS); do \
		(cd adapters/$$adp && go vet ./...) || exit 1; \
	done
	@echo "$(GREEN)Linting complete$(NC)"

backup: ## Cria backup do sistema
	@bash scripts/backup.sh

restore: ## Restaura de backup
	@bash scripts/restore.sh
