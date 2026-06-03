# Step Rendering

Exercises step rendering, matrix job names, runs-on normalization, and secret aggregation.

| Property | Value |
|----------|-------|
| File | `steps.yml` |
| Triggers | `push`, `workflow_dispatch` |

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

### Java ${{ matrix.java }} (`build`)

| Property | Value |
|----------|-------|
| Runs on | `self-hosted, linux, x64` |
| Matrix | `java`: 17, 21, 24 |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@8f4b7f84864484a7bf31766abe9204da3cbe65b3` (v4.1.1)
   - With:
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.CHECKOUT_TOKEN }}`

2. **actions/setup-java@v4.2.0**
   - Uses: `actions/setup-java@5896cecc08fd8a1fbdfaf517e29b571164b031f7` (v4.2.0)
   - With:
     - `distribution`: `temurin`
     - `java-version`: `${{ matrix.java }}`

3. **actions/cache@v4**
   - Uses: `actions/cache@v4`
   - With:
     - `path`: `~/.m2`
     - `key`: `maven-${{ hashFiles('**/pom.xml') }}`

4. **./mvnw -B clean verify** `[continue-on-error]`

5. **upload**
   - ID: `upload`

### Deploy ${{ matrix.target.env }} (`deploy`)

| Property | Value |
|----------|-------|
| Runs on | `group: deploy-runners, labels: linux, arm64` |
| Matrix | `target.env`: staging, production; `target.url`: https://staging.example.com, https://example.com |
| Depends on | `build` |
| Condition | `github.event_name == 'push' &&<br>startsWith(github.ref, 'refs/heads/main')` |

#### Steps

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

### Verify ${{ matrix.case }} (`verify`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Run checks**

