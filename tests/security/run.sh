#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FAILED=0

echo "security: go vet em todos os modulos"
for m in "$ROOT"/services/*/ "$ROOT"/adapters/*/; do
    (cd "$m" && go vet ./...) || FAILED=1
done

if command -v govulncheck >/dev/null 2>&1; then
    echo "security: govulncheck"
    for m in "$ROOT"/services/*/ "$ROOT"/adapters/*/; do
        (cd "$m" && govulncheck ./...) || FAILED=1
    done
else
    echo "security: govulncheck ausente; execucao estatica limitada" >&2
fi

if command -v gitleaks >/dev/null 2>&1; then
    echo "security: gitleaks"
    (cd "$ROOT" && gitleaks detect --no-git --redact) || FAILED=1
fi

exit $FAILED
