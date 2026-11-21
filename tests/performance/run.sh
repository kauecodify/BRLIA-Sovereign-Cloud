#!/usr/bin/env bash
# performance: mede latencia mediana e p95 contra /health localhost via urllib
set -euo pipefail

export PERFORMANCE_PORTS="${PORTS:-8080 8081}"
export PERFORMANCE_SAMPLES="${SAMPLES:-50}"

python3 - <<'RUNPY'
import os
import sys
import time
import urllib.request

ports = os.environ.get("PERFORMANCE_PORTS", "").split()
samples = int(os.environ.get("PERFORMANCE_SAMPLES", "50"))
failed = 0

for port in ports:
    url = "http://127.0.0.1:%s/health" % port
    latencies = []
    for _ in range(samples):
        start = time.monotonic()
        try:
            urllib.request.urlopen(url, timeout=5).read()
        except Exception as exc:
            print("performance: FAIL %s (%s)" % (port, exc), file=sys.stderr)
            failed = 1
            break
        latencies.append((time.monotonic() - start) * 1000)
    if failed:
        break
    latencies.sort()
    p50 = latencies[len(latencies) // 2]
    p95 = latencies[int(len(latencies) * 0.95)]
    print("performance: %s p50=%.1fms p95=%.1fms" % (port, p50, p95))

sys.exit(failed)
RUNPY
