#!/usr/bin/env bash
# Auto-generated test harness — run after cloning to verify the build.
set -uo pipefail
FAIL=0

echo "-- syntax/build check --"
if go vet ./... && go build ./...; then echo "  ok"; else echo "  FAIL"; FAIL=1; fi

echo "-- functional check --"
OUT="$(go run main.go access.log 2>&1)"
if echo "$OUT" | grep -qF -- 'Total requests:    4'; then
  echo "  ok   - contains: Total requests:    4"
else
  echo "  FAIL - expected to contain: Total requests:    4"
  echo "  got: $OUT"
  FAIL=1
fi

if [ "$FAIL" = "1" ]; then echo; echo "FAILED"; exit 1; fi
echo; echo "All tests passed."
