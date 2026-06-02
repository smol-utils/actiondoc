# Contents

- [Build and Deploy Snapshot](#build-and-deploy-snapshot)
- [Build Pull Request](#build-pull-request)
- [CI](#ci)
- [Distribute](#distribute)
- [Release Milestone](#release-milestone)
- [Release](#release)
- [Run CodeQL Analysis](#run-codeql-analysis)
- [Run System Tests](#run-system-tests)
- [Trigger Docs Build](#trigger-docs-build)
- [Verify](#verify)
- [Await HTTP Resource](#await-http-resource)
- [Build](#build)
- [Create GitHub Release](#create-github-release)
- [Prepare Gradle Build](#prepare-gradle-build)
- [Print JVM thread dumps](#print-jvm-thread-dumps)
- [Publish Gradle Plugin](#publish-gradle-plugin)
- [Publish to SDKMAN!](#publish-to-sdkman)
- [Send Notification](#send-notification)
- [Sync to Maven Central](#sync-to-maven-central)
- [Update Homebrew Tap](#update-homebrew-tap)

# Build and Deploy Snapshot

| Property | Value |
|----------|-------|
| File | `build-and-deploy-snapshot.yml` |
| Triggers | `workflow_dispatch`, `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`

## Call graph (rooted at this workflow)

```
build-and-deploy-snapshot.yml [workflow_dispatch, push]
+-- build-and-deploy-snapshot / Build and Publish (uses ./.github/actions/build)
+-- build-and-deploy-snapshot / Send Notification (uses ./.github/actions/send-notification)
+-- verify (uses verify.yml)
    +-- verify / Send Notification (uses ./send-notification/.github/actions/send-notification (unresolved))
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `commercial-repository-password`, `commercial-repository-username`, `google-chat-webhook-url`, `opensource-repository-password`, `opensource-repository-username`, `token`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_ARTIFACTORY_PASSWORD` | job `build-and-deploy-snapshot` step `Build and Publish` with `commercial-repository-password`; job `build-and-deploy-snapshot` step `Deploy` with `password`; job `verify` secrets `commercial-repository-password` |
| `COMMERCIAL_ARTIFACTORY_USERNAME` | job `build-and-deploy-snapshot` step `Build and Publish` with `commercial-repository-username`; job `build-and-deploy-snapshot` step `Deploy` with `username`; job `verify` secrets `commercial-repository-username` |
| `DEVELOCITY_ACCESS_KEY` | job `build-and-deploy-snapshot` step `Build and Publish` with `develocity-access-key` |
| `ARTIFACTORY_PASSWORD` | job `build-and-deploy-snapshot` step `Deploy` with `password`; job `verify` secrets `opensource-repository-password` |
| `GPG_PRIVATE_KEY` | job `build-and-deploy-snapshot` step `Deploy` with `signing-key` |
| `GPG_PASSPHRASE` | job `build-and-deploy-snapshot` step `Deploy` with `signing-passphrase` |
| `ARTIFACTORY_USERNAME` | job `build-and-deploy-snapshot` step `Deploy` with `username`; job `verify` secrets `opensource-repository-username` |
| `GOOGLE_CHAT_WEBHOOK_URL` | job `build-and-deploy-snapshot` step `Send Notification` with `webhook-url`; job `verify` secrets `google-chat-webhook-url` |
| `GITHUB_TOKEN` | job `trigger-docs-build` step `Run Deploy Docs Workflow` env `GH_TOKEN` |
| `GH_ACTIONS_REPO_TOKEN` | job `verify` secrets `token` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_RELEASE_REPO_URL` | job `build-and-deploy-snapshot` step `Build and Publish` with `commercial-release-repository-url` |
| `COMMERCIAL_SNAPSHOT_REPO_URL` | job `build-and-deploy-snapshot` step `Build and Publish` with `commercial-snapshot-repository-url` |
| `COMMERCIAL` | job `build-and-deploy-snapshot` step `Deploy` with `build-name`; job `build-and-deploy-snapshot` step `Deploy` with `password`; job `build-and-deploy-snapshot` step `Deploy` with `project`; job `build-and-deploy-snapshot` step `Deploy` with `repository`; job `build-and-deploy-snapshot` step `Deploy` with `username` |
| `COMMERCIAL_DEPLOY_REPO_URL` | job `build-and-deploy-snapshot` step `Deploy` with `uri` |

## Jobs

### Build and Deploy Snapshot (`build-and-deploy-snapshot`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_MEDIUM \|\| 'ubuntu-latest' }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' \|\| github.repository == 'spring-projects/spring-boot-commercial' }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Build and Publish**
   - ID: `build-and-publish`
   - Uses: `./.github/actions/build`
   - With:
     - `commercial-release-repository-url`: `${{ vars.COMMERCIAL_RELEASE_REPO_URL }}` - URL of the release repository
     - `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_PASSWORD }}` - Password for authentication with the commercial repository
     - `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_USERNAME }}` - Username for authentication with the commercial repository
     - `commercial-snapshot-repository-url`: `${{ vars.COMMERCIAL_SNAPSHOT_REPO_URL }}` - URL of the snapshot repository
     - `develocity-access-key`: `${{ secrets.DEVELOCITY_ACCESS_KEY }}` - Access key for authentication with ge.spring.io
     - `gradle-cache-read-only`: `false` - Whether Gradle's cache should be read only
     - `publish`: `true` - Whether to publish artifacts ready for deployment to Artifactory

3. **Deploy**
   - Uses: `spring-io/artifactory-deploy-action@926d7f7cc810569395346bf3a4d91b380b3e355b` (v0.0.4)
   - With:
     - `build-name`: `${{ vars.COMMERCIAL && format('spring-boot-commercial-{0}', '4.1.x') || format('spring-boot-{0}', '4.1.x') }}`
     - `folder`: `deployment-repository`
     - `password`: `${{ vars.COMMERCIAL && secrets.COMMERCIAL_ARTIFACTORY_PASSWORD || secrets.ARTIFACTORY_PASSWORD }}`
     - `project`: `${{ vars.COMMERCIAL && 'spring' }}`
     - `repository`: `${{ vars.COMMERCIAL && 'spring-enterprise-maven-dev-local' || 'libs-snapshot-local' }}`
     - `signing-key`: `${{ secrets.GPG_PRIVATE_KEY }}`
     - `signing-passphrase`: `${{ secrets.GPG_PASSPHRASE }}`
     - `threads`: `8`
     - `uri`: `${{ vars.COMMERCIAL_DEPLOY_REPO_URL || 'https://repo.spring.io' }}`
     - `username`: `${{ vars.COMMERCIAL && secrets.COMMERCIAL_ARTIFACTORY_USERNAME || secrets.ARTIFACTORY_USERNAME }}`

4. **Send Notification**
   - Uses: `./.github/actions/send-notification`
   - Condition: `always()`
   - With:
     - `build-scan-url`: `${{ steps.build-and-publish.outputs.build-scan-url }}` - URL of the build scan to include in the notification
     - `run-name`: `${{ format('{0} | Linux | Java 25', github.ref_name) }}` - Name of the run to include in the notification
     - `status`: `${{ job.status }}` - Status of the job (required)
     - `webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}` - Google Chat Webhook URL (required)

### Trigger Docs Build (`trigger-docs-build`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-deploy-snapshot` |

**Permissions:**

- `actions`: `write`

#### Steps

1. **Run Deploy Docs Workflow**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

### Verify (`verify`)

| Property | Value |
|----------|-------|
| Uses workflow | [Verify](#verify) |
| Depends on | `build-and-deploy-snapshot` |

#### Inputs forwarded

- `version`: `${{ needs.build-and-deploy-snapshot.outputs.version }}`

#### Secrets forwarded

- `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_PASSWORD }}`
- `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_USERNAME }}`
- `google-chat-webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}`
- `opensource-repository-password`: `${{ secrets.ARTIFACTORY_PASSWORD }}`
- `opensource-repository-username`: `${{ secrets.ARTIFACTORY_USERNAME }}`
- `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}`

# Build Pull Request

| Property | Value |
|----------|-------|
| File | `build-pull-request.yml` |
| Triggers | `pull_request` |

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
build-pull-request.yml [pull_request]
+-- build / Build (uses ./.github/actions/build)
+-- build / Print JVM Thread Dumps When Cancelled (uses ./.github/actions/print-jvm-thread-dumps)
```

## Jobs

### Build Pull Request (`build`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_MEDIUM \|\| 'ubuntu-latest' }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Build**
   - ID: `build`
   - Uses: `./.github/actions/build`

3. **Print JVM Thread Dumps When Cancelled**
   - Uses: `./.github/actions/print-jvm-thread-dumps`
   - Condition: `cancelled()`

4. **Upload Build Reports**
   - Uses: `actions/upload-artifact@v7`
   - Condition: `failure()`
   - With:
     - `name`: `build-reports`
     - `path`: `**/build/reports/`

# CI

| Property | Value |
|----------|-------|
| File | `ci.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
ci.yml [push]
+-- ci / Build (uses ./.github/actions/build)
+-- ci / Send Notification (uses ./.github/actions/send-notification)
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_ARTIFACTORY_PASSWORD` | job `ci` step `Build` with `commercial-repository-password` |
| `COMMERCIAL_ARTIFACTORY_USERNAME` | job `ci` step `Build` with `commercial-repository-username` |
| `DEVELOCITY_ACCESS_KEY` | job `ci` step `Build` with `develocity-access-key` |
| `GOOGLE_CHAT_WEBHOOK_URL` | job `ci` step `Send Notification` with `webhook-url` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_RELEASE_REPO_URL` | job `ci` step `Build` with `commercial-release-repository-url` |
| `COMMERCIAL_SNAPSHOT_REPO_URL` | job `ci` step `Build` with `commercial-snapshot-repository-url` |

## Jobs

### ${{ matrix.os.name}} | Java ${{ matrix.java.version}} (`ci`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os.id }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' \|\| github.repository == 'spring-projects/spring-boot-commercial' }}` |

#### Steps

1. **Prepare Windows runner**
   - Condition: `${{ runner.os == 'Windows' }}`

2. **Check Out Code**
   - Uses: `actions/checkout@v6`

3. **Build**
   - ID: `build`
   - Uses: `./.github/actions/build`
   - With:
     - `commercial-release-repository-url`: `${{ vars.COMMERCIAL_RELEASE_REPO_URL }}` - URL of the release repository
     - `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_PASSWORD }}` - Password for authentication with the commercial repository
     - `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_USERNAME }}` - Username for authentication with the commercial repository
     - `commercial-snapshot-repository-url`: `${{ vars.COMMERCIAL_SNAPSHOT_REPO_URL }}` - URL of the snapshot repository
     - `develocity-access-key`: `${{ secrets.DEVELOCITY_ACCESS_KEY }}` - Access key for authentication with ge.spring.io
     - `gradle-cache-read-only`: `false` - Whether Gradle's cache should be read only
     - `java-early-access`: `${{ matrix.java.early-access || 'false' }}` - Whether the Java version is in early access
     - `java-distribution`: `${{ matrix.java.distribution }}` - Java distribution to use
     - `java-toolchain`: `${{ matrix.java.toolchain }}` - Whether a Java toolchain should be used
     - `java-version`: `${{ matrix.java.version }}` - Java version to compile and test with

4. **Send Notification**
   - Uses: `./.github/actions/send-notification`
   - Condition: `always()`
   - With:
     - `build-scan-url`: `${{ steps.build.outputs.build-scan-url }}` - URL of the build scan to include in the notification
     - `run-name`: `${{ format('{0} | {1} | Java {2}', github.ref_name, matrix.os.name, matrix.java.version) }}` - Name of the run to include in the notification
     - `status`: `${{ job.status }}` - Status of the job (required)
     - `webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}` - Google Chat Webhook URL (required)

# Distribute

| Property | Value |
|----------|-------|
| File | `distribute.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `build-number` | string | Yes | - | Number of the build to use to create the bundle |
| `create-bundle` | boolean | Yes | `true` | Whether to create the bundle. If unchecked, only the bundle distribution is executed |
| `version` | string | Yes | - | Version to bundle and distribute |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_ARTIFACTORY_USERNAME` | job `distribute-spring-enterprise-release-bundle` step `Create Bundle` (run); job `distribute-spring-enterprise-release-bundle` step `Distribute Bundle` (run) |
| `COMMERCIAL_ARTIFACTORY_PASSWORD` | job `distribute-spring-enterprise-release-bundle` step `Create Bundle` (run); job `distribute-spring-enterprise-release-bundle` step `Distribute Bundle` (run) |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL` | job `distribute-spring-enterprise-release-bundle` step `Create Bundle` (if); job `distribute-spring-enterprise-release-bundle` step `Sleep` (if); job `distribute-spring-enterprise-release-bundle` step `Distribute Bundle` (if) |

## Jobs

### `distribute-spring-enterprise-release-bundle`

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |

#### Steps

1. **Create Bundle**
   - Condition: `${{ vars.COMMERCIAL && inputs.create-bundle }}`

2. **Sleep**
   - Condition: `${{ vars.COMMERCIAL && inputs.create-bundle }}`

3. **Distribute Bundle**
   - Condition: `${{ vars.COMMERCIAL }}`

# Release Milestone

| Property | Value |
|----------|-------|
| File | `release-milestone.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - tags: `v4.1.0-M[0-9]`, `v4.1.0-RC[0-9]`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`

## Call graph (rooted at this workflow)

```
release-milestone.yml [push]
+-- build-and-stage-release / Build and Publish (uses ./.github/actions/build)
+-- verify (uses verify.yml)
|   +-- verify / Send Notification (uses ./send-notification/.github/actions/send-notification (unresolved))
+-- sync-to-maven-central / Sync to Maven Central (uses ./.github/actions/sync-to-maven-central)
+-- publish-gradle-plugin / Publish (uses ./.github/actions/publish-gradle-plugin)
+-- create-github-release / Create GitHub Release (uses ./.github/actions/create-github-release)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `commercial-repository-password`, `commercial-repository-username`, `google-chat-webhook-url`, `opensource-repository-password`, `opensource-repository-username`, `token`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DEVELOCITY_ACCESS_KEY` | job `build-and-stage-release` step `Build and Publish` with `develocity-access-key` |
| `ARTIFACTORY_PASSWORD` | job `build-and-stage-release` step `Stage Release` with `password`; job `verify` secrets `opensource-repository-password` |
| `GPG_PRIVATE_KEY` | job `build-and-stage-release` step `Stage Release` with `signing-key` |
| `GPG_PASSPHRASE` | job `build-and-stage-release` step `Stage Release` with `signing-passphrase` |
| `ARTIFACTORY_USERNAME` | job `build-and-stage-release` step `Stage Release` with `username`; job `verify` secrets `opensource-repository-username` |
| `COMMERCIAL_ARTIFACTORY_RO_PASSWORD` | job `verify` secrets `commercial-repository-password` |
| `COMMERCIAL_ARTIFACTORY_RO_USERNAME` | job `verify` secrets `commercial-repository-username` |
| `GOOGLE_CHAT_WEBHOOK_URL` | job `verify` secrets `google-chat-webhook-url` |
| `GH_ACTIONS_REPO_TOKEN` | job `verify` secrets `token`; job `create-github-release` step `Create GitHub Release` with `token` |
| `CENTRAL_TOKEN_PASSWORD` | job `sync-to-maven-central` step `Sync to Maven Central` with `central-token-password` |
| `CENTRAL_TOKEN_USERNAME` | job `sync-to-maven-central` step `Sync to Maven Central` with `central-token-username` |
| `JF_ARTIFACTORY_SPRING` | job `sync-to-maven-central` step `Sync to Maven Central` with `jfrog-cli-config-token`; job `promote-release` step `Set up JFrog CLI` env `JF_ENV_SPRING`; job `publish-gradle-plugin` step `Publish` with `jfrog-cli-config-token` |
| `GRADLE_PLUGIN_PUBLISH_KEY` | job `publish-gradle-plugin` step `Publish` with `gradle-plugin-publish-key` |
| `GRADLE_PLUGIN_PUBLISH_SECRET` | job `publish-gradle-plugin` step `Publish` with `gradle-plugin-publish-secret` |
| `GITHUB_TOKEN` | job `trigger-docs-build` step `Run Deploy Docs Workflow` env `GH_TOKEN` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL` | job `sync-to-maven-central` (if); job `publish-gradle-plugin` (if); job `create-github-release` step `Create GitHub Release` with `commercial` |

## Jobs

### Build and Stage Release (`build-and-stage-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_MEDIUM \|\| 'ubuntu-latest' }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Build and Publish**
   - ID: `build-and-publish`
   - Uses: `./.github/actions/build`
   - With:
     - `develocity-access-key`: `${{ secrets.DEVELOCITY_ACCESS_KEY }}` - Access key for authentication with ge.spring.io
     - `gradle-cache-read-only`: `false` - Whether Gradle's cache should be read only
     - `publish`: `true` - Whether to publish artifacts ready for deployment to Artifactory

3. **Stage Release**
   - Uses: `spring-io/artifactory-deploy-action@926d7f7cc810569395346bf3a4d91b380b3e355b` (v0.0.4)
   - With:
     - `build-name`: `${{ format('spring-boot-{0}', steps.build-and-publish.outputs.version)}}`
     - `folder`: `deployment-repository`
     - `password`: `${{ secrets.ARTIFACTORY_PASSWORD }}`
     - `repository`: `libs-staging-local`
     - `signing-key`: `${{ secrets.GPG_PRIVATE_KEY }}`
     - `signing-passphrase`: `${{ secrets.GPG_PASSPHRASE }}`
     - `threads`: `8`
     - `uri`: `https://repo.spring.io`
     - `username`: `${{ secrets.ARTIFACTORY_USERNAME }}`

### Verify (`verify`)

| Property | Value |
|----------|-------|
| Uses workflow | [Verify](#verify) |
| Depends on | `build-and-stage-release` |

#### Inputs forwarded

- `staging`: `true`
- `version`: `${{ needs.build-and-stage-release.outputs.version }}`

#### Secrets forwarded

- `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_RO_PASSWORD }}`
- `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_RO_USERNAME }}`
- `google-chat-webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}`
- `opensource-repository-password`: `${{ secrets.ARTIFACTORY_PASSWORD }}`
- `opensource-repository-username`: `${{ secrets.ARTIFACTORY_USERNAME }}`
- `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}`

### Sync to Maven Central (`sync-to-maven-central`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `verify` |
| Condition | `${{ !vars.COMMERCIAL }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Sync to Maven Central**
   - Uses: `./.github/actions/sync-to-maven-central`
   - With:
     - `central-token-password`: `${{ secrets.CENTRAL_TOKEN_PASSWORD }}` - Password for authentication with central.sonatype.com (required)
     - `central-token-username`: `${{ secrets.CENTRAL_TOKEN_USERNAME }}` - Username for authentication with central.sonatype.com (required)
     - `jfrog-cli-config-token`: `${{ secrets.JF_ARTIFACTORY_SPRING }}` - Config token for the JFrog CLI (required)
     - `spring-boot-version`: `${{ needs.build-and-stage-release.outputs.version }}` - Version of Spring Boot that is being synced to Central (required)

### Promote Release (`promote-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |

#### Steps

1. **Set up JFrog CLI**
   - Uses: `jfrog/setup-jfrog-cli@1641575d87647fb969c0545f0b6a76873e328b7c` (v5.0.0)
   - Env:
     - `JF_ENV_SPRING`: `${{ secrets.JF_ARTIFACTORY_SPRING }}`

2. **Promote build**

### Publish Gradle Plugin (`publish-gradle-plugin`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |
| Condition | `${{ !vars.COMMERCIAL }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Publish**
   - Uses: `./.github/actions/publish-gradle-plugin`
   - With:
     - `gradle-plugin-publish-key`: `${{ secrets.GRADLE_PLUGIN_PUBLISH_KEY }}` - Gradle publishing key (required)
     - `gradle-plugin-publish-secret`: `${{ secrets.GRADLE_PLUGIN_PUBLISH_SECRET }}` - Gradle publishing secret (required)
     - `jfrog-cli-config-token`: `${{ secrets.JF_ARTIFACTORY_SPRING }}` - Config token for the JFrog CLI (required)
     - `plugin-version`: `${{ needs.build-and-stage-release.outputs.version }}` - Version of the plugin (required)

### Trigger Docs Build (`trigger-docs-build`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-and-stage-release`, `promote-release` |

**Permissions:**

- `actions`: `write`

#### Steps

1. **Run Deploy Docs Workflow**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

### Create GitHub Release (`create-github-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `promote-release`, `publish-gradle-plugin`, `trigger-docs-build` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Create GitHub Release**
   - Uses: `./.github/actions/create-github-release`
   - With:
     - `commercial`: `${{ vars.COMMERCIAL }}` - Whether to generate the changelog for the commercial release (required)
     - `milestone`: `${{ needs.build-and-stage-release.outputs.version }}` - Name of the GitHub milestone for which a release will be created (required)
     - `pre-release`: `true` - Whether the release is a pre-release (a milestone or release candidate)
     - `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}` - Token to use for authentication with GitHub (required)

# Release

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - tags: `v4.1.[0-9]+`

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`

## Call graph (rooted at this workflow)

```
release.yml [push]
+-- build-and-stage-release / Build and Publish (uses ./.github/actions/build)
+-- build-and-stage-release / Send Notification (uses ./.github/actions/send-notification)
+-- verify (uses verify.yml)
|   +-- verify / Send Notification (uses ./send-notification/.github/actions/send-notification (unresolved))
+-- sync-to-maven-central / Sync to Maven Central (uses ./.github/actions/sync-to-maven-central)
+-- publish-gradle-plugin / Publish (uses ./.github/actions/publish-gradle-plugin)
+-- publish-to-sdkman / Publish to SDKMAN! (uses ./.github/actions/publish-to-sdkman)
+-- update-homebrew-tap / Update Homebrew Tap (uses ./.github/actions/update-homebrew-tap)
+-- create-github-release / Create GitHub Release (uses ./.github/actions/create-github-release)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `commercial-repository-password`, `commercial-repository-username`, `google-chat-webhook-url`, `opensource-repository-password`, `opensource-repository-username`, `token`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_ARTIFACTORY_PASSWORD` | job `build-and-stage-release` step `Build and Publish` with `commercial-repository-password`; job `build-and-stage-release` step `Stage Release` with `password`; job `verify` secrets `commercial-repository-password` |
| `COMMERCIAL_ARTIFACTORY_USERNAME` | job `build-and-stage-release` step `Build and Publish` with `commercial-repository-username`; job `build-and-stage-release` step `Stage Release` with `username`; job `verify` secrets `commercial-repository-username` |
| `DEVELOCITY_ACCESS_KEY` | job `build-and-stage-release` step `Build and Publish` with `develocity-access-key` |
| `ARTIFACTORY_PASSWORD` | job `build-and-stage-release` step `Stage Release` with `password`; job `verify` secrets `opensource-repository-password` |
| `GPG_PRIVATE_KEY` | job `build-and-stage-release` step `Stage Release` with `signing-key` |
| `GPG_PASSPHRASE` | job `build-and-stage-release` step `Stage Release` with `signing-passphrase` |
| `ARTIFACTORY_USERNAME` | job `build-and-stage-release` step `Stage Release` with `username`; job `verify` secrets `opensource-repository-username` |
| `GOOGLE_CHAT_WEBHOOK_URL` | job `build-and-stage-release` step `Send Notification` with `webhook-url`; job `verify` secrets `google-chat-webhook-url` |
| `GH_ACTIONS_REPO_TOKEN` | job `verify` secrets `token`; job `update-homebrew-tap` step `Update Homebrew Tap` with `token`; job `create-github-release` step `Create GitHub Release` with `token` |
| `CENTRAL_TOKEN_PASSWORD` | job `sync-to-maven-central` step `Sync to Maven Central` with `central-token-password` |
| `CENTRAL_TOKEN_USERNAME` | job `sync-to-maven-central` step `Sync to Maven Central` with `central-token-username` |
| `JF_ARTIFACTORY_SPRING` | job `sync-to-maven-central` step `Sync to Maven Central` with `jfrog-cli-config-token`; job `promote-release` step `Set up JFrog CLI` env `JF_ENV_SPRING`; job `publish-gradle-plugin` step `Publish` with `jfrog-cli-config-token` |
| `COMMERCIAL_JF_ARTIFACTORY_SPRING` | job `promote-release` step `Set up JFrog CLI` env `JF_ENV_SPRING` |
| `GRADLE_PLUGIN_PUBLISH_KEY` | job `publish-gradle-plugin` step `Publish` with `gradle-plugin-publish-key` |
| `GRADLE_PLUGIN_PUBLISH_SECRET` | job `publish-gradle-plugin` step `Publish` with `gradle-plugin-publish-secret` |
| `SDKMAN_CONSUMER_KEY` | job `publish-to-sdkman` step `Publish to SDKMAN!` with `sdkman-consumer-key` |
| `SDKMAN_CONSUMER_TOKEN` | job `publish-to-sdkman` step `Publish to SDKMAN!` with `sdkman-consumer-token` |
| `GITHUB_TOKEN` | job `trigger-docs-build` step `Run Deploy Docs Workflow` env `GH_TOKEN` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL_RELEASE_REPO_URL` | job `build-and-stage-release` step `Build and Publish` with `commercial-release-repository-url` |
| `COMMERCIAL_SNAPSHOT_REPO_URL` | job `build-and-stage-release` step `Build and Publish` with `commercial-snapshot-repository-url` |
| `COMMERCIAL` | job `build-and-stage-release` step `Stage Release` with `build-name`; job `build-and-stage-release` step `Stage Release` with `password`; job `build-and-stage-release` step `Stage Release` with `project`; job `build-and-stage-release` step `Stage Release` with `repository`; job `build-and-stage-release` step `Stage Release` with `username`; job `sync-to-maven-central` (if); job `promote-release` step `Set up JFrog CLI` env `JF_ENV_SPRING`; job `promote-release` step `Promote open source build` (if); job `promote-release` step `Promote commercial build` (if); job `publish-gradle-plugin` (if); job `publish-to-sdkman` (if); job `create-github-release` step `Create GitHub Release` with `commercial` |
| `COMMERCIAL_DEPLOY_REPO_URL` | job `build-and-stage-release` step `Stage Release` with `uri` |

## Jobs

### Build and Stage Release (`build-and-stage-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_MEDIUM \|\| 'ubuntu-latest' }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' \|\| github.repository == 'spring-projects/spring-boot-commercial' }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Build and Publish**
   - ID: `build-and-publish`
   - Uses: `./.github/actions/build`
   - With:
     - `commercial-release-repository-url`: `${{ vars.COMMERCIAL_RELEASE_REPO_URL }}` - URL of the release repository
     - `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_PASSWORD }}` - Password for authentication with the commercial repository
     - `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_USERNAME }}` - Username for authentication with the commercial repository
     - `commercial-snapshot-repository-url`: `${{ vars.COMMERCIAL_SNAPSHOT_REPO_URL }}` - URL of the snapshot repository
     - `develocity-access-key`: `${{ secrets.DEVELOCITY_ACCESS_KEY }}` - Access key for authentication with ge.spring.io
     - `gradle-cache-read-only`: `false` - Whether Gradle's cache should be read only
     - `publish`: `true` - Whether to publish artifacts ready for deployment to Artifactory

3. **Stage Release**
   - Uses: `spring-io/artifactory-deploy-action@926d7f7cc810569395346bf3a4d91b380b3e355b` (v0.0.4)
   - With:
     - `build-name`: `${{ vars.COMMERCIAL && format('spring-boot-commercial-{0}', steps.build-and-publish.outputs.version) || format('spring-boot-{0}', steps.build-and-publish.outputs.version) }}`
     - `folder`: `deployment-repository`
     - `password`: `${{ vars.COMMERCIAL && secrets.COMMERCIAL_ARTIFACTORY_PASSWORD || secrets.ARTIFACTORY_PASSWORD }}`
     - `project`: `${{ vars.COMMERCIAL && 'spring' }}`
     - `repository`: `${{ vars.COMMERCIAL && 'spring-enterprise-maven-stage-local' || 'libs-staging-local' }}`
     - `signing-key`: `${{ secrets.GPG_PRIVATE_KEY }}`
     - `signing-passphrase`: `${{ secrets.GPG_PASSPHRASE }}`
     - `threads`: `8`
     - `uri`: `${{ vars.COMMERCIAL_DEPLOY_REPO_URL || 'https://repo.spring.io' }}`
     - `username`: `${{ vars.COMMERCIAL && secrets.COMMERCIAL_ARTIFACTORY_USERNAME || secrets.ARTIFACTORY_USERNAME }}`

4. **Send Notification**
   - Uses: `./.github/actions/send-notification`
   - Condition: `failure()`
   - With:
     - `run-name`: `${{ format('{0} | Release Staging | {1}', github.ref_name, inputs.version) }}` - Name of the run to include in the notification
     - `status`: `${{ job.status }}` - Status of the job (required)
     - `webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}` - Google Chat Webhook URL (required)

### Verify (`verify`)

| Property | Value |
|----------|-------|
| Uses workflow | [Verify](#verify) |
| Depends on | `build-and-stage-release` |

#### Inputs forwarded

- `staging`: `true`
- `version`: `${{ needs.build-and-stage-release.outputs.version }}`

#### Secrets forwarded

- `commercial-repository-password`: `${{ secrets.COMMERCIAL_ARTIFACTORY_PASSWORD }}`
- `commercial-repository-username`: `${{ secrets.COMMERCIAL_ARTIFACTORY_USERNAME }}`
- `google-chat-webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}`
- `opensource-repository-password`: `${{ secrets.ARTIFACTORY_PASSWORD }}`
- `opensource-repository-username`: `${{ secrets.ARTIFACTORY_USERNAME }}`
- `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}`

### Sync to Maven Central (`sync-to-maven-central`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `verify` |
| Condition | `${{ !vars.COMMERCIAL }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Sync to Maven Central**
   - Uses: `./.github/actions/sync-to-maven-central`
   - With:
     - `central-token-password`: `${{ secrets.CENTRAL_TOKEN_PASSWORD }}` - Password for authentication with central.sonatype.com (required)
     - `central-token-username`: `${{ secrets.CENTRAL_TOKEN_USERNAME }}` - Username for authentication with central.sonatype.com (required)
     - `jfrog-cli-config-token`: `${{ secrets.JF_ARTIFACTORY_SPRING }}` - Config token for the JFrog CLI (required)
     - `spring-boot-version`: `${{ needs.build-and-stage-release.outputs.version }}` - Version of Spring Boot that is being synced to Central (required)

### Promote Release (`promote-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |

#### Steps

1. **Set up JFrog CLI**
   - Uses: `jfrog/setup-jfrog-cli@1641575d87647fb969c0545f0b6a76873e328b7c` (v5.0.0)
   - Env:
     - `JF_ENV_SPRING`: `${{ vars.COMMERCIAL && secrets.COMMERCIAL_JF_ARTIFACTORY_SPRING || secrets.JF_ARTIFACTORY_SPRING }}`

2. **Promote open source build**
   - Condition: `${{ !vars.COMMERCIAL }}`

3. **Promote commercial build**
   - Condition: `${{ vars.COMMERCIAL }}`

### Publish Gradle Plugin (`publish-gradle-plugin`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |
| Condition | `${{ !vars.COMMERCIAL }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Publish**
   - Uses: `./.github/actions/publish-gradle-plugin`
   - With:
     - `gradle-plugin-publish-key`: `${{ secrets.GRADLE_PLUGIN_PUBLISH_KEY }}` - Gradle publishing key (required)
     - `gradle-plugin-publish-secret`: `${{ secrets.GRADLE_PLUGIN_PUBLISH_SECRET }}` - Gradle publishing secret (required)
     - `jfrog-cli-config-token`: `${{ secrets.JF_ARTIFACTORY_SPRING }}` - Config token for the JFrog CLI (required)
     - `plugin-version`: `${{ needs.build-and-stage-release.outputs.version }}` - Version of the plugin (required)

### Publish to SDKMAN! (`publish-to-sdkman`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |
| Condition | `${{ !vars.COMMERCIAL }}` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Publish to SDKMAN!**
   - Uses: `./.github/actions/publish-to-sdkman`
   - With:
     - `make-default`: `true` - Whether the release should be made the default version
     - `sdkman-consumer-key`: `${{ secrets.SDKMAN_CONSUMER_KEY }}` - Key for publishing to SDKMAN! (required)
     - `sdkman-consumer-token`: `${{ secrets.SDKMAN_CONSUMER_TOKEN }}` - Token for publishing to SDKMAN! (required)
     - `spring-boot-version`: `${{ needs.build-and-stage-release.outputs.version }}` - Version to publish (required)

### Update Homebrew Tap (`update-homebrew-tap`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `sync-to-maven-central` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Update Homebrew Tap**
   - Uses: `./.github/actions/update-homebrew-tap`
   - With:
     - `spring-boot-version`: `${{ needs.build-and-stage-release.outputs.version }}` - The version to publish (required)
     - `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}` - Token to use for GitHub authentication (required)

### Trigger Docs Build (`trigger-docs-build`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-and-stage-release`, `promote-release` |

**Permissions:**

- `actions`: `write`

#### Steps

1. **Run Deploy Docs Workflow**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

### Create GitHub Release (`create-github-release`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Depends on | `build-and-stage-release`, `promote-release`, `publish-gradle-plugin`, `publish-to-sdkman`, `trigger-docs-build`, `update-homebrew-tap` |

#### Steps

1. **Check Out Code**
   - Uses: `actions/checkout@v6`

2. **Create GitHub Release**
   - Uses: `./.github/actions/create-github-release`
   - With:
     - `commercial`: `${{ vars.COMMERCIAL }}` - Whether to generate the changelog for the commercial release (required)
     - `milestone`: `${{ needs.build-and-stage-release.outputs.version }}` - Name of the GitHub milestone for which a release will be created (required)
     - `token`: `${{ secrets.GH_ACTIONS_REPO_TOKEN }}` - Token to use for authentication with GitHub (required)

# Run CodeQL Analysis

| Property | Value |
|----------|-------|
| File | `run-codeql-analysis.yml` |
| Triggers | `push`, `pull_request`, `workflow_dispatch` |

## Permissions

All scopes: `read-all`.

## Call graph (rooted at this workflow)

```
run-codeql-analysis.yml [push, pull_request, workflow_dispatch]
+-- run-analysis (uses spring-io/github-actions/.github/workflows/codeql-analysis.yml@7dc305df87410aa851b873d2f1fd33ccbb7d0aa8)
```

## Transitive requirements (from full call graph)

External workflows referenced: `spring-io/github-actions/.github/workflows/codeql-analysis.yml@7dc305df87410aa851b873d2f1fd33ccbb7d0aa8`

## Jobs

### `run-analysis`

| Property | Value |
|----------|-------|
| Uses workflow | `spring-io/github-actions/.github/workflows/codeql-analysis.yml@7dc305df87410aa851b873d2f1fd33ccbb7d0aa8` (external) |

**Permissions:**

- `actions`: `read`
- `contents`: `read`
- `security-events`: `write`

# Run System Tests

| Property | Value |
|----------|-------|
| File | `run-system-tests.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Call graph (rooted at this workflow)

```
run-system-tests.yml [push]
+-- run-system-tests / Prepare Gradle Build (uses ./.github/actions/prepare-gradle-build)
+-- run-system-tests / Send Notification (uses ./.github/actions/send-notification)
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DEVELOCITY_ACCESS_KEY` | job `run-system-tests` step `Prepare Gradle Build` with `develocity-access-key` |
| `GOOGLE_CHAT_WEBHOOK_URL` | job `run-system-tests` step `Send Notification` with `webhook-url` |

## Jobs

### Java 17, 21 (`run-system-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_MEDIUM \|\| 'ubuntu-latest' }}` |
| Condition | `${{ github.repository == 'spring-projects/spring-boot' }}` |

#### Steps

1. **Switch Docker to Overlay2**

2. **Check Out Code**
   - Uses: `actions/checkout@v6`

3. **Prepare Gradle Build**
   - Uses: `./.github/actions/prepare-gradle-build`
   - With:
     - `develocity-access-key`: `${{ secrets.DEVELOCITY_ACCESS_KEY }}` - Access key for authentication with ge.spring.io
     - `java-toolchain`: `${{ matrix.java.toolchain }}` - Whether a Java toolchain should be used
     - `java-version`: `${{ matrix.java.version }}` - Java version to use for the build

4. **Run System Tests**
   - ID: `run-system-tests`

5. **Show docker info**
   - Condition: `always()`

6. **List docker images**
   - Condition: `always()`

7. **Send Notification**
   - Uses: `./.github/actions/send-notification`
   - Condition: `always()`
   - With:
     - `build-scan-url`: `${{ steps.run-system-tests.outputs.build-scan-url }}` - URL of the build scan to include in the notification
     - `run-name`: `${{ format('{0} | System Tests | Java {1}', github.ref_name, matrix.java.version) }}` - Name of the run to include in the notification
     - `status`: `${{ job.status }}` - Status of the job (required)
     - `webhook-url`: `${{ secrets.GOOGLE_CHAT_WEBHOOK_URL }}` - Google Chat Webhook URL (required)

# Trigger Docs Build

| Property | Value |
|----------|-------|
| File | `trigger-docs-build.yml` |
| Triggers | `push`, `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `build-version` | - | No | - | Version being build (e.g. 1.0.3-SNAPSHOT) |
| `build-sha` | - | No | - | Enter the SHA to build (e.g. 82c97891569821a7f91a77ca074232e0b54ca7a5) |
| `build-refname` | - | No | - | Git refname to build (e.g., 1.0.x) |

## Event filters

- **push**
  - branches: `main`
  - paths: `antora/*`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `trigger-docs-build` step `Trigger Workflow` env `GH_TOKEN` |

## Jobs

### Trigger Docs Build (`trigger-docs-build`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |
| Condition | `github.repository_owner == 'spring-projects'` |

**Permissions:**

- `actions`: `write`

#### Steps

1. **Check Out**
   - Uses: `actions/checkout@v6`
   - With:
     - `ref`: `docs-build`

2. **Trigger Workflow**
   - Env:
     - `GH_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

# Verify

| Property | Value |
|----------|-------|
| File | `verify.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `staging` | boolean | No | `false` | Whether the release to verify is in the staging repository |
| `version` | string | Yes | - | Version to verify |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `commercial-repository-password` | No | Password for authentication with the commercial repository |
| `commercial-repository-username` | No | Username for authentication with the commercial repository |
| `google-chat-webhook-url` | Yes | Google Chat Webhook URL |
| `opensource-repository-password` | No | Password for authentication with the open-source repository |
| `opensource-repository-username` | No | Username for authentication with the open-source repository |
| `token` | Yes | Token to use for authentication with GitHub |

## Permissions

- `contents`: `read`

## Called by

```
verify.yml
+-- build-and-deploy-snapshot.yml (job: verify)  <- entry point
+-- release-milestone.yml (job: verify)  <- entry point
+-- release.yml (job: verify)  <- entry point
```

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `token` | job `verify` step `Check Out Release Verification Tests` with `token` |
| `commercial-repository-password` | job `verify` step `Run Release Verification Tests` env `RVT_COMMERCIAL_REPOSITORY_PASSWORD` |
| `commercial-repository-username` | job `verify` step `Run Release Verification Tests` env `RVT_COMMERCIAL_REPOSITORY_USERNAME` |
| `opensource-repository-password` | job `verify` step `Run Release Verification Tests` env `RVT_OSS_REPOSITORY_PASSWORD` |
| `opensource-repository-username` | job `verify` step `Run Release Verification Tests` env `RVT_OSS_REPOSITORY_USERNAME` |
| `google-chat-webhook-url` | job `verify` step `Send Notification` with `webhook-url` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMERCIAL` | job `verify` step `Set Up Homebrew` (if); job `verify` step `Run Release Verification Tests` env `RVT_RELEASE_TYPE` |

## Jobs

### Verify (`verify`)

| Property | Value |
|----------|-------|
| Runs on | `${{ vars.UBUNTU_SMALL \|\| 'ubuntu-latest' }}` |

#### Steps

1. **Check Out Release Verification Tests**
   - Uses: `actions/checkout@v6`
   - With:
     - `ref`: `v0.0.15`
     - `repository`: `spring-projects/spring-boot-release-verification`
     - `token`: `${{ secrets.token }}`

2. **Check Out Send Notification Action**
   - Uses: `actions/checkout@v6`
   - With:
     - `path`: `send-notification`
     - `sparse-checkout`: `.github/actions/send-notification`

3. **Set Up Java**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `liberica`
     - `java-version`: `17`

4. **Set Up Homebrew**
   - Uses: `Homebrew/actions/setup-homebrew@7657c9512f50e1c35b640971116425935bab3eea`
   - Condition: `${{ !vars.COMMERCIAL }}`
   - With:
     - `stable`: `true`

5. **Set Up Gradle**
   - Uses: `gradle/actions/setup-gradle@0723195856401067f7a2779048b490ace7a47d7c` (v5.0.2)
   - With:
     - `cache-read-only`: `false`

6. **Configure Gradle Properties**

7. **Run Release Verification Tests**
   - Env:
     - `RVT_COMMERCIAL_REPOSITORY_PASSWORD`: `${{ secrets.commercial-repository-password }}`
     - `RVT_COMMERCIAL_REPOSITORY_USERNAME`: `${{ secrets.commercial-repository-username }}`
     - `RVT_OSS_REPOSITORY_PASSWORD`: `${{ secrets.opensource-repository-password }}`
     - `RVT_OSS_REPOSITORY_USERNAME`: `${{ secrets.opensource-repository-username }}`
     - `RVT_RELEASE_TYPE`: `${{ vars.COMMERCIAL && 'commercial' || 'oss' }}`
     - `RVT_STAGING`: `${{ inputs.staging }}`
     - `RVT_VERSION`: `${{ inputs.version }}`

8. **Upload Build Reports on Failure**
   - Uses: `actions/upload-artifact@v7`
   - Condition: `failure()`
   - With:
     - `name`: `build-reports`
     - `path`: `**/build/reports/`

9. **Send Notification**
   - Uses: `./send-notification/.github/actions/send-notification`
   - Condition: `always()`
   - With:
     - `run-name`: `${{ format('{0} | Verification | {1}', github.ref_name, inputs.version) }}`
     - `status`: `${{ job.status }}`
     - `webhook-url`: `${{ secrets.google-chat-webhook-url }}`

# Await HTTP Resource

Waits for an HTTP resource to be available (a HEAD request succeeds)

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `url` | URL of the resource to await | Yes | - |

# Build

Builds the project, optionally publishing it to a local deployment repository

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `commercial-release-repository-url` | URL of the release repository | No | - |
| `commercial-repository-password` | Password for authentication with the commercial repository | No | - |
| `commercial-repository-username` | Username for authentication with the commercial repository | No | - |
| `commercial-snapshot-repository-url` | URL of the snapshot repository | No | - |
| `develocity-access-key` | Access key for authentication with ge.spring.io | No | - |
| `gradle-cache-read-only` | Whether Gradle's cache should be read only | No | `true` |
| `java-distribution` | Java distribution to use | No | `liberica` |
| `java-early-access` | Whether the Java version is in early access | No | `false` |
| `java-toolchain` | Whether a Java toolchain should be used | No | `false` |
| `java-version` | Java version to compile and test with | No | `25` |
| `publish` | Whether to publish artifacts ready for deployment to Artifactory | No | `false` |

## Outputs

| Name | Description |
|------|-------------|
| `build-scan-url` | URL, if any, of the build scan produced by the build |
| `version` | Version that was built |

# Create GitHub Release

Create the release on GitHub with a changelog

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `commercial` | Whether to generate the changelog for the commercial release | Yes | - |
| `milestone` | Name of the GitHub milestone for which a release will be created | Yes | - |
| `pre-release` | Whether the release is a pre-release (a milestone or release candidate) | No | `false` |
| `token` | Token to use for authentication with GitHub | Yes | - |

# Prepare Gradle Build

Prepares a Gradle build. Sets up Java and Gradle and configures Gradle properties

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `cache-read-only` | Whether Gradle's cache should be read only | No | `true` |
| `develocity-access-key` | Access key for authentication with ge.spring.io | No | - |
| `java-distribution` | Java distribution to use | No | `liberica` |
| `java-early-access` | Whether the Java version is in early access. When true, forces java-distribution to temurin | No | `false` |
| `java-toolchain` | Whether a Java toolchain should be used | No | `false` |
| `java-version` | Java version to use for the build | No | `25` |

# Print JVM thread dumps

Prints a thread dump for all running JVMs

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

# Publish Gradle Plugin

Publishes Spring Boot's Gradle plugin to the Plugin Portal

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `build-number` | Build number to use when downloading plugin artifacts | No | `${{ github.run_number }}` |
| `gradle-plugin-publish-key` | Gradle publishing key | Yes | - |
| `gradle-plugin-publish-secret` | Gradle publishing secret | Yes | - |
| `jfrog-cli-config-token` | Config token for the JFrog CLI | Yes | - |
| `plugin-version` | Version of the plugin | Yes | - |

# Publish to SDKMAN!

Publishes the release as a new candidate version on SDKMAN!

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `make-default` | Whether the release should be made the default version | No | `false` |
| `sdkman-consumer-key` | Key for publishing to SDKMAN! | Yes | - |
| `sdkman-consumer-token` | Token for publishing to SDKMAN! | Yes | - |
| `spring-boot-version` | Version to publish | Yes | - |

# Send Notification

Sends a Google Chat message as a notification of the job's outcome

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `build-scan-url` | URL of the build scan to include in the notification | No | - |
| `run-name` | Name of the run to include in the notification | No | `${{ format('{0} {1}', github.ref_name, github.job) }}` |
| `status` | Status of the job | Yes | - |
| `webhook-url` | Google Chat Webhook URL | Yes | - |

# Sync to Maven Central

Syncs a release to Maven Central and waits for it to be available for use

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `central-token-password` | Password for authentication with central.sonatype.com | Yes | - |
| `central-token-username` | Username for authentication with central.sonatype.com | Yes | - |
| `jfrog-cli-config-token` | Config token for the JFrog CLI | Yes | - |
| `spring-boot-version` | Version of Spring Boot that is being synced to Central | Yes | - |

# Update Homebrew Tap

Updates the Homebrew Tap for the Spring Boot CLI

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `spring-boot-version` | The version to publish | Yes | - |
| `token` | Token to use for GitHub authentication | Yes | - |

