# slsa-verifier

14 workflows

## Contents

**Workflows**

- [CodeQL - push, pull_request, schedule](#codeql)
- [Dependency Review - pull_request](#dependency-review)
- [Schedule cli - workflow_run](#schedule-cli)
- [verifier action - schedule, workflow_dispatch](#verifier-action)
- [PR Title - pull_request](#pr-title)
- [Actions pre submits - pull_request, workflow_dispatch](#actions-pre-submits)
- [Pre submits cli - pull_request, workflow_dispatch, schedule](#pre-submits-cli)
- [Pre submits e2e - pull_request, workflow_dispatch](#pre-submits-e2e)
- [LFS Warning - pull_request](#lfs-warning)
- [Pre submits Lint - pull_request](#pre-submits-lint)
- [References pre submits - pull_request](#references-pre-submits)
- [Verifier releaser - workflow_dispatch, push, schedule](#verifier-releaser)
- [Scorecards supply-chain security - branch_protection_rule, schedule, push](#scorecards-supply-chain-security)
- [Update actions dist post-commit - workflow_dispatch](#update-actions-dist-post-commit)

# CodeQL

**Triggers:** `push`, `pull_request`, `schedule`

For most projects, this workflow file will not need changing; you simply need to commit it to your repository. You may wish to alter this file to override the set of languages analyzed, or to provide custom queries or build logic. ******** NOTE ******** We have attempted to detect the languages in your repository. Please check the `language` matrix defined below to confirm you have the correct set of supported CodeQL languages.

| Property | Value |
|----------|-------|
| File | `codeql-analysis.yml` |

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
| Matrix | `language`: go, javascript |

**Permissions:**

- `actions`: `read`
- `contents`: `read`
- `security-events`: `write`

<details>
<summary>Steps (5)</summary>

1. **Checkout repository**
   - Uses: `actions/checkout@v4.2.2`

2. **setup-go**
   - Uses: `actions/setup-go@v5.3.0`
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@v3.28.1`
   - With:
     - `languages`: `${{ matrix.language }}`

4. **Autobuild**
   - Uses: `github/codeql-action/autobuild@v3.28.1`

5. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@v3.28.1`

</details>

[Back to top](#contents)

# Dependency Review

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `depsreview.yml` |

## Permissions

- `contents`: `read`

## Jobs

### `dependency-review`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (2)</summary>

1. **Checkout Repository**
   - Uses: `actions/checkout@v4.2.2`

2. **Dependency Review**
   - Uses: `actions/dependency-review-action@v4.5.0`

</details>

[Back to top](#contents)

# Schedule cli

**Triggers:** `workflow_run`

| Property | Value |
|----------|-------|
| File | `e2e.schedule.cli.yml` |

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

<details>
<summary>Steps (5)</summary>

1. **actions/download-artifact@v4.1.8**
   - With:
     - `name`: `event_name`

2. **Check event name**
   - ID: `name`

3. **actions/checkout@v4.2.2**
   - Condition: `steps.name.outputs.continue == 'true'`
   - With:
     - `ref`: `main`
     - `repository`: `slsa-framework/example-package`

4. **./.github/workflows/scripts/e2e-report-failure.sh**
   - Condition: `steps.name.outputs.continue == 'true' && github.event.workflow_run.conclusion != 'success'`

5. **./.github/workflows/scripts/e2e-report-success.sh**
   - Condition: `steps.name.outputs.continue == 'true' && github.event.workflow_run.conclusion == 'success'`

</details>

[Back to top](#contents)

# verifier action

**Triggers:** `schedule`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `e2e.schedule.installer.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`list-verifiers`](#list-verifiers), [`verifier-run`](#verifier-run), [`if-succeed`](#if-succeed), [`if-failed`](#if-failed-1)

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
+-- verifier-run / Run the Action at tag (uses ./actions/installer (outside scan scope))
+-- verifier-run / Run the Action at commit (uses ./actions/installer (outside scan scope))
+-- verifier-run / Install invalid commit (uses ./actions/installer (outside scan scope))
+-- verifier-run / Install non-existent tag (uses ./actions/installer (outside scan scope))
+-- verifier-run / Install empty tag (uses ./actions/installer (outside scan scope))
```

## Jobs

### `list-verifiers`

<details>
<summary>Steps (5)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.2.2`
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **Checkout**
   - Uses: `actions/checkout@v4.2.2`
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

</details>

### `verifier-run`

| Property | Value |
|----------|-------|
| Matrix | `version`: ${{ fromJson(needs.list-verifiers.outputs.version) }} |
| Depends on | `list-verifiers` |

<details>
<summary>Steps (15)</summary>

1. **Debug**
   - Env:
     - `VERSION`: `${{ matrix.version }}`

2. **Checkout this repository**
   - Uses: `actions/checkout@v4.2.2`
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
     - `SLSA_VERIFIER_CI_ACTION_REF`: -

15. **[ "$SUCCESS" == "true" ]**
   - Env:
     - `SUCCESS`: `${{ steps.empty-tag.outcome == 'failure' }}`

</details>

### `if-succeed`

| Property | Value |
|----------|-------|
| Depends on | `verifier-run`, `list-verifiers` |
| Condition | `inputs.version == '' && needs.verifier-run.result != 'failure' && needs.list-verifiers.result != 'failure'` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-success.sh**

</details>

### `if-failed`

| Property | Value |
|----------|-------|
| Depends on | `verifier-run`, `list-verifiers` |
| Condition | `always() && inputs.version == '' && (needs.verifier-run.result == 'failure' \|\| needs.list-verifiers.result == 'failure')` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-failure.sh**

</details>

[Back to top](#contents)

# PR Title

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `pr-title.yml` |

## Event filters

- **pull_request**
  - types: `opened`, `edited`, `reopened`, `synchronize`

## Permissions

All scopes: `read-all`.

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `validate` step `thehanimo/pr-title-checker@v1.4.3` with `GITHUB_TOKEN` |

## Jobs

### `validate`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (1)</summary>

1. **thehanimo/pr-title-checker@v1.4.3**
   - With:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `configuration_path`: `.github/pr-title-checker-config.json`

</details>

[Back to top](#contents)

# Actions pre submits

**Triggers:** `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `pre-submit.actions.yml` |

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

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v4.2.2**

2. **Set Node.js 20**
   - Uses: `actions/setup-node@v4.1.0`
   - With:
     - `node-version`: `20`

3. **Rebuild the dist/ directory**

4. **Compare the expected and actual dist/ directories**
   - ID: `diff`

5. **actions/upload-artifact@v4.6.0**
   - Condition: `${{ failure() && steps.diff.conclusion == 'failure' }}`
   - With:
     - `name`: `dist`
     - `path`: `dist/`

</details>

[Back to top](#contents)

# Pre submits cli

**Triggers:** `pull_request`, `workflow_dispatch`, `schedule`

| Property | Value |
|----------|-------|
| File | `pre-submit.cli.yml` |

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

<details>
<summary>Steps (5)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.2.2`

2. **setup-go**
   - Uses: `actions/setup-go@v5.3.0`
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **Save event name**
   - Env:
     - `EVENT_NAME`: `${{ github.event_name }}`

4. **actions/upload-artifact@v4.6.0**
   - With:
     - `name`: `event_name`
     - `path`: `./event_name.txt`

5. **Run tests for verifier**

</details>

[Back to top](#contents)

# Pre submits e2e

**Triggers:** `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `pre-submit.e2e.yml` |

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

<details>
<summary>Steps (5)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.2.2`
   - With:
     - `path`: `__THIS_REPO__`

2. **setup-go**
   - Uses: `actions/setup-go@v5.3.0`
   - With:
     - `go-version-file`: `__THIS_REPO__/go.mod`
     - `cache`: `false`

3. **Build verifier at HEAD**

4. **Checkout e2e verification script**
   - Uses: `actions/checkout@v4.2.2`
   - With:
     - `path`: `__EXAMPLE_PACKAGE__`
     - `repository`: `slsa-framework/example-package`

5. **Run verification script with testdata and slsa-verifier HEAD**
   - Env:
     - `SLSA_VERIFIER_TESTING`: `true`
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

</details>

[Back to top](#contents)

# LFS Warning

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `pre-submit.lfs.yml` |

## Event filters

- **pull_request**
  - types: `assigned`, `opened`, `edited`, `reopened`, `synchronize`, `ready_for_review`

## Permissions

All scopes: `read-all`.

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `large-file-check` step `actionsdesk/lfs-warning@v3.3` with `token` |

## Jobs

### `large-file-check`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (2)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.2.2`

2. **actionsdesk/lfs-warning@v3.3**
   - With:
     - `token`: `${{ secrets.GITHUB_TOKEN }}`
     - `filesizelimit`: `10MB`
     - `labelName`: `lfs-warning`
     - `exclusionPatterns`: `cli/slsa-verifier/testdata/**`

</details>

[Back to top](#contents)

# Pre submits Lint

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `pre-submit.lint.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`golangci-lint`](#golangci-lint), [`yamllint`](#yamllint), [`eslint`](#eslint), [`renovate-config-validator`](#renovate-config-validator)

## Permissions

- `contents`: `read`

## Jobs

### `golangci-lint`

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v4.2.2**

2. **actions/setup-go@v5.3.0**
   - With:
     - `go-version-file`: `go.mod`
     - `cache`: `false`

3. **golangci-lint**
   - Uses: `golangci/golangci-lint-action@v6.2.0`
   - With:
     - `version`: `v1.61.0`

</details>

### `yamllint`

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**

2. **set -euo pipefail**

</details>

### `eslint`

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v4.2.2**

2. **actions/setup-node@v4.1.0**
   - With:
     - `node-version`: `20`

3. **make eslint**

</details>

### `renovate-config-validator`

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v4.2.2**

2. **actions/setup-node@v4.1.0**
   - With:
     - `node-version`: `20`

3. **make renovate-config-validator**

</details>

[Back to top](#contents)

# References pre submits

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `pre-submit.references.yml` |

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

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**

2. **Check documentation is up-to-date**

</details>

[Back to top](#contents)

# Verifier releaser

**Triggers:** `workflow_dispatch`, `push`, `schedule`

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`args`](#args), [builder-${{matrix.os}}-${{matrix.arch}}](#builder-matrixos-matrixarch-builder), [`verification`](#verification), [`if-succeed`](#if-succeed-1), [`if-failed`](#if-failed-2)

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

Secrets referenced (literal names): `GITHUB_TOKEN`

Permissions declared across the chain: `actions: read`, `contents: read`, `contents: write`, `id-token: write (OIDC)`, `issues: write`, `read-all`

External workflows referenced: `slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GH_TOKEN`; job `verification` step `Download assets` env `GH_TOKEN` |

## Jobs

### `args`

<details>
<summary>Steps (2)</summary>

1. **checkout**
   - ID: `checkout`
   - Uses: `actions/checkout@v4.2.2`
   - With:
     - `fetch-depth`: `0`

2. **ldflags**
   - ID: `ldflags`

</details>

### builder-${{matrix.os}}-${{matrix.arch}} (`builder`)

| Property | Value |
|----------|-------|
| Uses workflow | `slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0` (external) |
| Matrix | `os`: linux, windows, darwin; `arch`: amd64, arm64 |
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
| Depends on | `builder` |
| Condition | `github.event_name != 'schedule' && github.event_name != 'workflow_dispatch'` |

**Permissions:**

All scopes: `read-all`.

<details>
<summary>Steps (3)</summary>

1. **Install the verifier**
   - Uses: `slsa-framework/slsa-verifier/actions/installer@v2.6.0`

2. **Download assets**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `ATT_FILE_NAME`: `${{ needs.builder.outputs.go-binary-name }}.intoto.jsonl`
     - `ARTIFACT`: `${{ needs.builder.outputs.go-binary-name }}`

3. **Verify assets**
   - Env:
     - `ARTIFACT`: `${{ needs.builder.outputs.go-binary-name }}`
     - `ATT_FILE_NAME`: `${{ needs.builder.outputs.go-binary-name }}.intoto.jsonl`

</details>

### `if-succeed`

| Property | Value |
|----------|-------|
| Depends on | `args`, `builder` |
| Condition | `github.event_name == 'schedule' && needs.args.result != 'failure' && needs.builder.result != 'failure'` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-success.sh**

</details>

### `if-failed`

| Property | Value |
|----------|-------|
| Depends on | `args`, `builder` |
| Condition | `always() && github.event_name == 'schedule' && (needs.args.result == 'failure' \|\| needs.builder.result == 'failure')` |

**Permissions:**

- `contents`: `read`
- `issues`: `write`

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v4.2.2**
   - With:
     - `repository`: `slsa-framework/example-package`
     - `ref`: `main`

2. **./.github/workflows/scripts/e2e-report-failure.sh**

</details>

[Back to top](#contents)

# Scorecards supply-chain security

**Triggers:** `branch_protection_rule`, `schedule`, `push`

| Property | Value |
|----------|-------|
| File | `scorecards.yml` |

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

- `security-events`: `write`
- `id-token`: `write` (OIDC)
- `contents`: `read`
- `actions`: `read`

<details>
<summary>Steps (4)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v4.2.2`
   - With:
     - `persist-credentials`: `false`

2. **Run analysis**
   - Uses: `ossf/scorecard-action@v2.4.0`
   - With:
     - `results_file`: `results.sarif`
     - `results_format`: `sarif`
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@v4.6.0`
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@v3.28.1`
   - With:
     - `sarif_file`: `results.sarif`

</details>

[Back to top](#contents)

# Update actions dist post-commit

**Triggers:** `workflow_dispatch`

A workflow to run against renovate-bot's PRs, such as `make package` after it updates the package.json and package-lock.json files. The potentially untrusted code is first run inside a low-privilege Job, and the diff is uploaded as an artifact. Then a higher-privilege Job applies the diff and pushes the changes to the PR. It's important to only run this workflow against PRs from trusted sources, after also reviewing the changes! There have been vulnerabilities with using `git apply` https://github.blog/2023-04-25-git-security-vulnerabilities-announced-4/ At this point a compromised git binary cannot modify any of this repo's branches, only the PR fork's branch, due to our branch protection rules and CODEOWNERS. It aslso cannot submit a new release or modify exsiting releases due to tag protection rules.

| Property | Value |
|----------|-------|
| File | `update-actions-dist-post-commit.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`diff`](#diff), [`push`](#push)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `pr_number` | number | Yes | - | The pull request number. |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `diff`

**Permissions:**

- `pull-requests`: `read`

<details>
<summary>Steps (5)</summary>

1. **checkout**
   - Uses: `actions/checkout@v4.2.2`
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
   - Uses: `actions/upload-artifact@v4.6.0`
   - With:
     - `name`: `changes.patch`
     - `path`: `changes.patch`

</details>

### `push`

| Property | Value |
|----------|-------|
| Depends on | `diff` |
| Condition | `needs.diff.outputs.patch_not_empty == 'true'` |

**Permissions:**

- `pull-requests`: `read`
- `contents`: `write`

<details>
<summary>Steps (5)</summary>

1. **checkout**
   - Uses: `actions/checkout@v4.2.2`

2. **checkout-pr**
   - Env:
     - `GH_TOKEN`: `${{ github.token }}`
     - `PR_NUMBER`: `${{ inputs.pr_number }}`

3. **download-patch**
   - Uses: `actions/download-artifact@v4.1.8`
   - With:
     - `name`: `changes.patch`

4. **apply**
   - ID: `apply`

5. **push**

</details>

[Back to top](#contents)

