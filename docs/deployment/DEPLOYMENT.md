# BRLIA Sovereign Cloud - Guia de Deploy

## Visao Geral

Guia para levar a pilha do ambiente de desenvolvimento (Docker Compose) ate producao. Os caminhos suportados sao:

1. **Docker Compose** - implantacao de referencia em host unico (com overlay de producao)
2. **Kubernetes** - implantacao em cluster (manifestos em `infra/kubernetes/`)
3. **OpenTofu + Ansible** - provisionamento de nos e configuracao de sistema (`infra/opentofu`, `infra/ansible`)

## Pre-requisitos

- Docker Engine 24+ com plugin compose
- Make, Git, OpenTofu 1.6+ e Ansible 2.15+ (apenas para via IaC)
- DNS do operador apontando para o balanceador de carga (nao coberto nesta base)
- TLS terminado no balanceador ou certificados no Nginx (`infra/nginx`)

## Caminho 1 - Docker Compose (Host Unico)

### 1. Bootstrap

```bash
git clone <repositorio> brlia-cloud && cd brlia-cloud
bash scripts/bootstrap.sh       # valida toolchain, gera .env, build e testes
```

### 2. Configuracao

Edite `.env` e ajuste no minimo:

- `POSTGRES_PASSWORD` - senha do banco de controle
- `SERVER_MODE=release` no overlay de producao (ja definido em `infra/compose/docker-compose.production.yml`)
- Credenciais do `keycloak` simulado ou endpoint real do provedor de identidade

### 3. Banco de Dados

As migracoes em `database/migrations/` sao aplicadas automaticamente via `/docker-entrypoint-initdb.d` na primeira subida do Postgres. Para reaplicar seeds:

```bash
make migrate
make seed
```

### 4. Subida

```bash
docker compose -f docker-compose.yml -f infra/compose/docker-compose.production.yml up -d
docker compose ps
```

Mapa de portas do host:

| Servico | Porta |
| --- | --- |
| Nginx (porta de entrada) | 80 |
| control-api | 8080 |
| waterfall | 8081 |
| keycloak | 8082 |
| identity a government | 8083-8093 |
| Portais (admin/customer/developer/operations) | 3001-3004 |
| Prometheus / Grafana | 9090 / 3000 |

### 5. Validacao Pos-Deploy

```bash
bash tests/e2e/run.sh
```

## Caminho 2 - Kubernetes

1. Criar o namespace: `kubectl apply -f infra/kubernetes/namespace.yaml`
2. Aplicar a carga de control-api: `kubectl apply -f infra/kubernetes/`
3. Validar compatibilidade e capacidade: `kubectl get pods -n brlia`

Limitacoes conhecidas: apenas o `control-api` tem manifesto pronto; os demais servicos seguem o mesmo padrao (Deployment + Service + probes em `/health` e `/ready`), com 3 replicas e limites de recursos conforme `control-api-deployment.yaml`.

## Caminho 3 - OpenTofu + Ansible

```bash
cd infra/opentofu
tofu init
tofu plan  -var "region=br-sul-1"
tofu apply -var "region=br-sul-1"
```

O provisionamento prepara os nos de computacao e alimenta o inventario do Ansible:

```bash
cd ../ansible
ansible-playbook -i inventory/hosts.yml playbook.yml
```

O playbook instala e configura containerd e chrony (sincronizacao NTP) nos nos do cluster soberano. Acoes de rede (nftables, VLANs, segmentacao) ficam em `infra/networking/SOVEREIGN_NETWORK.md`.

## Backup e Continuidade

```bash
bash scripts/backup.sh                                                    # dump Postgres + snapshot Redis em backups/
RESTORE_CONFIRM=restaurar bash scripts/restore.sh backups/brlia-<timestamp>.tar.gz
```

Agendar `backup.sh` via cron (recomendado a cada 15 minutos). Testar a restauracao mensalmente.

## Atualizacoes e Rollback

1. Estrategia: implantacao azul-verde por servico, um servico por vez.
2. Ordem recomendada: `keycloak`/`identity` -> `control-api` -> `waterfall` -> demais servicos -> portais
3. Antes de iniciar: `bash scripts/backup.sh` (ponto de retorno garantido)
4. Rollback: `docker compose stop <servico>` e redeploy da tag anterior da imagem

## Hardening de Producao

- [ ] `SERVER_MODE=release` em todos os servicos
- [ ] Credenciais padrao removidas do `.env`
- [ ] TLS no Nginx e no balanceador
- [ ] `bash tests/security/run.sh` completo antes de ir para producao
- [ ] Limites de recursos (CPU/RAM) definidos nos containers
- [ ] Monitoramento com alertas configurados (Prometheus/Grafana ja inclusos)
