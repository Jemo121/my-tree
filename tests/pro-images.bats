#!/usr/bin/env bats
# tests/pro-images.bats — TEMPLATE-EXTERNAL-01
# Smoke-tests that all picsum.photos CDN URLs used in pro templates
# return HTTP 200. This is an integration test — requires network.
# Run: bats tests/pro-images.bats
#      Skip in offline / air-gapped environments.

REPO_ROOT="${BATS_TEST_DIRNAME}/.."

setup() {
  if ! command -v curl >/dev/null 2>&1; then
    skip "curl not found — required for image URL smoke test"
  fi
  # Allow skipping via env var for offline CI runners
  if [ "${SKIP_NETWORK_TESTS:-}" = "1" ]; then
    skip "SKIP_NETWORK_TESTS=1 — skipping network integration tests"
  fi
}

# ---------------------------------------------------------------------------
# 1. Extract all picsum URLs from pro templates and check HTTP 200
# ---------------------------------------------------------------------------
@test "all picsum.photos URLs in pro templates return HTTP 200" {
  local failed=0
  local checked=0

  while IFS= read -r url; do
    [ -z "${url}" ] && continue
    http_code=$(curl -sfL --max-time 15 --write-out '%{http_code}' --output /dev/null "${url}" 2>/dev/null || echo "000")
    checked=$((checked + 1))
    if [ "${http_code}" != "200" ]; then
      echo "FAIL: ${url} returned HTTP ${http_code}"
      failed=$((failed + 1))
    fi
  done < <(grep -rh 'src="https://picsum.photos' "${REPO_ROOT}/templates/pro/" | \
           grep -oE 'https://picsum\.photos/seed/[^"]+' | sort -u)

  if [ "${checked}" -eq 0 ]; then
    echo "FAIL: no picsum URLs found in templates/pro/"
    return 1
  fi

  echo "Checked ${checked} picsum URLs, ${failed} failed"
  [ "${failed}" -eq 0 ]
}

# ---------------------------------------------------------------------------
# 2. Every pro template contains exactly one accent-img figure
# ---------------------------------------------------------------------------
@test "every pro template has exactly one accent-img figure" {
  while IFS= read -r f; do
    count=$(grep -c 'class="accent-img"' "${f}" || true)
    if [ "${count}" -ne 1 ]; then
      echo "Expected 1 accent-img in ${f}, found ${count}"
      return 1
    fi
  done < <(find "${REPO_ROOT}/templates/pro" -name "*.html")
}

# ---------------------------------------------------------------------------
# 3. All picsum URLs use the /seed/ path format (not /id/ — seeds are stable)
# ---------------------------------------------------------------------------
@test "all pro template picsum URLs use /seed/ path (stable references)" {
  while IFS= read -r f; do
    # If the file has a picsum URL that does NOT use /seed/, fail
    if grep -q 'picsum.photos/' "${f}"; then
      if grep -q 'picsum.photos/[0-9]' "${f}"; then
        echo "Non-seed picsum URL found in: ${f}"
        return 1
      fi
    fi
  done < <(find "${REPO_ROOT}/templates/pro" -name "*.html")
}

# ---------------------------------------------------------------------------
# 4. picsum URLs include explicit width and height parameters
# ---------------------------------------------------------------------------
@test "all pro template picsum URLs include width and height" {
  while IFS= read -r url; do
    [ -z "${url}" ] && continue
    # URL format: https://picsum.photos/seed/<seed>/<width>/<height>
    if ! printf '%s' "${url}" | grep -qE 'picsum\.photos/seed/[^/]+/[0-9]+/[0-9]+'; then
      echo "picsum URL missing width/height: ${url}"
      return 1
    fi
  done < <(grep -rh 'src="https://picsum.photos' "${REPO_ROOT}/templates/pro/" | \
           grep -oE 'https://picsum\.photos/seed/[^"]+')
}
