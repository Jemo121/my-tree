#!/bin/sh
# tests/pre-commit-hook.sh
# Pre-commit hook script. Install via: make install-hooks
# Runs the same checks as CI validate stage locally.
# Fails fast on first error. Skips checks whose tools are not installed.

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
EXIT=0

# ---------------------------------------------------------------------------
# 1. HTML validation (HTML-VAL-01)
# ---------------------------------------------------------------------------
if command -v npx >/dev/null 2>&1; then
  echo "[pre-commit] Running html-validate..."
  cd "${REPO_ROOT}"
  if ! npx --yes html-validate@9 "templates/**/*.html" "errors/*.html"; then
    echo "[pre-commit] FAIL: HTML validation failed"
    EXIT=1
  fi
else
  echo "[pre-commit] SKIP: npx not found — install Node.js to enable HTML validation"
fi

# ---------------------------------------------------------------------------
# 2. PII / personal email check (TEMPLATE-PII-01)
# ---------------------------------------------------------------------------
echo "[pre-commit] Checking for personal email addresses..."
PERSONAL_PATTERNS="@gmail.com @yahoo.com @hotmail.com @outlook.com @icloud.com @protonmail.com"
# shellcheck disable=SC2086  # intentional word splitting on space-separated list
for pattern in ${PERSONAL_PATTERNS}; do
  matches=$(grep -rl "${pattern}" \
    "${REPO_ROOT}/templates/" \
    "${REPO_ROOT}/errors/" \
    "${REPO_ROOT}/data/" \
    2>/dev/null || true)
  if [ -n "${matches}" ]; then
    echo "[pre-commit] FAIL: personal email domain '${pattern}' found in:"
    echo "${matches}"
    EXIT=1
  fi
done

# ---------------------------------------------------------------------------
# 3. YAML schema validation (YAML-VAL-01)
# ---------------------------------------------------------------------------
if command -v yq >/dev/null 2>&1; then
  echo "[pre-commit] Validating data/family/members.yaml..."
  cd "${REPO_ROOT}"
  members="${REPO_ROOT}/data/family/members.yaml"
  count=$(yq '. | length' "${members}")
  idx=0
  while [ "${idx}" -lt "${count}" ]; do
    for field in name subdomain template status; do
      val=$(yq ".[${idx}].${field}" "${members}")
      if [ -z "${val}" ] || [ "${val}" = "null" ]; then
        echo "[pre-commit] FAIL: members.yaml entry ${idx} missing '${field}'"
        EXIT=1
      fi
    done
    sub=$(yq ".[${idx}].subdomain" "${members}")
    if ! printf '%s' "${sub}" | grep -qE '^[a-z0-9-]+$'; then
      echo "[pre-commit] FAIL: subdomain '${sub}' at index ${idx} fails ^[a-z0-9-]+\$"
      EXIT=1
    fi
    tmpl=$(yq ".[${idx}].template" "${members}")
    if [ ! -f "${REPO_ROOT}/templates/base/${tmpl}.html" ]; then
      echo "[pre-commit] FAIL: template '${tmpl}' at index ${idx} has no matching file"
      EXIT=1
    fi
    idx=$((idx + 1))
  done
else
  echo "[pre-commit] SKIP: yq not found — install yq v4 to enable YAML schema validation"
fi

# ---------------------------------------------------------------------------
# 4. Size check (SIZE-CHECK-01)
# ---------------------------------------------------------------------------
echo "[pre-commit] Running size check on HTML files..."
cd "${REPO_ROOT}"
result=$(find templates errors -name '*.html' -exec wc -l {} + 2>/dev/null | awk -v warn=300 -v fail=400 '
  BEGIN { rc=0 }
  /total$/ { next }
  $1+0 >= fail { print "FAIL: " $2 " — " $1 " lines"; rc=1; next }
  $1+0 >= warn { print "WARN: " $2 " — " $1 " lines" }
  END { exit rc+0 }
' || true)
if printf '%s' "${result}" | grep -q "^FAIL:"; then
  echo "[pre-commit] Size check failures:"
  printf '%s\n' "${result}"
  EXIT=1
elif [ -n "${result}" ]; then
  echo "[pre-commit] Size warnings:"
  printf '%s\n' "${result}"
fi

exit "${EXIT}"
