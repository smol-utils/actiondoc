# Contents

- [CodeQL](#codeql)
- [Release](#release)
- [Validate GitHub Actions](#validate-github-actions)
- [Validations](#validations)
- [Bootstrap](#bootstrap)

# CodeQL

| Property | Value |
|----------|-------|
| File | `codeql.yaml` |
| Triggers | `push`, `pull_request`, `schedule` |

## Schedule

- `38 11 * * 3`

## Event filters

- **push**
  - branches: `main`
- **pull_request**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Call graph (rooted at this workflow)

```
codeql.yaml [push, pull_request, schedule]
+-- analyze (uses anchore/workflows/.github/workflows/codeql.yaml@15122524ced7906bfa9685eeae12e22647773ea6)
```

## Transitive requirements (from full call graph)

Permissions declared across the chain: `actions: read`, `contents: read`, `packages: read`, `security-events: write`

External workflows referenced: `anchore/workflows/.github/workflows/codeql.yaml@15122524ced7906bfa9685eeae12e22647773ea6`

## Jobs

### Analyze (`analyze`)

| Property | Value |
|----------|-------|
| Uses workflow | `anchore/workflows/.github/workflows/codeql.yaml@15122524ced7906bfa9685eeae12e22647773ea6` (external) |

**Permissions:**

- `security-events`: `write`
- `packages`: `read`
- `actions`: `read`
- `contents`: `read`

# Release

| Property | Value |
|----------|-------|
| File | `release.yaml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | - | Yes | - | tag the latest commit on main with the given version (prefixed with v) |
| `phase` | choice | Yes | `all` | the specific workflow phase to run or all<br>Options: `all`, `install-script-only` |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `release`, cancel-in-progress: `false`

## Call graph (rooted at this workflow)

```
release.yaml [workflow_dispatch]
+-- version-available (uses anchore/workflows/.github/workflows/check-version-available.yaml@15122524ced7906bfa9685eeae12e22647773ea6)
+-- check-gate (uses anchore/workflows/.github/workflows/check-gate.yaml@15122524ced7906bfa9685eeae12e22647773ea6)
+-- release / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- release-install-script (uses anchore/workflows/.github/workflows/release-install-script.yaml@15122524ced7906bfa9685eeae12e22647773ea6)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ANCHOREOPS_GITHUB_OSS_WRITE_TOKEN`, `ANCHOREOSSWRITE_DH_PAT`, `ANCHOREOSSWRITE_DH_USERNAME`, `ANCHORE_APPLE_DEVELOPER_ID_CERT_CHAIN`, `ANCHORE_APPLE_DEVELOPER_ID_CERT_PASS`, `APPLE_NOTARY_ISSUER`, `APPLE_NOTARY_KEY`, `APPLE_NOTARY_KEY_ID`, `DEPLOY_KEY`, `GITHUB_TOKEN`, `OSS_R2_INSTALL_ACCESS_KEY_ID`, `OSS_R2_INSTALL_SECRET_ACCESS_KEY`, `R2_ENDPOINT`, `R2_INSTALL_ACCESS_KEY_ID`, `R2_INSTALL_SECRET_ACCESS_KEY`, `S3_INSTALL_AWS_ACCESS_KEY_ID`, `S3_INSTALL_AWS_SECRET_ACCESS_KEY`, `TOOLBOX_AWS_ACCESS_KEY_ID`, `TOOLBOX_AWS_SECRET_ACCESS_KEY`, `TOOLBOX_CLOUDFLARE_R2_ENDPOINT`

Permissions declared across the chain: `checks: read`, `contents: read`, `contents: write`, `id-token: write (OIDC)`, `packages: write`

External workflows referenced: `anchore/workflows/.github/workflows/check-gate.yaml@15122524ced7906bfa9685eeae12e22647773ea6`, `anchore/workflows/.github/workflows/check-version-available.yaml@15122524ced7906bfa9685eeae12e22647773ea6`, `anchore/workflows/.github/workflows/release-install-script.yaml@15122524ced7906bfa9685eeae12e22647773ea6`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `ANCHOREOSSWRITE_DH_USERNAME` | job `release` step `Login to Docker Hub` with `username` |
| `ANCHOREOSSWRITE_DH_PAT` | job `release` step `Login to Docker Hub` with `password` |
| `GITHUB_TOKEN` | job `release` step `Login to GitHub Container Registry` with `password`; job `release` step `Build & publish release artifacts` env `GITHUB_TOKEN` |
| `DEPLOY_KEY` | job `release` step `Build & publish release artifacts` env `DEPLOY_KEY` |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_CHAIN` | job `release` step `Build & publish release artifacts` env `QUILL_SIGN_P12` |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_PASS` | job `release` step `Build & publish release artifacts` env `QUILL_SIGN_PASSWORD` |
| `APPLE_NOTARY_ISSUER` | job `release` step `Build & publish release artifacts` env `QUILL_NOTARY_ISSUER` |
| `APPLE_NOTARY_KEY_ID` | job `release` step `Build & publish release artifacts` env `QUILL_NOTARY_KEY_ID` |
| `APPLE_NOTARY_KEY` | job `release` step `Build & publish release artifacts` env `QUILL_NOTARY_KEY` |
| `ANCHOREOPS_GITHUB_OSS_WRITE_TOKEN` | job `release` step `Build & publish release artifacts` env `GITHUB_BREW_TOKEN` |
| `OSS_R2_INSTALL_ACCESS_KEY_ID` | job `release-install-script` secrets `R2_INSTALL_ACCESS_KEY_ID` |
| `OSS_R2_INSTALL_SECRET_ACCESS_KEY` | job `release-install-script` secrets `R2_INSTALL_SECRET_ACCESS_KEY` |
| `TOOLBOX_CLOUDFLARE_R2_ENDPOINT` | job `release-install-script` secrets `R2_ENDPOINT` |
| `TOOLBOX_AWS_ACCESS_KEY_ID` | job `release-install-script` secrets `S3_INSTALL_AWS_ACCESS_KEY_ID` |
| `TOOLBOX_AWS_SECRET_ACCESS_KEY` | job `release-install-script` secrets `S3_INSTALL_AWS_SECRET_ACCESS_KEY` |

## Jobs

### `version-available`

| Property | Value |
|----------|-------|
| Uses workflow | `anchore/workflows/.github/workflows/check-version-available.yaml@15122524ced7906bfa9685eeae12e22647773ea6` (external) |
| Condition | `${{ github.event.inputs.phase == 'all' }}` |

**Permissions:**

- `contents`: `read` - required for fetching tags

#### Inputs forwarded

- `version`: `${{ github.event.inputs.version }}`

### `check-gate`

| Property | Value |
|----------|-------|
| Uses workflow | `anchore/workflows/.github/workflows/check-gate.yaml@15122524ced7906bfa9685eeae12e22647773ea6` (external) |
| Condition | `${{ github.event.inputs.phase == 'all' }}` |

**Permissions:**

- `checks`: `read` - required for getting the status of specific check names

#### Inputs forwarded

- `checks`: `["Acceptance tests (Linux)", "Acceptance tests (Mac)", "Build snapshot artifacts", "CLI tests (Linux)", "Integration tests", "Static analysis", "Unit tests"]`

### `release`

| Property | Value |
|----------|-------|
| Runs on | `runs-on=${{ github.run_id }}/cpu=16+32/ram=32+128/family=c5+c6+c7+c8/spot=false/extras=s3-cache+tmpfs` |
| Depends on | `check-gate`, `version-available` |
| Condition | `${{ github.event.inputs.phase == 'all' }}` |

**Deploys to environment:** `release` [gated]

> Environment protection rules (required reviewers, wait timers, branch policies) are configured in the repository's Settings -> Environments and are not represented here.

**Permissions:**

- `contents`: `write` - required for creating the GitHub release and pushing the version tag
- `packages`: `write` - required for publishing release artifacts to GitHub packages
- `id-token`: `write` (OIDC) - required for keyless signing (cosign/sigstore OIDC)

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `true`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`

3. **Login to Docker Hub**
   - Uses: `docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121` (v4.1.0)
   - With:
     - `username`: `${{ secrets.ANCHOREOSSWRITE_DH_USERNAME }}`
     - `password`: `${{ secrets.ANCHOREOSSWRITE_DH_PAT }}`

4. **Login to GitHub Container Registry**
   - Uses: `docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121` (v4.1.0)
   - With:
     - `registry`: `ghcr.io`
     - `username`: `${{ github.actor }}`
     - `password`: `${{ secrets.GITHUB_TOKEN }}`

5. **Build & publish release artifacts**
   - Env:
     - `DEPLOY_KEY`: `${{ secrets.DEPLOY_KEY }}`
     - `RELEASE_VERSION`: `${{ github.event.inputs.version }}`
     - `QUILL_SIGN_P12`: `${{ secrets.ANCHORE_APPLE_DEVELOPER_ID_CERT_CHAIN }}`
     - `QUILL_SIGN_PASSWORD`: `${{ secrets.ANCHORE_APPLE_DEVELOPER_ID_CERT_PASS }}`
     - `QUILL_NOTARY_ISSUER`: `${{ secrets.APPLE_NOTARY_ISSUER }}`
     - `QUILL_NOTARY_KEY_ID`: `${{ secrets.APPLE_NOTARY_KEY_ID }}`
     - `QUILL_NOTARY_KEY`: `${{ secrets.APPLE_NOTARY_KEY }}`
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_BREW_TOKEN`: `${{ secrets.ANCHOREOPS_GITHUB_OSS_WRITE_TOKEN }}`

6. **anchore/sbom-action@v0.24.0** `[continue-on-error]`
   - Uses: `anchore/sbom-action@e22c389904149dbc22b58101806040fa8d37a610` (v0.24.0)
   - With:
     - `file`: `go.mod`
     - `artifact-name`: `sbom.spdx.json`

### `release-install-script`

| Property | Value |
|----------|-------|
| Uses workflow | `anchore/workflows/.github/workflows/release-install-script.yaml@15122524ced7906bfa9685eeae12e22647773ea6` (external) |
| Depends on | `release` |
| Condition | `${{ always() && (needs.release.result == 'success' \|\| github.event.inputs.phase == 'install-script-only') }}` |

**Permissions:**

- `contents`: `read` - required for the reusable workflow to check out the repo and publish the install script

#### Inputs forwarded

- `tag`: `${{ github.event.inputs.version }}`

#### Secrets forwarded

- `R2_INSTALL_ACCESS_KEY_ID`: `${{ secrets.OSS_R2_INSTALL_ACCESS_KEY_ID }}`
- `R2_INSTALL_SECRET_ACCESS_KEY`: `${{ secrets.OSS_R2_INSTALL_SECRET_ACCESS_KEY }}`
- `R2_ENDPOINT`: `${{ secrets.TOOLBOX_CLOUDFLARE_R2_ENDPOINT }}`
- `S3_INSTALL_AWS_ACCESS_KEY_ID`: `${{ secrets.TOOLBOX_AWS_ACCESS_KEY_ID }}`
- `S3_INSTALL_AWS_SECRET_ACCESS_KEY`: `${{ secrets.TOOLBOX_AWS_SECRET_ACCESS_KEY }}`

# Validate GitHub Actions

| Property | Value |
|----------|-------|
| File | `validate-github-actions.yaml` |
| Triggers | `workflow_dispatch`, `pull_request`, `push` |

## Event filters

- **push**
  - branches: `main`
  - paths: `.github/workflows/**`, `.github/actions/**`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### Lint (`zizmor`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`
- `security-events`: `write` - for uploading SARIF results

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Run zizmor**
   - Uses: `zizmorcore/zizmor-action@a16621b09c6db4281f81a93cb393b05dcd7b7165` (v0.5.5)
   - With:
     - `advanced-security`: `true`
     - `inputs`: `.github`

# Validations

| Property | Value |
|----------|-------|
| File | `validations.yaml` |
| Triggers | `workflow_dispatch`, `pull_request`, `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
validations.yaml [workflow_dispatch, pull_request, push]
+-- Static-Analysis / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Unit-Test / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Integration-Test / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Build-Snapshot-Artifacts / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Acceptance-Linux / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Acceptance-Mac / Bootstrap environment (uses ./.github/actions/bootstrap)
+-- Cli-Linux / Bootstrap environment (uses ./.github/actions/bootstrap)
```

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

## Jobs

### Static analysis (`Static-Analysis`)

| Property | Value |
|----------|-------|
| Runs on | `&test-runner "runs-on=${{ github.run_id }}/cpu=4+8/ram=32+128/family=r5+r6+r7+r8+m4+m5+m6+m7+m8/spot=price-capacity-optimized/extras=tmpfs"` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Run static analysis**

### Unit tests (`Unit-Test`)

| Property | Value |
|----------|-------|
| Runs on | `*test-runner` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Run unit tests**

4. **Check for capability drift**

### Integration tests (`Integration-Test`)

| Property | Value |
|----------|-------|
| Runs on | `*test-runner` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Validate syft output against the CycloneDX schema**

4. **Run integration tests**

### Build snapshot artifacts (`Build-Snapshot-Artifacts`)

| Property | Value |
|----------|-------|
| Runs on | `runs-on=${{ github.run_id }}/cpu=16+32/ram=32+128/family=c5+c6+c7+c8/spot=false/extras=tmpfs` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `bootstrap-apt-packages`: - - Space delimited list of tools to install via apt

3. **Build snapshot artifacts**

4. **Smoke test snapshot build**

5. **Upload snapshot artifacts**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot/`
     - `retention-days`: `30`

### Acceptance tests (Linux) (`Acceptance-Linux`)

| Property | Value |
|----------|-------|
| Runs on | `*test-runner` |
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot`

4. **Restore binary permissions**

5. **Run comparison tests (Linux)**

6. **Load test image cache**
   - Condition: `steps.install-test-image-cache.outputs.cache-hit == 'true'`

7. **Run install.sh tests (Linux)**

8. **(cache-miss) Create test image cache**
   - Condition: `steps.install-test-image-cache.outputs.cache-hit != 'true'`

### Acceptance tests (Mac) (`Acceptance-Mac`)

| Property | Value |
|----------|-------|
| Runs on | `macos-latest` |
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **Install Cosign**
   - Uses: `sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6` (v4.1.2)

2. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `bootstrap-apt-packages`: - - Space delimited list of tools to install via apt
     - `go-dependencies`: `false`
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

4. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot`

5. **Restore binary permissions**

6. **Run comparison tests (Mac)**

7. **Run install.sh tests (Mac)**

### CLI tests (Linux) (`Cli-Linux`)

| Property | Value |
|----------|-------|
| Runs on | `*test-runner` |
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot`

4. **Restore binary permissions**

5. **Run CLI Tests (Linux)**

# Bootstrap

Bootstrap all syft tools and dependencies on top of go-make's setup action

| Property | Value |
|----------|-------|
| File | `action.yaml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `go-version` | Go version to install (passed to go-make/setup) | Yes | `1.26.2` |
| `cache-key-prefix` | Prefix all cache keys with this value (passed to go-make/setup) | Yes | `v1` |
| `cache-enabled` | Enable build/mod and tool caching (passed to go-make/setup) | Yes | `true` |
| `download-test-fixture-cache` | Download test fixture cache from OCI and github actions | Yes | `false` |
| `bootstrap-apt-packages` | Space delimited list of tools to install via apt | No | `libxml2-utils` |

