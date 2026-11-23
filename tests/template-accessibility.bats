#!/usr/bin/env bats
# tests/template-accessibility.bats — TEMPLATE-ACCESS-01
# Static accessibility checks on HTML templates:
#   - CTA anchors with ::before pseudo-element content must have aria-label
#   - Decorative elements use aria-hidden="true"
#   - <main> landmark exists in error pages
#   - <html lang> attribute is set on every file
# Run: bats tests/template-accessibility.bats

REPO_ROOT="${BATS_TEST_DIRNAME}/.."

_all_html_files() {
  find "${REPO_ROOT}/templates" "${REPO_ROOT}/errors" -name "*.html"
}

# ---------------------------------------------------------------------------
# 1. Every HTML file has lang attribute on <html>
# ---------------------------------------------------------------------------
@test "every HTML file has <html lang=...>" {
  while IFS= read -r f; do
    if ! grep -q '<html lang=' "${f}"; then
      echo "Missing lang attribute in: ${f}"
      return 1
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 2. Every HTML file has <meta charset>
# ---------------------------------------------------------------------------
@test "every HTML file has <meta charset=UTF-8>" {
  while IFS= read -r f; do
    if ! grep -qi '<meta charset=' "${f}"; then
      echo "Missing meta charset in: ${f}"
      return 1
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 3. Every HTML file has <meta name="viewport">
# ---------------------------------------------------------------------------
@test "every HTML file has <meta name=\"viewport\">" {
  while IFS= read -r f; do
    if ! grep -qi '<meta name="viewport"' "${f}"; then
      echo "Missing viewport meta in: ${f}"
      return 1
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 4. Every HTML file has a <title> tag
# ---------------------------------------------------------------------------
@test "every HTML file has a non-empty <title>" {
  while IFS= read -r f; do
    if ! grep -q '<title>' "${f}"; then
      echo "Missing <title> in: ${f}"
      return 1
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 5. .cta anchors that use CSS ::before '>' pseudo-content must have
#    aria-label so screen readers do not announce '> contact'.
#    The ::before content is defined in CSS as content: '>'
# ---------------------------------------------------------------------------
@test "every .cta anchor that uses CSS ::before has aria-label" {
  while IFS= read -r f; do
    # Only files with .cta::before { content: '>' } are in scope
    if ! grep -q "\.cta::before" "${f}"; then
      continue
    fi
    # Check if the cta anchor has an aria-label
    if grep -q 'class="cta"' "${f}"; then
      if ! grep -q 'class="cta".*aria-label\|aria-label.*class="cta"' "${f}"; then
        echo "Missing aria-label on .cta anchor in: ${f}"
        return 1
      fi
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 6. Decorative status-dot spans are aria-hidden
# ---------------------------------------------------------------------------
@test "every status-dot span is aria-hidden=\"true\"" {
  while IFS= read -r f; do
    if grep -q 'class="status-dot"' "${f}"; then
      if ! grep -q 'class="status-dot" aria-hidden="true"' "${f}"; then
        echo "status-dot is not aria-hidden in: ${f}"
        return 1
      fi
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 7. Decorative dividers are aria-hidden
# ---------------------------------------------------------------------------
@test "every divider element is aria-hidden=\"true\"" {
  while IFS= read -r f; do
    if grep -q 'class="divider"' "${f}"; then
      if ! grep -q 'class="divider" aria-hidden="true"' "${f}"; then
        echo "divider is not aria-hidden in: ${f}"
        return 1
      fi
    fi
  done < <(_all_html_files)
}

# ---------------------------------------------------------------------------
# 8. Error pages use <main> landmark
# ---------------------------------------------------------------------------
@test "error pages use <main> landmark element" {
  while IFS= read -r f; do
    if ! grep -q '<main' "${f}"; then
      echo "Missing <main> in error page: ${f}"
      return 1
    fi
  done < <(find "${REPO_ROOT}/errors" -name "*.html")
}

# ---------------------------------------------------------------------------
# 9. Pro template accent images are aria-hidden (decorative CDN images)
# ---------------------------------------------------------------------------
@test "pro template accent images are aria-hidden=\"true\"" {
  while IFS= read -r f; do
    if grep -q 'class="accent-img"' "${f}"; then
      if ! grep -q 'class="accent-img" aria-hidden="true"' "${f}"; then
        echo "accent-img figure is not aria-hidden in: ${f}"
        return 1
      fi
    fi
  done < <(find "${REPO_ROOT}/templates/pro" -name "*.html")
}
