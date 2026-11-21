# BRLIA Sovereign Cloud - Runbook de Operacoes

## SLOs e Indicadores

| Indicador | Objetivo | Medicao |
| --- | --- | --- |
| Disponibilidade do Control API | 99,9% mensal | Prometheus: probes em `/health` |
| Disponibilidade do Waterfall | 99,5% mensal | Prometheus: probes em `/health` |
| Latencia p95 de API | < 200 ms | `tests/performance/run.sh` (p95) |
| Taxa de erro 5xx | < 0,5% | Logs agregados e metricas |
| RTO (desastre) | < 4 horas | `scripts/disaster-recovery.sh` (ensaio trimestral) |
| RPO (dados) | < 15 minutos | Frequencia de backup (`scripts/backup.sh` via cron) |

## Estado de Saude por Servico

Cada microsservico expoe:

- `GET /health` - probe de vida (sem dependencias)
- `GET /ready` - probe de prontidao (PostgreSQL e Redis obrigados)

O Nginx expoe `GET /api/health` agregando o `control-api` (ver `infra/nginx/nginx.conf`).

### Comandos Rapidos

```bash
# Estado geral da pilha
make docker-up          # iniciar todos os servicos e imprimir mapa de portas
docker compose ps       # estado individual dos containers

# Saude por servico (verifica o codigo HTTP de cada porta)
for port in 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089 8090 8091 8092 8093; do
  python3 -c "import urllib.request; print('${port}', urllib.request.urlopen(f'http://localhost:{${port}}/health').status)"
done
```

## Runbook: Servico Indisponivel

**Sintoma:** `docker compose ps` mostra `unhealthy` ou probe de vida falha.

1. `docker compose logs --tail 200 <servico>` - identificar erro de inicializacao
2. Verificar dependencias: Postgres (`5432`), Redis (`6379`), NATS (`4222`)
3. Restart controlado: `docker compose restart <servico>`
4. Se persistir: `docker compose stop <servico> && docker compose up -d <servico>`
5. Validar localmente com `tests/resilience/run.sh`

**Escalacao:** se mais de 2 servicos estiverem instaveis, tratar como incidente de plataforma (S2).

## Runbook: Banco de Dados Cheio

**Sintoma:** alerta de disco no no Postgres ou erros `disk full` nos logs.

1. `docker compose exec postgres psql -U brlia -d brlia_control -c "SELECT pg_database_size('brlia_control');"`
2. Executar manutencao: `bash scripts/maintenance.sh` (inclui `VACUUM ANALYZE`)
3. Corrigir crescimento anormal: a tabela `audit_logs` e o maior crescimento esperado; revisar politica de arquivamento
4. Ampliar volume do host se a causa for legitima

## Runbook: Perda de Dados

**Sintoma:** corrompimento logico ou descarte acidental de tabelas do plano de controle.

1. Congelar escritas: `docker compose stop control-api`
2. Localizar snapshot mais recente: `ls -t backups/brlia-*.tar.gz | head -1`
3. `RESTORE_CONFIRM=restaurar bash scripts/restore.sh backups/brlia-<timestamp>.tar.gz`
4. Validar integridade: `docker compose exec postgres psql -U brlia -d brlia_control -c "\dt"`
5. Reabrir escritas e registrar o incidente no diario de auditoria

**Escalacao:** use `bash scripts/disaster-recovery.sh` apenas quando o objetivo for reiniciar toda a pilha a partir de um backup.

## Runbook: Incidente de Seguranca

**Sintoma:** alertas do gitleaks no CI, padroes anormais no `services/audit` ou tentativas de invasao.

1. Seguir o fluxo de resposta a incidentes em `docs/security/SECURITY.md` (secao Resposta a Incidentes)
2. Revogar credenciais afetadas via `services/security` ou Vault
3. Bloquear a origem no Nginx/WAF
4. Rotacionar todos os segredos do `.env` e do Vault
5. Documentar no relatorio pos-incidente

## Rotina Diaria

- Verificar `docker compose ps` - todos `healthy`
- Conferir alertas no Grafana (porta 3000)
- Checar fila de eventos NATS (atrasos em assuntos `brlia.*`)

## Rotina Semanal

- `bash tests/integration/run.sh` contra ambiente de homologacao
- `bash tests/security/run.sh` completo (com govulncheck se disponivel)
- Revisar espaco em disco e tamanho de `backups/`

## Rotina Mensal

- Ensaio de desastre: `scripts/disaster-recovery.sh` em ambiente de teste
- Revisao dos SLOs apurados no periodo
- Rotacao de credenciais de longa duracao (banco, chaves de unseal do Vault)

## Contatos de Escalacao

| Papel | Responsabilidade | Acionamento |
| --- | --- | --- |
| Operador de plantao | Primeira resposta, triagem S1/S2 | Sempre |
| Engenharia de plataforma | Causa raiz em servicos Go | Incidentes S2+ |
| Responsavel de seguranca | Incidentes com dados ou invasao | Sempre que envolver dados |
