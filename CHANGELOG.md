# Changelog

All notable changes to actiondoc are documented here.

## [Unreleased]

### Added
- Redacted output mode (`--redact`, `--redact-aggressive`, `--redact-map`): consistent pseudonymization of secret/variable/env names, hostnames, URLs, deploy environments, and self-hosted runner labels so generated docs can be shared externally; honored by both Markdown and JSON output, with a deterministic placeholder scheme and an optional local reverse map
- ActionDoc comment standard (v2) with 9 tags: @desc, @secret, @input, @env, @output, @deprecated, @see, @since, @example
- YAML parser using goccy/go-yaml with full comment preservation
- Markdown renderer with workflow, job, and step documentation
- JSON output mode (`--json`) for tooling integration
- CLI with `generate` subcommand, `-o` file output, directory or single-file input
- Golden file end-to-end test
- Standalone spec document (spec/actiondoc-spec.md)
- Sample workflow with all tag types (testdata/sample-workflow.yml)
- Makefile with build, test, lint, and install targets
- CI workflow for GitHub Actions
- Release workflow with goreleaser for cross-platform binaries

### Design Decisions
- Single dependency: github.com/goccy/go-yaml (zero transitive runtime deps)
- DRY tags: only document what YAML can't express (name, needs, runs-on extracted automatically)
- Data model as firewall: only parser.go imports the YAML library
- Errors propagate; only main.go calls os.Exit
