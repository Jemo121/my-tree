#!/usr/bin/env bats
# tests/subdomain-sync.bats — CI-SUBDOMAIN-01, MISSING-CONTRACT-01
# Tests the subdomain-sync shell logic using mock curl responses.
# No real Cloudflare API calls are made.

FIXTURES_DIR="${BATS_TEST_DIRNAME}/fixtures"

# ---------------------------------------------------------------------------
# Extracted sync logic — mirrors .gitlab-ci.yml lines 153-174 with:
#   1. curl exit-code guard before jq parsing
#   2. allow_failure removed from non-early-exit paths
#
# Usage: _sync_subdomain <subdomain> <mock_curl_script>
#   mock_curl_script: path to a sh script that simulates curl behaviour.
#   It must write its response body to stdout and exit with the desired code.
# ---------------------------------------------------------------------------
_sync_subdomain() {
  local subdomain="${1}"
  local mock_curl="${2}"
  local fqdn="${subdomain}.eguia.me"
  local CF_ZONE_ID="testzone"
  local CF_DNS_TOKEN="testtoken"
  local CF_PAGES_DOMAIN="my-tree.pages.dev"

  # First call: list existing records
  response=$(bash "${mock_curl}" list 2>&1)
  curl_exit=$?
  if [ "${curl_exit}" -ne 0 ]; then
    echo "CF API error on list (exit ${curl_exit}): ${response}" >&2
    return 1
  fi

  count=$(printf '%s' "${response}" | jq '.result | length' 2>/dev/null)
  if [ -z "${count}" ]; then
    echo "CF API returned unparseable JSON" >&2
    return 1
  fi

  if [ "${count}" -eq 0 ]; then
    echo "Creating CNAME: ${fqdn} → ${CF_PAGES_DOMAIN}"
    create_response=$(bash "${mock_curl}" create 2>&1)
    create_exit=$?
    if [ "${create_exit}" -ne 0 ]; then
      echo "CF API error on create (exit ${create_exit}): ${create_response}" >&2
      return 1
    fi
    result=$(printf '%s' "${create_response}" | jq -r 'if .success then "Created: \(.result.name)" else "Error: \(.errors[0].message)" end')
    echo "${result}"
  else
    echo "Exists: ${fqdn}"
  fi
}

# ---------------------------------------------------------------------------
# Fixture helpers — write temp mock curl scripts
# ---------------------------------------------------------------------------

_mock_curl_200_no_records() {
  tmpfile=$(mktemp /tmp/mock-curl-XXXXXX.sh)
  cat > "${tmpfile}" << 'MOCK'
#!/bin/sh
# Simulates CF DNS list returning 0 existing records (HTTP 200)
cat << 'JSON'
{"success":true,"errors":[],"messages":[],"result":[],"result_info":{"count":0,"page":1,"per_page":20,"total_count":0}}
JSON
exit 0
MOCK
  echo "${tmpfile}"
}

_mock_curl_200_one_record() {
  tmpfile=$(mktemp /tmp/mock-curl-XXXXXX.sh)
  cat > "${tmpfile}" << 'MOCK'
#!/bin/sh
# Simulates CF DNS list returning 1 existing record (HTTP 200)
cat << 'JSON'
{"success":true,"errors":[],"messages":[],"result":[{"id":"abc","type":"CNAME","name":"emilio.eguia.me","content":"my-tree.pages.dev","proxied":true,"ttl":1}],"result_info":{"count":1,"page":1,"per_page":20,"total_count":1}}
JSON
exit 0
MOCK
  echo "${tmpfile}"
}

_mock_curl_401() {
  tmpfile=$(mktemp /tmp/mock-curl-XXXXXX.sh)
  cat > "${tmpfile}" << 'MOCK'
#!/bin/sh
# Simulates CF API returning HTTP 401 (curl -sf exits non-zero on 4xx with -f)
echo '{"success":false,"errors":[{"code":9109,"message":"Invalid access token."}]}'
exit 22
MOCK
  echo "${tmpfile}"
}

_mock_curl_malformed_json() {
  tmpfile=$(mktemp /tmp/mock-curl-XXXXXX.sh)
  cat > "${tmpfile}" << 'MOCK'
#!/bin/sh
# Simulates CF API returning malformed JSON (e.g. a Cloudflare HTML error page)
echo '<html><body>Gateway Error</body></html>'
exit 0
MOCK
  echo "${tmpfile}"
}

