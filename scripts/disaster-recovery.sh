#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "DR: 1/4 verificando integridade do ultimo backup"
bash "$(dirname "$0")/backup.sh"

echo "DR: 2/4 derrubando a stack corrente"
docker compose down || true

echo "DR: 3/4 reconstruindo imagens"
docker compose build

echo "DR: 4/4 restabelecendo a stack"
docker compose up -d

echo "DR: failover de emergencia concluido; valide /health de cada servico"
echo "DR: proximos passos manuais: replicacao inter-regiao e validacao de dados"
