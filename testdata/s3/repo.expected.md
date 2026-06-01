# Contents

- [CI](#ci)
- [Reusable Build](#reusable-build)
- [Deploy](#deploy)

# CI

| Property | Value |
|----------|-------|
| File | `ci.yml` |
| Triggers | `push` |

## Call graph (rooted at this workflow)

```
ci.yml [push]
+-- build (uses reusable.yml)
+-- deploy / Deploy to staging (uses ./.github/actions/deploy)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `BUILD_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CI_BUILD_TOKEN` | job `build` secrets `BUILD_TOKEN` |
| `DEPLOY_TOKEN` | job `deploy` step `Deploy to staging` with `token` |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Uses workflow | [Reusable Build](#reusable-build) |

#### Inputs forwarded

- `target`: `linux`

#### Secrets forwarded

- `BUILD_TOKEN`: `${{ secrets.CI_BUILD_TOKEN }}`

### `deploy`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build` |

#### Steps

1. **Deploy to staging** `[continue-on-error]`
   - Uses: `./.github/actions/deploy`
   - With:
     - `environment`: `staging` - Target environment name (required)
     - `token`: `${{ secrets.DEPLOY_TOKEN }}` - Deployment token

# Reusable Build

| Property | Value |
|----------|-------|
| File | `reusable.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `target` | string | Yes | - | Build target |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `BUILD_TOKEN` | Yes | - |

## Called by

```
reusable.yml
+-- ci.yml (job: build)  <- entry point
```

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Build**

# Deploy

Deploys the application to an environment.

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `environment` | Target environment name | Yes | - |
| `token` | Deployment token | No | - |

