#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$ROOT/backups}"

LATEST="$(ls -1t "$BACKUP_DIR"/brlia-*.tar.gz 2>/dev/null | head -1 || true)"
if [ -z "$LATEST" ]; then
    echo "restore: nenhum backup encontrado em $BACKUP_DIR" >&2
    exit 1
fi

if [ -z "${RESTORE_CONFIRM:-}" ]; then
    echo "restore: este comando sobrescreve o estado atual." >&2
    echo "restore: defina RESTORE_CONFIRM=restaurar para prosseguir." >&2
    exit 1
fi

[ "$RESTORE_CONFIRM" = "restaurar" ] || { echo "restore: confirmacao invalida" >&2; exit 1; }

WORK="$(mktemp -d)"
tar -xzf "$LATEST" -C "$WORK"

if [ -f "$WORK"/*/postgres.sql ]; then
    echo "restore: reimportando PostgreSQL de $(basename "$LATEST")"
    SQL="$(ls "$WORK"/*/postgres.sql)"
    docker compose -f "$ROOT/docker-compose.yml" exec -T postgres psql -U brlia -d brlia_control < "$SQL"
else
    echo "restore: sem dump de PostgreSQL no pacote"
fi

cp -r "$WORK"/database/* "$ROOT/database/" 2>/dev/null || true
rm -rf "$WORK"

echo "restore: concluido a partir de $LATEST"
