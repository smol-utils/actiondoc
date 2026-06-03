# Contents

- [Add model like runner](#add-model-like-runner)
- [Anti-Slop](#anti-slop)
- [Assign PR Reviewers](#assign-pr-reviewers)
- [Self-hosted runner (benchmark)](#self-hosted-runner-benchmark)
- [Benchmark v2 Framework](#benchmark-v2-framework)
- [Benchmark v2 Scheduled Runner - A10 Single-GPU](#benchmark-v2-scheduled-runner---a10-single-gpu)
- [Benchmark v2 Scheduled Runner - MI325 Single-GPU](#benchmark-v2-scheduled-runner---mi325-single-gpu)
- [Build pr ci-docker](#build-pr-ci-docker)
- [Build docker images (scheduled)](#build-docker-images-scheduled)
- [Build docker images (Nightly CI)](#build-docker-images-nightly-ci)
- [Build docker images (Past CI)](#build-docker-images-past-ci)
- [Build documentation](#build-documentation)
- [Build PR Documentation](#build-pr-documentation)
- [Check Permissions Advisor](#check-permissions-advisor)
- [Process failed tests](#process-failed-tests)
- [Check Tiny Models](#check-tiny-models)
- [CircleCI Failure Summary Comment](#circleci-failure-summary-comment)
- [CodeQL Security Analysis](#codeql-security-analysis)
- [CI collated reports](#ci-collated-reports)
- [Doctest job](#doctest-job)
- [Doctests](#doctests)
- [Extras Smoke Test](#extras-smoke-test)
- [Get PR commit SHA](#get-pr-commit-sha)
- [Get PR number](#get-pr-number)
- [model jobs](#model-jobs)
- [model jobs](#model-jobs-1)
- [New model PR merged notification](#new-model-pr-merged-notification)
- [PR CI](#pr-ci)
- [PR Repo. Consistency Bot](#pr-repo-consistency-bot)
- [PR - build doc via comment](#pr---build-doc-via-comment)
- [PR slow CI - Suggestion](#pr-slow-ci---suggestion)
- [Slow tests on important models (on Push - A10)](#slow-tests-on-important-models-on-push---a10)
- [Release - Conda](#release---conda)
- [Release](#release)
- [PR comment GitHub CI](#pr-comment-github-ci)
- [Nvidia CI with nightly torch](#nvidia-ci-with-nightly-torch)
- [Self-hosted runner (nightly-past-ci-caller)](#self-hosted-runner-nightly-past-ci-caller)
- [Self-hosted runner (past-ci)](#self-hosted-runner-past-ci)
- [Self-hosted runner (AMD scheduled CI caller)](#self-hosted-runner-amd-scheduled-ci-caller)
- [Self-hosted runner (AMD mi250 scheduled CI caller)](#self-hosted-runner-amd-mi250-scheduled-ci-caller)
- [Self-hosted runner scale set (AMD mi325 scheduled CI caller)](#self-hosted-runner-scale-set-amd-mi325-scheduled-ci-caller)
- [Self-hosted runner scale set (AMD mi355 scheduled CI caller)](#self-hosted-runner-scale-set-amd-mi355-scheduled-ci-caller)
- [Nvidia CI](#nvidia-ci)
- [Nvidia CI - Flash Attn](#nvidia-ci---flash-attn)
- [Self-hosted runner (scheduled-intel-gaudi)](#self-hosted-runner-scheduled-intel-gaudi)
- [Self-hosted runner (Intel Gaudi3 scheduled CI caller)](#self-hosted-runner-intel-gaudi3-scheduled-ci-caller)
- [Nvidia CI (job definitions)](#nvidia-ci-job-definitions)
- [CI slack report](#ci-slack-report)
- [SSH into our runners](#ssh-into-our-runners)
- [Stale Bot](#stale-bot)
- [TRL CI bot](#trl-ci-bot)
- [Secret Leaks](#secret-leaks)
- [Update Transformers metadata](#update-transformers-metadata)
- [Upload PR Documentation](#upload-pr-documentation)

# Add model like runner

| Property | Value |
|----------|-------|
| File | `add-model-like.yml` |
| Triggers | `push` |

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

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install dependencies**

3. **Load cached virtual environment**
   - ID: `cache`
   - Uses: `actions/cache@0057852bfaa89a56745cba8c7296529d2fc39830` (v4.3.0)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `run_all_tests_new_models_test_reports`
     - `path`: `reports/tests_new_models`

# Anti-Slop

| Property | Value |
|----------|-------|
| File | `anti-slop.yml` |
| Triggers | `pull_request_target` |

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

#### Steps

1. **peakoss/anti-slop@v0.2.1**
   - Uses: `peakoss/anti-slop@85daca1880e9e1af197fc06ea03349daf08f4202` (v0.2.1)
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

# Assign PR Reviewers

| Property | Value |
|----------|-------|
| File | `assign-reviewers.yml` |
| Triggers | `pull_request_target` |

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

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
   - With:
     - `python-version`: `3.13`

3. **Install dependencies**

4. **Run assignment script**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

# Self-hosted runner (benchmark)

| Property | Value |
|----------|-------|
| File | `benchmark.yml` |
| Triggers | `push`, `pull_request` |

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

#### Steps

1. **Get repo**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
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

# Benchmark v2 Framework

| Property | Value |
|----------|-------|
| File | `benchmark_v2.yml` |
| Triggers | `workflow_dispatch` |

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

```
benchmark_v2.yml
+-- benchmark_v2_a10_caller.yml (job: benchmark-v2-default)  <- entry point
+-- benchmark_v2_mi325_caller.yml (job: benchmark-v2-default)  <- entry point
```

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

#### Steps

1. **Get repo**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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

# Benchmark v2 Scheduled Runner - A10 Single-GPU

| Property | Value |
|----------|-------|
| File | `benchmark_v2_a10_caller.yml` |
| Triggers | `workflow_dispatch` |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
benchmark_v2_a10_caller.yml [workflow_dispatch]
+-- benchmark-v2-default (uses benchmark_v2.yml)
```

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

# Benchmark v2 Scheduled Runner - MI325 Single-GPU

| Property | Value |
|----------|-------|
| File | `benchmark_v2_mi325_caller.yml` |
| Triggers | `workflow_dispatch` |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
benchmark_v2_mi325_caller.yml [workflow_dispatch]
+-- benchmark-v2-default (uses benchmark_v2.yml)
```

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

# Build pr ci-docker

| Property | Value |
|----------|-------|
| File | `build-ci-docker-images.yml` |
| Triggers | `push`, `repository_dispatch`, `workflow_call`, `schedule` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| Runs on | `ubuntu-22.04` |
| Matrix | `file`: quality, consistency, custom-tokenizers, torch-light, exotic-models, examples-torch |
| Condition | `${{ contains(github.event.head_commit.message, '[build-ci-image]') \|\| contains(github.event.head_commit.message, '[push-ci-image]') && '!cancelled()' \|\| github.event_name == 'schedule' }}` |

#### Steps

1. **Set tag**
   - Env:
     - `COMMIT_MESSAGE`: `${{ github.event.head_commit.message }}`

2. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

3. **Check out code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

4. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

5. **Build ${{ matrix.file }}.dockerfile**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
   - With:
     - `context`: `./docker`
     - `build-args`: `REF=${{ github.sha }}`
     - `file`: `./docker/${{ matrix.file }}.dockerfile`
     - `push`: `${{ contains(github.event.head_commit.message, 'ci-image]') ||  github.event_name == 'schedule' }}`
     - `tags`: `${{ env.TAG }}`

### `notify`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `${{ contains(github.event.head_commit.message, '[build-ci-image]') \|\| contains(github.event.head_commit.message, '[push-ci-image]') && '!cancelled()' \|\| github.event_name == 'schedule' }}` |

#### Steps

1. **Post to Slack**
   - Uses: `huggingface/hf-workflows/.github/actions/post-slack@a88e7fa2eaee28de5a4d6142381b1fb792349b67`
   - Condition: `${{ contains(github.event.head_commit.message, '[push-ci-image]') && github.event_name != 'schedule' }}`
   - With:
     - `slack_channel`: `#transformers-ci-circleci-images`
     - `title`: `🤗 New docker images for CircleCI are pushed.`
     - `status`: `${{ job.status }}`
     - `slack_token`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

# Build docker images (scheduled)

| Property | Value |
|----------|-------|
| File | `build-docker-images.yml` |
| Triggers | `push`, `repository_dispatch`, `workflow_dispatch`, `workflow_call`, `schedule` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### PyTorch with Flash Attn [dev] (`flash-attn-ci-image`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### Latest PyTorch + DeepSpeed (`latest-torch-deepspeed-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### Doc builder (`doc-builder`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### Latest PyTorch (AMD) [dev] (`latest-pytorch-amd`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-highcpu-32-priv` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### Cache Latest Pytorch (AMD) Image (`cache-latest-pytorch-amd`)

| Property | Value |
|----------|-------|
| Runs on | `group: amd-mi325-1gpu` |
| Depends on | `latest-pytorch-amd` |

#### Steps

1. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

2. **Pull and save docker image to cache**

### PyTorch + DeepSpeed (AMD) [dev] (`latest-pytorch-deepspeed-amd`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

### Latest Pytorch + Quantization [dev] (`latest-quantization-torch-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@8d2750c68a42422c14e847fe6c8ac0403b4cbd6f` (v3.12.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@c94ce9fb468520275223c153574b00df6fe4bcc9` (v3.7.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@ca052bb54ab0790a636c9b5f226502c73d547a25` (v5.4.0)
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

# Build docker images (Nightly CI)

| Property | Value |
|----------|-------|
| File | `build-nightly-ci-docker-images.yml` |
| Triggers | `workflow_call`, `push` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
build-nightly-ci-docker-images.yml
+-- self-nightly-caller.yml (job: build_nightly_torch_ci_images)  <- entry point
```

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

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@885d1462b80bc1c1c7f0b00334ad271f09369c55` (v2.10.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@465a07811f14bebb1938fbed4728c6a1ff8901fc` (v2.2.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@1104d471370f9806843c095c1db02b5a90c5f8b6` (v3.3.1)
   - With:
     - `context`: `./docker/transformers-all-latest-gpu`
     - `build-args`: `REF=main PYTORCH=pre`
     - `push`: `true`
     - `tags`: `huggingface/transformers-all-latest-torch-nightly-gpu`

### Nightly PyTorch + DeepSpeed (`nightly-torch-deepspeed-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g4dn-2xlarge-cache` |
| Condition | `inputs.job == 'nightly-torch-deepspeed-docker' \|\| inputs.job == ''` |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@885d1462b80bc1c1c7f0b00334ad271f09369c55` (v2.10.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Login to DockerHub**
   - Uses: `docker/login-action@465a07811f14bebb1938fbed4728c6a1ff8901fc` (v2.2.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

4. **Build and push**
   - Uses: `docker/build-push-action@1104d471370f9806843c095c1db02b5a90c5f8b6` (v3.3.1)
   - With:
     - `context`: `./docker/transformers-pytorch-deepspeed-nightly-gpu`
     - `build-args`: `REF=main`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-deepspeed-nightly-gpu`

# Build docker images (Past CI)

| Property | Value |
|----------|-------|
| File | `build-past-ci-docker-images.yml` |
| Triggers | `push` |

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
| Runs on | `group: aws-general-8-plus` |
| Matrix | `version`: 1.13, 1.12, 1.11 |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@885d1462b80bc1c1c7f0b00334ad271f09369c55` (v2.10.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Get Base Image**
   - ID: `get-base-image`
   - Env:
     - `framework_version`: `${{ matrix.version }}`

4. **Print Base Image**

5. **Login to DockerHub**
   - Uses: `docker/login-action@465a07811f14bebb1938fbed4728c6a1ff8901fc` (v2.2.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

6. **Build and push**
   - Uses: `docker/build-push-action@1104d471370f9806843c095c1db02b5a90c5f8b6` (v3.3.1)
   - With:
     - `context`: `./docker/transformers-past-gpu`
     - `build-args`: `REF=main BASE_DOCKER_IMAGE=${{ steps.get-base-image.outputs.base_image }} FRAMEWORK=pytorch VERSION=${{ matrix.version }}`
     - `push`: `true`
     - `tags`: `huggingface/transformers-pytorch-past-${{ matrix.version }}-gpu`

### Past TensorFlow Docker (`past-tensorflow-docker`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-general-8-plus` |
| Matrix | `version`: 2.11, 2.10, 2.9, 2.8, 2.7, 2.6, 2.5 |

#### Steps

1. **Set up Docker Buildx**
   - Uses: `docker/setup-buildx-action@885d1462b80bc1c1c7f0b00334ad271f09369c55` (v2.10.0)

2. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Get Base Image**
   - ID: `get-base-image`
   - Env:
     - `framework_version`: `${{ matrix.version }}`

4. **Print Base Image**

5. **Login to DockerHub**
   - Uses: `docker/login-action@465a07811f14bebb1938fbed4728c6a1ff8901fc` (v2.2.0)
   - With:
     - `username`: `${{ secrets.DOCKERHUB_USERNAME }}`
     - `password`: `${{ secrets.DOCKERHUB_PASSWORD }}`

6. **Build and push**
   - Uses: `docker/build-push-action@1104d471370f9806843c095c1db02b5a90c5f8b6` (v3.3.1)
   - With:
     - `context`: `./docker/transformers-past-gpu`
     - `build-args`: `REF=main BASE_DOCKER_IMAGE=${{ steps.get-base-image.outputs.base_image }} FRAMEWORK=tensorflow VERSION=${{ matrix.version }}`
     - `push`: `true`
     - `tags`: `huggingface/transformers-tensorflow-past-${{ matrix.version }}-gpu`

# Build documentation

| Property | Value |
|----------|-------|
| File | `build_documentation.yml` |
| Triggers | `workflow_dispatch`, `push` |

## Event filters

- **push**
  - branches: `main`, `doc-builder*`, `v*-release`, `use_templates`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
build_documentation.yml [workflow_dispatch, push]
+-- build (uses huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391)
+-- build_other_lang (uses huggingface/doc-builder/.github/workflows/build_main_documentation.yml@2430c1ec91d04667414e2fa31ecfc36c153ea391)
```

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

# Build PR Documentation

| Property | Value |
|----------|-------|
| File | `build_pr_documentation.yml` |
| Triggers | `pull_request`, `merge_group` |

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.head_ref || github.run_id }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
build_pr_documentation.yml [pull_request, merge_group]
+-- build (uses huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@90b4ee2c10b81b5c1a6367c4e6fc9e2fb510a7e3)
```

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
| Runs on | `ubuntu-latest` |
| Condition | `github.event_name == 'merge_group'` |

#### Steps

1. **echo "Skipping doc build in merge queue"**

### `doc_build_status_check`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build`, `skip_merge_queue` |
| Condition | `always()` |

#### Steps

1. **if [[ "${{ needs.build.result }}" == "success" || "${{ ne...**

# Check Permissions Advisor

| Property | Value |
|----------|-------|
| File | `check-workflow-permissions.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `workflow_name` | string | No | - | Workflow file name |
| `run_count` | string | No | `10` | Number of runs to analyze |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
check-workflow-permissions.yml [workflow_dispatch]
+-- advisor (uses huggingface/security-workflows/.github/workflows/permissions-advisor-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c)
```

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

# Process failed tests

| Property | Value |
|----------|-------|
| File | `check_failed_tests.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
check_failed_tests.yml
+-- self-scheduled.yml (job: check_new_failures)
    +-- push-important-models.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
    +-- self-nightly-caller.yml (job: model-ci)  <- entry point
    +-- self-past-caller.yml (job: model-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-past-caller.yml (job: deepspeed-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
    +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
    +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `GITHUB_TOKEN` | job `check_new_failures` step `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` with `github-token`; job `process_new_failures_with_commit_info` step `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` with `github-token` |
| `ACCESS_REPO_INFO_TOKEN` | job `check_new_failures` step `Get `END_SHA` from previous CI runs of the same workflow` env `ACCESS_TOKEN`; job `process_new_failures_with_commit_info` step `Process report` env `ACCESS_REPO_INFO_TOKEN` |
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | job `process_new_failures_with_commit_info` step `Process report` env `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` |
| `SLACK_CIFEEDBACK_BOT_TOKEN` | job `process_new_failures_with_commit_info` step `Send processed report` env `SLACK_BOT_TOKEN` |

## Jobs

### Setup matrix for finding commits (`setup_check_new_failures`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **actions/download-artifact@v8.0.1** `[continue-on-error]`
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `ci_results_${{ inputs.job }}`

2. **Set matrix**
   - ID: `set-matrix`
   - Env:
     - `job`: `${{ inputs.job }}`
     - `max_num_runners`: `${{ inputs.max_num_runners }}`

### Find commits for new failing tests (`check_new_failures`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g5-4xlarge-cache` |
| Depends on | `setup_check_new_failures` |
| Condition | `needs.setup_check_new_failures.outputs.process == 'true'` |

#### Steps

1. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `/transformers/ci_results_${{ inputs.job }}`

2. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
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
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `new_failures_with_bad_commit_${{ inputs.job }}_${{ matrix.run_idx }}`
     - `path`: `/transformers/new_failures_with_bad_commit_${{ inputs.job }}_${{ matrix.run_idx }}.json`

### process bad commit reports (`process_new_failures_with_commit_info`)

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g5-4xlarge-cache` |
| Depends on | `check_new_failures` |
| Condition | `needs.check_new_failures.outputs.process == 'true'` |

#### Steps

1. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `/transformers/ci_results_${{ inputs.job }}`

2. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
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

# Check Tiny Models

| Property | Value |
|----------|-------|
| File | `check_tiny_models.yml` |
| Triggers | `push`, `repository_dispatch`, `schedule` |

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

#### Steps

1. **Checkout transformers**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `fetch-depth`: `2`
     - `persist-credentials`: `false`

2. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

3. **Set up Python 3.10**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
   - With:
     - `python-version`: `3.10`
     - `architecture`: `x64`

4. **Install**

5. **Create all tiny models (locally)**

6. **Local tiny model reports artifacts**
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `tiny_local_model_creation_reports`
     - `path`: `tiny_local_models/reports`

# CircleCI Failure Summary Comment

| Property | Value |
|----------|-------|
| File | `circleci-failure-summary-comment.yml` |
| Triggers | `pull_request_target` |

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

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Setup Python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
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
   - Uses: `actions/github-script@f28e40c7f34bde8b3046d885e986cb6290c5673b` (v7.1.0)
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

# CodeQL Security Analysis

| Property | Value |
|----------|-------|
| File | `codeql.yml` |
| Triggers | `push`, `workflow_dispatch` |

## Event filters

- **push**
  - branches: `main`, `fix_security_issue_*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
codeql.yml [push, workflow_dispatch]
+-- codeql (uses huggingface/security-workflows/.github/workflows/codeql-reusable.yml@1b6a139c28db347498b30338da6a602e0a06f56c)
```

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

# CI collated reports

| Property | Value |
|----------|-------|
| File | `collated-reports.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
collated-reports.yml
+-- model_jobs.yml (job: collated_reports)
    +-- self-scheduled.yml (job: run_models_gpu)
    |   +-- push-important-models.yml (job: model-ci)  <- entry point
    |   +-- self-comment-ci.yml (job: model-ci)  <- entry point
    |   +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
    |   +-- self-nightly-caller.yml (job: model-ci)  <- entry point
    |   +-- self-past-caller.yml (job: model-ci)
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    |   +-- self-past-caller.yml (job: deepspeed-ci)
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    |   +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
    |   +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
    |   +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
    |   +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
    |   +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
    |   +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
    |   +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
    |   +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
    +-- self-scheduled.yml (job: run_trainer_and_fsdp_gpu)
        +-- push-important-models.yml (job: model-ci)  <- entry point
        +-- self-comment-ci.yml (job: model-ci)  <- entry point
        +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
        +-- self-nightly-caller.yml (job: model-ci)  <- entry point
        +-- self-past-caller.yml (job: model-ci)
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
        +-- self-past-caller.yml (job: deepspeed-ci)
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
        |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
        +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
        +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
        +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
        +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
        +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
        +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
        +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
        +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
```

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

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **actions/download-artifact@v4.3.0**
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)

3. **Collated reports**
   - Env:
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `CI_SHA`: `${{ github.sha }}`
     - `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`: `${{ secrets.TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN }}`
     - `MACHINE_TYPE`: `${{ inputs.machine_type }}`
     - `JOB`: `${{ inputs.job }}`
     - `REPORT_REPO_ID`: `${{ inputs.report_repo_id }}`
     - `GPU_NAME`: `${{ inputs.gpu_name }}`

# Doctest job

| Property | Value |
|----------|-------|
| File | `doctest_job.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
doctest_job.yml
+-- doctests.yml (job: call_doctest_job)  <- entry point
```

## Jobs

### `run_doctests`

| Property | Value |
|----------|-------|
| Runs on | `group: aws-g5-4xlarge-cache` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `doc_tests_gpu_test_reports_${{ env.split_keys }}`
     - `path`: `/transformers/reports/doc_tests_gpu_${{ env.split_keys }}`

# Doctests

| Property | Value |
|----------|-------|
| File | `doctests.yml` |
| Triggers | `push`, `repository_dispatch`, `schedule` |

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

```
doctests.yml [push, repository_dispatch, schedule]
+-- call_doctest_job (uses doctest_job.yml)
```

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

#### Steps

1. **Update clone**

2. **Reinstall transformers in edit mode (remove the one installed during docker image build)**

3. **Show installed libraries and their versions**

4. **Check values for matrix**

5. **Set values for matrix**
   - ID: `set-matrix`

### Call doctest jobs (`call_doctest_job`)

| Property | Value |
|----------|-------|
| Uses workflow | [Doctest job](#doctest-job) |
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

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **actions/download-artifact@v4.3.0**
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)

3. **Send message to Slack**
   - Env:
     - `CI_SLACK_BOT_TOKEN`: `${{ secrets.CI_SLACK_BOT_TOKEN }}`
     - `ACCESS_REPO_INFO_TOKEN`: `${{ secrets.ACCESS_REPO_INFO_TOKEN }}`
     - `SLACK_REPORT_CHANNEL`: `${{ secrets.CI_SLACK_CHANNEL_ID_DAILY_DOCS }}`

4. **Upload results**
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `doc_test_results`
     - `path`: `doc_test_results`

# Extras Smoke Test

| Property | Value |
|----------|-------|
| File | `extras-smoke-test.yml` |
| Triggers | `schedule` |

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

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install setuptools**

3. **Extract Python versions from setup.py**
   - ID: `extract-versions`

### Test extras on Python ${{ matrix.python-version }} (`test-extras`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `get-python-versions` |

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Set up Python ${{ matrix.python-version }}**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `always()`
   - With:
     - `name`: `failure-report-${{ matrix.python-version }}`
     - `path`: `failure_reports/`
     - `retention-days`: `1`
     - `if-no-files-found`: `ignore`

### Check Slack token availability (`precheck-slack`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **chk**
   - ID: `chk`
   - Env:
     - `SLACK_BOT_TOKEN`: `${{ secrets.SLACK_CIFEEDBACK_BOT_TOKEN }}`

### Notify failures to Slack (`notify-failures`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `test-extras`, `precheck-slack` |
| Condition | `always() && needs.precheck-slack.outputs.has_slack_token == 'true' && needs.test-extras.result != 'success'` |

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
   - With:
     - `python-version`: `3.11`

3. **Download all failure reports** `[continue-on-error]`
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)
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

# Get PR commit SHA

| Property | Value |
|----------|-------|
| File | `get-pr-info.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
get-pr-info.yml
+-- pr-repo-consistency-bot.yml (job: get-pr-info)  <- entry point
+-- pr_build_doc_with_comment.yml (job: get-pr-info)  <- entry point
+-- pr_slow_ci_suggestion.yml (job: get-pr-info)  <- entry point
+-- self-comment-ci.yml (job: get-pr-info)  <- entry point
```

## Jobs

### Get PR commit SHA better (`get-pr-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Condition | `${{ inputs.pr_number != '' }}` |

#### Steps

1. **Extract PR details**
   - ID: `pr_info`
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - With:
     - `script`: `const pull_number = parseInt(process.env.PR_NUMBER, 10);  const { data: pr } = await github.rest.pulls.get({   owner: context.repo.owner,   repo: context.repo.repo,   pull_number, });  const { data: head_commit } = await github.rest.repos.getCommit({   owner: pr.head.repo.owner.login,   repo: pr.head.repo.name,   ref: pr.head.ref });  const { data: merge_commit } = await github.rest.repos.getCommit({   owner: pr.base.repo.owner.login,   repo: pr.base.repo.name,   ref: pr.merge_commit_sha, });  const { data: files } = await github.rest.pulls.listFiles({   owner: context.repo.owner,   repo: context.repo.repo,   pull_number, });  core.setOutput('head_repo_full_name', pr.head.repo.full_name); core.setOutput('base_repo_full_name', pr.base.repo.full_name); core.setOutput('head_repo_owner', pr.head.repo.owner.login); core.setOutput('base_repo_owner', pr.base.repo.owner.login); core.setOutput('head_repo_name', pr.head.repo.name); core.setOutput('base_repo_name', pr.base.repo.name); core.setOutput('head_ref', pr.head.ref); core.setOutput('base_ref', pr.base.ref); core.setOutput('head_sha', pr.head.sha); core.setOutput('base_sha', pr.base.sha); core.setOutput('merge_commit_base_sha', merge_commit.parents[0].sha); core.setOutput('merge_commit_sha', pr.merge_commit_sha); core.setOutput('pr', pr);  core.setOutput('head_commit_date', head_commit.commit.committer.date); core.setOutput('merge_commit_date', merge_commit.commit.committer.date);  core.setOutput('files', files);              console.log('PR head commit:', {   head_commit: head_commit,   commit: head_commit.commit,   date: head_commit.commit.committer.date });  console.log('PR merge commit:', {   merge_commit: merge_commit,   commit: merge_commit.commit,   date: merge_commit.commit.committer.date });  console.log('PR Info:', {   pr_info: pr });`
   - Env:
     - `PR_NUMBER`: `${{ inputs.pr_number }}`

2. **Convert dates to timestamps**
   - ID: `get_timestamps`
   - Env:
     - `head_commit_date`: `${{ steps.pr_info.outputs.head_commit_date }}`
     - `merge_commit_date`: `${{ steps.pr_info.outputs.merge_commit_date }}`

# Get PR number

| Property | Value |
|----------|-------|
| File | `get-pr-number.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `PR_NUMBER` | The extracted PR number | `${{ jobs.get-pr-number.outputs.PR_NUMBER }}` |

## Permissions

- `contents`: `read`

## Called by

```
get-pr-number.yml
+-- pr-repo-consistency-bot.yml (job: get-pr-number)  <- entry point
+-- pr_build_doc_with_comment.yml (job: get-pr-number)  <- entry point
+-- pr_slow_ci_suggestion.yml (job: get-pr-number)  <- entry point
+-- self-comment-ci.yml (job: get-pr-number)  <- entry point
```

## Jobs

### Get PR number (`get-pr-number`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

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

# model jobs

| Property | Value |
|----------|-------|
| File | `model_jobs.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
model_jobs.yml
+-- self-scheduled.yml (job: run_models_gpu)
|   +-- push-important-models.yml (job: model-ci)  <- entry point
|   +-- self-comment-ci.yml (job: model-ci)  <- entry point
|   +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
|   +-- self-nightly-caller.yml (job: model-ci)  <- entry point
|   +-- self-past-caller.yml (job: model-ci)
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
|   +-- self-past-caller.yml (job: deepspeed-ci)
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
|   |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
|   +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
|   +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
|   +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
|   +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
|   +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
|   +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
|   +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
|   +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
+-- self-scheduled.yml (job: run_trainer_and_fsdp_gpu)
    +-- push-important-models.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
    +-- self-nightly-caller.yml (job: model-ci)  <- entry point
    +-- self-past-caller.yml (job: model-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-past-caller.yml (job: deepspeed-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
    +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
    +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
```

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

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`

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

# model jobs

| Property | Value |
|----------|-------|
| File | `model_jobs_intel_gaudi.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
model_jobs_intel_gaudi.yml
+-- self-scheduled-intel-gaudi.yml (job: run_models_gpu)
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: model-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: pipeline-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: example-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: deepspeed-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: trainer-fsdp-ci)  <- entry point
+-- self-scheduled-intel-gaudi.yml (job: run_trainer_and_fsdp_gpu)
    +-- self-scheduled-intel-gaudi3-caller.yml (job: model-ci)  <- entry point
    +-- self-scheduled-intel-gaudi3-caller.yml (job: pipeline-ci)  <- entry point
    +-- self-scheduled-intel-gaudi3-caller.yml (job: example-ci)  <- entry point
    +-- self-scheduled-intel-gaudi3-caller.yml (job: deepspeed-ci)  <- entry point
    +-- self-scheduled-intel-gaudi3-caller.yml (job: trainer-fsdp-ci)  <- entry point
```

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

#### Steps

1. **Echo input and matrix info**
   - Env:
     - `FOLDER_SLICES`: `${{ inputs.folder_slices }}`
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`
     - `SLICE`: `${{ toJson(fromJson(inputs.folder_slices)[inputs.slice_id]) }}`

2. **Echo folder ${{ matrix.folders }}**
   - Env:
     - `MATRIX_FOLDERS`: `${{ matrix.folders }}`

3. **Checkout**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ env.matrix_folders }}_test_reports`
     - `path`: `reports/${{ env.machine_type }}_${{ inputs.report_name_prefix }}_${{ matrix.folders }}_test_reports`

# New model PR merged notification

Used to notify core maintainers about new model PR being merged

| Property | Value |
|----------|-------|
| File | `new_model_pr_merged_notification.yml` |
| Triggers | `push` |

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

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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

# PR CI

| Property | Value |
|----------|-------|
| File | `pr-ci-caller.yml` |
| Triggers | `pull_request` |

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.event.pull_request.number }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
pr-ci-caller.yml [pull_request]
+-- pr-ci (uses huggingface/transformers-test-ci/.github/workflows/pr-ci_dynamic_caller_example.yml@91d590c4f744e4564a8ae0d3810068c8a35b939e)
```

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

# PR Repo. Consistency Bot

| Property | Value |
|----------|-------|
| File | `pr-repo-consistency-bot.yml` |
| Triggers | `issue_comment` |

## Event filters

- **issue_comment**
  - types: `created`
  - branches-ignore: `main`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.event.issue.number }}-${{ startsWith(github.event.comment.body, '@bot /repo') || startsWith(github.event.comment.body, '@bot /style') }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
pr-repo-consistency-bot.yml [issue_comment]
+-- get-pr-number (uses get-pr-number.yml)
+-- get-pr-info (uses get-pr-info.yml)
```

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
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-info` |

#### Steps

1. **Verify \`merge\_commit\` timestamp is older than the issue comment timestamp**
   - Env:
     - `COMMENT_DATE`: `${{ github.event.comment.created_at }}`
     - `PR_MERGE_COMMIT_TIMESTAMP`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}`

### Init Comment on PR (`init_comment_with_url`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `check-timestamps` |

**Permissions:**

- `pull-requests`: `write`

#### Steps

1. **Delete existing bot comment if it exists**
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - With:
     - `script`: `` const PR_NUMBER = parseInt(process.env.PR_NUMBER, 10);  // Get all comments on the PR const { data: comments } = await github.rest.issues.listComments({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: PR_NUMBER });  // Find existing bot comments that start with "Repo. Consistency" or "Style fix" const existingComments = comments.filter(comment =>    comment.user.login === 'github-actions[bot]' &&    (comment.body.startsWith('Repo. Consistency') || comment.body.startsWith('Style fix')) );  if (existingComments.length > 0) {   // Get the most recent comment   const mostRecentComment = existingComments     .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))[0];      console.log(`Deleting most recent comment #${mostRecentComment.id}`);   await github.rest.issues.deleteComment({     owner: context.repo.owner,     repo: context.repo.repo,     comment_id: mostRecentComment.id   }); } ``
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

2. **Comment on PR with workflow run link**
   - ID: `init_comment`
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - With:
     - `script`: `` const PR_NUMBER = parseInt(process.env.PR_NUMBER, 10); const COMMENT_BODY = process.env.COMMENT_BODY; const runUrl = `${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`  // Determine which command was used const isStyleFix = COMMENT_BODY.startsWith('@bot /style'); const messagePrefix = isStyleFix ? 'Style fix' : 'Repo. Consistency fix';  const { data: botComment } = await github.rest.issues.createComment({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: PR_NUMBER,   body: `${messagePrefix} is beginning .... [View the workflow run here](${runUrl}).` }); core.setOutput('comment_id', botComment.id); ``
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
     - `COMMENT_BODY`: `${{ github.event.comment.body }}`

### `run-repo-consistency-checks`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-info`, `check-timestamps`, `init_comment_with_url` |

#### Steps

1. **Checkout base repository**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `ref`: `main`
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@7f4fc3e22c37d6ff65e88745f38bd3157c663f7c` (v4.9.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `steps.run_repo_checks.outputs.changes_detected == 'true' || steps.run_style_checks.outputs.changes_detected == 'true'`
   - With:
     - `name`: `modified-files`
     - `path`: `artifact-staging/`

### `commit-and-comment`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-pr-info`, `check-timestamps`, `init_comment_with_url`, `run-repo-consistency-checks` |
| Condition | `always()` |

**Permissions:**

- `pull-requests`: `write`
- `contents`: `write`

#### Steps

1. **Download modified files**
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)
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
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - Condition: `needs.init_comment_with_url.result == 'success'`
   - With:
     - `script`: `const pr_number = parseInt(process.env.PR_NUMBER, 10); const comment_id = parseInt(process.env.COMMENT_ID, 10); const body = process.env.FINAL_COMMENT; await github.rest.issues.updateComment({   owner: context.repo.owner,   repo: context.repo.repo,   comment_id,   body, });`
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`
     - `COMMENT_ID`: `${{ needs.init_comment_with_url.outputs.comment_id }}`
     - `FINAL_COMMENT`: `${{ steps.prepare_final_comment.outputs.final_comment }}`

# PR - build doc via comment

| Property | Value |
|----------|-------|
| File | `pr_build_doc_with_comment.yml` |
| Triggers | `issue_comment` |

## Event filters

- **issue_comment**
  - types: `created`
  - branches-ignore: `main`

## Permissions

No permissions granted (`permissions: {}` -- default-deny).

**Concurrency:** group `${{ github.workflow }}-${{ github.event.issue.number }}-${{ startsWith(github.event.comment.body, 'build-doc') }}`, cancel-in-progress: `true`

## Call graph (rooted at this workflow)

```
pr_build_doc_with_comment.yml [issue_comment]
+-- get-pr-number (uses get-pr-number.yml)
+-- get-pr-info (uses get-pr-info.yml)
+-- build-doc (uses huggingface/doc-builder/.github/workflows/build_pr_documentation.yml@093eb65f2e8745457987df060dc392e6bcf1347a)
```

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
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-info` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != ''}}` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `COMMENT_DATE` | `${{ github.event.comment.created_at }}` |
| `PR_MERGE_COMMIT_DATE` | `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_DATE }}` |
| `PR_MERGE_COMMIT_TIMESTAMP` | `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}` |

#### Steps

1. **COMMENT\_TIMESTAMP=$(date -d "${COMMENT\_DATE}" +"%s")**

### Create run (`create_run`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-pr-info` |
| Condition | `${{ needs.get-pr-number.outputs.PR_NUMBER != '' }}` |

**Permissions:**

- `statuses`: `write`

#### Steps

1. **Create Run**
   - ID: `create_run`
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `NEEDS_GET_PR_INFO_OUTPUTS_PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

### Reply to the comment (`reply_to_comment`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `create_run` |
| Condition | `${{ needs.create_run.result == 'success' }}` |

**Permissions:**

- `pull-requests`: `write`

#### Steps

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `NEEDS_GET_PR_NUMBER_OUTPUTS_PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

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
| Runs on | `ubuntu-22.04` |
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

#### Steps

1. **Get \`build-doc\` job status**

2. **Update PR commit statuses**
   - Env:
     - `NEEDS_GET_PR_INFO_OUTPUTS_PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

# PR slow CI - Suggestion

| Property | Value |
|----------|-------|
| File | `pr_slow_ci_suggestion.yml` |
| Triggers | `pull_request_target` |

## Event filters

- **pull_request_target**
  - types: `opened`, `synchronize`, `reopened`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
pr_slow_ci_suggestion.yml [pull_request_target]
+-- get-pr-number (uses get-pr-number.yml)
+-- get-pr-info (uses get-pr-info.yml)
```

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
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-pr-info` |

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Write pr\_files file**
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - With:
     - `script`: `const fs = require('node:fs'); const files = await github.paginate(github.rest.pulls.listFiles, {   owner: context.repo.owner,   repo: context.repo.repo,   pull_number: parseInt(process.env.PR_NUMBER, 10), }); fs.writeFileSync('pr_files.txt', JSON.stringify(files));`
   - Env:
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

3. **Get repository content**
   - ID: `repo_content`
   - Uses: `actions/github-script@d7906e4ad0b1822421a7e6a35d5ca353c962f410` (v6.4.1)
   - With:
     - `script`: `const fs = require('node:fs'); const { PR_HEAD_REPO_OWNER, PR_HEAD_REPO_NAME, PR_HEAD_SHA } = process.env;  const { data: tests_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests',   ref: PR_HEAD_SHA, });  const { data: tests_models_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests/models',   ref: PR_HEAD_SHA, });  const { data: tests_quantization_dir } = await github.rest.repos.getContent({   owner: PR_HEAD_REPO_OWNER,   repo: PR_HEAD_REPO_NAME,   path: 'tests/quantization',   ref: PR_HEAD_SHA, });  // Write to files instead of outputs fs.writeFileSync('tests_dir.txt', JSON.stringify(tests_dir, null, 2)); fs.writeFileSync('tests_models_dir.txt', JSON.stringify(tests_models_dir, null, 2)); fs.writeFileSync('tests_quantization_dir.txt', JSON.stringify(tests_quantization_dir, null, 2));`
   - Env:
     - `PR_HEAD_REPO_OWNER`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_OWNER }}`
     - `PR_HEAD_REPO_NAME`: `${{ needs.get-pr-info.outputs.PR_HEAD_REPO_NAME }}`
     - `PR_HEAD_SHA`: `${{ needs.get-pr-info.outputs.PR_HEAD_SHA }}`

4. **Run script to get jobs to run**
   - ID: `get_jobs`

### Send a comment to suggest jobs to run (`send_comment`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-jobs` |
| Condition | `${{ needs.get-jobs.outputs.jobs != '' }}` |

**Permissions:**

- `pull-requests`: `write`

#### Steps

1. **Check and update comment if needed**
   - Uses: `actions/github-script@f28e40c7f34bde8b3046d885e986cb6290c5673b` (v7.1.0)
   - With:
     - `script`: `` const prNumber = parseInt(process.env.PR_NUMBER, 10); const commentPrefix = "**[For maintainers]** Suggested jobs to run (before merge)"; const thirtyMinutesAgo = new Date(Date.now() - 30 * 60 * 1000); // 30 minutes ago const newBody = `${commentPrefix}${process.env.BODY}`;  // Get all comments on the PR const { data: comments } = await github.rest.issues.listComments({   owner: context.repo.owner,   repo: context.repo.repo,   issue_number: prNumber });  // Find existing comments that start with our prefix const existingComments = comments.filter(comment =>   comment.user.login === 'github-actions[bot]' &&   comment.body.startsWith(commentPrefix) );  let shouldCreateNewComment = true; let commentsToDelete = [];  if (existingComments.length > 0) {   // Get the most recent comment   const mostRecentComment = existingComments     .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))[0];    const commentDate = new Date(mostRecentComment.created_at);   const isOld = commentDate < thirtyMinutesAgo;   const isDifferentContent = mostRecentComment.body !== newBody;    console.log(`Most recent comment created: ${mostRecentComment.created_at}`);   console.log(`Is older than 30 minutes: ${isOld}`);   console.log(`Has different content: ${isDifferentContent}`);    if (isOld || isDifferentContent) {     // Delete all existing comments and create new one     commentsToDelete = existingComments;     console.log(`Will delete ${commentsToDelete.length} existing comment(s) and create new one`);   } else {     // Content is same and comment is recent, skip     shouldCreateNewComment = false;     console.log('Comment is recent and content unchanged, skipping update');   } } else {   console.log('No existing comments found, will create new one'); }  // Delete old comments if needed for (const comment of commentsToDelete) {   console.log(`Deleting comment #${comment.id} (created: ${comment.created_at})`);   await github.rest.issues.deleteComment({     owner: context.repo.owner,     repo: context.repo.repo,     comment_id: comment.id   }); }  // Create new comment if needed if (shouldCreateNewComment) {   await github.rest.issues.createComment({     owner: context.repo.owner,     repo: context.repo.repo,     issue_number: prNumber,     body: newBody   });   console.log('✅ New comment created'); } else {   console.log('ℹ️ No comment update needed'); } ``
   - Env:
     - `BODY`: `run-slow: ${{ needs.get-jobs.outputs.jobs }}`
     - `PR_NUMBER`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

# Slow tests on important models (on Push - A10)

| Property | Value |
|----------|-------|
| File | `push-important-models.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
push-important-models.yml [push]
+-- model-ci (uses self-scheduled.yml)
    +-- run_models_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- send_results (uses slack-report.yml)
    +-- check_new_failures (uses check_failed_tests.yml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Get all modified files (`get_modified_models`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Check out code**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **Get changed files using \`actions/github-script\`**
   - ID: `get-changed-files`
   - Uses: `actions/github-script@f28e40c7f34bde8b3046d885e986cb6290c5673b` (v7.1.0)
   - With:
     - `script`: `` let files = [];  // Only handle push events if (context.eventName === 'push') {   const afterSha = context.payload.after;   const branchName = context.payload.ref.replace('refs/heads/', '');      let baseSha;      if (branchName === 'main') {     console.log('Push to main branch, comparing to parent commit');     // Get the parent commit of the pushed commit     const { data: commit } = await github.rest.repos.getCommit({       owner: context.repo.owner,       repo: context.repo.repo,       ref: afterSha     });     baseSha = commit.parents[0]?.sha;     if (!baseSha) {       throw new Error('No parent commit found for the pushed commit');     }   } else {     console.log(`Push to branch ${branchName}, comparing to main`);     baseSha = 'main';   }      const { data: comparison } = await github.rest.repos.compareCommits({     owner: context.repo.owner,     repo: context.repo.repo,     base: baseSha,     head: afterSha   });      // Include added, modified, and renamed files   files = comparison.files     .filter(file => file.status === 'added' || file.status === 'modified' || file.status === 'renamed')     .map(file => file.filename); }  // Include all files under src/transformers/ (not just models subdirectory) const filteredFiles = files.filter(file =>    file.startsWith('src/transformers/') );  core.setOutput('changed_files', filteredFiles.join(' ')); core.setOutput('any_changed', filteredFiles.length > 0 ? 'true' : 'false'); ``

3. **Parse changed files with Python**
   - ID: `set-matrix`
   - Condition: `steps.get-changed-files.outputs.any_changed == 'true'`
   - Env:
     - `CHANGED_FILES`: `${{ steps.get-changed-files.outputs.changed_files }}`

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

# Release - Conda

| Property | Value |
|----------|-------|
| File | `release-conda.yml` |
| Triggers | `push` |

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

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Install miniconda**
   - Uses: `conda-incubator/setup-miniconda@9f54435e0e72c53962ee863144e47a4b094bfd35` (v2.3.0)
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

# Release

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - tags: `v*`
  - branches: `v*-release`

## Permissions

- `contents`: `read`

## Jobs

### build release (`build_and_test`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **set up python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `python-dist`
     - `path`: `dist/** build/**`

### `upload_package`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build_and_test` |
| Condition | `startsWith(github.ref, 'refs/tags/')` |

**Deploys to environment:** `pypi-release` [gated]

> Environment protection rules (required reviewers, wait timers, branch policies) are configured in the repository's Settings -> Environments and are not represented here.

**Permissions:**

- `id-token`: `write` (OIDC)

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Download build artifacts**
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)
   - With:
     - `name`: `python-dist`
     - `path`: `.`

3. **Publish package distributions to TestPyPI**
   - Uses: `pypa/gh-action-pypi-publish@ed0c53931b1dc9bd32cbe73a98c7f6766f8a527e`
   - With:
     - `verbose`: `true`

# PR comment GitHub CI

| Property | Value |
|----------|-------|
| File | `self-comment-ci.yml` |
| Triggers | `issue_comment` |

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

```
self-comment-ci.yml [issue_comment]
+-- get-pr-number (uses get-pr-number.yml)
+-- get-pr-info (uses get-pr-info.yml)
+-- model-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- quantization-ci (uses self-scheduled.yml)
    +-- run_models_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- send_results (uses slack-report.yml)
    +-- check_new_failures (uses check_failed_tests.yml)
```

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
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-info` |

#### Steps

1. **Verify \`merge\_commit\` timestamp is older than the issue comment timestamp**
   - Env:
     - `COMMENT_DATE`: `${{ github.event.comment.created_at }}`
     - `PR_MERGE_COMMIT_TIMESTAMP`: `${{ needs.get-pr-info.outputs.PR_MERGE_COMMIT_TIMESTAMP }}`

### `get-tests`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `check-timestamps` |

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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

### Report error earlier (`report_error_earlier`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-pr-info`, `get-tests` |
| Condition | `${{ always() && needs.get-pr-info.result == 'success' && needs.get-tests.result != 'success' }}` |

**Permissions:**

- `pull-requests`: `write`

#### Steps

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `PREFIX`: `[Workflow Run ⚙️](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})\n\n`
     - `github_repository`: `${{ github.repository }}`
     - `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Reply to the comment (`reply_to_comment`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-tests` |
| Condition | `${{ needs.get-tests.outputs.models != '[]'  \|\| needs.get-tests.outputs.quantizations != '[]' }}` |

**Permissions:**

- `pull-requests`: `write`

#### Steps

1. **Reply to the comment**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `PREFIX`: `[Workflow Run ⚙️](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})`
     - `INFO`: `` \n\nThis comment contains `run-slow`, running the specified jobs ``
     - `BODY`: `\n\nmodels: ${{ needs.get-tests.outputs.models }}\nquantizations: ${{ needs.get-tests.outputs.quantizations }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_number`: `${{ needs.get-pr-number.outputs.PR_NUMBER }}`

### Create run (`create_run`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `check-timestamps`, `reply_to_comment` |

**Permissions:**

- `statuses`: `write`

#### Steps

1. **Create Run**
   - ID: `create_run`
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`
     - `GITHUB_RUN_URL`: `https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}`
     - `github_repository`: `${{ github.repository }}`
     - `pr_head_sha`: `${{ needs.check-timestamps.outputs.PR_HEAD_SHA }}`

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
| Runs on | `ubuntu-22.04` |
| Depends on | `get-pr-number`, `get-pr-info`, `check-timestamps`, `create_run`, `model-ci`, `quantization-ci` |
| Condition | `${{ always() && needs.create_run.result == 'success' }}` |

**Permissions:**

- `pull-requests`: `write`
- `statuses`: `write`

#### Steps

1. **actions/download-artifact@v4.3.0**
   - Uses: `actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093` (v4.3.0)
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

# Nvidia CI with nightly torch

| Property | Value |
|----------|-------|
| File | `self-nightly-caller.yml` |
| Triggers | `repository_dispatch`, `workflow_run`, `push` |

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

```
self-nightly-caller.yml [repository_dispatch, workflow_run, push]
+-- build_nightly_torch_ci_images (uses build-nightly-ci-docker-images.yml)
+-- model-ci (uses self-scheduled.yml)
    +-- run_models_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- send_results (uses slack-report.yml)
    +-- check_new_failures (uses check_failed_tests.yml)
```

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

#### Steps

1. **Setup**
   - Env:
     - `PREV_WORKFLOW_RUN_ID`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `OTHER_WORKFLOW_RUN_ID`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

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

# Self-hosted runner (nightly-past-ci-caller)

| Property | Value |
|----------|-------|
| File | `self-nightly-past-ci-caller.yml` |
| Triggers | `schedule`, `push` |

## Schedule

- `17 2,14 * * *`

## Event filters

- **push**
  - branches: `run_past_ci*`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
self-nightly-past-ci-caller.yml [schedule, push]
+-- run_past_ci_tensorflow_2-11 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-10 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-9 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-8 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-7 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-6 (uses self-past-caller.yml)
|   +-- model-ci (uses self-scheduled.yml)
|   |   +-- run_models_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   |   +-- send_results (uses slack-report.yml)
|   |   +-- check_new_failures (uses check_failed_tests.yml)
|   +-- deepspeed-ci (uses self-scheduled.yml)
|       +-- run_models_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|       |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|       +-- send_results (uses slack-report.yml)
|       +-- check_new_failures (uses check_failed_tests.yml)
+-- run_past_ci_tensorflow_2-5 (uses self-past-caller.yml)
    +-- model-ci (uses self-scheduled.yml)
    |   +-- run_models_gpu (uses model_jobs.yml)
    |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    |   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    |   +-- send_results (uses slack-report.yml)
    |   +-- check_new_failures (uses check_failed_tests.yml)
    +-- deepspeed-ci (uses self-scheduled.yml)
        +-- run_models_gpu (uses model_jobs.yml)
        |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
        +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
        |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
        +-- send_results (uses slack-report.yml)
        +-- check_new_failures (uses check_failed_tests.yml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Get number (`get_number`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Get number**
   - ID: `get_number`

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

# Self-hosted runner (past-ci)

| Property | Value |
|----------|-------|
| File | `self-past-caller.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `framework` | string | Yes | - | - |
| `version` | string | Yes | - | - |
| `sha` | string | No | `main` | - |

## Permissions

- `contents`: `read`

## Called by

```
self-past-caller.yml
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
+-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
```

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

# Self-hosted runner (AMD scheduled CI caller)

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-caller.yml` |
| Triggers | `schedule` |

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

#### Steps

1. **Trigger scheduled AMD CI via workflow\_run**

# Self-hosted runner (AMD mi250 scheduled CI caller)

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi250-caller.yml` |
| Triggers | `workflow_run`, `push` |

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

```
self-scheduled-amd-mi250-caller.yml [workflow_run, push]
+-- model-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- torch-pipeline (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- example-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- deepspeed-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
```

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

# Self-hosted runner scale set (AMD mi325 scheduled CI caller)

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi325-caller.yml` |
| Triggers | `workflow_run`, `push` |

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

```
self-scheduled-amd-mi325-caller.yml [workflow_run, push]
+-- model-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- torch-pipeline (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- example-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- deepspeed-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
```

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

# Self-hosted runner scale set (AMD mi355 scheduled CI caller)

| Property | Value |
|----------|-------|
| File | `self-scheduled-amd-mi355-caller.yml` |
| Triggers | `workflow_run`, `push` |

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

```
self-scheduled-amd-mi355-caller.yml [workflow_run, push]
+-- model-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- torch-pipeline (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- example-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
+-- deepspeed-ci (uses huggingface/hf-workflows/.github/workflows/transformers_amd_ci_scheduled_arc_scale_set.yaml@63657f571a92cc9759159442936061c51d6d9ae4)
```

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

# Nvidia CI

| Property | Value |
|----------|-------|
| File | `self-scheduled-caller.yml` |
| Triggers | `repository_dispatch`, `schedule`, `push`, `workflow_dispatch` |

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

```
self-scheduled-caller.yml [repository_dispatch, schedule, push, workflow_dispatch]
+-- model-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- torch-pipeline (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- example-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- trainer-fsdp-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- deepspeed-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- quantization-ci (uses self-scheduled.yml)
|   +-- run_models_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
|   |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
|   +-- send_results (uses slack-report.yml)
|   +-- check_new_failures (uses check_failed_tests.yml)
+-- kernels-ci (uses self-scheduled.yml)
    +-- run_models_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- send_results (uses slack-report.yml)
    +-- check_new_failures (uses check_failed_tests.yml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Setup**
   - Env:
     - `prev_workflow_run_id`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `other_workflow_run_id`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

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

# Nvidia CI - Flash Attn

| Property | Value |
|----------|-------|
| File | `self-scheduled-flash-attn-caller.yml` |
| Triggers | `repository_dispatch`, `schedule`, `push`, `workflow_dispatch` |

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

```
self-scheduled-flash-attn-caller.yml [repository_dispatch, schedule, push, workflow_dispatch]
+-- model-ci (uses self-scheduled.yml)
    +-- run_models_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs.yml)
    |   +-- collated_reports (uses collated-reports.yml@6abd9725ee7d809dc974991f8ff6c958afb63a3a)
    +-- send_results (uses slack-report.yml)
    +-- check_new_failures (uses check_failed_tests.yml)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `ACCESS_REPO_INFO_TOKEN`, `CI_SLACK_BOT_TOKEN`, `CI_SLACK_CHANNEL_DUMMY_TESTS`, `CI_SLACK_CHANNEL_ID`, `CI_SLACK_CHANNEL_ID_DAILY`, `GITHUB_TOKEN`, `HF_HUB_READ_TOKEN`, `SLACK_CIFEEDBACK_BOT_TOKEN`, `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN`

Permissions declared across the chain: `contents: read`

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |

#### Steps

1. **Setup**
   - Env:
     - `PREV_WORKFLOW_RUN_ID`: `${{ inputs.prev_workflow_run_id || env.prev_workflow_run_id }}`
     - `OTHER_WORKFLOW_RUN_ID`: `${{ inputs.other_workflow_run_id || env.other_workflow_run_id }}`

2. **Upload artifacts**
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `setup_values`
     - `path`: `setup_values`

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

# Self-hosted runner (scheduled-intel-gaudi)

| Property | Value |
|----------|-------|
| File | `self-scheduled-intel-gaudi.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
self-scheduled-intel-gaudi.yml
+-- self-scheduled-intel-gaudi3-caller.yml (job: model-ci)  <- entry point
+-- self-scheduled-intel-gaudi3-caller.yml (job: pipeline-ci)  <- entry point
+-- self-scheduled-intel-gaudi3-caller.yml (job: example-ci)  <- entry point
+-- self-scheduled-intel-gaudi3-caller.yml (job: deepspeed-ci)  <- entry point
+-- self-scheduled-intel-gaudi3-caller.yml (job: trainer-fsdp-ci)  <- entry point
```

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

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Set up Python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
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

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs-1) |
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
| Runs on | `group: ${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}` |
| Matrix | `machine_type`: 1gaudi, 2gaudi |
| Condition | `${{ inputs.job == 'run_pipelines_torch_gpu' }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`

### Examples directory (`run_examples_gpu`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}` |
| Matrix | `machine_type`: 1gaudi |
| Condition | `${{ inputs.job == 'run_examples_gpu' }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_examples_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_examples_gpu_test_reports`

### Intel Gaudi deepspeed tests (`run_torch_cuda_extensions_gpu`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ inputs.runner_scale_set }}-${{ matrix.machine_type }}` |
| Matrix | `machine_type`: 1gaudi, 2gaudi |
| Condition | `${{ inputs.job == 'run_torch_cuda_extensions_gpu' }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`
     - `path`: `reports/${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`

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

# Self-hosted runner (Intel Gaudi3 scheduled CI caller)

| Property | Value |
|----------|-------|
| File | `self-scheduled-intel-gaudi3-caller.yml` |
| Triggers | `repository_dispatch`, `workflow_dispatch`, `schedule` |

## Schedule

- `17 2 * * *`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
self-scheduled-intel-gaudi3-caller.yml [repository_dispatch, workflow_dispatch, schedule]
+-- model-ci (uses self-scheduled-intel-gaudi.yml)
|   +-- run_models_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- send_results (uses slack-report.yml)
+-- pipeline-ci (uses self-scheduled-intel-gaudi.yml)
|   +-- run_models_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- send_results (uses slack-report.yml)
+-- example-ci (uses self-scheduled-intel-gaudi.yml)
|   +-- run_models_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- send_results (uses slack-report.yml)
+-- deepspeed-ci (uses self-scheduled-intel-gaudi.yml)
|   +-- run_models_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- run_trainer_and_fsdp_gpu (uses model_jobs_intel_gaudi.yml)
|   +-- send_results (uses slack-report.yml)
+-- trainer-fsdp-ci (uses self-scheduled-intel-gaudi.yml)
    +-- run_models_gpu (uses model_jobs_intel_gaudi.yml)
    +-- run_trainer_and_fsdp_gpu (uses model_jobs_intel_gaudi.yml)
    +-- send_results (uses slack-report.yml)
```

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

# Nvidia CI (job definitions)

| Property | Value |
|----------|-------|
| File | `self-scheduled.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
self-scheduled.yml
+-- push-important-models.yml (job: model-ci)  <- entry point
+-- self-comment-ci.yml (job: model-ci)  <- entry point
+-- self-comment-ci.yml (job: quantization-ci)  <- entry point
+-- self-nightly-caller.yml (job: model-ci)  <- entry point
+-- self-past-caller.yml (job: model-ci)
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
+-- self-past-caller.yml (job: deepspeed-ci)
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
|   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
+-- self-scheduled-caller.yml (job: model-ci)  <- entry point
+-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
+-- self-scheduled-caller.yml (job: example-ci)  <- entry point
+-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
+-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
+-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
+-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
+-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `HF_HUB_READ_TOKEN` | workflow env `HF_TOKEN` |
| `GITHUB_TOKEN` | job `run_extract_warnings` step `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` with `github-token` |
| `ACCESS_REPO_INFO_TOKEN` | job `run_extract_warnings` step `Extract warnings in CI artifacts` env `access_token` |

## Jobs

### Setup (`setup`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.machine_type }}` |
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `contains(fromJSON('["run_models_gpu", "run_trainer_and_fsdp_gpu", "run_quantization_torch_gpu"]'), inputs.job)` |

#### Steps

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

### `run_models_gpu`

| Property | Value |
|----------|-------|
| Uses workflow | [model jobs](#model-jobs) |
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
| Runs on | `group: ${{ matrix.machine_type }}` |
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `${{ inputs.job == 'run_pipelines_torch_gpu' }}` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_pipelines_torch_gpu_test_reports`

### Examples directory (`run_examples_gpu`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.machine_type }}` |
| Matrix | `machine_type`: aws-g5-4xlarge-cache |
| Condition | `${{ inputs.job == 'run_examples_gpu' }}` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_examples_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_examples_gpu_test_reports`

### Torch CUDA extension tests (`run_torch_cuda_extensions_gpu`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.machine_type }}` |
| Matrix | `machine_type`: aws-g5-4xlarge-cache, aws-g5-12xlarge-cache |
| Condition | `${{ inputs.job == 'run_torch_cuda_extensions_gpu' }}` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`
     - `path`: `${{ inputs.working-directory-prefix }}/transformers/reports/${{ env.machine_type }}_run_torch_cuda_extensions_gpu_test_reports`

### `run_quantization_torch_gpu`

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.machine_type }}` |
| Depends on | `setup` |
| Condition | `${{ inputs.job == 'run_quantization_torch_gpu' }}` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_quantization_torch_gpu_${{ env.matrix_folders }}_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_quantization_torch_gpu_${{ env.matrix_folders }}_test_reports`

### Kernel tests (`run_kernels_gpu`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ matrix.machine_type }}` |
| Matrix | `machine_type`: aws-g5-4xlarge-cache |
| Condition | `${{ inputs.job == 'run_kernels_gpu' }}` |

#### Steps

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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `${{ env.machine_type }}_run_kernels_gpu_test_reports`
     - `path`: `/transformers/reports/${{ env.machine_type }}_run_kernels_gpu_test_reports`

### Extract warnings in CI artifacts (`run_extract_warnings`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-22.04` |
| Depends on | `setup`, `run_models_gpu` |
| Condition | `${{ always() && inputs.job == 'run_models_gpu' }}` |

#### Steps

1. **Checkout transformers**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **Install transformers**

3. **Show installed libraries and their versions**

4. **Create output directory**

5. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - Condition: `${{ always() }}`
   - With:
     - `name`: `warnings_in_ci`
     - `path`: `warnings_in_ci/selected_warnings.json`

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

# CI slack report

| Property | Value |
|----------|-------|
| File | `slack-report.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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

```
slack-report.yml
+-- self-scheduled-intel-gaudi.yml (job: send_results)
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: model-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: pipeline-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: example-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: deepspeed-ci)  <- entry point
|   +-- self-scheduled-intel-gaudi3-caller.yml (job: trainer-fsdp-ci)  <- entry point
+-- self-scheduled.yml (job: send_results)
    +-- push-important-models.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: model-ci)  <- entry point
    +-- self-comment-ci.yml (job: quantization-ci)  <- entry point
    +-- self-nightly-caller.yml (job: model-ci)  <- entry point
    +-- self-past-caller.yml (job: model-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-past-caller.yml (job: deepspeed-ci)
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-11)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-10)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-9)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-8)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-7)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-6)  <- entry point
    |   +-- self-nightly-past-ci-caller.yml (job: run_past_ci_tensorflow_2-5)  <- entry point
    +-- self-scheduled-caller.yml (job: model-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: torch-pipeline)  <- entry point
    +-- self-scheduled-caller.yml (job: example-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: trainer-fsdp-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: deepspeed-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: quantization-ci)  <- entry point
    +-- self-scheduled-caller.yml (job: kernels-ci)  <- entry point
    +-- self-scheduled-flash-attn-caller.yml (job: model-ci)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` | workflow env `TRANSFORMERS_CI_RESULTS_UPLOAD_TOKEN` |
| `GITHUB_TOKEN` | job `send_results` step `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` with `github-token` |
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

#### Steps

1. **Preliminary job status**
   - Env:
     - `setup_status`: `${{ inputs.setup_status }}`

2. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `fetch-depth`: `2`
     - `ref`: `${{ (github.event_name == 'issue_comment' || github.event_name == 'pull_request_target') && 'main' || (inputs.commit_sha || github.sha) }}`
     - `persist-credentials`: `false`

3. **actions/download-artifact@v8.0.1**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
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
   - Uses: `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
   - With:
     - `name`: `ci_results_${{ inputs.job }}`
     - `path`: `ci_results_${{ inputs.job }}`

# SSH into our runners

| Property | Value |
|----------|-------|
| File | `ssh-runner.yml` |
| Triggers | `workflow_dispatch` |

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

#### Steps

1. **Get runner to use**
   - Env:
     - `NUM_GPUS`: `${{ github.event.inputs.num_gpus }}`
     - `RUNNER_TYPE`: `${{ github.event.inputs.runner_type }}`

2. **Set runner to use**
   - ID: `set_runner`

### SSH (`ssh_runner`)

| Property | Value |
|----------|-------|
| Runs on | `group: ${{ needs.get_runner.outputs.RUNNER }}` |
| Depends on | `get_runner` |

#### Steps

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

# Stale Bot

| Property | Value |
|----------|-------|
| File | `stale.yml` |
| Triggers | `schedule` |

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

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **Setup Python**
   - Uses: `actions/setup-python@a26af69be951a213d495a4c3e4e4022e16d87065` (v5.6.0)
   - With:
     - `python-version`: `3.8`

3. **Install requirements**

4. **Close stale issues**

# TRL CI bot

This workflow allows trusted contributors to trigger TRL CI runs against specific Transformers commits by commenting `/trl-ci` on a PR in the TRL repo. It is meant to be used during the ongoing Trainer refactor/unbloat in Transformers, to help evaluate the downstream impact on TRL.

| Property | Value |
|----------|-------|
| File | `trl-ci-bot.yml` |
| Triggers | `issue_comment` |

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

#### Steps

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

# Secret Leaks

| Property | Value |
|----------|-------|
| File | `trufflehog.yml` |
| Triggers | `push` |

## Permissions

- `contents`: `read`

## Jobs

### `trufflehog`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `fetch-depth`: `0`
     - `persist-credentials`: `false`

2. **Secret Scanning**
   - Uses: `trufflesecurity/trufflehog@6bd2d14f7a4bc1e569fa3550efa7ec632a4fa67b`
   - With:
     - `extra_args`: `--results=verified,unknown`

# Update Transformers metadata

| Property | Value |
|----------|-------|
| File | `update_metdata.yml` |
| Triggers | `push` |

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

#### Steps

1. **actions/checkout@v4.3.1**
   - Uses: `actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5` (v4.3.1)
   - With:
     - `persist-credentials`: `false`

2. **Setup environment**

3. **Update metadata**
   - Env:
     - `HF_TOKEN`: `${{ secrets.LYSANDRE_HF_TOKEN }}`

# Upload PR Documentation

| Property | Value |
|----------|-------|
| File | `upload_pr_documentation.yml` |
| Triggers | `workflow_run` |

## Event filters

- **workflow_run**
  - workflows: `Build PR Documentation`
  - types: `completed`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
upload_pr_documentation.yml [workflow_run]
+-- build (uses huggingface/doc-builder/.github/workflows/upload_pr_documentation.yml@9ad2de8582b56c017cb530c1165116d40433f1c6)
```

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

