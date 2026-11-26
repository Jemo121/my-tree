# Makefile — local test runner and hook installer
# Mirrors the CI validate stage so local runs catch the same issues.

.PHONY: test test-unit test-schema test-html test-pii test-access test-images \
        install-hooks uninstall-hooks copy-wasm-exec wasm wasm-build sync-members help

BATS     ?= bats
REPO_ROOT := $(shell git rev-parse --show-toplevel)

## Run all unit/static tests (no network required)
test: test-schema test-html test-pii test-access test-ci

## Validate data/family/members.yaml schema (YAML-VAL-01)
test-schema:
	@echo "==> members-schema"
	$(BATS) tests/members-schema.bats

## Validate HTML templates with html-validate@9 (HTML-VAL-01)
test-html:
	@echo "==> validate-html"
	$(BATS) tests/validate-html.bats

## Check for personal email addresses in templates (TEMPLATE-PII-01)
test-pii:
	@echo "==> template-pii"
	$(BATS) tests/template-pii.bats

## Static accessibility checks on templates (TEMPLATE-ACCESS-01)
test-access:
	@echo "==> template-accessibility"
	$(BATS) tests/template-accessibility.bats

## CI shell logic unit tests (CI-DEPLOY-01, CI-GOCOMPAT-01, SIZE-CHECK-01, CI-SAST-01, CI-SUBDOMAIN-01)
test-ci:
	@echo "==> ci-guards"
	$(BATS) tests/ci-guards.bats
	@echo "==> subdomain-sync"
	$(BATS) tests/subdomain-sync.bats

## CDN image smoke tests — requires network (TEMPLATE-EXTERNAL-01)
test-images:
	@echo "==> pro-images"
	$(BATS) tests/pro-images.bats

## Install pre-commit and pre-push hooks
install-hooks:
	@echo "==> Installing git hooks"
	cp tests/pre-commit-hook.sh $(REPO_ROOT)/.git/hooks/pre-commit
	chmod +x $(REPO_ROOT)/.git/hooks/pre-commit
	cp tests/pre-push-hook.sh $(REPO_ROOT)/.git/hooks/pre-push
	chmod +x $(REPO_ROOT)/.git/hooks/pre-push
	@echo "    Installed: .git/hooks/pre-commit"
	@echo "    Installed: .git/hooks/pre-push"

## Remove git hooks
uninstall-hooks:
	@echo "==> Removing git hooks"
	rm -f $(REPO_ROOT)/.git/hooks/pre-commit
	rm -f $(REPO_ROOT)/.git/hooks/pre-push

## Fetch wasm_exec.js from Go SDK (requires WASM_GO_VERSION, default 1.27.4)
copy-wasm-exec:
	@echo "==> copy-wasm-exec"
	wget -qO static/wasm_exec.js \
	  "https://raw.githubusercontent.com/golang/go/go$${WASM_GO_VERSION:-1.27.4}/lib/wasm/wasm_exec.js"

## Sync data/family/members.yaml → wasm/assets/members.json (requires yq v4)
sync-members:
	@echo "==> sync-members"
	yq -o=json '.' data/family/members.yaml > wasm/assets/members.json

## Build Go WASM binary → static/wasm/tree.wasm (requires Go 1.27+)
wasm: sync-members
	@echo "==> wasm"
	cd wasm && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ../static/wasm/tree.wasm .
	@ls -lh static/wasm/tree.wasm

## Build WASM binary locally (requires Go 1.27+, yq v4) — mirrors CI wasm-build job
wasm-build:
	@echo "==> wasm-build"
	yq -o=json eval '.' data/family/members.yaml > wasm/assets/members.json
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o static/wasm/tree.wasm ./wasm/

help:
	@grep -E '^##' Makefile | sed 's/^## //'
