# transformers

43 workflows, 11 reusable workflows

## Contents

**Workflows**

- [Add model like runner - push](#add-model-like-runner)
- [Anti-Slop - pull_request_target](#anti-slop)
- [Assign PR Reviewers - pull_request_target](#assign-pr-reviewers)
- [Benchmark v2 Framework - workflow_dispatch](#benchmark-v2-framework)
- [Benchmark v2 Scheduled Runner - A10 Single-GPU - workflow_dispatch](#benchmark-v2-scheduled-runner---a10-single-gpu)
- [Benchmark v2 Scheduled Runner - MI325 Single-GPU - workflow_dispatch](#benchmark-v2-scheduled-runner---mi325-single-gpu)
- [Build docker images (Nightly CI) - workflow_call, push](#build-docker-images-nightly-ci)
- [Build docker images (Past CI) - push](#build-docker-images-past-ci)
- [Build docker images (scheduled) - push, repository_dispatch, workflow_dispatch, workflow_call, schedule](#build-docker-images-scheduled)
- [Build documentation - workflow_dispatch, push](#build-documentation)
- [Build pr ci-docker - push, repository_dispatch, workflow_call, schedule](#build-pr-ci-docker)
- [Build PR Documentation - pull_request, merge_group](#build-pr-documentation)
- [Check Permissions Advisor - workflow_dispatch](#check-permissions-advisor)
- [Check Tiny Models - push, repository_dispatch, schedule](#check-tiny-models)
- [CircleCI Failure Summary Comment - pull_request_target](#circleci-failure-summary-comment)
- [CodeQL Security Analysis - push, workflow_dispatch](#codeql-security-analysis)
- [Doctests - push, repository_dispatch, schedule](#doctests)
- [Extras Smoke Test - schedule](#extras-smoke-test)
- [New model PR merged notification - push](#new-model-pr-merged-notification)
- [Nvidia CI - repository_dispatch, schedule, push, workflow_dispatch](#nvidia-ci)
- [Nvidia CI - Flash Attn - repository_dispatch, schedule, push, workflow_dispatch](#nvidia-ci---flash-attn)
- [Nvidia CI with nightly torch - repository_dispatch, workflow_run, push](#nvidia-ci-with-nightly-torch)
- [PR - build doc via comment - issue_comment](#pr---build-doc-via-comment)
- [PR CI - pull_request](#pr-ci)
- [PR comment GitHub CI - issue_comment](#pr-comment-github-ci)
- [PR Repo. Consistency Bot - issue_comment](#pr-repo-consistency-bot)
- [PR slow CI - Suggestion - pull_request_target](#pr-slow-ci---suggestion)
- [Release - push](#release)
- [Release - Conda - push](#release---conda)
- [Secret Leaks - push](#secret-leaks)
- [Self-hosted runner (AMD mi250 scheduled CI caller) - workflow_run, push](#self-hosted-runner-amd-mi250-scheduled-ci-caller)
- [Self-hosted runner (AMD scheduled CI caller) - schedule](#self-hosted-runner-amd-scheduled-ci-caller)
- [Self-hosted runner (benchmark) - push, pull_request](#self-hosted-runner-benchmark)
- [Self-hosted runner (Intel Gaudi3 scheduled CI caller) - repository_dispatch, workflow_dispatch, schedule](#self-hosted-runner-intel-gaudi3-scheduled-ci-caller)
- [Self-hosted runner (nightly-past-ci-caller) - schedule, push](#self-hosted-runner-nightly-past-ci-caller)
- [Self-hosted runner scale set (AMD mi325 scheduled CI caller) - workflow_run, push](#self-hosted-runner-scale-set-amd-mi325-scheduled-ci-caller)
- [Self-hosted runner scale set (AMD mi355 scheduled CI caller) - workflow_run, push](#self-hosted-runner-scale-set-amd-mi355-scheduled-ci-caller)
- [Slow tests on important models (on Push - A10) - push](#slow-tests-on-important-models-on-push---a10)
- [SSH into our runners - workflow_dispatch](#ssh-into-our-runners)
- [Stale Bot - schedule](#stale-bot)
- [TRL CI bot - issue_comment](#trl-ci-bot)
- [Update Transformers metadata - push](#update-transformers-metadata)
- [Upload PR Documentation - workflow_run](#upload-pr-documentation)

**Reusable workflows**

- [CI collated reports](#ci-collated-reports)
- [CI slack report](#ci-slack-report)
- [Doctest job](#doctest-job)
- [Get PR commit SHA](#get-pr-commit-sha)
- [Get PR number](#get-pr-number)
- [model jobs (model_jobs.yml)](#model-jobs)
- [model jobs (model_jobs_intel_gaudi.yml)](#model-jobs-1)
- [Nvidia CI (job definitions)](#nvidia-ci-job-definitions)
- [Process failed tests](#process-failed-tests)
- [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci)
- [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi)

## Secrets and variables used across this repository

**Secrets:**

| Name | Used by |
|------|---------|
| `ACCESS_REPO_INFO_TOKEN` | [CI collated reports](#ci-collated-reports), [CI slack report](#ci-slack-report), [Doctests](#doctests), [Nvidia CI (job definitions)](#nvidia-ci-job-definitions), [Process failed tests](#process-failed-tests) |
| `ANACONDA_API_TOKEN` | [Release - Conda](#release---conda) |
| `CI_SLACK_BOT_TOKEN` | [CI slack report](#ci-slack-report), [Doctests](#doctests) |
| `CI_SLACK_CHANNEL_DOCKER` | [Build docker images (scheduled)](#build-docker-images-scheduled) |
| `CI_SLACK_CHANNEL_DUMMY_TESTS` | [CI slack report](#ci-slack-report) |
| `CI_SLACK_CHANNEL_ID` | [CI slack report](#ci-slack-report) |
| `CI_SLACK_CHANNEL_ID_DAILY` | [CI slack report](#ci-slack-report) |
| `CI_SLACK_CHANNEL_ID_DAILY_DOCS` | [Doctests](#doctests) |
| `COMMENT_BOT_TOKEN` | [Upload PR Documentation](#upload-pr-documentation) |
| `DOCKERHUB_PASSWORD` | [Build docker images (Nightly CI)](#build-docker-images-nightly-ci), [Build docker images (Past CI)](#build-docker-images-past-ci), [Build docker images (scheduled)](#build-docker-images-scheduled), [Build pr ci-docker](#build-pr-ci-docker) |
| `DOCKERHUB_USERNAME` | [Build docker images (Nightly CI)](#build-docker-images-nightly-ci), [Build docker images (Past CI)](#build-docker-images-past-ci), [Build docker images (scheduled)](#build-docker-images-scheduled), [Build pr ci-docker](#build-pr-ci-docker) |
| `HF_CI_WRITE_TOKEN` | [CircleCI Failure Summary Comment](#circleci-failure-summary-comment) |
| `HF_DOC_BUILD_PUSH` | [Build documentation](#build-documentation), [Upload PR Documentation](#upload-pr-documentation) |
| `HF_HUB_READ_TOKEN` | [Benchmark v2 Framework](#benchmark-v2-framework), [Check Tiny Models](#check-tiny-models), [model jobs](#model-jobs), [model jobs](#model-jobs-1), [Nvidia CI (job definitions)](#nvidia-ci-job-definitions), [PR comment GitHub CI](#pr-comment-github-ci), [Process failed tests](#process-failed-tests), [Self-hosted runner (benchmark)](#self-hosted-runner-benchmark), [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi), [SSH into our runners](#ssh-into-our-runners) |
| `HF_STYLE_BOT_ACTION` | [PR Repo. Consistency Bot](#pr-repo-consistency-bot) |
| `HUGGINGFACE_PUSH` | [Build documentation](#build-documentation) |
| `LYSANDRE_HF_TOKEN` | [Update Transformers metadata](#update-transformers-metadata) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | [PR CI](#pr-ci) |
| `OTEL_TOKEN` | [PR CI](#pr-ci) |
| `PUSH_TO_HUB_TOKEN` | [Self-hosted runner (benchmark)](#self-hosted-runner-benchmark) |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | [Build docker images (scheduled)](#build-docker-images-scheduled), [Build pr ci-docker](#build-pr-ci-docker), [Extras Smoke Test](#extras-smoke-test), [New model PR merged notification](#new-model-pr-merged-notification), [Process failed tests](#process-failed-tests), [SSH into our runners](#ssh-into-our-runners) |
| `SLACK_CIFEEDBACK_CHANNEL` | [SSH into our runners](#ssh-into-our-runners) |
| `TAILSCALE_SSH_AUTHKEY` | [SSH into our runners](#ssh-into-our-runners) |
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | [Benchmark v2 Framework](#benchmark-v2-framework), [CI collated reports](#ci-collated-reports), [CI slack report](#ci-slack-report), [Process failed tests](#process-failed-tests) |
| `TRANSFORMERS_HUB_BOT_HF_TOKEN` | [Check Tiny Models](#check-tiny-models) |
| `TRL_CI_DISPATCH_TOKEN` | [TRL CI bot](#trl-ci-bot) |

## Permissions across this repository

| Scope | Level |
|-------|-------|
| `actions` | `read` |
| `contents` | `write` (also granted as `read` elsewhere) |
| `id-token` | `write` (OIDC) |
| `issues` | `write` (also granted as `read` elsewhere) |
| `packages` | `read` |
| `pull-requests` | `write` (also granted as `read` elsewhere) |
| `security-events` | `write` |
| `statuses` | `write` |

# Add model like runner

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `add-model-like.yml` |

## Event filters

- **push**
  - branches: `none`

## Permissions

- `contents`: `read`

## Jobs

### Add new model like template tests (`run_tests_templates_like`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (10)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Install dependencies**

3. **Load cached virtual environment**
   - ID: `cache`
   - Uses: `actions/cache@v4.3.0`
   - With:
     - `path`: `~/venv/`
     - `key`: `v4-tests_model_like-${{ hashFiles('setup.py') }}`

4. **Create virtual environment on cache miss**
   - Condition: `steps.cache.outputs.cache-hit != 'true'`

5. **Check transformers location**

6. **Create model files**

7. **Run all PyTorch modeling test**

8. **Run style changes**

9. **Failure short reports**
   - Condition: `${{ always() }}`

10. **Test suite reports artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `run_all_tests_new_models_test_reports`
     - `path`: `reports/tests_new_models`

</details>

[Back to top](#contents)

# Anti-Slop

**Triggers:** `pull_request_target`

| Property | Value |
|----------|-------|
| File | `anti-slop.yml` |

## Event filters

- **pull_request_target**
  - types: `opened`, `reopened`

## Permissions

- `contents`: `read`
- `issues`: `read`
- `pull-requests`: `write`

## Jobs

### `anti-slop`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (1)</summary>

1. **peakoss/anti-slop@v0.2.1**
   - With:
     - `max-failures`: `2`
     - `close-pr`: `false`
     - `lock-pr`: `false`
     - `failure-add-pr-labels`: `Code agent slop`
     - `failure-pr-message`: `This PR was flagged by our automated quality checks. If you're a genuine contributor, please reply here and a maintainer will review your PR.  Common reasons for flagging: - New GitHub account - Unusually high number of repository forks in a 24-hour window  We appreciate your contribution and apologize if this is a false positive!`
     - `min-account-age`: `30`
     - `max-daily-forks`: `7`
     - `blocked-source-branches`: -
     - `blocked-paths`: -
     - `detect-spam-usernames`: `false`
     - `min-profile-completeness`: `0`
     - `require-description`: `false`
     - `require-linked-issue`: `false`
     - `require-conventional-title`: `false`
     - `require-pr-template`: `false`
     - `strict-pr-template-sections`: -
     - `optional-pr-template-sections`: -
     - `max-additional-pr-template-sections`: `0`
     - `max-description-length`: `0`
     - `require-conventional-commits`: `false`
     - `require-commit-author-match`: `false`
     - `require-maintainer-can-modify`: `false`
     - `require-final-newline`: `false`
     - `max-added-comments`: `0`
     - `max-emoji-count`: `0`
     - `max-code-references`: `0`
     - `max-commit-message-length`: `0`
     - `min-repo-merged-prs`: `0`
     - `min-repo-merge-ratio`: `0`
     - `min-global-merge-ratio`: `0`
     - `exempt-author-association`: `OWNER,MEMBER,COLLABORATOR`
     - `exempt-label`: `exempt`

</details>

[Back to top](#contents)

# Assign PR Reviewers

**Triggers:** `pull_request_target`

| Property | Value |
|----------|-------|
| File | `assign-reviewers.yml` |

## Event filters

- **pull_request_target**
  - branches: `main`
  - types: `ready_for_review`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `assign_reviewers` step `Run assignment script` env `GITHUB_TOKEN` |

## Jobs

### `assign_reviewers`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.13`

3. **Install dependencies**

4. **Run assignment script**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

</details>

[Back to top](#contents)

# Benchmark v2 Framework

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `benchmark_v2.yml` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runner` | string | Yes | - | Runner to use for the benchmark job |
| `container_image` | string | Yes | - | Container image to use |
| `container_options` | string | No | - | Container options |
| `commit_sha` | string | No | - | Commit SHA to benchmark |
| `run_id` | string | No | - | Run ID for tracking |
| `benchmark_repo_id` | string | Yes | - | Repository ID to push benchmark results |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |

## Called by

`benchmark_v2.yml`

- [benchmark_v2_a10_caller.yml](#benchmark-v2---default-models-benchmark-v2-default) (job: `benchmark-v2-default`) - entry point
- [benchmark_v2_mi325_caller.yml](#benchmark-v2---default-models-benchmark-v2-default-1) (job: `benchmark-v2-default`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN`; job `benchmark-v2` step `Run benchmark v2` env `HF_TOKEN` |
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | job `benchmark-v2` step `Run benchmark v2` env `UPLOAD_TOKEN` |

## Jobs

### Benchmark v2 (`benchmark-v2`)

| Property | Value |
|----------|-------|
| Runs on | `${{ inputs.runner }}` |
| Condition | `(github.event_name == 'pull_request' && contains( github.event.pull_request.labels.*.name, 'run-benchmark')) \|\|<br>(github.event_name == 'schedule')` |

<details>
<summary>Steps (5)</summary>

1. **Get repo**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `ref`: `${{ inputs.commit_sha || github.sha }}`
     - `persist-credentials`: `false`

2. **Install benchmark dependencies**

3. **Reinstall transformers in edit mode**

4. **Show installed libraries and their versions**

5. **Run benchmark v2**
   - Env:
     - `HF_TOKEN`: `${{ secrets.HF_HUB_READ_TOKEN }}`
     - `COMMIT_ID`: `${{ inputs.commit_sha || github.sha }}`
     - `RUN_ID`: `${{ inputs.run_id }}`
     - `BENCHMARK_REPO_ID`: `${{ inputs.benchmark_repo_id }}`
     - `UPLOAD_TOKEN`: `${{ secrets.TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN }}`

</details>

[Back to top](#contents)

# Benchmark v2 Scheduled Runner - A10 Single-GPU

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `benchmark_v2_a10_caller.yml` |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`benchmark_v2_a10_caller.yml` [workflow_dispatch]

- `benchmark-v2-default` uses [benchmark_v2.yml](#benchmark-v2-framework)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `HF_HUB_READ_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Benchmark v2 - Default Models (`benchmark-v2-default`)

| Property | Value |
|----------|-------|
| Uses workflow | [Benchmark v2 Framework](#benchmark-v2-framework) |

#### Inputs forwarded

- `runner`: `aws-g5-4xlarge-cache-use1-public-80`
- `container_image`: `huggingface/transformers-all-latest-gpu`
- `container_options`: `--gpus all --privileged --ipc host --shm-size "16gb"`
- `commit_sha`: `${{ github.sha }}`
- `run_id`: `${{ github.run_id }}`
- `benchmark_repo_id`: `hf-internal-testing/transformers-daily-benchmarks`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Benchmark v2 Scheduled Runner - MI325 Single-GPU

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `benchmark_v2_mi325_caller.yml` |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`benchmark_v2_mi325_caller.yml` [workflow_dispatch]

- `benchmark-v2-default` uses [benchmark_v2.yml](#benchmark-v2-framework)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `HF_HUB_READ_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Benchmark v2 - Default Models (`benchmark-v2-default`)

| Property | Value |
|----------|-------|
| Uses workflow | [Benchmark v2 Framework](#benchmark-v2-framework) |

#### Inputs forwarded

- `runner`: `amd-mi325-ci-1gpu`
- `container_image`: `huggingface/transformers-pytorch-amd-gpu`
- `container_options`: `--device /dev/kfd --device /dev/dri --env ROCR_VISIBLE_DEVICES --shm-size "16gb" --ipc host -v /mnt/cache/.cache/huggingface:/mnt/cache`
- `commit_sha`: `${{ github.sha }}`
- `run_id`: `${{ github.run_id }}`
- `benchmark_repo_id`: `hf-internal-testing/transformers-daily-benchmarks`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Build docker images (Nightly CI)

**Triggers:** `workflow_call`, `push`

| Property | Value |
|----------|-------|
| File | `build-nightly-ci-docker-images.yml` |

**Jobs:** [Nightly PyTorch](#nightly-pytorch-latest-with-torch-nightly-docker), [Nightly PyTorch + DeepSpeed](#nightly-pytorch--deepspeed-nightly-torch-deepspeed-docker)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job` | string | Yes | - | - |

## Event filters

- **push**
  - branches: `build_nightly_ci_docker_image*`

## Permissions

- `contents`: `read`

**Concurrency:** group `docker-images-builds`, cancel-in-progress: `false`

## Called by

`build-nightly-ci-docker-images.yml`

- [self-nightly-caller.yml](#build-ci-docker-images-with-nightly-torch-build_nightly_torch_ci_images) (job: `build_nightly_torch_ci_images`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCKERHUB_USERNAME` | job `latest-with-torch-nightly-docker` step `Login to DockerHub` with `username`; job `nightly-torch-deepspeed-docker` step `Login to DockerHub` with `username` |
| `DOCKERHUB_PASSWORD` | job `latest-with-torch-nightly-docker` step `Login to DockerHub` with `password`; job `nightly-torch-deepspeed-docker` step `Login to DockerHub` with `password` |

## Jobs

### Nightly PyTorch (`latest-with-torch-nightly-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |
| Condition | `inputs.job == 'latest-with-torch-nightly-docker' \|\| inputs.job == ''` |

<details>
<summary>Steps (4)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v2.10.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v2.2.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v3.3.1`
   - With:
     - `context`: `./docker/transformers-all-latest-gpu`
     - `build-args`: `REF=main PYTORCH=pre`
     - `push`: `true`
     - `tags`: `huggingface/transformers-all-latest-torch-nightly-gpu`

</details>

### Nightly PyTorch + DeepSpeed (`nightly-torch-deepspeed-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g4dn-2xlarge-cache` |
| Condition | `inputs.job == 'nightly-torch-deepspeed-docker' \|\| inputs.job == ''` |

<details>
<summary>Steps (4)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v2.10.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v2.2.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v3.3.1`
   - With:
     - `context`: `./docker/transformers-pytorch-deepspeed-nightly-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-deepspeed-nightly-gpu`

</details>

[Back to top](#contents)

# Build docker images (Past CI)

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `build-past-ci-docker-images.yml` |
| Default runs-on | `group: aws-general-8-plus` |

**Jobs:** [Past PyTorch Docker](#past-pytorch-docker-past-pytorch-docker), [Past TensorFlow Docker](#past-tensorflow-docker-past-tensorflow-docker)

## Event filters

- **push**
  - branches: `build_past_ci_docker_image*`

## Permissions

- `contents`: `read`

**Concurrency:** group `docker-images-builds`, cancel-in-progress: `false`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCKERHUB_USERNAME` | job `past-pytorch-docker` step `Login to DockerHub` with `username`; job `past-tensorflow-docker` step `Login to DockerHub` with `username` |
| `DOCKERHUB_PASSWORD` | job `past-pytorch-docker` step `Login to DockerHub` with `password`; job `past-tensorflow-docker` step `Login to DockerHub` with `password` |

## Jobs

### Past PyTorch Docker (`past-pytorch-docker`)

| Property | Value |
|----------|-------|
| Matrix | `version`: 1.13, 1.12, 1.11 |

<details>
<summary>Steps (6)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v2.10.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Get Base Image**
   - ID: `get-base-image`
   - Env:
     - `framework_version`: `${{ matrix.version }}`

4. **Print Base Image**

5. **Login to DockerHub**
   - Uses: `docker/login-action@v2.2.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

6. **Build and push**
   - Uses: `docker/build-push-action@v3.3.1`
   - With:
     - `context`: `./docker/transformers-past-gpu`
     - `build-args`: `REF=main BASE_DOCKER_IMAGE=${{ steps.get-base-image.outputs.base_image }} FRAMEWORK=pytorch VERSION=${{ matrix.version }}`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-past-${{ matrix.version }}-gpu`

</details>

### Past TensorFlow Docker (`past-tensorflow-docker`)

| Property | Value |
|----------|-------|
| Matrix | `version`: 2.11, 2.10, 2.9, 2.8, 2.7, 2.6, 2.5 |

<details>
<summary>Steps (6)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v2.10.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Get Base Image**
   - ID: `get-base-image`
   - Env:
     - `framework_version`: `${{ matrix.version }}`

4. **Print Base Image**

5. **Login to DockerHub**
   - Uses: `docker/login-action@v2.2.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

6. **Build and push**
   - Uses: `docker/build-push-action@v3.3.1`
   - With:
     - `context`: `./docker/transformers-past-gpu`
     - `build-args`: `REF=main BASE_DOCKER_IMAGE=${{ steps.get-base-image.outputs.base_image }} FRAMEWORK=tensorflow VERSION=${{ matrix.version }}`
     - `push`: `true`
     - `tags`: `huggingface/transformers-tensorflow-past-${{ matrix.version }}-gpu`

</details>

[Back to top](#contents)

# Build docker images (scheduled)

**Triggers:** `push`, `repository_dispatch`, `workflow_dispatch`, `workflow_call`, `schedule`

| Property | Value |
|----------|-------|
| File | `build-docker-images.yml` |
| Default runs-on | `group: aws-general-8-plus` |

**Jobs:** [Latest PyTorch \[dev\]](#latest-pytorch-dev-latest-docker), [PyTorch with Flash Attn \[dev\]](#pytorch-with-flash-attn-dev-flash-attn-ci-image), [Latest PyTorch + DeepSpeed](#latest-pytorch--deepspeed-latest-torch-deepspeed-docker), [Doc builder](#doc-builder-doc-builder), [Latest PyTorch (AMD) \[dev\]](#latest-pytorch-amd-dev-latest-pytorch-amd), [Cache Latest Pytorch (AMD) Image](#cache-latest-pytorch-amd-image-cache-latest-pytorch-amd), [PyTorch + DeepSpeed (AMD) \[dev\]](#pytorch--deepspeed-amd-dev-latest-pytorch-deepspeed-amd), [Latest Pytorch + Quantization \[dev\]](#latest-pytorch--quantization-dev-latest-quantization-torch-docker)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `image_postfix` | string | Yes | - | - |

## Schedule

- `17 0 * * *`

## Event filters

- **push**
  - branches: `build_ci_docker_image*`

## Permissions

- `contents`: `read`

**Concurrency:** group `docker-images-builds`, cancel-in-progress: `false`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCKERHUB_USERNAME` | job `latest-docker` step `Login to DockerHub` with `username`; job `flash-attn-ci-image` step `Login to DockerHub` with `username`; job `latest-torch-deepspeed-docker` step `Login to DockerHub` with `username`; job `doc-builder` step `Login to DockerHub` with `username`; job `latest-pytorch-amd` step `Login to DockerHub` with `username`; job `cache-latest-pytorch-amd` step `Login to DockerHub` with `username`; job `latest-pytorch-deepspeed-amd` step `Login to DockerHub` with `username`; job `latest-quantization-torch-docker` step `Login to DockerHub` with `username` |
| `DOCKERHUB_PASSWORD` | job `latest-docker` step `Login to DockerHub` with `password`; job `flash-attn-ci-image` step `Login to DockerHub` with `password`; job `latest-torch-deepspeed-docker` step `Login to DockerHub` with `password`; job `doc-builder` step `Login to DockerHub` with `password`; job `latest-pytorch-amd` step `Login to DockerHub` with `password`; job `cache-latest-pytorch-amd` step `Login to DockerHub` with `password`; job `latest-pytorch-deepspeed-amd` step `Login to DockerHub` with `password`; job `latest-quantization-torch-docker` step `Login to DockerHub` with `password` |
| `CI_SLACK_CHANNEL_DOCKER` | job `latest-docker` step `Post to Slack` with `slack_channel`; job `flash-attn-ci-image` step `Post to Slack` with `slack_channel`; job `latest-torch-deepspeed-docker` step `Post to Slack` with `slack_channel`; job `doc-builder` step `Post to Slack` with `slack_channel`; job `latest-pytorch-amd` step `Post to Slack` with `slack_channel`; job `latest-pytorch-deepspeed-amd` step `Post to Slack` with `slack_channel`; job `latest-quantization-torch-docker` step `Post to Slack` with `slack_channel` |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `latest-docker` step `Post to Slack` with `slack_token`; job `flash-attn-ci-image` step `Post to Slack` with `slack_token`; job `latest-torch-deepspeed-docker` step `Post to Slack` with `slack_token`; job `doc-builder` step `Post to Slack` with `slack_token`; job `latest-pytorch-amd` step `Post to Slack` with `slack_token`; job `latest-pytorch-deepspeed-amd` step `Post to Slack` with `slack_token`; job `latest-quantization-torch-docker` step `Post to Slack` with `slack_token` |

## Jobs

### Latest PyTorch [dev] (`latest-docker`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-all-latest-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-all-latest-gpu${{ inputs.image_postfix }}`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the transformers-all-latest-gpu docker build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### PyTorch with Flash Attn [dev] (`flash-attn-ci-image`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-all-latest-gpu`
     - `build-args`: `REF=main PYTORCH=2.8.0 TORCHCODEC=0.7.0 FLASH_ATTN=yes`
     - `push`: `true`
     - `tags`: `huggingface/transformers-all-latest-gpu${{ inputs.image_postfix }}:flash-attn`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the transformers-all-latest-gpu docker build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Latest PyTorch + DeepSpeed (`latest-torch-deepspeed-docker`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-pytorch-deepspeed-latest-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-deepspeed-latest-gpu${{ inputs.image_postfix }}`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER}}`
     - `title`: `🤗 Results of the transformers-pytorch-deepspeed-latest-gpu docker build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Doc builder (`doc-builder`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-doc-builder`
     - `push`: `true`
     - `tags`: `huggingface/transformers-doc-builder`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the huggingface/transformers-doc-builder docker build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Latest PyTorch (AMD) [dev] (`latest-pytorch-amd`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-highcpu-32-priv` |

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-pytorch-amd-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-amd-gpu${{ inputs.image_postfix }}`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the huggingface/transformers-pytorch-amd-gpu build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Cache Latest Pytorch (AMD) Image (`cache-latest-pytorch-amd`)

| Property | Value |
|----------|-------|
| Runs on | `group: amd-mi325-1gpu` |
| Depends on | `latest-pytorch-amd` |

<details>
<summary>Steps (2)</summary>

1. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

2. **Pull and save docker image to cache**

</details>

### PyTorch + DeepSpeed (AMD) [dev] (`latest-pytorch-deepspeed-amd`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-pytorch-deepspeed-amd-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-deepspeed-amd-gpu${{ inputs.image_postfix }}`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the transformers-pytorch-deepspeed-amd-gpu build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Latest Pytorch + Quantization [dev] (`latest-quantization-torch-docker`)

<details>
<summary>Steps (5)</summary>

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

2. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker/transformers-quantization-latest-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-quantization-latest-gpu${{ inputs.image_postfix }}`

5. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@63657f571a92cc9759159442936061c51d6d9ae4`
   - Condition: `always()`
   - With:
     - `slack_channel`: `${{ secrets.CI_SLACK_CHANNEL_DOCKER }}`
     - `title`: `🤗 Results of the transformers-quantization-latest-gpu build`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

[Back to top](#contents)

# Build documentation

**Triggers:** `workflow_dispatch`, `push`

| Property | Value |
|----------|-------|
| File | `build_documentation.yml` |

**Jobs:** [`build`](#build), [`build_other_lang`](#build_other_lang)

## Event filters

- **push**
  - branches: `main`, `doc-builder*`, `v*-release`, `use_templates`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`build_documentation.yml` [workflow_dispatch, push]

- uses **`huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391`** (x2)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `HF_DOC_BUILD_PUSH`, `HUGGINGFACE_PUSH`, `hf_token`, `token`

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HUGGINGFACE_PUSH` | job `build` secrets `token`; job `build_other_lang` secrets `token` |
| `HF_DOC_BUILD_PUSH` | job `build` secrets `hf_token`; job `build_other_lang` secrets `hf_token` |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391` (external) |

#### Inputs forwarded

- `commit_sha`: `${{ github.sha }}`
- `package`: `transformers`
- `notebook_folder`: `transformers_doc`
- `languages`: `en`
- `custom_container`: `huggingface/transformers-doc-builder`

#### Secrets forwarded

- `token`: `${{ secrets.HUGGINGFACE_PUSH }}`
- `hf_token`: `${{ secrets.HF_DOC_BUILD_PUSH }}`

### `build_other_lang`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391` (external) |

#### Inputs forwarded

- `commit_sha`: `${{ github.sha }}`
- `package`: `transformers`
- `notebook_folder`: `transformers_doc`
- `languages`: `ar de es fr hi it ja ko pt tr zh`
- `custom_container`: `huggingface/transformers-doc-builder`

#### Secrets forwarded

- `token`: `${{ secrets.HUGGINGFACE_PUSH }}`
- `hf_token`: `${{ secrets.HF_DOC_BUILD_PUSH }}`

[Back to top](#contents)

# Build pr ci-docker

**Triggers:** `push`, `repository_dispatch`, `workflow_call`, `schedule`

| Property | Value |
|----------|-------|
| File | `build-ci-docker-images.yml` |
| Default runs-on | `ubuntu-22.04` |

**Jobs:** [`build`](#build-1), [`notify`](#notify)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `image_postfix` | string | Yes | - | - |

## Schedule

- `6 0 * * *`

## Event filters

- **push**
  - branches: `push-ci-image`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}`, cancel-in-progress: `true`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCKERHUB_USERNAME` | job `build` step `Login to DockerHub` with `username` |
| `DOCKERHUB_PASSWORD` | job `build` step `Login to DockerHub` with `password` |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `notify` step `Post to Slack` with `slack_token` |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Matrix | `file`: quality, consistency, custom-tokenizers, torch-light, exotic-models, examples-torch |
| Condition | `${{ contains(github.event.head_commit.message, '[build-ci-image]') \|\| contains(github.event.head_commit.message, '[push-ci-image]') && '!cancelled()' \|\| github.event_name == 'schedule' }}` |

<details>
<summary>Steps (5)</summary>

1. **Set tag**
   - Env:
     - `COMMIT_MESSAGE`: `${{ github.event.head_commit.message }}`

2. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@v3.12.0`

3. **Check out code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

4. **Login to DockerHub**
   - Uses: `docker/login-action@v3.7.0`
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

5. **Build ${{ matrix.file }}.dockerfile**
   - Uses: `docker/build-push-action@v5.4.0`
   - With:
     - `context`: `./docker`
     - `build-args`: `REF=${{ github.sha }}`
     - `file`: `./docker/${{ matrix.file }}.dockerfile`
     - `push`: `${{ contains(github.event.head_commit.message, 'ci-image]') ||  github.event_name == 'schedule' }}`
     - `tags`: `${{ env.TAG }}`

</details>

### `notify`

| Property | Value |
|----------|-------|
| Condition | `${{ contains(github.event.head_commit.message, '[build-ci-image]') \|\| contains(github.event.head_commit.message, '[push-ci-image]') && '!cancelled()' \|\| github.event_name == 'schedule' }}` |

<details>
<summary>Steps (1)</summary>

1. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@a88e7fa2eaee28de5a4d6142381b1fb792349b67`
   - Condition: `${{ contains(github.event.head_commit.message, '[push-ci-image]') && github.event_name != 'schedule' }}`
   - With:
     - `slack_channel`: `#transformers-ci-circleci-images`
     - `title`: `🤗 New docker images for CircleCI are pushed.`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

[Back to top](#contents)

# Build PR Documentation

**Triggers:** `pull_request`, `merge_group`

| Property | Value |
|----------|-------|
| File | `build_pr_documentation.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`build`](#build-2), [`skip_merge_queue`](#skip_merge_queue), [`doc_build_status_check`](#doc_build_status_check)

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.head_ref || github.run_id }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

`build_pr_documentation.yml` [pull_request, merge_group]

- `build` uses `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@90b4ee2c10b81b5c1a6367c4e6fc9e2fb510a7e3`

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@90b4ee2c10b81b5c1a6367c4e6fc9e2fb510a7e3`

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@90b4ee2c10b81b5c1a6367c4e6fc9e2fb510a7e3` (external) |
| Condition | `github.event_name == 'pull_request'` |

#### Inputs forwarded

- `commit_sha`: `${{ github.event.pull_request.head.sha }}`
- `pr_number`: `${{ github.event.number }}`
- `package`: `transformers`
- `languages`: `en`

### `skip_merge_queue`

| Property | Value |
|----------|-------|
| Condition | `github.event_name == 'merge_group'` |

<details>
<summary>Steps (1)</summary>

1. **echo "Skipping doc build in merge queue"**

</details>

### `doc_build_status_check`

| Property | Value |
|----------|-------|
| Depends on | `build`, `skip_merge_queue` |
| Condition | `always()` |

<details>
<summary>Steps (1)</summary>

1. **if [[ "${{ needs.build.result }}" == "success" || "${{ ne...**

</details>

[Back to top](#contents)

# Check Permissions Advisor

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `check-workflow-permissions.yml` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow_name` | string | No | - | Workflow file name |
| `run_count` | string | No | `10` | Number of runs to analyze |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`check-workflow-permissions.yml` [workflow_dispatch]

- `advisor` uses `huggingface/security-workflows/.github/workflows/permissions-advisor-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c`

## Transitive requirements (from full call graph)

Permissions declared across the chain: `actions: read`, `contents: read`

External workflows referenced: `huggingface/security-workflows/.github/workflows/permissions-advisor-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c`

## Jobs

### `advisor`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/security-workflows/.github/workflows/permissions-advisor-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c` (external) |

**Permissions:**

- `actions`: `read`
- `contents`: `read`

#### Inputs forwarded

- `workflow_name`: `${{ inputs.workflow_name }}`
- `run_count`: `${{ fromJSON(inputs.run_count) }}`

[Back to top](#contents)

# Check Tiny Models

**Triggers:** `push`, `repository_dispatch`, `schedule`

| Property | Value |
|----------|-------|
| File | `check_tiny_models.yml` |

## Schedule

- `0 2 * * *`

## Event filters

- **push**
  - branches: `check_tiny_models*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `TOKEN` | `${{ secrets.TRANSFORMERS_HUB_BOT_HF_TOKEN }}` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `TRANSFORMERS_HUB_BOT_HF_TOKEN` | workflow env `TOKEN` |
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |

## Jobs

### Check tiny models (`check_tiny_models`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (6)</summary>

1. **Checkout transformers**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

2. **actions/checkout@v4.3.1**
   - With:
     - `persist-credentials`: `false`

3. **Set up Python 3.10**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.10`
     - `architecture`: `x64`

4. **Install**

5. **Create all tiny models (locally)**

6. **Local tiny model reports artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `tiny_local_model_creation_reports`
     - `path`: `tiny_local_models/reports`

</details>

[Back to top](#contents)

# CircleCI Failure Summary Comment

**Triggers:** `pull_request_target`

| Property | Value |
|----------|-------|
| File | `circleci-failure-summary-comment.yml` |

## Event filters

- **pull_request_target**
  - types: `opened`, `synchronize`, `reopened`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `comment` step `Wait for CircleCI check suite completion` env `GH_TOKEN`; job `comment` step `Get CircleCI run's artifacts and upload them to Hub` env `GITHUB_TOKEN`; job `comment` step `Post comment with helper link` env `GH_TOKEN` |
| `HF_CI_WRITE_TOKEN` | job `comment` step `Upload summaries to Hub` env `HF_TOKEN` |

## Jobs

### `comment`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (8)</summary>

1. **Checkout repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Setup Python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.13`

3. **Install dependencies**

4. **Wait for CircleCI check suite completion**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `COMMIT_SHA`: `${{ github.event.pull_request.head.sha }}`
     - `GITHUB_REPOSITORY`: `${{ github.repository }}`

5. **Get CircleCI run's artifacts and upload them to Hub**
   - ID: `circleci`
   - Env:
     - `COMMIT_SHA`: `${{ github.event.pull_request.head.sha }}`
     - `REPO`: `${{ github.repository }}`
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

6. **Upload summaries to Hub**
   - Condition: `steps.circleci.outputs.artifact_found == 'true'`
   - Env:
     - `HF_TOKEN`: `${{ secrets.HF_CI_WRITE_TOKEN }}`
     - `CIRCLECI_RESULTS_DATASET_ID`: `transformers-community/circleci-test-results`
     - `PR_NUMBER`: `${{ github.event.pull_request.number }}`
     - `COMMIT_SHA`: `${{ github.event.pull_request.head.sha }}`

7. **Delete existing CircleCI summary comments**
   - Uses: `actions/github-script@v7.1.0`
   - Condition: `steps.circleci.outputs.artifact_found == 'true'`
   - With:
     - `script`: `` const PR_NUMBER = parseInt(process.env.PR_NUMBER, 10);  // Get all comments on the PR const { data: comments } = await github.rest.issues.listComments({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: PR_NUMBER });  // Find existing bot comments that start with "View the CircleCI Test Summary for this PR:" const existingComments = comments.filter(comment =>    comment.user.login === 'github-actions[bot]' &&    comment.body.startsWith('View the CircleCI Test Summary for this PR:') );  // Delete all matching comments for (const comment of existingComments) {   console.log(`Deleting comment #${comment.id}`);   await github.rest.issues.deleteComment({     owner: context.repo.owner,     repo: context.repo.repo,     comment_id: comment.id   }); }  console.log(`Deleted ${existingComments.length} old CircleCI summary comment(s)`); ``
   - Env:
     - `PR_NUMBER`: `${{ github.event.pull_request.number }}`

8. **Post comment with helper link**
   - Condition: `steps.circleci.outputs.artifact_found == 'true'`
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_REPOSITORY`: `${{ github.repository }}`
     - `PR_NUMBER`: `${{ github.event.pull_request.number }}`
     - `PR_SHA`: `${{ github.event.pull_request.head.sha }}`

</details>

[Back to top](#contents)

# CodeQL Security Analysis

**Triggers:** `push`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `codeql.yml` |

## Event filters

- **push**
  - branches: `main`, `fix_security_issue_*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`codeql.yml` [push, workflow_dispatch]

- `codeql` uses `huggingface/security-workflows/.github/workflows/codeql-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c`

## Transitive requirements (from full call graph)

Permissions declared across the chain: `actions: read`, `contents: read`, `packages: read`, `security-events: write`

External workflows referenced: `huggingface/security-workflows/.github/workflows/codeql-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c`

## Jobs

### CodeQL Analysis (`codeql`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/security-workflows/.github/workflows/codeql-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c` (external) |

**Permissions:**

- `security-events`: `write`
- `packages`: `read`
- `actions`: `read`
- `contents`: `read`

#### Inputs forwarded

- `languages`: `["actions"]`
- `queries`: `security-extended,security-and-quality`
- `runner`: `ubuntu-latest`

[Back to top](#contents)

# Doctests

**Triggers:** `push`, `repository_dispatch`, `schedule`

| Property | Value |
|----------|-------|
| File | `doctests.yml` |

**Jobs:** [Setup](#setup-setup), [Call doctest jobs](#call-doctest-jobs-call_doctest_job), [Send results to webhook](#send-results-to-webhook-send_results)

## Schedule

- `17 2 * * *`

## Event filters

- **push**
  - branches: `run_doctest*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `NUM_SLICES` | `3` |

## Call graph (rooted at this workflow)

`doctests.yml` [push, repository_dispatch, schedule]

- `call_doctest_job` uses [doctest_job.yml](#doctest-job)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_ID_DAILY_DOCS`

Permissions declared across the chain: `contents: read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CI_SLACK_BOT_TOKEN` | job `send_results` step `Send message to Slack` env `CI_SLACK_BOT_TOKEN` |
| `ACCESS_REPO_INFO_TOKEN` | job `send_results` step `Send message to Slack` env `ACCESS_REPO_INFO_TOKEN` |
| `CI_SLACK_CHANNEL_ID_DAILY_DOCS` | job `send_results` step `Send message to Slack` env `SLACK_REPORT_CHANNEL` |

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g5-4xlarge-cache` |

<details>
<summary>Steps (5)</summary>

1. **Update clone**

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **Show installed libraries and their versions**

4. **Check values for matrix**

5. **Set values for matrix**
   - ID: `set-matrix`

</details>

### Call doctest jobs (`call_doctest_job`)

| Property | Value |
|----------|-------|
| Uses workflow | [Doctest job](#doctest-job) |
| Matrix | `split_keys`: ${{ fromJson(needs.setup.outputs.split_keys) }} |
| Depends on | `setup` |

#### Inputs forwarded

- `job_splits`: `${{ needs.setup.outputs.job_splits }}`
- `split_keys`: `${{ toJson(matrix.split_keys) }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Send results to webhook (`send_results`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `call_doctest_job` |
| Condition | `always()` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `persist-credentials`: `false`

2. **actions/download-artifact@v4.3.0**

3. **Send message to Slack**
   - Env:
     - `CI_SLACK_BOT_TOKEN`: `${{ secrets.CI_SLACK_BOT_TOKEN }}`
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `SLACK_REPORT_CHANNEL`: `${{ secrets.CI_SLACK_CHANNEL_ID_DAILY_DOCS }}`

4. **Upload results**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `doc_test_results`
     - `path`: `doc_test_results`

</details>

[Back to top](#contents)

# Extras Smoke Test

**Triggers:** `schedule`

| Property | Value |
|----------|-------|
| File | `extras-smoke-test.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [Get supported Python versions](#get-supported-python-versions-get-python-versions), [Test extras on Python (python-version)](#test-extras-on-python-python-version-test-extras), [Check Slack token availability](#check-slack-token-availability-precheck-slack), [Notify failures to Slack](#notify-failures-to-slack-notify-failures)

## Schedule

- `0 3 * * *`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `SLACK_CHANNEL_ID` | `#transformers-gh-ci-central` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `precheck-slack` step `chk` env `SLACK_BOT_TOKEN`; job `notify-failures` step `Send Slack notification` env `SLACK_BOT_TOKEN` |

## Jobs

### Get supported Python versions (`get-python-versions`)

<details>
<summary>Steps (3)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Install setuptools**

3. **Extract Python versions from setup.py**
   - ID: `extract-versions`

</details>

### Test extras on Python (python-version) (`test-extras`)

| Property | Value |
|----------|-------|
| Matrix | `python-version`: ${{ fromJson(needs.get-python-versions.outputs.versions) }} |
| Depends on | `get-python-versions` |

<details>
<summary>Steps (8)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Set up Python ${{ matrix.python-version }}**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `${{ matrix.python-version }}`
     - `allow-prereleases`: `true`

3. **Install base dependencies**

4. **Extract extras for this Python version**
   - ID: `get-extras`
   - Env:
     - `MATRIX_PYTHON_VERSION`: `${{ matrix.python-version }}`

5. **Install base package**

6. **Test all extras**
   - ID: `test-extras`
   - Env:
     - `MATRIX_PYTHON_VERSION`: `${{ matrix.python-version }}`

7. **Verify installation**

8. **Upload failure report**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `always()`
   - With:
     - `name`: `failure-report-${{ matrix.python-version }}`
     - `path`: `failure_reports/`
     - `retention-days`: `1`
     - `if-no-files-found`: `ignore`

</details>

### Check Slack token availability (`precheck-slack`)

<details>
<summary>Steps (1)</summary>

1. **chk**
   - ID: `chk`
   - Env:
     - `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

### Notify failures to Slack (`notify-failures`)

| Property | Value |
|----------|-------|
| Depends on | `test-extras`, `precheck-slack` |
| Condition | `always() && needs.precheck-slack.outputs.has_slack_token == 'true' && needs.test-extras.result != 'success'` |

<details>
<summary>Steps (6)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.11`

3. **Download all failure reports** `[continue-on-error]`
   - Uses: `actions/download-artifact@v4.3.0`
   - With:
     - `pattern`: `failure-report-*`
     - `path`: `failure_reports/`
     - `merge-multiple`: `true`

4. **Aggregate failures**

5. **Format Slack message**
   - Env:
     - `FAILURES_FILE`: `all_failures.json`
     - `WORKFLOW_URL`: `${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}`

6. **Send Slack notification**
   - Uses: `slackapi/slack-github-action@6c661ce58804a1a20f6dc5fbee7f0381b469e001`
   - Condition: `env.SLACK_MESSAGE != ''`
   - With:
     - `channel-id`: `${{ env.SLACK_CHANNEL_ID }}`
     - `payload`: `{   "blocks": [     {       "type": "header",       "text": {         "type": "plain_text",         "text": "${{ env.SLACK_TITLE }}"       }     },     {       "type": "section",       "text": {         "type": "mrkdwn",         "text": "${{ env.SLACK_MESSAGE }}"       }     },     {       "type": "divider"     },     {       "type": "section",       "text": {         "type": "mrkdwn",         "text": "<${{ env.SLACK_WORKFLOW_URL }}|View workflow run>"       }     }   ] }`
   - Env:
     - `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

[Back to top](#contents)

# New model PR merged notification

**Triggers:** `push`

Used to notify core maintainers about new model PR being merged

| Property | Value |
|----------|-------|
| File | `new_model_pr_merged_notification.yml` |

## Event filters

- **push**
  - branches: `main`
  - paths: `src/transformers/models/*/modeling_*`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `notify_new_model` step `Notify` env `SLACK_BOT_TOKEN` |

## Jobs

### Notify new model (`notify_new_model`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Check new model**

3. **print commit sha**
   - Condition: `${{ env.NEW_MODEL != ''}}`

4. **print new model**
   - Condition: `${{ env.NEW_MODEL != ''}}`

5. **Notify**
   - Uses: `slackapi/slack-github-action@6c661ce58804a1a20f6dc5fbee7f0381b469e001`
   - Condition: `${{ env.NEW_MODEL != ''}}`
   - With:
     - `channel-id`: `transformers-new-model-notification`
     - `payload`: `{   "blocks": [     {       "type": "header",       "text": {         "type": "plain_text",         "text": "New model!",         "emoji": true       }     },     {       "type": "section",       "text": {         "type": "mrkdwn",         "text": "<https://github.com/huggingface/transformers/commit/${{ env.COMMIT_SHA }}|New model: ${{ env.NEW_MODEL }}> GH_ArthurZucker, GH_lysandrejik, GH_ydshieh\ncommit SHA: ${{ env.COMMIT_SHA }}"       }     }   ] }`
   - Env:
     - `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

[Back to top](#contents)

# Nvidia CI

**Triggers:** `repository_dispatch`, `schedule`, `push`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `self-scheduled-caller.yml` |

**Jobs:** [Setup](#setup-setup-1), [Model CI](#model-ci-model-ci), [Torch pipeline CI](#torch-pipeline-ci-torch-pipeline), [Example CI](#example-ci-example-ci), [Trainer/FSDP CI](#trainerfsdp-ci-trainer-fsdp-ci), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci), [Quantization CI](#quantization-ci-quantization-ci), [Kernels CI](#kernels-ci-kernels-ci)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `prev_workflow_run_id` | string | No | - | previous workflow run id to compare |
| `other_workflow_run_id` | string | No | - | other workflow run id to compare |

## Schedule

- `17 2 * * *`

## Event filters

- **push**
  - branches: `run_nvidia_ci*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `prev_workflow_run_id` | - |
| `other_workflow_run_id` | - |

## Call graph (rooted at this workflow)

`self-scheduled-caller.yml` [repository_dispatch, schedule, push, workflow_dispatch]

- uses **[self-scheduled.yml](#nvidia-ci-job-definitions)** (x7)
  - uses **[model_jobs.yml](#model-jobs)** (x2)
    - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
  - `send_results` uses [slack-report.yml](#ci-slack-report)
  - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (2)</summary>

1. **Setup**
   - Env:
     - `prev_workflow_run_id`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `other_workflow_run_id`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

</details>

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-daily-models`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `ci_event`: `Daily CI`
- `runner_type`: `a10`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Torch pipeline CI (`torch-pipeline`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_pipelines_torch_gpu`
- `slack_report_channel`: `#transformers-ci-daily-pipeline-torch`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `ci_event`: `Daily CI`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Example CI (`example-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_examples_gpu`
- `slack_report_channel`: `#transformers-ci-daily-examples`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `ci_event`: `Daily CI`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Trainer/FSDP CI (`trainer-fsdp-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_trainer_and_fsdp_gpu`
- `slack_report_channel`: `#transformers-ci-daily-training`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `runner_type`: `a10`
- `ci_event`: `Daily CI`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `slack_report_channel`: `#transformers-ci-daily-training`
- `docker`: `huggingface/transformers-pytorch-deepspeed-latest-gpu`
- `ci_event`: `Daily CI`
- `working-directory-prefix`: `/workspace`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Quantization CI (`quantization-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_quantization_torch_gpu`
- `slack_report_channel`: `#transformers-ci-daily-quantization`
- `docker`: `huggingface/transformers-quantization-latest-gpu`
- `ci_event`: `Daily CI`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Kernels CI (`kernels-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_kernels_gpu`
- `slack_report_channel`: `#transformers-ci-daily-kernels`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `ci_event`: `Daily CI`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci`
- `commit_sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Nvidia CI - Flash Attn

**Triggers:** `repository_dispatch`, `schedule`, `push`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `self-scheduled-flash-attn-caller.yml` |

**Jobs:** [Setup](#setup-setup-2), [Model CI](#model-ci-model-ci-1)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `prev_workflow_run_id` | string | No | - | previous workflow run id to compare |
| `other_workflow_run_id` | string | No | - | other workflow run id to compare |

## Schedule

- `17 2 * * *`

## Event filters

- **push**
  - branches: `run_nvidia_ci_flash_attn*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `prev_workflow_run_id` | - |
| `other_workflow_run_id` | - |

## Call graph (rooted at this workflow)

`self-scheduled-flash-attn-caller.yml` [repository_dispatch, schedule, push, workflow_dispatch]

- `model-ci` uses [self-scheduled.yml](#nvidia-ci-job-definitions)
  - uses **[model_jobs.yml](#model-jobs)** (x2)
    - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
  - `send_results` uses [slack-report.yml](#ci-slack-report)
  - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (2)</summary>

1. **Setup**
   - Env:
     - `PREV_WORKFLOW_RUN_ID`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `OTHER_WORKFLOW_RUN_ID`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

</details>

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-flash-attn`
- `docker`: `huggingface/transformers-all-latest-gpu:flash-attn`
- `ci_event`: `Daily CI`
- `runner_type`: `a10`
- `report_repo_id`: `hf-internal-testing/transformers_flash_attn_ci`
- `commit_sha`: `${{ github.sha }}`
- `pytest_marker`: `flash_attn_test or flash_attn_3_test or flash_attn_4_test or all_flash_attn_test`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Nvidia CI with nightly torch

**Triggers:** `repository_dispatch`, `workflow_run`, `push`

| Property | Value |
|----------|-------|
| File | `self-nightly-caller.yml` |

**Jobs:** [Build CI Docker Images with nightly torch](#build-ci-docker-images-with-nightly-torch-build_nightly_torch_ci_images), [Setup](#setup-setup-3), [Model CI](#model-ci-model-ci-2)

## Event filters

- **workflow_run**
  - workflows: `Nvidia CI`
  - branches: `main`
  - types: `completed`
- **push**
  - branches: `run_ci_with_nightly_torch*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `prev_workflow_run_id` | - |
| `other_workflow_run_id` | - |

## Call graph (rooted at this workflow)

`self-nightly-caller.yml` [repository_dispatch, workflow_run, push]

- `build_nightly_torch_ci_images` uses [build-nightly-ci-docker-images.yml](#build-docker-images-nightly-ci)
- `model-ci` uses [self-scheduled.yml](#nvidia-ci-job-definitions)
  - uses **[model_jobs.yml](#model-jobs)** (x2)
    - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
  - `send_results` uses [slack-report.yml](#ci-slack-report)
  - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `DOCKERHUB_PASSWORD`, `DOCKERHUB_USERNAME`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Build CI Docker Images with nightly torch (`build_nightly_torch_ci_images`)

| Property | Value |
|----------|-------|
| Uses workflow | [Build docker images (Nightly CI)](#build-docker-images-nightly-ci) |

#### Inputs forwarded

- `job`: `latest-with-torch-nightly-docker`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (2)</summary>

1. **Setup**
   - Env:
     - `PREV_WORKFLOW_RUN_ID`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `OTHER_WORKFLOW_RUN_ID`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

</details>

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |
| Depends on | `build_nightly_torch_ci_images` |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-past-future`
- `docker`: `huggingface/transformers-all-latest-torch-nightly-gpu`
- `ci_event`: `Nightly CI`
- `runner_type`: `a10`
- `report_repo_id`: `hf-internal-testing/transformers_daily_ci_with_torch_nightly`
- `commit_sha`: `${{ github.event.workflow_run.head_sha || github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# PR - build doc via comment

**Triggers:** `issue_comment`

| Property | Value |
|----------|-------|
| File | `pr_build_doc_with_comment.yml` |
| Default runs-on | `ubuntu-22.04` |

**Jobs:** [Get PR number](#get-pr-number-get-pr-number), [Get PR commit SHA](#get-pr-commit-sha-get-pr-info), [Verity PR commit corresponds to a specific event by comparing timestamps](#verity-pr-commit-corresponds-to-a-specific-event-by-comparing-timestamps-verity_pr_commit), [Create run](#create-run-create_run), [Reply to the comment](#reply-to-the-comment-reply_to_comment), [Build doc](#build-doc-build-doc), [Update Check Run Status](#update-check-run-status-update_run_status)

## Event filters

- **issue_comment**
  - types: `created`
  - branches-ignore: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `${{ github.workflow }}-${{ github.event.issue.number }}-${{ startsWith(github.event.comment.body, 'build-doc') }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

`pr_build_doc_with_comment.yml` [issue_comment]

- `get-pr-number` uses [get-pr-number.yml](#get-pr-number)
- `get-pr-info` uses [get-pr-info.yml](#get-pr-commit-sha)
- `build-doc` uses `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@093eb65f2e8745457987df060dc392e6bcf1347a`

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `GITHUB_TOKEN`

Permissions declared across the chain: `contents: read`, `pull-requests: write`, `statuses: write`

External workflows referenced: `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@093eb65f2e8745457987df060dc392e6bcf1347a`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `create_run` step `Create Run` env `GH_TOKEN`; job `reply_to_comment` step `Reply to the comment` env `GH_TOKEN`; job `update_run_status` env `GH_TOKEN` |

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR number](#get-pr-number) |
| Condition | `${{ github.event.issue.state == 'open' && contains(fromJSON('["ydshieh", "ArthurZucker", "zucchini-nlp", "molbap", "gante", "LysandreJik", "Cyrilvallez", "Rocketknight1", "SunMarc", "eustlb", "MekkCyber", "vasqu", "ivarflakstad", "stevhliu", "ebezzam", "itazap", "tarekziade"]'), github.actor) && (startsWith(github.event.comment.body, 'build-doc')) }}` |

### Get PR commit SHA (`get-pr-info`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR commit SHA](#get-pr-commit-sha) |
| Depends on | `get-pr-number` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

#### Inputs forwarded

- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Verity PR commit corresponds to a specific event by comparing timestamps (`verity_pr_commit`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-info` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COMMENT_DATE` | `${{ github.event.comment.created_at }}` |
| `PR_MERGE_COMMIT_DATE` | `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_DATE }}` |
| `PR_MERGE_COMMIT_TIMESTAMP` | `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}` |

<details>
<summary>Steps (1)</summary>

1. **COMMENT\_TIMESTAMP=$(date -d "${COMMENT\_DATE}" +"%s")**

</details>

### Create run (`create_run`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-pr-info` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != '' }}` |

**Permissions:**

- `statuses`: `write`

<details>
<summary>Steps (1)</summary>

1. **Create Run**
   - ID: `create_run`
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `NEEDS_GET_PR_INFO_OUTPUTS_PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

</details>

### Reply to the comment (`reply_to_comment`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `create_run` |
| Condition | `${{ needs.create_run.result == 'success' }}` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (1)</summary>

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `NEEDS_GET_PR_NUMBER_OUTPUTS_PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

</details>

### Build doc (`build-doc`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@093eb65f2e8745457987df060dc392e6bcf1347a` (external) |
| Depends on | `get-pr-number`, `get-pr-info` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != '' }}` |

#### Inputs forwarded

- `commit_sha`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`
- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
- `package`: `transformers`
- `languages`: `ar de en es fr hi it ja ko pt zh`

### Update Check Run Status (`update_run_status`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-info`, `create_run`, `build-doc` |
| Condition | `${{ always() && needs.create_run.result == 'success' }}` |

**Permissions:**

- `statuses`: `write`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GH_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |
| `GITHUB_RUN_URL` | `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}` |
| `STATUS_OK` | `${{ contains(fromJSON('["skipped", "success"]'), needs.build-doc.result) }}` |

<details>
<summary>Steps (2)</summary>

1. **Get \`build-doc\` job status**

2. **Update PR commit statuses**
   - Env:
     - `NEEDS_GET_PR_INFO_OUTPUTS_PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

</details>

[Back to top](#contents)

# PR CI

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `pr-ci-caller.yml` |

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.event.pull_request.number }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

`pr-ci-caller.yml` [pull_request]

- `pr-ci` uses `huggingface/transformers-test-ci/.github/workflows/pr-ci_dynamic_caller_example.yml@91d590c4f744e4564a8ae0d3810068c8a35b939e`

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_TOKEN`

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/transformers-test-ci/.github/workflows/pr-ci_dynamic_caller_example.yml@91d590c4f744e4564a8ae0d3810068c8a35b939e`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | job `pr-ci` secrets `OTEL_EXPORTER_OTLP_ENDPOINT` |
| `OTEL_TOKEN` | job `pr-ci` secrets `OTEL_TOKEN` |

## Jobs

### `pr-ci`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/transformers-test-ci/.github/workflows/pr-ci_dynamic_caller_example.yml@91d590c4f744e4564a8ae0d3810068c8a35b939e` (external) |
| Condition | `contains(fromJSON('["MEMBER","OWNER","COLLABORATOR"]'), github.event.pull_request.author_association) \|\| github.event.pull_request.user.login == 'ydshieh2'` |

#### Secrets forwarded

- `OTEL_EXPORTER_OTLP_ENDPOINT`: `${{ secrets.OTEL_EXPORTER_OTLP_ENDPOINT }}`
- `OTEL_TOKEN`: `${{ secrets.OTEL_TOKEN }}`

[Back to top](#contents)

# PR comment GitHub CI

**Triggers:** `issue_comment`

| Property | Value |
|----------|-------|
| File | `self-comment-ci.yml` |
| Default runs-on | `ubuntu-22.04` |

**Jobs:**

- [Get PR number](#get-pr-number-get-pr-number-1)
- [Get PR commit SHA](#get-pr-commit-sha-get-pr-info-1)
- [Check timestamps (security check)](#check-timestamps-security-check-check-timestamps)
- [`get-tests`](#get-tests)
- [Report error earlier](#report-error-earlier-report_error_earlier)
- [Reply to the comment](#reply-to-the-comment-reply_to_comment-1)
- [Create run](#create-run-create_run-1)
- [Model CI](#model-ci-model-ci-3)
- [Quantization CI](#quantization-ci-quantization-ci-1)
- [Check & Report](#check--report-report)

## Event filters

- **issue_comment**
  - types: `created`
  - branches-ignore: `main`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `OMP_NUM_THREADS` | `8` |
| `MKL_NUM_THREADS` | `8` |
| `RUN_SLOW` | `yes` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |
| `CUDA_VISIBLE_DEVICES` | `0,1` |

**Concurrency:** group `${{ github.workflow }}-${{ github.event.issue.number }}-${{ startsWith(github.event.comment.body, 'run-slow') || startsWith(github.event.comment.body, 'run slow') || startsWith(github.event.comment.body, 'run_slow') }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

`self-comment-ci.yml` [issue_comment]

- `get-pr-number` uses [get-pr-number.yml](#get-pr-number)
- `get-pr-info` uses [get-pr-info.yml](#get-pr-commit-sha)
- uses **[self-scheduled.yml](#nvidia-ci-job-definitions)** (x2)
  - uses **[model_jobs.yml](#model-jobs)** (x2)
    - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
  - `send_results` uses [slack-report.yml](#ci-slack-report)
  - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`, `pull-requests: write`, `statuses: write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `GITHUB_TOKEN` | job `report_error_earlier` step `Reply to the comment` env `GH_TOKEN`; job `reply_to_comment` step `Reply to the comment` env `GH_TOKEN`; job `create_run` step `Create Run` env `GH_TOKEN`; job `report` step `Post results as PR comment` env `GH_TOKEN`; job `report` step `Update PR commit statuses` env `GH_TOKEN` |

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR number](#get-pr-number) |
| Condition | `${{ github.event.issue.state == 'open' && contains(fromJSON('["ydshieh", "ArthurZucker", "zucchini-nlp", "molbap", "LysandreJik", "Cyrilvallez", "Rocketknight1", "SunMarc", "eustlb", "vasqu", "ivarflakstad", "stevhliu", "ebezzam", "remi-or", "itazap", "3outeille", "IlyasMoutawwakil", "tarekziade", "yonigozlan", "guarin"]'), github.actor) && (startsWith(github.event.comment.body, 'run-slow') \|\| startsWith(github.event.comment.body, 'run slow') \|\| startsWith(github.event.comment.body, 'run_slow')) }}` |

### Get PR commit SHA (`get-pr-info`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR commit SHA](#get-pr-commit-sha) |
| Depends on | `get-pr-number` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

#### Inputs forwarded

- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Check timestamps (security check) (`check-timestamps`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-info` |

<details>
<summary>Steps (1)</summary>

1. **Verify \`merge\_commit\` timestamp is older than the issue comment timestamp**
   - Env:
     - `COMMENT_DATE`: `${{ github.event.comment.created_at }}`
     - `PR_MERGE_COMMIT_TIMESTAMP`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}`

</details>

### `get-tests`

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `check-timestamps` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `fetch-depth`: `0`
     - `ref`: `refs/pull/${{ needs.get-pr-number.outputs.PR_NUMBER }}/merge`
     - `persist-credentials`: `false`

2. **Verify merge commit SHA**
   - Env:
     - `VERIFIED_PR_MERGE_SHA`: `${{ needs.check-timestamps.outputs.PR_MERGE_SHA }}`

3. **Get models to test**
   - Env:
     - `PR_COMMENT`: `${{ github.event.comment.body }}`

4. **Show models to test**
   - ID: `models_to_run`

</details>

### Report error earlier (`report_error_earlier`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-pr-info`, `get-tests` |
| Condition | `${{ always() && needs.get-pr-info.result == 'success' && needs.get-tests.result != 'success' }}` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (1)</summary>

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `PREFIX`: `[Workflow Run ⚙️](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})\n\n`
     - `github_repository`: `${{ github.repository }}`
     - `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

</details>

### Reply to the comment (`reply_to_comment`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-tests` |
| Condition | `${{ needs.get-tests.outputs.models != '[]'  \|\| needs.get-tests.outputs.quantizations != '[]' }}` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (1)</summary>

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `PREFIX`: `[Workflow Run ⚙️](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})`
     - `INFO`: `` \n\nThis comment contains `run-slow`, running the specified jobs ``
     - `BODY`: `\n\nmodels: ${{ needs.get-tests.outputs.models }}\nquantizations: ${{ needs.get-tests.outputs.quantizations }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

</details>

### Create run (`create_run`)

| Property | Value |
|----------|-------|
| Depends on | `check-timestamps`, `reply_to_comment` |

**Permissions:**

- `statuses`: `write`

<details>
<summary>Steps (1)</summary>

1. **Create Run**
   - ID: `create_run`
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_head_sha`: `${{ needs.check-timestamps.outputs.PR_HEAD_SHA }}`

</details>

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |
| Depends on | `get-pr-number`, `check-timestamps`, `get-tests`, `create_run` |
| Condition | `${{ needs.get-tests.outputs.models != '[]' }}` |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-pr`
- `docker`: `huggingface/transformers-all-latest-gpu`
- `ci_event`: `PR Comment CI`
- `report_repo_id`: `hf-internal-testing/transformers_pr_ci`
- `commit_sha`: `${{ needs.check-timestamps.outputs.PR_MERGE_SHA }}`
- `subdirs`: `${{ needs.get-tests.outputs.models }}`
- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Quantization CI (`quantization-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |
| Depends on | `get-pr-number`, `check-timestamps`, `get-tests`, `create_run` |
| Condition | `${{ needs.get-tests.outputs.quantizations != '[]' }}` |

#### Inputs forwarded

- `job`: `run_quantization_torch_gpu`
- `slack_report_channel`: `#transformers-ci-pr`
- `docker`: `huggingface/transformers-quantization-latest-gpu`
- `ci_event`: `PR Comment CI`
- `report_repo_id`: `hf-internal-testing/transformers_pr_ci`
- `commit_sha`: `${{ needs.check-timestamps.outputs.PR_MERGE_SHA }}`
- `subdirs`: `${{ needs.get-tests.outputs.quantizations }}`
- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Check & Report (`report`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-pr-info`, `check-timestamps`, `create_run`, `model-ci`, `quantization-ci` |
| Condition | `${{ always() && needs.create_run.result == 'success' }}` |

**Permissions:**

- `pull-requests`: `write`
- `statuses`: `write`

<details>
<summary>Steps (6)</summary>

1. **actions/download-artifact@v4.3.0**
   - With:
     - `pattern`: `new_failures_with_bad_commit_{run_models_gpu,run_quantization_torch_gpu}`
     - `path`: `./new_failures`
     - `merge-multiple`: `false`

2. **List downloaded artifacts**

3. **Show reports from jobs**

4. **Process and filter reports**

5. **Post results as PR comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
     - `pr_head_repo`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_FULL_NAME }}`
     - `run_commit`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_SHA }}`
     - `pr_commit`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`
     - `main_commit`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_BASE_SHA }}`
     - `model_ci_result`: `${{ needs.model-ci.result }}`
     - `quantization_ci_result`: `${{ needs.quantization-ci.result }}`
     - `model_infrastructure_ok`: `${{ needs.model-ci.outputs.is_infrastructure_ok }}`
     - `quant_infrastructure_ok`: `${{ needs.quantization-ci.outputs.is_infrastructure_ok }}`

6. **Update PR commit statuses**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_head_sha`: `${{ needs.check-timestamps.outputs.PR_HEAD_SHA }}`

</details>

[Back to top](#contents)

# PR Repo. Consistency Bot

**Triggers:** `issue_comment`

| Property | Value |
|----------|-------|
| File | `pr-repo-consistency-bot.yml` |
| Default runs-on | `ubuntu-22.04` |

**Jobs:** [Get PR number](#get-pr-number-get-pr-number-2), [Get PR commit SHA](#get-pr-commit-sha-get-pr-info-2), [Check timestamps (security check)](#check-timestamps-security-check-check-timestamps-1), [Init Comment on PR](#init-comment-on-pr-init_comment_with_url), [`run-repo-consistency-checks`](#run-repo-consistency-checks), [`commit-and-comment`](#commit-and-comment)

## Event filters

- **issue_comment**
  - types: `created`
  - branches-ignore: `main`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.event.issue.number }}-${{ startsWith(github.event.comment.body, '@bot /repo') || startsWith(github.event.comment.body, '@bot /style') }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

`pr-repo-consistency-bot.yml` [issue_comment]

- `get-pr-number` uses [get-pr-number.yml](#get-pr-number)
- `get-pr-info` uses [get-pr-info.yml](#get-pr-commit-sha)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `HF_STYLE_BOT_ACTION`

Permissions declared across the chain: `contents: read`, `contents: write`, `pull-requests: write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_STYLE_BOT_ACTION` | job `commit-and-comment` step `Push changes to fork using git` env `GITHUB_TOKEN` |

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR number](#get-pr-number) |
| Condition | `${{ github.event.issue.state == 'open' && contains(fromJSON('["ydshieh", "ArthurZucker", "zucchini-nlp", "molbap", "gante", "LysandreJik", "Cyrilvallez", "Rocketknight1", "SunMarc", "eustlb", "MekkCyber", "vasqu", "ivarflakstad", "stevhliu", "ebezzam", "remi-or", "itazap", "3outeille", "IlyasMoutawwakil", "tarekziade"]'), github.actor) && (startsWith(github.event.comment.body, '@bot /repo') \|\| startsWith(github.event.comment.body, '@bot /style')) }}` |

### Get PR commit SHA (`get-pr-info`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR commit SHA](#get-pr-commit-sha) |
| Depends on | `get-pr-number` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

#### Inputs forwarded

- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Check timestamps (security check) (`check-timestamps`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-info` |

<details>
<summary>Steps (1)</summary>

1. **Verify \`merge\_commit\` timestamp is older than the issue comment timestamp**
   - Env:
     - `COMMENT_DATE`: `${{ github.event.comment.created_at }}`
     - `PR_MERGE_COMMIT_TIMESTAMP`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}`

</details>

### Init Comment on PR (`init_comment_with_url`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `check-timestamps` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (2)</summary>

1. **Delete existing bot comment if it exists**
   - Uses: `actions/github-script@v6.4.1`
   - With:
     - `script`: `` const PR_NUMBER = parseInt(process.env.PR_NUMBER, 10);  // Get all comments on the PR const { data: comments } = await github.rest.issues.listComments({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: PR_NUMBER });  // Find existing bot comments that start with "Repo. Consistency" or "Style fix" const existingComments = comments.filter(comment =>    comment.user.login === 'github-actions[bot]' &&    (comment.body.startsWith('Repo. Consistency') || comment.body.startsWith('Style fix')) );  if (existingComments.length > 0) {   // Get the most recent comment   const mostRecentComment = existingComments     .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))[0];      console.log(`Deleting most recent comment #${mostRecentComment.id}`);   await github.rest.issues.deleteComment({     owner: context.repo.owner,     repo: context.repo.repo,     comment_id: mostRecentComment.id   }); } ``
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

2. **Comment on PR with workflow run link**
   - ID: `init_comment`
   - Uses: `actions/github-script@v6.4.1`
   - With:
     - `script`: `` const PR_NUMBER = parseInt(process.env.PR_NUMBER, 10); const COMMENT_BODY = process.env.COMMENT_BODY; const runUrl = `${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`  // Determine which command was used const isStyleFix = COMMENT_BODY.startsWith('@bot /style'); const messagePrefix = isStyleFix ? 'Style fix' : 'Repo. Consistency fix';  const { data: botComment } = await github.rest.issues.createComment({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: PR_NUMBER,   body: `${messagePrefix} is beginning .... [View the workflow run here](${runUrl}).` }); core.setOutput('comment_id', botComment.id); ``
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
     - `COMMENT_BODY`: `${{ github.event.comment.body }}`

</details>

### `run-repo-consistency-checks`

| Property | Value |
|----------|-------|
| Depends on | `get-pr-info`, `check-timestamps`, `init_comment_with_url` |

<details>
<summary>Steps (10)</summary>

1. **Checkout base repository**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `ref`: `main`
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@v4.9.1`
   - With:
     - `python-version`: `3.10`

3. **Install dependencies from trusted main branch**

4. **Fetch and checkout PR code manually**
   - Env:
     - `PR_HEAD_REPO_FULL_NAME`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_FULL_NAME }}`
     - `PR_HEAD_REF`: `${{ needs.get-pr-info.outputs.PR_HEAD_REF }}`
     - `PR_HEAD_SHA`: `${{ needs.check-timestamps.outputs.VERIFIED_PR_HEAD_SHA }}`

5. **Copy trusted scripts from main branch**

6. **Install editable transformers from PR branch with copied scripts**

7. **Run repo consistency checks with trusted script**
   - ID: `run_repo_checks`
   - Condition: `startsWith(github.event.comment.body, '@bot /repo')`

8. **Run style checks with trusted script**
   - ID: `run_style_checks`
   - Condition: `startsWith(github.event.comment.body, '@bot /style')`

9. **Save modified files**
   - Condition: `steps.run_repo_checks.outputs.changes_detected == 'true' || steps.run_style_checks.outputs.changes_detected == 'true'`

10. **Upload modified files**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `steps.run_repo_checks.outputs.changes_detected == 'true' || steps.run_style_checks.outputs.changes_detected == 'true'`
   - With:
     - `name`: `modified-files`
     - `path`: `artifact-staging/`

</details>

### `commit-and-comment`

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-pr-info`, `check-timestamps`, `init_comment_with_url`, `run-repo-consistency-checks` |
| Condition | `always()` |

**Permissions:**

- `pull-requests`: `write`
- `contents`: `write`

<details>
<summary>Steps (4)</summary>

1. **Download modified files**
   - Uses: `actions/download-artifact@v4.3.0`
   - Condition: `needs.run-repo-consistency-checks.outputs.changes_detected == 'true'`
   - With:
     - `name`: `modified-files`

2. **Push changes to fork using git**
   - Condition: `needs.run-repo-consistency-checks.outputs.changes_detected == 'true'`
   - Env:
     - `PR_HEAD_REF`: `${{ needs.get-pr-info.outputs.PR_HEAD_REF }}`
     - `PR_HEAD_SHA`: `${{ needs.check-timestamps.outputs.VERIFIED_PR_HEAD_SHA }}`
     - `PR_HEAD_REPO_FULL_NAME`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_FULL_NAME }}`
     - `GITHUB_TOKEN`: `${{ secrets.HF_STYLE_BOT_ACTION }}`

3. **Prepare final comment message**
   - ID: `prepare_final_comment`
   - Condition: `needs.init_comment_with_url.result == 'success'`
   - Env:
     - `CHANGES_DETECTED`: `${{ needs.run-repo-consistency-checks.outputs.changes_detected }}`
     - `COMMENT_BODY`: `${{ github.event.comment.body }}`

4. **Comment on PR**
   - Uses: `actions/github-script@v6.4.1`
   - Condition: `needs.init_comment_with_url.result == 'success'`
   - With:
     - `script`: `const pr_number = parseInt(process.env.PR_NUMBER, 10); const comment_id = parseInt(process.env.COMMENT_ID, 10); const body = process.env.FINAL_COMMENT; await github.rest.issues.updateComment({   owner: context.repo.owner,   repo: context.repo.repo,   comment_id,   body, });`
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
     - `COMMENT_ID`: `${{ needs.init_comment_with_url.outputs.comment_id }}`
     - `FINAL_COMMENT`: `${{ steps.prepare_final_comment.outputs.final_comment }}`

</details>

[Back to top](#contents)

# PR slow CI - Suggestion

**Triggers:** `pull_request_target`

| Property | Value |
|----------|-------|
| File | `pr_slow_ci_suggestion.yml` |
| Default runs-on | `ubuntu-22.04` |

**Jobs:** [Get PR number](#get-pr-number-get-pr-number-3), [Get PR commit SHA](#get-pr-commit-sha-get-pr-info-3), [Get test files to run](#get-test-files-to-run-get-jobs), [Send a comment to suggest jobs to run](#send-a-comment-to-suggest-jobs-to-run-send_comment)

## Event filters

- **pull_request_target**
  - types: `opened`, `synchronize`, `reopened`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`pr_slow_ci_suggestion.yml` [pull_request_target]

- `get-pr-number` uses [get-pr-number.yml](#get-pr-number)
- `get-pr-info` uses [get-pr-info.yml](#get-pr-commit-sha)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`, `pull-requests: write`

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR number](#get-pr-number) |

### Get PR commit SHA (`get-pr-info`)

| Property | Value |
|----------|-------|
| Uses workflow | [Get PR commit SHA](#get-pr-commit-sha) |
| Depends on | `get-pr-number` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

#### Inputs forwarded

- `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Get test files to run (`get-jobs`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-pr-info` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Write pr\_files file**
   - Uses: `actions/github-script@v6.4.1`
   - With:
     - `script`: `const fs = require('node:fs'); const files = await github.paginate(github.rest.pulls.listFiles, {   owner: context.repo.owner,   repo: context.repo.repo,   pull_number: parseInt(process.env.PR_NUMBER, 10), }); fs.writeFileSync('pr_files.txt', JSON.stringify(files));`
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

3. **Get repository content**
   - ID: `repo_content`
   - Uses: `actions/github-script@v6.4.1`
   - With:
     - `script`: `const fs = require('node:fs'); const { PR_HEAD_REPO_OWNER, PR_HEAD_REPO_NAME, PR_HEAD_SHA } = process.env;  const { data: tests_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests',   ref: PR_HEAD_SHA, });  const { data: tests_models_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests/models',   ref: PR_HEAD_SHA, });  const { data: tests_quantization_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests/quantization',   ref: PR_HEAD_SHA, });  // Write to files instead of outputs fs.writeFileSync('tests_dir.txt', JSON.stringify(tests_dir, null, 2)); fs.writeFileSync('tests_models_dir.txt', JSON.stringify(tests_models_dir, null, 2)); fs.writeFileSync('tests_quantization_dir.txt', JSON.stringify(tests_quantization_dir, null, 2));`
   - Env:
     - `PR_HEAD_REPO_OWNER`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_OWNER }}`
     - `PR_HEAD_REPO_NAME`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_NAME }}`
     - `PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

4. **Run script to get jobs to run**
   - ID: `get_jobs`

</details>

### Send a comment to suggest jobs to run (`send_comment`)

| Property | Value |
|----------|-------|
| Depends on | `get-pr-number`, `get-jobs` |
| Condition | `${{ needs.get-jobs.outputs.jobs != '' }}` |

**Permissions:**

- `pull-requests`: `write`

<details>
<summary>Steps (1)</summary>

1. **Check and update comment if needed**
   - Uses: `actions/github-script@v7.1.0`
   - With:
     - `script`: `` const prNumber = parseInt(process.env.PR_NUMBER, 10); const commentPrefix = "**[For maintainers]** Suggested jobs to run (before merge)"; const thirtyMinutesAgo = new Date(Date.now() - 30 * 60 * 1000); // 30 minutes ago const newBody = `${commentPrefix}${process.env.BODY}`;  // Get all comments on the PR const { data: comments } = await github.rest.issues.listComments({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: prNumber });  // Find existing comments that start with our prefix const existingComments = comments.filter(comment =>   comment.user.login === 'github-actions[bot]' &&   comment.body.startsWith(commentPrefix) );  let shouldCreateNewComment = true; let commentsToDelete = [];  if (existingComments.length > 0) {   // Get the most recent comment   const mostRecentComment = existingComments     .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))[0];    const commentDate = new Date(mostRecentComment.created_at);   const isOld = commentDate < thirtyMinutesAgo;   const isDifferentContent = mostRecentComment.body !== newBody;    console.log(`Most recent comment created: ${mostRecentComment.created_at}`);   console.log(`Is older than 30 minutes: ${isOld}`);   console.log(`Has different content: ${isDifferentContent}`);    if (isOld || isDifferentContent) {     // Delete all existing comments and create new one     commentsToDelete = existingComments;     console.log(`Will delete ${commentsToDelete.length} existing comment(s) and create new one`);   } else {     // Content is same and comment is recent, skip     shouldCreateNewComment = false;     console.log('Comment is recent and content unchanged, skipping update');   } } else {   console.log('No existing comments found, will create new one'); }  // Delete old comments if needed for (const comment of commentsToDelete) {   console.log(`Deleting comment #${comment.id} (created: ${comment.created_at})`);   await github.rest.issues.deleteComment({     owner: context.repo.owner,     repo: context.repo.repo,     comment_id: comment.id   }); }  // Create new comment if needed if (shouldCreateNewComment) {   await github.rest.issues.createComment({     owner: context.repo.owner,     repo: context.repo.repo,     issue_number: prNumber,     body: newBody   });   console.log('✅ New comment created'); } else {   console.log('ℹ️ No comment update needed'); } ``
   - Env:
     - `BODY`: `run-slow: ${{ needs.get-jobs.outputs.jobs }}`
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

</details>

[Back to top](#contents)

# Release

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [build release](#build-release-build_and_test), [`upload_package`](#upload_package)

## Event filters

- **push**
  - tags: `v*`
  - branches: `v*-release`

## Permissions

- `contents`: `read`

## Jobs

### build release (`build_and_test`)

<details>
<summary>Steps (13)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **set up python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.13`

3. **pip install setuptools**

4. **pip install -e .**

5. **make build-release**

6. **pip uninstall -y transformers**

7. **pip install dist/\*.whl**

8. **python -c "from transformers import \*"**

9. **pip install -e .[torch]**

10. **python -c "from transformers import pipeline; classifier ...**

11. **pip install twine**

12. **twine check --strict dist/\***

13. **Upload build artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `python-dist`
     - `path`: `dist/** build/**`

</details>

### `upload_package`

| Property | Value |
|----------|-------|
| Depends on | `build_and_test` |
| Condition | `startsWith(github.ref, 'refs/tags/')` |

**Deploys to environment:** `pypi-release` [gated]

**Permissions:**

- `id-token`: `write` (OIDC)

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Download build artifacts**
   - Uses: `actions/download-artifact@v4.3.0`
   - With:
     - `name`: `python-dist`
     - `path`: `.`

3. **Publish package distributions to TestPyPI**
   - Uses: `pypa/gh-action-pypi-publish@ed0c53931b1dc9bd32cbe73a98c7f6766f8a527e`
   - With:
     - `verbose`: `true`

</details>

[Back to top](#contents)

# Release - Conda

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `release-conda.yml` |

## Event filters

- **push**
  - tags: `v*`
  - branches: `conda_*`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `ANACONDA_API_TOKEN` | `${{ secrets.ANACONDA_API_TOKEN }}` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `ANACONDA_API_TOKEN` | workflow env `ANACONDA_API_TOKEN` |

## Jobs

### `build_and_package`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

**Defaults:** shell `bash -l {0}`

<details>
<summary>Steps (6)</summary>

1. **Checkout repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Install miniconda**
   - Uses: `conda-incubator/setup-miniconda@v2.3.0`
   - With:
     - `auto-update-conda`: `true`
     - `auto-activate-base`: `false`
     - `python-version`: `3.8`
     - `activate-environment`: `build-transformers`
     - `channels`: `huggingface`

3. **Setup conda env**

4. **Extract version**

5. **Build conda packages**

6. **Upload to Anaconda**

</details>

[Back to top](#contents)

# Secret Leaks

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `trufflehog.yml` |

## Permissions

- `contents`: `read`

## Jobs

### `trufflehog`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (2)</summary>

1. **Checkout code**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Secret Scanning**
   - Uses: `trufflesecurity/trufflehog@6bd2d14f7a4bc1e569fa3550efa7ec632a4fa67b`
   - With:
     - `extra_args`: `--results=verified,unknown`

</details>

[Back to top](#contents)

# Self-hosted runner (AMD mi250 scheduled CI caller)

**Triggers:** `workflow_run`, `push`

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi250-caller.yml` |

**Jobs:** [Model CI](#model-ci-model-ci-4), [Torch pipeline CI](#torch-pipeline-ci-torch-pipeline-1), [Example CI](#example-ci-example-ci-1), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci-1)

## Event filters

- **workflow_run**
  - workflows: `Self-hosted runner (AMD scheduled CI caller)`
  - branches: `main`
  - types: `completed`
- **push**
  - branches: `run_amd_scheduled_ci_caller*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`self-scheduled-amd-mi250-caller.yml` [workflow_run, push]

- uses **`huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4`** (x4)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4`

## Jobs

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-daily-amd`
- `runner`: `mi250`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi250`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Torch pipeline CI (`torch-pipeline`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_pipelines_torch_gpu`
- `slack_report_channel`: `#transformers-ci-daily-amd`
- `runner`: `mi250`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi250`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Example CI (`example-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_examples_gpu`
- `slack_report_channel`: `#transformers-ci-daily-amd`
- `runner`: `mi250`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi250`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `slack_report_channel`: `#transformers-ci-daily-amd`
- `runner`: `mi250`
- `docker`: `huggingface/transformers-pytorch-deepspeed-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi250`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Self-hosted runner (AMD scheduled CI caller)

**Triggers:** `schedule`

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-caller.yml` |

## Schedule

- `17 5 * * *`

## Permissions

- `contents`: `read`

## Jobs

### Trigger Scheduled AMD CI (`run_scheduled_amd_ci`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `${{ always() }}` |

<details>
<summary>Steps (1)</summary>

1. **Trigger scheduled AMD CI via workflow\_run**

</details>

[Back to top](#contents)

# Self-hosted runner (benchmark)

**Triggers:** `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `benchmark.yml` |

## Event filters

- **push**
  - branches: `main`
- **pull_request**
  - types: `opened`, `labeled`, `reopened`, `synchronize`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `DATASET_ID` | `hf-benchmarks/transformers` |
| `MODEL_ID` | `meta-llama/Llama-3.1-8B-Instruct` |

**Concurrency:** group `${{ github.workflow }}-${{ github.head_ref || github.run_id }}`, cancel-in-progress: `true`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | job `benchmark` step `Run benchmark` env `HF_TOKEN` |
| `PUSH_TO_HUB_TOKEN` | job `benchmark` step `Run benchmark` env `PUSH_TO_HUB_TOKEN` |

## Jobs

### Benchmark (`benchmark`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.group }}` |
| Matrix | `group`: aws-g5-4xlarge-cache |
| Condition | `(github.event_name == 'pull_request' && contains( github.event.pull_request.labels.*.name, 'run-benchmark') )\|\|<br>(github.event_name == 'push' && github.ref == 'refs/heads/main')` |

<details>
<summary>Steps (4)</summary>

1. **Get repo**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `fetch-depth`: `1`
     - `persist-credentials`: `false`

2. **Install benchmark script dependencies**

3. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

4. **Run benchmark**
   - Env:
     - `HF_TOKEN`: `${{ secrets.HF_HUB_READ_TOKEN }}`
     - `PUSH_TO_HUB_TOKEN`: `${{ secrets.PUSH_TO_HUB_TOKEN }}`
     - `BRANCH_NAME`: `${{ github.head_ref || github.ref_name }}`

</details>

[Back to top](#contents)

# Self-hosted runner (Intel Gaudi3 scheduled CI caller)

**Triggers:** `repository_dispatch`, `workflow_dispatch`, `schedule`

| Property | Value |
|----------|-------|
| File | `self-scheduled-intel-gaudi3-caller.yml` |

**Jobs:** [Model CI](#model-ci-model-ci-5), [Pipeline CI](#pipeline-ci-pipeline-ci), [Example CI](#example-ci-example-ci-2), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci-2), [Trainer/FSDP CI](#trainerfsdp-ci-trainer-fsdp-ci-1)

## Schedule

- `17 2 * * *`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`self-scheduled-intel-gaudi3-caller.yml` [repository_dispatch, workflow_dispatch, schedule]

- uses **[self-scheduled-intel-gaudi.yml](#self-hosted-runner-scheduled-intel-gaudi)** (x5)
  - uses **[model_jobs_intel_gaudi.yml](#model-jobs-1)** (x2)
  - `send_results` uses [slack-report.yml](#ci-slack-report)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `ci_event`: `Scheduled CI (Intel) - Gaudi3`
- `runner_scale_set`: `itac-bm-emr-gaudi3-dell`
- `slack_report_channel`: `#transformers-ci-daily-intel-gaudi3`
- `report_repo_id`: `optimum-intel/transformers_daily_ci_intel_gaudi3`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Pipeline CI (`pipeline-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi) |

#### Inputs forwarded

- `job`: `run_pipelines_torch_gpu`
- `ci_event`: `Scheduled CI (Intel) - Gaudi3`
- `runner_scale_set`: `itac-bm-emr-gaudi3-dell`
- `slack_report_channel`: `#transformers-ci-daily-intel-gaudi3`
- `report_repo_id`: `optimum-intel/transformers_daily_ci_intel_gaudi3`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Example CI (`example-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi) |

#### Inputs forwarded

- `job`: `run_examples_gpu`
- `ci_event`: `Scheduled CI (Intel) - Gaudi3`
- `runner_scale_set`: `itac-bm-emr-gaudi3-dell`
- `slack_report_channel`: `#transformers-ci-daily-intel-gaudi3`
- `report_repo_id`: `optimum-intel/transformers_daily_ci_intel_gaudi3`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `ci_event`: `Scheduled CI (Intel) - Gaudi3`
- `runner_scale_set`: `itac-bm-emr-gaudi3-dell`
- `slack_report_channel`: `#transformers-ci-daily-intel-gaudi3`
- `report_repo_id`: `optimum-intel/transformers_daily_ci_intel_gaudi3`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Trainer/FSDP CI (`trainer-fsdp-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi) |

#### Inputs forwarded

- `job`: `run_trainer_and_fsdp_gpu`
- `ci_event`: `Scheduled CI (Intel) - Gaudi3`
- `runner_scale_set`: `itac-bm-emr-gaudi3-dell`
- `slack_report_channel`: `#transformers-ci-daily-intel-gaudi3`
- `report_repo_id`: `optimum-intel/transformers_daily_ci_intel_gaudi3`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Self-hosted runner (nightly-past-ci-caller)

**Triggers:** `schedule`, `push`

| Property | Value |
|----------|-------|
| File | `self-nightly-past-ci-caller.yml` |

**Jobs:** [Get number](#get-number-get_number), [TensorFlow 2.11](#tensorflow-211-run_past_ci_tensorflow_2-11), [TensorFlow 2.10](#tensorflow-210-run_past_ci_tensorflow_2-10), [TensorFlow 2.9](#tensorflow-29-run_past_ci_tensorflow_2-9), [TensorFlow 2.8](#tensorflow-28-run_past_ci_tensorflow_2-8), [TensorFlow 2.7](#tensorflow-27-run_past_ci_tensorflow_2-7), [TensorFlow 2.6](#tensorflow-26-run_past_ci_tensorflow_2-6), [TensorFlow 2.5](#tensorflow-25-run_past_ci_tensorflow_2-5)

## Schedule

- `17 2,14 * * *`

## Event filters

- **push**
  - branches: `run_past_ci*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`self-nightly-past-ci-caller.yml` [schedule, push]

- uses **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x7)
  - uses **[self-scheduled.yml](#nvidia-ci-job-definitions)** (x2)
    - uses **[model_jobs.yml](#model-jobs)** (x2)
      - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
    - `send_results` uses [slack-report.yml](#ci-slack-report)
    - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Get number (`get_number`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (1)</summary>

1. **Get number**
   - ID: `get_number`

</details>

### TensorFlow 2.11 (`run_past_ci_tensorflow_2-11`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 3 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.11`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.10 (`run_past_ci_tensorflow_2-10`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 4 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.10`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.9 (`run_past_ci_tensorflow_2-9`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 5 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.9`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.8 (`run_past_ci_tensorflow_2-8`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 6 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.8`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.7 (`run_past_ci_tensorflow_2-7`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 7 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.7`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.6 (`run_past_ci_tensorflow_2-6`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 8 && (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.6`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### TensorFlow 2.5 (`run_past_ci_tensorflow_2-5`)

| Property | Value |
|----------|-------|
| Uses workflow | [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci) |
| Depends on | `get_number` |
| Condition | `needs.get_number.outputs.run_number == 9 &&  (cancelled() != true) && ((github.event_name == 'push') && startsWith(github.ref_name, 'run_past_ci'))` |

#### Inputs forwarded

- `framework`: `tensorflow`
- `version`: `2.5`
- `sha`: `${{ github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Self-hosted runner scale set (AMD mi325 scheduled CI caller)

**Triggers:** `workflow_run`, `push`

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi325-caller.yml` |

**Jobs:** [Model CI](#model-ci-model-ci-6), [Torch pipeline CI](#torch-pipeline-ci-torch-pipeline-2), [Example CI](#example-ci-example-ci-3), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci-3)

## Event filters

- **workflow_run**
  - workflows: `Self-hosted runner (AMD scheduled CI caller)`
  - branches: `main`
  - types: `completed`
- **push**
  - branches: `run_amd_scheduled_ci_caller*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`self-scheduled-amd-mi325-caller.yml` [workflow_run, push]

- uses **`huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4`** (x4)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4`

## Jobs

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `amd-mi325`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi325`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`
- `env_file`: `/etc/podinfo/gha-gpu-isolation-settings`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Torch pipeline CI (`torch-pipeline`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_pipelines_torch_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `amd-mi325`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi325`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`
- `env_file`: `/etc/podinfo/gha-gpu-isolation-settings`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Example CI (`example-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_examples_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `amd-mi325`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi325`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`
- `env_file`: `/etc/podinfo/gha-gpu-isolation-settings`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `amd-mi325`
- `docker`: `huggingface/transformers-pytorch-deepspeed-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi325`
- `report_repo_id`: `optimum-amd/transformers_daily_ci`
- `env_file`: `/etc/podinfo/gha-gpu-isolation-settings`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Self-hosted runner scale set (AMD mi355 scheduled CI caller)

**Triggers:** `workflow_run`, `push`

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi355-caller.yml` |

**Jobs:** [Model CI](#model-ci-model-ci-7), [Torch pipeline CI](#torch-pipeline-ci-torch-pipeline-3), [Example CI](#example-ci-example-ci-4), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci-4)

## Event filters

- **workflow_run**
  - workflows: `Self-hosted runner (AMD scheduled CI caller)`
  - branches: `main`
  - types: `completed`
- **push**
  - branches: `run_amd_scheduled_ci_caller*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`self-scheduled-amd-mi355-caller.yml` [workflow_run, push]

- uses **`huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4`** (x4)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4`

## Jobs

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `hfc-amd-mi355`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi355`
- `report_repo_id`: `hf-transformers-bot/transformers-ci-dummy`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Torch pipeline CI (`torch-pipeline`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_pipelines_torch_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `hfc-amd-mi355`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi355`
- `report_repo_id`: `hf-transformers-bot/transformers-ci-dummy`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Example CI (`example-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_examples_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `hfc-amd-mi355`
- `docker`: `huggingface/transformers-pytorch-amd-gpu`
- `ci_event`: `Scheduled CI (AMD) - mi355`
- `report_repo_id`: `hf-transformers-bot/transformers-ci-dummy`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4` (external) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `slack_report_channel`: `#amd-hf-ci`
- `runner_group`: `hfc-amd-mi355`
- `docker`: `huggingface/testing-rocm7.0-preview`
- `ci_event`: `Scheduled CI (AMD) - mi355`
- `report_repo_id`: `hf-transformers-bot/transformers-ci-dummy`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Slow tests on important models (on Push - A10)

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `push-important-models.yml` |

**Jobs:** [Get all modified files](#get-all-modified-files-get_modified_models), [Model CI](#model-ci-model-ci-8)

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`push-important-models.yml` [push]

- `model-ci` uses [self-scheduled.yml](#nvidia-ci-job-definitions)
  - uses **[model_jobs.yml](#model-jobs)** (x2)
    - `collated_reports` uses [collated-reports.yml](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`)
  - `send_results` uses [slack-report.yml](#ci-slack-report)
  - `check_new_failures` uses [check_failed_tests.yml](#process-failed-tests)

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Get all modified files (`get_modified_models`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (3)</summary>

1. **Check out code**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

2. **Get changed files using \`actions/github-script\`**
   - ID: `get-changed-files`
   - Uses: `actions/github-script@v7.1.0`
   - With:
     - `script`: `` let files = [];  // Only handle push events if (context.eventName === 'push') {   const afterSha = context.payload.after;   const branchName = context.payload.ref.replace('refs/heads/', '');      let baseSha;      if (branchName === 'main') {     console.log('Push to main branch, comparing to parent commit');     // Get the parent commit of the pushed commit     const { data: commit } = await github.rest.repos.getCommit({       owner: context.repo.owner,       repo: context.repo.repo,       ref: afterSha     });     baseSha = commit.parents[0]?.sha;     if (!baseSha) {       throw new Error('No parent commit found for the pushed commit');     }   } else {     console.log(`Push to branch ${branchName}, comparing to main`);     baseSha = 'main';   }      const { data: comparison } = await github.rest.repos.compareCommits({     owner: context.repo.owner,     repo: context.repo.repo,     base: baseSha,     head: afterSha   });      // Include added, modified, and renamed files   files = comparison.files     .filter(file => file.status === 'added' || file.status === 'modified' || file.status === 'renamed')     .map(file => file.filename); }  // Include all files under src/transformers/ (not just models subdirectory) const filteredFiles = files.filter(file =>    file.startsWith('src/transformers/') );  core.setOutput('changed_files', filteredFiles.join(' ')); core.setOutput('any_changed', filteredFiles.length > 0 ? 'true' : 'false'); ``

3. **Parse changed files with Python**
   - ID: `set-matrix`
   - Condition: `steps.get-changed-files.outputs.any_changed == 'true'`
   - Env:
     - `CHANGED_FILES`: `${{ steps.get-changed-files.outputs.changed_files }}`

</details>

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |
| Depends on | `get_modified_models` |
| Condition | `needs.get_modified_models.outputs.matrix != '' && needs.get_modified_models.outputs.matrix != '[]'` |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-push`
- `docker`: `huggingface/transformers-all-latest-gpu:flash-attn`
- `ci_event`: `push`
- `report_repo_id`: `hf-internal-testing/transformers_ci_push`
- `commit_sha`: `${{ github.sha }}`
- `subdirs`: `${{ needs.get_modified_models.outputs.matrix }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# SSH into our runners

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `ssh-runner.yml` |

**Jobs:** [Get runner to use](#get-runner-to-use-get_runner), [SSH](#ssh-ssh_runner)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `runner_type` | - | Yes | - | Type of runner to test (a10) |
| `docker_image` | - | Yes | - | Name of the Docker image |
| `num_gpus` | - | Yes | - | Type of the number of gpus to use (`single` or `multi`) |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `OMP_NUM_THREADS` | `8` |
| `MKL_NUM_THREADS` | `8` |
| `RUN_SLOW` | `yes` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |
| `CUDA_VISIBLE_DEVICES` | `0,1` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `SLACK_CIFEEDBACK_CHANNEL` | job `ssh_runner` step `Store Slack infos` env `default_slack_channel` |
| `TAILSCALE_SSH_AUTHKEY` | job `ssh_runner` step `Tailscale` with `authkey` |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `ssh_runner` step `Tailscale` with `slackToken` |

## Jobs

### Get runner to use (`get_runner`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (2)</summary>

1. **Get runner to use**
   - Env:
     - `NUM_GPUS`: `${{ github.event.inputs.num_gpus }}`
     - `RUNNER_TYPE`: `${{ github.event.inputs.runner_type }}`

2. **Set runner to use**
   - ID: `set_runner`

</details>

### SSH (`ssh_runner`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ needs.get_runner.outputs.RUNNER }}` |
| Depends on | `get_runner` |

<details>
<summary>Steps (13)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ github.sha }}`

2. **Cleanup**

3. **Show installed libraries and their versions**

4. **NVIDIA-SMI**

5. **Create python alias**

6. **Install psutil for memory monitor**

7. **Download memory monitor script**

8. **Start memory monitor** `[continue-on-error]`

9. **Install utilities**

10. **Store Slack infos**
   - Env:
     - `GITHUB_ACTOR`: `${{ github.actor }}`

11. **Setup automatic environment for SSH login**

12. **Store Slack infos**
   - Env:
     - `user_slack_id`: `${{ secrets[format('{0}_{1}', env.github_actor, 'SLACK_ID')] }}`
     - `default_slack_channel`: `${{ secrets.SLACK_CIFEEDBACK_CHANNEL }}`

13. **Tailscale**
   - Uses: `huggingface/tailscale-action@7d53c9737e53934c30290b5524d1c9b4a7c98c8a`
   - With:
     - `authkey`: `${{ secrets.TAILSCALE_SSH_AUTHKEY }}`
     - `slackChannel`: `${{ env.SLACKCHANNEL }}`
     - `slackToken`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`
     - `waitForSSH`: `true`
     - `sshTimeout`: `15m`

</details>

[Back to top](#contents)

# Stale Bot

**Triggers:** `schedule`

| Property | Value |
|----------|-------|
| File | `stale.yml` |

## Schedule

- `0 8 * * *`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `close_stale_issues` env `GITHUB_TOKEN` |

## Jobs

### Close Stale Issues (`close_stale_issues`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `github.repository == 'huggingface/transformers'` |

**Permissions:**

- `issues`: `write`

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `GITHUB_TOKEN` | `${{ secrets.GITHUB_TOKEN }}` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `persist-credentials`: `false`

2. **Setup Python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.8`

3. **Install requirements**

4. **Close stale issues**

</details>

[Back to top](#contents)

# TRL CI bot

**Triggers:** `issue_comment`

This workflow allows trusted contributors to trigger TRL CI runs against specific Transformers commits by commenting `/trl-ci` on a PR in the TRL repo. It is meant to be used during the ongoing Trainer refactor/unbloat in Transformers, to help evaluate the downstream impact on TRL.

| Property | Value |
|----------|-------|
| File | `trl-ci-bot.yml` |

## Event filters

- **issue_comment**
  - types: `created`

## Permissions

- `contents`: `read`
- `pull-requests`: `read`
- `issues`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `TRL_CI_DISPATCH_TOKEN` | job `dispatch` step `Dispatch TRL workflow` env `GH_TOKEN`; job `dispatch` step `Find TRL workflow run URL` env `GH_TOKEN` |

## Jobs

### `dispatch`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.event.issue.pull_request && contains(github.event.comment.body, '/trl-ci')` |

<details>
<summary>Steps (6)</summary>

1. **Gate on trusted commenter**
   - ID: `trust`

2. **Reject untrusted commenter**
   - Condition: `steps.trust.outputs.trusted != 'true'`

3. **Fetch PR head SHA + number**
   - ID: `pr`
   - Condition: `steps.trust.outputs.trusted == 'true'`
   - Env:
     - `GH_TOKEN`: `${{ github.token }}`
     - `PR_URL`: `${{ github.event.issue.pull_request.url }}`

4. **Dispatch TRL workflow**
   - ID: `dispatch`
   - Condition: `steps.trust.outputs.trusted == 'true'`
   - Env:
     - `GH_TOKEN`: `${{ secrets.TRL_CI_DISPATCH_TOKEN }}`
     - `STEPS_PR_OUTPUTS_SHA`: `${{ steps.pr.outputs.sha }}`

5. **Find TRL workflow run URL**
   - ID: `find_run`
   - Condition: `steps.trust.outputs.trusted == 'true'`
   - Env:
     - `GH_TOKEN`: `${{ secrets.TRL_CI_DISPATCH_TOKEN }}`

6. **Comment back on PR with link**
   - Condition: `steps.trust.outputs.trusted == 'true'`
   - Env:
     - `GH_TOKEN`: `${{ github.token }}`
     - `STEPS_PR_OUTPUTS_SHA`: `${{ steps.pr.outputs.sha }}`
     - `STEPS_FIND_RUN_OUTPUTS_URL`: `${{ steps.find_run.outputs.url }}`

</details>

[Back to top](#contents)

# Update Transformers metadata

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `update_metdata.yml` |

## Event filters

- **push**
  - branches: `main`, `update_transformers_metadata*`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `LYSANDRE_HF_TOKEN` | job `build_and_package` step `Update metadata` env `HF_TOKEN` |

## Jobs

### `build_and_package`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

**Defaults:** shell `bash -l {0}`

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `persist-credentials`: `false`

2. **Setup environment**

3. **Update metadata**
   - Env:
     - `HF_TOKEN`: `${{ secrets.LYSANDRE_HF_TOKEN }}`

</details>

[Back to top](#contents)

# Upload PR Documentation

**Triggers:** `workflow_run`

| Property | Value |
|----------|-------|
| File | `upload_pr_documentation.yml` |

## Event filters

- **workflow_run**
  - workflows: `Build PR Documentation`
  - types: `completed`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

`upload_pr_documentation.yml` [workflow_run]

- `build` uses `huggingface/doc-builder/.github/workflows/upload_pr_documentation.yml@9ad2de8582b56c017cb530c1165116d40433f1c6`

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `COMMENT_BOT_TOKEN`, `HF_DOC_BUILD_PUSH`, `comment_bot_token`, `hf_token`

Permissions declared across the chain: `contents: read`

External workflows referenced: `huggingface/doc-builder/.github/workflows/upload_pr_documentation.yml@9ad2de8582b56c017cb530c1165116d40433f1c6`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_DOC_BUILD_PUSH` | job `build` secrets `hf_token` |
| `COMMENT_BOT_TOKEN` | job `build` secrets `comment_bot_token` |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Uses workflow | `huggingface/doc-builder/.github/workflows/upload_pr_documentation.yml@9ad2de8582b56c017cb530c1165116d40433f1c6` (external) |

#### Inputs forwarded

- `package_name`: `transformers`

#### Secrets forwarded

- `hf_token`: `${{ secrets.HF_DOC_BUILD_PUSH }}`
- `comment_bot_token`: `${{ secrets.COMMENT_BOT_TOKEN }}`

[Back to top](#contents)

# CI collated reports

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `collated-reports.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job` | string | Yes | - | - |
| `report_repo_id` | string | Yes | - | - |
| `machine_type` | string | Yes | - | - |
| `gpu_name` | string | Yes | - | Name of the GPU used for the job. Its enough that the value contains the name of the GPU, e.g. "noise-h100-more-noise". Case insensitive. |

## Permissions

- `contents`: `read`

## Called by

`collated-reports.yml`

- [model_jobs.yml](#collated-reports-collated_reports-1) (job: `collated_reports`)
  - **[self-scheduled.yml](#nvidia-ci-job-definitions)** (x2)
    - [push-important-models.yml](#model-ci-model-ci-8) (job: `model-ci`) - entry point
    - **[self-comment-ci.yml](#pr-comment-github-ci)** - entry point (x2)
    - [self-nightly-caller.yml](#model-ci-model-ci-2) (job: `model-ci`) - entry point
    - **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x2)
      - **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)
    - **[self-scheduled-caller.yml](#nvidia-ci)** - entry point (x7)
    - [self-scheduled-flash-attn-caller.yml](#model-ci-model-ci-1) (job: `model-ci`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `ACCESS_REPO_INFO_TOKEN` | job `collated_reports` step `Collated reports` env `ACCESS_REPO_INFO_TOKEN` |
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | job `collated_reports` step `Collated reports` env `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` |

## Jobs

### Collated reports (`collated_reports`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `always()` |

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v4.3.1**
   - With:
     - `persist-credentials`: `false`

2. **actions/download-artifact@v4.3.0**

3. **Collated reports**
   - Env:
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `CI_SHA`: `${{ github.sha }}`
     - `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`: `${{ secrets.TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN }}`
     - `MACHINE_TYPE`: `${{ inputs.machine_type }}`
     - `JOB`: `${{ inputs.job }}`
     - `REPORT_REPO_ID`: `${{ inputs.report_repo_id }}`
     - `GPU_NAME`: `${{ inputs.gpu_name }}`

</details>

[Back to top](#contents)

# CI slack report

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `slack-report.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job` | string | Yes | - | - |
| `slack_report_channel` | string | Yes | - | - |
| `setup_status` | string | Yes | - | - |
| `folder_slices` | string | Yes | - | - |
| `quantization_matrix` | string | Yes | - | - |
| `ci_event` | string | Yes | - | - |
| `report_repo_id` | string | Yes | - | - |
| `commit_sha` | string | No | - | - |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `is_slack_reporting_job_ok` | Whether the send_results job succeeded (not failed) | `${{ jobs.send_results.result != 'failure' }}` |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | `${{ secrets.TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN }}` |

## Called by

`slack-report.yml`

- [self-scheduled-intel-gaudi.yml](#slack-report-send_results-1) (job: `send_results`)
  - **[self-scheduled-intel-gaudi3-caller.yml](#self-hosted-runner-intel-gaudi3-scheduled-ci-caller)** - entry point (x5)
- [self-scheduled.yml](#slack-report-send_results) (job: `send_results`)
  - [push-important-models.yml](#model-ci-model-ci-8) (job: `model-ci`) - entry point
  - **[self-comment-ci.yml](#pr-comment-github-ci)** - entry point (x2)
  - [self-nightly-caller.yml](#model-ci-model-ci-2) (job: `model-ci`) - entry point
  - **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x2)
    - **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)
  - **[self-scheduled-caller.yml](#nvidia-ci)** - entry point (x7)
  - [self-scheduled-flash-attn-caller.yml](#model-ci-model-ci-1) (job: `model-ci`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | workflow env `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` |
| `GITHUB_TOKEN` | job `send_results` step `actions/download-artifact@v8.0.1` with `github-token` |
| `CI_SLACK_BOT_TOKEN` | job `send_results` step `Send message to Slack` env `CI_SLACK_BOT_TOKEN` |
| `CI_SLACK_CHANNEL_ID` | job `send_results` step `Send message to Slack` env `CI_SLACK_CHANNEL_ID` |
| `CI_SLACK_CHANNEL_ID_DAILY` | job `send_results` step `Send message to Slack` env `CI_SLACK_CHANNEL_ID_DAILY` |
| `CI_SLACK_CHANNEL_DUMMY_TESTS` | job `send_results` step `Send message to Slack` env `CI_SLACK_CHANNEL_DUMMY_TESTS` |
| `ACCESS_REPO_INFO_TOKEN` | job `send_results` step `Send message to Slack` env `ACCESS_REPO_INFO_TOKEN` |

## Jobs

### Send results to webhook (`send_results`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `always() && !cancelled()` |

<details>
<summary>Steps (6)</summary>

1. **Preliminary job status**
   - Env:
     - `setup_status`: `${{ inputs.setup_status }}`

2. **actions/checkout@v4.3.1**
   - With:
     - `fetch-depth`: `2`
     - `ref`: `${{ (github.event_name == 'issue_comment' || github.event_name == 'pull_request_target') && 'main' || (inputs.commit_sha || github.sha) }}`
     - `persist-credentials`: `false`

3. **actions/download-artifact@v8.0.1**
   - With:
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
   - Env:
     - `ACTIONS_ARTIFACT_MAX_ARTIFACT_COUNT`: `2000`

4. **Prepare some setup values**

5. **Send message to Slack**
   - Env:
     - `CI_SLACK_BOT_TOKEN`: `${{ secrets.CI_SLACK_BOT_TOKEN }}`
     - `CI_SLACK_CHANNEL_ID`: `${{ secrets.CI_SLACK_CHANNEL_ID }}`
     - `CI_SLACK_CHANNEL_ID_DAILY`: `${{ secrets.CI_SLACK_CHANNEL_ID_DAILY }}`
     - `CI_SLACK_CHANNEL_DUMMY_TESTS`: `${{ secrets.CI_SLACK_CHANNEL_DUMMY_TESTS }}`
     - `SLACK_REPORT_CHANNEL`: `${{ inputs.slack_report_channel }}`
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `CI_EVENT`: `${{ inputs.ci_event }}`
     - `CI_TITLE`: `${{ github.event.head_commit.message }}`
     - `CI_SHA`: `${{ inputs.commit_sha || github.sha }}`
     - `CI_TEST_JOB`: `${{ inputs.job }}`
     - `SETUP_STATUS`: `${{ inputs.setup_status }}`
     - `REPORT_REPO_ID`: `${{ inputs.report_repo_id }}`
     - `quantization_matrix`: `${{ inputs.quantization_matrix }}`
     - `folder_slices`: `${{ inputs.folder_slices }}`

6. **Failure table artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `ci_results_${{ inputs.job }}`

</details>

[Back to top](#contents)

# Doctest job

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `doctest_job.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job_splits` | string | Yes | - | - |
| `split_keys` | string | Yes | - | - |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `RUN_SLOW` | `yes` |
| `OMP_NUM_THREADS` | `16` |
| `MKL_NUM_THREADS` | `16` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |

## Called by

`doctest_job.yml`

- [doctests.yml](#call-doctest-jobs-call_doctest_job) (job: `call_doctest_job`) - entry point

## Jobs

### `run_doctests`

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g5-4xlarge-cache` |
| Matrix | `split_keys`: ${{ fromJson(inputs.split_keys) }} |

<details>
<summary>Steps (9)</summary>

1. **Update clone**

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **GPU visibility**

4. **Show installed libraries and their versions**

5. **Get doctest files**

6. **Set \`split\_keys\`**
   - Env:
     - `MATRIX_SPLIT_KEYS`: `${{ matrix.split_keys }}`

7. **Run doctests**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: doc\_tests\_gpu\_test\_reports\_${{ env.split\_keys }}**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `doc_tests_gpu_test_reports_${{ env.split_keys }}`
     - `path`: `/transformers/reports/doc_tests_gpu_${{ env.split_keys }}`

</details>

[Back to top](#contents)

# Get PR commit SHA

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `get-pr-info.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `pr_number` | string | Yes | - | - |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `PR_HEAD_REPO_FULL_NAME` | The full name of the repository from which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_HEAD_REPO_FULL_NAME }}` |
| `PR_BASE_REPO_FULL_NAME` | The full name of the repository to which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_BASE_REPO_FULL_NAME }}` |
| `PR_HEAD_REPO_OWNER` | The owner of the repository from which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_HEAD_REPO_OWNER }}` |
| `PR_BASE_REPO_OWNER` | The owner of the repository to which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_BASE_REPO_OWNER }}` |
| `PR_HEAD_REPO_NAME` | The name of the repository from which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_HEAD_REPO_NAME }}` |
| `PR_BASE_REPO_NAME` | The name of the repository to which the pull request is created | `${{ jobs.get-pr-info.outputs.PR_BASE_REPO_NAME }}` |
| `PR_HEAD_REF` | The branch name of the pull request in the head repository | `${{ jobs.get-pr-info.outputs.PR_HEAD_REF }}` |
| `PR_BASE_REF` | The branch name in the base repository (to merge into) | `${{ jobs.get-pr-info.outputs.PR_BASE_REF }}` |
| `PR_HEAD_SHA` | The head sha of the pull request branch in the head repository | `${{ jobs.get-pr-info.outputs.PR_HEAD_SHA }}` |
| `PR_BASE_SHA` | The head sha of the target branch in the base repository | `${{ jobs.get-pr-info.outputs.PR_BASE_SHA }}` |
| `PR_MERGE_COMMIT_SHA` | The sha of the merge commit for the pull request (created by GitHub) in the base repository | `${{ jobs.get-pr-info.outputs.PR_MERGE_COMMIT_SHA }}` |
| `PR_MERGE_COMMIT_BASE_SHA` | The sha of the parent commit of the merge commit on the target branch in the base repository | `${{ jobs.get-pr-info.outputs.PR_MERGE_COMMIT_BASE_SHA }}` |
| `PR_HEAD_COMMIT_DATE` | The date of the head sha of the pull request branch in the head repository | `${{ jobs.get-pr-info.outputs.PR_HEAD_COMMIT_DATE }}` |
| `PR_MERGE_COMMIT_DATE` | The date of the merge commit for the pull request (created by GitHub) in the base repository | `${{ jobs.get-pr-info.outputs.PR_MERGE_COMMIT_DATE }}` |
| `PR_HEAD_COMMIT_TIMESTAMP` | The timestamp of the head sha of the pull request branch in the head repository | `${{ jobs.get-pr-info.outputs.PR_HEAD_COMMIT_TIMESTAMP }}` |
| `PR_MERGE_COMMIT_TIMESTAMP` | The timestamp of the merge commit for the pull request (created by GitHub) in the base repository | `${{ jobs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}` |
| `PR` | The PR | `${{ jobs.get-pr-info.outputs.PR }}` |
| `PR_FILES` | The files touched in the PR | `${{ jobs.get-pr-info.outputs.PR_FILES }}` |

## Permissions

- `contents`: `read`

## Called by

`get-pr-info.yml`

- [pr-repo-consistency-bot.yml](#get-pr-commit-sha-get-pr-info-2) (job: `get-pr-info`) - entry point
- [pr_build_doc_with_comment.yml](#get-pr-commit-sha-get-pr-info) (job: `get-pr-info`) - entry point
- [pr_slow_ci_suggestion.yml](#get-pr-commit-sha-get-pr-info-3) (job: `get-pr-info`) - entry point
- [self-comment-ci.yml](#get-pr-commit-sha-get-pr-info-1) (job: `get-pr-info`) - entry point

## Jobs

### Get PR commit SHA better (`get-pr-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `${{ inputs.pr_number != '' }}` |

<details>
<summary>Steps (2)</summary>

1. **Extract PR details**
   - ID: `pr_info`
   - Uses: `actions/github-script@v6.4.1`
   - With:
     - `script`: `const pull_number = parseInt(process.env.PR_NUMBER, 10);  const { data: pr } = await github.rest.pulls.get({   owner: context.repo.owner,   repo: context.repo.repo,   pull_number, });  const { data: head_commit } = await github.rest.repos.getCommit({   owner: pr.head.repo.owner.login,   repo: pr.head.repo.name,   ref: pr.head.ref });  const { data: merge_commit } = await github.rest.repos.getCommit({   owner: pr.base.repo.owner.login,   repo: pr.base.repo.name,   ref: pr.merge_commit_sha, });  const { data: files } = await github.rest.pulls.listFiles({   owner: context.repo.owner,   repo: context.repo.repo,   pull_number, });  core.setOutput('head_repo_full_name', pr.head.repo.full_name); core.setOutput('base_repo_full_name', pr.base.repo.full_name); core.setOutput('head_repo_owner', pr.head.repo.owner.login); core.setOutput('base_repo_owner', pr.base.repo.owner.login); core.setOutput('head_repo_name', pr.head.repo.name); core.setOutput('base_repo_name', pr.base.repo.name); core.setOutput('head_ref', pr.head.ref); core.setOutput('base_ref', pr.base.ref); core.setOutput('head_sha', pr.head.sha); core.setOutput('base_sha', pr.base.sha); core.setOutput('merge_commit_base_sha', merge_commit.parents[0].sha); core.setOutput('merge_commit_sha', pr.merge_commit_sha); core.setOutput('pr', pr);  core.setOutput('head_commit_date', head_commit.commit.committer.date); core.setOutput('merge_commit_date', merge_commit.commit.committer.date);  core.setOutput('files', files);              console.log('PR head commit:', {   head_commit: head_commit,   commit: head_commit.commit,   date: head_commit.commit.committer.date });  console.log('PR merge commit:', {   merge_commit: merge_commit,   commit: merge_commit.commit,   date: merge_commit.commit.committer.date });  console.log('PR Info:', {   pr_info: pr });`
   - Env:
     - `PR_NUMBER`: `${{ inputs.pr_number }}`

2. **Convert dates to timestamps**
   - ID: `get_timestamps`
   - Env:
     - `head_commit_date`: `${{ steps.pr_info.outputs.head_commit_date }}`
     - `merge_commit_date`: `${{ steps.pr_info.outputs.merge_commit_date }}`

</details>

[Back to top](#contents)

# Get PR number

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `get-pr-number.yml` |

## Workflow call API

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `PR_NUMBER` | The extracted PR number | `${{ jobs.get-pr-number.outputs.PR_NUMBER }}` |

## Permissions

- `contents`: `read`

## Called by

`get-pr-number.yml`

- [pr-repo-consistency-bot.yml](#get-pr-number-get-pr-number-2) (job: `get-pr-number`) - entry point
- [pr_build_doc_with_comment.yml](#get-pr-number-get-pr-number) (job: `get-pr-number`) - entry point
- [pr_slow_ci_suggestion.yml](#get-pr-number-get-pr-number-3) (job: `get-pr-number`) - entry point
- [self-comment-ci.yml](#get-pr-number-get-pr-number-1) (job: `get-pr-number`) - entry point

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (3)</summary>

1. **Get PR number**
   - Env:
     - `issue_number`: `${{ github.event.issue.number }}`
     - `is_pull_request_issue`: `${{ github.event.issue.pull_request != null }}`
     - `pr_number`: `${{ github.event.pull_request.number }}`
     - `is_pull_request`: `${{ github.event.pull_request != null }}`
     - `event_number`: `${{ github.event.number }}`

2. **Check PR number**

3. **Set PR number**
   - ID: `set_pr_number`

</details>

[Back to top](#contents)

# model jobs

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `model_jobs.yml` |

**Jobs:** [`run_models_gpu`](#run_models_gpu), [Collated Reports](#collated-reports-collated_reports-1)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `folder_slices` | string | Yes | - | - |
| `machine_type` | string | Yes | - | - |
| `slice_id` | number | Yes | - | - |
| `docker` | string | Yes | - | - |
| `commit_sha` | string | No | - | - |
| `report_name_prefix` | string | No | `run_models_gpu` | - |
| `runner_type` | string | No | - | - |
| `report_repo_id` | string | No | - | - |
| `pytest_marker` | string | No | - | - |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `OMP_NUM_THREADS` | `8` |
| `MKL_NUM_THREADS` | `8` |
| `RUN_SLOW` | `yes` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |
| `CUDA_VISIBLE_DEVICES` | `0,1` |

## Called by

`model_jobs.yml`

- **[self-scheduled.yml](#nvidia-ci-job-definitions)** (x2)
  - [push-important-models.yml](#model-ci-model-ci-8) (job: `model-ci`) - entry point
  - **[self-comment-ci.yml](#pr-comment-github-ci)** - entry point (x2)
  - [self-nightly-caller.yml](#model-ci-model-ci-2) (job: `model-ci`) - entry point
  - **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x2)
    - **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)
  - **[self-scheduled-caller.yml](#nvidia-ci)** - entry point (x7)
  - [self-scheduled-flash-attn-caller.yml](#model-ci-model-ci-1) (job: `model-ci`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |

## Jobs

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ inputs.machine_type }}` |
| Matrix | `folders`: ${{ fromJson(inputs.folder_slices)[inputs.slice_id] }} |

<details>
<summary>Steps (16)</summary>

1. **Echo input and matrix info**
   - Env:
     - `folder_slices`: `${{ inputs.folder_slices }}`
     - `matrix_folders`: `${{ matrix.folders }}`
     - `slice_data`: `${{ toJson(fromJson(inputs.folder_slices)[inputs.slice_id]) }}`

2. **Echo folder ${{ matrix.folders }}**
   - Env:
     - `matrix_folders_raw`: `${{ matrix.folders }}`

3. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

4. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

5. **Update / Install some packages (for Past CI)**
   - Condition: `${{ contains(inputs.docker, '-past-') }}`

6. **Update / Install some packages (for Past CI)**
   - Condition: `${{ contains(inputs.docker, '-past-') && contains(inputs.docker, '-pytorch-') }}`

7. **NVIDIA-SMI**

8. **Environment**

9. **Show installed libraries and their versions**

10. **Set \`machine\_type\` for report and artifact names**
   - ID: `set_machine_type`
   - Env:
     - `input_machine_type`: `${{ inputs.machine_type }}`

11. **Create report directory if it doesn't exist**
   - Env:
     - `report_name_prefix`: `${{ inputs.report_name_prefix }}`

12. **Run all tests on GPU**
   - Env:
     - `report_name_prefix`: `${{ inputs.report_name_prefix }}`
     - `pytest_marker`: `${{ inputs.pytest_marker }}`
     - `model`: `${{ matrix.folders }}`

13. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`
   - Env:
     - `report_name_prefix`: `${{ inputs.report_name_prefix }}`

14. **Captured information** `[continue-on-error]`
   - Condition: `${{ failure() }}`
   - Env:
     - `report_name_prefix`: `${{ inputs.report_name_prefix }}`

15. **Copy test\_outputs.txt** `[continue-on-error]`
   - Condition: `${{ always() }}`
   - Env:
     - `report_name_prefix`: `${{ inputs.report_name_prefix }}`

16. **Test suite reports artifacts: ${{ env.machine\_type }}\_${{ inputs.report\_name\_prefix }}\_${{ env.matrix\_folders }}\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`

</details>

### Collated Reports (`collated_reports`)

| Property | Value |
|----------|-------|
| Uses workflow | [CI collated reports](#ci-collated-reports) (`@6abd9725ee7d809dc974991f8ff6c958afb63a3a`) |
| Depends on | `run_models_gpu` |
| Condition | `${{ always() && inputs.runner_type != '' }}` |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `gpu_name`: `${{ inputs.runner_type }}`
- `machine_type`: `${{ needs.run_models_gpu.outputs.machine_type }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# model jobs

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `model_jobs_intel_gaudi.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `folder_slices` | string | Yes | - | - |
| `slice_id` | number | Yes | - | - |
| `runner` | string | Yes | - | - |
| `machine_type` | string | Yes | - | - |
| `report_name_prefix` | string | No | `run_models_gpu` | - |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `RUN_SLOW` | `yes` |
| `PT_HPU_LAZY_MODE` | `0` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `PT_ENABLE_INT64_SUPPORT` | `1` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `HF_HOME` | `/mnt/cache/.cache/huggingface` |

## Called by

`model_jobs_intel_gaudi.yml`

- **[self-scheduled-intel-gaudi.yml](#self-hosted-runner-scheduled-intel-gaudi)** (x2)
  - **[self-scheduled-intel-gaudi3-caller.yml](#self-hosted-runner-intel-gaudi3-scheduled-ci-caller)** - entry point (x5)

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |

## Jobs

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ inputs.runner }}` |
| Matrix | `folders`: ${{ fromJson(inputs.folder_slices)[inputs.slice_id] }} |

<details>
<summary>Steps (12)</summary>

1. **Echo input and matrix info**
   - Env:
     - `FOLDER_SLICES`: `${{ inputs.folder_slices }}`
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`
     - `SLICE`: `${{ toJson(fromJson(inputs.folder_slices)[inputs.slice_id]) }}`

2. **Echo folder ${{ matrix.folders }}**
   - Env:
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`

3. **Checkout**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

4. **Install dependencies**

5. **HL-SMI**

6. **Environment**

7. **Show installed libraries and their versions**

8. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `MACHINE_TYPE`: `${{ inputs.machine_type }}`

9. **Run all tests on Gaudi**
   - Env:
     - `REPORT_NAME_PREFIX`: `${{ inputs.report_name_prefix }}`
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`

10. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`
   - Env:
     - `REPORT_NAME_PREFIX`: `${{ inputs.report_name_prefix }}`
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`

11. **Run test**
   - Env:
     - `REPORT_NAME_PREFIX`: `${{ inputs.report_name_prefix }}`
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`

12. **Test suite reports artifacts: ${{ env.machine\_type }}\_${{ inputs.report\_name\_prefix }}\_${{ env.matrix\_folders }}\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`
     - `path`: `reports/${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ matrix.folders }}_test_reports`

</details>

[Back to top](#contents)

# Nvidia CI (job definitions)

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `self-scheduled.yml` |
| Default runs-on | `group: ${{ matrix.machine_type }}` |

**Jobs:**

- [Setup](#setup-setup-4)
- [`run_models_gpu`](#run_models_gpu-2)
- [`run_trainer_and_fsdp_gpu`](#run_trainer_and_fsdp_gpu)
- [PyTorch pipelines](#pytorch-pipelines-run_pipelines_torch_gpu)
- [Examples directory](#examples-directory-run_examples_gpu)
- [Torch CUDA extension tests](#torch-cuda-extension-tests-run_torch_cuda_extensions_gpu)
- [`run_quantization_torch_gpu`](#run_quantization_torch_gpu)
- [Kernel tests](#kernel-tests-run_kernels_gpu)
- [Extract warnings in CI artifacts](#extract-warnings-in-ci-artifacts-run_extract_warnings)
- [Slack Report](#slack-report-send_results)
- [Check new failures](#check-new-failures-check_new_failures)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job` | string | Yes | - | - |
| `slack_report_channel` | string | Yes | - | - |
| `docker` | string | Yes | - | - |
| `ci_event` | string | Yes | - | - |
| `working-directory-prefix` | string | No | - | - |
| `report_repo_id` | string | Yes | - | - |
| `commit_sha` | string | No | - | - |
| `runner_type` | string | No | - | - |
| `subdirs` | string | No | - | - |
| `pytest_marker` | string | No | - | - |
| `pr_number` | string | No | - | - |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `is_infrastructure_ok` | Whether the CI infrastructure (slack reporting and failure checking) succeeded | `${{ jobs.send_results.outputs.is_slack_reporting_job_ok == 'true' && jobs.check_new_failures.outputs.is_check_failures_ok == 'true' }}` |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `OMP_NUM_THREADS` | `8` |
| `MKL_NUM_THREADS` | `8` |
| `RUN_SLOW` | `yes` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |
| `CUDA_VISIBLE_DEVICES` | `0,1` |

## Called by

`self-scheduled.yml`

- [push-important-models.yml](#model-ci-model-ci-8) (job: `model-ci`) - entry point
- **[self-comment-ci.yml](#pr-comment-github-ci)** - entry point (x2)
- [self-nightly-caller.yml](#model-ci-model-ci-2) (job: `model-ci`) - entry point
- **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x2)
  - **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)
- **[self-scheduled-caller.yml](#nvidia-ci)** - entry point (x7)
- [self-scheduled-flash-attn-caller.yml](#model-ci-model-ci-1) (job: `model-ci`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `GITHUB_TOKEN` | job `run_extract_warnings` step `actions/download-artifact@v8.0.1` with `github-token` |
| `ACCESS_REPO_INFO_TOKEN` | job `run_extract_warnings` step `Extract warnings in CI artifacts` env `access_token` |

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `contains(fromJSON('["run_models_gpu", "run_trainer_and_fsdp_gpu", "run_quantization_torch_gpu"]'), inputs.job)` |

<details>
<summary>Steps (6)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

2. **Cleanup**

3. **Show installed libraries and their versions**

4. **Identify models to test**
   - ID: `set-matrix`
   - Condition: `contains(fromJSON('["run_models_gpu", "run_trainer_and_fsdp_gpu"]'), inputs.job)`
   - Env:
     - `job`: `${{ inputs.job }}`
     - `subdirs`: `${{ inputs.subdirs }}`
     - `NUM_SLICES`: `2`

5. **Identify quantization method to test**
   - ID: `set-matrix-quantization`
   - Condition: `${{ inputs.job == 'run_quantization_torch_gpu' }}`
   - Env:
     - `subdirs`: `${{ inputs.subdirs || 'None' }}`

6. **NVIDIA-SMI**

</details>

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs) |
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache; `slice_id`: ${{ fromJSON(needs.setup.outputs.slice_ids) }} |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_models_gpu' }}` |

#### Inputs forwarded

- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `machine_type`: `${{ matrix.machine_type }}`
- `slice_id`: `${{ matrix.slice_id }}`
- `docker`: `${{ inputs.docker }}`
- `commit_sha`: `${{ inputs.commit_sha || github.sha }}`
- `runner_type`: `${{ inputs.runner_type }}`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `pytest_marker`: `${{ inputs.pytest_marker }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `run_trainer_and_fsdp_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs) |
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache; `slice_id`: 0, 1, 2 |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_trainer_and_fsdp_gpu' }}` |

#### Inputs forwarded

- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `machine_type`: `${{ matrix.machine_type }}`
- `slice_id`: `${{ matrix.slice_id }}`
- `docker`: `${{ inputs.docker }}`
- `commit_sha`: `${{ inputs.commit_sha || github.sha }}`
- `runner_type`: `${{ inputs.runner_type }}`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `report_name_prefix`: `run_trainer_and_fsdp_gpu`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### PyTorch pipelines (`run_pipelines_torch_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `${{ inputs.job == 'run_pipelines_torch_gpu' }}` |

<details>
<summary>Steps (9)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **NVIDIA-SMI**

4. **Environment**

5. **Show installed libraries and their versions**

6. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `matrix_machine_type`: `${{ matrix.machine_type }}`

7. **Run all pipeline tests on GPU**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_pipelines\_torch\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`

</details>

### Examples directory (`run_examples_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: aws-g5-4xlarge-cache |
| Condition | `${{ inputs.job == 'run_examples_gpu' }}` |

<details>
<summary>Steps (9)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **NVIDIA-SMI**

4. **Environment**

5. **Show installed libraries and their versions**

6. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `matrix_machine_type`: `${{ matrix.machine_type }}`

7. **Run examples tests on GPU**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_examples\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_examples_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_examples_gpu_test_reports`

</details>

### Torch CUDA extension tests (`run_torch_cuda_extensions_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `${{ inputs.job == 'run_torch_cuda_extensions_gpu' }}` |

<details>
<summary>Steps (13)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **Update / Install some packages (for Past CI)**
   - Condition: `${{ contains(inputs.docker, '-past-') && contains(inputs.docker, '-pytorch-') }}`

4. **Remove cached torch extensions**

5. **Pre build DeepSpeed \*again\* (for daily CI)**
   - Condition: `${{ contains(inputs.ci_event, 'Daily CI') }}`

6. **Pre build DeepSpeed \*again\* (for nightly & Past CI)**
   - Condition: `${{ contains(inputs.ci_event, 'Nightly CI') || contains(inputs.ci_event, 'Past CI') }}`

7. **NVIDIA-SMI**

8. **Environment**

9. **Show installed libraries and their versions**

10. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `matrix_machine_type`: `${{ matrix.machine_type }}`

11. **Run all tests on GPU**

12. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`
   - Env:
     - `working_directory_prefix`: `${{ inputs.working-directory-prefix }}`

13. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_torch\_cuda\_extensions\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`
     - `path`: `${{ inputs.working-directory-prefix }}/transformers/reports/${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`

</details>

### `run_quantization_torch_gpu`

| Property | Value |
|----------|-------|
| Matrix | `folders`: ${{ fromJson(needs.setup.outputs.quantization_matrix) }}; `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_quantization_torch_gpu' }}` |

<details>
<summary>Steps (10)</summary>

1. **Echo folder ${{ matrix.folders }}**
   - Env:
     - `matrix_folders_raw`: `${{ matrix.folders }}`

2. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

3. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

4. **NVIDIA-SMI**

5. **Environment**

6. **Show installed libraries and their versions**

7. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `matrix_machine_type`: `${{ matrix.machine_type }}`

8. **Run quantization tests on GPU**
   - Env:
     - `folders`: `${{ matrix.folders }}`

9. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

10. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_quantization\_torch\_gpu\_${{ env.matrix\_folders }}\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_quantization_torch_gpu_${{ env.matrix_folders }}_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_quantization_torch_gpu_${{ env.matrix_folders }}_test_reports`

</details>

### Kernel tests (`run_kernels_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: aws-g5-4xlarge-cache |
| Condition | `${{ inputs.job == 'run_kernels_gpu' }}` |

<details>
<summary>Steps (10)</summary>

1. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

2. **Reinstall transformers in edit mode**

3. **Install kernels**

4. **NVIDIA-SMI**

5. **Environment**

6. **Show installed libraries and their versions**

7. **Set \`machine\_type\` for report and artifact names**
   - Env:
     - `matrix_machine_type`: `${{ matrix.machine_type }}`

8. **Run kernel tests on GPU**

9. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

10. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_kernels\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_kernels_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_kernels_gpu_test_reports`

</details>

### Extract warnings in CI artifacts (`run_extract_warnings`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `setup`, `run_models_gpu` |
| Condition | `${{ always() && inputs.job == 'run_models_gpu' }}` |

<details>
<summary>Steps (8)</summary>

1. **Checkout transformers**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `persist-credentials`: `false`

2. **Install transformers**

3. **Show installed libraries and their versions**

4. **Create output directory**

5. **actions/download-artifact@v8.0.1**
   - With:
     - `path`: `warnings_in_ci`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
   - Env:
     - `ACTIONS_ARTIFACT_MAX_ARTIFACT_COUNT`: `2000`

6. **Show artifacts**

7. **Extract warnings in CI artifacts**
   - Env:
     - `github_run_id`: `${{ github.run_id }}`
     - `access_token`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`

8. **Upload artifact**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `warnings_in_ci`
     - `path`: `warnings_in_ci/selected_warnings.json`

</details>

### Slack Report (`send_results`)

| Property | Value |
|----------|-------|
| Uses workflow | [CI slack report](#ci-slack-report) |
| Depends on | `setup`, `run_models_gpu`, `run_trainer_and_fsdp_gpu`, `run_pipelines_torch_gpu`, `run_examples_gpu`, `run_torch_cuda_extensions_gpu`, `run_quantization_torch_gpu`, `run_kernels_gpu`, `run_extract_warnings` |
| Condition | `always() && !cancelled()` |

#### Inputs forwarded

- `job`: `${{ inputs.job }}`
- `setup_status`: `${{ needs.setup.result }}`
- `slack_report_channel`: `${{ inputs.slack_report_channel }}`
- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `quantization_matrix`: `${{ needs.setup.outputs.quantization_matrix }}`
- `ci_event`: `${{ inputs.ci_event }}`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Check new failures (`check_new_failures`)

| Property | Value |
|----------|-------|
| Uses workflow | [Process failed tests](#process-failed-tests) |
| Depends on | `send_results` |
| Condition | `${{ always() && needs.send_results.result == 'success' }}` |

#### Inputs forwarded

- `docker`: `${{ inputs.docker }}`
- `commit_sha`: `${{ inputs.commit_sha || github.sha }}`
- `job`: `${{ inputs.job }}`
- `slack_report_channel`: `${{ inputs.slack_report_channel }}`
- `ci_event`: `${{ inputs.ci_event }}`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `pr_number`: `${{ inputs.pr_number }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Process failed tests

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `check_failed_tests.yml` |
| Default runs-on | `group: aws-g5-4xlarge-cache` |

**Jobs:** [Setup matrix for finding commits](#setup-matrix-for-finding-commits-setup_check_new_failures), [Find commits for new failing tests](#find-commits-for-new-failing-tests-check_new_failures), [process bad commit reports](#process-bad-commit-reports-process_new_failures_with_commit_info)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `docker` | string | Yes | - | - |
| `job` | string | Yes | - | - |
| `slack_report_channel` | string | Yes | - | - |
| `ci_event` | string | Yes | - | - |
| `report_repo_id` | string | Yes | - | - |
| `commit_sha` | string | No | - | - |
| `pr_number` | string | No | - | - |
| `max_num_runners` | number | No | `4` | - |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `is_check_failures_ok` | Whether the failure checking infrastructure succeeded | `${{ jobs.check_new_failures.result != 'failure' && jobs.process_new_failures_with_commit_info.result != 'failure' }}` |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `HF_HOME` | `/mnt/cache` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `OMP_NUM_THREADS` | `8` |
| `MKL_NUM_THREADS` | `8` |
| `RUN_SLOW` | `yes` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `TF_FORCE_GPU_ALLOW_GROWTH` | `true` |
| `CUDA_VISIBLE_DEVICES` | `0,1` |

## Called by

`check_failed_tests.yml`

- [self-scheduled.yml](#check-new-failures-check_new_failures) (job: `check_new_failures`)
  - [push-important-models.yml](#model-ci-model-ci-8) (job: `model-ci`) - entry point
  - **[self-comment-ci.yml](#pr-comment-github-ci)** - entry point (x2)
  - [self-nightly-caller.yml](#model-ci-model-ci-2) (job: `model-ci`) - entry point
  - **[self-past-caller.yml](#self-hosted-runner-past-ci)** (x2)
    - **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)
  - **[self-scheduled-caller.yml](#nvidia-ci)** - entry point (x7)
  - [self-scheduled-flash-attn-caller.yml](#model-ci-model-ci-1) (job: `model-ci`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `GITHUB_TOKEN` | job `check_new_failures` step `actions/download-artifact@v8.0.1` with `github-token`; job `process_new_failures_with_commit_info` step `actions/download-artifact@v8.0.1` with `github-token` |
| `ACCESS_REPO_INFO_TOKEN` | job `check_new_failures` step `Get `END_SHA` from previous CI runs of the same workflow` env `ACCESS_TOKEN`; job `process_new_failures_with_commit_info` step `Process report` env `ACCESS_REPO_INFO_TOKEN` |
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | job `process_new_failures_with_commit_info` step `Process report` env `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `process_new_failures_with_commit_info` step `Send processed report` env `SLACK_BOT_TOKEN` |

## Jobs

### Setup matrix for finding commits (`setup_check_new_failures`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

<details>
<summary>Steps (2)</summary>

1. **actions/download-artifact@v8.0.1** `[continue-on-error]`
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `ci_results_${{ inputs.job }}`

2. **Set matrix**
   - ID: `set-matrix`
   - Env:
     - `job`: `${{ inputs.job }}`
     - `max_num_runners`: `${{ inputs.max_num_runners }}`

</details>

### Find commits for new failing tests (`check_new_failures`)

| Property | Value |
|----------|-------|
| Matrix | `run_idx`: ${{ fromJson(needs.setup_check_new_failures.outputs.matrix) }} |
| Depends on | `setup_check_new_failures` |
| Condition | `needs.setup_check_new_failures.outputs.process == 'true'` |

<details>
<summary>Steps (16)</summary>

1. **actions/download-artifact@v8.0.1**
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `/transformers/ci_results_${{ inputs.job }}`

2. **actions/download-artifact@v8.0.1**
   - With:
     - `pattern`: `setup_values*`
     - `path`: `setup_values`
     - `merge-multiple`: `true`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
   - Env:
     - `ACTIONS_ARTIFACT_MAX_ARTIFACT_COUNT`: `2000`

3. **Prepare some setup values**

4. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

5. **Get \`START\_SHA\`**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

6. **Extract the base commit on \`main\` (of the merge commit created by Github) if it is a PR**
   - ID: `pr_info`
   - Uses: `actions/github-script@v6.4.1`
   - Condition: `${{ inputs.pr_number != '' }}`
   - With:
     - `script`: `const pull_number = parseInt(process.env.PR_NUMBER, 10); const commit_sha = process.env.COMMIT_SHA;  const { data: pr } = await github.rest.pulls.get({   owner: context.repo.owner,   repo: context.repo.repo,   pull_number, });  const { data: merge_commit } = await github.rest.repos.getCommit({   owner: pr.base.repo.owner.login,   repo: pr.base.repo.name,   ref: commit_sha, });  core.setOutput('merge_commit_base_sha', merge_commit.parents[0].sha);`
   - Env:
     - `PR_NUMBER`: `${{ inputs.pr_number }}`
     - `COMMIT_SHA`: `${{ inputs.commit_sha }}`

7. **Get \`END\_SHA\` from previous CI runs of the same workflow**
   - Condition: `${{ inputs.pr_number == '' }}`
   - Env:
     - `ACCESS_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`

8. **Set \`END\_SHA\`**
   - Condition: `${{ inputs.pr_number != '' }}`
   - Env:
     - `merge_commit_base_sha`: `${{ steps.pr_info.outputs.merge_commit_base_sha }}`

9. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

10. **NVIDIA-SMI**

11. **Environment**

12. **Install pytest-flakefinder**

13. **Show installed libraries and their versions**

14. **Check failed tests**
   - Env:
     - `job`: `${{ inputs.job }}`
     - `n_runners`: `${{ needs.setup_check_new_failures.outputs.n_runners }}`
     - `run_idx`: `${{ matrix.run_idx }}`
     - `pr_number`: `${{ inputs.pr_number }}`

15. **Show results**
   - Env:
     - `job`: `${{ inputs.job }}`
     - `run_idx`: `${{ matrix.run_idx }}`

16. **Upload artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `new_failures_with_bad_commit_${{ inputs.job }}_${{ matrix.run_idx }}`
     - `path`: `/transformers/new_failures_with_bad_commit_${{ inputs.job }}_${{ matrix.run_idx }}.json`

</details>

### process bad commit reports (`process_new_failures_with_commit_info`)

| Property | Value |
|----------|-------|
| Depends on | `check_new_failures` |
| Condition | `needs.check_new_failures.outputs.process == 'true'` |

<details>
<summary>Steps (10)</summary>

1. **actions/download-artifact@v8.0.1**
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `/transformers/ci_results_${{ inputs.job }}`

2. **actions/download-artifact@v8.0.1**
   - With:
     - `pattern`: `new_failures_with_bad_commit_${{ inputs.job }}*`
     - `path`: `/transformers/new_failures_with_bad_commit_${{ inputs.job }}`
     - `merge-multiple`: `true`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
   - Env:
     - `ACTIONS_ARTIFACT_MAX_ARTIFACT_COUNT`: `2000`

3. **Check files**
   - Env:
     - `job`: `${{ inputs.job }}`

4. **Merge files**
   - Env:
     - `job`: `${{ inputs.job }}`

5. **Update clone**
   - Env:
     - `commit_sha`: `${{ inputs.commit_sha || github.sha }}`

6. **Process report**
   - Env:
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`: `${{ secrets.TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN }}`
     - `JOB_NAME`: `${{ inputs.job }}`
     - `REPORT_REPO_ID`: `${{ inputs.report_repo_id }}`

7. **Show results**

8. **Upload artifacts**
   - Uses: `actions/upload-artifact@v4.6.2`
   - With:
     - `name`: `new_failures_with_bad_commit_${{ inputs.job }}`
     - `path`: `/transformers/new_failures_with_bad_commit.json /transformers/new_failures_with_bad_commit_url.txt`

9. **Prepare Slack report title**
   - Env:
     - `ci_event`: `${{ inputs.ci_event }}`
     - `job`: `${{ inputs.job }}`

10. **Send processed report**
   - Uses: `slackapi/slack-github-action@6c661ce58804a1a20f6dc5fbee7f0381b469e001`
   - Condition: `${{ !endsWith(env.REPORT_TEXT, '{}') }}`
   - With:
     - `channel-id`: `#${{ inputs.slack_report_channel }}`
     - `payload`: `{   "blocks": [     {       "type": "header",       "text": {         "type": "plain_text",         "text": "${{ env.title }}"       }     },     {       "type": "section",       "text": {         "type": "mrkdwn",         "text": "${{ env.REPORT_TEXT }}"       }     }   ] }`
   - Env:
     - `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

</details>

[Back to top](#contents)

# Self-hosted runner (past-ci)

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `self-past-caller.yml` |

**Jobs:** [Model CI](#model-ci-model-ci-9), [DeepSpeed CI](#deepspeed-ci-deepspeed-ci-5)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `framework` | string | Yes | - | - |
| `version` | string | Yes | - | - |
| `sha` | string | No | `main` | - |

## Permissions

- `contents`: `read`

## Called by

`self-past-caller.yml`

- **[self-nightly-past-ci-caller.yml](#self-hosted-runner-nightly-past-ci-caller)** - entry point (x7)

## Jobs

### Model CI (`model-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_models_gpu`
- `slack_report_channel`: `#transformers-ci-past-future`
- `runner`: `past-ci`
- `docker`: `huggingface/transformers-${{ inputs.framework }}-past-${{ inputs.version }}-gpu`
- `ci_event`: `Past CI - ${{ inputs.framework }}-${{ inputs.version }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### DeepSpeed CI (`deepspeed-ci`)

| Property | Value |
|----------|-------|
| Uses workflow | [Nvidia CI (job definitions)](#nvidia-ci-job-definitions) |

#### Inputs forwarded

- `job`: `run_torch_cuda_extensions_gpu`
- `slack_report_channel`: `#transformers-ci-past-future`
- `runner`: `past-ci`
- `docker`: `huggingface/transformers-${{ inputs.framework }}-past-${{ inputs.version }}-gpu`
- `ci_event`: `Past CI - ${{ inputs.framework }}-${{ inputs.version }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

# Self-hosted runner (scheduled-intel-gaudi)

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `self-scheduled-intel-gaudi.yml` |
| Default runs-on | `group: ${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}` |

**Jobs:** [Setup](#setup-setup-5), [`run_models_gpu`](#run_models_gpu-3), [`run_trainer_and_fsdp_gpu`](#run_trainer_and_fsdp_gpu-1), [Pipelines](#pipelines-run_pipelines_torch_gpu), [Examples directory](#examples-directory-run_examples_gpu-1), [Intel Gaudi deepspeed tests](#intel-gaudi-deepspeed-tests-run_torch_cuda_extensions_gpu), [Slack Report](#slack-report-send_results-1)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `job` | string | Yes | - | - |
| `slack_report_channel` | string | Yes | - | - |
| `runner_scale_set` | string | Yes | - | - |
| `ci_event` | string | Yes | - | - |
| `report_repo_id` | string | Yes | - | - |

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `NUM_SLICES` | `2` |
| `RUN_SLOW` | `yes` |
| `PT_HPU_LAZY_MODE` | `0` |
| `TRANSFORMERS_IS_CI` | `yes` |
| `PT_ENABLE_INT64_SUPPORT` | `1` |
| `HF_TOKEN` | `${{ secrets.HF_HUB_READ_TOKEN }}` |
| `HF_HOME` | `/mnt/cache/.cache/huggingface` |

## Called by

`self-scheduled-intel-gaudi.yml`

- **[self-scheduled-intel-gaudi3-caller.yml](#self-hosted-runner-intel-gaudi3-scheduled-ci-caller)** - entry point (x5)

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `contains(fromJSON('["run_models_gpu", "run_trainer_and_fsdp_gpu"]'), inputs.job)` |

<details>
<summary>Steps (4)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@v5.6.0`
   - With:
     - `python-version`: `3.10`

3. **Identify models to test**
   - ID: `set-matrix`
   - Condition: `contains(fromJSON('["run_models_gpu", "run_trainer_and_fsdp_gpu"]'), inputs.job)`
   - Env:
     - `JOB`: `${{ inputs.job }}`

4. **Identify quantization method to test**
   - ID: `set-matrix-quantization`
   - Condition: `${{ inputs.job == 'run_quantization_torch_gpu' }}`

</details>

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs-1) |
| Matrix | `machine_type`: 1gaudi, 2gaudi; `slice_id`: ${{ fromJSON(needs.setup.outputs.slice_ids) }} |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_models_gpu' }}` |

#### Inputs forwarded

- `slice_id`: `${{ matrix.slice_id }}`
- `machine_type`: `${{ matrix.machine_type }}`
- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `runner`: `${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `run_trainer_and_fsdp_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs-1) |
| Matrix | `machine_type`: 1gaudi, 2gaudi; `slice_id`: ${{ fromJSON(needs.setup.outputs.slice_ids) }} |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_trainer_and_fsdp_gpu' }}` |

#### Inputs forwarded

- `slice_id`: `${{ matrix.slice_id }}`
- `machine_type`: `${{ matrix.machine_type }}`
- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `runner`: `${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}`
- `report_name_prefix`: `run_trainer_and_fsdp_gpu`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### Pipelines (`run_pipelines_torch_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: 1gaudi, 2gaudi |
| Condition | `${{ inputs.job == 'run_pipelines_torch_gpu' }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Install dependencies**

3. **HL-SMI**

4. **Environment**

5. **Show installed libraries and their versions**

6. **Set \`machine\_type\` for report and artifact names**

7. **Run all pipeline tests on Intel Gaudi**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_pipelines\_torch\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`

</details>

### Examples directory (`run_examples_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: 1gaudi |
| Condition | `${{ inputs.job == 'run_examples_gpu' }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Install dependencies**

3. **HL-SMI**

4. **Environment**

5. **Show installed libraries and their versions**

6. **Set \`machine\_type\` for report and artifact names**

7. **Run examples tests on Intel Gaudi**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_examples\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_examples_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_examples_gpu_test_reports`

</details>

### Intel Gaudi deepspeed tests (`run_torch_cuda_extensions_gpu`)

| Property | Value |
|----------|-------|
| Matrix | `machine_type`: 1gaudi, 2gaudi |
| Condition | `${{ inputs.job == 'run_torch_cuda_extensions_gpu' }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.3.1`
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Install dependencies**

3. **HL-SMI**

4. **Environment**

5. **Show installed libraries and their versions**

6. **Set \`machine\_type\` for report and artifact names**

7. **Run all deepspeed tests on intel Gaudi**

8. **Failure short reports** `[continue-on-error]`
   - Condition: `${{ failure() }}`

9. **Test suite reports artifacts: ${{ env.machine\_type }}\_run\_torch\_cuda\_extensions\_gpu\_test\_reports**
   - Uses: `actions/upload-artifact@v4.6.2`
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`

</details>

### Slack Report (`send_results`)

| Property | Value |
|----------|-------|
| Uses workflow | [CI slack report](#ci-slack-report) |
| Depends on | `setup`, `run_models_gpu`, `run_examples_gpu`, `run_torch_cuda_extensions_gpu`, `run_pipelines_torch_gpu`, `run_trainer_and_fsdp_gpu` |
| Condition | `${{ always() }}` |

#### Inputs forwarded

- `job`: `${{ inputs.job }}`
- `setup_status`: `${{ needs.setup.result }}`
- `slack_report_channel`: `${{ inputs.slack_report_channel }}`
- `quantization_matrix`: `${{ needs.setup.outputs.quantization_matrix }}`
- `folder_slices`: `${{ needs.setup.outputs.folder_slices }}`
- `report_repo_id`: `${{ inputs.report_repo_id }}`
- `ci_event`: `${{ inputs.ci_event }}`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

[Back to top](#contents)

