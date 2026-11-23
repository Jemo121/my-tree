# my-tree

Personal website — [emilio.eguia.me](https://emilio.eguia.me)

Built with [Hugo](https://gohugo.io) · Pure CSS · [Cloudflare Pages](https://pages.cloudflare.com)

No JS frameworks, no TypeScript. The only permitted JavaScript is `wasm_exec.js` + `WebAssembly.instantiateStreaming` for Go WASM bootstrapping.

---

## Stack

| Layer | Choice |
|-------|--------|
| SSG | Hugo 0.147.9 |
| Styling | Pure CSS — `@keyframes`, custom properties, no frameworks |
| Hosting | Cloudflare Pages |
| Workers config | Wrangler 4 (`wrangler.toml`) |
| CI/CD | GitLab CI → Cloudflare Pages |
| Mirror | GitLab (primary) → GitHub |

---

## Project layout

```
my-tree/
├── templates/           # Under-construction page designs
│   ├── base/            # 10 originals — no images
│   ├── pro/             # + picsum.photos accent image (CDN, zero binaries)
│   └── meme/            # + inline SVG corner meme (~200–600 bytes)
├── errors/              # Error pages (404)
├── data/
│   └── family/
│       └── members.yaml # Family member records — drives subdomain DNS sync
├── tests/               # Bats test suites and pre-commit hook
└── (hugo site — Phase 1, not started)
```

---

## Local preview

Requirements: any browser. No build step needed for templates.

```sh
git clone git@gitlab.com:jemo121/my-tree.git
cd my-tree
xdg-open templates/base/warm.html
```

Hugo site (Phase 1, not started):

```sh
hugo server -D
```

---

## Tests

Tests run with [Bats](https://bats-core.readthedocs.io). The Makefile mirrors the CI validate stage exactly.

```sh
# Install Bats if needed
# https://bats-core.readthedocs.io/en/stable/installation.html

# Run all offline tests (no network)
make test

# Run individual suites
make test-schema    # YAML schema for data/family/members.yaml
make test-html      # HTML validity (requires Node.js / npx)
make test-pii       # Personal email address scan
make test-access    # Static accessibility checks
make test-ci        # CI shell logic and subdomain-sync contract

# Network test (requires internet — skipped in CI by default)
make test-images    # CDN URL smoke test for pro templates
```

Install the pre-commit hook to run the same checks before every commit:

```sh
make install-hooks
```

The hook requires `npx` (HTML validation) and `yq` v4 (YAML schema). It skips checks whose tools are absent rather than blocking the commit.

---

## data/family/members.yaml

Each record drives one Cloudflare DNS CNAME during the `subdomain-sync` CI job.

```yaml
- name: Emilio Eguia
  subdomain: emilio        # used as <subdomain>.eguia.me — must match ^[a-z0-9-]+$
  template: sharp          # maps to templates/base/<template>.html
  generation: 1
  parent: null
  status: active           # active | inactive | pending
```

Required fields: `name`, `subdomain`, `template`, `status`. Duplicate subdomains fail CI.

---

## CI pipeline

Four stages run on merge requests and the default branch:

| Stage | Job | What it does |
|-------|-----|--------------|
| security | `secrets-scan` | Gitleaks — exits 1 on any detected secret |
| security | `sast` | Semgrep `p/ci` ruleset on `.gitlab-ci.yml` |
| validate | `html-validate` | html-validate@9 on all templates and error pages |
| validate | `size-checker` | WARN at 300 lines, FAIL at 400 lines per file |
| validate | `go-compat` | Asserts Hugo binary's embedded Go version meets `GO_VERSION` |
| build | `hugo-build` | Builds the Hugo site; runs only when `hugo.toml` exists |
| deploy | `production-deploy` | Triggers Cloudflare Pages via deploy hook on `v*` tags |
| deploy | `staging-deploy` | Pushes to `staging` preview on the `dev` branch |
| deploy | `subdomain-sync` | Creates missing Cloudflare DNS CNAMEs from `members.yaml` |

### Required CI variables

Set these in GitLab project settings before the deploy stage can run:

| Variable | Used by |
|----------|---------|
| `CF_DEPLOY_HOOK` | `production-deploy` |
| `CF_ZONE_ID` | `subdomain-sync` |
| `CF_DNS_TOKEN` | `subdomain-sync` |
| `CF_PAGES_DOMAIN` | `subdomain-sync` (skips DNS sync when unset) |

`subdomain-sync` exits 0 and prints a skip message when `CF_PAGES_DOMAIN` is unset, so the pipeline does not block before Cloudflare Pages is provisioned.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) if it exists, or open an issue on GitLab.

MR templates are available under `.gitlab/merge_request_templates/`.

---

## License

[AGPL-3.0-or-later](LICENSE)
