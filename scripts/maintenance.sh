#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "maintenance: limpando artefatos de build antigos"
rm -rf bin

echo "maintenance: pruning de recursos Docker nao utilizados"
docker system prune -f >/dev/null 2>&1 || true

if [ -n "$(docker compose ps -q postgres 2>/dev/null)" ]; then
    echo "maintenance: VACUUM/ANALYZE do PostgreSQL"
    docker compose exec -T postgres psql -U brlia -d brlia_control -c "VACUUM ANALYZE;" || true
fi

echo "maintenance: verificando saude basica dos binarios locais"
if [ -x bin/control-api ]; then
    echo "maintenance: bin/control-api presente"
else
    echo "maintenance: build ausente; rode 'make build-services'" >&2
fi

echo "maintenance: concluido"
