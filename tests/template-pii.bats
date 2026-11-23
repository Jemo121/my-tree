#!/usr/bin/env bats
# tests/template-pii.bats — TEMPLATE-PII-01
# Asserts no personal email addresses (@gmail.com or other personal domains)
# are present in HTML templates or error pages.
# Run: bats tests/template-pii.bats

REPO_ROOT="${BATS_TEST_DIRNAME}/.."

# Personal email domains to block. Add more as needed.
PERSONAL_DOMAINS="@gmail.com @yahoo.com @hotmail.com @outlook.com @icloud.com @protonmail.com"

@test "templates/base/ contains no personal email addresses" {
  for domain in ${PERSONAL_DOMAINS}; do
    matches=$(grep -rl "${domain}" "${REPO_ROOT}/templates/base/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "Personal email domain '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}

@test "templates/meme/ contains no personal email addresses" {
  for domain in ${PERSONAL_DOMAINS}; do
    matches=$(grep -rl "${domain}" "${REPO_ROOT}/templates/meme/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "Personal email domain '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}

@test "templates/pro/ contains no personal email addresses" {
  for domain in ${PERSONAL_DOMAINS}; do
    matches=$(grep -rl "${domain}" "${REPO_ROOT}/templates/pro/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "Personal email domain '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}

@test "errors/ contains no personal email addresses" {
  for domain in ${PERSONAL_DOMAINS}; do
    matches=$(grep -rl "${domain}" "${REPO_ROOT}/errors/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "Personal email domain '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}

@test "data/ contains no personal email addresses" {
  for domain in ${PERSONAL_DOMAINS}; do
    matches=$(grep -rl "${domain}" "${REPO_ROOT}/data/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "Personal email domain '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}

@test "no mailto: links with personal email domains anywhere in repo" {
  for domain in ${PERSONAL_DOMAINS}; do
    # mailto:*@domain.com pattern
    pattern="mailto:[^\"']*${domain}"
    matches=$(grep -rl "${pattern}" "${REPO_ROOT}/templates/" "${REPO_ROOT}/errors/" 2>/dev/null || true)
    if [ -n "${matches}" ]; then
      echo "mailto: link with '${domain}' found in:"
      echo "${matches}"
      return 1
    fi
  done
}
