# Scripts de Automacao BRLIA

| Script                | Funcao                                                        |
|-----------------------|---------------------------------------------------------------|
| `bootstrap.sh`        | Preparar um ambiente local do zero (build + lint + testes)    |
| `backup.sh`           | Backup de volumes Postgres/Redis e artefatos criticos          |
| `restore.sh`          | Restaurar o backup mais recente                               |
| `disaster-recovery.sh`| Plano de DR executavel (verificacao e failover basico)         |
| `maintenance.sh`      | Manutencao periodica (prune, rotacao e verificacao de saude)   |

Todos os scripts sao idempotentes, aceitam variaveis de ambiente e sao
seguros para cron/systemd timers.
