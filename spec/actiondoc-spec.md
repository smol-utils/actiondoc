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

## Tag Reference

### `@desc <text>`

Free-form description. Supports multi-line. Applied at workflow, job, or step level.

### `@secret <name> [- <description>]`

Documents a secret that must be configured in the repository or organization settings for this workflow/job to function.

```yaml
# @secret DEPLOY_KEY - SSH key used for deployment to production
```

### `@input <name> [{type}] [- <description>]`

Documents a `workflow_dispatch` input or a reusable workflow input. The optional `{type}` indicates the expected format.

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

## Design Principles

1. **Zero risk.** ActionDoc comments are standard YAML comments. They are invisible to GitHub Actions and all YAML parsers. Adding them cannot break your CI.
2. **Graceful degradation.** If a workflow has no ActionDoc comments, the parser still produces useful documentation from the native YAML properties.
3. **Familiar syntax.** Developers who have used JSDoc, JavaDoc, or Python docstrings will recognize the `@tag` pattern immediately.
