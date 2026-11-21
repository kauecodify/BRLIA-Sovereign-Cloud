#!/usr/bin/env bash
# e2e: verifica /health de todos os servicos. HTTP feito via python urllib,
# restringido a localhost.
set -euo pipefail

export E2E_PORTS="${E2E_PORTS:-8080 8081 8083 8084 8085 8086 8087 8088 8089 8090 8091 8092 8093}"

python3 - <<'RUNPY'
import json
import os
import sys
import urllib.request

ports = os.environ.get("E2E_PORTS", "").split()
base = os.environ.get("E2E_HOST", "127.0.0.1")
failed = 0

for port in ports:
    url = "http://%s:%s/health" % (base, port)
    try:
        with urllib.request.urlopen(url, timeout=5) as resp:
            body = json.loads(resp.read().decode())
            print("e2e: ok %s (%s)" % (port, body.get("service", "?")))
    except Exception as exc:
        print("e2e: FAIL %s (%s)" % (port, exc), file=sys.stderr)
        failed = 1

sys.exit(failed)
RUNPY
