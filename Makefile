# actiondoc Makefile
# ---------------------------
# Wraps common Go commands for convenience.
# The help target auto-generates a list of available commands.

.DEFAULT_GOAL := help

## help: Show available commands with descriptions
help:
	@echo "Available targets:"
	@grep -E '^## [a-zA-Z0-9_.-]+:' $(MAKEFILE_LIST) \
	  | sed 's/^## \([^:]*\): \(.*\)/\1:\2/' \
	  | sort \
	  | awk -F':' 'BEGIN { max=0 } { if (length($$1)>max) max=length($$1) } { lines[NR]=$$0 } END { for (i=1; i<=NR; i++) { split(lines[i], parts, ":"); printf "  %-"max+2"s %s\n", parts[1]":", parts[2] } }'

## build: Build the actiondoc binary
build:
	go build -o actiondoc .

## test: Run all tests
test:
	go test ./... -count=1

## fmt-check: Fail if any file is not gofmt-clean (matches CI)
fmt-check:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files are not gofmt-clean:"; echo "$$unformatted"; \
		echo "Run: gofmt -w ."; exit 1; \
	fi

## lint: Check formatting (gofmt) and run Go vet (static analysis)
lint: fmt-check
	go vet ./...

## ci: Local CI check (vet + test)
ci: lint test

## install: Install actiondoc to $GOPATH/bin
install:
	go install .

## clean: Remove build artifacts
clean:
	rm -f actiondoc
	rm -rf dist/

## golden: Regenerate the golden test files (workflow + action)
golden:
	go run . generate testdata/sample-workflow.yml > testdata/expected-output.md
	go run . generate testdata/action.yml > testdata/expected-action-output.md

## demo: Run actiondoc against the sample workflow
demo:
	go run . generate testdata/sample-workflow.yml

DOGFOOD_DIR ?= dogfood/repos

## dogfood-fetch: Shallow-clone the corpus from dogfood/manifest.txt at pinned SHAs
dogfood-fetch:
	@mkdir -p dogfood/repos
	@grep -v '^#' dogfood/manifest.txt | while IFS="$$(printf '\t')" read -r name url sha; do \
		[ -n "$$name" ] || continue; \
		dir="dogfood/repos/$$name"; \
		if [ -d "$$dir/.git" ]; then echo "have $$name"; continue; fi; \
		echo "fetch $$name ($$sha)"; \
		git init -q "$$dir" && git -C "$$dir" remote add origin "$$url" 2>/dev/null; \
		git -C "$$dir" fetch -q --depth 1 origin "$$sha" && git -C "$$dir" checkout -q FETCH_HEAD || \
			{ echo "  ERROR: could not fetch pinned SHA $$sha for $$name"; exit 1; }; \
	done

## dogfood: Run actiondoc against each corpus repo; fail if any repo errors
dogfood: build
	@ok=0; fail=0; failed=""; \
	for d in $(DOGFOOD_DIR)/*/; do \
		[ -d "$$d.github/workflows" ] || continue; \
		name=$$(basename "$$d"); \
		if ./actiondoc generate "$$d.github/workflows" >/dev/null 2>&1; then \
			ok=$$((ok+1)); \
		else \
			fail=$$((fail+1)); failed="$$failed $$name"; \
		fi; \
	done; \
	echo "dogfood: $$ok ok, $$fail failed"; \
	if [ $$fail -gt 0 ]; then echo "failed:$$failed"; exit 1; fi
