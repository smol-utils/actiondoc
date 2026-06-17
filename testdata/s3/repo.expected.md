# repo

1 workflow, 1 reusable workflow, 1 composite action

## Contents

**Workflows**

- [CI - push](#ci)

**Reusable workflows**

- [Reusable Build](#reusable-build)

**Composite actions**

- [Deploy](#deploy-1)

## Workflows by trigger

- **push**: [CI](#ci)

## Secrets and variables used across this repository

**Secrets:**

| Name | Used by |
|------|---------|
| `CI_BUILD_TOKEN` | [CI](#ci) |
| `DEPLOY_TOKEN` | [CI](#ci) |

## Permissions across this repository

| Scope | Level |
|-------|-------|
| `contents` | `read` |
| `id-token` | `write` (OIDC) |

# CI

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `ci.yml` |

**Jobs:** [`build`](#build), [`deploy`](#deploy)

## Call graph (rooted at this workflow)

- `build` uses [reusable.yml](#reusable-build)
- `deploy / Deploy to staging` uses [./.github/actions/deploy](#deploy-1)

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `BUILD_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CI_BUILD_TOKEN` | `build`: (`BUILD_TOKEN`) |
| `DEPLOY_TOKEN` | `deploy`: Deploy to staging (`token`) |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Uses workflow | [Reusable Build](#reusable-build) |
| Matrix | `target`: linux, darwin (combinations adjusted by include/exclude) |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

#### Inputs forwarded

- `target`: `${{ matrix.target }}`

#### Secrets forwarded

- `BUILD_TOKEN`: `${{ secrets.CI_BUILD_TOKEN }}`

### `deploy`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build` |

<details>
<summary>Steps (1)</summary>

1. **Deploy to staging** `[continue-on-error]`
   - Uses: `./.github/actions/deploy`
   - With:
     - `environment`: `staging` - Target environment name (required)
     - `token`: `${{ secrets.DEPLOY_TOKEN }}` - Deployment token

</details>

[Back to contents](#contents)

# Reusable Build

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `reusable.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `target` | string | Yes | - | Build target |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `BUILD_TOKEN` | Yes | - |

## Called by

- [ci.yml](#build) (job: `build`) - entry point

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (1)</summary>

1. **Build**

</details>

[Back to contents](#contents)

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

[Back to contents](#contents)

