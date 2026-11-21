#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "bootstrap: validando toolchain"
command -v go >/dev/null 2>&1 || { echo "Go 1.21+ e obrigatorio" >&2; exit 1; }

if [ ! -f .env ] && [ -f .env.example ]; then
    cp .env.example .env
    echo "bootstrap: .env criado a partir de .env.example (revise os valores)"
fi

echo "bootstrap: build dos servicos"
make build-services

echo "bootstrap: build dos adaptadores"
make build-adapters

echo "bootstrap: testes unitarios e lint"
make test-unit
make lint

echo "bootstrap: concluido; para subir a stack: make docker-up"
