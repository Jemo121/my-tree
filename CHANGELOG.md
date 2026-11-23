# Changelog

All notable changes to this project are documented here.
Format: [Keep a Changelog](https://keepachangelog.com) — versions follow `vYYYY.N.PATCH`.

---

## [v2026.2.1] — 2026-11-22

### Added

- Staging deploy job runs on the `dev` branch: builds the Hugo site with Wrangler 4 and publishes a Cloudflare Pages preview at `https://staging.dev.eguia.me`.

### Changed

- Wrangler config migrated from `wrangler.jsonc` to `wrangler.toml`. The new file adds a `[previews.vars]` block that sets `ENVIRONMENT = "preview"` on non-production deployments.

### Fixed

- Worker name corrected from `ackerman` to `my-tree-ackerman` in `wrangler.toml`.

---

## [v2026.1.0] — 2026-11-22

### Added

- `wrangler.jsonc` — initial Cloudflare Workers Builds configuration (`name = "ackerman"`, `compatibility_date = "2026-11-22"`).

---

## [v2026.0.0] — 2026-11-22

### Added

- Full GitLab CI pipeline with four stages: `security`, `validate`, `build`, `deploy`.
  - `secrets-scan` (Gitleaks v8.18.4) — blocks on detected secrets in MRs and default branch.
  - `sast` (Semgrep 1.72.0) — scans `.gitlab-ci.yml` with the `p/ci` ruleset.
  - `html-validate` — runs html-validate@9 on all templates and error pages.
  - `size-checker` — warns at 300 lines per file, fails at 400 lines.
  - `go-compat` — verifies the Hugo binary was compiled with Go >= 1.23.
  - `hugo-build` — validates the site builds without errors on the default branch when `hugo.toml` exists.
  - `production-deploy` — triggers a Cloudflare Pages build via deploy hook on `v*` tags.
  - `subdomain-sync` — reads `data/family/members.yaml` and creates missing Cloudflare DNS CNAMEs. Skips gracefully when `CF_PAGES_DOMAIN` is unset.
- `data/family/members.yaml` — family member registry. Each record maps a subdomain to a template and drives the DNS sync job.
- Bats test suites under `tests/`:
  - `members-schema.bats` — validates required fields, subdomain format, template file existence, status values, and duplicate detection.
  - `validate-html.bats` — runs html-validate@9 locally against all template sets and error pages.
  - `template-pii.bats` — scans templates, error pages, and data files for personal email addresses.
  - `template-accessibility.bats` — checks `lang`, `charset`, `viewport`, `title`, aria attributes on decorative elements, and `<main>` landmark presence.
  - `ci-guards.bats` — unit tests for the CI shell logic: secret variable guard, Go version extraction, and awk size-checker.
  - `subdomain-sync.bats` — tests the DNS sync shell logic against mock curl responses; includes a fixture contract test suite.
  - `pro-images.bats` — network smoke test for all picsum.photos CDN URLs in pro templates.
- `Makefile` — `make test` runs all offline suites; `make test-images` runs the network test; `make install-hooks` installs the pre-commit hook.
- `tests/pre-commit-hook.sh` — pre-commit hook running HTML validation, PII scan, YAML schema check, and size check locally before push.
- GitLab MR templates: `default`, `bugfix`, `feature`, `hotfix`, `infrastructure`, `refactor`, `release`.
- `README.md` and `LICENSE` (AGPL-3.0-or-later).

---

[v2026.2.1]: https://gitlab.com/jemo121/my-tree/-/compare/v2026.1.0...v2026.2.1
[v2026.1.0]: https://gitlab.com/jemo121/my-tree/-/compare/v2026.0.0...v2026.1.0
[v2026.0.0]: https://gitlab.com/jemo121/my-tree/-/releases/v2026.0.0
