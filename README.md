# actiondoc

> Generate documentation for GitHub Actions workflows using JSDoc-style comments.

actiondoc reads GitHub Actions workflow (`.github/workflows/*.yml`) and action files
(`action.yml`), extracts both native YAML properties and ActionDoc comment tags, and
generates Markdown documentation.

For the comment standard, see [spec/actiondoc-spec.md](spec/actiondoc-spec.md).
For a detailed history of changes, see [CHANGELOG.md](CHANGELOG.md).

---

## Quick Start

```bash
# Install from source
go install github.com/smol-utils/actiondoc@latest

# Generate docs for all workflows in the current repo
actiondoc generate

# Generate for a specific file
actiondoc generate .github/workflows/ci.yml

# Write to a file instead of stdout
actiondoc generate -o WORKFLOWS.md
```

### Before / After

**Your workflow (input):**

```yaml
# @desc Deploy the application to production.
# @secret AWS_ACCESS_KEY_ID - AWS credentials for deployment
# @secret AWS_SECRET_ACCESS_KEY - AWS credentials for deployment
name: Deploy

on:
  push:
    branches: [main]

jobs:
  # @desc Build and push the Docker image to ECR.
  # @env ECR_REGISTRY - Must be set in repository variables
  build:
    name: Build Image
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Build and push
        run: docker build -t $ECR_REGISTRY/app .
```

**Generated docs (output):**

```markdown
# Deploy

Deploy the application to production.

| Property | Value |
|----------|-------|
| File | `deploy.yml` |
| Triggers | `push` |

## Secrets

| Name | Type | Description |
|------|------|-------------|
| `AWS_ACCESS_KEY_ID` | - | AWS credentials for deployment |
| `AWS_SECRET_ACCESS_KEY` | - | AWS credentials for deployment |

## Jobs

### Build Image (`build`)

Build and push the Docker image to ECR.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Environment Variables:**

| Name | Type | Description |
|------|------|-------------|
| `ECR_REGISTRY` | - | Must be set in repository variables |
```

---

## Comment Tags

ActionDoc uses JSDoc-style comments placed above the YAML node they document. Only
use tags for information that can't be inferred from the YAML itself -- names, dependencies,
and runner info are extracted automatically.

| Tag | Purpose | Example |
|-----|---------|---------|
| `@desc` | Description (multi-line) | `@desc Build and deploy the app.` |
| `@secret` | Required secret | `@secret AWS_KEY - AWS access key` |
| `@input` | Workflow input | `@input env {string} - Target environment` |
| `@env` | External env var | `@env DB_URL - Connection string` |
| `@output` | Produced output | `@output tag {semver} - Image tag` |
| `@deprecated` | Deprecation notice | `@deprecated Use deploy-v2 instead.` |
| `@see` | Related link | `@see https://docs.example.com/deploy` |
| `@since` | Version introduced | `@since v2.0.0` |
| `@example` | Usage example (multi-line) | `@example` followed by indented lines |

Tags can be placed at three levels: above the first key (workflow-level), above a job
key (job-level), or above a step entry (step-level).

For the full specification, see [spec/actiondoc-spec.md](spec/actiondoc-spec.md).

---

## Usage

```bash
# Markdown to stdout (default)
actiondoc generate

# JSON output for tooling integration
actiondoc generate --json

# Specific file or directory
actiondoc generate .github/workflows/deploy.yml
actiondoc generate .github/workflows/

# Write to file
actiondoc generate -o docs/WORKFLOWS.md

# Version
actiondoc version
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-o` | stdout | Write output to a file |
| `--json` | false | Output JSON instead of Markdown |
| path (positional) | `.github/workflows` | File or directory to process |

---

## Install

### From source (requires Go 1.22+)

```bash
go install github.com/smol-utils/actiondoc@latest
```

### From GitHub Releases

Download a pre-built binary from the
[Releases](https://github.com/smol-utils/actiondoc/releases) page.

### Build locally

```bash
git clone https://github.com/smol-utils/actiondoc.git
cd actiondoc
make build
./actiondoc generate
```

---

## Development

```bash
make help       # Show all available targets
make build      # Build the binary
make test       # Run all tests
make lint       # Run go vet
make ci         # Lint + test (local CI check)
make demo       # Run against the sample workflow
make golden     # Regenerate the golden test file
make clean      # Remove build artifacts
```

---

## FAQ

**Why not just read the YAML directly?**

You can, but large repos accumulate dozens of workflow files with hundreds of steps. The
biggest pain point is figuring out which secrets and environment variables need to be
configured before a workflow will run. ActionDoc surfaces that information in one place.

---

**What happens if I don't add any ActionDoc comments?**

The tool still works. It extracts names, triggers, job dependencies, runner info, and
step actions from the YAML itself. Comments add descriptions and surface hidden
requirements (secrets, env vars) -- but they're optional.

---

**How is this different from action-docs?**

Both tools document `action.yml` and workflow files. The differences are footprint and
the comment tag system:

- action-docs extracts structured YAML fields (inputs, outputs, description) from both
  file types.
- actiondoc does the same, plus supports JSDoc-style comment tags (`@secret`, `@env`,
  `@deprecated`, `@example`, etc.) for documenting things the YAML can't express --
  like which secrets need to be configured before a workflow will run.

|  | actiondoc | action-docs |
|--|-----------|-------------|
| Workflows | Yes + comment tags | Yes |
| Actions (`action.yml`) | Yes + comment tags | Yes |
| Comment tag system | @desc, @secret, @env, @input, @output, @deprecated, @see, @since, @example | No |
| Language | Go | Node.js |
| Direct deps | 1 | 6 |
| Transitive deps | 0 | 28 |
| Runtime needed | None (static binary) | Node.js |
| Install | Download binary or `go install` | `npm install` |
