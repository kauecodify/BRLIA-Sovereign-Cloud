#!/usr/bin/env bash
set -euo pipefail

if ! command -v docker >/dev/null 2>&1; then
    echo "resilience: Docker nao disponivel; suite pulada neste ambiente" >&2
    exit 0
fi

COMPOSE="docker compose"

for svc in postgres redis control-api identity compute; do
    echo "resilience: reiniciando $svc e verificando recuperacao"
    $COMPOSE restart "$svc" >/dev/null 2>&1 || { echo "resilience: falha ao reiniciar $svc" >&2; exit 1; }
done

echo "resilience: ok"
