# Step Rendering

**Triggers:** `push`, `workflow_dispatch`

Exercises step rendering, matrix job names, runs-on normalization, and secret aggregation.

| Property | Value |
|----------|-------|
| File | `steps.yml` |

**Jobs:** [Java (java)](#java-java-build), [Deploy (env)](#deploy-env-deploy), [Verify (case)](#verify-case-verify)

## Event filters

- **push**
  - branches: `main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CHECKOUT_TOKEN` | job `build` step `Checkout` with `token` |
| `REGISTRY_PASSWORD` | job `deploy` step `Push image` (run) |
| `IMAGE_SIGNING_KEY` | job `deploy` step `Push image` env `IMAGE_SIGNING_KEY` |
| `SLACK_WEBHOOK` | job `deploy` step `Notify` (if); job `deploy` step `Notify` (run) |

**Variables:**

| Name | Used by |
|------|---------|
| `UPLOAD_BUCKET` | job `build` step `upload` (run) |
| `REGISTRY_USER` | job `deploy` step `Push image` (run) |
| `NOTIFY_CHANNEL` | job `deploy` step `Notify` (if) |

## Jobs

### Java (java) (`build`)

| Property | Value |
|----------|-------|
| Runs on | `self-hosted, linux, x64` |
| Matrix | `java`: 17, 21, 24 |

<details>
<summary>Steps (5)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v4.1.1`
   - With:
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.CHECKOUT_TOKEN }}`

2. **actions/setup-java@v4.2.0**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `${{ matrix.java }}`

3. **actions/cache@v4**
   - With:
     - `path`: `~/.m2`
     - `key`: `maven-${{ hashFiles('**/pom.xml') }}`

4. **./mvnw -B clean verify** `[continue-on-error]`

5. **upload**
   - ID: `upload`

</details>

### Deploy (env) (`deploy`)

| Property | Value |
|----------|-------|
| Runs on | `group: deploy-runners, labels: linux, arm64` |
| Matrix | `target.env`: staging, production; `target.url`: https://staging.example.com, https://example.com |
| Depends on | `build` |
| Condition | `github.event_name == 'push' &&<br>startsWith(github.ref, 'refs/heads/main')` |

<details>
<summary>Steps (4)</summary>

1. **Push image**
   - Env:
     - `DOCKER_BUILDKIT`: `1`
     - `IMAGE_SIGNING_KEY`: `${{ secrets.IMAGE_SIGNING_KEY }}`

2. **Notify** `[continue-on-error]`
   - Condition: `${{ vars.NOTIFY_CHANNEL && secrets.SLACK_WEBHOOK }}`

3. **Report deploy status: ${{ matrix.target.env }}**

4. **Comment on PR**
   - Uses: `actions/github-script@v7`
   - With:
     - `script`: `` const env = `${{ matrix.target.env }}`; github.rest.issues.createComment({ body: `Deployed to ${env}` }); ``
     - `result-encoding`: -

</details>

### Verify (case) (`verify`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Matrix | `case`: a, b, c (combinations adjusted by include/exclude) |

<details>
<summary>Steps (1)</summary>

1. **Run checks**

</details>

