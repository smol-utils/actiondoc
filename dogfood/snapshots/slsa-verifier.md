# Contents

- [CodeQL](#codeql)
- [Dependency Review](#dependency-review)
- [Schedule cli](#schedule-cli)
- [verifier action](#verifier-action)
- [PR Title](#pr-title)
- [Actions pre submits](#actions-pre-submits)
- [Pre submits cli](#pre-submits-cli)
- [Pre submits e2e](#pre-submits-e2e)
- [LFS Warning](#lfs-warning)
- [Pre submits Lint](#pre-submits-lint)
- [References pre submits](#references-pre-submits)
- [Verifier releaser](#verifier-releaser)
- [Scorecards supply-chain security](#scorecards-supply-chain-security)
- [Update actions dist post-commit](#update-actions-dist-post-commit)

# CodeQL

For most projects, this workflow file will not need changing; you simply need to commit it to your repository.  You may wish to alter this file to override the set of languages analyzed, or to provide custom queries or build logic.  ******** NOTE ******** We have attempted to detect the languages in your repository. Please check the `language` matrix defined below to confirm you have the correct set of supported CodeQL languages.

| Property | Value |
|----------|-------|
| File | `codeql-analysis.yml` |
| Triggers | `push`, `pull_request`, `schedule` |

## Schedule

- `30 0 * * 2`

## Event filters

- **push**
  - branches: `main`, `*`
- **pull_request**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Jobs

### Analyze (`analyze`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `actions`: `read`
- `contents`: `read`
- `security-events`: `write`

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **setup-go**
   - Uses: `actions/setup-go@f111f3307d8850f501ac008e886eec1fd1932a34` (v5.3.0)
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@b6a472f63d85b9c78a3ac5e89422239fc15e9b3c` (v3.28.1)
   - With:
     - `languages`: `${{ matrix.language }}`

4. **Autobuild**
   - Uses: `github/codeql-action/autobuild@b6a472f63d85b9c78a3ac5e89422239fc15e9b3c` (v3.28.1)

5. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@b6a472f63d85b9c78a3ac5e89422239fc15e9b3c` (v3.28.1)

# Dependency Review

| Property | Value |
|----------|-------|
| File | `depsreview.yml` |
| Triggers | `pull_request` |

## Permissions

- `contents`: `read`

## Jobs

### `dependency-review`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **Dependency Review**
   - Uses: `actions/dependency-review-action@3b139cfc5fae8b618d3eae3675e383bb1769c019` (v4.5.0)

# Schedule cli

| Property | Value |
|----------|-------|
| File | `e2e.schedule.cli.yml` |
| Triggers | `workflow_run` |

## Event filters

- **workflow_run**
  - workflows: `Pre submits cli`
  - types: `completed`
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GH_TOKEN` | `${{ github.token }}` |
| `ISSUE_REPOSITORY` | `${{ github.repository }}` |

## Jobs

### `if-failed`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/download-artifact@v4.1.8**
   - Uses: `actions/download-artifact@fa0a91b85d4f404e444e00e005971372dc801d16` (v4.1.8)
   - With:
     - `name`: `event_name`

2. **Check event name**
   - ID: `name`

3. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - Condition: `steps.name.outputs.continue == 'true'`
   - With:
     - `ref`: `main`
     - `repository`: `slsa-framework/example-package`

4. **./.github/workflows/scripts/e2e-report-failure.sh**
   - Condition: `steps.name.outputs.continue == 'true' && github.event.workflow_run.conclusion != 'success'`

5. **./.github/workflows/scripts/e2e-report-success.sh**
   - Condition: `steps.name.outputs.continue == 'true' && github.event.workflow_run.conclusion == 'success'`

# verifier action

| Property | Value |
|----------|-------|
| File | `e2e.schedule.installer.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | No | - | The version to to test for pre-release. |

## Schedule

- `0 4 * * *`

## Permissions

All scopes: `read-all`.

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GH_TOKEN` | `${{ github.token }}` |
| `ISSUE_REPOSITORY` | `${{ github.repository }}` |
| `MINIMUM_INSTALLER_VERSION` | `v2.0.1` |

## Call graph (rooted at this workflow)

```
e2e.schedule.installer.yml [schedule, workflow_dispatch]
+-- verifier-run / Run the Action at tag (uses ./actions/installer (unresolved))
+-- verifier-run / Run the Action at commit (uses ./actions/installer (unresolved))
+-- verifier-run / Install invalid commit (uses ./actions/installer (unresolved))
+-- verifier-run / Install non-existent tag (uses ./actions/installer (unresolved))
+-- verifier-run / Install empty tag (uses ./actions/installer (unresolved))
```

## Jobs

### `list-verifiers`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **Checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `path`: `__THIS_REPO__`

3. **Generate verifier list**
   - ID: `generate-list`
   - Condition: `inputs.version == ''`

4. **Generate pre-release list**
   - ID: `generate-prerelease`
   - Condition: `inputs.version != ''`
   - Env:
     - `PRE_RELEASE_VERSION`: `${{ inputs.version }}`

5. **Generate pre-release list**
   - ID: `generate-versions`
   - Env:
     - `PRE_RELEASE_VERSION`: `${{ steps.generate-prerelease.outputs.version }}`
     - `LIST_VERSION`: `${{ steps.generate-list.outputs.version }}`

### `verifier-run`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `list-verifiers` |

#### Steps

1. **Debug**
   - Env:
     - `VERSION`: `${{ matrix.version }}`

2. **Checkout this repository**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - With:
     - `ref`: `${{ matrix.version }}`

3. **Run the Action at tag**
   - Uses: `./actions/installer`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `SLSA_VERIFIER_CI_ACTION_REF`: `${{ matrix.version }}`

4. **Verify the version**
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `VERSION`: `${{ matrix.version }}`

5. **Delete the binary**
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`

6. **Get sha1**
   - ID: `commit`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `VERSION`: `${{ matrix.version }}`

7. **Run the Action at commit**
   - Uses: `./actions/installer`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `SLSA_VERIFIER_CI_ACTION_REF`: `${{ steps.commit.outputs.commit_sha }}`

8. **Verify the version**
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `VERSION`: `${{ matrix.version }}`

9. **Delete the binary**
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`

10. **Install invalid commit** `[continue-on-error]`
   - ID: `invalid-commit`
   - Uses: `./actions/installer`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `SLSA_VERIFIER_CI_ACTION_REF`: `55ca6286e3e4f4fba5d0448333fa99fc5a404a73`

11. **[ "$SUCCESS" == "true" ]**
   - Env:
     - `SUCCESS`: `${{ steps.invalid-commit.outcome == 'failure' }}`

12. **Install non-existent tag** `[continue-on-error]`
   - ID: `nonexistent-tag`
   - Uses: `./actions/installer`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `SLSA_VERIFIER_CI_ACTION_REF`: `v100.3.5`

13. **[ "$SUCCESS" == "true" ]**
   - Env:
     - `SUCCESS`: `${{ steps.nonexistent-tag.outcome == 'failure' }}`

14. **Install empty tag** `[continue-on-error]`
   - ID: `empty-tag`
   - Uses: `./actions/installer`
   - Condition: `${{ inputs.version != '' || ! contains(matrix.version, '-rc' ) }}`
   - Env:
     - `SLSA_VERIFIER_CI_ACTION_REF`: ``

15. **[ "$SUCCESS" == "true" ]**
   - Env:
     - `SUCCESS`: `${{ steps.empty-tag.outcome == 'failure' }}`

### `if-succeed`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `verifier-run`, `list-verifiers` |
| Condition | `inputs.version == '' && needs.verifier-run.result != 'failure' && needs.list-verifiers.result != 'failure'` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-success.sh**

### `if-failed`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `verifier-run`, `list-verifiers` |
| Condition | `always() && inputs.version == '' && (needs.verifier-run.result == 'failure' \|\| needs.list-verifiers.result == 'failure')` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-failure.sh**

# PR Title

| Property | Value |
|----------|-------|
| File | `pr-title.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - types: `opened`, `edited`, `reopened`, `synchronize`

## Permissions

All scopes: `read-all`.

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `validate` step `thehanimo/pr-title-checker@7fbfe05602bdd86f926d3fb3bccb6f3aed43bc70` with `GITHUB_TOKEN` |

## Jobs

### `validate`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **thehanimo/pr-title-checker@v1.4.3**
   - Uses: `thehanimo/pr-title-checker@7fbfe05602bdd86f926d3fb3bccb6f3aed43bc70` (v1.4.3)
   - With:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `configuration_path`: `.github/pr-title-checker-config.json`

# Actions pre submits

| Property | Value |
|----------|-------|
| File | `pre-submit.actions.yml` |
| Triggers | `pull_request`, `workflow_dispatch` |

## Event filters

- **pull_request**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Jobs

### `check-dist`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **Set Node.js 20**
   - Uses: `actions/setup-node@39370e3970a6d050c480ffad4ff0ed4d3fdee5af` (v4.1.0)
   - With:
     - `node-version`: `20`

3. **Rebuild the dist/ directory**

4. **Compare the expected and actual dist/ directories**
   - ID: `diff`

5. **actions/upload-artifact@v4.6.0**
   - Uses: `actions/upload-artifact@65c4c4a1ddee5b72f698fdd19549f0f0fb45cf08` (v4.6.0)
   - Condition: `${{ failure() && steps.diff.conclusion == 'failure' }}`
   - With:
     - `name`: `dist`
     - `path`: `dist/`

# Pre submits cli

| Property | Value |
|----------|-------|
| File | `pre-submit.cli.yml` |
| Triggers | `pull_request`, `workflow_dispatch`, `schedule` |

## Schedule

- `25 6 * * 5`

## Event filters

- **pull_request**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Jobs

### `pre-submit`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **setup-go**
   - Uses: `actions/setup-go@f111f3307d8850f501ac008e886eec1fd1932a34` (v5.3.0)
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **Save event name**
   - Env:
     - `EVENT_NAME`: `${{ github.event_name }}`

4. **actions/upload-artifact@v4.6.0**
   - Uses: `actions/upload-artifact@65c4c4a1ddee5b72f698fdd19549f0f0fb45cf08` (v4.6.0)
   - With:
     - `name`: `event_name`
     - `path`: `./event_name.txt`

5. **Run tests for verifier**

# Pre submits e2e

| Property | Value |
|----------|-------|
| File | `pre-submit.e2e.yml` |
| Triggers | `pull_request`, `workflow_dispatch` |

## Event filters

- **pull_request**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `pre-submit` step `Run verification script with testdata and slsa-verifier HEAD` env `GH_TOKEN` |

## Jobs

### `pre-submit`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `path`: `__THIS_REPO__`

2. **setup-go**
   - Uses: `actions/setup-go@f111f3307d8850f501ac008e886eec1fd1932a34` (v5.3.0)
   - With:
     - `go-version-file`: `__THIS_REPO__/go.mod`
     - `cache`: `false`

3. **Build verifier at HEAD**

4. **Checkout e2e verification script**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `path`: `__EXAMPLE_PACKAGE__`
     - `repository`: `slsa-framework/example-package`

5. **Run verification script with testdata and slsa-verifier HEAD**
   - Env:
     - `SLSA_VERIFIER_TESTING`: `true`
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

# LFS Warning

| Property | Value |
|----------|-------|
| File | `pre-submit.lfs.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - types: `assigned`, `opened`, `edited`, `reopened`, `synchronize`, `ready_for_review`

## Permissions

All scopes: `read-all`.

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `large-file-check` step `actionsdesk/lfs-warning@4b98a8a5e6c429c23c34eee02d71553bca216425` with `token` |

## Jobs

### `large-file-check`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **actionsdesk/lfs-warning@v3.3**
   - Uses: `actionsdesk/lfs-warning@4b98a8a5e6c429c23c34eee02d71553bca216425` (v3.3)
   - With:
     - `token`: `${{ secrets.GITHUB_TOKEN }}`
     - `filesizelimit`: `10MB`
     - `labelName`: `lfs-warning`
     - `exclusionPatterns`: `cli/slsa-verifier/testdata/**`

# Pre submits Lint

| Property | Value |
|----------|-------|
| File | `pre-submit.lint.yml` |
| Triggers | `pull_request` |

## Permissions

- `contents`: `read` - Needed to check out the repo.

## Jobs

### `golangci-lint`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **actions/setup-go@v5.3.0**
   - Uses: `actions/setup-go@f111f3307d8850f501ac008e886eec1fd1932a34` (v5.3.0)
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@ec5d18412c0aeab7936cb16880d708ba2a64e1ae` (v6.2.0)
   - With:
     - `version`: `v1.61.0`

### `yamllint`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **set -euo pipefail**

### `eslint`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **actions/setup-node@v4.1.0**
   - Uses: `actions/setup-node@39370e3970a6d050c480ffad4ff0ed4d3fdee5af` (v4.1.0)
   - With:
     - `node-version`: `20`

3. **make eslint**

### `renovate-config-validator`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **actions/setup-node@v4.1.0**
   - Uses: `actions/setup-node@39370e3970a6d050c480ffad4ff0ed4d3fdee5af` (v4.1.0)
   - With:
     - `node-version`: `20`

3. **make renovate-config-validator**

# References pre submits

| Property | Value |
|----------|-------|
| File | `pre-submit.references.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - types: `opened`, `edited`, `reopened`, `synchronize`

## Permissions

All scopes: `read-all`.

## Jobs

### `check-docs`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `${{ contains(github.event.pull_request.body, '#label:release') }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BODY` | `${{ github.event.pull_request.body }}` |

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **Check documentation is up-to-date**

# Verifier releaser

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Triggers | `workflow_dispatch`, `push`, `schedule` |

## Schedule

- `0 1 * * *`

## Event filters

- **push**
  - tags: `*`

## Permissions

All scopes: `read-all`.

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GH_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `ISSUE_REPOSITORY` | `slsa-framework/slsa-verifier` |
| `HEADER` | `release` |

## Call graph (rooted at this workflow)

```
release.yml [workflow_dispatch, push, schedule]
+-- builder (uses slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0)
```

## Transitive requirements (from full call graph)

External workflows referenced: `slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GH_TOKEN`; job `verification` step `Download assets` env `GH_TOKEN` |

## Jobs

### `args`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **checkout**
   - ID: `checkout`
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `fetch-depth`: `0`

2. **ldflags**
   - ID: `ldflags`

### builder-linux, windows, darwin-amd64, arm64 (`builder`)

| Property | Value |
|----------|-------|
| Uses workflow | `slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0` (external) |
| Depends on | `args` |

**Permissions:**

- `actions`: `read` - For the detection of GitHub Actions environment.
- `id-token`: `write` (OIDC) - For signing.
- `contents`: `write` - For asset uploads.

#### Inputs forwarded

- `go-version-file`: `go.mod`
- `config-file`: `.slsa-goreleaser/${{matrix.os}}-${{matrix.arch}}.yml`
- `compile-builder`: `true`
- `evaluated-envs`: `VERSION:${{needs.args.outputs.version}}`

### `verification`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `builder` |
| Condition | `github.event_name != 'schedule' && github.event_name != 'workflow_dispatch'` |

**Permissions:**

All scopes: `read-all`.

#### Steps

1. **Install the verifier**
   - Uses: `slsa-framework/slsa-verifier/actions/installer@3714a2a4684014deb874a0e737dffa0ee02dd647` (v2.6.0)

2. **Download assets**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `ATT_FILE_NAME`: `${{ needs.builder.outputs.go-binary-name }}.intoto.jsonl`
     - `ARTIFACT`: `${{ needs.builder.outputs.go-binary-name }}`

3. **Verify assets**
   - Env:
     - `ARTIFACT`: `${{ needs.builder.outputs.go-binary-name }}`
     - `ATT_FILE_NAME`: `${{ needs.builder.outputs.go-binary-name }}.intoto.jsonl`

### `if-succeed`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `args`, `builder` |
| Condition | `github.event_name == 'schedule' && needs.args.result != 'failure' && needs.builder.result != 'failure'` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-success.sh**

### `if-failed`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `args`, `builder` |
| Condition | `always() && github.event_name == 'schedule' && (needs.args.result == 'failure' \|\| needs.builder.result == 'failure')` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

#### Steps

1. **actions/checkout@v4.2.2**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-failure.sh**

# Scorecards supply-chain security

| Property | Value |
|----------|-------|
| File | `scorecards.yml` |
| Triggers | `branch_protection_rule`, `schedule`, `push` |

## Schedule

- `25 6 * * 5`

## Event filters

- **push**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Jobs

### Scorecards analysis (`analysis`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `security-events`: `write` - Needed to upload the results to code-scanning dashboard.
- `id-token`: `write` (OIDC) - Used to receive a badge. (Upcoming feature)
- `contents`: `read` - Needs for private repositories.
- `actions`: `read`

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `persist-credentials`: `false`

2. **Run analysis**
   - Uses: `ossf/scorecard-action@62b2cac7ed8198b15735ed49ab1e5cf35480ba46` (v2.4.0)
   - With:
     - `results_file`: `results.sarif`
     - `results_format`: `sarif`
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@65c4c4a1ddee5b72f698fdd19549f0f0fb45cf08` (v4.6.0)
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@b6a472f63d85b9c78a3ac5e89422239fc15e9b3c` (v3.28.1)
   - With:
     - `sarif_file`: `results.sarif`

# Update actions dist post-commit

A workflow to run against renovate-bot's PRs, such as `make package` after it updates the package.json and package-lock.json files. The potentially untrusted code is first run inside a low-privilege Job, and the diff is uploaded as an artifact. Then a higher-privilege Job applies the diff and pushes the changes to the PR. It's important to only run this workflow against PRs from trusted sources, after also reviewing the changes! There have been vulnerabilities with using `git apply` https://github.blog/2023-04-25-git-security-vulnerabilities-announced-4/ At this point a compromised git binary cannot modify any of this repo's branches, only the PR fork's branch, due to our branch protection rules and CODEOWNERS. It aslso cannot submit a new release or modify exsiting releases due to tag protection rules.

| Property | Value |
|----------|-------|
| File | `update-actions-dist-post-commit.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `pr_number` | number | Yes | - | The pull request number. |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `diff`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `pull-requests`: `read` - This Job executes the PR's untrusted code, so it must how low permissions.

#### Steps

1. **checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
   - With:
     - `repository`: `${{ github.repository }}`
     - `persist-credentials`: `false`

2. **checkout-pr**
   - Env:
     - `GH_TOKEN`: `${{ github.token }}`
     - `PR_NUMBER`: `${{ inputs.pr_number }}`

3. **run-command**

4. **diff**
   - ID: `diff`

5. **upload**
   - Uses: `actions/upload-artifact@65c4c4a1ddee5b72f698fdd19549f0f0fb45cf08` (v4.6.0)
   - With:
     - `name`: `changes.patch`
     - `path`: `changes.patch`

### `push`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `diff` |
| Condition | `needs.diff.outputs.patch_not_empty == 'true'` |

**Permissions:**

- `pull-requests`: `read` - This Job does not run untrusted code, but it does need to push changes to the PR's branch.
- `contents`: `write`

#### Steps

1. **checkout**
   - Uses: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)

2. **checkout-pr**
   - Env:
     - `GH_TOKEN`: `${{ github.token }}`
     - `PR_NUMBER`: `${{ inputs.pr_number }}`

3. **download-patch**
   - Uses: `actions/download-artifact@fa0a91b85d4f404e444e00e005971372dc801d16` (v4.1.8)
   - With:
     - `name`: `changes.patch`

4. **apply**
   - ID: `apply`

5. **push**

