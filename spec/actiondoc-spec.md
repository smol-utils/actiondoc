# ActionDoc Comment Standard (v2)

A JSDoc-inspired documentation standard for GitHub Actions workflow files.

## Philosophy

**DRY (Don't Repeat Yourself).** ActionDoc tags document only what cannot be inferred from the YAML itself. Properties like `name:`, `needs:`, and `runs-on:` are extracted automatically by the parser. Tags are reserved for context that the machine cannot derive -- descriptions, required secrets, external dependencies, and deprecation notices.

## Comment Syntax

ActionDoc comments are standard YAML comments (`#`) placed directly above the node they document. There must be no blank line between the comment block and the YAML node.

```yaml
# @desc Runs the full test suite.
# @secret NPM_TOKEN - Required for private package access
test:
  name: Run Tests
  runs-on: ubuntu-latest
```

## Multi-line Values

A continuation line is any comment line that does not start with `@`. It is appended to the most recently declared tag:

```yaml
# @desc Runs the full test suite including
#   unit tests, integration tests, and
#   end-to-end browser tests.
```

## Tag Allowlist

ActionDoc recognizes a **fixed allowlist** of `@`-prefixed tags in YAML comments:

`@desc`, `@secret`, `@input`, `@env`, `@output`, `@deprecated`, `@see`, `@since`,
`@example`

Any other `@`-prefixed string in a comment is **ignored** -- not interpreted, not
warned about. This guarantee lets ActionDoc coexist with other tooling that embeds its
own `@`-markers in YAML comments (e.g. Lula compliance markers `@lulaStart` / `@lulaEnd`
in defenseunicorns/uds-core). The allowlist is fixed by design; it is not configurable,
so downstream tools can rely on exactly which markers ActionDoc consumes.

## Tag Reference

### `@desc <text>`

Free-form description. Supports multi-line. Applied at workflow, job, or step level.

### `@secret <name> [- <description>]`

Documents a secret that must be configured in the repository or organization settings for this workflow/job to function.

```yaml
# @secret DEPLOY_KEY - SSH key used for deployment to production
```

### `@input <name> [{type}] [- <description>]`

Documents an input the annotated scope expects: a `workflow_dispatch` / reusable workflow input at workflow level, or a forwarded / `with:` input at job and step level. The optional `{type}` indicates the expected format.

```yaml
# @input environment {string} - Target deployment environment (staging or production)
```

### `@env <name> [- <description>]`

Documents an environment variable that must be provided externally (e.g., from repository settings or a calling workflow). Do not use this for env vars defined inline in the YAML -- those are extracted automatically.

```yaml
# @env DATABASE_URL - Connection string for the test database
```

### `@output <name> [{format}] [- <description>]`

Documents an output produced by a workflow, job, or step. The optional `{format}` describes the value shape (e.g., `{url}`, `{semver}`, `{json}`, `{path}`).

```yaml
# @output image-tag {semver} - The Docker image tag that was built
```

### `@deprecated [reason]`

Marks a workflow, job, or step as deprecated. Generates a warning banner in the documentation.

```yaml
# @deprecated Use the new "deploy-v2" workflow instead.
```

### `@see <url-or-reference>`

Links to related resources -- other workflows, runbooks, architecture docs, dashboards. Can appear multiple times.

```yaml
# @see https://docs.example.com/deploy-runbook
# @see .github/workflows/rollback.yml
```

### `@since <version>`

Indicates when a workflow, job, or step was introduced. Useful for reusable workflows consumed by multiple repositories.

```yaml
# @since v2.1.0
```

### `@example`

Provides a usage example. All continuation lines form the example body.

```yaml
# @example
#   gh workflow run deploy.yml \
#     -f environment=staging \
#     -f version=1.2.3
```

## Scope

### Workflow files (`.github/workflows/*.yml`)

Tags can be placed at three levels:

| Level    | Placement                                  |
|----------|--------------------------------------------|
| Workflow | Above the first key in the file (e.g., above `name:`) |
| Job      | Above the job key inside `jobs:`           |
| Step     | Above the step entry in `steps:`           |

### Action metadata files (`action.yml`)

Tags are placed above the first key in the file (action-level). The structured `inputs`,
`outputs`, `description`, and `runs` fields are extracted automatically. Tags supplement
them with information the YAML can't express: `@secret`, `@env`, `@deprecated`, `@since`,
`@see`, and `@example`.

## What the Parser Extracts Automatically

### From workflow files

- `name:` -- display name for workflows, jobs, and steps
- `on:` -- trigger events
- `needs:` -- job dependencies
- `runs-on:` -- runner environment
- `if:` -- conditional expressions
- `uses:` -- action reference for steps
- `id:` -- step identifier

### From action.yml files

- `name:` -- action name
- `description:` -- action description
- `inputs:` -- input names, descriptions, required flags, defaults
- `outputs:` -- output names and descriptions
- `runs:` -- execution method (node20, docker, composite)
- `branding:` -- icon and color

## Redacted Output

Generated documentation can carry sensitive-but-unclassified (SBU) material: secret and
variable names, environment-variable names and values, hostnames and URLs, deploy
environment names, and self-hosted runner labels. Redacted mode replaces that material
with stable placeholders so the output can be shared outside the owning team or handed to
an external service -- for example, an LLM asked to write prose documentation from the
structure.

```bash
# Conservative redaction (default profile)
actiondoc generate --redact

# Also write a local reverse map so real names can be restored afterward
actiondoc generate --redact --redact-map redaction-map.json

# Aggressive: additionally replace every literal env:/with: value
actiondoc generate --redact-aggressive
```

Redaction applies to both Markdown and JSON output.

### What is redacted

Redaction is **consistent pseudonymization, not blanking**: the same original always maps
to the same placeholder, so the cross-references that make the output useful -- the
secret/variable inventories, the call-graph transitive requirements, and the per-step
reference tables -- stay intact.

| Category | Placeholder | Source |
|----------|-------------|--------|
| Secret names | `SECRET_n` | `@secret`, `workflow_call` secrets, forwarded `secrets:` keys, `${{ secrets.X }}` |
| Variable names | `VAR_n` | `${{ vars.X }}` |
| Env-var names | `ENV_n` | `env:` keys, `@env`, `${{ env.X }}` |
| Self-hosted runner labels | `RUNNER_n` | `runs-on:` labels (and `group:` names) outside the GitHub-hosted set |
| Deploy environment names | `ENVIRONMENT_n` | `environment:` |
| URLs | `URL_n` | scheme-qualified URLs in any field's literal text (outside `${{ }}`) |
| Hostnames | `HOST_n` | bare dotted hostnames in any field's literal text (outside `${{ }}`) |
| Literal values | `VALUE_n` | every literal `env:`/`with:` value (aggressive only) |

Placeholder numbering is **deterministic**: originals are sorted before they are numbered,
so the same input always produces the same output and diffs stay reviewable.

### What is NOT redacted

- Pipeline structure: jobs, triggers, permissions, the call-graph shape, matrix shape.
- Workflow / job / step / action **names** (titles), which are structural anchors. They
  are kept verbatim (redacting them would desync the call-graph cross-links), so a name is
  the one place a hostname is not swept -- avoid putting sensitive hosts in names if you
  intend to share redacted output.
- `uses:` references, including local paths and `owner/repo` refs -- redacting them would
  break call-graph resolution. (Hosts inside a `uses:`-style value are not treated as URLs.)
- `${{ }}` expression structure and well-known contexts (`github.*`, `matrix.*`,
  `inputs.*`): only the sensitive identifiers inside an expression are replaced.
- `GITHUB_TOKEN` and GitHub-hosted runner labels (`ubuntu-latest`, etc.), which are
  universal and carry no private information.
- Under the conservative profile, harmless literal values such as `node-version: 20`.

### Fields that accept a literal or an expression

GitHub Actions fields routinely accept either a literal or a `${{ expression }}`. Redaction
handles both uniformly: a literal is matched for hosts/URLs (and, under the aggressive
profile, replaced wholesale), while an expression has only its sensitive identifiers
replaced and its structure preserved. A value mixing the two (`prefix-${{ secrets.X }}`)
is handled span by span.

### The reverse map

`--redact-map <path>` writes a JSON map of placeholder to original to a local file. This
closes the external round-trip: redact, let the external service write prose, then
substitute the real names back locally. The map is the one piece of redaction output that
is **not** safe to share -- it is written only to the explicit path given, never to stdout
or mixed into the documentation, and should not be committed.

### Limits

Hostname detection in free text (descriptions, `run:` scripts) is regex-based and
best-effort. Scheme-qualified URLs are always redacted whole. A bare dotted token is
treated as a hostname only when it is all-lowercase, its final label is not a known file
extension, and its final label is a recognized public TLD or internal pseudo-TLD
(`.internal`, `.corp`, `.local`, `.svc`, and similar). This keeps the dotted tokens that
crowd shell scripts -- Java `-D` system properties, Maven coordinates
(`org.codehaus.mojo`), archive names (`*.jar`), git config keys (`user.email`) --
readable. The residual limits run the other way: a single-label internal name with no
dot, or a host under an unrecognized TLD, passes through unredacted. `.name` and `.email`
are deliberately unrecognized because they collide with git config keys; put such hosts
behind a scheme (`https://...`) if they must be caught. Identifier redaction (secret/variable/env names)
is exact, because those come from parsed fields, and covers both the dot and quoted-bracket
expression forms (`secrets.NAME` and `secrets['NAME']`). Under the aggressive profile, a
bare host that appears both as a literal value and elsewhere in free text is pseudonymized
consistently only in the conservative profile; identifier consistency holds in both.

The set of GitHub-hosted runner labels that are kept readable is maintained by hand, so a
newly introduced hosted label may not be recognized and would be redacted to `RUNNER_n`.
That is safe -- over-redaction, never a leak -- but means an unfamiliar `RUNNER_n` is not
proof a runner is self-hosted.

## Design Principles

1. **No CI impact.** ActionDoc comments are standard YAML comments. They are invisible to GitHub Actions and all YAML parsers. Adding them cannot break your CI.
2. **Familiar syntax.** The `@tag` comment pattern is familiar and widespread (JSDoc, JavaDoc, Python docstrings).
