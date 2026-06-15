# cosign

17 workflows

## Contents

**Workflows**

- [CI-Container-Build - push](#ci-container-build)
- [CodeQL - push](#codeql)
- [Conformance Tests Nightly - schedule, workflow_dispatch](#conformance-tests-nightly)
- [Conformance Tests - push, pull_request](#conformance-tests)
- [Cut Release - workflow_dispatch](#cut-release)
- [Dependency Review - pull_request](#dependency-review)
- [Do Not Submit - pull_request](#do-not-submit)
- [e2e-tests - push, pull_request, workflow_dispatch](#e2e-tests)
- [e2e-with-binary - push, workflow_dispatch](#e2e-with-binary)
- [Test GitHub OIDC - push, schedule, workflow_dispatch](#test-github-oidc)
- [golangci-lint - push, pull_request](#golangci-lint)
- [Test attest / verify-attestation - pull_request, workflow_dispatch](#test-attest--verify-attestation)
- [Scorecards supply-chain security - branch_protection_rule, schedule, push](#scorecards-supply-chain-security)
- [CI-Tests - workflow_dispatch, push, pull_request](#ci-tests)
- [CI-Validate-Release-Job - pull_request](#ci-validate-release-job)
- [Docgen - workflow_dispatch, push, pull_request](#docgen)
- [Whitespace - pull_request](#whitespace)

# CI-Container-Build

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `build.yaml` |

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

<details>
<summary>Steps (9)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **sigstore/cosign-installer@v4.1.2**

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **ko-build/setup-ko@v0.9**

6. **Set up Cloud SDK**
   - Uses: `google-github-actions/auth@v3.0.0`
   - With:
     - `workload_identity_provider`: `projects/498091336538/locations/global/workloadIdentityPools/githubactions/providers/sigstore-cosign`
     - `service_account`: `github-actions@projectsigstore.iam.gserviceaccount.com`

7. **creds**

8. **Login to GitHub Container Registry**
   - Uses: `docker/login-action@v4.1.0`
   - With:
     - `registry`: `ghcr.io`
     - `username`: `${{ github.actor }}`
     - `password`: `${{ secrets.GITHUB_TOKEN }}`

9. **containers-cosign**
   - Env:
     - `KO_PREFIX`: `ghcr.io/sigstore/cosign/cosign/ci`
     - `COSIGN_PASSWORD`: `${{secrets.COSIGN_PASSWORD}}`

</details>

[Back to top](#contents)

# CodeQL

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `codeql-analysis.yml` |

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
| Matrix | `language`: go |
| Condition | `github.repository == 'sigstore/cosign'` |

**Permissions:**

- `security-events`: `write`
- `actions`: `read`
- `contents`: `read`

<details>
<summary>Steps (7)</summary>

1. **Checkout repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Utilize Go Module Cache**
   - Uses: `actions/cache@v5.0.5`
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@v2.16.1`
   - With:
     - `languages`: `${{ matrix.language }}`

6. **Build cosign for CodeQL**

7. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@v2.16.1`

</details>

[Back to top](#contents)

# Conformance Tests Nightly

**Triggers:** `schedule`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `conformance-nightly.yml` |

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

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **make cosign conformance**

5. **sigstore/sigstore-conformance@main**
   - With:
     - `entrypoint`: `${{ github.workspace }}/conformance`
     - `xfail`: `test_verify*PATH-message-digest-mismatch_fail]`

6. **Create Issue on Failure**
   - Uses: `actions/github-script@v8.0.0`
   - Condition: `failure()`
   - With:
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
     - `script`: `` const { owner, repo } = context.repo; const runId = context.runId; const issueTitle = 'Conformance Tests Failed'; const issueBody = `The nightly conformance tests have failed. Please check the logs for more details.\n\nWorkflow run: https://github.com/${owner}/${repo}/actions/runs/${runId}\n\ncc @sigstore/security-response-team @sigstore/cosign-codeowners`; const issueLabel = 'bug';  const existingIssues = await github.rest.issues.listForRepo({   owner,   repo,   state: 'open',   labels: issueLabel, });  const issueExists = existingIssues.data.some(issue => issue.title === issueTitle);  if (!issueExists) {   await github.rest.issues.create({     owner,     repo,     title: issueTitle,     body: issueBody,     labels: [issueLabel],   }); } ``

</details>

[Back to top](#contents)

# Conformance Tests

**Triggers:** `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `conformance.yml` |

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

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **make cosign conformance**

5. **sigstore/sigstore-conformance@v0.0.27**
   - With:
     - `entrypoint`: `${{ github.workspace }}/conformance`
     - `xfail`: `test_verify*PATH-message-digest-mismatch_fail]`

</details>

[Back to top](#contents)

# Cut Release

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `cut-release.yml` |

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

[Back to top](#contents)

# Dependency Review

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `depsreview.yml` |

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

[Back to top](#contents)

# Do Not Submit

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `donotsubmit.yaml` |

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

<details>
<summary>Steps (2)</summary>

1. **Check out code**
   - Uses: `actions/checkout@v2.4.0`
   - With:
     - `persist-credentials`: `false`

2. **Do Not Submit**
   - Uses: `chainguard-dev/actions/donotsubmit@v1.6.19`

</details>

[Back to top](#contents)

# e2e-tests

**Triggers:** `push`, `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `e2e-tests.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`e2e-cross`](#e2e-cross), [`e2e-test-pkcs11`](#e2e-test-pkcs11), [`e2e-kms`](#e2e-kms), [`e2e-registry`](#e2e-registry)

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`

## Jobs

### `e2e-cross`

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Matrix | `os`: macos-latest, ubuntu-latest |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Run cross platform e2e tests**

</details>

### `e2e-test-pkcs11`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Run pkcs11 end-to-end tests**

</details>

### `e2e-kms`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VAULT_TOKEN` | `root` |
| `VAULT_ADDR` | `http://localhost:8200` |
| `COSIGN_YES` | `true` |
| `SCAFFOLDING_RELEASE_VERSION` | `v0.7.24` |

<details>
<summary>Steps (8)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **setup vault**
   - Uses: `cpanato/vault-installer@v1.4.0`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **imjasonh/setup-crane@v0.5**

6. **Install cluster + sigstore**
   - Uses: `sigstore/scaffolding/actions/setup@main`
   - With:
     - `version`: `${{ env.SCAFFOLDING_RELEASE_VERSION }}`

7. **enable vault transit**

8. **Acceptance Tests**

</details>

### `e2e-registry`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `SCAFFOLDING_RELEASE_VERSION` | `v0.7.24` |

<details>
<summary>Steps (12)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Setup mirror**
   - Uses: `chainguard-dev/actions/setup-mirror@v1.6.19`
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
   - Uses: `chainguard-dev/actions/kind-diag@v1.6.19`
   - Condition: `${{ failure() }}`

</details>

[Back to top](#contents)

# e2e-with-binary

**Triggers:** `push`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `e2e-with-binary.yml` |

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
| Matrix | `os`: macos-latest, ubuntu-latest, windows-latest |
| Condition | `${{ github.repository == 'sigstore/cosign' }}` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COSIGN_YES` | `true` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **build cosign and check sign-blob and verify-blob**

</details>

[Back to top](#contents)

# Test GitHub OIDC

**Triggers:** `push`, `schedule`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `github-oidc.yaml` |

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

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **ko-build/setup-ko@v0.9**

5. **build cosign from the HEAD**

6. **Build and sign a container image**

</details>

[Back to top](#contents)

# golangci-lint

**Triggers:** `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `golangci-lint.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [lint](#lint-golangci), [lint-test-e2e](#lint-test-e2e-golangci-test-e2e)

## Event filters

- **push**
  - branches: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### lint (`golangci`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@v9.2.0`
   - With:
     - `version`: `v2.12`

</details>

### lint-test-e2e (`golangci-test-e2e`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@v9.2.0`
   - With:
     - `version`: `v2.9`
     - `args`: `--build-tags e2e ./test`

</details>

[Back to top](#contents)

# Test attest / verify-attestation

**Triggers:** `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `kind-verify-attestation.yaml` |

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
| Matrix | `k8s-version`: v1.30.x, v1.31.x, v1.32.x, v1.33.x; `tuf-root`: remote, air-gap |

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

<details>
<summary>Steps (24)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **ko-build/setup-ko@v0.9**

5. **Install yq**
   - Uses: `mikefarah/yq@v4.53.2`

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
   - Uses: `chainguard-dev/actions/kind-diag@v1.6.19`
   - Condition: `${{ failure() }}`

22. **Create vuln attestation for it**

23. **Verify vuln attestation with cosign, works**

24. **Verify vuln attestation with cosign, fails**

</details>

[Back to top](#contents)

# Scorecards supply-chain security

**Triggers:** `branch_protection_rule`, `schedule`, `push`

| Property | Value |
|----------|-------|
| File | `scorecard-action.yml` |

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

- `security-events`: `write`
- `actions`: `read`
- `contents`: `read`
- `id-token`: `write` (OIDC)

<details>
<summary>Steps (4)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Run analysis**
   - Uses: `ossf/scorecard-action@v2.4.3`
   - With:
     - `results_file`: `results.sarif`
     - `results_format`: `sarif`
     - `repo_token`: `${{ secrets.SCORECARD_TOKEN }}`
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@v7.0.1`
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@v2.16.1`
   - With:
     - `sarif_file`: `results.sarif`

</details>

[Back to top](#contents)

# CI-Tests

**Triggers:** `workflow_dispatch`, `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `tests.yaml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [Run unit tests](#run-unit-tests-unit-tests), [Run e2e tests](#run-e2e-tests-e2e-tests), [Run PowerShell E2E tests](#run-powershell-e2e-tests-e2e-windows-powershell-tests), [license boilerplate check](#license-boilerplate-check-license-check)

## Event filters

- **push**
  - paths: `**`, `!**.md`, `!doc/**`, `!**.txt`, `!images/**`, `!LICENSE`, `test/**`
  - branches: `main`, `release-*`

## Jobs

### Run unit tests (`unit-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Matrix | `os`: macos-latest, ubuntu-latest, windows-latest |

**Permissions:**

- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `OS` | `${{ matrix.os }}` |

<details>
<summary>Steps (7)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **actions/cache@v5.0.5**
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build ~/Library/Caches/go-build %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **Run Go tests**

6. **Upload Coverage Report**
   - Uses: `codecov/codecov-action@v6.0.0`
   - With:
     - `env_vars`: `OS`

7. **Run Go tests w/ \`-race\`**
   - Condition: `${{ runner.os == 'Linux' }}`

</details>

### Run e2e tests (`e2e-tests`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (10)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **free up disk space**

3. **check disk space**

4. **actions/cache@v5.0.5**
   - With:
     - `path`: `~/go/pkg/mod ~/.cache/go-build ~/Library/Caches/go-build %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

5. **Extract version of Go to use**

6. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

7. **ko-build/setup-ko@v0.9**

8. **setup kind cluster**

9. **Run end-to-end tests**

10. **Collect diagnostics**
   - Uses: `chainguard-dev/actions/kind-diag@v1.6.19`
   - Condition: `${{ failure() }}`

</details>

### Run PowerShell E2E tests (`e2e-windows-powershell-tests`)

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **actions/cache@v5.0.5**
   - With:
     - `path`: `~/go/pkg/mod %LocalAppData%\go-build`
     - `key`: `${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}`
     - `restore-keys`: `${{ runner.os }}-go-`

5. **Run e2e\_test.ps1**

</details>

### license boilerplate check (`license-check`)

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Extract version of Go to use**

3. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

4. **Install addlicense**

5. **Check license headers**

</details>

[Back to top](#contents)

# CI-Validate-Release-Job

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `validate-release.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`check-signature`](#check-signature), [`validate-release-job`](#validate-release-job)

## Event filters

- **pull_request**
  - branches: `main`, `release-*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `check-signature`

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (1)</summary>

1. **Check Signature**
   - Env:
     - `TUF_ROOT`: `/tmp`

</details>

### `validate-release-job`

| Property | Value |
|----------|-------|
| Depends on | `check-signature` |

**Permissions:**

- `contents`: `read`

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **git config --system --add safe.directory /\_\_w/cosign/cosign**

3. **free up disk space for the release**

4. **check disk space**

5. **goreleaser snapshot**
   - Env:
     - `PROJECT_ID`: `honk-fake-project`
     - `RUNTIME_IMAGE`: `gcr.io/distroless/static-debian13:nonroot`

6. **check binaries**

</details>

[Back to top](#contents)

# Docgen

**Triggers:** `workflow_dispatch`, `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `verify-docgen.yaml` |

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

<details>
<summary>Steps (5)</summary>

1. **deps**

2. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

3. **Extract version of Go to use**

4. **actions/setup-go@v6.4.0**
   - With:
     - `go-version`: `${{ env.GOVERSION }}`
     - `check-latest`: `true`
     - `cache`: `false`

5. **./cmd/help/verify.sh**

</details>

[Back to top](#contents)

# Whitespace

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `whitespace.yaml` |

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

<details>
<summary>Steps (3)</summary>

1. **Check out code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **chainguard-dev/actions/trailing-space@v1.6.19**
   - Condition: `${{ always() }}`

3. **chainguard-dev/actions/eof-newline@v1.6.19**
   - Condition: `${{ always() }}`

</details>

[Back to top](#contents)

