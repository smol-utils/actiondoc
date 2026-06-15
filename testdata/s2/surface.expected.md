# Release

**Triggers:** `workflow_dispatch`, `push`, `pull_request`, `repository_dispatch`, `schedule`

Release pipeline: builds artifacts and deploys to the production environment. Triggered manually, on a weekly schedule, and on tagged pushes.

| Property | Value |
|----------|-------|
| File | `surface.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [Build](#build-build), [Deploy](#deploy-deploy)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `distribution` | choice | Yes | - | Which build to release.<br>Options: `mandrel`, `graalvm-community`, `graalvm`, `liberica` |
| `dry-run` | boolean | No | `false` | Skip the publish step. |
| `retries` | number | No | `3` | - |

## Schedule

- `0 16 * * THU` - weekly security scan
- `30 2 * * *`

## Event filters

- **push**
  - branches: `main`, `release/*`
  - paths: `src/**`, `!docs/**`
- **pull_request**
  - types: `opened`, `labeled`, `synchronize`
- **repository_dispatch**
  - types: `deploy-request`

## Permissions

- `contents`: `read` - for actions/checkout to fetch code
- `id-token`: `write` (OIDC) - OIDC token for keyless signing
- `packages`: `write`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `LANG` | `en_US.UTF-8` |
| `GRADLE_OPTS` | `-Dorg.gradle.daemon=false` |

**Concurrency:** group `release-${{ github.ref }}`, cancel-in-progress: `true`

**Defaults:** shell `bash`, working-directory `./app`

## Jobs

### Build (`build`)

**Permissions:**

No permissions granted (`permissions: {}` -- default-deny).

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BUILD_MODE` | `release` |

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4**

2. **make build**

</details>

### Deploy (`deploy`)

| Property | Value |
|----------|-------|
| Depends on | `build` |

**Deploys to environment:** `production` (`https://app.example.com`) [gated]

**Permissions:**

- `contents`: `read`
- `id-token`: `write` (OIDC)

**Concurrency:** group `deploy-prod`, cancel-in-progress: `false`

**Defaults:** working-directory `./deploy`

<details>
<summary>Steps (1)</summary>

1. **./publish.sh**

</details>

