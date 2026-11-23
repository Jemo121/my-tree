#!/usr/bin/env bats
# tests/members-schema.bats — YAML-VAL-01
# Validates data/family/members.yaml schema before subdomain-sync fires.
# Required fields: name, subdomain, template, status (all non-empty).
# subdomain must match ^[a-z0-9-]+$.
# template value must correspond to an existing file in templates/base/.
# Run: bats tests/members-schema.bats
#      (requires yq v4 on PATH)

MEMBERS_FILE="${BATS_TEST_DIRNAME}/../data/family/members.yaml"
TEMPLATES_BASE="${BATS_TEST_DIRNAME}/../templates/base"
REQUIRED_FIELDS=(name subdomain template status)

# ---------------------------------------------------------------------------
# helper: index of current member under test (populated per-test by setup)
# ---------------------------------------------------------------------------
setup() {
  if ! command -v yq >/dev/null 2>&1; then
    skip "yq not found — install yq v4 to run schema tests"
  fi
  # Total member count — used to drive index-based tests
  MEMBER_COUNT=$(yq '. | length' "${MEMBERS_FILE}")
}

# ---------------------------------------------------------------------------
# 1. File exists and is non-empty YAML
# ---------------------------------------------------------------------------
@test "members.yaml exists and parses as valid YAML" {
  [ -f "${MEMBERS_FILE}" ]
  count=$(yq '. | length' "${MEMBERS_FILE}")
  [ "${count}" -gt 0 ]
}

# ---------------------------------------------------------------------------
# 2. Every member has all required fields with non-empty values
# ---------------------------------------------------------------------------
@test "every member has a non-empty 'name' field" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    val=$(yq ".[${idx}].name" "${MEMBERS_FILE}")
    if [ -z "${val}" ] || [ "${val}" = "null" ]; then
      echo "member[${idx}] missing or null 'name'"
      return 1
    fi
    idx=$((idx + 1))
  done
}

@test "every member has a non-empty 'subdomain' field" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    val=$(yq ".[${idx}].subdomain" "${MEMBERS_FILE}")
    if [ -z "${val}" ] || [ "${val}" = "null" ]; then
      echo "member[${idx}] missing or null 'subdomain'"
      return 1
    fi
    idx=$((idx + 1))
  done
}

@test "every member has a non-empty 'template' field" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    val=$(yq ".[${idx}].template" "${MEMBERS_FILE}")
    if [ -z "${val}" ] || [ "${val}" = "null" ]; then
      echo "member[${idx}] missing or null 'template'"
      return 1
    fi
    idx=$((idx + 1))
  done
}

@test "every member has a non-empty 'status' field" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    val=$(yq ".[${idx}].status" "${MEMBERS_FILE}")
    if [ -z "${val}" ] || [ "${val}" = "null" ]; then
      echo "member[${idx}] missing or null 'status'"
      return 1
    fi
    idx=$((idx + 1))
  done
}

# ---------------------------------------------------------------------------
# 3. subdomain format: ^[a-z0-9-]+$  (no spaces, no slashes, no uppercase)
# ---------------------------------------------------------------------------
@test "every subdomain matches ^[a-z0-9-]+\$" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    sub=$(yq ".[${idx}].subdomain" "${MEMBERS_FILE}")
    # skip null entries — the non-empty test already catches those
    if [ "${sub}" = "null" ]; then
      idx=$((idx + 1))
      continue
    fi
    if ! printf '%s' "${sub}" | grep -qE '^[a-z0-9-]+$'; then
      echo "member[${idx}] subdomain '${sub}' fails ^[a-z0-9-]+\$"
      return 1
    fi
    idx=$((idx + 1))
  done
}

@test "subdomain must not contain spaces" {
  yq '.[].subdomain' "${MEMBERS_FILE}" | while IFS= read -r sub; do
    case "${sub}" in
      *" "*) echo "subdomain contains space: '${sub}'"; return 1;;
    esac
  done
}

@test "subdomain must not contain forward slashes" {
  yq '.[].subdomain' "${MEMBERS_FILE}" | while IFS= read -r sub; do
    case "${sub}" in
      *"/"*) echo "subdomain contains slash: '${sub}'"; return 1;;
    esac
  done
}

@test "subdomain must not be empty string (pipe-safety)" {
  # yq can emit empty string for a present-but-empty value
  yq '.[].subdomain' "${MEMBERS_FILE}" | while IFS= read -r sub; do
    if [ -z "${sub}" ]; then
      echo "subdomain is empty string"
      return 1
    fi
  done
}

# ---------------------------------------------------------------------------
# 4. template value maps to a real file in templates/base/
# ---------------------------------------------------------------------------
@test "every member template maps to an existing file in templates/base/" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    tmpl=$(yq ".[${idx}].template" "${MEMBERS_FILE}")
    if [ "${tmpl}" = "null" ]; then
      idx=$((idx + 1))
      continue
    fi
    target="${TEMPLATES_BASE}/${tmpl}.html"
    if [ ! -f "${target}" ]; then
      echo "member[${idx}] template '${tmpl}' has no matching file at ${target}"
      return 1
    fi
    idx=$((idx + 1))
  done
}

# ---------------------------------------------------------------------------
# 5. status field has a known value (active | inactive | pending)
# ---------------------------------------------------------------------------
@test "every member status is one of: active, inactive, pending" {
  local idx=0
  while [ "${idx}" -lt "${MEMBER_COUNT}" ]; do
    st=$(yq ".[${idx}].status" "${MEMBERS_FILE}")
    case "${st}" in
      active|inactive|pending) ;;
      *) echo "member[${idx}] status '${st}' is not a known value"; return 1;;
    esac
    idx=$((idx + 1))
  done
}

# ---------------------------------------------------------------------------
# 6. No duplicate subdomains
# ---------------------------------------------------------------------------
@test "no duplicate subdomains" {
  dupes=$(yq '.[].subdomain' "${MEMBERS_FILE}" | sort | uniq -d)
  if [ -n "${dupes}" ]; then
    echo "duplicate subdomains found: ${dupes}"
    return 1
  fi
}
