#!/usr/bin/env bats
# tests/ci-guards.bats
# Gaps: CI-DEPLOY-01, CI-GOCOMPAT-01, SIZE-CHECK-01, CI-SAST-01, CI-SUBDOMAIN-01

# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------

# Guard: asserts all three CF deployment secrets are non-empty.
# Mirrors the proposed validate-secrets CI job logic for CI-DEPLOY-01.
_guard_secrets() {
  local missing=""
  for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
    eval "val=\${${var}:-}"
    if [ -z "${val}" ]; then
      missing="${missing} ${var}"
    fi
  done
  if [ -n "${missing}" ]; then
    printf 'FAIL: required variables not set:%s\n' "${missing}" >&2
    return 1
  fi
  echo "OK"
}

# Go version compatibility check.
# Mirrors .gitlab-ci.yml go-compat lines 73-85 with the empty-extraction guard added.
_check_go_compat() {
  local strings_fixture="${1}"
  local required="${2}"

  hugo_go=$(grep -m1 '^go[0-9]' "${strings_fixture}" | sed 's/^go//')

  if [ -z "${hugo_go}" ]; then
    printf 'FAIL: could not extract Go version from binary strings\n' >&2
    return 1
  fi

  printf 'Extracted Go version: %s\n' "${hugo_go}"

  maj_got=$(printf '%s' "${hugo_go}"  | cut -d. -f1)
  min_got=$(printf '%s' "${hugo_go}"  | cut -d. -f2)
  maj_req=$(printf '%s' "${required}" | cut -d. -f1)
  min_req=$(printf '%s' "${required}" | cut -d. -f2)

  fail=0
  [ "${maj_got}" -lt "${maj_req}" ] && fail=1
  [ "${maj_got}" -eq "${maj_req}" ] && [ "${min_got}" -lt "${min_req}" ] && fail=1

  if [ "${fail}" -eq 1 ]; then
    printf 'FAIL: Hugo uses Go %s, below required %s\n' "${hugo_go}" "${required}" >&2
    return 1
  fi
  echo "OK"
}

# Size-checker awk script (fixed with BEGIN { rc=0 }).
_size_check_awk() {
  local input_file="${1}"
  awk -v warn=300 -v fail=400 '
    BEGIN { rc=0 }
    /total$/ { next }
    $1+0 >= fail { print "FAIL: " $2 " -- " $1 " lines"; rc=1; next }
    $1+0 >= warn { print "WARN: " $2 " -- " $1 " lines" }
    END { exit rc+0 }
  ' "${input_file}"
}

# ---------------------------------------------------------------------------
# CI-DEPLOY-01 — required secret variable guard
# ---------------------------------------------------------------------------

@test "guard passes when all three secret vars are set" {
  run env CF_DEPLOY_HOOK="https://hooks.example.com/abc" CF_ZONE_ID="zone1" CF_DNS_TOKEN="tok1" \
    bash -c '
      missing=""
      for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
        eval "val=\${${var}:-}"
        [ -z "${val}" ] && missing="${missing} ${var}"
      done
      [ -z "${missing}" ] && echo "OK" && exit 0
      printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
    '
  [ "${status}" -eq 0 ]
  [[ "${output}" == "OK" ]]
}

@test "guard fails when CF_DEPLOY_HOOK is unset" {
  run env CF_ZONE_ID="zone1" CF_DNS_TOKEN="tok1" \
    bash -c '
      missing=""
      for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
        eval "val=\${${var}:-}"
        [ -z "${val}" ] && missing="${missing} ${var}"
      done
      [ -z "${missing}" ] && echo "OK" && exit 0
      printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
    '
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"CF_DEPLOY_HOOK"* ]]
}

@test "guard fails when CF_ZONE_ID is unset" {
  run env CF_DEPLOY_HOOK="https://hooks.example.com/abc" CF_DNS_TOKEN="tok1" \
    bash -c '
      missing=""
      for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
        eval "val=\${${var}:-}"
        [ -z "${val}" ] && missing="${missing} ${var}"
      done
      [ -z "${missing}" ] && echo "OK" && exit 0
      printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
    '
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"CF_ZONE_ID"* ]]
}

@test "guard fails when CF_DNS_TOKEN is unset" {
  run env CF_DEPLOY_HOOK="https://hooks.example.com/abc" CF_ZONE_ID="zone1" \
    bash -c '
      missing=""
      for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
        eval "val=\${${var}:-}"
        [ -z "${val}" ] && missing="${missing} ${var}"
      done
      [ -z "${missing}" ] && echo "OK" && exit 0
      printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
    '
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"CF_DNS_TOKEN"* ]]
}

@test "guard fails when CF_DEPLOY_HOOK is empty string" {
  run env CF_DEPLOY_HOOK="" CF_ZONE_ID="zone1" CF_DNS_TOKEN="tok1" \
    bash -c '
      missing=""
      for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
        eval "val=\${${var}:-}"
        [ -z "${val}" ] && missing="${missing} ${var}"
      done
      [ -z "${missing}" ] && echo "OK" && exit 0
      printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
    '
  [ "${status}" -ne 0 ]
}

@test "guard fails and reports all three vars when all are unset" {
  run bash -c '
    unset CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN
    missing=""
    for var in CF_DEPLOY_HOOK CF_ZONE_ID CF_DNS_TOKEN; do
      eval "val=\${${var}:-}"
      [ -z "${val}" ] && missing="${missing} ${var}"
    done
    [ -z "${missing}" ] && echo "OK" && exit 0
    printf "FAIL: not set:%s\n" "${missing}" >&2; exit 1
  '
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"CF_DEPLOY_HOOK"* ]]
  [[ "${output}" == *"CF_ZONE_ID"* ]]
  [[ "${output}" == *"CF_DNS_TOKEN"* ]]
}