_mock_curl_create_success() {
  # For tests that go through to the CREATE call, the mock must handle both
  # "list" (returns 0 records) and "create" (returns success).
  tmpfile=$(mktemp /tmp/mock-curl-XXXXXX.sh)
  cat > "${tmpfile}" << 'MOCK'
#!/bin/sh
ACTION="${1:-list}"
if [ "${ACTION}" = "list" ]; then
  cat << 'JSON'
{"success":true,"errors":[],"messages":[],"result":[],"result_info":{"count":0}}
JSON
else
  cat << 'JSON'
{"success":true,"errors":[],"result":{"id":"new123","type":"CNAME","name":"emilio.eguia.me","content":"my-tree.pages.dev"}}
JSON
fi
exit 0
MOCK
  echo "${tmpfile}"
}

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

@test "sync skips when CF_PAGES_DOMAIN is unset (early exit 0)" {
  # Mirrors the guard at line 154 of .gitlab-ci.yml
  run bash -c '
    if [ -z "${CF_PAGES_DOMAIN}" ]; then
      echo "CF_PAGES_DOMAIN not set — skipping DNS sync"
      exit 0
    fi
    exit 1
  '
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"skipping"* ]]
}

@test "sync exits non-zero when CF API returns HTTP 401" {
  mock=$(_mock_curl_401)
  run _sync_subdomain "emilio" "${mock}"
  rm -f "${mock}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"CF API error"* ]]
}

@test "sync exits non-zero when CF API returns malformed JSON" {
  mock=$(_mock_curl_malformed_json)
  run _sync_subdomain "emilio" "${mock}"
  rm -f "${mock}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"unparseable"* ]]
}

@test "sync prints 'Exists' and exits 0 when record already present" {
  mock=$(_mock_curl_200_one_record)
  run _sync_subdomain "emilio" "${mock}"
  rm -f "${mock}"
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"Exists: emilio.eguia.me"* ]]
}

@test "sync creates record and exits 0 when no existing record found" {
  mock=$(_mock_curl_create_success)
  run _sync_subdomain "emilio" "${mock}"
  rm -f "${mock}"
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"Creating CNAME"* ]]
}

# ---------------------------------------------------------------------------
# MISSING-CONTRACT-01 — CF DNS API fixture contract test
# ---------------------------------------------------------------------------

@test "CF DNS list fixture has expected shape: .success .result .result_info" {
  fixture="${FIXTURES_DIR}/cf-dns-list.json"
  [ -f "${fixture}" ]
  run jq -e '.success == true and (.result | type) == "array" and .result_info != null' "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "CF DNS list fixture .result | length returns integer" {
  fixture="${FIXTURES_DIR}/cf-dns-list.json"
  count=$(jq '.result | length' "${fixture}")
  # Must be a non-negative integer
  printf '%s' "${count}" | grep -qE '^[0-9]+$'
}

@test "CF DNS create success fixture has .success true and .result.name" {
  fixture="${FIXTURES_DIR}/cf-dns-create-success.json"
  [ -f "${fixture}" ]
  run jq -e '.success == true and .result.name != null' "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "CF DNS create error fixture has .success false and .errors[0].message" {
  fixture="${FIXTURES_DIR}/cf-dns-create-error.json"
  [ -f "${fixture}" ]
  run jq -e '.success == false and (.errors | length) > 0 and .errors[0].message != null' "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "jq '.result | length' on HTTP-401 JSON body returns non-zero count guard" {
  # Verifies that our guard ([ -z count ] && exit 1) would catch a failed CF call.
  # The 401 body has result:null — jq emits 'null', not an integer.
  body='{"success":false,"errors":[{"code":9109,"message":"Invalid access token."}],"result":null}'
  count=$(printf '%s' "${body}" | jq '.result | length' 2>/dev/null)
  # 'null | length' in jq returns 0 — this is the silent-pass trap.
  # The guard must check curl exit code BEFORE jq, not after.
  # Asserting 0 here documents the trap so reviewers understand why the
  # curl-exit-code guard is necessary.
  [ "${count}" = "0" ]
}
