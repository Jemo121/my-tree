#!/bin/sh
# tests/pre-push-hook.sh
# Pre-push hook — catches issues before CI burns shared minutes.
# Install via: make install-hooks
# Mirrors the CI validate stage for fast, free local checks.

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
EXIT=0

cd "${REPO_ROOT}"

# ---------------------------------------------------------------------------
# 1. Size checker — all CI-scoped file types (SIZE-CHECK-01)
# ---------------------------------------------------------------------------
echo "[pre-push] Running size check..."
result=$(find . -not -path './.git/*' -not -name '*.wasm' \
  \( -name '*.html' -o -name '*.go' -o -name '*.css' \
     -o -name '*.yaml' -o -name '*.yml' -o -name '*.toml' \
     -o -name '*.sh' -o -name '*.bats' -o -name '*.js' \
     -o -name '*.json' -o -name 'Makefile' \) \
  -exec wc -l {} + 2>/dev/null | awk -v warn=300 -v fail=400 '
  BEGIN { rc=0 }
  /total$/ { next }
  $1+0 >= fail { print "FAIL: " $2 " — " $1 " lines"; rc=1; next }
  $1+0 >= warn { print "WARN: " $2 " — " $1 " lines" }
  END { exit rc+0 }
' 2>/dev/null || true)
if printf '%s' "${result}" | grep -q "^FAIL:"; then
  echo "[pre-push] Size check failures:"
  printf '%s\n' "${result}"
  EXIT=1
elif [ -n "${result}" ]; then
  printf '%s\n' "${result}"
fi

# ---------------------------------------------------------------------------
# 2. YAML lint (YAML-VAL-01)
# ---------------------------------------------------------------------------
if command -v yamllint >/dev/null 2>&1; then
  echo "[pre-push] Running yamllint..."
  if ! yamllint .; then
    echo "[pre-push] FAIL: yamllint found issues"
    EXIT=1
  fi
else
  echo "[pre-push] SKIP: yamllint not found — pip install yamllint"
fi

# ---------------------------------------------------------------------------
# 3. Shell check on hook scripts (SHELL-CHECK-01)
# ---------------------------------------------------------------------------
if command -v shellcheck >/dev/null 2>&1; then
  echo "[pre-push] Running shellcheck..."
  if ! shellcheck tests/pre-commit-hook.sh tests/pre-push-hook.sh; then
    echo "[pre-push] FAIL: shellcheck found issues"
    EXIT=1
  fi
else
  echo "[pre-push] SKIP: shellcheck not found — install shellcheck"
fi

exit "${EXIT}"
