#!/usr/bin/env bash
set -euo pipefail

base_url=${1:-http://127.0.0.1:4010}
summary_file="e2e-summary.md"
failed_tests=0
total_tests=0

# Initialize summary
echo "# E2E Mock Test Results" > "$summary_file"
echo "" >> "$summary_file"

request() {
  local method=$1
  local url=$2
  local data=${3-}
  local test_name="$method $(basename "$url")"
  shift 3

  echo "==> ${method} ${url}"
  if [[ -n "${data}" ]]; then
    echo "    payload: ${data}"
  fi

  local status_code
  local success=true
  
  if ! curl -sS --fail-with-body -o /tmp/e2e-body.json -w "%{http_code}" \
    -X "${method}" \
    "$@" \
    ${data:+-d "$data"} \
    "${url}" > /tmp/e2e-status.txt 2>/dev/null; then
    success=false
  fi
  
  status_code=$(cat /tmp/e2e-status.txt 2>/dev/null || echo "000")
  local body=$(cat /tmp/e2e-body.json 2>/dev/null || echo "No response")
  
  echo "    status: $status_code"
  echo "    body: $body"
  
  ((total_tests++))
  
  # Add to summary
  if [[ "$success" == "true" && "$status_code" =~ ^2[0-9][0-9]$ ]]; then
    echo "✅ $test_name - Status: $status_code" >> "$summary_file"
  else
    echo "❌ $test_name - Status: $status_code" >> "$summary_file"
    echo "::error title=E2E Test Failed::$test_name returned status $status_code"
    ((failed_tests++))
  fi
}

ready=false
for i in {1..30}; do
  if curl -sf "$base_url/prices/singles" >/dev/null 2>&1; then
    ready=true
    break
  fi
  sleep 2
done

if [[ "${ready}" != "true" ]]; then
  echo "::error title=Prism Mock Server Failed::Prism mock did not become ready at ${base_url}"
  if [[ -n "${PRISM_LOG:-}" && -f "${PRISM_LOG}" ]]; then
    echo "Prism log output:"
    cat "${PRISM_LOG}"
  fi
  exit 1
fi

request GET "$base_url/prices/singles"
request GET "$base_url/prices/variants"
request GET "$base_url/prices/sealed"

request GET "$base_url/account" "" \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com"

request POST "$base_url/buyer/optimizer" \
  '{"cart":[{"type":"mtg_single","name":"Polar Kraken","quantity_requested":1,"language_ids":["EN"],"finish_ids":["NF"],"condition_ids":["NM"]}]}' \
  -H "Content-Type: application/json"

request POST "$base_url/deck" \
  '{"commander_names":["Atraxa, Praetors\u0027 Voice"],"other_cards":[{"name":"Lightning Bolt","quantity":4}]}' \
  -H "Content-Type: application/json"

request POST "$base_url/card_info" \
  '{"card_names":["Lightning Bolt"]}' \
  -H "Content-Type: application/json"

request GET "$base_url/buyer/orders?since=2024-04-01T00:00:00Z&limit=1" "" \
  -H "X-ManaPool-Access-Token: test-token" \
  -H "X-ManaPool-Email: test@example.com"

# Generate final summary
echo "" >> "$summary_file"
passed_tests=$((total_tests - failed_tests))

echo "## Summary" >> "$summary_file"
echo "- Total tests: $total_tests" >> "$summary_file"
echo "- Passed: $passed_tests" >> "$summary_file"
echo "- Failed: $failed_tests" >> "$summary_file"

echo "::notice title=E2E Results::$passed_tests/$total_tests tests passed"

if [[ "$failed_tests" -gt 0 ]]; then
  echo "::error title=E2E Tests Failed::$failed_tests out of $total_tests tests failed"
  exit 1
fi

echo "All E2E tests passed!"
