# Changelog

All notable changes to actiondoc are documented here. This project adheres to
[Semantic Versioning](https://semver.org).

## [Unreleased]

### Added
- Redacted output mode (`--redact`, `--redact-aggressive`, `--redact-map`): consistent
  pseudonymization of secret/variable/env names, hostnames, URLs, deploy environments, and
  self-hosted runner labels so generated docs can be shared externally. Honored by both
  Markdown and JSON output, with a deterministic placeholder scheme and an optional local
  reverse map for restoring real names.

## [0.2.1] - 2026-06-03

### Added
- Transitive requirements: aggregate every referenced secret/variable and permission
  across the full call graph into one per-entry-point summary.
- Step-level `env:` blocks are parsed, rendered, and scanned for secret/variable references.
- Output snapshots for a diverse real-world corpus, diffed on every run (dogfood
  regression harness).

### Fixed
- Render Matrix rows for every matrix form, not just literal axes; show matrix values in a
  job property row and render job names verbatim.
- Markdown safety: escape markup in step titles and job headings, escape pipes in table
  cells, and make inline code spans safe for values containing backticks.
- Resolve YAML anchors and aliases document-wide before reading fields.
- Resolve same-repo cross-repo workflow refs to in-scope nodes; cross-links reuse the same
  duplicate-name anchor disambiguation as the table of contents.
- Render `@input` tags at job and step scope.
- Skip fully commented-out workflow files with a note instead of failing the run.
- Relabel unresolved external references as "(outside scan scope)".
- Normalize multi-line names to a single line; strip banner framing from comment-derived
  descriptions.

## [0.2.0] - 2026-06-01

### Added
- ActionDoc comment standard (v2) with nine tags: `@desc`, `@secret`, `@input`, `@env`,
  `@output`, `@deprecated`, `@see`, `@since`, `@example`.
- YAML parser (goccy/go-yaml) with full comment preservation.
- Markdown renderer for workflow, job, and step documentation.
- JSON output mode (`--json`) for tooling integration.
- Call-graph builder and reusable-workflow rendering: call-graph trees rooted at each entry
  point, a "Workflow call API" (inputs/secrets/outputs), and reverse "Called by" links.
- Declared-surface rendering: triggers and event filters, permissions, environments,
  concurrency, and matrices.
- Composite-action discovery under `.github/actions/`.
- CLI `generate` subcommand with `-o` file output and directory or single-file input.
- Standalone spec (`spec/actiondoc-spec.md`), dogfood corpus harness, and a goreleaser
  release workflow for cross-platform binaries.

### Design notes
- Single runtime dependency (goccy/go-yaml); zero transitive runtime dependencies.
- DRY tags: document only what the YAML cannot express (names, `needs:`, `runs-on:`, and
  inputs are extracted automatically).
- The data model is a firewall: only the parser imports the YAML library.
