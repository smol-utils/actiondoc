# syft

4 workflows, 1 composite action

## Contents

**Workflows**

- [CodeQL - push, pull_request, schedule](#codeql)
- [Release - workflow_dispatch](#release)
- [Validate GitHub Actions - workflow_dispatch, pull_request, push](#validate-github-actions)
- [Validations - workflow_dispatch, pull_request, push](#validations)

**Composite actions**

- [Bootstrap](#bootstrap)

## Workflows by trigger

- **pull_request**: [CodeQL](#codeql), [Validate GitHub Actions](#validate-github-actions), [Validations](#validations)
- **push**: [CodeQL](#codeql), [Validate GitHub Actions](#validate-github-actions), [Validations](#validations)
- **workflow_dispatch**: [Release](#release), [Validate GitHub Actions](#validate-github-actions), [Validations](#validations)
- **schedule**: [CodeQL](#codeql)

## Secrets and variables used across this repository

**Secrets:**

| Name | Used by |
|------|---------|
| `ANCHOREOPS_GITHUB_OSS_WRITE_TOKEN` | [Release](#release) |
| `ANCHOREOSSWRITE_DH_PAT` | [Release](#release) |
| `ANCHOREOSSWRITE_DH_USERNAME` | [Release](#release) |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_CHAIN` | [Release](#release) |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_PASS` | [Release](#release) |
| `APPLE_NOTARY_ISSUER` | [Release](#release) |
| `APPLE_NOTARY_KEY` | [Release](#release) |
| `APPLE_NOTARY_KEY_ID` | [Release](#release) |
| `DEPLOY_KEY` | [Release](#release) |
| `OSS_R2_INSTALL_ACCESS_KEY_ID` | [Release](#release) |
| `OSS_R2_INSTALL_SECRET_ACCESS_KEY` | [Release](#release) |
| `TOOLBOX_AWS_ACCESS_KEY_ID` | [Release](#release) |
| `TOOLBOX_AWS_SECRET_ACCESS_KEY` | [Release](#release) |
| `TOOLBOX_CLOUDFLARE_R2_ENDPOINT` | [Release](#release) |

## Permissions across this repository

| Scope | Level |
|-------|-------|
| `actions` | `read` |
| `checks` | `read` |
| `contents` | `write` (also granted as `read` elsewhere) |
| `id-token` | `write` (OIDC) |
| `packages` | `write` (also granted as `read` elsewhere) |
| `security-events` | `write` |

# CodeQL

**Triggers:** `push`, `pull_request`, `schedule`

| Property | Value |
|----------|-------|
| File | `codeql.yaml` |

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

- `analyze` uses `anchore/workflows/.github/workflows/codeql.yaml@15122524ced7906bfa9685eeae12e22647773ea6`

## Transitive requirements (from full call graph)

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

[Back to contents](#contents)

# Release

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `release.yaml` |

**Jobs:** [`version-available`](#version-available), [`check-gate`](#check-gate), [`release`](#release), [`release-install-script`](#release-install-script)

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

- `version-available` uses `anchore/workflows/.github/workflows/check-version-available.yaml@15122524ced7906bfa9685eeae12e22647773ea6`
- `check-gate` uses `anchore/workflows/.github/workflows/check-gate.yaml@15122524ced7906bfa9685eeae12e22647773ea6`
- `release / Bootstrap environment` uses [./.github/actions/bootstrap](#bootstrap)
- `release-install-script` uses `anchore/workflows/.github/workflows/release-install-script.yaml@15122524ced7906bfa9685eeae12e22647773ea6`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `R2_ENDPOINT`, `R2_INSTALL_ACCESS_KEY_ID`, `R2_INSTALL_SECRET_ACCESS_KEY`, `S3_INSTALL_AWS_ACCESS_KEY_ID`, `S3_INSTALL_AWS_SECRET_ACCESS_KEY`

External workflows referenced: `anchore/workflows/.github/workflows/check-gate.yaml@15122524ced7906bfa9685eeae12e22647773ea6`, `anchore/workflows/.github/workflows/check-version-available.yaml@15122524ced7906bfa9685eeae12e22647773ea6`, `anchore/workflows/.github/workflows/release-install-script.yaml@15122524ced7906bfa9685eeae12e22647773ea6`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `ANCHOREOSSWRITE_DH_USERNAME` | `release`: Login to Docker Hub (`username`) |
| `ANCHOREOSSWRITE_DH_PAT` | `release`: Login to Docker Hub (`password`) |
| `GITHUB_TOKEN` | `release`: Login to GitHub Container Registry (`password`), Build & publish release artifacts (`GITHUB_TOKEN`) |
| `DEPLOY_KEY` | `release`: Build & publish release artifacts (`DEPLOY_KEY`) |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_CHAIN` | `release`: Build & publish release artifacts (`QUILL_SIGN_P12`) |
| `ANCHORE_APPLE_DEVELOPER_ID_CERT_PASS` | `release`: Build & publish release artifacts (`QUILL_SIGN_PASSWORD`) |
| `APPLE_NOTARY_ISSUER` | `release`: Build & publish release artifacts (`QUILL_NOTARY_ISSUER`) |
| `APPLE_NOTARY_KEY_ID` | `release`: Build & publish release artifacts (`QUILL_NOTARY_KEY_ID`) |
| `APPLE_NOTARY_KEY` | `release`: Build & publish release artifacts (`QUILL_NOTARY_KEY`) |
| `ANCHOREOPS_GITHUB_OSS_WRITE_TOKEN` | `release`: Build & publish release artifacts (`GITHUB_BREW_TOKEN`) |
| `OSS_R2_INSTALL_ACCESS_KEY_ID` | `release-install-script`: (`R2_INSTALL_ACCESS_KEY_ID`) |
| `OSS_R2_INSTALL_SECRET_ACCESS_KEY` | `release-install-script`: (`R2_INSTALL_SECRET_ACCESS_KEY`) |
| `TOOLBOX_CLOUDFLARE_R2_ENDPOINT` | `release-install-script`: (`R2_ENDPOINT`) |
| `TOOLBOX_AWS_ACCESS_KEY_ID` | `release-install-script`: (`S3_INSTALL_AWS_ACCESS_KEY_ID`) |
| `TOOLBOX_AWS_SECRET_ACCESS_KEY` | `release-install-script`: (`S3_INSTALL_AWS_SECRET_ACCESS_KEY`) |

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

**Permissions:**

- `contents`: `write` - required for creating the GitHub release and pushing the version tag
- `packages`: `write` - required for publishing release artifacts to GitHub packages
- `id-token`: `write` (OIDC) - required for keyless signing (cosign/sigstore OIDC)

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `true`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`

3. **Login to Docker Hub**
   - Uses: `docker/login-action@v4.1.0`
   - With:
     - `username`: `${{ secrets.ANCHOREOSSWRITE_DH_USERNAME }}`
     - `password`: `${{ secrets.ANCHOREOSSWRITE_DH_PAT }}`

4. **Login to GitHub Container Registry**
   - Uses: `docker/login-action@v4.1.0`
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
   - With:
     - `file`: `go.mod`
     - `artifact-name`: `sbom.spdx.json`

</details>

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

[Back to contents](#contents)

# Validate GitHub Actions

**Triggers:** `workflow_dispatch`, `pull_request`, `push`

| Property | Value |
|----------|-------|
| File | `validate-github-actions.yaml` |

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

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Run zizmor**
   - Uses: `zizmorcore/zizmor-action@v0.5.5`
   - With:
     - `advanced-security`: `true`
     - `inputs`: `.github`

</details>

[Back to contents](#contents)

# Validations

**Triggers:** `workflow_dispatch`, `pull_request`, `push`

| Property | Value |
|----------|-------|
| File | `validations.yaml` |
| Default runs-on | `runs-on=${{ github.run_id }}/cpu=4+8/ram=32+128/family=r5+r6+r7+r8+m4+m5+m6+m7+m8/spot=price-capacity-optimized/extras=tmpfs` |

**Jobs:** [Static analysis](#static-analysis-static-analysis), [Unit tests](#unit-tests-unit-test), [Integration tests](#integration-tests-integration-test), [Build snapshot artifacts](#build-snapshot-artifacts-build-snapshot-artifacts), [Acceptance tests (Linux)](#acceptance-tests-linux-acceptance-linux), [Acceptance tests (Mac)](#acceptance-tests-mac-acceptance-mac), [CLI tests (Linux)](#cli-tests-linux-cli-linux)

## Event filters

- **push**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

- uses **[./.github/actions/bootstrap](#bootstrap)** (x7)

## Jobs

### Static analysis (`Static-Analysis`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Run static analysis**

</details>

### Unit tests (`Unit-Test`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Run unit tests**

4. **Check for capability drift**

</details>

### Integration tests (`Integration-Test`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Validate syft output against the CycloneDX schema**

4. **Run integration tests**

</details>

### Build snapshot artifacts (`Build-Snapshot-Artifacts`)

| Property | Value |
|----------|-------|
| Runs on | `runs-on=${{ github.run_id }}/cpu=16+32/ram=32+128/family=c5+c6+c7+c8/spot=false/extras=tmpfs` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`

3. **Build snapshot artifacts**

4. **Smoke test snapshot build**

5. **Upload snapshot artifacts**
   - Uses: `actions/upload-artifact@v7.0.1`
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot/`
     - `retention-days`: `30`

</details>

### Acceptance tests (Linux) (`Acceptance-Linux`)

| Property | Value |
|----------|-------|
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (8)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
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

</details>

### Acceptance tests (Mac) (`Acceptance-Mac`)

| Property | Value |
|----------|-------|
| Runs on | `macos-latest` |
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (7)</summary>

1. **Install Cosign**
   - Uses: `sigstore/cosign-installer@v4.1.2`

2. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

3. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `go-dependencies`: `false`
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

4. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot`

5. **Restore binary permissions**

6. **Run comparison tests (Mac)**

7. **Run install.sh tests (Mac)**

</details>

### CLI tests (Linux) (`Cli-Linux`)

| Property | Value |
|----------|-------|
| Depends on | `Build-Snapshot-Artifacts` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Bootstrap environment**
   - Uses: `./.github/actions/bootstrap`
   - With:
     - `download-test-fixture-cache`: `true` - Download test fixture cache from OCI and github actions (required)

3. **Download snapshot artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `snapshot`
     - `path`: `snapshot`

4. **Restore binary permissions**

5. **Run CLI Tests (Linux)**

</details>

[Back to contents](#contents)

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

[Back to contents](#contents)

