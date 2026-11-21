# Suites de Teste BRLIA

| Suite        | Execucao                        | Requisitos                       |
|--------------|---------------------------------|----------------------------------|
| unit         | `make test-unit` ou `bash tests/unit/run.sh`        | Go 1.21+      |
| integration  | `bash tests/integration/run.sh` | Docker Compose                   |
| e2e          | `bash tests/e2e/run.sh`         | Stack do compose no ar           |
| security     | `bash tests/security/run.sh`    | govulncheck e gitleaks opcionais |
| performance  | `bash tests/performance/run.sh` | Stack no ar; `hey` opcional      |
| resilience   | `bash tests/resilience/run.sh`  | Docker Compose                   |

O CI executa unit, security (estatico) e validacao de contratos por padrao.
Suites integration/e2e/performance/resilience sao gates de pipeline completo
e exigem a stack disponivel.
