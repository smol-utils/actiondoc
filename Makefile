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

## lint: Run Go vet (static analysis)
lint:
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

## golden: Regenerate the golden test file
golden:
	go run . generate testdata/sample-workflow.yml > testdata/expected-output.md

## demo: Run actiondoc against the sample workflow
demo:
	go run . generate testdata/sample-workflow.yml
