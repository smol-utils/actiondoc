# CI Pipeline

Main CI pipeline for building and testing the application.

| Property | Value |
|----------|-------|
| File | `sample-workflow.yml` |
| Triggers | `push`, `pull_request` |
| Since | v1.0.0 |

**See also:** https://docs.example.com/ci

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

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout** - Check out the repository at the triggering commit.
   - Uses: `actions/checkout@v4`

2. **Build**

### Run Tests (`test`)

Runs the full test suite against the compiled artifacts.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
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

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Run tests** - Execute unit and integration tests.
   - ID: `tests`
   - Output: `test-results` {path} - Path to the JUnit XML report

### Deploy (`deploy`)

> **Deprecated**: Use the new "build-v2" workflow instead.

Builds the production Docker image and pushes to registry.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build`, `test` |

**Example:**

```
  To trigger manually:
  gh workflow run ci.yml -f deploy=true
```

#### Steps

1. **Build image**

2. **Push image**

