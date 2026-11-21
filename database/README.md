# BRLIA Sovereign Cloud - Database

- `migrations/`: schema do PostgreSQL. Aplicado automaticamente no primeiro
  boot do `docker compose` (montado em `/docker-entrypoint-initdb.d`) ou via
  `make migrate`.
- `seeds/`: dados iniciais (regioes, roles e organizacao de demonstracao),
  aplicados via `make seed`.

Todos os scripts sao idempotentes e podem ser reexecutados com seguranca.
