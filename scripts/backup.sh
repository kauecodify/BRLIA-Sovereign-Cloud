#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$ROOT/backups}"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
TARGET="$BACKUP_DIR/brlia-$TIMESTAMP"

mkdir -p "$TARGET"

dump_postgres() {
    if docker compose -f "$ROOT/docker-compose.yml" ps --quiet postgres >/dev/null 2>&1 && \
       [ -n "$(docker compose -f "$ROOT/docker-compose.yml" ps -q postgres 2>/dev/null)" ]; then
        echo "backup: dump PostgreSQL"
        docker compose -f "$ROOT/docker-compose.yml" exec -T pg_dump pg_dump -U brlia brlia_control > "$TARGET/postgres.sql" 2>/dev/null \
            || docker compose -f "$ROOT/docker-compose.yml" exec -T postgres pg_dump -U brlia brlia_control > "$TARGET/postgres.sql"
    else
        echo "backup: PostgreSQL nao esta em execucao; dump ignorado"
    fi
}

dump_redis() {
    if [ -n "$(docker compose -f "$ROOT/docker-compose.yml" ps -q redis 2>/dev/null)" ]; then
        echo "backup: snapshot Redis"
        docker compose -f "$ROOT/docker-compose.yml" exec -T redis redis-cli BGSAVE >/dev/null || true
    fi
}

dump_postgres
dump_redis

echo "backup: copiando contratos, migracoes e seeds"
cp -r "$ROOT/contracts" "$TARGET/contracts"
cp -r "$ROOT/database" "$TARGET/database"

tar -czf "$TARGET.tar.gz" -C "$BACKUP_DIR" "brlia-$TIMESTAMP"
rm -rf "$TARGET"

echo "backup: concluido em $TARGET.tar.gz"
