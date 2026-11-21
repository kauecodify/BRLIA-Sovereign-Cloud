#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FAILED=0

for svc in "$ROOT"/services/*/; do
    name="$(basename "$svc")"
    echo "--- unit: $name"
    (cd "$svc" && go test ./... -count=1) || FAILED=1
done

for adp in "$ROOT"/adapters/*/; do
    name="$(basename "$adp")"
    echo "--- unit: adapter $name"
    (cd "$adp" && go test ./... -count=1) || FAILED=1
done

exit $FAILED
