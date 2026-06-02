# Contents

- [CI-Container-Build](#ci-container-build)
- [CodeQL](#codeql)
- [Conformance Tests Nightly](#conformance-tests-nightly)
- [Conformance Tests](#conformance-tests)
- [Cut Release](#cut-release)
- [Dependency Review](#dependency-review)
- [Do Not Submit](#do-not-submit)
- [e2e-tests](#e2e-tests)
- [e2e-with-binary](#e2e-with-binary)
- [Test GitHub OIDC](#test-github-oidc)
- [golangci-lint](#golangci-lint)
- [Test attest / verify-attestation](#test-attest--verify-attestation)
- [Scorecards supply-chain security](#scorecards-supply-chain-security)
- [CI-Tests](#ci-tests)
- [CI-Validate-Release-Job](#ci-validate-release-job)
- [Docgen](#docgen)
- [Whitespace](#whitespace)

# CI-Container-Build

| Property | Value |
|----------|-------|
| File | `build.yaml` |
| Triggers | `push` |

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `build` step `Login to GitHub Container Registry` with `password` |
| `COSIGN_PASSWORD` | job `build` step `containers-cosign` env `COSIGN_PASSWORD` |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`
- `packages`: `write`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **sigstore/cosign-installer@v4.1.2**
   - Uses: `sigstore/cosign-installer@6f9f17788090df1f26f669e9d70d6ae9567deba6` (v4.1.2)

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **ko-build/setup-ko@v0.9**
   - Uses: `ko-build/setup-ko@d006021bd0c28d1ce33a07e7943d48b079944c8d` (v0.9)

6. **Set up Cloud SDK**
   - Uses: `google-github-actions/auth@7c6bc770dae815cd3e89ee6cdf493a5fab2cc093` (v3.0.0)
   - With:
     - `workload_identity_provider`: `projects/498091336538/locations/global/workloadIdentityPools/githubactions/providers/sigstore-cosign`
     - `service_account`: `github-actions@projectsigstore.iam.gserviceaccount.com`

7. **creds**

8. **Login to GitHub Container Registry**
   - Uses: `docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121` (v4.1.0)
   - With:
     - `registry`: `ghcr.io`
     - `username`: `${{ github.actor }}`
     - `password`: `${{ secrets.GITHUB_TOKEN }}`

9. **containers-cosign**
   - Env:
     - `KO_PREFIX`: `ghcr.io/sigstore/cosign/cosign/ci`
     - `COSIGN_PASSWORD`: `${{secrets.COSIGN_PASSWORD}}`

# CodeQL

| Property | Value |
|----------|-------|
| File | `codeql-analysis.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `CODEQL_EXTRACTOR_GO_BUILD_TRACING` | `true` |

## Jobs

### Analyze (`analyze`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `security-events`: `write`
- `actions`: `read`
- `contents`: `read`

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Utilize Go Module Cache**
   - Uses: `actions/cache@27d5ce7f107fe9357f9df03efb73ab90386fccae` (v5.0.5)
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@65c74964a9ed8c44ed9f19d4bbc5757a6a8e9ab9` (v2.16.1)
   - With:
     - `languages`: `${{ matrix.language }}`

6. **Build cosign for CodeQL**

7. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@65c74964a9ed8c44ed9f19d4bbc5757a6a8e9ab9` (v2.16.1)

# Conformance Tests Nightly

| Property | Value |
|----------|-------|
| File | `conformance-nightly.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 0 * * *` - 12:00 AM UTC

## Permissions

- `contents`: `read`
- `issues`: `write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `conformance` step `Create Issue on Failure` with `github-token` |

## Jobs

### `conformance`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **make cosign conformance**

5. **sigstore/sigstore-conformance@main**
   - Uses: `sigstore/sigstore-conformance@main`
   - With:
     - `entrypoint`: `${{ github.workspace }}/conformance`
     - `xfail`: `test_verify*PATH-message-digest-mismatch_fail]`

6. **Create Issue on Failure**
   - Uses: `actions/github-script@ed597411d8f924073f98dfc5c65a23a2325f34cd` (v8.0.0)
   - Condition: `failure()`
   - With:
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
     - `script`: `const { owner, repo } = context.repo; const runId = context.runId; const issueTitle = 'Conformance Tests Failed'; const issueBody = `The nightly conformance tests have failed. Please check the logs for more details.\n\nWorkflow run: https://github.com/${owner}/${repo}/actions/runs/${runId}\n\ncc @sigstore/security-response-team @sigstore/cosign-codeowners`; const issueLabel = 'bug';  const existingIssues = await github.rest.issues.listForRepo({   owner,   repo,   state: 'open',   labels: issueLabel, });  const issueExists = existingIssues.data.some(issue => issue.title === issueTitle);  if (!issueExists) {   await github.rest.issues.create({     owner,     repo,     title: issueTitle,     body: issueBody,     labels: [issueLabel],   }); }`

# Conformance Tests

| Property | Value |
|----------|-------|
| File | `conformance.yml` |
| Triggers | `push`, `pull_request` |

## Event filters

- **push**
  - branches: `main`
- **pull_request**
  - branches: `main`

## Permissions

- `contents`: `read`

## Jobs

### `conformance`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **make cosign conformance**

5. **sigstore/sigstore-conformance@v0.0.27**
   - Uses: `sigstore/sigstore-conformance@4d66ba3cb0c9c95f705c757c0f5e226d3f4d5151` (v0.0.27)
   - With:
     - `entrypoint`: `${{ github.workspace }}/conformance`
     - `xfail`: `test_verify*PATH-message-digest-mismatch_fail]`

# Cut Release

| Property | Value |
|----------|-------|
| File | `cut-release.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `release_tag` | string | Yes | - | Release tag |
| `key_ring` | string | Yes | - | Key ring for cosign key |
| `key_name` | string | Yes | - | Key name for cosign key |

**Concurrency:** group `cut-release`

## Call graph (rooted at this workflow)

```
cut-release.yml [workflow_dispatch]
+-- cut-release (uses sigstore/community/.github/workflows/reusable-release.yml@main)
```

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`, `id-token: write (OIDC)`

External workflows referenced: `sigstore/community/.github/workflows/reusable-release.yml@main`

## Jobs

### Cut release (`cut-release`)

| Property | Value |
|----------|-------|
| Uses workflow | `sigstore/community/.github/workflows/reusable-release.yml@main` (external) |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

#### Inputs forwarded

- `release_tag`: `${{ github.event.inputs.release_tag }}`
- `key_ring`: `${{ github.event.inputs.key_ring }}`
- `key_name`: `${{ github.event.inputs.key_name }}`
- `workload_identity_provider`: `projects/498091336538/locations/global/workloadIdentityPools/githubactions/providers/sigstore-cosign`
- `service_account`: `github-actions-cosign@projectsigstore.iam.gserviceaccount.com`
- `repo`: `cosign`

# Dependency Review

| Property | Value |
|----------|-------|
| File | `depsreview.yml` |
| Triggers | `pull_request` |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Call graph (rooted at this workflow)

```
depsreview.yml [pull_request]
+-- dependency-review (uses sigstore/community/.github/workflows/reusable-dependency-review.yml@main)
```

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

External workflows referenced: `sigstore/community/.github/workflows/reusable-dependency-review.yml@main`

## Jobs

### License and Vulnerability Scan (`dependency-review`)

| Property | Value |
|----------|-------|
| Uses workflow | `sigstore/community/.github/workflows/reusable-dependency-review.yml@main` (external) |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `contents`: `read`

# Do Not Submit

| Property | Value |
|----------|-------|
| File | `donotsubmit.yaml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### Do Not Submit (`donotsubmit`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **Check out code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v2.4.0)
   - With:
     - `persist-credentials`: `false`

2. **Do Not Submit**
   - Uses: `chainguard-dev/actions/donotsubmit@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)

# e2e-tests

| Property | Value |
|----------|-------|
| File | `e2e-tests.yml` |
| Triggers | `push`, `pull_request`, `workflow_dispatch` |

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`

## Jobs

### `e2e-cross`

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Run cross platform e2e tests**

### `e2e-test-pkcs11`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Run pkcs11 end-to-end tests**

### `e2e-kms`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VAULT_TOKEN` | `root` |
| `VAULT_ADDR` | `http://localhost:8200` |
| `COSIGN_YES` | `true` |
| `SCAFFOLDING_RELEASE_VERSION` | `v0.7.24` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **setup vault**
   - Uses: `cpanato/vault-installer@fe568170412f5d81202ec528148f05176efbecc1` (v1.4.0)

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **imjasonh/setup-crane@v0.5**
   - Uses: `imjasonh/setup-crane@6da1ae018866400525525ce74ff892880c099987` (v0.5)

6. **Install cluster + sigstore**
   - Uses: `sigstore/scaffolding/actions/setup@main`
   - With:
     - `version`: `${{ env.SCAFFOLDING_RELEASE_VERSION }}`

7. **enable vault transit**

8. **Acceptance Tests**

### `e2e-registry`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `SCAFFOLDING_RELEASE_VERSION` | `v0.7.24` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Setup mirror**
   - Uses: `chainguard-dev/actions/setup-mirror@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - With:
     - `mirror`: `mirror.gcr.io`

5. **Install cluster + sigstore**
   - Uses: `sigstore/scaffolding/actions/setup@main`
   - With:
     - `version`: `${{ env.SCAFFOLDING_RELEASE_VERSION }}`

6. **Setup local insecure registry**
   - Env:
     - `INSECURE_REGISTRY_NAME`: `insecure-registry.notlocal`
     - `INSECURE_REGISTRY_PORT`: `5001`

7. **Run Insecure Registry Tests**
   - Env:
     - `COSIGN_TEST_REPO`: `insecure-registry.notlocal:5001`
     - `TUF_ROOT_JSON`: `${{ github.workspace }}/root.json`

8. **Setup local insecure OCI 1.1 registry**
   - Env:
     - `ZOT_VERSION`: `v2.0.0-rc6`
     - `INSECURE_OCI_REGISTRY_NAME`: `insecure-oci-registry.notlocal`
     - `INSECURE_OCI_REGISTRY_PORT`: `5002`

9. **Run Insecure OCI 1.1 Registry Tests**
   - Env:
     - `OCI11`: `yes`
     - `COSIGN_TEST_REPO`: `insecure-oci-registry.notlocal:5002`
     - `TUF_ROOT_JSON`: `${{ github.workspace }}/root.json`

10. **Set up local HTTP registry**
   - Env:
     - `HTTP_REGISTRY_NAME`: `http-registry.notlocal`
     - `HTTP_REGISTRY_PORT`: `5003`

11. **Run HTTP registry tests**
   - Env:
     - `COSIGN_TEST_REPO`: `http-registry.notlocal:5003`
     - `TUF_ROOT_JSON`: `${{ github.workspace }}/root.json`

12. **Collect diagnostics**
   - Uses: `chainguard-dev/actions/kind-diag@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - Condition: `${{ failure() }}`

# e2e-with-binary

| Property | Value |
|----------|-------|
| File | `e2e-with-binary.yml` |
| Triggers | `push`, `workflow_dispatch` |

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### Run tests (`e2e-tests-with-binary`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Condition | `${{ github.repository == 'sigstore/cosign' }}` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COSIGN_YES` | `true` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **build cosign and check sign-blob and verify-blob**

# Test GitHub OIDC

| Property | Value |
|----------|-------|
| File | `github-oidc.yaml` |
| Triggers | `push`, `schedule`, `workflow_dispatch` |

## Schedule

- `0 1 * * *` - 1AM UTC

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `packages`: `write`
- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GIT_HASH` | `${{ github.sha }}` |
| `GIT_VERSION` | `unstable` |
| `GITHUB_RUN_ID` | `${{ github.run_id }}` |
| `GITHUB_RUN_ATTEMPT` | `${{ github.run_attempt }}` |
| `KO_PREFIX` | `ghcr.io/${{ github.repository }}` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **ko-build/setup-ko@v0.9**
   - Uses: `ko-build/setup-ko@d006021bd0c28d1ce33a07e7943d48b079944c8d` (v0.9)

5. **build cosign from the HEAD**

6. **Build and sign a container image**

# golangci-lint

| Property | Value |
|----------|-------|
| File | `golangci-lint.yml` |
| Triggers | `push`, `pull_request` |

## Event filters

- **push**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### lint (`golangci`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20` (v9.2.0)
   - With:
     - `version`: `v2.12`

### lint-test-e2e (`golangci-test-e2e`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20` (v9.2.0)
   - With:
     - `version`: `v2.9`
     - `args`: `--build-tags e2e ./test`

# Test attest / verify-attestation

| Property | Value |
|----------|-------|
| File | `kind-verify-attestation.yaml` |
| Triggers | `pull_request`, `workflow_dispatch` |

## Event filters

- **pull_request**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Defaults:** shell `bash`

## Jobs

### attest / verify-attestation test (`cip-test`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `KO_DOCKER_REPO` | `registry.local:5000/policy-controller` |
| `SCAFFOLDING_RELEASE_VERSION` | `v0.7.24` |
| `GO111MODULE` | `on` |
| `GOFLAGS` | `-ldflags=-s -ldflags=-w` |
| `KOCACHE` | `~/ko` |
| `COSIGN_YES` | `true` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **ko-build/setup-ko@v0.9**
   - Uses: `ko-build/setup-ko@d006021bd0c28d1ce33a07e7943d48b079944c8d` (v0.9)

5. **Install yq**
   - Uses: `mikefarah/yq@751d8ad57b84f1794661bc70c0afb92a22ad7b3c` (v4.53.2)

6. **build cosign**

7. **Install cluster + sigstore**
   - Uses: `sigstore/scaffolding/actions/setup@main`
   - With:
     - `legacy-variables`: `false`
     - `k8s-version`: `${{ matrix.k8s-version }}`
     - `version`: `${{ env.SCAFFOLDING_RELEASE_VERSION }}`

8. **Create sample image - demoimage**

9. **Initialize with our custom TUF root pointing to remote root**
   - Condition: `${{ matrix.tuf-root == 'remote' }}`

10. **Get copy of TUF repository**

11. **Initialize with custom TUF root pointing to local filesystem**
   - Condition: `${{ matrix.tuf-root == 'air-gap' }}`

12. **Set TrustedRoot**

13. **Create SigningConfig**

14. **Sign demoimage with cosign**

15. **Create attestation for it**

16. **Sign a blob**

17. **Verify with cosign**

18. **Verify custom attestation with cosign, works**

19. **Verify custom attestation with cosign, fails**

20. **Verify a blob**

21. **Collect diagnostics**
   - Uses: `chainguard-dev/actions/kind-diag@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - Condition: `${{ failure() }}`

22. **Create vuln attestation for it**

23. **Verify vuln attestation with cosign, works**

24. **Verify vuln attestation with cosign, fails**

# Scorecards supply-chain security

| Property | Value |
|----------|-------|
| File | `scorecard-action.yml` |
| Triggers | `branch_protection_rule`, `schedule`, `push` |

## Schedule

- `30 1 * * 6`

## Event filters

- **push**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SCORECARD_TOKEN` | job `analysis` step `Run analysis` with `repo_token` |

## Jobs

### Scorecards analysis (`analysis`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `security-events`: `write` - Needed to upload the results to code-scanning dashboard.
- `actions`: `read`
- `contents`: `read`
- `id-token`: `write` (OIDC)

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Run analysis**
   - Uses: `ossf/scorecard-action@4eaacf0543bb3f2c246792bd56e8cdeffafb205a` (v2.4.3)
   - With:
     - `results_file`: `results.sarif`
     - `results_format`: `sarif`
     - `repo_token`: `${{ secrets.SCORECARD_TOKEN }}`
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@65c74964a9ed8c44ed9f19d4bbc5757a6a8e9ab9` (v2.16.1)
   - With:
     - `sarif_file`: `results.sarif`

# CI-Tests

| Property | Value |
|----------|-------|
| File | `tests.yaml` |
| Triggers | `workflow_dispatch`, `push`, `pull_request` |

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`, `release-*`

## Jobs

### Run unit tests (`unit-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |

**Permissions:**

- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `OS` | `${{ matrix.os }}` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **actions/cache@v5.0.5**
   - Uses: `actions/cache@27d5ce7f107fe9357f9df03efb73ab90386fccae` (v5.0.5)
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build ~/Library/Caches/go-build %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **Run Go tests**

6. **Upload Coverage Report**
   - Uses: `codecov/codecov-action@57e3a136b779b570ffcdbf80b3bdc90e7fab3de2` (v6.0.0)
   - With:
     - `env_vars`: `OS`

7. **Run Go tests w/ `-race`**
   - Condition: `${{ runner.os == 'Linux' }}`

### Run e2e tests (`e2e-tests`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **free up disk space**

3. **check disk space**

4. **actions/cache@v5.0.5**
   - Uses: `actions/cache@27d5ce7f107fe9357f9df03efb73ab90386fccae` (v5.0.5)
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build ~/Library/Caches/go-build %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

5. **Extract version of Go to use**

6. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

7. **ko-build/setup-ko@v0.9**
   - Uses: `ko-build/setup-ko@d006021bd0c28d1ce33a07e7943d48b079944c8d` (v0.9)

8. **setup kind cluster**

9. **Run end-to-end tests**

10. **Collect diagnostics**
   - Uses: `chainguard-dev/actions/kind-diag@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - Condition: `${{ failure() }}`

### Run PowerShell E2E tests (`e2e-windows-powershell-tests`)

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **actions/cache@v5.0.5**
   - Uses: `actions/cache@27d5ce7f107fe9357f9df03efb73ab90386fccae` (v5.0.5)
   - With:
     - `path`: `~/go/pkg/mod %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

5. **Run e2e_test.ps1**

### license boilerplate check (`license-check`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Install addlicense**

5. **Check license headers**

# CI-Validate-Release-Job

| Property | Value |
|----------|-------|
| File | `validate-release.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `check-signature`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **Check Signature**
   - Env:
     - `TUF_ROOT`: `/tmp`

### `validate-release-job`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `check-signature` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **git config --system --add safe.directory /__w/cosign/cosign**

3. **free up disk space for the release**

4. **check disk space**

5. **goreleaser snapshot**
   - Env:
     - `PROJECT_ID`: `honk-fake-project`
     - `RUNTIME_IMAGE`: `gcr.io/distroless/static-debian13:nonroot`

6. **check binaries**

# Docgen

| Property | Value |
|----------|-------|
| File | `verify-docgen.yaml` |
| Triggers | `workflow_dispatch`, `push`, `pull_request` |

## Event filters

- **push**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### Verify Docgen (`docgen`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **deps**

2. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **./cmd/help/verify.sh**

# Whitespace

| Property | Value |
|----------|-------|
| File | `whitespace.yaml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### Check Whitespace (`whitespace`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `read`

#### Steps

1. **Check out code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **chainguard-dev/actions/trailing-space@v1.6.19**
   - Uses: `chainguard-dev/actions/trailing-space@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - Condition: `${{ always() }}`

3. **chainguard-dev/actions/eof-newline@v1.6.19**
   - Uses: `chainguard-dev/actions/eof-newline@c69a264ec2a5934c3186c618f368fc1c86f16cff` (v1.6.19)
   - Condition: `${{ always() }}`

