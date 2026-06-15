# CI Pipeline

**Triggers:** `push`, `pull_request`

Main CI pipeline for building and testing the application.

| Property | Value |
|----------|-------|
| File | `sample-workflow.yml` |
| Default runs-on | `ubuntu-latest` |
| Since | v1.0.0 |

**See also:** https://docs.example.com/ci

**Jobs:** [Build](#build-build), [Run Tests](#run-tests-test), [Deploy](#deploy-deploy)

## Event filters

- **push**
  - branches: `main`
- **pull_request**
  - branches: `main`

## Secrets

| Name | Type | Description |
|------|------|-------------|
| `DEPLOY_KEY` | - | SSH key used for deployment |

## Environment Variables

| Name | Type | Description |
|------|------|-------------|
| `CI` | - | Must be set to "true" in repository settings |

## Jobs

### Build (`build`)

Compile the application and produce build artifacts.

<details>
<summary>Steps (2)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Build**

</details>

### Run Tests (`test`)

Runs the full test suite against the compiled artifacts.

| Property | Value |
|----------|-------|
| Depends on | `build` |
| Condition | `github.event_name == 'push'` |

**Secrets:**

| Name | Type | Description |
|------|------|-------------|
| `NPM_TOKEN` | - | Required for private package access |

**Environment Variables:**

| Name | Type | Description |
|------|------|-------------|
| `DATABASE_URL` | - | Connection string for the test database |

<details>
<summary>Steps (2)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Run tests** - Execute unit and integration tests.
   - ID: `tests`
   - Output: `test-results` {path} - Path to the JUnit XML report

</details>

### Deploy (`deploy`)

> **Deprecated**: Use the new "build-v2" workflow instead.

Builds the production Docker image and pushes to registry.

| Property | Value |
|----------|-------|
| Depends on | `build`, `test` |

**Example:**

```
  To trigger manually:
  gh workflow run ci.yml -f deploy=true
```

<details>
<summary>Steps (2)</summary>

1. **Build image**

2. **Push image**

</details>

