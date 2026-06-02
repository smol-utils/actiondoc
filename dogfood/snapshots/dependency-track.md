# Contents

- [_meta-build.yaml](#_meta-buildyaml)
- [Build CI](#build-ci)
- [PR Template Check](#pr-template-check)
- [Publish CI](#publish-ci)
- [Release CI](#release-ci)
- [Report PR Test Coverage](#report-pr-test-coverage)
- [Tests CI](#tests-ci)
- [Dependency Review](#dependency-review)
- [Lock Threads](#lock-threads)

# _meta-build.yaml

| Property | Value |
|----------|-------|
| File | `_meta-build.yaml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `app-version` | string | No | `snapshot` | the version that should be set/used as tag for the container image |
| `publish-container` | boolean | No | `false` | publish and scan the container image once its built |
| `ref-name` | string | Yes | - | Short ref name of the branch or tag that triggered the workflow run |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `registry-0-usr` | Yes | - |
| `registry-0-psw` | Yes | - |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Called by

```
_meta-build.yaml
+-- ci-build.yaml (job: call-build)  <- entry point
+-- ci-publish.yaml (job: call-build)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `registry-0-usr` | job `build-container` step `Login to Docker.io` with `username` |
| `registry-0-psw` | job `build-container` step `Login to Docker.io` with `password` |

## Jobs

### `build-java`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Set up JDK**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `21`
     - `cache`: `maven`

3. **Setup CycloneDX CLI**

4. **Build with Maven**

5. **Upload Artifacts**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a`
   - With:
     - `name`: `assembled-wars`
     - `path`: `target/*.jar target/bom.json`

### `build-container`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-java` |

**Permissions:**

- `security-events`: `write` - Required to upload trivy's SARIF output

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Download Artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c`
   - With:
     - `name`: `assembled-wars`
     - `path`: `target`

3. **Set up QEMU**
   - Uses: `docker/setup-qemu-action@ce360397dd3f832beb865e1373c09c0e9f86d70a`

4. **Set up Docker Buildx**
   - ID: `buildx`
   - Uses: `docker/setup-buildx-action@4d04d5d9486b7bd6fa91e7baf45bbb4f8b9deedd`
   - With:
     - `install`: `true`

5. **Login to Docker.io**
   - Uses: `docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121`
   - Condition: `${{ inputs.publish-container }}`
   - With:
     - `registry`: `docker.io`
     - `username`: `${{ secrets.registry-0-usr }}`
     - `password`: `${{ secrets.registry-0-psw }}`

6. **Set Container Tags**
   - ID: `tags`
   - Env:
     - `REF_NAME`: `${{ inputs.ref-name }}`
     - `APP_VERSION`: `${{ inputs.app-version }}`
     - `DISTRIBUTION`: `${{ matrix.distribution }}`

7. **Build multi-arch Container Image**
   - Uses: `docker/build-push-action@bcafcacb16a39f128d818304e6c9c0c18556b85f`
   - With:
     - `tags`: `${{ steps.tags.outputs.tags }}`
     - `build-args`: `APP_VERSION=${{ inputs.app-version }} COMMIT_SHA=${{ github.sha }} WAR_FILENAME=dependency-track-${{ matrix.distribution }}.jar`
     - `platforms`: `linux/amd64,linux/arm64`
     - `push`: `${{ inputs.publish-container }}`
     - `context`: `.`
     - `file`: `src/main/docker/Dockerfile`

8. **Build Alpine multi-arch Container Image**
   - Uses: `docker/build-push-action@bcafcacb16a39f128d818304e6c9c0c18556b85f`
   - With:
     - `tags`: `${{ steps.tags.outputs.tags-alpine }}`
     - `build-args`: `APP_VERSION=${{ inputs.app-version }} COMMIT_SHA=${{ github.sha }} WAR_FILENAME=dependency-track-${{ matrix.distribution }}.jar`
     - `platforms`: `linux/amd64,linux/arm64`
     - `push`: `${{ inputs.publish-container }}`
     - `context`: `.`
     - `file`: `src/main/docker/Dockerfile.alpine`

# Build CI

| Property | Value |
|----------|-------|
| File | `ci-build.yaml` |
| Triggers | `push`, `pull_request`, `workflow_dispatch` |

## Event filters

- **push**
  - branches: `master`, `feature-**`, `[0-9]+.[0-9]+.x`
  - paths-ignore: `**/*.md`, `docs/**`
- **pull_request**
  - branches: `master`, `feature-**`
  - paths-ignore: `**/*.md`, `docs/**`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Call graph (rooted at this workflow)

```
ci-build.yaml [push, pull_request, workflow_dispatch]
+-- call-build (uses _meta-build.yaml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `HUB_ACCESSS_TOKEN`, `HUB_USERNAME`, `registry-0-psw`, `registry-0-usr`

Permissions declared across the chain: `security-events: write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HUB_USERNAME` | job `call-build` secrets `registry-0-usr` |
| `HUB_ACCESSS_TOKEN` | job `call-build` secrets `registry-0-psw` |

## Jobs

### `call-build`

| Property | Value |
|----------|-------|
| Uses workflow | [_meta-build.yaml](#_meta-buildyaml) |

**Permissions:**

- `security-events`: `write` - Required to upload trivy's SARIF output

#### Inputs forwarded

- `app-version`: `snapshot`
- `publish-container`: `${{ github.ref_name == 'master' || startsWith(github.ref_name, 'feature-') }}`
- `ref-name`: `${{ github.ref_name }}`

#### Secrets forwarded

- `registry-0-usr`: `${{ secrets.HUB_USERNAME }}`
- `registry-0-psw`: `${{ secrets.HUB_ACCESSS_TOKEN }}`

# PR Template Check

| Property | Value |
|----------|-------|
| File | `ci-pr-template.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - types: `opened`, `edited`, `synchronize`, `reopened`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `validate`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Validate PR body**
   - Env:
     - `PR_BODY`: `${{ github.event.pull_request.body }}`
     - `PR_AUTHOR_TYPE`: `${{ github.event.pull_request.user.type }}`
     - `PR_AUTHOR_LOGIN`: `${{ github.event.pull_request.user.login }}`

# Publish CI

This workflow is responsible to build and publish a release build It triggers once a new GitHub Release is published

| Property | Value |
|----------|-------|
| File | `ci-publish.yaml` |
| Triggers | `push`, `workflow_dispatch` |

## Event filters

- **push**
  - tags: `[0-9]*`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Call graph (rooted at this workflow)

```
ci-publish.yaml [push, workflow_dispatch]
+-- call-build (uses _meta-build.yaml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `GITHUB_TOKEN`, `HUB_ACCESSS_TOKEN`, `HUB_USERNAME`, `registry-0-psw`, `registry-0-usr`

Permissions declared across the chain: `contents: write`, `security-events: write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HUB_USERNAME` | job `call-build` secrets `registry-0-usr` |
| `HUB_ACCESSS_TOKEN` | job `call-build` secrets `registry-0-psw` |
| `GITHUB_TOKEN` | job `update-github-release` step `Update Release` env `GITHUB_TOKEN`; job `update-github-release` step `Publish Release` env `GITHUB_TOKEN` |

## Jobs

### `read-version`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Assert ref type**

2. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

3. **Parse Version from POM**
   - ID: `parse`

### `call-build`

| Property | Value |
|----------|-------|
| Uses workflow | [_meta-build.yaml](#_meta-buildyaml) |
| Depends on | `read-version` |

**Permissions:**

- `security-events`: `write` - Required to upload trivy's SARIF output

#### Inputs forwarded

- `app-version`: `${{ needs.read-version.outputs.version }}`
- `publish-container`: `true`
- `ref-name`: `${{ github.ref_name }}`

#### Secrets forwarded

- `registry-0-usr`: `${{ secrets.HUB_USERNAME }}`
- `registry-0-psw`: `${{ secrets.HUB_ACCESSS_TOKEN }}`

### `update-github-release`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `read-version`, `call-build` |

**Permissions:**

- `contents`: `write` - Required to update GitHub release assets and notes

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Download Artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c`
   - With:
     - `name`: `assembled-wars`
     - `path`: `target`

3. **Create Checksums and SBOM**

4. **Update Release**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `VERSION`: `${{ needs.read-version.outputs.version }}`

5. **Publish Release**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `VERSION`: `${{ needs.read-version.outputs.version }}`

# Release CI

| Property | Value |
|----------|-------|
| File | `ci-release.yaml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version-overwrite` | string | No | - | Use this to overwrite the version number to release, otherwise uses the current SNAPSHOT version (expected format x.y.z) |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `BOT_RELEASE_GITHUB_TOKEN` | job `create-release` step `Checkout Repository` with `token`; job `create-release` step `Create GitHub Release` env `GITHUB_TOKEN` |

## Jobs

### `prepare-release`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Setup Environment**
   - ID: `variables`
   - Env:
     - `VERSION_OVERWRITE`: `${{ github.event.inputs.version-overwrite }}`

### `create-release`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `prepare-release` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VERSION` | `${{ needs.prepare-release.outputs.version }}` |
| `BRANCH_NAME` | `${{ needs.prepare-release.outputs.release-branch }}` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `token`: `${{ secrets.BOT_RELEASE_GITHUB_TOKEN }}`

2. **Set up JDK**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `21`
     - `cache`: `maven`

3. **Set Version and Commit**

4. **Create GitHub Release**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.BOT_RELEASE_GITHUB_TOKEN }}`
     - `RELEASE_VERSION`: `${{ needs.prepare-release.outputs.version }}`
     - `RELEASE_BRANCH`: `${{ needs.prepare-release.outputs.release-branch }}`

### `post-release`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `prepare-release`, `create-release` |

**Permissions:**

- `contents`: `write` - Required to push pom.xml update

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `NEXT_VERSION` | `${{ needs.prepare-release.outputs.next-version }}` |
| `BRANCH_NAME` | `${{ needs.prepare-release.outputs.release-branch }}` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `ref`: `${{ needs.prepare-release.outputs.release-branch }}`

2. **Set up JDK**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `21`
     - `cache`: `maven`

3. **Set SNAPSHOT Version after Release**

# Report PR Test Coverage

| Property | Value |
|----------|-------|
| File | `ci-test-pr-coverage.yml` |
| Triggers | `workflow_run` |

## Event filters

- **workflow_run**
  - workflows: `Tests CI`
  - types: `completed`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `publish` step `Download PR test coverage report` with `github-token` |
| `CODACY_PROJECT_TOKEN` | job `publish` step `Report Coverage to Codacy` env `CODACY_PROJECT_TOKEN` |

## Jobs

### Report Coverage (`publish`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.event.workflow_run.event == 'pull_request'<br>  && github.event.workflow_run.conclusion == 'success'` |

#### Steps

1. **Download PR test coverage report**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c`
   - With:
     - `name`: `pr-test-coverage-report`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
     - `run-id`: `${{ github.event.workflow_run.id }}`

2. **Report Coverage to Codacy**
   - Env:
     - `CODACY_PROJECT_TOKEN`: `${{ secrets.CODACY_PROJECT_TOKEN }}`
     - `HEAD_SHA`: `${{ github.event.workflow_run.head_sha }}`

# Tests CI

| Property | Value |
|----------|-------|
| File | `ci-test.yaml` |
| Triggers | `push`, `pull_request`, `workflow_dispatch` |

## Event filters

- **push**
  - branches: `master`, `feature-**`, `[0-9]+.[0-9]+.x`
  - paths-ignore: `**/*.md`, `docs/**`
- **pull_request**
  - branches: `master`, `feature-**`, `[0-9]+.[0-9]+.x`
  - paths-ignore: `**/*.md`, `docs/**`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `${{ github.workflow }}-${{ github.head_ref || github.run_id }}`, cancel-in-progress: `true`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CODACY_PROJECT_TOKEN` | job `test` step `Publish test coverage` with `project-token` |

## Jobs

### `test`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Set up JDK**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `21`
     - `cache`: `maven`

3. **Execute unit tests**

4. **Publish test coverage**
   - Uses: `codacy/codacy-coverage-reporter-action@89d6c85cfafaec52c72b6c5e8b2878d33104c699`
   - Condition: `${{ github.event_name != 'pull_request' && github.repository_owner == 'DependencyTrack' }}`
   - With:
     - `project-token`: `${{ secrets.CODACY_PROJECT_TOKEN }}`
     - `language`: `Java`
     - `coverage-reports`: `target/jacoco-ut/jacoco.xml`

5. **Save PR details**
   - Condition: `${{ github.event_name == 'pull_request' }}`

6. **Upload PR test coverage report**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a`
   - Condition: `${{ github.event_name == 'pull_request' }}`
   - With:
     - `name`: `pr-test-coverage-report`
     - `path`: `pr-commit.txt pr-number.txt target/jacoco-ut/jacoco.xml`

# Dependency Review

| Property | Value |
|----------|-------|
| File | `dependency-review.yaml` |
| Triggers | `pull_request` |

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `dependency-review`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout Repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd`
   - With:
     - `persist-credentials`: `false`

2. **Dependency Review**
   - Uses: `actions/dependency-review-action@2031cfc080254a8a887f58cffee85186f0e49e48`

# Lock Threads

| Property | Value |
|----------|-------|
| File | `lock.yaml` |
| Triggers | `schedule` |

## Schedule

- `0 10 * * *`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

## Jobs

### `action`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `${{ contains(github.repository, 'DependencyTrack/') }}` |

**Permissions:**

- `issues`: `write` - Required to lock issues
- `pull-requests`: `write` - Required to lock PRs

#### Steps

1. **dessant/lock-threads**
   - Uses: `dessant/lock-threads@7266a7ce5c1df01b1c6db85bf8cd86c737dadbe7`
   - With:
     - `github-token`: `${{ github.token }}`
     - `issue-inactive-days`: `30`
     - `exclude-issue-created-before`: -
     - `exclude-any-issue-labels`: -
     - `add-issue-labels`: -
     - `issue-comment`: `This thread has been automatically locked since there has not been any recent activity after it was closed. Please open a new issue for related bugs.`
     - `issue-lock-reason`: `resolved`
     - `pr-inactive-days`: `30`
     - `exclude-pr-created-before`: -
     - `exclude-any-pr-labels`: -
     - `add-pr-labels`: -
     - `pr-comment`: -
     - `pr-lock-reason`: `resolved`
     - `process-only`: -

