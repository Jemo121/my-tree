#!/usr/bin/env bats
# tests/validate-html.bats — HTML-VAL-01
# Gates HTML validity locally using the same html-validate@9 invocation
# as the CI job, so a malformed template is caught before push.
# Run: bats tests/validate-html.bats
#      (requires Node.js / npx on PATH)

REPO_ROOT="${BATS_TEST_DIRNAME}/.."

setup() {
  if ! command -v npx >/dev/null 2>&1; then
    skip "npx not found — install Node.js to run HTML validation tests"
  fi
  if ! command -v node >/dev/null 2>&1; then
    skip "node not found — install Node.js to run HTML validation tests"
  fi
}

# ---------------------------------------------------------------------------
# 1. html-validate can be resolved by npx (network or cache)
# ---------------------------------------------------------------------------
@test "html-validate@9 is resolvable via npx" {
  run npx --yes html-validate@9 --version
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 2. All base templates pass html-validate
# ---------------------------------------------------------------------------
@test "templates/base/*.html all pass html-validate@9" {
  cd "${REPO_ROOT}"
  run npx --yes html-validate@9 "templates/base/*.html"
  if [ "${status}" -ne 0 ]; then
    echo "${output}"
  fi
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 3. All meme templates pass html-validate
# ---------------------------------------------------------------------------
@test "templates/meme/*.html all pass html-validate@9" {
  cd "${REPO_ROOT}"
  run npx --yes html-validate@9 "templates/meme/*.html"
  if [ "${status}" -ne 0 ]; then
    echo "${output}"
  fi
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 4. All pro templates pass html-validate
# ---------------------------------------------------------------------------
@test "templates/pro/*.html all pass html-validate@9" {
  cd "${REPO_ROOT}"
  run npx --yes html-validate@9 "templates/pro/*.html"
  if [ "${status}" -ne 0 ]; then
    echo "${output}"
  fi
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 5. Error pages pass html-validate
# ---------------------------------------------------------------------------
@test "errors/*.html all pass html-validate@9" {
  cd "${REPO_ROOT}"
  run npx --yes html-validate@9 "errors/*.html"
  if [ "${status}" -ne 0 ]; then
    echo "${output}"
  fi
  [ "${status}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 6. Intentionally malformed HTML fragment is rejected (self-test of tool)
# ---------------------------------------------------------------------------
@test "html-validate@9 correctly rejects malformed HTML" {
  tmpfile=$(mktemp /tmp/invalid-XXXXXX.html)
  # Missing closing tag for <p> inside a <div> — html-validate catches this
  cat > "${tmpfile}" << 'ENDHTML'
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Bad</title>
</head>
<body>
  <div><p>unclosed
</body>
</html>
ENDHTML
  run npx --yes html-validate@9 "${tmpfile}"
  rm -f "${tmpfile}"
  # We expect non-zero — the tool must reject bad HTML
  [ "${status}" -ne 0 ]
}