# ---------------------------------------------------------------------------
# CI-GOCOMPAT-01 — Go version extraction logic
# ---------------------------------------------------------------------------

@test "go-compat: extracts version and passes when hugo go equals required" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'some_symbol\nother_stuff\ngo1.23.4\nmore_stuff\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"OK"* ]]
}

@test "go-compat: passes when hugo go minor is greater than required" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go1.24.0\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "go-compat: passes when hugo go major is greater than required" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go2.0.0\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "go-compat: fails when hugo go minor is below required" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go1.21.9\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"FAIL"* ]]
}

@test "go-compat: fails when hugo go major is below required" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go0.23.0\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -ne 0 ]
}

@test "go-compat: fails with error when no go version line present in strings output" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'some_symbol\nother_stuff\nno_version_here\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"could not extract"* ]]
}

@test "go-compat: extracts only the first go version line (-m1 semantics)" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go1.23.4\ngo1.19.0\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "go-compat: passes on exact boundary 1.23.0" {
  fixture=$(mktemp /tmp/hugo-strings-XXXXXX)
  printf 'go1.23.0\n' > "${fixture}"
  run _check_go_compat "${fixture}" "1.23"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# SIZE-CHECK-01 — awk rc initialization (BEGIN { rc=0 })
# ---------------------------------------------------------------------------

@test "size-checker: exits 0 when all files are under warn threshold" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '100 templates/base/sharp.html\n200 templates/base/swiss.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "size-checker: exits 0 and prints WARN when file hits warn threshold (300)" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '300 templates/base/cinematic.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"WARN"* ]]
}

@test "size-checker: exits 1 when file meets fail threshold (400)" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '400 templates/base/toobig.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"FAIL"* ]]
}

@test "size-checker: exits 1 when file exceeds fail threshold (401 lines)" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '401 templates/pro/toobig.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -ne 0 ]
  [[ "${output}" == *"FAIL"* ]]
}

@test "size-checker: exits 0 and shows no output for empty file list" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [ -z "${output}" ]
}

@test "size-checker: 'total' line from wc -l is skipped (not counted as failure)" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '100 file.html\n7865 total\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [[ "${output}" != *"total"* ]]
}

@test "size-checker: rc is initialized — deterministic exit 0 when all files pass" {
  # Regression for SIZE-CHECK-01: without BEGIN{rc=0}, awk exit behaviour
  # on unset var is implementation-defined. With BEGIN{rc=0} it is deterministic.
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '50 a.html\n60 b.html\n70 c.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
}

@test "size-checker: 399 lines is WARN not FAIL" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '399 templates/base/edge.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [[ "${output}" == *"WARN"* ]]
  [[ "${output}" != *"FAIL"* ]]
}

@test "size-checker: 299 lines produces no output" {
  fixture=$(mktemp /tmp/size-check-XXXXXX)
  printf '299 templates/base/small.html\n' > "${fixture}"
  run _size_check_awk "${fixture}"
  rm -f "${fixture}"
  [ "${status}" -eq 0 ]
  [ -z "${output}" ]
}

# ---------------------------------------------------------------------------
# CI-SAST-01 — sast job must not have allow_failure: true
# Tests catch the current gap; they will pass once the CI is fixed.
# ---------------------------------------------------------------------------

@test "sast CI job does not have allow_failure: true (CI-SAST-01)" {
  ci_file="${BATS_TEST_DIRNAME}/../.gitlab-ci.yml"
  [ -f "${ci_file}" ]

  in_sast=0
  found=0
  while IFS= read -r line; do
    if printf '%s' "${line}" | grep -qE '^sast:'; then
      in_sast=1
      continue
    fi
    if [ "${in_sast}" -eq 1 ] && printf '%s' "${line}" | grep -qE '^[a-zA-Z]'; then
      break
    fi
    if [ "${in_sast}" -eq 1 ] && printf '%s' "${line}" | grep -qE 'allow_failure:\s*true'; then
      found=1
      break
    fi
  done < "${ci_file}"

  if [ "${found}" -eq 1 ]; then
    echo "FAIL: sast job has allow_failure: true — security gate is advisory-only (CI-SAST-01)"
    return 1
  fi
}

@test "subdomain-sync CI job does not have blanket allow_failure: true (CI-SUBDOMAIN-01)" {
  ci_file="${BATS_TEST_DIRNAME}/../.gitlab-ci.yml"
  [ -f "${ci_file}" ]

  in_job=0
  found=0
  while IFS= read -r line; do
    if printf '%s' "${line}" | grep -qE '^subdomain-sync:'; then
      in_job=1
      continue
    fi
    if [ "${in_job}" -eq 1 ] && printf '%s' "${line}" | grep -qE '^[a-zA-Z]'; then
      break
    fi
    if [ "${in_job}" -eq 1 ] && printf '%s' "${line}" | grep -qE 'allow_failure:\s*true'; then
      found=1
      break
    fi
  done < "${ci_file}"

  if [ "${found}" -eq 1 ]; then
    echo "FAIL: subdomain-sync has allow_failure: true — DNS failures are silently swallowed (CI-SUBDOMAIN-01)"
    return 1
  fi
}
