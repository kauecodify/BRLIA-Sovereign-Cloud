#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

if ! command -v docker >/dev/null 2>&1; then
    echo "integration: Docker nao disponivel; suite pulada neste ambiente" >&2
    exit 0
fi

echo "integration: subindo dependencias (postgres, redis, nats)"
docker compose -f "$ROOT/docker-compose.yml" up -d postgres redis nats

cleanup() {
    docker compose -f "$ROOT/docker-compose.yml" stop postgres redis nats >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "aguardando postgres..."
for i in $(seq 1 30); do
    if docker compose -f "$ROOT/docker-compose.yml" exec -T postgres pg_isready -U brlia >/dev/null 2>&1; then
        break
    fi
    sleep 1
done

echo "integration: validacao de schema"
tables="$(docker compose -f "$ROOT/docker-compose.yml" exec -T postgres psql -U brlia -d brlia_control -c '\\dt')"
echo "$tables"
case "$tables" in
    *organizations*) echo "integration: ok" ;;
    *) echo "integration: FAIL tabela organizations ausente" >&2; exit 1 ;;
esac
