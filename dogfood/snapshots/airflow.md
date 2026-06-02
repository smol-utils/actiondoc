# Contents

- [Additional CI image checks](#additional-ci-image-checks)
- [Additional PROD image tests](#additional-prod-image-tests)
- [Non-core Distribution tests](#non-core-distribution-tests)
- [Airflow E2E Tests](#airflow-e2e-tests)
- [ASF Allowlist Check](#asf-allowlist-check)
- [Automatic Backport](#automatic-backport)
- [Backport Commit](#backport-commit)
- [Basic tests](#basic-tests)
- [Check newsfragment PR number](#check-newsfragment-pr-number)
- [Tests (AMD)](#tests-amd)
- [Tests (ARM)](#tests-arm)
- [Build CI images](#build-ci-images)
- [CI Image Checks](#ci-image-checks)
- [CI Notification](#ci-notification)
- [CodeQL](#codeql)
- [E2E Flaky Tests Report](#e2e-flaky-tests-report)
- [Finalize tests](#finalize-tests)
- [Generate constraints](#generate-constraints)
- [Helm tests](#helm-tests)
- [Integration and system tests](#integration-and-system-tests)
- [K8s tests](#k8s-tests)
- [Milestone Tag Assistant](#milestone-tag-assistant)
- [Notify uv.lock conflicts](#notify-uvlock-conflicts)
- [Build PROD images](#build-prod-images)
- [PROD images extra checks](#prod-images-extra-checks)
- [Publish Docs to S3](#publish-docs-to-s3)
- [Push image cache](#push-image-cache)
- [Recheck old bug reports](#recheck-old-bug-reports)
- [Registry Backfill](#registry-backfill)
- [Build & Publish Registry](#build--publish-registry)
- [Registry Tests](#registry-tests)
- [Release PROD images](#release-prod-images)
- [Release single PROD image](#release-single-prod-image)
- [Unit tests](#unit-tests)
- [\[main\] Scheduled CI upgrade check](#main-scheduled-ci-upgrade-check)
- [\[v3-2-test\] Scheduled CI upgrade check](#v3-2-test-scheduled-ci-upgrade-check)
- [Scheduled verify release calendar](#scheduled-verify-release-calendar)
- [Special tests](#special-tests)
- [Close stale PRs & Issues](#close-stale-prs--issues)
- [Provider tests](#provider-tests)
- [UI End-to-End Tests](#ui-end-to-end-tests)
- [Update constraints on push for stable branch (always)](#update-constraints-on-push-for-stable-branch-always)
- [Update constraints on push for main (only when uv.lock changes)](#update-constraints-on-push-for-main-only-when-uvlock-changes)
- [Upgrade check](#upgrade-check)
- [Setup Breeze](#setup-breeze)
- [Install prek](#install-prek)
- [Run migration tests](#run-migration-tests)
- [Post tests on failure](#post-tests-on-failure)
- [Post tests on success](#post-tests-on-success)
- [Prepare all CI images](#prepare-all-ci-images)
- [Prepare breeze && current image (CI or PROD)](#prepare-breeze--current-image-ci-or-prod)
- [Prepare single CI image](#prepare-single-ci-image)

# Additional CI image checks

| Property | Value |
|----------|-------|
| File | `additional-ci-image-checks.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions` | string | Yes | - | The list of python versions (stringified JSON array) to run the tests on. |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `constraints-branch` | string | Yes | - | Branch used to get constraints from |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies (true/false) |
| `skip-prek-hooks` | string | Yes | - | Whether to skip prek hooks (true/false) |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |
| `canary-run` | string | Yes | - | Whether this is a canary run (true/false) |
| `latest-versions-only` | string | Yes | - | Whether to run only latest versions (true/false) |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs (true/false) |
| `debug-resources` | string | Yes | - | Whether to debug resources (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
additional-ci-image-checks.yml
+-- ci-amd.yml (job: additional-ci-image-checks)  <- entry point
+-- ci-arm.yml (job: additional-ci-image-checks)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `check-that-image-builds-quickly` env `GITHUB_TOKEN` |

## Jobs

### Push Early Image Cache (`push-early-buildx-cache-to-github-registry`)

| Property | Value |
|----------|-------|
| Uses workflow | [Push image cache](#push-image-cache) |
| Condition | `inputs.canary-run == 'true' && (github.event_name == 'schedule' \|\| github.event_name == 'workflow_dispatch')` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `cache-type`: `Early`
- `include-prod-images`: `false`
- `push-latest-images`: `false`
- `platform`: `${{ inputs.platform }}`
- `python-versions`: `${{ inputs.python-versions }}`
- `branch`: `${{ inputs.branch }}`
- `constraints-branch`: `${{ inputs.constraints-branch }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `docker-cache`: `${{ inputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ inputs.disable-airflow-repo-cache }}`

### Check that image builds quickly (`check-that-image-builds-quickly`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.branch == 'main'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `UPGRADE_TO_NEWER_DEPENDENCIES` | `false` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `PYTHON_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |
| `PLATFORM` | `${{ inputs.platform }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Check that image builds quickly**

# Additional PROD image tests

| Property | Value |
|----------|-------|
| File | `additional-prod-image-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-branch` | string | Yes | - | The default branch for the repository |
| `run-task-sdk-integration-tests` | string | Yes | - | Whether to run Task SDK integration tests (true/false) |
| `run-remote-logging-s3-e2e-tests` | string | Yes | - | Whether to run S3 remote logging e2e tests (true/false) |
| `run-remote-logging-elasticsearch-e2e-tests` | string | Yes | - | Whether to run Elasticsearch remote logging e2e tests (true/false) |
| `run-remote-logging-opensearch-e2e-tests` | string | Yes | - | Whether to run OpenSearch remote logging e2e tests (true/false) |
| `run-event-driven-e2e-tests` | string | Yes | - | Whether to run event driven e2e tests (true/false) |
| `constraints-branch` | string | Yes | - | Branch used to construct constraints URL from. |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies (true/false) |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |
| `canary-run` | string | Yes | - | Whether to run the canary run (true/false) |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `use-uv` | string | Yes | - | Whether to use uv |
| `run-ui-e2e-tests` | string | Yes | - | Whether to run UI e2e tests (true/false) |
| `run-airflow-ctl-integration-tests` | string | Yes | - | Whether to run Airflow CTL integration tests (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
additional-prod-image-tests.yml
+-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
+-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `test-examples-of-prod-image-building` env `GITHUB_TOKEN`; job `test-docker-compose-quick-start` env `GITHUB_TOKEN`; job `task-sdk-integration-tests` env `GITHUB_TOKEN`; job `airflow-ctl-integration-tests` env `GITHUB_TOKEN` |

## Jobs

### PROD image extra checks (main) (`prod-image-extra-checks-main`)

| Property | Value |
|----------|-------|
| Uses workflow | [PROD images extra checks](#prod-images-extra-checks) |
| Condition | `inputs.default-branch == 'main' && inputs.canary-run == 'true'` |

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `python-versions`: `[ '${{ inputs.default-python-version }}' ]`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `branch`: `${{ inputs.default-branch }}`
- `upgrade-to-newer-dependencies`: `${{ inputs.upgrade-to-newer-dependencies }}`
- `constraints-branch`: `${{ inputs.constraints-branch }}`
- `docker-cache`: `${{ inputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ inputs.disable-airflow-repo-cache }}`

### PROD image extra checks (release) (`prod-image-extra-checks-release-branch`)

| Property | Value |
|----------|-------|
| Uses workflow | [PROD images extra checks](#prod-images-extra-checks) |
| Condition | `inputs.default-branch != 'main' && inputs.canary-run == 'true'` |

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `python-versions`: `[ '${{ inputs.default-python-version }}' ]`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `branch`: `${{ inputs.default-branch }}`
- `upgrade-to-newer-dependencies`: `${{ inputs.upgrade-to-newer-dependencies }}`
- `constraints-branch`: `${{ inputs.constraints-branch }}`
- `docker-cache`: `${{ inputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ inputs.disable-airflow-repo-cache }}`

### Test examples of PROD image building (`test-examples-of-prod-image-building`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Test examples of PROD image building**

### Docker Compose quick start with PROD image verifying (`test-docker-compose-quick-start`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Test docker-compose quick start**

### Task SDK integration tests with PROD image (`task-sdk-integration-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-task-sdk-integration-tests == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Run Task SDK integration tests**

### Test e2e integration tests with PROD image (`test-e2e-integration-tests-basic`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |

#### Inputs forwarded

- `workflow-name`: `Regular e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`

### Remote logging tests with PROD image (`test-e2e-integration-tests-remote-log`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |
| Condition | `inputs.canary-run == 'true' \|\| inputs.run-remote-logging-s3-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `Remote logging e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `e2e_test_mode`: `remote_log`

### Elasticsearch remote logging tests with PROD image (`test-e2e-integration-tests-remote-log-elasticsearch`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |
| Condition | `inputs.canary-run == 'true' \|\| inputs.run-remote-logging-elasticsearch-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `Elasticsearch remote logging e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `e2e_test_mode`: `remote_log_elasticsearch`

### OpenSearch remote logging tests with PROD image (`test-e2e-integration-tests-remote-log-opensearch`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |
| Condition | `inputs.canary-run == 'true' \|\| inputs.run-remote-logging-opensearch-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `OpenSearch remote logging e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `e2e_test_mode`: `remote_log_opensearch`

### XCom object storage backend tests with PROD image (`test-e2e-integration-tests-xcom-object-storage`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |

#### Inputs forwarded

- `workflow-name`: `XCom object storage backend e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `e2e_test_mode`: `xcom_object_storage`

### Event driven tests with PROD image (`test-e2e-integration-tests-event-driven`)

| Property | Value |
|----------|-------|
| Uses workflow | [Airflow E2E Tests](#airflow-e2e-tests) |
| Condition | `inputs.canary-run == 'true' \|\| inputs.run-event-driven-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `Event driven e2e test`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `e2e_test_mode`: `event_driven`

### Chromium UI e2e tests with PROD image (`test-ui-e2e-chromium`)

| Property | Value |
|----------|-------|
| Uses workflow | [UI End-to-End Tests](#ui-end-to-end-tests) |
| Condition | `inputs.run-ui-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `Chromium UI e2e tests`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `browser`: `chromium`

### Firefox UI e2e tests with PROD image (`test-ui-e2e-firefox`)

| Property | Value |
|----------|-------|
| Uses workflow | [UI End-to-End Tests](#ui-end-to-end-tests) |
| Condition | `inputs.run-ui-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `Firefox UI e2e tests`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `browser`: `firefox`

### WebKit UI e2e tests with PROD image (`test-ui-e2e-webkit`)

| Property | Value |
|----------|-------|
| Uses workflow | [UI End-to-End Tests](#ui-end-to-end-tests) |
| Condition | `inputs.run-ui-e2e-tests == 'true'` |

#### Inputs forwarded

- `workflow-name`: `WebKit UI e2e tests`
- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `browser`: `webkit`

### Airflow CTL integration tests with PROD image (`airflow-ctl-integration-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-airflow-ctl-integration-tests == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Run airflowctl integration tests**

# Non-core Distribution tests

| Property | Value |
|----------|-------|
| File | `airflow-distributions-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `distribution-name` | string | Yes | - | The name of the distribution to test |
| `distribution-cmd-format` | string | Yes | - | The type of distribution to test |
| `test-type` | string | Yes | - | distribution test type |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to build images from |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `canary-run` | string | Yes | - | Whether this is a canary run (true/false) |
| `use-local-venv` | string | Yes | - | Whether local venv should be used for tests (true/false) |
| `test-timeout` | number | No | `60` | - |

## Permissions

- `contents`: `read`

## Called by

```
airflow-distributions-tests.yml
+-- ci-amd.yml (job: tests-task-sdk)  <- entry point
+-- ci-amd.yml (job: tests-airflow-ctl)  <- entry point
+-- ci-arm.yml (job: tests-task-sdk)  <- entry point
+-- ci-arm.yml (job: tests-airflow-ctl)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `distributions-tests` env `GITHUB_TOKEN` |

## Jobs

### ${{ inputs.distribution-name }}:P${{ matrix.python-version }} tests (`distributions-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ matrix.python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - Condition: `${{ inputs.use-local-venv != 'true' }}`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ matrix.python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`
   - Condition: `${{ inputs.use-local-venv == 'true' }}`

5. **Cleanup dist files**
   - Condition: `${{ matrix.python-version == inputs.default-python-version }}`

6. **Prepare Airflow ${{inputs.distribution-name}}: wheel**
   - Condition: `${{ matrix.python-version == inputs.default-python-version }}`

7. **Verify wheel packages with twine**
   - Condition: `${{ matrix.python-version == inputs.default-python-version }}`

8. **Run unit tests for Airflow ${{inputs.distribution-name}}:Python ${{ matrix.python-version }}
**

# Airflow E2E Tests

| Property | Value |
|----------|-------|
| File | `airflow-e2e-tests.yml` |
| Triggers | `workflow_dispatch`, `workflow_call` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow-name` | string | Yes | - | Name of the test |
| `runners` | string | No | `["ubuntu-24.04"]` | The array of labels (in json form) determining runners. |
| `platform` | string | No | `linux/amd64` | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-python-version` | string | No | `3.10` | Which version of python should be used by default |
| `use-uv` | string | No | `true` | Whether to use uv to build the image (true/false) |
| `docker-image-tag` | string | Yes | - | Tag of the Docker image to test |
| `e2e_test_mode` | string | No | `basic` | Test mode - basic, remote_log, remote_log_elasticsearch, remote_log_opensearch, xcom_object_storage, or event_driven |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow-name` | string | Yes | - | Name of the test |
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `docker-image-tag` | string | No | - | Tag of the Docker image to test |
| `e2e_test_mode` | string | No | `basic` | Test mode - basic, remote_log, remote_log_elasticsearch, remote_log_opensearch, xcom_object_storage, or event_driven |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
airflow-e2e-tests.yml [workflow_dispatch, workflow_call]
+-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
```

## Called by

```
airflow-e2e-tests.yml
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-basic)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-remote-log)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-remote-log-elasticsearch)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-remote-log-opensearch)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-xcom-object-storage)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-e2e-integration-tests-event-driven)
    +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
    +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `test-e2e-integration-tests` env `GITHUB_TOKEN` |

## Jobs

### ${{ inputs.workflow-name }} (`test-e2e-integration-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Test e2e integration tests**

5. **Zip logs**
   - Condition: `always()`

6. **Upload logs**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `always()`
   - With:
     - `name`: `e2e-test-logs-${{ inputs.e2e_test_mode }}`
     - `path`: `./airflow-e2e-tests/logs.zip`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

# ASF Allowlist Check

| Property | Value |
|----------|-------|
| File | `asf-allowlist-check.yml` |
| Triggers | `pull_request`, `push` |

## Event filters

- **pull_request**
  - paths: `.github/**`
- **push**
  - branches: `main`, `v*-test`
  - paths: `.github/**`

## Permissions

- `contents`: `read`

## Jobs

### `asf-allowlist-check`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **apache/infrastructure-actions/allowlist-check**
   - Uses: `apache/infrastructure-actions/allowlist-check@4e9c961f587f72b170874b6f5cd4ac15f7f26eb8`

# Automatic Backport

| Property | Value |
|----------|-------|
| File | `automatic-backport.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
automatic-backport.yml [push]
+-- trigger-backport (uses backport-cli.yml)
```

## Jobs

### Get PR information (`get-pr-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Get commit SHA**
   - ID: `get-sha`

2. **Add delay for GitHub to process PR merge**

3. **Find PR information**
   - ID: `pr-info`
   - Uses: `actions/github-script@3a2844b7e9c422d3c10d287c895573f7108da1b3` (v9.0.0)
   - With:
     - `script`: `const { data: pullRequest } = await github.rest.repos.listPullRequestsAssociatedWithCommit({     owner: context.repo.owner,     repo: context.repo.repo,     commit_sha: process.env.GITHUB_SHA }); if (pullRequest.length > 0) {     const pr = pullRequest[0];     const backportBranches = pr.labels           .filter(label => label.name.startsWith('backport-to-'))           .map(label => label.name.replace('backport-to-', ''));      console.log(`Commit ${process.env.GITHUB_SHA} is associated with PR ${pr.number}`);     console.log(`Backport branches: ${backportBranches}`);     core.setOutput('branches', JSON.stringify(backportBranches)); } else {     console.log('⚠️ No pull request found for this commit.');     core.setOutput('branches', '[]'); }`

### Trigger Backport (`trigger-backport`)

| Property | Value |
|----------|-------|
| Uses workflow | [Backport Commit](#backport-commit) |
| Depends on | `get-pr-info` |
| Condition | `${{ needs.get-pr-info.outputs.branches != '[]' }}` |

**Permissions:**

- `contents`: `write`
- `pull-requests`: `write`

#### Inputs forwarded

- `target-branch`: `${{ matrix.branch }}`
- `commit-sha`: `${{ needs.get-pr-info.outputs.commit-sha }}`

# Backport Commit

| Property | Value |
|----------|-------|
| File | `backport-cli.yml` |
| Triggers | `workflow_dispatch`, `workflow_call` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `commit-sha` | string | Yes | - | Commit sha to backport. |
| `target-branch` | string | Yes | - | Target branch to backport. |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `commit-sha` | string | Yes | - | Commit sha to backport. |
| `target-branch` | string | Yes | - | Target branch to backport. |

## Permissions

- `contents`: `write` - Those permissions are only active for workflow dispatch (only committers can trigger it) and workflow call Which is triggered automatically by "automatic-backport" push workflow (only when merging by committer) Branch protection  prevents from pushing to the "code" branches
- `pull-requests`: `write`

## Called by

```
backport-cli.yml
+-- automatic-backport.yml (job: trigger-backport)  <- entry point
```

## Jobs

### `backport`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - ID: `checkout-for-backport`
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

2. **Install Python dependencies**

3. **Run backport script** `[continue-on-error]`
   - ID: `execute-backport`

4. **Parse backport output** `[continue-on-error]`
   - ID: `parse-backport-output`

5. **Update Status**
   - ID: `backport-status`

# Basic tests

| Property | Value |
|----------|-------|
| File | `basic-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `run-ui-tests` | string | Yes | - | Whether to run UI tests (true/false) |
| `run-www-tests` | string | Yes | - | Whether to run WWW tests (true/false) |
| `run-api-codegen` | string | Yes | - | Whether to run API codegen (true/false) |
| `run-breeze-integration-tests` | string | Yes | - | Whether to run breeze integration tests (true/false) |
| `run-scripts-tests` | string | Yes | - | Whether to run scripts tests (true/false) |
| `basic-checks-only` | string | Yes | - | Whether to run only basic checks (true/false) |
| `skip-prek-hooks` | string | Yes | - | Whether to skip prek hooks (true/false) |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `shared-distributions-as-json` | string | Yes | - | Json array of shared distributions to run tests for |
| `canary-run` | string | Yes | - | Whether to run canary tests (true/false) |
| `latest-versions-only` | string | Yes | - | Whether to run only latest version checks (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv in the image |
| `platform` | string | Yes | - | Platform for the build - linux/amd64 or linux/arm64 |

## Permissions

- `contents`: `read`

## Called by

```
basic-tests.yml
+-- ci-amd.yml (job: basic-tests)  <- entry point
+-- ci-arm.yml (job: basic-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `test-airflow-release-commands` env `GITHUB_TOKEN` |

## Jobs

### Breeze unit tests (`run-breeze-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

#### Steps

1. **Cleanup repo**

2. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Run unit tests**

### Breeze integration tests (`run-breeze-integration-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-breeze-integration-tests == 'true'` |

#### Steps

1. **Cleanup repo**

2. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Install SVN**

5. **Install Java (for Apache RAT)**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`

6. **Install hatch**

7. **Run integration tests**

### Shared ${{ matrix.shared-distribution }} tests (`tests-shared-distributions`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `1`
     - `persist-credentials`: `false`

2. **Install uv**

3. **Run shared ${{ matrix.shared-distribution }} tests**

### Scripts tests (`tests-scripts`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-scripts-tests == 'true'` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `1`
     - `persist-credentials`: `false`

2. **Install uv**

3. **Run scripts tests**

### React UI tests (`tests-ui`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-ui-tests == 'true'` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Setup pnpm**
   - Uses: `pnpm/action-setup@0e279bb959325dab635dd2c09392533439d90093` (v6.0.8)
   - With:
     - `version`: `9`
     - `run_install`: `false`

4. **Setup node**
   - Uses: `actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e` (v6.4.0)
   - With:
     - `node-version`: `24`
     - `cache`: `pnpm`
     - `cache-dependency-path`: `airflow-core/src/airflow/**/pnpm-lock.yaml`

5. **Restore eslint cache (ui)**
   - ID: `restore-eslint-cache-ui`
   - Uses: `apache/infrastructure-actions/stash/restore@49df447b39b18354895520e0a63731b7cad7cbec`
   - With:
     - `path`: `airflow-core/src/airflow/ui/node_modules/`
     - `key`: `cache-ui-node-modules-v1-${{ runner.os }}-${{ hashFiles('airflow-core/src/airflow/ui/**/pnpm-lock.yaml') }}`

6. **cd airflow-core/src/airflow/ui && pnpm install --frozen-l...**

7. **cd airflow-core/src/airflow/ui && pnpm test**

8. **Save eslint cache (ui)**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `steps.restore-eslint-cache-ui.outputs.stash-hit != 'true'`
   - With:
     - `path`: `airflow-core/src/airflow/ui/node_modules/`
     - `key`: `cache-ui-node-modules-v1-${{ runner.os }}-${{ hashFiles('airflow/ui/**/pnpm-lock.yaml') }}`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

9. **Restore eslint cache (simple auth manager UI)**
   - ID: `restore-eslint-cache-simple-am-ui`
   - Uses: `apache/infrastructure-actions/stash/restore@49df447b39b18354895520e0a63731b7cad7cbec`
   - With:
     - `path`: `airflow-core/src/airflow/api_fastapi/auth/managers/simple/ui/node_modules/`
     - `key`: `cache-simple-am-ui-node-modules-v1- ${{ runner.os }}-${{ hashFiles('airflow/api_fastapi/auth/managers/simple/ui/**/pnpm-lock.yaml') }}`

10. **cd airflow-core/src/airflow/api_fastapi/auth/managers/sim...**

11. **cd airflow-core/src/airflow/api_fastapi/auth/managers/sim...**

12. **Save eslint cache (ui)**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `steps.restore-eslint-cache-simple-am-ui.outputs.stash-hit != 'true'`
   - With:
     - `path`: `airflow-core/src/airflow/api_fastapi/auth/managers/simple/ui/node_modules/`
     - `key`: `cache-simple-am-ui-node-modules-v1- ${{ runner.os }}-${{ hashFiles('airflow/api_fastapi/auth/managers/simple/ui/**/pnpm-lock.yaml') }}`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

### Check translation completeness (`check-translation-completness`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install Breeze**
   - Uses: `./.github/actions/breeze`

3. **Check translation completeness**

### Static checks: basic checks only (`static-checks-basic-checks-only`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.basic-checks-only == 'true'` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ steps.breeze.outputs.host-python-version }}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `true` - Whether to save prek cache (required)

5. **Fetch incoming commit ${{ github.sha }} with its parent**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ github.sha }}`
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

6. **Static checks: basic checks only**

### Test git clone on Windows (`test-git-clone-on-windows`)

| Property | Value |
|----------|-------|
| Runs on | `windows-2025` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

### Test Airflow release commands (`test-airflow-release-commands`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.canary-run == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Cleanup dist files**

5. **Setup git for tagging**

6. **Install twine**

7. **Check Airflow create minor branch command**

8. **Check Airflow RC process command**

9. **Check Airflow release process command**

10. **Test providers metadata generation**

11. **Fetch all git tags for origin**

12. **Test airflow core issue generation automatically**

### Test Airflow standalone commands (`test-airflow-standalone`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `AIRFLOW_HOME` | `~/airflow` |
| `FORCE_COLOR` | `1` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install uv**

3. **Set up Airflow home directory**

4. **Install Airflow from current repo (simulating user installation)**

5. **Test airflow standalone command**

# Check newsfragment PR number

| Property | Value |
|----------|-------|
| File | `check-newsfragment-pr-number.yml` |
| Triggers | `pull_request` |

## Event filters

- **pull_request**
  - branches: `main`
  - types: `opened`, `reopened`, `synchronize`

## Permissions

- `contents`: `read`
- `pull-requests`: `read`

**Concurrency:** group `check-newsfragment-${{ github.event.pull_request.number }}`, cancel-in-progress: `false`

## Jobs

### `check-newsfragment-pr-number`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `${{ !contains(github.event.pull_request.labels.*.name, 'skip newsfragment check') }}` |

#### Steps

1. **Check newsfragment PR number**

# Tests (AMD)

| Property | Value |
|----------|-------|
| File | `ci-amd.yml` |
| Triggers | `schedule`, `pull_request`, `push`, `workflow_dispatch` |

## Schedule

- `58 1,7,13,19 * * *`

## Event filters

- **pull_request**
  - branches: `main`, `v[0-9]+-[0-9]+-test`, `v[0-9]+-[0-9]+-stable`, `providers-[a-z]+-?[a-z]*/v[0-9]+-[0-9]+`
  - types: `opened`, `reopened`, `synchronize`, `ready_for_review`
- **push**
  - branches: `v[0-9]+-[0-9]+-test`, `providers-[a-z]+-?[a-z]*/v[0-9]+-[0-9]+`

## Permissions

- `contents`: `read` - All other permissions are set to none by default

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |
| `VERBOSE` | `true` |

**Concurrency:** group `ci-amd-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
ci-amd.yml [schedule, pull_request, push, workflow_dispatch]
+-- build-info / Install Breeze (uses ./.github/actions/breeze)
+-- basic-tests (uses basic-tests.yml)
|   +-- run-breeze-tests / Install Breeze (uses ./.github/actions/breeze)
|   +-- run-breeze-integration-tests / Install Breeze (uses ./.github/actions/breeze)
|   +-- check-translation-completness / Install Breeze (uses ./.github/actions/breeze)
|   +-- static-checks-basic-checks-only / Install Breeze (uses ./.github/actions/breeze)
|   +-- static-checks-basic-checks-only / Install prek (uses ./.github/actions/install-prek)
|   +-- test-airflow-release-commands / Install Breeze (uses ./.github/actions/breeze)
+-- build-ci-images (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- additional-ci-image-checks (uses additional-ci-image-checks.yml)
|   +-- push-early-buildx-cache-to-github-registry (uses push-image-cache.yml)
|   |   +-- push-ci-image-cache / Install Breeze (uses ./.github/actions/breeze)
|   |   +-- push-prod-image-cache / Install Breeze (uses ./.github/actions/breeze)
|   +-- check-that-image-builds-quickly / Install Breeze (uses ./.github/actions/breeze)
+-- generate-constraints (uses generate-constraints.yml)
|   +-- generate-constraints-matrix / Install prek (uses ./.github/actions/install-prek)
|   +-- generate-constraints-matrix / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- ci-image-checks (uses ci-image-checks.yml)
|   +-- static-checks / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- static-checks / Install prek (uses ./.github/actions/install-prek)
|   +-- build-docs / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- publish-docs / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-python-api-client / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- mypy-providers / Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- mypy-providers / Install prek (uses ./.github/actions/install-prek)
+-- migration-round-trip / Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- migration-round-trip / Install prek (uses ./.github/actions/install-prek)
+-- providers (uses test-providers.yml)
|   +-- prepare-install-verify-provider-distributions / Install prek (uses ./.github/actions/install-prek)
|   +-- prepare-install-verify-provider-distributions / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- providers-compatibility-tests-matrix / Install prek (uses ./.github/actions/install-prek)
|   +-- providers-compatibility-tests-matrix / Prepare breeze & CI image: ${{ matrix.compat.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-helm (uses helm-tests.yml)
|   +-- tests-helm / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-helm-release / Install Breeze (uses ./.github/actions/breeze)
+-- tests-postgres-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-postgres-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-mysql-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-mysql-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-sqlite-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-sqlite-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-non-db-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-non-db-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-special (uses special-tests.yml)
|   +-- tests-min-sqlalchemy (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-min-sqlalchemy-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-latest-sqlalchemy (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-latest-sqlalchemy-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-boto-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-boto-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-pendulum-2-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-pendulum-2-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-quarantined-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-quarantined-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-system-core (uses run-unit-tests.yml)
|       +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|       +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|       +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|       +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-integration-system (uses integration-system-tests.yml)
|   +-- tests-core-integration / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-core-integration / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-core-integration / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-providers-integration / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-providers-integration / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-providers-integration / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-system / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-system / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-system / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-with-lowest-direct-resolution-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-with-lowest-direct-resolution-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- build-prod-images (uses prod-image-build.yml)
|   +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
+-- additional-prod-image-tests (uses additional-prod-image-tests.yml)
|   +-- prod-image-extra-checks-main (uses prod-image-extra-checks.yml)
|   |   +-- pip-image (uses prod-image-build.yml)
|   |       +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   |       +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   |       +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
|   +-- prod-image-extra-checks-release-branch (uses prod-image-extra-checks.yml)
|   |   +-- pip-image (uses prod-image-build.yml)
|   |       +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   |       +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   |       +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
|   +-- test-examples-of-prod-image-building / Prepare breeze & PROD image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-docker-compose-quick-start / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- task-sdk-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-basic (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log-elasticsearch (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log-opensearch (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-xcom-object-storage (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-event-driven (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-ui-e2e-chromium (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- test-ui-e2e-firefox (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- test-ui-e2e-webkit (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- airflow-ctl-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-kubernetes (uses k8s-tests.yml)
|   +-- tests-kubernetes / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-task-sdk (uses airflow-distributions-tests.yml)
|   +-- distributions-tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- distributions-tests / Install Breeze (uses ./.github/actions/breeze)
+-- tests-airflow-ctl (uses airflow-distributions-tests.yml)
|   +-- distributions-tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- distributions-tests / Install Breeze (uses ./.github/actions/breeze)
+-- finalize-tests (uses finalize-tests.yml)
    +-- dependency-upgrade-summary / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
    +-- push-buildx-cache-to-github-registry (uses push-image-cache.yml)
        +-- push-ci-image-cache / Install Breeze (uses ./.github/actions/breeze)
        +-- push-prod-image-cache / Install Breeze (uses ./.github/actions/breeze)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `DOCS_AWS_ACCESS_KEY_ID`, `DOCS_AWS_SECRET_ACCESS_KEY`, `SLACK_BOT_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN`; job `mypy-providers` env `GITHUB_TOKEN`; job `migration-round-trip` env `GITHUB_TOKEN`; job `tests-go-sdk` env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | workflow env `SLACK_BOT_TOKEN`; job `ci-image-checks` secrets `SLACK_BOT_TOKEN` |
| `DOCS_AWS_ACCESS_KEY_ID` | job `ci-image-checks` secrets `DOCS_AWS_ACCESS_KEY_ID` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `ci-image-checks` secrets `DOCS_AWS_SECRET_ACCESS_KEY` |

## Jobs

### Build info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Fetch incoming commit ${{ github.sha }} with its parent**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ github.sha }}`
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

4. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

5. **Save github context to file**

6. **Get information about the Workflow**
   - ID: `source-run-info`

7. **Selective checks**
   - ID: `selective-checks`

8. **env**

### Platform: AMD (`print-platform`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info` |

#### Steps

1. **Print architecture**

### Basic tests (`basic-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Basic tests](#basic-tests) |
| Depends on | `build-info` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `run-ui-tests`: `${{needs.build-info.outputs.run-ui-tests}}`
- `run-www-tests`: `${{needs.build-info.outputs.run-www-tests}}`
- `run-api-codegen`: `${{needs.build-info.outputs.run-api-codegen}}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `basic-checks-only`: `${{ needs.build-info.outputs.basic-checks-only }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `canary-run`: `${{needs.build-info.outputs.canary-run}}`
- `run-breeze-integration-tests`: `${{needs.build-info.outputs.run-breeze-integration-tests}}`
- `run-scripts-tests`: `${{needs.build-info.outputs.run-scripts-tests}}`
- `latest-versions-only`: `${{needs.build-info.outputs.latest-versions-only}}`
- `use-uv`: `${{needs.build-info.outputs.use-uv}}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `shared-distributions-as-json`: `${{needs.build-info.outputs.shared-distributions-as-json}}`

### Build CI images (`build-ci-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.ci-image-build == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `${{ needs.build-info.outputs.canary-run }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`

### Additional CI image checks (`additional-ci-image-checks`)

| Property | Value |
|----------|-------|
| Uses workflow | [Additional CI image checks](#additional-ci-image-checks) |
| Depends on | `build-info`, `build-ci-images` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`
- `id-token`: `write` (OIDC)

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `latest-versions-only`: `${{ needs.build-info.outputs.latest-versions-only }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Generate constraints (`generate-constraints`)

| Property | Value |
|----------|-------|
| Uses workflow | [Generate constraints](#generate-constraints) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.ci-image-build == 'true'` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `generate-pypi-constraints`: `true`
- `generate-no-providers-constraints`: `${{ needs.build-info.outputs.canary-run }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### CI image checks (`ci-image-checks`)

| Property | Value |
|----------|-------|
| Uses workflow | [CI Image Checks](#ci-image-checks) |
| Depends on | `build-info`, `build-ci-images` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `docs-list-as-string`: `${{ needs.build-info.outputs.docs-list-as-string }}`
- `latest-versions-only`: `${{ needs.build-info.outputs.latest-versions-only }}`
- `basic-checks-only`: `${{ needs.build-info.outputs.basic-checks-only }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `ci-image-build`: `${{ needs.build-info.outputs.ci-image-build }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `docs-build`: `${{ needs.build-info.outputs.docs-build }}`
- `run-api-codegen`: `${{ needs.build-info.outputs.run-api-codegen }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `source-head-repo`: `${{ needs.build-info.outputs.source-head-repo }}`
- `source-head-ref`: `${{ needs.build-info.outputs.source-head-ref }}`

#### Secrets forwarded

- `DOCS_AWS_ACCESS_KEY_ID`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
- `DOCS_AWS_SECRET_ACCESS_KEY`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
- `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_BOT_TOKEN }}`

### MyPy providers checks (`mypy-providers`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-mypy-providers == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ needs.build-info.outputs.default-python-version }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ needs.build-info.outputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ needs.build-info.outputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{steps.breeze.outputs.host-python-version}}` - Python version to use
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

5. **MyPy checks for providers**

### Migration round-trip check (`migration-round-trip`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.has-migrations == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ needs.build-info.outputs.default-python-version }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ needs.build-info.outputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ needs.build-info.outputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{steps.breeze.outputs.host-python-version}}` - Python version to use
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

5. **Migration round-trip check**

### provider distributions tests (`providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Provider tests](#provider-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.skip-providers-tests != 'true' && needs.build-info.outputs.latest-versions-only != 'true' && needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `selected-providers-list-as-string`: `${{ needs.build-info.outputs.selected-providers-list-as-string }}`
- `providers-compatibility-tests-matrix`: `${{ needs.build-info.outputs.providers-compatibility-tests-matrix }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `providers-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Helm tests (`tests-helm`)

| Property | Value |
|----------|-------|
| Uses workflow | [Helm tests](#helm-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-helm-tests == 'true' && needs.build-info.outputs.default-branch == 'main' && needs.build-info.outputs.latest-versions-only != 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `helm-test-packages`: `${{ needs.build-info.outputs.helm-test-packages }}`
- `helm-test-kubernetes-versions`: `${{ needs.build-info.outputs.helm-test-kubernetes-versions }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Postgres tests: core (`tests-postgres-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `postgres`
- `test-name`: `Postgres`
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.postgres-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.postgres-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-migration-tests`: `true`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Postgres tests: providers (`tests-postgres-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `postgres`
- `test-name`: `Postgres`
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.postgres-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.postgres-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-migration-tests`: `true`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### MySQL tests: core (`tests-mysql-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && needs.build-info.outputs.platform == 'linux/amd64'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `mysql`
- `test-name`: `MySQL`
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.mysql-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.mysql-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### MySQL tests: providers (`tests-mysql-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && needs.build-info.outputs.platform == 'linux/amd64'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `mysql`
- `test-name`: `MySQL`
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.mysql-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.mysql-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Sqlite tests: core (`tests-sqlite-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: `Sqlite`
- `test-name-separator`: ``
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Sqlite tests: providers (`tests-sqlite-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: `Sqlite`
- `test-name-separator`: ``
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Non-DB tests: core (`tests-non-db-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: ``
- `test-name-separator`: ``
- `test-scope`: `Non-DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Non-DB tests: providers (`tests-non-db-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: ``
- `test-name-separator`: ``
- `test-scope`: `Non-DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Special tests (`tests-special`)

| Property | Value |
|----------|-------|
| Uses workflow | [Special tests](#special-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && (needs.build-info.outputs.canary-run == 'true' \|\|<br> needs.build-info.outputs.upgrade-to-newer-dependencies != 'false' \|\|<br> needs.build-info.outputs.full-tests-needed == 'true')` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`
- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `core-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `providers-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Integration and System Tests (`tests-integration-system`)

| Property | Value |
|----------|-------|
| Uses workflow | [Integration and system tests](#integration-and-system-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `testable-core-integrations`: `${{ needs.build-info.outputs.testable-core-integrations }}`
- `testable-providers-integrations`: `${{ needs.build-info.outputs.testable-providers-integrations }}`
- `run-system-tests`: `${{ needs.build-info.outputs.run-system-tests }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `default-mysql-version`: `${{ needs.build-info.outputs.default-mysql-version }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Low dep tests:core (`tests-with-lowest-direct-resolution-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `test-name`: `LowestDeps`
- `force-lowest-dependencies`: `true`
- `test-scope`: `All`
- `test-group`: `core`
- `backend`: `sqlite`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['${{ needs.build-info.outputs.default-postgres-version }}']`
- `excluded-providers-as-string`: ``
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `monitor-delay-time-in-seconds`: `120`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Low dep tests: providers (`tests-with-lowest-direct-resolution-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `test-name`: `LowestDeps`
- `force-lowest-dependencies`: `true`
- `test-scope`: `All`
- `test-group`: `providers`
- `backend`: `sqlite`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['${{ needs.build-info.outputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.individual-providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `monitor-delay-time-in-seconds`: `120`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Build PROD images (`build-prod-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build PROD images](#build-prod-images) |
| Depends on | `build-info`, `build-ci-images`, `generate-constraints` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `build-type`: `Regular`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-package-artifact`: `true`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `prod-image-build`: `${{ needs.build-info.outputs.prod-image-build }}`

### Additional PROD image tests (`additional-prod-image-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Additional PROD image tests](#additional-prod-image-tests) |
| Depends on | `build-info`, `build-prod-images`, `generate-constraints` |
| Condition | `needs.build-info.outputs.prod-image-build == 'true'` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `run-task-sdk-integration-tests`: `${{ needs.build-info.outputs.run-task-sdk-integration-tests }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `run-remote-logging-elasticsearch-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-elasticsearch-e2e-tests }}`
- `run-remote-logging-opensearch-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-opensearch-e2e-tests }}`
- `run-remote-logging-s3-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-s3-e2e-tests }}`
- `run-event-driven-e2e-tests`: `${{ needs.build-info.outputs.run-event-driven-e2e-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `run-ui-e2e-tests`: `${{ needs.build-info.outputs.run-ui-e2e-tests }}`
- `run-airflow-ctl-integration-tests`: `${{ needs.build-info.outputs.run-airflow-ctl-integration-tests }}`

### Kubernetes tests (`tests-kubernetes`)

| Property | Value |
|----------|-------|
| Uses workflow | [K8s tests](#k8s-tests) |
| Depends on | `build-info`, `build-prod-images` |
| Condition | `( needs.build-info.outputs.run-kubernetes-tests == 'true' \|\| needs.build-info.outputs.run-helm-tests == 'true')` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `kubernetes-combos`: `${{ needs.build-info.outputs.kubernetes-combos }}`

### Task SDK tests (`tests-task-sdk`)

| Property | Value |
|----------|-------|
| Uses workflow | [Non-core Distribution tests](#non-core-distribution-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-task-sdk-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `distribution-name`: `task-sdk`
- `distribution-cmd-format`: `prepare-task-sdk-distributions`
- `test-type`: `task-sdk-tests`
- `use-local-venv`: `false`
- `test-timeout`: `20`

### Go SDK tests (`tests-go-sdk`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.run-go-sdk-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Setup Go**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `1.24`
     - `cache-dependency-path`: `go-sdk/go.sum`

3. **Setup Gotestsum**

4. **Cleanup dist files**

5. **Run Go tests**

### Airflow CTL tests (`tests-airflow-ctl`)

| Property | Value |
|----------|-------|
| Uses workflow | [Non-core Distribution tests](#non-core-distribution-tests) |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.run-airflow-ctl-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `distribution-name`: `airflow-ctl`
- `distribution-cmd-format`: `prepare-airflow-ctl-distributions`
- `test-type`: `airflow-ctl-tests`
- `use-local-venv`: `true`
- `test-timeout`: `20`

### Finalize tests (`finalize-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Finalize tests](#finalize-tests) |
| Depends on | `additional-ci-image-checks`, `additional-prod-image-tests`, `basic-tests`, `build-info`, `build-prod-images`, `ci-image-checks`, `generate-constraints`, `migration-round-trip`, `mypy-providers`, `providers`, `tests-helm`, `tests-integration-system`, `tests-kubernetes`, `tests-mysql-core`, `tests-mysql-providers`, `tests-non-db-core`, `tests-non-db-providers`, `tests-postgres-core`, `tests-postgres-providers`, `tests-sqlite-core`, `tests-sqlite-providers`, `tests-task-sdk`, `tests-airflow-ctl`, `tests-go-sdk`, `tests-with-lowest-direct-resolution-core`, `tests-with-lowest-direct-resolution-providers` |
| Condition | `always() && !failure() && !cancelled()` |

**Permissions:**

- `contents`: `write`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`

### Notify Slack (`notify-slack`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `finalize-tests` |
| Condition | `always() && !cancelled() && github.event_name == 'schedule' && github.run_attempt == 1` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Get failing jobs**
   - ID: `get-failures`

3. **Determine notification action**
   - ID: `notification`

4. **Upload notification state**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `slack-state-tests-${{ github.ref_name }}-amd`
     - `path`: `./slack-state/`
     - `retention-days`: `7`
     - `overwrite`: `true`

5. **Notify Slack (new/changed failures)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_new'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨 Failure Alert: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨 Failure Alert: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

6. **Notify Slack (still not fixed)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_reminder'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨🔁 Still not fixed: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨🔁 Still not fixed: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

7. **Notify Slack (all tests passing)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_recovery'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "✅ All tests passing: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the run log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "✅ All tests passing: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the run log>"`

### Summarize warnings (`summarize-warnings`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `tests-mysql-core`, `tests-mysql-providers`, `tests-non-db-core`, `tests-non-db-providers`, `tests-postgres-core`, `tests-postgres-providers`, `tests-sqlite-core`, `tests-sqlite-providers`, `tests-task-sdk`, `tests-airflow-ctl`, `tests-special`, `tests-with-lowest-direct-resolution-core`, `tests-with-lowest-direct-resolution-providers` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Free up disk space**

4. **Download all test warning artifacts from the current build**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `path`: `./artifacts`
     - `pattern`: `test-warnings-*`

5. **Setup python**
   - Uses: `actions/setup-python@a309ff8b426b58ec0e2a45f0f869d46889d02405` (v6.2.0)
   - With:
     - `python-version`: `${{ inputs.default-python-version }}`

6. **Summarize all warnings**

7. **Upload artifact for summarized warnings**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `test-summarized-warnings`
     - `path`: `./files/warn-summary-*.txt`
     - `retention-days`: `7`
     - `if-no-files-found`: `ignore`
     - `overwrite`: `true`

# Tests (ARM)

| Property | Value |
|----------|-------|
| File | `ci-arm.yml` |
| Triggers | `schedule`, `push`, `workflow_dispatch` |

## Schedule

- `28 1,3,7,9,13,15,19,21 * * *`

## Event filters

- **push**
  - branches: `v[0-9]+-[0-9]+-test`, `providers-[a-z]+-?[a-z]*/v[0-9]+-[0-9]+`

## Permissions

- `contents`: `read` - All other permissions are set to none by default

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |
| `VERBOSE` | `true` |

**Concurrency:** group `ci-arm-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
ci-arm.yml [schedule, push, workflow_dispatch]
+-- build-info / Install Breeze (uses ./.github/actions/breeze)
+-- basic-tests (uses basic-tests.yml)
|   +-- run-breeze-tests / Install Breeze (uses ./.github/actions/breeze)
|   +-- run-breeze-integration-tests / Install Breeze (uses ./.github/actions/breeze)
|   +-- check-translation-completness / Install Breeze (uses ./.github/actions/breeze)
|   +-- static-checks-basic-checks-only / Install Breeze (uses ./.github/actions/breeze)
|   +-- static-checks-basic-checks-only / Install prek (uses ./.github/actions/install-prek)
|   +-- test-airflow-release-commands / Install Breeze (uses ./.github/actions/breeze)
+-- build-ci-images (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- additional-ci-image-checks (uses additional-ci-image-checks.yml)
|   +-- push-early-buildx-cache-to-github-registry (uses push-image-cache.yml)
|   |   +-- push-ci-image-cache / Install Breeze (uses ./.github/actions/breeze)
|   |   +-- push-prod-image-cache / Install Breeze (uses ./.github/actions/breeze)
|   +-- check-that-image-builds-quickly / Install Breeze (uses ./.github/actions/breeze)
+-- generate-constraints (uses generate-constraints.yml)
|   +-- generate-constraints-matrix / Install prek (uses ./.github/actions/install-prek)
|   +-- generate-constraints-matrix / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- ci-image-checks (uses ci-image-checks.yml)
|   +-- static-checks / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- static-checks / Install prek (uses ./.github/actions/install-prek)
|   +-- build-docs / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- publish-docs / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-python-api-client / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- mypy-providers / Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- mypy-providers / Install prek (uses ./.github/actions/install-prek)
+-- migration-round-trip / Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- migration-round-trip / Install prek (uses ./.github/actions/install-prek)
+-- providers (uses test-providers.yml)
|   +-- prepare-install-verify-provider-distributions / Install prek (uses ./.github/actions/install-prek)
|   +-- prepare-install-verify-provider-distributions / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- providers-compatibility-tests-matrix / Install prek (uses ./.github/actions/install-prek)
|   +-- providers-compatibility-tests-matrix / Prepare breeze & CI image: ${{ matrix.compat.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-helm (uses helm-tests.yml)
|   +-- tests-helm / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-helm-release / Install Breeze (uses ./.github/actions/breeze)
+-- tests-postgres-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-postgres-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-mysql-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-mysql-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-sqlite-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-sqlite-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-non-db-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-non-db-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-special (uses special-tests.yml)
|   +-- tests-min-sqlalchemy (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-min-sqlalchemy-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-latest-sqlalchemy (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-latest-sqlalchemy-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-boto-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-boto-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-pendulum-2-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-pendulum-2-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-quarantined-core (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-quarantined-providers (uses run-unit-tests.yml)
|   |   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   |   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   |   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-system-core (uses run-unit-tests.yml)
|       +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|       +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|       +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|       +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-integration-system (uses integration-system-tests.yml)
|   +-- tests-core-integration / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-core-integration / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-core-integration / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-providers-integration / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-providers-integration / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-providers-integration / Post Tests failure (uses ./.github/actions/post_tests_failure)
|   +-- tests-system / Prepare breeze & CI image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests-system / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests-system / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-with-lowest-direct-resolution-core (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- tests-with-lowest-direct-resolution-providers (uses run-unit-tests.yml)
|   +-- tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- tests / Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
 (uses ./.github/actions/migration_tests)
|   +-- tests / Post Tests success (uses ./.github/actions/post_tests_success)
|   +-- tests / Post Tests failure (uses ./.github/actions/post_tests_failure)
+-- build-prod-images (uses prod-image-build.yml)
|   +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
+-- additional-prod-image-tests (uses additional-prod-image-tests.yml)
|   +-- prod-image-extra-checks-main (uses prod-image-extra-checks.yml)
|   |   +-- pip-image (uses prod-image-build.yml)
|   |       +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   |       +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   |       +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
|   +-- prod-image-extra-checks-release-branch (uses prod-image-extra-checks.yml)
|   |   +-- pip-image (uses prod-image-build.yml)
|   |       +-- build-prod-packages / Install prek (uses ./.github/actions/install-prek)
|   |       +-- build-prod-packages / Install Breeze (uses ./.github/actions/breeze)
|   |       +-- build-prod-images / Install Breeze (uses ./.github/actions/breeze)
|   +-- test-examples-of-prod-image-building / Prepare breeze & PROD image: ${{ inputs.default-python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-docker-compose-quick-start / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- task-sdk-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-basic (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log-elasticsearch (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-remote-log-opensearch (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-xcom-object-storage (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-e2e-integration-tests-event-driven (uses airflow-e2e-tests.yml)
|   |   +-- test-e2e-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- test-ui-e2e-chromium (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- test-ui-e2e-firefox (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- test-ui-e2e-webkit (uses ui-e2e-tests.yml)
|   |   +-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
|   |   +-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
|   +-- airflow-ctl-integration-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-kubernetes (uses k8s-tests.yml)
|   +-- tests-kubernetes / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
+-- tests-task-sdk (uses airflow-distributions-tests.yml)
|   +-- distributions-tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- distributions-tests / Install Breeze (uses ./.github/actions/breeze)
+-- tests-airflow-ctl (uses airflow-distributions-tests.yml)
|   +-- distributions-tests / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
|   +-- distributions-tests / Install Breeze (uses ./.github/actions/breeze)
+-- finalize-tests (uses finalize-tests.yml)
    +-- dependency-upgrade-summary / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
    +-- push-buildx-cache-to-github-registry (uses push-image-cache.yml)
        +-- push-ci-image-cache / Install Breeze (uses ./.github/actions/breeze)
        +-- push-prod-image-cache / Install Breeze (uses ./.github/actions/breeze)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `DOCS_AWS_ACCESS_KEY_ID`, `DOCS_AWS_SECRET_ACCESS_KEY`, `SLACK_BOT_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN`; job `mypy-providers` env `GITHUB_TOKEN`; job `migration-round-trip` env `GITHUB_TOKEN`; job `tests-go-sdk` env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | workflow env `SLACK_BOT_TOKEN`; job `ci-image-checks` secrets `SLACK_BOT_TOKEN` |
| `DOCS_AWS_ACCESS_KEY_ID` | job `ci-image-checks` secrets `DOCS_AWS_ACCESS_KEY_ID` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `ci-image-checks` secrets `DOCS_AWS_SECRET_ACCESS_KEY` |

## Jobs

### Build info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Fetch incoming commit ${{ github.sha }} with its parent**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ github.sha }}`
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

4. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

5. **Save github context to file**

6. **Get information about the Workflow**
   - ID: `source-run-info`

7. **Selective checks**
   - ID: `selective-checks`

8. **env**

### Platform: ARM (`print-platform`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info` |

#### Steps

1. **Print architecture**

### Basic tests (`basic-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Basic tests](#basic-tests) |
| Depends on | `build-info` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `run-ui-tests`: `${{needs.build-info.outputs.run-ui-tests}}`
- `run-www-tests`: `${{needs.build-info.outputs.run-www-tests}}`
- `run-api-codegen`: `${{needs.build-info.outputs.run-api-codegen}}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `basic-checks-only`: `${{ needs.build-info.outputs.basic-checks-only }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `canary-run`: `${{needs.build-info.outputs.canary-run}}`
- `run-breeze-integration-tests`: `${{needs.build-info.outputs.run-breeze-integration-tests}}`
- `run-scripts-tests`: `${{needs.build-info.outputs.run-scripts-tests}}`
- `latest-versions-only`: `${{needs.build-info.outputs.latest-versions-only}}`
- `use-uv`: `${{needs.build-info.outputs.use-uv}}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `shared-distributions-as-json`: `${{needs.build-info.outputs.shared-distributions-as-json}}`

### Build CI images (`build-ci-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.ci-image-build == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `${{ needs.build-info.outputs.canary-run }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`

### Additional CI image checks (`additional-ci-image-checks`)

| Property | Value |
|----------|-------|
| Uses workflow | [Additional CI image checks](#additional-ci-image-checks) |
| Depends on | `build-info`, `build-ci-images` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`
- `id-token`: `write` (OIDC)

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `latest-versions-only`: `${{ needs.build-info.outputs.latest-versions-only }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Generate constraints (`generate-constraints`)

| Property | Value |
|----------|-------|
| Uses workflow | [Generate constraints](#generate-constraints) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.ci-image-build == 'true'` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `generate-pypi-constraints`: `true`
- `generate-no-providers-constraints`: `${{ needs.build-info.outputs.canary-run }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### CI image checks (`ci-image-checks`)

| Property | Value |
|----------|-------|
| Uses workflow | [CI Image Checks](#ci-image-checks) |
| Depends on | `build-info`, `build-ci-images` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `docs-list-as-string`: `${{ needs.build-info.outputs.docs-list-as-string }}`
- `latest-versions-only`: `${{ needs.build-info.outputs.latest-versions-only }}`
- `basic-checks-only`: `${{ needs.build-info.outputs.basic-checks-only }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `skip-prek-hooks`: `${{ needs.build-info.outputs.skip-prek-hooks }}`
- `ci-image-build`: `${{ needs.build-info.outputs.ci-image-build }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `docs-build`: `${{ needs.build-info.outputs.docs-build }}`
- `run-api-codegen`: `${{ needs.build-info.outputs.run-api-codegen }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `source-head-repo`: `${{ needs.build-info.outputs.source-head-repo }}`
- `source-head-ref`: `${{ needs.build-info.outputs.source-head-ref }}`

#### Secrets forwarded

- `DOCS_AWS_ACCESS_KEY_ID`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
- `DOCS_AWS_SECRET_ACCESS_KEY`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
- `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_BOT_TOKEN }}`

### MyPy providers checks (`mypy-providers`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-mypy-providers == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ needs.build-info.outputs.default-python-version }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ needs.build-info.outputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ needs.build-info.outputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{steps.breeze.outputs.host-python-version}}` - Python version to use
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

5. **MyPy checks for providers**

### Migration round-trip check (`migration-round-trip`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.has-migrations == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ needs.build-info.outputs.default-python-version }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ needs.build-info.outputs.default-python-version }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ needs.build-info.outputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ needs.build-info.outputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{steps.breeze.outputs.host-python-version}}` - Python version to use
     - `platform`: `${{ needs.build-info.outputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

5. **Migration round-trip check**

### provider distributions tests (`providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Provider tests](#provider-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.skip-providers-tests != 'true' && needs.build-info.outputs.latest-versions-only != 'true' && needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `selected-providers-list-as-string`: `${{ needs.build-info.outputs.selected-providers-list-as-string }}`
- `providers-compatibility-tests-matrix`: `${{ needs.build-info.outputs.providers-compatibility-tests-matrix }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `providers-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Helm tests (`tests-helm`)

| Property | Value |
|----------|-------|
| Uses workflow | [Helm tests](#helm-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-helm-tests == 'true' && needs.build-info.outputs.default-branch == 'main' && needs.build-info.outputs.latest-versions-only != 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `helm-test-packages`: `${{ needs.build-info.outputs.helm-test-packages }}`
- `helm-test-kubernetes-versions`: `${{ needs.build-info.outputs.helm-test-kubernetes-versions }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Postgres tests: core (`tests-postgres-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `postgres`
- `test-name`: `Postgres`
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.postgres-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.postgres-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-migration-tests`: `true`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Postgres tests: providers (`tests-postgres-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `postgres`
- `test-name`: `Postgres`
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.postgres-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.postgres-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-migration-tests`: `true`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### MySQL tests: core (`tests-mysql-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && needs.build-info.outputs.platform == 'linux/amd64'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `mysql`
- `test-name`: `MySQL`
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.mysql-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.mysql-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### MySQL tests: providers (`tests-mysql-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && needs.build-info.outputs.platform == 'linux/amd64'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `mysql`
- `test-name`: `MySQL`
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `${{ needs.build-info.outputs.mysql-versions }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.mysql-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Sqlite tests: core (`tests-sqlite-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: `Sqlite`
- `test-name-separator`: ``
- `test-scope`: `DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Sqlite tests: providers (`tests-sqlite-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: `Sqlite`
- `test-name-separator`: ``
- `test-scope`: `DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `run-migration-tests`: `true`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Non-DB tests: core (`tests-non-db-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: ``
- `test-name-separator`: ``
- `test-scope`: `Non-DB`
- `test-group`: `core`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Non-DB tests: providers (`tests-non-db-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `backend`: `sqlite`
- `test-name`: ``
- `test-name-separator`: ``
- `test-scope`: `Non-DB`
- `test-group`: `providers`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `${{ needs.build-info.outputs.sqlite-exclude }}`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Special tests (`tests-special`)

| Property | Value |
|----------|-------|
| Uses workflow | [Special tests](#special-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true' && (needs.build-info.outputs.canary-run == 'true' \|\|<br> needs.build-info.outputs.upgrade-to-newer-dependencies != 'false' \|\|<br> needs.build-info.outputs.full-tests-needed == 'true')` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`
- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `core-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `providers-test-types-list-as-strings-in-json`: `${{ needs.build-info.outputs.providers-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Integration and System Tests (`tests-integration-system`)

| Property | Value |
|----------|-------|
| Uses workflow | [Integration and system tests](#integration-and-system-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `testable-core-integrations`: `${{ needs.build-info.outputs.testable-core-integrations }}`
- `testable-providers-integrations`: `${{ needs.build-info.outputs.testable-providers-integrations }}`
- `run-system-tests`: `${{ needs.build-info.outputs.run-system-tests }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `default-postgres-version`: `${{ needs.build-info.outputs.default-postgres-version }}`
- `default-mysql-version`: `${{ needs.build-info.outputs.default-mysql-version }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`

### Low dep tests:core (`tests-with-lowest-direct-resolution-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `test-name`: `LowestDeps`
- `force-lowest-dependencies`: `true`
- `test-scope`: `All`
- `test-group`: `core`
- `backend`: `sqlite`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['${{ needs.build-info.outputs.default-postgres-version }}']`
- `excluded-providers-as-string`: ``
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `monitor-delay-time-in-seconds`: `120`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Low dep tests: providers (`tests-with-lowest-direct-resolution-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `test-name`: `LowestDeps`
- `force-lowest-dependencies`: `true`
- `test-scope`: `All`
- `test-group`: `providers`
- `backend`: `sqlite`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `backend-versions`: `['${{ needs.build-info.outputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ needs.build-info.outputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ needs.build-info.outputs.individual-providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `run-coverage`: `${{ needs.build-info.outputs.run-coverage }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `monitor-delay-time-in-seconds`: `120`
- `skip-providers-tests`: `${{ needs.build-info.outputs.skip-providers-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`

### Build PROD images (`build-prod-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build PROD images](#build-prod-images) |
| Depends on | `build-info`, `build-ci-images`, `generate-constraints` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `build-type`: `Regular`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-package-artifact`: `true`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `prod-image-build`: `${{ needs.build-info.outputs.prod-image-build }}`

### Additional PROD image tests (`additional-prod-image-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Additional PROD image tests](#additional-prod-image-tests) |
| Depends on | `build-info`, `build-prod-images`, `generate-constraints` |
| Condition | `needs.build-info.outputs.prod-image-build == 'true'` |

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `run-task-sdk-integration-tests`: `${{ needs.build-info.outputs.run-task-sdk-integration-tests }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `run-remote-logging-elasticsearch-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-elasticsearch-e2e-tests }}`
- `run-remote-logging-opensearch-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-opensearch-e2e-tests }}`
- `run-remote-logging-s3-e2e-tests`: `${{ needs.build-info.outputs.run-remote-logging-s3-e2e-tests }}`
- `run-event-driven-e2e-tests`: `${{ needs.build-info.outputs.run-event-driven-e2e-tests }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `run-ui-e2e-tests`: `${{ needs.build-info.outputs.run-ui-e2e-tests }}`
- `run-airflow-ctl-integration-tests`: `${{ needs.build-info.outputs.run-airflow-ctl-integration-tests }}`

### Kubernetes tests (`tests-kubernetes`)

| Property | Value |
|----------|-------|
| Uses workflow | [K8s tests](#k8s-tests) |
| Depends on | `build-info`, `build-prod-images` |
| Condition | `( needs.build-info.outputs.run-kubernetes-tests == 'true' \|\| needs.build-info.outputs.run-helm-tests == 'true')` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`
- `kubernetes-combos`: `${{ needs.build-info.outputs.kubernetes-combos }}`

### Task SDK tests (`tests-task-sdk`)

| Property | Value |
|----------|-------|
| Uses workflow | [Non-core Distribution tests](#non-core-distribution-tests) |
| Depends on | `build-info`, `build-ci-images` |
| Condition | `needs.build-info.outputs.run-task-sdk-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `distribution-name`: `task-sdk`
- `distribution-cmd-format`: `prepare-task-sdk-distributions`
- `test-type`: `task-sdk-tests`
- `use-local-venv`: `false`
- `test-timeout`: `20`

### Go SDK tests (`tests-go-sdk`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.run-go-sdk-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Setup Go**
   - Uses: `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0)
   - With:
     - `go-version`: `1.24`
     - `cache-dependency-path`: `go-sdk/go.sum`

3. **Setup Gotestsum**

4. **Cleanup dist files**

5. **Run Go tests**

### Airflow CTL tests (`tests-airflow-ctl`)

| Property | Value |
|----------|-------|
| Uses workflow | [Non-core Distribution tests](#non-core-distribution-tests) |
| Depends on | `build-info` |
| Condition | `needs.build-info.outputs.run-airflow-ctl-tests == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `distribution-name`: `airflow-ctl`
- `distribution-cmd-format`: `prepare-airflow-ctl-distributions`
- `test-type`: `airflow-ctl-tests`
- `use-local-venv`: `true`
- `test-timeout`: `20`

### Finalize tests (`finalize-tests`)

| Property | Value |
|----------|-------|
| Uses workflow | [Finalize tests](#finalize-tests) |
| Depends on | `additional-ci-image-checks`, `additional-prod-image-tests`, `basic-tests`, `build-info`, `build-prod-images`, `ci-image-checks`, `generate-constraints`, `migration-round-trip`, `mypy-providers`, `providers`, `tests-helm`, `tests-integration-system`, `tests-kubernetes`, `tests-mysql-core`, `tests-mysql-providers`, `tests-non-db-core`, `tests-non-db-providers`, `tests-postgres-core`, `tests-postgres-providers`, `tests-sqlite-core`, `tests-sqlite-providers`, `tests-task-sdk`, `tests-airflow-ctl`, `tests-go-sdk`, `tests-with-lowest-direct-resolution-core`, `tests-with-lowest-direct-resolution-providers` |
| Condition | `always() && !failure() && !cancelled()` |

**Permissions:**

- `contents`: `write`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `${{ needs.build-info.outputs.runner-type }}`
- `platform`: `${{ needs.build-info.outputs.platform }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `default-python-version`: `${{ needs.build-info.outputs.default-python-version }}`
- `upgrade-to-newer-dependencies`: `${{ needs.build-info.outputs.upgrade-to-newer-dependencies }}`
- `include-success-outputs`: `${{ needs.build-info.outputs.include-success-outputs }}`
- `docker-cache`: `${{ needs.build-info.outputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ needs.build-info.outputs.disable-airflow-repo-cache }}`
- `canary-run`: `${{ needs.build-info.outputs.canary-run }}`
- `use-uv`: `${{ needs.build-info.outputs.use-uv }}`
- `debug-resources`: `${{ needs.build-info.outputs.debug-resources }}`

### Notify Slack (`notify-slack`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `finalize-tests` |
| Condition | `always() && !cancelled() && github.event_name == 'schedule' && github.run_attempt == 1` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Get failing jobs**
   - ID: `get-failures`

3. **Determine notification action**
   - ID: `notification`

4. **Upload notification state**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `slack-state-tests-${{ github.ref_name }}-arm`
     - `path`: `./slack-state/`
     - `retention-days`: `7`
     - `overwrite`: `true`

5. **Notify Slack (new/changed failures)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_new'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨 Failure Alert: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨 Failure Alert: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

6. **Notify Slack (still not fixed)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_reminder'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨🔁 Still not fixed: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨🔁 Still not fixed: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\nFailing jobs:\n${{ steps.get-failures.outputs.failed-jobs }}\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

7. **Notify Slack (all tests passing)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_recovery'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "✅ All tests passing: Scheduled CI (${{ needs.build-info.outputs.platform }}) on branch *${{ github.ref_name }}*\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the run log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "✅ All tests passing: Scheduled CI (${{ needs.build-info.outputs.platform }}) on *${{ github.ref_name }}*\n\n*Details:* <https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the run log>"`

### Summarize warnings (`summarize-warnings`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(needs.build-info.outputs.runner-type) }}` |
| Depends on | `build-info`, `tests-mysql-core`, `tests-mysql-providers`, `tests-non-db-core`, `tests-non-db-providers`, `tests-postgres-core`, `tests-postgres-providers`, `tests-sqlite-core`, `tests-sqlite-providers`, `tests-task-sdk`, `tests-airflow-ctl`, `tests-special`, `tests-with-lowest-direct-resolution-core`, `tests-with-lowest-direct-resolution-providers` |
| Condition | `needs.build-info.outputs.run-unit-tests == 'true'` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Free up disk space**

4. **Download all test warning artifacts from the current build**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `path`: `./artifacts`
     - `pattern`: `test-warnings-*`

5. **Setup python**
   - Uses: `actions/setup-python@a309ff8b426b58ec0e2a45f0f869d46889d02405` (v6.2.0)
   - With:
     - `python-version`: `${{ inputs.default-python-version }}`

6. **Summarize all warnings**

7. **Upload artifact for summarized warnings**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `test-summarized-warnings`
     - `path`: `./files/warn-summary-*.txt`
     - `retention-days`: `7`
     - `if-no-files-found`: `ignore`
     - `overwrite`: `true`

# Build CI images

| Property | Value |
|----------|-------|
| File | `ci-image-build.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `target-commit-sha` | string | No | - | The commit SHA to checkout for the build |
| `pull-request-target` | string | No | `false` | Whether we are running this from pull-request-target workflow (true/false) |
| `is-committer-build` | string | No | `false` | Whether the build is executed by committer (true/false) |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `push-image` | string | No | `true` | Whether to push image to the registry (true/false) |
| `upload-image-artifact` | string | Yes | - | Whether to upload docker image artifact |
| `upload-mount-cache-artifact` | string | Yes | - | Whether to upload mount-cache artifact |
| `debian-version` | string | No | `bookworm` | Base Debian distribution to use for the build (bookworm) |
| `install-mysql-client-type` | string | No | `mariadb` | MySQL client type to use during build (mariadb/mysql) |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `python-versions` | string | Yes | `[""]` | JSON-formatted array of Python versions to build images from |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `constraints-branch` | string | Yes | - | Branch used to construct constraints URL from. |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to attempt to upgrade image to newer dependencies (false/RANDOM_VALUE) |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |

## Permissions

- `contents`: `read`

## Called by

```
ci-image-build.yml
+-- ci-amd.yml (job: build-ci-images)  <- entry point
+-- ci-arm.yml (job: build-ci-images)  <- entry point
+-- registry-backfill.yml (job: build-ci-image)  <- entry point
+-- registry-build.yml (job: build-ci-image)  <- entry point
|   +-- publish-docs-to-s3.yml (job: update-registry)  <- entry point
+-- update-constraints-on-push-stable.yml (job: build-ci-images)  <- entry point
+-- update-constraints-on-push.yml (job: build-ci-images)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `build-ci-images` env `GITHUB_TOKEN` |

## Jobs

### Build CI ${{ inputs.platform }} image ${{ matrix.python-version }} (`build-ci-images`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `sqlite` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.python-version }}` |
| `DEFAULT_BRANCH` | `${{ inputs.branch }}` |
| `DEFAULT_CONSTRAINTS_BRANCH` | `${{ inputs.constraints-branch }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout target branch**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Free up disk space**

4. **Make /mnt writeable**

5. **Move docker to /mnt**

6. **Install Breeze**
   - Uses: `./.github/actions/breeze`

7. **Restore ci-cache mount image ${{ inputs.platform }}:${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `restore-cache-mount`
   - Uses: `apache/infrastructure-actions/stash/restore@49df447b39b18354895520e0a63731b7cad7cbec`
   - With:
     - `key`: `ci-cache-mount-save-v3-${{ inputs.platform }}-${{ env.PYTHON_MAJOR_MINOR_VERSION }}`
     - `path`: `/tmp/`

8. **Verify ci-cache file exists**
   - Condition: `steps.restore-cache-mount.outputs.stash-hit == 'true'`

9. **Import mount-cache ${{ inputs.platform }}:${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Condition: `steps.restore-cache-mount.outputs.stash-hit == 'true'`

10. **Login to ghcr.io**

11. **Build ${{ inputs.push-image == 'true' && ' & push ' || '' }} ${{ inputs.platform }}:${{ env.PYTHON_MAJOR_MINOR_VERSION }} image
**

12. **Export CI docker image ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Condition: `inputs.upload-image-artifact == 'true'`

13. **Stash CI docker image ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `inputs.upload-image-artifact == 'true'`
   - With:
     - `key`: `ci-image-save-v3-${{ inputs.platform }}-${{ env.PYTHON_MAJOR_MINOR_VERSION }}`
     - `path`: `/mnt/ci-image-save-*-${{ env.PYTHON_MAJOR_MINOR_VERSION }}.tar`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

14. **Export mount cache ${{ inputs.platform }}:${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Condition: `inputs.upload-mount-cache-artifact == 'true'`

15. **Stash cache mount ${{ inputs.platform }}:${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `inputs.upload-mount-cache-artifact == 'true'`
   - With:
     - `key`: `ci-cache-mount-save-v3-${{ inputs.platform }}-${{ env.PYTHON_MAJOR_MINOR_VERSION }}`
     - `path`: `/tmp/ci-cache-mount-save-v3-${{ env.PYTHON_MAJOR_MINOR_VERSION }}.tar.gz`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

16. **Check disk space after build**

# CI Image Checks

| Property | Value |
|----------|-------|
| File | `ci-image-checks.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions-list-as-string` | string | Yes | - | The list of python versions as string separated by spaces |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `canary-run` | string | Yes | - | Whether this is a canary run (true/false) |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `docs-list-as-string` | string | Yes | - | Stringified list of docs to build (space separated) |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies (true/false) |
| `basic-checks-only` | string | Yes | - | Whether to run only basic checks (true/false) |
| `latest-versions-only` | string | Yes | - | Whether to run only latest versions (true/false) |
| `ci-image-build` | string | Yes | - | Whether to build CI images (true/false) |
| `skip-prek-hooks` | string | Yes | - | Whether to skip prek hooks (true/false) |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs |
| `debug-resources` | string | Yes | - | Whether to debug resources (true/false) |
| `docs-build` | string | Yes | - | Whether to build docs (true/false) |
| `run-api-codegen` | string | Yes | - | Whether to run API codegen (true/false) |
| `default-postgres-version` | string | Yes | - | The default version of the postgres to use |
| `run-coverage` | string | Yes | - | Whether to run coverage or not (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `source-head-repo` | string | No | `apache/airflow` | The source head repository to use for back-references |
| `source-head-ref` | string | No | `main` | The source head ref to use for back-references |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `DOCS_AWS_ACCESS_KEY_ID` | Yes | - |
| `DOCS_AWS_SECRET_ACCESS_KEY` | Yes | - |
| `SLACK_BOT_TOKEN` | No | - |

## Permissions

- `contents`: `read`

## Called by

```
ci-image-checks.yml
+-- ci-amd.yml (job: ci-image-checks)  <- entry point
+-- ci-arm.yml (job: ci-image-checks)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `static-checks` env `GITHUB_TOKEN`; job `build-docs` env `GITHUB_TOKEN`; job `publish-docs` env `GITHUB_TOKEN`; job `test-python-api-client` env `GITHUB_TOKEN` |
| `DOCS_AWS_ACCESS_KEY_ID` | job `publish-docs` step `Configure AWS credentials` with `aws-access-key-id` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `publish-docs` step `Configure AWS credentials` with `aws-secret-access-key` |

## Jobs

### Static checks (`static-checks`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.basic-checks-only == 'false' && inputs.latest-versions-only != 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `UPGRADE_TO_NEWER_DEPENDENCIES` | `${{ inputs.upgrade-to-newer-dependencies }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{steps.breeze.outputs.host-python-version}}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `true` - Whether to save prek cache (required)

5. **Static checks**

6. **Show prek log on failure**
   - Condition: `failure()`

### Build documentation (`build-docs`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.docs-build == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `INCLUDE_SUCCESS_OUTPUTS` | `${{ inputs.include-success-outputs }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Restore docs inventory cache**
   - ID: `restore-docs-inventory-cache`
   - Uses: `apache/infrastructure-actions/stash/restore@49df447b39b18354895520e0a63731b7cad7cbec`
   - With:
     - `path`: `./generated/_inventory_cache/`
     - `key`: `cache-docs-inventory-v1`

5. **Building docs with ${{ matrix.flag }} flag**

6. **Check for missing third-party inventories**
   - ID: `check-missing-inventories`
   - Condition: `always()`

7. **Get docs build job URL**
   - ID: `get-job-url`
   - Condition: `always() && inputs.canary-run == 'true' && matrix.flag == '--docs-only'`

8. **Determine inventory notification action**
   - ID: `inventory-notification`
   - Condition: `always() && inputs.canary-run == 'true' && matrix.flag == '--docs-only'`

9. **Upload inventory notification state**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `always() && inputs.canary-run == 'true' && matrix.flag == '--docs-only'`
   - With:
     - `name`: `slack-state-inventory-${{ inputs.branch }}-${{ contains(inputs.platform, 'arm') && 'arm' || 'amd' }}`
     - `path`: `./slack-state/`
     - `retention-days`: `7`
     - `overwrite`: `true`

10. **Notify Slack about missing inventories (new/changed)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `inputs.canary-run == 'true' && matrix.flag == '--docs-only' && steps.inventory-notification.outputs.action == 'notify_new'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "⚠️ Missing 3rd-party doc inventories in canary build on *${{ github.ref_name }}*: ${{ steps.check-missing-inventories.outputs.packages }}\n\n<${{ steps.get-job-url.outputs.url }}|View job log>"`

11. **Notify Slack about missing inventories (still not fixed)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `inputs.canary-run == 'true' && matrix.flag == '--docs-only' && steps.inventory-notification.outputs.action == 'notify_reminder'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "⚠️🔁 Still not fixed: Missing 3rd-party doc inventories in canary build on *${{ github.ref_name }}*: ${{ steps.check-missing-inventories.outputs.packages }}\n\n<${{ steps.get-job-url.outputs.url }}|View job log>"`

12. **Notify Slack about inventory recovery**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `inputs.canary-run == 'true' && matrix.flag == '--docs-only' && steps.inventory-notification.outputs.action == 'notify_recovery'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "✅ All 3rd-party doc inventories are now available in canary build on *${{ github.ref_name }}*\n\n<${{ steps.get-job-url.outputs.url }}|View job log>"`

13. **Save docs inventory cache**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `steps.restore-docs-inventory-cache.outputs.stash-hit != 'true' && matrix.flag == '--docs-only'`
   - With:
     - `path`: `./generated/_inventory_cache/`
     - `key`: `cache-docs-inventory-v1`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

14. **Upload build docs**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `matrix.flag == '--docs-only'`
   - With:
     - `name`: `airflow-docs`
     - `path`: `./generated/_build`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

### Publish documentation and validate versions (`publish-docs`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Depends on | `build-docs` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `INCLUDE_SUCCESS_OUTPUTS` | `${{ inputs.include-success-outputs }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `VERBOSE` | `true` |
| `HEAD_REPO` | `${{ inputs.source-head-repo }}` |
| `HEAD_REF` | `${{ inputs.source-head-ref }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Download docs prepared as artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `airflow-docs`
     - `path`: `./generated/_build`

5. **Make sure SBOM dir exists and has the right permissions**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

6. **Determine Airflow version for SBOMs**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

7. **Prepare SBOMs**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

8. **Generated SBOM files**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

9. **Check disk space available**

10. **Create /mnt/airflow-site directory**

11. **Clone airflow-site**

12. **Publish docs**

13. **Check disk space available**

14. **Generate back references for providers**

15. **Generate back references for apache-airflow**

16. **Generate back references for docker-stack**

17. **Generate back references for helm-chart**

18. **Generate back references for apache-airflow-ctl**

19. **Validate published doc versions**
   - ID: `validate-docs-versions`

20. **Install AWS CLI v2**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

21. **Configure AWS credentials**
   - Uses: `aws-actions/configure-aws-credentials@d979d5b3a71173a29b74b5b88418bfda9437d885` (v6.1.1)
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`
   - With:
     - `aws-access-key-id`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
     - `aws-secret-access-key`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
     - `aws-region`: `eu-central-1`

22. **Upload documentation to AWS S3**
   - Condition: `inputs.canary-run == 'true' && (github.event_name == 'schedule' || github.event_name == 'workflow_dispatch')`

### Test Python API client (`test-python-api-client`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-api-codegen == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `postgres` |
| `BACKEND_VERSION` | `${{ inputs.default-postgres-version }}` |
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `ENABLE_COVERAGE` | `${{ inputs.run-coverage }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `JOB_ID` | `python-api-client-tests` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `repository`: `apache/airflow-client-python`
     - `fetch-depth`: `1`
     - `persist-credentials`: `false`
     - `path`: `./airflow-client-python`

4. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

5. **Generate airflow python client**

6. **Show diff**

7. **Python API client tests**

# CI Notification

| Property | Value |
|----------|-------|
| File | `ci-notification.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 6,17 * * *`

## Permissions

- `contents`: `read` - All other permissions are set to none by default

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |
| `VERBOSE` | `true` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | workflow env `SLACK_BOT_TOKEN` |

## Jobs

### `workflow-status`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Find workflow run status**
   - ID: `find-workflow-run-status`

3. **Determine notification action**
   - ID: `notification`

4. **Upload notification state**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `slack-state-ci-${{ matrix.branch }}-${{ matrix.workflow-id }}`
     - `path`: `./slack-state/`
     - `retention-days`: `7`
     - `overwrite`: `true`

5. **Send Slack notification (new/changed failures)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_new'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨 Failure Alert: ${{ env.workflow_id }} on branch *${{ env.branch }}*\n\nFailing jobs:\n${{ steps.find-workflow-run-status.outputs.failed-jobs }}\n\n*Details:* <${{ env.run_url }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨 Failure Alert: ${{ env.workflow_id }} on *${{ env.branch }}*\n\nFailing jobs:\n${{ steps.find-workflow-run-status.outputs.failed-jobs }}\n\n*Details:* <${{ env.run_url }}|View the failure log>"`

6. **Send Slack notification (still not fixed)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_reminder'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨🔁 Still not fixed: ${{ env.workflow_id }} on branch *${{ env.branch }}*\n\nFailing jobs:\n${{ steps.find-workflow-run-status.outputs.failed-jobs }}\n\n*Details:* <${{ env.run_url }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨🔁 Still not fixed: ${{ env.workflow_id }} on *${{ env.branch }}*\n\nFailing jobs:\n${{ steps.find-workflow-run-status.outputs.failed-jobs }}\n\n*Details:* <${{ env.run_url }}|View the failure log>"`

7. **Send Slack notification (all passing)**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `steps.notification.outputs.action == 'notify_recovery'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "✅ All passing: ${{ env.workflow_id }} on branch *${{ env.branch }}*\n\n*Details:* <${{ env.run_url }}|View the run log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "✅ All passing: ${{ env.workflow_id }} on *${{ env.branch }}*\n\n*Details:* <${{ env.run_url }}|View the run log>"`

# CodeQL

| Property | Value |
|----------|-------|
| File | `codeql-analysis.yml` |
| Triggers | `pull_request`, `push`, `schedule` |

## Schedule

- `0 2 * * *`

## Event filters

- **pull_request**
  - branches: `main`, `v[0-9]+-[0-9]+-test`, `v[0-9]+-[0-9]+-stable`
- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

**Concurrency:** group `codeql-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Jobs

### Analyze (`analyze`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

**Permissions:**

- `actions`: `read`
- `contents`: `read`
- `pull-requests`: `read`
- `security-events`: `write`

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@9e0d7b8d25671d64c341c19c0152d693099fb5ba` (v4.35.5)
   - With:
     - `languages`: `${{ matrix.language }}`

3. **Autobuild**
   - Uses: `github/codeql-action/autobuild@9e0d7b8d25671d64c341c19c0152d693099fb5ba` (v4.35.5)

4. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@9e0d7b8d25671d64c341c19c0152d693099fb5ba` (v4.35.5)
   - With:
     - `category`: `/language:${{matrix.language}}`

# E2E Flaky Tests Report

| Property | Value |
|----------|-------|
| File | `e2e-flaky-tests-report.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 0 * * *`

## Permissions

- `contents`: `read`
- `actions`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | workflow env `SLACK_BOT_TOKEN` |

## Jobs

### Analyze E2E flaky tests (`analyze-flaky-tests`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Analyze E2E test results**
   - ID: `analyze`

3. **Post report to Slack**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - Condition: `always() && steps.analyze.outcome == 'success'`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload-file-path`: `slack-message.json`

4. **Upload analysis results**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `always()`
   - With:
     - `name`: `e2e-flaky-test-analysis`
     - `path`: `slack-message.json`
     - `retention-days`: `14`

# Finalize tests

| Property | Value |
|----------|-------|
| File | `finalize-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to test |
| `python-versions-list-as-string` | string | Yes | - | Stringified array of all Python versions to test - separated by spaces. |
| `branch` | string | Yes | - | The default branch to use for the build |
| `constraints-branch` | string | Yes | - | The branch to use for constraints |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies (true/false) |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs (true/false) |
| `canary-run` | string | Yes | - | Whether this is a canary run (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `debug-resources` | string | Yes | - | Whether to debug resources or not (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
finalize-tests.yml
+-- ci-amd.yml (job: finalize-tests)  <- entry point
+-- ci-arm.yml (job: finalize-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `update-constraints` env `GITHUB_TOKEN`; job `dependency-upgrade-summary` env `GITHUB_TOKEN` |

## Jobs

### Update constraints (`update-constraints`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.upgrade-to-newer-dependencies != 'false' && inputs.platform == 'linux/amd64'` |

**Permissions:**

- `contents`: `write`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources}}` |
| `PYTHON_VERSIONS` | `${{ inputs.python-versions-list-as-string }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Set constraints branch name**
   - ID: `constraints-branch`

4. **Checkout ${{ steps.constraints-branch.outputs.branch }}**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `path`: `constraints`
     - `ref`: `${{ steps.constraints-branch.outputs.branch }}`
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

5. **Download constraints from the constraints generated by build CI image**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `constraints-*`
     - `path`: `./files`

6. **Diff in constraints for Python: ${{ inputs.python-versions-list-as-string }}**

### Deps ${{ matrix.python-version }}:${{ matrix.constraints-mode }} (`dependency-upgrade-summary`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Depends on | `update-constraints` |
| Condition | `inputs.upgrade-to-newer-dependencies == 'true' && inputs.platform == 'linux/amd64'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ matrix.python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ matrix.python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Deps: ${{ matrix.python-version }}:${{ matrix.constraints-mode }}**

### Push Regular Image Cache ${{ inputs.platform }} (`push-buildx-cache-to-github-registry`)

| Property | Value |
|----------|-------|
| Uses workflow | [Push image cache](#push-image-cache) |
| Depends on | `update-constraints` |
| Condition | `inputs.canary-run == 'true' && github.event_name != 'pull_request'` |

**Permissions:**

- `contents`: `read`
- `packages`: `write` - This write is only given here for `push` events from "apache/airflow" repo. It is not given for PRs from forks. This is to prevent malicious PRs from creating images in the "apache/airflow" repo.

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `cache-type`: `Regular AMD`
- `include-prod-images`: `true`
- `push-latest-images`: `true`
- `python-versions`: `${{ inputs.python-versions }}`
- `branch`: `${{ inputs.branch }}`
- `constraints-branch`: `${{ inputs.constraints-branch }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `docker-cache`: `${{ inputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ inputs.disable-airflow-repo-cache }}`

# Generate constraints

| Property | Value |
|----------|-------|
| File | `generate-constraints.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions-list-as-string` | string | Yes | - | Stringified array of all Python versions to test - separated by spaces. |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to generate constraints for |
| `generate-no-providers-constraints` | string | Yes | - | Whether to generate constraints without providers (true/false) |
| `generate-pypi-constraints` | string | Yes | - | Whether to generate PyPI constraints (true/false) |
| `debug-resources` | string | Yes | - | Whether to run in debug mode (true/false) |
| `use-uv` | string | Yes | - | Whether to use uvloop (true/false) |

## Called by

```
generate-constraints.yml
+-- ci-amd.yml (job: generate-constraints)  <- entry point
+-- ci-arm.yml (job: generate-constraints)  <- entry point
+-- update-constraints-on-push-stable.yml (job: generate-constraints)  <- entry point
+-- update-constraints-on-push.yml (job: generate-constraints)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `generate-constraints-matrix` env `GITHUB_TOKEN` |

## Jobs

### Generate constraints for ${{ matrix.python-version }} on ${{ inputs.platform }} (`generate-constraints-matrix`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Permissions:**

- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `true` |
| `PYTHON_VERSION` | `${{ matrix.python-version }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ matrix.python-version }}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

4. **Prepare breeze & CI image: ${{ matrix.python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ matrix.python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

5. **Source constraints**

6. **No providers constraints**
   - Condition: `inputs.generate-no-providers-constraints == 'true'`

7. **Prepare updated provider distributions**
   - Condition: `inputs.generate-pypi-constraints == 'true'`

8. **Prepare airflow distributions**
   - Condition: `inputs.generate-pypi-constraints == 'true'`

9. **Prepare task-sdk distribution**
   - Condition: `inputs.generate-pypi-constraints == 'true'`

10. **PyPI constraints**
   - Condition: `inputs.generate-pypi-constraints == 'true'`

11. **Upload constraint artifacts**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `constraints-${{ matrix.python-version }}`
     - `path`: `./files/constraints-${{ matrix.python-version }}/constraints-*.txt`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

12. **Dependency upgrade summary**

# Helm tests

| Property | Value |
|----------|-------|
| File | `helm-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `helm-test-packages` | string | Yes | - | Stringified JSON array of helm test packages to test |
| `helm-test-kubernetes-versions` | string | Yes | - | Stringified JSON array of Kubernetes versions to validate against |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `use-uv` | string | Yes | - | Whether to use uvloop (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
helm-tests.yml
+-- ci-amd.yml (job: tests-helm)  <- entry point
+-- ci-arm.yml (job: tests-helm)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `tests-helm` env `GITHUB_TOKEN`; job `tests-helm-release` env `GITHUB_TOKEN` |

## Jobs

### Unit tests Helm: ${{ matrix.helm-test-package }} (K8S ${{ matrix.kubernetes-version }}) (`tests-helm`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `PARALLEL_TEST_TYPES` | `Helm` |
| `BACKEND` | `none` |
| `DB_RESET` | `false` |
| `JOB_ID` | `helm-tests` |
| `USE_XDIST` | `true` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Helm Unit Tests: ${{ matrix.helm-test-package }} (K8S ${{ matrix.kubernetes-version }})**

### Release Helm (`tests-helm-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{inputs.default-python-version}}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Setup git for tagging**

5. **Remove old artifacts**

6. **Setup k8s/helm environment**

7. **Install helm gpg plugin**

8. **Helm release tarball**

9. **Generate GPG key for signing**

10. **Helm release package**

11. **Verify packaged chart contents and lint**

12. **Sign artifacts for ASF distribution**

13. **Fetch Git Tags**

14. **Test helm chart issue generation**

15. **Upload Helm artifacts**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `Helm artifacts`
     - `path`: `./dist/airflow-*`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

# Integration and system tests

| Property | Value |
|----------|-------|
| File | `integration-system-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining public runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `testable-core-integrations` | string | Yes | - | The list of testable core integrations as JSON array. |
| `testable-providers-integrations` | string | Yes | - | The list of testable providers integrations as JSON array. |
| `run-system-tests` | string | Yes | - | Run system tests (true/false) |
| `default-postgres-version` | string | Yes | - | Default version of Postgres to use |
| `default-mysql-version` | string | Yes | - | Default version of MySQL to use |
| `skip-providers-tests` | string | Yes | - | Skip provider tests (true/false) |
| `run-coverage` | string | Yes | - | Run coverage (true/false) |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `debug-resources` | string | Yes | - | Debug resources (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv |

## Permissions

- `contents`: `read`

## Called by

```
integration-system-tests.yml
+-- ci-amd.yml (job: tests-integration-system)  <- entry point
+-- ci-arm.yml (job: tests-integration-system)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `tests-core-integration` env `GITHUB_TOKEN`; job `tests-providers-integration` env `GITHUB_TOKEN`; job `tests-system` env `GITHUB_TOKEN` |
| `CODECOV_TOKEN` | job `tests-core-integration` step `Post Tests success` with `codecov-token`; job `tests-providers-integration` step `Post Tests success` with `codecov-token`; job `tests-system` step `Post Tests success` with `codecov-token` |

## Jobs

### Integration core ${{ matrix.integration }} (`tests-core-integration`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.testable-core-integrations != '[]'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `postgres` |
| `BACKEND_VERSION` | `${{ inputs.default-postgres-version }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `JOB_ID` | `integration-core-${{ matrix.integration }}` |
| `SKIP_PROVIDERS_TESTS` | `${{ inputs.skip-providers-tests }}` |
| `ENABLE_COVERAGE` | `${{ inputs.run-coverage}}` |
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Integration: core ${{ matrix.integration }}**

5. **Post Tests success**
   - Uses: `./.github/actions/post_tests_success`
   - With:
     - `codecov-token`: `${{ secrets.CODECOV_TOKEN }}` - Codecov token (required)
     - `python-version`: `${{ inputs.default-python-version }}` - Python version (required)

6. **Post Tests failure**
   - Uses: `./.github/actions/post_tests_failure`
   - Condition: `failure()`

### Integration: providers ${{ matrix.integration }} (`tests-providers-integration`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.testable-providers-integrations != '[]' && inputs.skip-providers-tests != 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `postgres` |
| `BACKEND_VERSION` | `${{ inputs.default-postgres-version }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `JOB_ID` | `integration-providers-${{ matrix.integration }}` |
| `SKIP_PROVIDERS_TESTS` | `${{ inputs.skip-providers-tests }}` |
| `ENABLE_COVERAGE` | `${{ inputs.run-coverage}}` |
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Integration: providers ${{ matrix.integration }}**

5. **Post Tests success**
   - Uses: `./.github/actions/post_tests_success`
   - With:
     - `codecov-token`: `${{ secrets.CODECOV_TOKEN }}` - Codecov token (required)
     - `python-version`: `${{ inputs.default-python-version }}` - Python version (required)

6. **Post Tests failure**
   - Uses: `./.github/actions/post_tests_failure`
   - Condition: `failure()`

### System Tests (`tests-system`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.run-system-tests == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `postgres` |
| `BACKEND_VERSION` | `${{ inputs.default-postgres-version }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `JOB_ID` | `system` |
| `SKIP_PROVIDERS_TESTS` | `${{ inputs.skip-providers-tests }}` |
| `ENABLE_COVERAGE` | `${{ inputs.run-coverage}}` |
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **System Tests**

5. **Post Tests success**
   - Uses: `./.github/actions/post_tests_success`
   - With:
     - `codecov-token`: `${{ secrets.CODECOV_TOKEN }}` - Codecov token (required)
     - `python-version`: `${{ inputs.default-python-version }}` - Python version (required)

6. **Post Tests failure**
   - Uses: `./.github/actions/post_tests_failure`
   - Condition: `failure()`

# K8s tests

| Property | Value |
|----------|-------|
| File | `k8s-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions-list-as-string` | string | Yes | - | List of Python versions to test: space separated string |
| `kubernetes-combos` | string | Yes | - | Array of combinations of Kubernetes and Python versions to test |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs |
| `use-uv` | string | Yes | - | Whether to use uv |
| `debug-resources` | string | Yes | - | Whether to debug resources |

## Permissions

- `contents`: `read`

## Called by

```
k8s-tests.yml
+-- ci-amd.yml (job: tests-kubernetes)  <- entry point
+-- ci-arm.yml (job: tests-kubernetes)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `tests-kubernetes` env `GITHUB_TOKEN` |

## Jobs

### K8S System:${{ matrix.executor }}-${{ matrix.kubernetes-combo }}-${{ matrix.use-standard-naming }} (`tests-kubernetes`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `${{ inputs.include-success-outputs }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Prepare PYTHON_MAJOR_MINOR_VERSION and KUBERNETES_VERSION**
   - ID: `prepare-versions`

3. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

4. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

5. **Run complete K8S tests ${{ matrix.executor }}-${{ env.PYTHON_MAJOR_MINOR_VERSION }}-${{env.KUBERNETES_VERSION}}-${{ matrix.use-standard-naming }}**

6. **Print logs ${{ matrix.executor }}-${{ matrix.kubernetes-combo }}-${{ matrix.use-standard-naming }}**
   - Condition: `failure() || cancelled() || inputs.include-success-outputs == 'true'`

7. **Upload KinD logs ${{ matrix.executor }}-${{ matrix.kubernetes-combo }}-${{ matrix.use-standard-naming }}**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `failure() || cancelled() || inputs.include-success-outputs == 'true'`
   - With:
     - `name`: `kind-logs-${{ matrix.kubernetes-combo }}-${{ matrix.executor }}-${{ matrix.use-standard-naming }}`
     - `path`: `/tmp/kind_logs_*`
     - `retention-days`: `7`

8. **Delete clusters just in case they are left**
   - Condition: `always()`

# Milestone Tag Assistant

| Property | Value |
|----------|-------|
| File | `milestone-tag-assistant.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`, `v3-2-test`, `v3-1-test`

## Permissions

- `contents`: `write` - zizmor: ignore[excessive-permissions]
- `pull-requests`: `write` - zizmor: ignore[excessive-permissions]

## Call graph (rooted at this workflow)

```
milestone-tag-assistant.yml [push]
+-- set-milestone / Install Breeze (uses ./.github/actions/breeze)
```

## Jobs

### Get PR information (`get-pr-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Add delay for GitHub to process PR merge**

2. **Find PR information**
   - ID: `pr-info`
   - Uses: `actions/github-script@3a2844b7e9c422d3c10d287c895573f7108da1b3` (v9.0.0)
   - With:
     - `script`: `const { data: pullRequests } = await github.rest.repos.listPullRequestsAssociatedWithCommit({     owner: context.repo.owner,     repo: context.repo.repo,     commit_sha: process.env.GITHUB_SHA });  if (pullRequests.length === 0) {     console.log('⚠️ No pull request found for this commit.');     core.setOutput('should-run', 'false');     return; }  const pr = pullRequests[0];  // Skip if PR already has a milestone if (pr.milestone !== null) {     console.log(`PR #${pr.number} already has milestone: ${pr.milestone.title}`);     core.setOutput('should-run', 'false');     return; }  const labels = pr.labels.map(label => label.name);  console.log(`Commit ${process.env.GITHUB_SHA} is associated with PR #${pr.number}`); console.log(`Title: ${pr.title}`); console.log(`Labels: ${JSON.stringify(labels)}`); console.log(`Base branch: ${pr.base.ref}`); console.log(`Merged by: ${pr.merged_by?.login || 'unknown'}`);  core.setOutput('should-run', 'true'); core.setOutput('pr-number', pr.number.toString()); core.setOutput('pr-title', pr.title); core.setOutput('pr-labels', JSON.stringify(labels)); core.setOutput('base-branch', pr.base.ref); core.setOutput('merged-by', pr.merged_by?.login || 'unknown');`

### Set milestone on merged PR (`set-milestone`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `get-pr-info` |
| Condition | `${{ needs.get-pr-info.outputs.should-run == 'true' }}` |

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

3. **Check criteria and set milestone**

# Notify uv.lock conflicts

| Property | Value |
|----------|-------|
| File | `notify-uv-lock-conflicts.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`
  - paths: `uv.lock`

## Permissions

- `contents`: `read`
- `pull-requests`: `write`

## Jobs

### Notify open PRs that conflict on uv.lock (`notify`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install uv**

3. **Notify open PRs**

# Build PROD images

| Property | Value |
|----------|-------|
| File | `prod-image-build.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `build-type` | string | Yes | - | Name of the 'type' of the build - usually 'Regular' but other types are used to test image variations. |
| `upload-package-artifact` | string | Yes | - | Whether to upload package artifacts (true/false). If false, the job will rely on artifacts prepared by the main prod-image build job. |
| `target-commit-sha` | string | No | - | The commit SHA to checkout for the build |
| `pull-request-target` | string | No | `false` | Whether we are running this from pull-request-target workflow (true/false) |
| `is-committer-build` | string | No | `false` | Whether the build is executed by committer (true/false) |
| `push-image` | string | Yes | - | Whether to push image to the registry (true/false) |
| `upload-image-artifact` | string | No | `false` | Whether to upload docker image artifact |
| `debian-version` | string | No | `bookworm` | Base Debian distribution to use for the build (bookworm) |
| `install-mysql-client-type` | string | No | `mariadb` | MySQL client type to use during build (mariadb/mysql) |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `python-versions` | string | Yes | `[""]` | JSON-formatted array of Python versions to build images from |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `constraints-branch` | string | Yes | - | Branch used to construct constraints URL from. |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to attempt to upgrade image to newer dependencies (true/false) |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |
| `prod-image-build` | string | Yes | - | Whether this is a prod-image build (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
prod-image-build.yml
+-- ci-amd.yml (job: build-prod-images)  <- entry point
+-- ci-arm.yml (job: build-prod-images)  <- entry point
+-- prod-image-extra-checks.yml (job: pip-image)
    +-- additional-prod-image-tests.yml (job: prod-image-extra-checks-main)
    |   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
    |   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
    +-- additional-prod-image-tests.yml (job: prod-image-extra-checks-release-branch)
        +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
        +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CONSTRAINTS_GITHUB_REPOSITORY` | job `build-prod-images` env `CONSTRAINTS_GITHUB_REPOSITORY` |
| `GITHUB_TOKEN` | job `build-prod-images` env `GITHUB_TOKEN` |

## Jobs

### Build Airflow and provider distributions (`build-prod-packages`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.prod-image-build == 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |

#### Steps

1. **Cleanup repo**
   - Condition: `inputs.upload-package-artifact == 'true'`

2. **Checkout target branch**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Make /mnt writeable**
   - Condition: `inputs.upload-package-artifact == 'true'`

4. **Move docker to /mnt**
   - Condition: `inputs.upload-package-artifact == 'true'`

5. **Cleanup dist and context file**
   - Condition: `inputs.upload-package-artifact == 'true'`

6. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ matrix.python-version }}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

7. **Install Breeze**
   - Uses: `./.github/actions/breeze`
   - Condition: `inputs.upload-package-artifact == 'true'`

8. **Prepare providers packages - all providers built from sources**
   - Condition: `inputs.upload-package-artifact == 'true' && inputs.branch == 'main'`

9. **Prepare providers packages with only new versions of providers**
   - Condition: `inputs.upload-package-artifact == 'true' && inputs.branch != 'main'`

10. **Prepare airflow package**
   - Condition: `inputs.upload-package-artifact == 'true'`

11. **Prepare task-sdk package**
   - Condition: `inputs.upload-package-artifact == 'true'`

12. **Prepare airflow-ctl package**
   - Condition: `inputs.upload-package-artifact == 'true'`

13. **Upload prepared packages as artifacts**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `inputs.upload-package-artifact == 'true'`
   - With:
     - `name`: `prod-packages`
     - `path`: `./dist`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

### Build PROD ${{ inputs.build-type }} image ${{ matrix.python-version }} (`build-prod-images`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Depends on | `build-prod-packages` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `sqlite` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.python-version }}` |
| `DEFAULT_BRANCH` | `${{ inputs.branch }}` |
| `DEFAULT_CONSTRAINTS_BRANCH` | `${{ inputs.constraints-branch }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `CONSTRAINTS_GITHUB_REPOSITORY` | `${{ secrets.CONSTRAINTS_GITHUB_REPOSITORY != '' && secrets.CONSTRAINTS_GITHUB_REPOSITORY \|\| 'apache/airflow' }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `PLATFORM` | `${{ inputs.platform }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout target branch**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Make /mnt writeable**

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`

5. **Cleanup dist and context file**

6. **Download packages prepared as artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `prod-packages`
     - `path`: `./docker-context-files`

7. **Download constraints**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `constraints-${{ matrix.python-version }}`
     - `path`: `./docker-context-files/constraints-${{ matrix.python-version }}`

8. **Show downloaded files**

9. **Show constraints**

10. **Login to ghcr.io**

11. **Build PROD images w/ source providers ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**

12. **Verify PROD image ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**

13. **Export PROD docker image ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Condition: `inputs.upload-image-artifact == 'true'`

14. **Stash PROD docker image ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `inputs.upload-image-artifact == 'true'`
   - With:
     - `key`: `prod-image-save-v3-${{ inputs.platform }}-${{ env.PYTHON_MAJOR_MINOR_VERSION }}`
     - `path`: `/mnt/prod-image-save-*-${{ env.PYTHON_MAJOR_MINOR_VERSION }}.tar`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

# PROD images extra checks

| Property | Value |
|----------|-------|
| File | `prod-image-extra-checks.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to build images from |
| `default-python-version` | string | Yes | - | Default Python version to use for the image |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to attempt to upgrade image to newer dependencies (false/RANDOM_VALUE) |
| `constraints-branch` | string | Yes | - | Branch used to construct constraints URL from. |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |

## Permissions

- `contents`: `read`

## Called by

```
prod-image-extra-checks.yml
+-- additional-prod-image-tests.yml (job: prod-image-extra-checks-main)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: prod-image-extra-checks-release-branch)
    +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
    +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
```

## Jobs

### `pip-image`

| Property | Value |
|----------|-------|
| Uses workflow | [Build PROD images](#build-prod-images) |

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `build-type`: `pip`
- `upload-image-artifact`: `false`
- `upload-package-artifact`: `false`
- `install-mysql-client-type`: `mariadb`
- `python-versions`: `${{ inputs.python-versions }}`
- `default-python-version`: `${{ inputs.default-python-version }}`
- `branch`: `${{ inputs.branch }}`
- `push-image`: `false`
- `use-uv`: `false`
- `upgrade-to-newer-dependencies`: `${{ inputs.upgrade-to-newer-dependencies }}`
- `constraints-branch`: `${{ inputs.constraints-branch }}`
- `docker-cache`: `${{ inputs.docker-cache }}`
- `disable-airflow-repo-cache`: `${{ inputs.disable-airflow-repo-cache }}`
- `prod-image-build`: `true`

# Publish Docs to S3

| Property | Value |
|----------|-------|
| File | `publish-docs-to-s3.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ref` | string | Yes | - | The branch or tag to checkout for the docs publishing |
| `destination` | choice | No | `auto` | The destination location in S3<br>Options: `auto`, `live`, `staging` |
| `include-docs` | string | Yes | - | Space separated list of packages to build |
| `exclude-docs` | string | No | `no-docs-excluded` | Comma separated list of docs to exclude |
| `skip-write-to-stable-folder` | boolean | No | `false` | Do not override stable version |
| `build-sboms` | boolean | No | `false` | Build SBOMs |
| `airflow-base-version` | string | No | - | Override the Airflow Base Version to use for the docs build |
| `airflow-version` | string | No | - | Override the Airflow Version to use for the docs build |
| `apply-commits` | string | No | - | Optionally apply commit hashes before building - to patch the docs (coma separated) |
| `ignore-missing-inventories` | boolean | No | `false` | Do not fail the build on missing third-party inventories |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
publish-docs-to-s3.yml [workflow_dispatch]
+-- build-docs / Install Breeze from the ${{ inputs.ref }} reference (uses ./.github/actions/breeze)
+-- publish-docs-to-s3 / Install Breeze (uses ./.github/actions/breeze)
+-- update-registry (uses registry-build.yml)
    +-- build-ci-image (uses ci-image-build.yml)
    |   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
    +-- build-and-publish-registry / Prepare breeze & CI image (uses ./.github/actions/prepare_breeze_and_image)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `DOCS_AWS_ACCESS_KEY_ID`, `DOCS_AWS_SECRET_ACCESS_KEY`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `build-docs` env `GITHUB_TOKEN`; job `publish-docs-to-s3` env `GITHUB_TOKEN` |
| `DOCS_AWS_ACCESS_KEY_ID` | job `publish-docs-to-s3` step `Configure AWS credentials` with `aws-access-key-id`; job `update-registry` secrets `DOCS_AWS_ACCESS_KEY_ID` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `publish-docs-to-s3` step `Configure AWS credentials` with `aws-secret-access-key`; job `update-registry` secrets `DOCS_AWS_SECRET_ACCESS_KEY` |

## Jobs

### Build Info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-24.04` |
| Condition | `contains(fromJSON('[ "ashb", "bugraoz93", "eladkal", "ephraimbuddy", "jedcunningham", "jscheffl", "kaxil", "pierrejeambrun", "shahar1", "potiuk", "utkarsharma2", "vincbeck", "vatsrahul1001", ]'), github.event.sender.login)` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VERBOSE` | `true` |
| `REF` | `${{ inputs.ref }}` |
| `INCLUDE_DOCS` | `${{ inputs.include-docs }}` |
| `EXCLUDE_DOCS` | `${{ inputs.exclude-docs }}` |
| `DESTINATION` | `${{ inputs.destination }}` |
| `SKIP_WRITE_TO_STABLE_FOLDER` | `${{ inputs.skip-write-to-stable-folder }}` |
| `BUILD_SBOMS` | `${{ inputs.build-sboms }}` |
| `AIRFLOW_BASE_VERSION` | `${{ inputs.airflow-base-version \|\| '' }}` |
| `AIRFLOW_VERSION` | `${{ inputs.airflow-version \|\| '' }}` |
| `APPLY_COMMITS` | `${{ inputs.apply-commits \|\| '' }}` |

#### Steps

1. **Checkout for wave provider derivation**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `${{ inputs.ref }}`
     - `fetch-tags`: `true`
     - `fetch-depth`: `0`

2. **Derive registry trigger inputs**
   - ID: `derive_registry_inputs`

3. **Input parameters summary**
   - ID: `parameters`

### Build documentation (`build-docs`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-info` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `false` |
| `VERBOSE` | `true` |
| `EXTRA_BUILD_OPTIONS` | `${{ needs.build-info.outputs.extra-build-options }}` |
| `APPLY_COMMITS` | `${{ inputs.apply-commits \|\| '' }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ needs.build-info.outputs.default-python-version }}` |
| `DOCKER_CACHE` | `registry` |

#### Steps

1. **Cleanup repo**

2. **Checkout current version first to clean-up stuff**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `path`: `current-version`

3. **Free up disk space**

4. **Make /mnt writeable**

5. **Move docker to /mnt**

6. **Copy the version retrieval script**

7. **Checkout ${{ inputs.ref }} **
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `${{ inputs.ref }}`
     - `fetch-tags`: `true`
     - `fetch-depth`: `0`

8. **Apply patch commits if provided**

9. **Install Breeze from the ${{ inputs.ref }} reference**
   - Uses: `./.github/actions/breeze`
   - With:
     - `python-version`: `${{ needs.build-info.outputs.default-python-version }}` - Python version to use

10. **Login to ghcr.io**

11. **Building image from the ${{ inputs.ref }} reference**

12. **Restore docs inventory cache**
   - ID: `restore-docs-inventory-cache`
   - Uses: `apache/infrastructure-actions/stash/restore@49df447b39b18354895520e0a63731b7cad7cbec`
   - With:
     - `path`: `./generated/_inventory_cache/`
     - `key`: `cache-docs-inventory-v1`

13. **Building docs with --docs-only flag using ${{ inputs.ref }} reference breeze**

14. **Save docs inventory cache**
   - Uses: `apache/infrastructure-actions/stash/save@49df447b39b18354895520e0a63731b7cad7cbec`
   - Condition: `steps.restore-docs-inventory-cache.outputs.stash-hit != 'true'`
   - With:
     - `path`: `./generated/_inventory_cache/`
     - `key`: `cache-docs-inventory-v1`
     - `if-no-files-found`: `error`
     - `retention-days`: `2`

15. **Store stable versions**

16. **Saving build docs folder**

17. **Upload build docs**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `airflow-docs`
     - `path`: `/mnt/_build`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`
     - `overwrite`: `true`

### Publish documentation to S3 (`publish-docs-to-s3`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-docs`, `build-info` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`
- `packages`: `write`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `false` |
| `PYTHON_MAJOR_MINOR_VERSION` | `3.10` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout current version with all history for SBOM**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

3. **Make /mnt writeable and cleanup**

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`

5. **Download docs prepared as artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `airflow-docs`
     - `path`: `/mnt/_build`

6. **Move docs to generated folder**

7. **Make sure SBOM dir exists and has the right permissions**
   - Condition: `inputs.build-sboms`

8. **Prepare SBOMs**
   - Condition: `inputs.build-sboms`

9. **Generated SBOM files**
   - Condition: `inputs.build-sboms`

10. **Check disk space available**

11. **Create /mnt/airflow-site directory**

12. **Publish docs to /mnt/airflow-site directory using ${{ inputs.ref }} reference breeze**

13. **Check disk space available**

14. **Update watermarks**
   - Condition: `needs.build-info.outputs.destination == 'staging'`

15. **Install AWS CLI v2**

16. **Configure AWS credentials**
   - Uses: `aws-actions/configure-aws-credentials@d979d5b3a71173a29b74b5b88418bfda9437d885` (v6.1.1)
   - With:
     - `aws-access-key-id`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
     - `aws-secret-access-key`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
     - `aws-region`: `us-east-2`

17. **Syncing docs to S3**

### Update Provider Registry (`update-registry`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build & Publish Registry](#build--publish-registry) |
| Depends on | `publish-docs-to-s3`, `build-info` |
| Condition | `needs.build-info.outputs.registry-providers != '' \|\| needs.build-info.outputs.registry-full-build == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

#### Inputs forwarded

- `destination`: `${{ needs.build-info.outputs.destination }}`
- `provider`: `${{ needs.build-info.outputs.registry-providers }}`

#### Secrets forwarded

- `DOCS_AWS_ACCESS_KEY_ID`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
- `DOCS_AWS_SECRET_ACCESS_KEY`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`

# Push image cache

| Property | Value |
|----------|-------|
| File | `push-image-cache.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `cache-type` | string | Yes | - | Type of cache to push (Early / Regular). |
| `include-prod-images` | string | Yes | - | Whether to build PROD image cache additionally to CI image cache (true/false). |
| `push-latest-images` | string | Yes | - | Whether to also push latest images (true/false). |
| `debian-version` | string | No | `bookworm` | Base Debian distribution to use for the build (bookworm) |
| `install-mysql-client-type` | string | No | `mariadb` | MySQL client type to use during build (mariadb/mysql) |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to build images from |
| `branch` | string | Yes | - | Branch used to run the CI jobs in (main/v*_*_test). |
| `constraints-branch` | string | Yes | - | Branch used to construct constraints URL from. |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs (true/false). |
| `docker-cache` | string | Yes | - | Docker cache specification to build the image (registry, local, disabled). |
| `disable-airflow-repo-cache` | string | Yes | - | Disable airflow repo cache read from main. |

## Called by

```
push-image-cache.yml
+-- additional-ci-image-checks.yml (job: push-early-buildx-cache-to-github-registry)
|   +-- ci-amd.yml (job: additional-ci-image-checks)  <- entry point
|   +-- ci-arm.yml (job: additional-ci-image-checks)  <- entry point
+-- finalize-tests.yml (job: push-buildx-cache-to-github-registry)
    +-- ci-amd.yml (job: finalize-tests)  <- entry point
    +-- ci-arm.yml (job: finalize-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CONSTRAINTS_GITHUB_REPOSITORY` | job `push-ci-image-cache` env `CONSTRAINTS_GITHUB_REPOSITORY`; job `push-prod-image-cache` env `CONSTRAINTS_GITHUB_REPOSITORY` |
| `GITHUB_TOKEN` | job `push-ci-image-cache` env `GITHUB_TOKEN`; job `push-prod-image-cache` env `GITHUB_TOKEN` |

## Jobs

### Push CI ${{ inputs.cache-type }}:${{ matrix.python }} image cache  (`push-ci-image-cache`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COMMIT_SHA` | `${{ github.sha }}` |
| `CONSTRAINTS_GITHUB_REPOSITORY` | `${{ secrets.CONSTRAINTS_GITHUB_REPOSITORY != '' && secrets.CONSTRAINTS_GITHUB_REPOSITORY \|\| 'apache/airflow' }}` |
| `DEBIAN_VERSION` | `${{ inputs.debian-version }}` |
| `DEFAULT_BRANCH` | `${{ inputs.branch }}` |
| `DEFAULT_CONSTRAINTS_BRANCH` | `${{ inputs.constraints-branch }}` |
| `DOCKER_CACHE` | `${{ inputs.docker-cache }}` |
| `DISABLE_AIRFLOW_REPO_CACHE` | `${{ inputs.disable-airflow-repo-cache }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `${{ inputs.include-success-outputs }}` |
| `INSTALL_MYSQL_CLIENT_TYPE` | `${{ inputs.install-mysql-client-type }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.python }}` |
| `UPGRADE_TO_NEWER_DEPENDENCIES` | `false` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Free up disk space**

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`

5. **Login to ghcr.io**

6. **Push CI latest images: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (linux/amd64 only)**
   - Condition: `inputs.push-latest-images == 'true' && inputs.platform == 'linux/amd64'`

7. **Push CI ${{ inputs.cache-type }} cache:${{ env.PYTHON_MAJOR_MINOR_VERSION }}:${{ inputs.platform }}**

### Push PROD ${{ inputs.cache-type }}:${{ matrix.python }} image cache (`push-prod-image-cache`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.include-prod-images == 'true'` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COMMIT_SHA` | `${{ github.sha }}` |
| `CONSTRAINTS_GITHUB_REPOSITORY` | `${{ secrets.CONSTRAINTS_GITHUB_REPOSITORY != '' && secrets.CONSTRAINTS_GITHUB_REPOSITORY \|\| 'apache/airflow' }}` |
| `DEBIAN_VERSION` | `${{ inputs.debian-version }}` |
| `DEFAULT_BRANCH` | `${{ inputs.branch }}` |
| `DEFAULT_CONSTRAINTS_BRANCH` | `${{ inputs.constraints-branch }}` |
| `DOCKER_CACHE` | `${{ inputs.docker-cache }}` |
| `DISABLE_AIRFLOW_REPO_CACHE` | `${{ inputs.disable-airflow-repo-cache }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INSTALL_MYSQL_CLIENT_TYPE` | `${{ inputs.install-mysql-client-type }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.python }}` |
| `UPGRADE_TO_NEWER_DEPENDENCIES` | `false` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Free up disk space**

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`

5. **Cleanup dist and context file**

6. **Download packages prepared as artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `prod-packages`
     - `path`: `./docker-context-files`

7. **Login to ghcr.io**

8. **Push PROD latest image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (linux/amd64 ONLY)**
   - Condition: `inputs.push-latest-images == 'true' && inputs.platform == 'linux/amd64'`

9. **Push PROD ${{ inputs.cache-type }} cache: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} ${{ inputs.platform }}**

# Recheck old bug reports

| Property | Value |
|----------|-------|
| File | `recheck-old-bug-report.yml` |
| Triggers | `schedule` |

## Schedule

- `0 7 * * *`

## Permissions

- `issues`: `write` - All other permissions are set to none

## Jobs

### `recheck-old-bug-report`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **actions/stale@v10.2.0**
   - Uses: `actions/stale@b5d41d4e1d5dceea10e7104786b73624c18a190f` (v10.2.0)
   - With:
     - `only-issue-labels`: `kind:bug`
     - `stale-issue-label`: `Stale Bug Report`
     - `days-before-issue-stale`: `365`
     - `days-before-issue-close`: `30`
     - `days-before-pr-stale`: `-1`
     - `days-before-pr-close`: `-1`
     - `remove-stale-when-updated`: `false`
     - `remove-issue-stale-when-updated`: `true`
     - `labels-to-add-when-unstale`: `needs-triage`
     - `labels-to-remove-when-unstale`: `Stale Bug Report`
     - `stale-issue-message`: `This issue has been automatically marked as stale because it has been open for 365 days without any activity. There has been several Airflow releases since last activity on this issue. Kindly asking to recheck the report against latest Airflow version and let us know if the issue is reproducible. The issue will be closed in next 30 days if no further activity occurs from the issue author.`
     - `close-issue-message`: `This issue has been closed because it has not received response from the issue author.`

# Registry Backfill

| Property | Value |
|----------|-------|
| File | `registry-backfill.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `destination` | choice | Yes | `staging` | Publish to live or staging S3 bucket<br>Options: `staging`, `live` |
| `provider-versions` | string | Yes | - | Space-separated provider/version pairs (e.g. 'amazon/9.24.0 google/21.0.0 celery/3.17.2'). Multiple versions per provider are grouped into one job. |

## Permissions

- `contents`: `read`
- `packages`: `read`

## Call graph (rooted at this workflow)

```
registry-backfill.yml [workflow_dispatch]
+-- build-ci-image (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- backfill / Prepare breeze & CI image (uses ./.github/actions/prepare_breeze_and_image)
+-- publish-versions / Install Breeze (uses ./.github/actions/breeze)
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCS_AWS_ACCESS_KEY_ID` | job `backfill` step `Configure AWS credentials` with `aws-access-key-id`; job `publish-versions` step `Configure AWS credentials` with `aws-access-key-id` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `backfill` step `Configure AWS credentials` with `aws-secret-access-key`; job `publish-versions` step `Configure AWS credentials` with `aws-secret-access-key` |

## Jobs

### Build CI image (`build-ci-image`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Condition | `contains(fromJSON('[<br>  "ashb",<br>  "bugraoz93",<br>  "eladkal",<br>  "ephraimbuddy",<br>  "jedcunningham",<br>  "jscheffl",<br>  "kaxil",<br>  "pierrejeambrun",<br>  "shahar1",<br>  "potiuk",<br>  "utkarsharma2",<br>  "vincbeck"<br>  ]'), github.event.sender.login)` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `false`
- `python-versions`: `["3.12"]`
- `branch`: `main`
- `constraints-branch`: `constraints-main`
- `use-uv`: `true`
- `upgrade-to-newer-dependencies`: `false`
- `docker-cache`: `registry`
- `disable-airflow-repo-cache`: `false`

### `prepare`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Build provider matrix**
   - ID: `matrix`

2. **Determine S3 destination**
   - ID: `destination`

### Backfill ${{ matrix.provider }} (${{ matrix.versions }}) (`backfill`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `prepare`, `build-ci-image` |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Fetch provider tags**

3. **Prepare breeze & CI image**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `python`: `3.12` - Python version for image to prepare (required)
     - `platform`: `linux/amd64` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `use-uv`: `true` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install AWS CLI v2**

5. **Configure AWS credentials**
   - Uses: `aws-actions/configure-aws-credentials@d979d5b3a71173a29b74b5b88418bfda9437d885` (v6.1.1)
   - With:
     - `aws-access-key-id`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
     - `aws-secret-access-key`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
     - `aws-region`: `us-east-2`

6. **Download existing providers.json**

7. **Run breeze registry backfill**

8. **Download data files from S3 for build**

9. **Patch providers.json with backfill version(s)**

10. **Setup pnpm**
   - Uses: `pnpm/action-setup@0e279bb959325dab635dd2c09392533439d90093` (v6.0.8)
   - With:
     - `version`: `10`

11. **Setup Node.js**
   - Uses: `actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e` (v6.4.0)
   - With:
     - `node-version`: `24`
     - `cache`: `pnpm`
     - `cache-dependency-path`: `registry/pnpm-lock.yaml`

12. **Install Node.js dependencies**

13. **Build registry site**

14. **Sync backfilled version pages to S3**

### Publish versions.json (`publish-versions`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `prepare`, `backfill` |

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install Breeze**
   - Uses: `./.github/actions/breeze`
   - With:
     - `python-version`: `3.12` - Python version to use

3. **Install AWS CLI v2**

4. **Configure AWS credentials**
   - Uses: `aws-actions/configure-aws-credentials@d979d5b3a71173a29b74b5b88418bfda9437d885` (v6.1.1)
   - With:
     - `aws-access-key-id`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
     - `aws-secret-access-key`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
     - `aws-region`: `us-east-2`

5. **Download providers.json from S3**

6. **Publish version metadata**

# Build & Publish Registry

| Property | Value |
|----------|-------|
| File | `registry-build.yml` |
| Triggers | `workflow_dispatch`, `workflow_call` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `destination` | choice | Yes | `staging` | Publish to live or staging S3 bucket<br>Options: `staging`, `live` |
| `provider` | string | No | - | Provider ID(s) for incremental build (space-separated, empty = full build) |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `destination` | string | No | `staging` | Publish to live or staging S3 bucket |
| `provider` | string | No | - | Provider ID(s) for incremental build (space-separated, empty = full build) |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `DOCS_AWS_ACCESS_KEY_ID` | Yes | - |
| `DOCS_AWS_SECRET_ACCESS_KEY` | Yes | - |

## Permissions

- `contents`: `read`
- `packages`: `read`

## Call graph (rooted at this workflow)

```
registry-build.yml [workflow_dispatch, workflow_call]
+-- build-ci-image (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- build-and-publish-registry / Prepare breeze & CI image (uses ./.github/actions/prepare_breeze_and_image)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `DOCS_AWS_ACCESS_KEY_ID`, `DOCS_AWS_SECRET_ACCESS_KEY`

## Called by

```
registry-build.yml
+-- publish-docs-to-s3.yml (job: update-registry)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCS_AWS_ACCESS_KEY_ID` | job `build-and-publish-registry` step `Configure AWS credentials` with `aws-access-key-id` |
| `DOCS_AWS_SECRET_ACCESS_KEY` | job `build-and-publish-registry` step `Configure AWS credentials` with `aws-secret-access-key` |

## Jobs

### Build CI image (`build-ci-image`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Condition | `github.event_name == 'workflow_call' \|\| contains(fromJSON('[<br>  "ashb",<br>  "bugraoz93",<br>  "eladkal",<br>  "ephraimbuddy",<br>  "jedcunningham",<br>  "jscheffl",<br>  "kaxil",<br>  "pierrejeambrun",<br>  "shahar1",<br>  "potiuk",<br>  "utkarsharma2",<br>  "vincbeck"<br>  ]'), github.event.sender.login)` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `false`
- `python-versions`: `["3.12"]`
- `branch`: `main`
- `constraints-branch`: `constraints-main`
- `use-uv`: `true`
- `upgrade-to-newer-dependencies`: `false`
- `docker-cache`: `registry`
- `disable-airflow-repo-cache`: `false`

### Build & Publish Registry (`build-and-publish-registry`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-ci-image` |

**Permissions:**

- `contents`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `SCARF_ANALYTICS` | `false` |
| `DO_NOT_TRACK` | `1` |
| `EXISTING_REGISTRY_DIR` | `/tmp/existing-registry` |
| `REGISTRY_DATA_DIR` | `dev/registry` |
| `REGISTRY_PROVIDERS_JSON` | `providers.json` |
| `REGISTRY_MODULES_JSON` | `modules.json` |
| `REGISTRY_SITE_DATA_DIR` | `registry/src/_data` |
| `REGISTRY_SITE_VERSIONS_DIR` | `registry/src/_data/versions` |
| `REGISTRY_SITE_LOGOS_DIR` | `registry/public/logos` |
| `REGISTRY_CACHE_CONTROL` | `public, max-age=300` |

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-tags`: `true`

2. **Prepare breeze & CI image**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `python`: `3.12` - Python version for image to prepare (required)
     - `platform`: `linux/amd64` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `use-uv`: `true` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

3. **Install AWS CLI v2**

4. **Configure AWS credentials**
   - Uses: `aws-actions/configure-aws-credentials@d979d5b3a71173a29b74b5b88418bfda9437d885` (v6.1.1)
   - With:
     - `aws-access-key-id`: `${{ secrets.DOCS_AWS_ACCESS_KEY_ID }}`
     - `aws-secret-access-key`: `${{ secrets.DOCS_AWS_SECRET_ACCESS_KEY }}`
     - `aws-region`: `us-east-2`

5. **Determine S3 destination**
   - ID: `destination`

6. **Download existing registry data from S3**
   - ID: `download-existing`
   - Condition: `inputs.provider != ''`

7. **Extract registry data (breeze)**

8. **Merge with existing registry data**
   - Condition: `inputs.provider != '' && steps.download-existing.outputs.found == 'true'`

9. **Copy breeze output to registry data**

10. **Setup pnpm**
   - Uses: `pnpm/action-setup@0e279bb959325dab635dd2c09392533439d90093` (v6.0.8)
   - With:
     - `version`: `10`

11. **Setup Node.js**
   - Uses: `actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e` (v6.4.0)
   - With:
     - `node-version`: `24`
     - `cache`: `pnpm`
     - `cache-dependency-path`: `registry/pnpm-lock.yaml`

12. **Install Node.js dependencies**

13. **Build registry site**

14. **Upload registry artifact**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `registry-site`
     - `path`: `registry/_site`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

15. **Verify build emitted expected content**

16. **Sync registry to S3**

17. **Publish version metadata**

# Registry Tests

| Property | Value |
|----------|-------|
| File | `registry-tests.yml` |
| Triggers | `pull_request`, `push` |

## Event filters

- **pull_request**
  - branches: `main`
  - paths: `dev/registry/**`, `registry/**`, `providers/*/provider.yaml`, `providers/*/*/provider.yaml`, `.github/workflows/registry-tests.yml`
- **push**
  - branches: `main`
  - paths: `dev/registry/**`

## Permissions

- `contents`: `read`

**Concurrency:** group `registry-tests-${{ github.event.pull_request.number || github.ref }}`, cancel-in-progress: `true`

## Jobs

### Registry extraction tests (`registry-tests`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install uv**
   - Uses: `astral-sh/setup-uv@08807647e7069bb48b6ef5acd8ec9567f424441b` (v8.1.0)
   - With:
     - `python-version`: `3.12`

3. **Run registry extraction tests**

# Release PROD images

| Property | Value |
|----------|-------|
| File | `release_dockerhub_image.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `airflowVersion` | - | Yes | - | Airflow version (e.g. 3.0.1, 3.0.1rc1, 3.0.1b1) |
| `amdOnly` | boolean | No | `false` | Limit to amd64 images |
| `limitPythonVersions` | string | No | - | Force python versions (e.g. "3.10 3.11") |

## Permissions

- `contents`: `read`
- `packages`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `VERBOSE` | `true` |

**Concurrency:** group `${{ github.event.inputs.airflowVersion }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
release_dockerhub_image.yml [workflow_dispatch]
+-- build-info / Install Breeze (uses ./.github/actions/breeze)
+-- release-images (uses release_single_dockerhub_image.yml)
    +-- build-images / Install Breeze (uses ./.github/actions/breeze)
    +-- merge-images / Install Breeze (uses ./.github/actions/breeze)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `DOCKERHUB_TOKEN`, `DOCKERHUB_USER`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `DOCKERHUB_USER` | job `release-images` secrets `DOCKERHUB_USER` |
| `DOCKERHUB_TOKEN` | job `release-images` secrets `DOCKERHUB_TOKEN` |

## Jobs

### Build Info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-24.04` |
| Condition | `contains(fromJSON('[ "ashb", "bugraoz93", "eladkal", "ephraimbuddy", "jedcunningham", "jscheffl", "kaxil", "pierrejeambrun", "potiuk", "utkarsharma2", "vincbeck", "vatsrahul1001", ]'), github.event.sender.login)` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VERBOSE` | `true` |
| `AIRFLOW_VERSION` | `${{ github.event.inputs.airflowVersion }}` |
| `AMD_ONLY` | `${{ github.event.inputs.amdOnly }}` |
| `LIMIT_PYTHON_VERSIONS` | `${{ github.event.inputs.limitPythonVersions }}` |

#### Steps

1. **Input parameters summary**

2. **Cleanup repo**

3. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

4. **Install Breeze**
   - Uses: `./.github/actions/breeze`

5. **Save github context to file**

6. **Selective checks**
   - ID: `selective-checks`

7. **Check airflow version**
   - ID: `check-airflow-version`

8. **Determine build matrix**
   - ID: `determine-matrix`

9. **Determine python versions**
   - ID: `determine-python-versions`

### Release images (`release-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Release single PROD image](#release-single-prod-image) |
| Depends on | `build-info` |

**Permissions:**

- `contents`: `read`

#### Inputs forwarded

- `pythonVersion`: `${{ matrix.python }}`
- `airflowVersion`: `${{ needs.build-info.outputs.airflowVersion }}`
- `platformMatrix`: `${{ needs.build-info.outputs.platformMatrix }}`
- `skipLatest`: `${{ needs.build-info.outputs.skipLatest }}`
- `armRunners`: `${{ needs.build-info.outputs.arm-runners }}`
- `amdRunners`: `${{ needs.build-info.outputs.amd-runners }}`

#### Secrets forwarded

- `DOCKERHUB_USER`: `${{ secrets.DOCKERHUB_USER }}`
- `DOCKERHUB_TOKEN`: `${{ secrets.DOCKERHUB_TOKEN }}`

# Release single PROD image

| Property | Value |
|----------|-------|
| File | `release_single_dockerhub_image.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `airflowVersion` | string | Yes | - | Airflow version (e.g. 3.0.1, 3.0.1rc1, 3.0.1b1) |
| `platformMatrix` | string | Yes | - | Platform matrix formatted as json (e.g. ["linux/amd64", "linux/arm64"]) |
| `pythonVersion` | string | Yes | - | Python version (e.g. 3.10, 3.11) |
| `skipLatest` | string | Yes | - | Skip tagging latest release (true/false) |
| `amdRunners` | string | Yes | - | Amd64 runners (e.g. ["ubuntu-22.04", "ubuntu-24.04"]) |
| `armRunners` | string | Yes | - | Arm64 runners (e.g. ["ubuntu-22.04", "ubuntu-24.04"]) |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `DOCKERHUB_USER` | Yes | - |
| `DOCKERHUB_TOKEN` | Yes | - |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `VERBOSE` | `true` |

## Called by

```
release_single_dockerhub_image.yml
+-- release_dockerhub_image.yml (job: release-images)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `DOCKERHUB_TOKEN` | job `build-images` step `Login to hub.docker.com` (run); job `merge-images` step `Login to hub.docker.com` (run) |
| `DOCKERHUB_USER` | job `build-images` step `Login to hub.docker.com` (run); job `merge-images` step `Login to hub.docker.com` (run) |

## Jobs

### Build: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}, ${{ matrix.platform }} (`build-images`)

| Property | Value |
|----------|-------|
| Runs on | `${{ (matrix.platform == 'linux/amd64') && fromJSON(inputs.amdRunners) \|\| fromJSON(inputs.armRunners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `AIRFLOW_VERSION` | `${{ inputs.airflowVersion }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.pythonVersion }}` |
| `PLATFORM` | `${{ matrix.platform }}` |
| `SKIP_LATEST` | `${{ inputs.skipLatest == 'true' && '--skip-latest' \|\| '' }}` |
| `COMMIT_SHA` | `${{ github.sha }}` |
| `REPOSITORY` | `${{ github.repository }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Free space**

5. **Cleanup dist and context file**

6. **Login to hub.docker.com**

7. **Get env vars for metadata**

8. **Login to ghcr.io**

9. **Install buildx plugin**

10. **Create airflow_cache builder**

11. **Build regular images: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}, ${{ matrix.platform }}
**

12. **Verify regular image: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}, ${{ matrix.platform }}
**

13. **Release slim images: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}, ${{ matrix.platform }}
**

14. **Verify slim image: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}, ${{ matrix.platform }}
**

15. **List upload-able artifacts**

16. **Upload metadata artifact ${{ env.ARTIFACT_NAME }}**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - With:
     - `name`: `${{ env.ARTIFACT_NAME }}`
     - `path`: `./dist/metadata-*`
     - `retention-days`: `7`
     - `if-no-files-found`: `error`

17. **Docker logout**
   - Condition: `always()`

### Merge: ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }} (`merge-images`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-images` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `AIRFLOW_VERSION` | `${{ inputs.airflowVersion }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.pythonVersion }}` |
| `SKIP_LATEST` | `${{ inputs.skipLatest == 'true' && '--skip-latest' \|\| '' }}` |
| `COMMIT_SHA` | `${{ github.sha }}` |
| `REPOSITORY` | `${{ github.repository }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install Breeze**
   - Uses: `./.github/actions/breeze`

4. **Free space**

5. **Cleanup dist and context file**

6. **Login to hub.docker.com**

7. **Login to ghcr.io**

8. **Download metadata artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `path`: `./dist`
     - `pattern`: `metadata-${{ inputs.pythonVersion }}-*`

9. **List downloaded artifacts**

10. **Install buildx plugin**

11. **Install regctl**

12. **Merge regular images ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}**

13. **Merge slim images ${{ inputs.airflowVersion }}, ${{ inputs.pythonVersion }}**

14. **Docker logout**
   - Condition: `always()`

# Unit tests

| Property | Value |
|----------|-------|
| File | `run-unit-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining public AMD runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `test-group` | string | Yes | - | Test group to run: ('core', 'providers') |
| `test-types-as-strings-in-json` | string | Yes | - | The list of list of test types to run (types in item are separated by spaces) as json |
| `backend` | string | Yes | - | The backend to run the tests on |
| `test-scope` | string | Yes | - | The scope of the test to run: ('DB', 'Non-DB', 'All') |
| `test-name` | string | Yes | - | The name of the test to run |
| `test-name-separator` | string | No | `:` | The separator to use after the test name |
| `python-versions` | string | Yes | - | The list of python versions (stringified JSON array) to run the tests on. |
| `backend-versions` | string | Yes | - | The list of backend versions (stringified JSON array) to run the tests on. |
| `excluded-providers-as-string` | string | Yes | - | Excluded providers (per Python version) as json string |
| `excludes` | string | Yes | - | Excluded combos (stringified JSON array of python-version/backend-version dicts) |
| `run-migration-tests` | string | No | `false` | Whether to run migration tests or not (true/false) |
| `run-coverage` | string | Yes | - | Whether to run coverage or not (true/false) |
| `debug-resources` | string | Yes | - | Whether to debug resources or not (true/false) |
| `include-success-outputs` | string | No | `false` | Whether to include success outputs or not (true/false) |
| `downgrade-sqlalchemy` | string | No | `false` | Whether to downgrade SQLAlchemy or not (true/false) |
| `upgrade-sqlalchemy` | string | No | `false` | Whether to upgrade SQLAlchemy or not (true/false) |
| `upgrade-boto` | string | No | `false` | Whether to upgrade boto or not (true/false) |
| `downgrade-pendulum` | string | No | `false` | Whether to downgrade pendulum or not (true/false) |
| `force-lowest-dependencies` | string | No | `false` | Whether to force lowest dependencies for the tests or not (true/false) |
| `monitor-delay-time-in-seconds` | number | No | `20` | How much time to wait between printing parallel monitor summary |
| `skip-providers-tests` | string | Yes | - | Whether to skip providers tests or not (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv |
| `default-branch` | string | Yes | - | The default branch of the repository |

## Permissions

- `contents`: `read`

## Called by

```
run-unit-tests.yml
+-- ci-amd.yml (job: tests-postgres-core)  <- entry point
+-- ci-amd.yml (job: tests-postgres-providers)  <- entry point
+-- ci-amd.yml (job: tests-mysql-core)  <- entry point
+-- ci-amd.yml (job: tests-mysql-providers)  <- entry point
+-- ci-amd.yml (job: tests-sqlite-core)  <- entry point
+-- ci-amd.yml (job: tests-sqlite-providers)  <- entry point
+-- ci-amd.yml (job: tests-non-db-core)  <- entry point
+-- ci-amd.yml (job: tests-non-db-providers)  <- entry point
+-- ci-amd.yml (job: tests-with-lowest-direct-resolution-core)  <- entry point
+-- ci-amd.yml (job: tests-with-lowest-direct-resolution-providers)  <- entry point
+-- ci-arm.yml (job: tests-postgres-core)  <- entry point
+-- ci-arm.yml (job: tests-postgres-providers)  <- entry point
+-- ci-arm.yml (job: tests-mysql-core)  <- entry point
+-- ci-arm.yml (job: tests-mysql-providers)  <- entry point
+-- ci-arm.yml (job: tests-sqlite-core)  <- entry point
+-- ci-arm.yml (job: tests-sqlite-providers)  <- entry point
+-- ci-arm.yml (job: tests-non-db-core)  <- entry point
+-- ci-arm.yml (job: tests-non-db-providers)  <- entry point
+-- ci-arm.yml (job: tests-with-lowest-direct-resolution-core)  <- entry point
+-- ci-arm.yml (job: tests-with-lowest-direct-resolution-providers)  <- entry point
+-- special-tests.yml (job: tests-min-sqlalchemy)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-min-sqlalchemy-providers)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-latest-sqlalchemy)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-latest-sqlalchemy-providers)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-boto-core)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-boto-providers)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-pendulum-2-core)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-pendulum-2-providers)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-quarantined-core)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-quarantined-providers)
|   +-- ci-amd.yml (job: tests-special)  <- entry point
|   +-- ci-arm.yml (job: tests-special)  <- entry point
+-- special-tests.yml (job: tests-system-core)
    +-- ci-amd.yml (job: tests-special)  <- entry point
    +-- ci-arm.yml (job: tests-special)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `tests` env `GITHUB_TOKEN` |
| `CODECOV_TOKEN` | job `tests` step `Post Tests success` with `codecov-token` |

## Jobs

### ${{ inputs.test-scope == 'All' && '' || inputs.test-scope == 'Quarantined' && 'Qrnt' || inputs.test-scope }}${{ inputs.test-scope == 'All' && '' || '-' }}${{ inputs.test-group == 'providers' && 'prov' || inputs.test-group}}:${{ inputs.test-name }}${{ inputs.test-name-separator }}${{ matrix.backend-version }}:${{ matrix.python-version}}:${{ matrix.test-types.description }} (`tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.test-group == 'core' \|\| inputs.skip-providers-tests != 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `BACKEND` | `${{ inputs.backend }}` |
| `BACKEND_VERSION` | `${{ matrix.backend-version }}` |
| `DB_RESET` | `true` |
| `DEBUG_RESOURCES` | `${{ inputs.debug-resources }}` |
| `DOWNGRADE_SQLALCHEMY` | `${{ inputs.downgrade-sqlalchemy }}` |
| `DOWNGRADE_PENDULUM` | `${{ inputs.downgrade-pendulum }}` |
| `ENABLE_COVERAGE` | `${{ inputs.run-coverage }}` |
| `EXCLUDED_PROVIDERS` | `${{ inputs.excluded-providers-as-string }}` |
| `FORCE_LOWEST_DEPENDENCIES` | `${{ inputs.force-lowest-dependencies }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_SUCCESS_OUTPUTS` | `${{ inputs.include-success-outputs }}` |
| `PLATFORM` | `${{ inputs.platform }}` |
| `JOB_ID` | `${{ inputs.test-group }}-${{ matrix.test-types.description }}-${{ inputs.test-scope }}-${{ inputs.test-name }}-${{inputs.backend}}-${{ matrix.backend-version }}-${{ matrix.python-version }}` |
| `MOUNT_SOURCES` | `skip` |
| `PARALLEL_TEST_TYPES` | `${{ matrix.test-types.test_types }}` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.python-version }}` |
| `UPGRADE_BOTO` | `${{ inputs.upgrade-boto }}` |
| `UPGRADE_SQLALCHEMY` | `${{ inputs.upgrade-sqlalchemy }}` |
| `AIRFLOW_MONITOR_DELAY_TIME_IN_SECONDS` | `${{inputs.monitor-delay-time-in-seconds}}` |
| `VERBOSE` | `true` |
| `DEFAULT_BRANCH` | `${{ inputs.default-branch }}` |
| `TOTAL_TEST_TIMEOUT` | `3600` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Make /mnt writeable**

4. **Move docker to /mnt**

5. **Prepare breeze & CI image: ${{ matrix.python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ matrix.python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `false` - Whether to cleanup /mnt (required)

6. **Migration Tests: ${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
**
   - Uses: `./.github/actions/migration_tests`
   - Condition: `inputs.run-migration-tests == 'true' && inputs.test-group == 'core' && matrix.python-version != '3.14'`
   - With:
     - `python-version`: `${{ matrix.python-version }}` - Python version to run the tests on (required)

7. **${{ inputs.test-group }}:${{ inputs.test-scope }} Tests ${{ inputs.test-name }} ${{ matrix.backend-version }} Py${{ matrix.python-version }}:${{ env.PARALLEL_TEST_TYPES }}
**

8. **Post Tests success**
   - Uses: `./.github/actions/post_tests_success`
   - Condition: `success()`
   - With:
     - `codecov-token`: `${{ secrets.CODECOV_TOKEN }}` - Codecov token (required)
     - `python-version`: `${{ matrix.python-version }}` - Python version (required)

9. **Post Tests failure**
   - Uses: `./.github/actions/post_tests_failure`
   - Condition: `failure() || cancelled()`

# [main] Scheduled CI upgrade check

| Property | Value |
|----------|-------|
| File | `scheduled-upgrade-check-main.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 6 * * 1,3,5`

## Permissions

- `contents`: `write`
- `pull-requests`: `write`

## Call graph (rooted at this workflow)

```
scheduled-upgrade-check-main.yml [schedule, workflow_dispatch]
+-- upgrade-main (uses upgrade-check.yml)
    +-- createupgrade-check / [${{ inputs.target-branch }}] Install Breeze (uses ./.github/actions/breeze)
    +-- createupgrade-check / [${{ inputs.target-branch }}] Install prek (uses ./.github/actions/install-prek)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `SLACK_BOT_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SLACK_BOT_TOKEN` | job `upgrade-main` secrets `SLACK_BOT_TOKEN` |

## Jobs

### [main] Upgrade (`upgrade-main`)

| Property | Value |
|----------|-------|
| Uses workflow | [Upgrade check](#upgrade-check) |

#### Inputs forwarded

- `target-branch`: `main`

#### Secrets forwarded

- `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_BOT_TOKEN }}`

# [v3-2-test] Scheduled CI upgrade check

| Property | Value |
|----------|-------|
| File | `scheduled-upgrade-check-v3-2-test.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 6 * * 2,4`

## Permissions

- `contents`: `write`
- `pull-requests`: `write`

## Call graph (rooted at this workflow)

```
scheduled-upgrade-check-v3-2-test.yml [schedule, workflow_dispatch]
+-- upgrade-v3-2-test (uses upgrade-check.yml)
    +-- createupgrade-check / [${{ inputs.target-branch }}] Install Breeze (uses ./.github/actions/breeze)
    +-- createupgrade-check / [${{ inputs.target-branch }}] Install prek (uses ./.github/actions/install-prek)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `SLACK_BOT_TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SLACK_BOT_TOKEN` | job `upgrade-v3-2-test` secrets `SLACK_BOT_TOKEN` |

## Jobs

### [v3-2-test] Upgrade (`upgrade-v3-2-test`)

| Property | Value |
|----------|-------|
| Uses workflow | [Upgrade check](#upgrade-check) |

#### Inputs forwarded

- `target-branch`: `v3-2-test`

#### Secrets forwarded

- `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_BOT_TOKEN }}`

# Scheduled verify release calendar

| Property | Value |
|----------|-------|
| File | `scheduled-verify-release-calendar.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 6 * * *`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SLACK_BOT_TOKEN` | job `verify-release-calendar` step `Notify Slack on failure` with `token` |

## Jobs

### Verify release calendar (`verify-release-calendar`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install uv**

3. **Verify release calendar**

4. **Notify Slack on failure**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c`
   - Condition: `failure()`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ secrets.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "release-management" text: >-   :warning: Release calendar verification failed.   See:   ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }} blocks:   - type: section     text:       type: mrkdwn       text: >-         :warning: *Release calendar verification failed*          The scheduled `verify_release_calendar.py` check         failed. Please review and fix the mismatch between         the Confluence release wiki and the Google         Calendar entries.          • <https://cwiki.apache.org/confluence/display/AIRFLOW/Release+Plan|Release Plan wiki>          • <https://calendar.google.com/calendar/u/0?cid=Y19kZTIxNGU5MmRmM2I3NTk3NzljYjY1ZjNlNDllNTYyNzk2YzYxMjZlNzUwMGNmYTdlNTI0YmY3ODE4NmQ4YjVlQGdyb3VwLmNhbGVuZGFyLmdvb2dsZS5jb20|Release Calendar>          • <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View failed run>`

# Special tests

| Property | Value |
|----------|-------|
| File | `special-tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-branch` | string | Yes | - | The default branch for the repository |
| `core-test-types-list-as-strings-in-json` | string | Yes | - | The list of core test types to run separated by spaces |
| `providers-test-types-list-as-strings-in-json` | string | Yes | - | The list of providers test types to run separated by spaces |
| `run-coverage` | string | Yes | - | Whether to run coverage or not (true/false) |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `excluded-providers-as-string` | string | Yes | - | Excluded providers (per Python version) as json string |
| `python-versions` | string | Yes | - | The list of python versions (stringified JSON array) to run the tests on. |
| `default-postgres-version` | string | Yes | - | The default version of the postgres to use |
| `canary-run` | string | Yes | - | Whether to run canary tests or not (true/false) |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies or not (true/false) |
| `include-success-outputs` | string | Yes | - | Whether to include success outputs or not (true/false) |
| `debug-resources` | string | Yes | - | Whether to debug resources or not (true/false) |
| `skip-providers-tests` | string | Yes | - | Whether to skip providers tests or not (true/false) |
| `use-uv` | string | Yes | - | Whether to use uv or not (true/false) |

## Permissions

- `contents`: `read`

## Called by

```
special-tests.yml
+-- ci-amd.yml (job: tests-special)  <- entry point
+-- ci-arm.yml (job: tests-special)  <- entry point
```

## Jobs

### Min SQLAlchemy test: core (`tests-min-sqlalchemy`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `downgrade-sqlalchemy`: `true`
- `test-name`: `MinSQLAlchemy-Postgres`
- `test-scope`: `DB`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Min SQLAlchemy test: providers (`tests-min-sqlalchemy-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `downgrade-sqlalchemy`: `true`
- `test-name`: `MinSQLAlchemy-Postgres`
- `test-scope`: `DB`
- `test-group`: `providers`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.providers-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Latest SQLAlchemy test: core (`tests-latest-sqlalchemy`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `upgrade-sqlalchemy`: `true`
- `test-name`: `LatestSQLAlchemy-Postgres`
- `test-scope`: `DB`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Latest SQLAlchemy test: providers (`tests-latest-sqlalchemy-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `upgrade-sqlalchemy`: `true`
- `test-name`: `LatestSQLAlchemy-Postgres`
- `test-scope`: `DB`
- `test-group`: `providers`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.providers-test-types-list-as-strings-in-json }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Latest Boto test: core (`tests-boto-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `upgrade-boto`: `true`
- `test-name`: `LatestBoto-Postgres`
- `test-scope`: `All`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Latest Boto test: providers (`tests-boto-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `upgrade-boto`: `true`
- `test-name`: `LatestBoto-Postgres`
- `test-scope`: `All`
- `test-group`: `providers`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Pendulum2 test: core (`tests-pendulum-2-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `downgrade-pendulum`: `true`
- `test-name`: `Pendulum2-Postgres`
- `test-scope`: `All`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Pendulum2 test: providers (`tests-pendulum-2-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `downgrade-pendulum`: `true`
- `test-name`: `Pendulum2-Postgres`
- `test-scope`: `All`
- `test-group`: `providers`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Quarantined test: core (`tests-quarantined-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `test-name`: `Postgres`
- `test-scope`: `Quarantined`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### Quarantined test: providers (`tests-quarantined-providers`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `test-name`: `Postgres`
- `test-scope`: `Quarantined`
- `test-group`: `providers`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.providers-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

### System test: ${{ matrix.test-group }} (`tests-system-core`)

| Property | Value |
|----------|-------|
| Uses workflow | [Unit tests](#unit-tests) |

**Permissions:**

- `contents`: `read`
- `packages`: `read`

#### Inputs forwarded

- `runners`: `${{ inputs.runners }}`
- `platform`: `${{ inputs.platform }}`
- `test-name`: `SystemTest`
- `test-scope`: `System`
- `test-group`: `core`
- `backend`: `postgres`
- `python-versions`: `['${{ inputs.default-python-version }}']`
- `backend-versions`: `['${{ inputs.default-postgres-version }}']`
- `excluded-providers-as-string`: `${{ inputs.excluded-providers-as-string }}`
- `excludes`: `[]`
- `test-types-as-strings-in-json`: `${{ inputs.core-test-types-list-as-strings-in-json }}`
- `include-success-outputs`: `${{ inputs.include-success-outputs }}`
- `run-coverage`: `${{ inputs.run-coverage }}`
- `debug-resources`: `${{ inputs.debug-resources }}`
- `skip-providers-tests`: `${{ inputs.skip-providers-tests }}`
- `use-uv`: `${{ inputs.use-uv }}`
- `default-branch`: `${{ inputs.default-branch }}`

# Close stale PRs & Issues

| Property | Value |
|----------|-------|
| File | `stale.yml` |
| Triggers | `schedule` |

## Schedule

- `0 0 * * *`

## Permissions

- `pull-requests`: `write` - All other permissions are set to none
- `issues`: `write`

## Jobs

### `stale`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **actions/stale@v10.2.0**
   - Uses: `actions/stale@b5d41d4e1d5dceea10e7104786b73624c18a190f` (v10.2.0)
   - With:
     - `stale-pr-message`: `This pull request has been automatically marked as stale because it has not had recent activity. It will be closed in 5 days if no further activity occurs. Thank you for your contributions.`
     - `days-before-pr-stale`: `45`
     - `days-before-pr-close`: `5`
     - `exempt-pr-labels`: `pinned,security,pending-response`
     - `only-issue-labels`: `pending-response`
     - `remove-stale-when-updated`: `true`
     - `days-before-issue-stale`: `14`
     - `days-before-issue-close`: `7`
     - `stale-issue-message`: `This issue has been automatically marked as stale because it has been open for 14 days with no response from the author. It will be closed in next 7 days if no further activity occurs from the issue author.`
     - `close-issue-message`: `This issue has been closed because it has not received response from the issue author.`

2. **actions/stale@v10.2.0**
   - Uses: `actions/stale@b5d41d4e1d5dceea10e7104786b73624c18a190f` (v10.2.0)
   - With:
     - `only-pr-labels`: `pending-response`
     - `days-before-pr-stale`: `7`
     - `days-before-pr-close`: `7`
     - `stale-pr-message`: `This pull request has been automatically marked as stale because the author has not responded to a request for more information. It will be closed in 7 days if no further activity occurs. Thank you for your contributions.`
     - `close-pr-message`: `This pull request has been closed because the author has not responded to a request for more information.`
     - `labels-to-remove-when-unstale`: `pending-response,stale`
     - `remove-stale-when-updated`: `true`
     - `days-before-issue-stale`: `-1`
     - `days-before-issue-close`: `-1`

# Provider tests

| Property | Value |
|----------|-------|
| File | `test-providers.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runners` | string | Yes | - | The array of labels (in json form) determining public AMD runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `canary-run` | string | Yes | - | Whether this is a canary run |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `upgrade-to-newer-dependencies` | string | Yes | - | Whether to upgrade to newer dependencies |
| `selected-providers-list-as-string` | string | No | - | List of affected providers as string |
| `providers-compatibility-tests-matrix` | string | Yes | - | JSON-formatted array of providers compatibility tests in the form of array of dicts (airflow-version, python-versions, remove-providers, run-unit-tests) |
| `providers-test-types-list-as-strings-in-json` | string | Yes | - | List of parallel provider test types as string |
| `skip-providers-tests` | string | Yes | - | Whether to skip provider tests (true/false) |
| `python-versions` | string | Yes | - | JSON-formatted array of Python versions to build images from |
| `use-uv` | string | Yes | - | Whether to use uv |

## Permissions

- `contents`: `read`

## Called by

```
test-providers.yml
+-- ci-amd.yml (job: providers)  <- entry point
+-- ci-arm.yml (job: providers)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `prepare-install-verify-provider-distributions` env `GITHUB_TOKEN`; job `providers-compatibility-tests-matrix` env `GITHUB_TOKEN` |

## Jobs

### Providers wheel, sdist tests (`prepare-install-verify-provider-distributions`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version }}` |
| `VERBOSE` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ inputs.default-python-version }}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

4. **Prepare breeze & CI image: ${{ inputs.default-python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ inputs.default-python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

5. **Cleanup dist files**

6. **Set current date as RELEASE_DATE variable**
   - ID: `date`

7. **Prepare provider documentation**
   - Condition: `matrix.package-format == 'wheel'`

8. **Prepare provider distributions: ${{ matrix.package-format }}**

9. **Prepare airflow package: ${{ matrix.package-format }}**

10. **Prepare task-sdk package: ${{ matrix.package-format }}**

11. **Prepare airflow-ctl package: ${{ matrix.package-format }}**

12. **Verify ${{ matrix.package-format }} packages with twine**

13. **Test providers issue generation automatically**
   - Condition: `matrix.package-format == 'wheel'`

14. **Generate source constraints from CI image**

15. **Install and verify wheel provider distributions**
   - Condition: `matrix.package-format == 'wheel'`

16. **Install all sdist provider distributions and airflow**
   - Condition: `matrix.package-format == 'sdist'`

### Compat ${{ matrix.compat.airflow-version }}:P${{ matrix.compat.python-version }}:${{ matrix.compat.test-types.description }} (`providers-compatibility-tests-matrix`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners) }}` |
| Condition | `inputs.skip-providers-tests != 'true'` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `INCLUDE_NOT_READY_PROVIDERS` | `true` |
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ matrix.compat.python-version }}` |
| `VERBOSE` | `true` |
| `CLEAN_AIRFLOW_INSTALLATION` | `true` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ matrix.compat.python-version }}` - Python version to use
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

4. **Prepare breeze & CI image: ${{ matrix.compat.python-version }}**
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `python`: `${{ matrix.compat.python-version }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

5. **Cleanup dist files**

6. **Prepare provider distributions: wheel**

7. **Remove incompatible Airflow ${{ matrix.compat.airflow-version }}:Python ${{ matrix.compat.python-version }} provider distributions**
   - Condition: `matrix.compat.remove-providers != ''`

8. **Download airflow package: wheel**

9. **Install and verify all provider distributions and airflow on Airflow ${{ matrix.compat.airflow-version }}:Python ${{ matrix.compat.python-version }}
**
   - Condition: `matrix.compat.run-unit-tests != 'true'`

10. **Check amount of disk space available**

11. **Run provider unit tests on Airflow ${{ matrix.compat.airflow-version }}:Python ${{ matrix.compat.python-version }}:${{ matrix.test-types.description }}
**
   - Condition: `matrix.compat.run-unit-tests == 'true'`

# UI End-to-End Tests

| Property | Value |
|----------|-------|
| File | `ui-e2e-tests.yml` |
| Triggers | `workflow_dispatch`, `workflow_call` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow-name` | string | Yes | - | Name of the test |
| `runners` | string | No | `["ubuntu-24.04"]` | The array of labels (in json form) determining runners. |
| `platform` | string | No | `linux/amd64` | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-python-version` | string | No | `3.10` | Which version of python should be used by default |
| `use-uv` | string | No | `true` | Whether to use uv to build the image (true/false) |
| `docker-image-tag` | string | Yes | - | Tag of the Docker image to test |
| `browser` | string | No | `all` | Browser to test (chromium, firefox, webkit, all) |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow-name` | string | Yes | - | Name of the test |
| `runners` | string | Yes | - | The array of labels (in json form) determining runners. |
| `platform` | string | Yes | - | Platform for the build - 'linux/amd64' or 'linux/arm64' |
| `default-python-version` | string | Yes | - | Which version of python should be used by default |
| `use-uv` | string | Yes | - | Whether to use uv to build the image (true/false) |
| `docker-image-tag` | string | No | - | Tag of the Docker image to test |
| `browser` | string | No | `all` | Browser to test (chromium, firefox, webkit, all) |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
ui-e2e-tests.yml [workflow_dispatch, workflow_call]
+-- test-ui-e2e-tests / Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }} (uses ./.github/actions/prepare_breeze_and_image)
+-- test-ui-e2e-tests / Install Breeze (manual trigger) (uses ./.github/actions/breeze)
```

## Called by

```
ui-e2e-tests.yml
+-- additional-prod-image-tests.yml (job: test-ui-e2e-chromium)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-ui-e2e-firefox)
|   +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
|   +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
+-- additional-prod-image-tests.yml (job: test-ui-e2e-webkit)
    +-- ci-amd.yml (job: additional-prod-image-tests)  <- entry point
    +-- ci-arm.yml (job: additional-prod-image-tests)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `test-ui-e2e-tests` env `GITHUB_TOKEN` |

## Jobs

### ${{ inputs.workflow-name || 'UI E2E Tests' }} (`test-ui-e2e-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ fromJSON(inputs.runners \|\| '["ubuntu-24.04"]') }}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_MAJOR_MINOR_VERSION` | `${{ inputs.default-python-version \|\| '3.10' }}` |
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |
| `BROWSER` | `${{ inputs.browser \|\| 'all' }}` |
| `PLATFORM` | `${{ inputs.platform \|\| 'linux/amd64' }}` |
| `USE_UV` | `${{ inputs.use-uv \|\| 'true' }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

3. **Prepare breeze & PROD image: ${{ env.PYTHON_MAJOR_MINOR_VERSION }}**
   - ID: `breeze`
   - Uses: `./.github/actions/prepare_breeze_and_image`
   - Condition: `github.event_name != 'workflow_dispatch'`
   - With:
     - `platform`: `${{ inputs.platform }}` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `image-type`: `prod` - Which image type to prepare (ci/prod)
     - `python`: `${{ env.PYTHON_MAJOR_MINOR_VERSION }}` - Python version for image to prepare (required)
     - `use-uv`: `${{ inputs.use-uv }}` - Whether to use uv (required)
     - `make-mnt-writeable-and-cleanup`: `true` - Whether to cleanup /mnt (required)

4. **Install Breeze (manual trigger)**
   - Uses: `./.github/actions/breeze`
   - Condition: `github.event_name == 'workflow_dispatch'`

5. **Setup pnpm**
   - Uses: `pnpm/action-setup@0e279bb959325dab635dd2c09392533439d90093` (v6.0.8)
   - With:
     - `version`: `9`
     - `run_install`: `false`

6. **Setup node**
   - Uses: `actions/setup-node@48b55a011bda9f5d6aeb4c2d9c7362e8dae4041e` (v6.4.0)
   - With:
     - `node-version`: `24`

7. **Compile UI assets (for image build fallback)**
   - Condition: `github.event_name == 'workflow_dispatch'`

8. **Install Playwright browsers and dependencies**

9. **Test UI e2e tests**

10. **Upload test results**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `always()`
   - With:
     - `name`: `playwright-report-${{ env.BROWSER }}`
     - `path`: `airflow-core/src/airflow/ui/playwright-report/ airflow-core/src/airflow/ui/test-results/`
     - `retention-days`: `7`
     - `if-no-files-found`: `warn`

11. **Extract E2E test failures and fixme tests**
   - Condition: `always()`

12. **Upload E2E test report**
   - Uses: `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (v7.0.1)
   - Condition: `always()`
   - With:
     - `name`: `e2e-test-report-${{ env.BROWSER }}`
     - `path`: `e2e-test-report/`
     - `retention-days`: `14`
     - `if-no-files-found`: `warn`

# Update constraints on push for stable branch (always)

| Property | Value |
|----------|-------|
| File | `update-constraints-on-push-stable.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `v[0-9]+-[0-9]+-stable`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
update-constraints-on-push-stable.yml [push]
+-- build-info / Install Breeze (uses ./.github/actions/breeze)
+-- build-ci-images (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- generate-constraints (uses generate-constraints.yml)
    +-- generate-constraints-matrix / Install prek (uses ./.github/actions/install-prek)
    +-- generate-constraints-matrix / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | job `notify-on-failure` env `SLACK_BOT_TOKEN` |

## Jobs

### Build info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Fetch incoming commit ${{ github.sha }} with its parent**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ github.sha }}`
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

4. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

5. **Save github context to file**

6. **Selective checks**
   - ID: `selective-checks`

### Build CI images (`build-ci-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Depends on | `build-info` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `false`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `use-uv`: `true`
- `upgrade-to-newer-dependencies`: `false`
- `docker-cache`: `registry`
- `disable-airflow-repo-cache`: `false`

### Generate constraints (`generate-constraints`)

| Property | Value |
|----------|-------|
| Uses workflow | [Generate constraints](#generate-constraints) |
| Depends on | `build-info`, `build-ci-images` |

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `generate-pypi-constraints`: `true`
- `generate-no-providers-constraints`: `true`
- `debug-resources`: `false`
- `use-uv`: `true`

### Commit and push constraints (`update-constraints`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `generate-constraints` |

**Permissions:**

- `contents`: `write`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_VERSIONS` | `${{ needs.build-info.outputs.python-versions-list-as-string }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Set constraints branch name**
   - ID: `constraints-branch`

4. **Checkout ${{ steps.constraints-branch.outputs.branch }}**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `path`: `constraints`
     - `ref`: `${{ steps.constraints-branch.outputs.branch }}`
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

5. **Download constraints from the generate-constraints job**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `constraints-*`
     - `path`: `./files`

6. **Diff in constraints for Python: ${{ needs.build-info.outputs.python-versions-list-as-string }}**

7. **Commit changed constraint files for Python: ${{ needs.build-info.outputs.python-versions-list-as-string }}**

8. **Push changes**

### Notify on failure (`notify-on-failure`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `build-ci-images`, `generate-constraints`, `update-constraints` |
| Condition | `failure()` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |

#### Steps

1. **Send Slack notification**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨 Update constraints workflow failed on branch *${{ github.ref_name }}*\n\n*Details:* <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨 Update constraints workflow failed on *${{ github.ref_name }}*\n\n*Details:* <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

# Update constraints on push for main (only when uv.lock changes)

| Property | Value |
|----------|-------|
| File | `update-constraints-on-push.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`, `v[0-9]+-[0-9]+-test`
  - paths: `uv.lock`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_REPOSITORY` | `${{ github.repository }}` |
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_USERNAME` | `${{ github.actor }}` |
| `VERBOSE` | `true` |

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
update-constraints-on-push.yml [push]
+-- build-info / Install Breeze (uses ./.github/actions/breeze)
+-- build-ci-images (uses ci-image-build.yml)
|   +-- build-ci-images / Install Breeze (uses ./.github/actions/breeze)
+-- generate-constraints (uses generate-constraints.yml)
    +-- generate-constraints-matrix / Install prek (uses ./.github/actions/install-prek)
    +-- generate-constraints-matrix / Prepare breeze & CI image: ${{ matrix.python-version }} (uses ./.github/actions/prepare_breeze_and_image)
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | job `notify-on-failure` env `SLACK_BOT_TOKEN` |

## Jobs

### Build info (`build-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Fetch incoming commit ${{ github.sha }} with its parent**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ github.sha }}`
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

4. **Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

5. **Save github context to file**

6. **Selective checks**
   - ID: `selective-checks`

### Build CI images (`build-ci-images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build CI images](#build-ci-images) |
| Depends on | `build-info` |

**Permissions:**

- `contents`: `read`
- `packages`: `write`

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `push-image`: `false`
- `upload-image-artifact`: `true`
- `upload-mount-cache-artifact`: `false`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `branch`: `${{ needs.build-info.outputs.default-branch }}`
- `constraints-branch`: `${{ needs.build-info.outputs.default-constraints-branch }}`
- `use-uv`: `true`
- `upgrade-to-newer-dependencies`: `false`
- `docker-cache`: `registry`
- `disable-airflow-repo-cache`: `false`

### Generate constraints (`generate-constraints`)

| Property | Value |
|----------|-------|
| Uses workflow | [Generate constraints](#generate-constraints) |
| Depends on | `build-info`, `build-ci-images` |

#### Inputs forwarded

- `runners`: `["ubuntu-22.04"]`
- `platform`: `linux/amd64`
- `python-versions-list-as-string`: `${{ needs.build-info.outputs.python-versions-list-as-string }}`
- `python-versions`: `${{ needs.build-info.outputs.python-versions }}`
- `generate-pypi-constraints`: `true`
- `generate-no-providers-constraints`: `true`
- `debug-resources`: `false`
- `use-uv`: `true`

### Commit and push constraints (`update-constraints`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `generate-constraints` |

**Permissions:**

- `contents`: `write`
- `packages`: `read`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `PYTHON_VERSIONS` | `${{ needs.build-info.outputs.python-versions-list-as-string }}` |

#### Steps

1. **Cleanup repo**

2. **Checkout ${{ github.ref }} ( ${{ github.sha }} )**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

3. **Set constraints branch name**
   - ID: `constraints-branch`

4. **Checkout ${{ steps.constraints-branch.outputs.branch }}**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `path`: `constraints`
     - `ref`: `${{ steps.constraints-branch.outputs.branch }}`
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

5. **Download constraints from the generate-constraints job**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `constraints-*`
     - `path`: `./files`

6. **Diff in constraints for Python: ${{ needs.build-info.outputs.python-versions-list-as-string }}**

7. **Commit changed constraint files for Python: ${{ needs.build-info.outputs.python-versions-list-as-string }}**

8. **Push changes**

### Notify on failure (`notify-on-failure`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `build-info`, `build-ci-images`, `generate-constraints`, `update-constraints` |
| Condition | `failure()` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |

#### Steps

1. **Send Slack notification**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c` (v3.0.3)
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: "🚨 Update constraints workflow failed on branch *${{ github.ref_name }}*\n\n*Details:* <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>" blocks:   - type: "section"     text:       type: "mrkdwn"       text: "🚨 Update constraints workflow failed on *${{ github.ref_name }}*\n\n*Details:* <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View the failure log>"`

# Upgrade check

| Property | Value |
|----------|-------|
| File | `upgrade-check.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `target-branch` | string | Yes | - | Branch to upgrade (e.g. 'main' or 'v3-2-test') |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `SLACK_BOT_TOKEN` | Yes | Slack bot token for notifications |

## Permissions

- `contents`: `write`
- `pull-requests`: `write`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `SLACK_BOT_TOKEN` | `${{ secrets.SLACK_BOT_TOKEN }}` |
| `TARGET_BRANCH` | `${{ inputs.target-branch }}` |

## Called by

```
upgrade-check.yml
+-- scheduled-upgrade-check-main.yml (job: upgrade-main)  <- entry point
+-- scheduled-upgrade-check-v3-2-test.yml (job: upgrade-v3-2-test)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env `GITHUB_TOKEN` |
| `SLACK_BOT_TOKEN` | workflow env `SLACK_BOT_TOKEN` |

## Jobs

### [${{ inputs.target-branch }}] Upgrade checks and PR (`createupgrade-check`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **[${{ inputs.target-branch }}] Cleanup repo**

2. **[${{ inputs.target-branch }}] Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `ref`: `${{ inputs.target-branch }}`
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

3. **[${{ inputs.target-branch }}] Configure git credentials**

4. **[${{ inputs.target-branch }}] Install Breeze**
   - ID: `breeze`
   - Uses: `./.github/actions/breeze`

5. **[${{ inputs.target-branch }}] Install prek**
   - ID: `prek`
   - Uses: `./.github/actions/install-prek`
   - With:
     - `python-version`: `${{ steps.breeze.outputs.host-python-version }}` - Python version to use
     - `platform`: `linux/amd64` - Platform for the build - linux/amd64 or linux/arm64 (required)
     - `save-cache`: `false` - Whether to save prek cache (required)

6. **[${{ inputs.target-branch }}] Run breeze ci upgrade**

7. **[${{ inputs.target-branch }}] Find upgrade PR**
   - ID: `find-pr`

8. **[${{ inputs.target-branch }}] Notify Slack on success**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c`
   - Condition: `success() && steps.find-pr.outputs.pr-url != ''`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: >-   🔧 [${{ inputs.target-branch }}] Scheduled CI upgrade PR   ready: ${{ steps.find-pr.outputs.pr-url }} blocks:   - type: section     text:       type: mrkdwn       text: >-         🔧 *[${{ inputs.target-branch }}] Scheduled CI upgrade         PR ready*          A new CI upgrade PR has been created as a draft on the         `${{ inputs.target-branch }}` branch. Please:            1. *Undraft* the PR to trigger CI           2. *Review* the changes           3. *Merge* it once CI passes          <${{ steps.find-pr.outputs.pr-url }}|View PR>`

9. **[${{ inputs.target-branch }}] Notify Slack on failure**
   - Uses: `slackapi/slack-github-action@45a88b9581bfab2566dc881e2cd66d334e621e2c`
   - Condition: `failure()`
   - With:
     - `method`: `chat.postMessage`
     - `token`: `${{ env.SLACK_BOT_TOKEN }}`
     - `payload`: `channel: "internal-airflow-ci-cd" text: >-   ⚠️ [${{ inputs.target-branch }}] Scheduled CI upgrade FAILED.   See: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }} blocks:   - type: section     text:       type: mrkdwn       text: >-         ⚠️ *[${{ inputs.target-branch }}] Scheduled CI upgrade         FAILED*          The `breeze ci upgrade` job on the         `${{ inputs.target-branch }}` branch did not complete         successfully. Please investigate the failed run and         re-run the workflow if needed.          <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View failed run>`

# Setup Breeze

Sets up Python and Breeze

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python-version` | Python version to use | No | `3.10` |

## Outputs

| Name | Description |
|------|-------------|
| `host-python-version` | Python version used in host |

# Install prek

Installs prek and related packages

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python-version` | Python version to use | No | `3.10` |
| `save-cache` | Whether to save prek cache | Yes | - |
| `platform` | Platform for the build - linux/amd64 or linux/arm64 | Yes | - |

# Run migration tests

Runs migration tests

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python-version` | Python version to run the tests on | Yes | - |

# Post tests on failure

Run post tests actions on failure

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

# Post tests on success

Run post tests actions on success

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `codecov-token` | Codecov token | Yes | - |
| `python-version` | Python version | Yes | - |

# Prepare all CI images

Recreates current python CI images from artifacts for all python versions

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python-versions-list-as-string` | Stringified array of all Python versions to test - separated by spaces. | Yes | - |
| `docker-volume-location` | File system location where to move docker space to | No | `/mnt/var-lib-docker` |
| `platform` | Platform for the build - linux/amd64 or linux/arm64 | Yes | - |

# Prepare breeze && current image (CI or PROD)

Installs breeze and recreates current python image from artifact

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python` | Python version for image to prepare | Yes | - |
| `image-type` | Which image type to prepare (ci/prod) | No | `ci` |
| `platform` | Platform for the build - linux/amd64 or linux/arm64 | Yes | - |
| `use-uv` | Whether to use uv | Yes | - |
| `make-mnt-writeable-and-cleanup` | Whether to cleanup /mnt | Yes | - |

## Outputs

| Name | Description |
|------|-------------|
| `host-python-version` | Python version used in host |

# Prepare single CI image

Recreates current python image from artifacts (needed for the hard-coded actions calling all possible Python versions in "prepare_all_ci_images" action. Hopefully we can get rid of it when the https://github.com/apache/airflow/issues/45268 is resolved and we contribute capability of downloading multiple keys to the stash action.


| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `python` | Python version for image to prepare | Yes | - |
| `python-versions-list-as-string` | Stringified array of all Python versions to prepare - separated by spaces. | Yes | - |
| `platform` | Platform for the build - linux/amd64 or linux/arm64 | Yes | - |

