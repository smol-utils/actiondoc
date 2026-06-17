# jreleaser

9 workflows, 6 reusable workflows

## Contents

**Workflows**

- [Build - pull_request](#build)
- [Clear cache - schedule, workflow_dispatch](#clear-cache)
- [CodeQL - workflow_dispatch, push, pull_request](#codeql)
- [EarlyAccess - push](#earlyaccess)
- [Lint - push](#lint)
- [OpenSSF Scorecard - branch_protection_rule, schedule, push, workflow_dispatch](#openssf-scorecard)
- [Release - workflow_dispatch](#release)
- [SmokeTests - push](#smoketests)
- [Trigger Early Access - workflow_dispatch](#trigger-early-access)

**Reusable workflows**

- [X-BachInfo](#x-bachinfo)
- [X-Jlink](#x-jlink)
- [X-JPackage](#x-jpackage)
- [X-NativeImage](#x-nativeimage)
- [X-Precheck](#x-precheck)
- [X-UpdateWiki](#x-updatewiki)

## Workflows by trigger

- **push**: [CodeQL](#codeql), [EarlyAccess](#earlyaccess), [Lint](#lint), [OpenSSF Scorecard](#openssf-scorecard), [SmokeTests](#smoketests)
- **workflow_dispatch**: [Clear cache](#clear-cache), [CodeQL](#codeql), [OpenSSF Scorecard](#openssf-scorecard), [Release](#release), [Trigger Early Access](#trigger-early-access)
- **pull_request**: [Build](#build), [CodeQL](#codeql)
- **schedule**: [Clear cache](#clear-cache), [OpenSSF Scorecard](#openssf-scorecard)
- **branch_protection_rule**: [OpenSSF Scorecard](#openssf-scorecard)

## Secrets and variables used across this repository

**Secrets:**

| Name | Used by |
|------|---------|
| `BLUESKY_HANDLE` | [Release](#release) |
| `BLUESKY_HOST` | [Release](#release) |
| `BLUESKY_PASSWORD` | [Release](#release) |
| `CODECOV_TOKEN` | [SmokeTests](#smoketests) |
| `COMMIT_EMAIL` | [Release](#release) |
| `COVERALLS_TOKEN` | [SmokeTests](#smoketests) |
| `GIT_ACCESS_TOKEN` | [EarlyAccess](#earlyaccess), [Release](#release), [SmokeTests](#smoketests), [Trigger Early Access](#trigger-early-access) |
| `GIT_PAT_TOKEN` | [SmokeTests](#smoketests) |
| `GPG_PASSPHRASE` | [EarlyAccess](#earlyaccess), [Release](#release), [SmokeTests](#smoketests) |
| `GPG_PUBLIC_KEY` | [EarlyAccess](#earlyaccess), [Release](#release), [SmokeTests](#smoketests) |
| `GPG_SECRET_KEY` | [EarlyAccess](#earlyaccess), [Release](#release), [SmokeTests](#smoketests) |
| `GRADLE_PUBLISH_KEY` | [Release](#release) |
| `GRADLE_PUBLISH_SECRET` | [Release](#release) |
| `JRELEASER_DOCKER_PASSWORD` | [EarlyAccess](#earlyaccess), [Release](#release) |
| `JRELEASER_OCI_COMPARTMENTID` | [EarlyAccess](#earlyaccess), [Release](#release), [SmokeTests](#smoketests) |
| `MASTODON_ACCESS_TOKEN` | [Release](#release) |
| `NOTICEABLE_APIKEY` | [Release](#release) |
| `OPENCOLLECTIVE_TOKEN` | [Release](#release) |
| `SDKMAN_CONSUMER_KEY` | [Release](#release) |
| `SDKMAN_CONSUMER_TOKEN` | [Release](#release) |
| `SONARCLOUD_TOKEN` | [SmokeTests](#smoketests) |
| `SONATYPE_PASSWORD` | [Release](#release) |
| `SONATYPE_USERNAME` | [Release](#release) |
| `gh-access-token` | [X-BachInfo](#x-bachinfo), [X-UpdateWiki](#x-updatewiki) |
| `github-token` | [X-Precheck](#x-precheck) |
| `gpg-passphrase` | [X-Jlink](#x-jlink) |
| `oci-compartment-id` | [X-Jlink](#x-jlink) |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMIT_EMAIL` | [Release](#release), [X-BachInfo](#x-bachinfo) |
| `GH_BOT_EMAIL` | [EarlyAccess](#earlyaccess), [Release](#release) |
| `GRAAL_JAVA_VERSION` | [SmokeTests](#smoketests), [X-NativeImage](#x-nativeimage) |
| `JAVA_DISTRO` | [Release](#release), [SmokeTests](#smoketests), [Trigger Early Access](#trigger-early-access), [X-Jlink](#x-jlink), [X-JPackage](#x-jpackage), [X-NativeImage](#x-nativeimage) |
| `JAVA_VERSION` | [Release](#release), [SmokeTests](#smoketests), [Trigger Early Access](#trigger-early-access), [X-Jlink](#x-jlink), [X-JPackage](#x-jpackage), [X-NativeImage](#x-nativeimage) |

## Permissions across this repository

| Scope | Level |
|-------|-------|
| `(all scopes)` | `read-all` |
| `actions` | `write` (also granted as `read` elsewhere) |
| `attestations` | `write` |
| `contents` | `write` (also granted as `read` elsewhere) |
| `id-token` | `write` (OIDC) |
| `security-events` | `write` |

# Build

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `build.yml` |

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`, cancel-in-progress: `true`

## Jobs

### Build (`build`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Matrix | `os`: ubuntu-latest, macos-15-intel, windows-latest |

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `21`
     - `distribution`: `zulu`
     - `cache`: `gradle`

3. **Build**

</details>

[Back to contents](#contents)

# Clear cache

**Triggers:** `schedule`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `cache.yml` |

## Schedule

- `0 3 * * *`

## Jobs

### Delete all caches (`clear`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (1)</summary>

1. **Clear caches**
   - Uses: `easimon/wipe-cache@e7ab82e64c328fd39c2e96933d426cd72ac2beba`

</details>

[Back to contents](#contents)

# CodeQL

**Triggers:** `workflow_dispatch`, `push`, `pull_request`

| Property | Value |
|----------|-------|
| File | `codeql.yml` |

**Jobs:** [Precheck](#precheck-precheck), [CodeQL](#codeql-codeql)

## Event filters

- **push**
  - branches: `main`
- **pull_request**
  - branches: `main`

## Permissions

- `security-events`: `write`
- `actions`: `read`
- `contents`: `read`

## Call graph (rooted at this workflow)

- `precheck` uses `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | `precheck`: (`github-token`)<br>`codeql`: Cancel previous run (`access_token`) |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main` (external) |

#### Secrets forwarded

- `github-token`: `${{ secrets.GITHUB_TOKEN }}`

### CodeQL (`codeql`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (6)</summary>

1. **Checkout repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Cancel previous run**
   - Uses: `styfle/cancel-workflow-action@v0.13.1`
   - With:
     - `access_token`: `${{ secrets.GITHUB_TOKEN }}`

3. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@v4.35.1`
   - With:
     - `languages`: `java`
     - `build-mode`: `manual`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `21`
     - `distribution`: `zulu`
     - `cache`: `gradle`

5. **Build**

6. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@v4.35.1`
   - With:
     - `category`: `/language:java`

</details>

[Back to contents](#contents)

# EarlyAccess

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `early-access.yml` |

**Jobs:** [Precheck](#precheck-precheck-1), [Jlink](#jlink-jlink), [JPackage](#jpackage-jpackage), [Native Image](#native-image-native-image), [Release](#release-release), [Provenance](#provenance-provenance), [Update Wiki](#update-wiki-update-wiki)

## Event filters

- **push**
  - branches: `main`

## Permissions

- `actions`: `write`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Call graph (rooted at this workflow)

- `precheck` uses `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`
- `jlink` uses `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`
- `jpackage` uses `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`
- `native-image` uses `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`
- `provenance` uses `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`
- `update-wiki` uses `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `gh-access-token`, `github-token`, `gpg-passphrase`, `oci-compartment-id`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`, `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`, `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`, `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`, `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`, `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | `precheck`: (`github-token`) |
| `GPG_PASSPHRASE` | `jlink`: (`gpg-passphrase`)<br>`release`: Release (`JRELEASER_GPG_PASSPHRASE`) |
| `JRELEASER_OCI_COMPARTMENTID` | `jlink`: (`oci-compartment-id`) |
| `GIT_ACCESS_TOKEN` | `native-image`: (`gh-access-token`)<br>`release`: Release (`JRELEASER_GITHUB_TOKEN`)<br>`update-wiki`: (`gh-access-token`) |
| `GPG_PUBLIC_KEY` | `release`: Release (`JRELEASER_GPG_PUBLIC_KEY`) |
| `GPG_SECRET_KEY` | `release`: Release (`JRELEASER_GPG_SECRET_KEY`) |
| `JRELEASER_DOCKER_PASSWORD` | `release`: Release (`JRELEASER_DOCKER_DEFAULT_PASSWORD`) |

**Variables:**

| Name | Used by |
|------|---------|
| `GH_BOT_EMAIL` | `update-wiki`: (`commit-email`) |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main` (external) |

#### Secrets forwarded

- `github-token`: `${{ secrets.GITHUB_TOKEN }}`

### Jlink (`jlink`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main` (external) |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.version }}`

#### Secrets forwarded

- `gpg-passphrase`: `${{ secrets.GPG_PASSPHRASE }}`
- `oci-compartment-id`: `${{ secrets.JRELEASER_OCI_COMPARTMENTID }}`

### JPackage (`jpackage`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main` (external) |
| Depends on | `precheck`, `jlink` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.version }}`
- `project-effective-version`: `early-access`

### Native Image (`native-image`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main` (external) |
| Depends on | `precheck`, `jlink` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.version }}`

#### Secrets forwarded

- `gh-access-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

### Release (`release`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `jlink`, `jpackage`, `native-image` |

**Permissions:**

- `id-token`: `write` (OIDC)
- `contents`: `read`
- `attestations`: `write`

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Download artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download java-archive**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive`

4. **Download jlink**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

5. **Download jpackage**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `jpackage-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage`

6. **Download native-image**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `native-image-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image`

7. **Release**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `full-release`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ needs.precheck.outputs.version }}`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `JRELEASER_DOCKER_DEFAULT_PASSWORD`: `${{ secrets.JRELEASER_DOCKER_PASSWORD }}`

8. **JReleaser release output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-release`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

9. **SLSA**
   - ID: `slsa`

</details>

### Provenance (`provenance`)

| Property | Value |
|----------|-------|
| Uses workflow | `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0` (external) |
| Depends on | `precheck`, `release` |

#### Inputs forwarded

- `base64-subjects`: `${{ needs.release.outputs.hashes }}`
- `upload-assets`: `true`
- `upload-tag-name`: `${{ needs.release.outputs.tagname }}`
- `provenance-name`: `jreleaser-all-${{ needs.release.outputs.tagname }}.intoto.jsonl`

### Update Wiki (`update-wiki`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main` (external) |
| Depends on | `precheck`, `release` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.version }}`
- `project-tag`: `${{ needs.release.outputs.tagname }}`
- `commit-email`: `${{ vars.GH_BOT_EMAIL }}`

#### Secrets forwarded

- `gh-access-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

[Back to contents](#contents)

# Lint

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `lint.yml` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Jobs

### Lint (`lint`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (2)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **actionlint**
   - ID: `actionlint`
   - Uses: `raven-actions/actionlint@v2.1.2`

</details>

[Back to contents](#contents)

# OpenSSF Scorecard

**Triggers:** `branch_protection_rule`, `schedule`, `push`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `openssf-scorecard.yml` |

**Jobs:** [Precheck](#precheck-precheck-2), [Scorecards analysis](#scorecards-analysis-analysis)

## Schedule

- `30 1 * * 6`

## Event filters

- **push**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Call graph (rooted at this workflow)

- `precheck` uses `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | `precheck`: (`github-token`) |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main` (external) |

#### Secrets forwarded

- `github-token`: `${{ secrets.GITHUB_TOKEN }}`

### Scorecards analysis (`analysis`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

**Permissions:**

- `security-events`: `write`
- `id-token`: `write` (OIDC)
- `actions`: `read`
- `contents`: `read`

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
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@v4.35.1`
   - With:
     - `sarif_file`: `results.sarif`

</details>

[Back to contents](#contents)

# Release

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [Precheck](#precheck-precheck-3), [Jlink](#jlink-jlink-1), [JPackage](#jpackage-jpackage-1), [Native Image](#native-image-native-image-1), [Release](#release-release-1), [Provenance](#provenance-provenance-1), [Update Wiki](#update-wiki-update-wiki-1), [Update Website](#update-website-update-website)

## Permissions

- `actions`: `write`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Call graph (rooted at this workflow)

- `jlink` uses `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`
- `jpackage` uses `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`
- `native-image` uses `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`
- `provenance` uses `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`
- `update-wiki` uses `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `gh-access-token`, `gpg-passphrase`, `oci-compartment-id`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`, `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`, `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`, `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`, `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMIT_EMAIL` | `precheck`: Commit version (`run`)<br>`release`: Bump version (`run`) |
| `GPG_PASSPHRASE` | `jlink`: (`gpg-passphrase`)<br>`release`: Release (`JRELEASER_GPG_PASSPHRASE`) |
| `JRELEASER_OCI_COMPARTMENTID` | `jlink`: (`oci-compartment-id`) |
| `GIT_ACCESS_TOKEN` | `native-image`: (`gh-access-token`)<br>`release`: Release (`JRELEASER_GITHUB_TOKEN`)<br>`update-wiki`: (`gh-access-token`)<br>`update-website`: Checkout (`token`) |
| `GRADLE_PUBLISH_KEY` | `release`: Deploy (`GRADLE_PUBLISH_KEY`) |
| `GRADLE_PUBLISH_SECRET` | `release`: Deploy (`GRADLE_PUBLISH_SECRET`) |
| `GPG_PUBLIC_KEY` | `release`: Release (`JRELEASER_GPG_PUBLIC_KEY`) |
| `GPG_SECRET_KEY` | `release`: Release (`JRELEASER_GPG_SECRET_KEY`) |
| `JRELEASER_DOCKER_PASSWORD` | `release`: Release (`JRELEASER_DOCKER_DEFAULT_PASSWORD`) |
| `SDKMAN_CONSUMER_KEY` | `release`: Release (`JRELEASER_SDKMAN_CONSUMER_KEY`) |
| `SDKMAN_CONSUMER_TOKEN` | `release`: Release (`JRELEASER_SDKMAN_CONSUMER_TOKEN`) |
| `MASTODON_ACCESS_TOKEN` | `release`: Release (`JRELEASER_MASTODON_ACCESS_TOKEN`) |
| `SONATYPE_USERNAME` | `release`: Release (`JRELEASER_MAVENCENTRAL_USERNAME`) |
| `SONATYPE_PASSWORD` | `release`: Release (`JRELEASER_MAVENCENTRAL_PASSWORD`) |
| `NOTICEABLE_APIKEY` | `release`: Release (`JRELEASER_HTTP_NOTICEABLE_PASSWORD`) |
| `OPENCOLLECTIVE_TOKEN` | `release`: Release (`JRELEASER_OPENCOLLECTIVE_TOKEN`) |
| `BLUESKY_HOST` | `release`: Release (`JRELEASER_BLUESKY_HOST`) |
| `BLUESKY_HANDLE` | `release`: Release (`JRELEASER_BLUESKY_HANDLE`) |
| `BLUESKY_PASSWORD` | `release`: Release (`JRELEASER_BLUESKY_PASSWORD`) |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | `release`: Setup Java (`java-version`)<br>`update-website`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `release`: Setup Java (`distribution`)<br>`update-website`: Setup Java (`distribution`) |
| `GH_BOT_EMAIL` | `update-wiki`: (`commit-email`) |
| `COMMIT_EMAIL` | `update-website`: Commit (`COMMIT_EMAIL`) |

## Jobs

### Precheck (`precheck`)

<details>
<summary>Steps (3)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`

2. **Version**
   - ID: `version`

3. **Commit version**

</details>

### Jlink (`jlink`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main` (external) |
| Depends on | `precheck` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.release-version }}`

#### Secrets forwarded

- `gpg-passphrase`: `${{ secrets.GPG_PASSPHRASE }}`
- `oci-compartment-id`: `${{ secrets.JRELEASER_OCI_COMPARTMENTID }}`

### JPackage (`jpackage`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main` (external) |
| Depends on | `precheck`, `jlink` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.release-version }}`
- `project-effective-version`: `${{ needs.precheck.outputs.release-version }}`

### Native Image (`native-image`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main` (external) |
| Depends on | `precheck`, `jlink` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.release-version }}`

#### Secrets forwarded

- `gh-access-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

### Release (`release`)

| Property | Value |
|----------|-------|
| Depends on | `precheck`, `jlink`, `jpackage`, `native-image` |

<details>
<summary>Steps (13)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`
     - `ref`: `main`
     - `fetch-depth`: `0`

2. **Download artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download java-archive**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive`

4. **Download jlink**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

5. **Download jpackage**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `jpackage-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage`

6. **Download native-image**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `native-image-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image`

7. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

8. **Deploy**
   - Env:
     - `GRADLE_PUBLISH_KEY`: `${{ secrets.GRADLE_PUBLISH_KEY }}`
     - `GRADLE_PUBLISH_SECRET`: `${{ secrets.GRADLE_PUBLISH_SECRET }}`

9. **Upload deploy artifacts**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `7`
     - `name`: `deploy`
     - `path`: `build/repos/local/release/`

10. **Release**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `full-release`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ needs.precheck.outputs.release-version }}`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `JRELEASER_DOCKER_DEFAULT_PASSWORD`: `${{ secrets.JRELEASER_DOCKER_PASSWORD }}`
     - `JRELEASER_SDKMAN_CONSUMER_KEY`: `${{ secrets.SDKMAN_CONSUMER_KEY }}`
     - `JRELEASER_SDKMAN_CONSUMER_TOKEN`: `${{ secrets.SDKMAN_CONSUMER_TOKEN }}`
     - `JRELEASER_MASTODON_ACCESS_TOKEN`: `${{ secrets.MASTODON_ACCESS_TOKEN }}`
     - `JRELEASER_MAVENCENTRAL_USERNAME`: `${{ secrets.SONATYPE_USERNAME }}`
     - `JRELEASER_MAVENCENTRAL_PASSWORD`: `${{ secrets.SONATYPE_PASSWORD }}`
     - `JRELEASER_HTTP_NOTICEABLE_PASSWORD`: `${{ secrets.NOTICEABLE_APIKEY }}`
     - `JRELEASER_OPENCOLLECTIVE_TOKEN`: `${{ secrets.OPENCOLLECTIVE_TOKEN }}`
     - `JRELEASER_BLUESKY_HOST`: `${{ secrets.BLUESKY_HOST }}`
     - `JRELEASER_BLUESKY_HANDLE`: `${{ secrets.BLUESKY_HANDLE }}`
     - `JRELEASER_BLUESKY_PASSWORD`: `${{ secrets.BLUESKY_PASSWORD }}`

11. **JReleaser release output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-release`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

12. **SLSA**
   - ID: `slsa`

13. **Bump version**

</details>

### Provenance (`provenance`)

| Property | Value |
|----------|-------|
| Uses workflow | `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0` (external) |
| Depends on | `precheck`, `release` |

#### Inputs forwarded

- `base64-subjects`: `${{ needs.release.outputs.hashes }}`
- `upload-assets`: `true`
- `upload-tag-name`: `${{ needs.release.outputs.tagname }}`
- `provenance-name`: `jreleaser-all-${{ needs.precheck.outputs.release-version }}.intoto.jsonl`

### Update Wiki (`update-wiki`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main` (external) |
| Depends on | `precheck`, `release` |

#### Inputs forwarded

- `project-version`: `${{ needs.precheck.outputs.release-version }}`
- `project-tag`: `${{ needs.release.outputs.tagname }}`
- `commit-email`: `${{ vars.GH_BOT_EMAIL }}`
- `template-params`: `-PincludeSboms`

#### Secrets forwarded

- `gh-access-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

### Update Website (`update-website`)

| Property | Value |
|----------|-------|
| Depends on | `precheck`, `release` |

<details>
<summary>Steps (4)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`
     - `repository`: `jreleaser/jreleaser.github.io`
     - `ref`: `main`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

2. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`

3. **Download assets**
   - Env:
     - `TAG`: `${{ needs.release.outputs.tagname }}`
     - `RELEASE_VERSION`: `${{ needs.precheck.outputs.release-version }}`

4. **Commit**
   - Env:
     - `TAG`: `${{ needs.release.outputs.tagname }}`
     - `RELEASE_VERSION`: `${{ needs.precheck.outputs.release-version }}`
     - `NEXT_VERSION`: `${{ needs.precheck.outputs.next-version }}`
     - `COMMIT_EMAIL`: `${{ vars.COMMIT_EMAIL }}`

</details>

[Back to contents](#contents)

# SmokeTests

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `smoke-tests.yml` |
| Default runs-on | `${{ matrix.job.os }}` |

**Jobs:**

- [Precheck](#precheck-precheck-4)
- [CLI (os)](#cli-os-build-cli)
- [Tool (os)](#tool-os-build-tool)
- [Ant (os)](#ant-os-build-ant)
- [Gradle (os)](#gradle-os-build-gradle)
- [Maven (os)](#maven-os-build-maven)
- [Unit Test (os)](#unit-test-os-unit-tests)
- [Coveralls](#coveralls-coveralls)
- [Codecov](#codecov-codecov)
- [Sonar](#sonar-sonar)

## Event filters

- **push**
  - branches: `main`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `CI` | `true` |
| `GPG_PASSPHRASE` | `${{ secrets.GPG_PASSPHRASE }}` |
| `JRELEASER_OCI_COMPARTMENTID` | `${{ secrets.JRELEASER_OCI_COMPARTMENTID }}` |

## Call graph (rooted at this workflow)

- `precheck` uses `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GPG_PASSPHRASE` | workflow env: (`GPG_PASSPHRASE`)<br>`build-cli`: JReleaser (`JRELEASER_GPG_PASSPHRASE`)<br>`build-tool`: JReleaser (`JRELEASER_GPG_PASSPHRASE`)<br>`build-ant`: JReleaser (`JRELEASER_GPG_PASSPHRASE`)<br>`build-gradle`: JReleaser (`JRELEASER_GPG_PASSPHRASE`), Clean (`JRELEASER_GPG_PASSPHRASE`)<br>`build-maven`: JReleaser (`JRELEASER_GPG_PASSPHRASE`) |
| `JRELEASER_OCI_COMPARTMENTID` | workflow env: (`JRELEASER_OCI_COMPARTMENTID`) |
| `GITHUB_TOKEN` | `precheck`: (`github-token`) |
| `GIT_ACCESS_TOKEN` | `build-cli`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)<br>`build-tool`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)<br>`build-ant`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)<br>`build-gradle`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)<br>`build-maven`: Setup Graal (`github-token`), Checkout smoketests repository (`token`) |
| `GIT_PAT_TOKEN` | `build-cli`: JReleaser (`JRELEASER_GITHUB_TOKEN`)<br>`build-tool`: JReleaser (`JRELEASER_GITHUB_TOKEN`)<br>`build-ant`: JReleaser (`JRELEASER_GITHUB_TOKEN`)<br>`build-gradle`: JReleaser (`JRELEASER_GITHUB_TOKEN`), Clean (`JRELEASER_GITHUB_TOKEN`)<br>`build-maven`: JReleaser (`JRELEASER_GITHUB_TOKEN`) |
| `GPG_PUBLIC_KEY` | `build-cli`: JReleaser (`JRELEASER_GPG_PUBLIC_KEY`)<br>`build-tool`: JReleaser (`JRELEASER_GPG_PUBLIC_KEY`)<br>`build-ant`: JReleaser (`JRELEASER_GPG_PUBLIC_KEY`)<br>`build-gradle`: JReleaser (`JRELEASER_GPG_PUBLIC_KEY`), Clean (`JRELEASER_GPG_PUBLIC_KEY`)<br>`build-maven`: JReleaser (`JRELEASER_GPG_PUBLIC_KEY`) |
| `GPG_SECRET_KEY` | `build-cli`: JReleaser (`JRELEASER_GPG_SECRET_KEY`)<br>`build-tool`: JReleaser (`JRELEASER_GPG_SECRET_KEY`)<br>`build-ant`: JReleaser (`JRELEASER_GPG_SECRET_KEY`)<br>`build-gradle`: JReleaser (`JRELEASER_GPG_SECRET_KEY`), Clean (`JRELEASER_GPG_SECRET_KEY`)<br>`build-maven`: JReleaser (`JRELEASER_GPG_SECRET_KEY`) |
| `COVERALLS_TOKEN` | `coveralls`: Upload coverage to Coveralls (`COVERALLS_REPO_TOKEN`) |
| `CODECOV_TOKEN` | `codecov`: Upload coverage to Codecov (`token`) |
| `SONARCLOUD_TOKEN` | `sonar`: Sonar (`run`) |

**Variables:**

| Name | Used by |
|------|---------|
| `GRAAL_JAVA_VERSION` | `build-cli`: Setup Graal (`java-version`)<br>`build-tool`: Setup Graal (`java-version`)<br>`build-ant`: Setup Graal (`java-version`)<br>`build-gradle`: Setup Graal (`java-version`)<br>`build-maven`: Setup Graal (`java-version`) |
| `JAVA_VERSION` | `build-cli`: Setup Java (`java-version`)<br>`build-tool`: Setup Java (`java-version`)<br>`build-ant`: Setup Java (`java-version`)<br>`build-gradle`: Setup Java (`java-version`)<br>`build-maven`: Setup Java (`java-version`)<br>`unit-tests`: Setup Java (`java-version`)<br>`coveralls`: Setup Java (`java-version`)<br>`codecov`: Setup Java (`java-version`)<br>`sonar`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `build-cli`: Setup Java (`distribution`)<br>`build-tool`: Setup Java (`distribution`)<br>`build-ant`: Setup Java (`distribution`)<br>`build-gradle`: Setup Java (`distribution`)<br>`build-maven`: Setup Java (`distribution`)<br>`unit-tests`: Setup Java (`distribution`)<br>`coveralls`: Setup Java (`distribution`)<br>`codecov`: Setup Java (`distribution`)<br>`sonar`: Setup Java (`distribution`) |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main` (external) |

#### Secrets forwarded

- `github-token`: `${{ secrets.GITHUB_TOKEN }}`

### CLI (os) (`build-cli`)

| Property | Value |
|----------|-------|
| Matrix | `job.os`: macos-15-intel, ubuntu-latest, windows-latest; `job.args`: -xp docker,  |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (12)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@v5.0.3`
   - With:
     - `path`: `~/.m2/repository`
     - `key`: `setup-java-${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}`
     - `restore-keys`: `${{ runner.os }}-m2`

8. **Prepare**

9. **JReleaser**
   - Env:
     - `JRELEASER_OUTPUT_DIRECTORY`: `out`
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `JAVA_OPTS`: `-javaagent:jacoco/jacocoagent.jar=includes=*jreleaser*,destfile=jreleaser-cli-${{ runner.os }}.exec`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-cli-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/out/jreleaser/trace.log` ... (+3 more lines)

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-cli-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

</details>

### Tool (os) (`build-tool`)

| Property | Value |
|----------|-------|
| Matrix | `job.os`: macos-15-intel, ubuntu-latest, windows-latest; `job.args`: -xp docker,  |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (12)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@v5.0.3`
   - With:
     - `path`: `~/.m2/repository`
     - `key`: `setup-java-${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}`
     - `restore-keys`: `${{ runner.os }}-m2`

8. **Prepare**

9. **JReleaser**
   - Env:
     - `JRELEASER_OUTPUT_DIRECTORY`: `out`
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `JAVA_OPTS`: `-javaagent:jacoco/jacocoagent.jar=includes=*jreleaser*,destfile=jreleaser-tool-${{ runner.os }}.exec`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-tool-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/out/jreleaser/trace.log` ... (+3 more lines)

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-tool-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

</details>

### Ant (os) (`build-ant`)

| Property | Value |
|----------|-------|
| Matrix | `job.os`: macos-15-intel, ubuntu-latest, windows-latest; `job.args`: -Djreleaser.excluded.packagers=docker,  |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (12)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@v5.0.3`
   - With:
     - `path`: `~/.m2/repository`
     - `key`: `setup-java-${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}`
     - `restore-keys`: `${{ runner.os }}-m2`

8. **Prepare**

9. **JReleaser**
   - Env:
     - `JRELEASER_OUTPUT_DIRECTORY`: `out`
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `ANT_OPTS`: `-javaagent:jacoco/jacocoagent.jar=includes=*jreleaser*,destfile=jreleaser-ant-${{ runner.os }}.exec`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-ant-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/build/jreleaser/trace.log` ... (+3 more lines)

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-ant-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

</details>

### Gradle (os) (`build-gradle`)

| Property | Value |
|----------|-------|
| Matrix | `job.os`: macos-15-intel, ubuntu-latest, windows-latest; `job.args`: --exclude-packager docker,  |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (13)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@v5.0.3`
   - With:
     - `path`: `~/.m2/repository`
     - `key`: `setup-java-${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}`
     - `restore-keys`: `${{ runner.os }}-m2`

8. **Prepare**

9. **JReleaser**
   - Env:
     - `JRELEASER_OUTPUT_DIRECTORY`: `out`
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `JAVA_OPTS`: `-javaagent:jacoco/jacocoagent.jar=includes=*jreleaser*,destfile=jreleaser-gradle-${{ runner.os }}.exec`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-gradle-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/build/jreleaser/trace.log` ... (+3 more lines)

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-gradle-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Clean**
   - Env:
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`

13. **Cleanup**
   - Condition: `always()`

</details>

### Maven (os) (`build-maven`)

| Property | Value |
|----------|-------|
| Matrix | `job.os`: macos-15-intel, ubuntu-latest, windows-latest; `job.args`: -Djreleaser.excluded.packagers=docker,  |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (12)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@v5.0.3`
   - With:
     - `path`: `~/.m2/repository`
     - `key`: `setup-java-${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}`
     - `restore-keys`: `${{ runner.os }}-m2`

8. **Prepare**

9. **JReleaser**
   - Env:
     - `JRELEASER_OUTPUT_DIRECTORY`: `out`
     - `JRELEASER_USER_HOME`: `${{ github.workspace }}/smoketests-jreleaser/.jreleaser`
     - `JRELEASER_PROJECT_VERSION`: `1.0.0`
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.GIT_PAT_TOKEN }}`
     - `JRELEASER_GPG_PASSPHRASE`: `${{ secrets.GPG_PASSPHRASE }}`
     - `JRELEASER_GPG_PUBLIC_KEY`: `${{ secrets.GPG_PUBLIC_KEY }}`
     - `JRELEASER_GPG_SECRET_KEY`: `${{ secrets.GPG_SECRET_KEY }}`
     - `MAVEN_OPTS`: `-javaagent:jacoco/jacocoagent.jar=includes=*jreleaser*,destfile=jreleaser-maven-${{ runner.os }}.exec`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-maven-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/target/jreleaser/trace.log` ... (+3 more lines)

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-maven-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

</details>

### Unit Test (os) (`unit-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Matrix | `os`: ubuntu-latest, macos-15-intel, windows-latest |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (7)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Test**

5. **Rename JaCoCo execution data**

6. **JaCoCo upload**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-${{ runner.os }}`
     - `path`: `*.exec`

7. **Cleanup**
   - Condition: `always()`

</details>

### Coveralls (`coveralls`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `jacoco-*`
     - `merge-multiple`: `true`
     - `path`: `jacoco`

6. **JaCoCo merge**

7. **JaCoCo report**

8. **Upload coverage to Coveralls**
   - Env:
     - `COVERALLS_REPO_TOKEN`: `${{ secrets.COVERALLS_TOKEN }}`

9. **Cleanup**
   - Condition: `always()`

</details>

### Codecov (`codecov`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `jacoco-*`
     - `merge-multiple`: `true`
     - `path`: `jacoco`

6. **JaCoCo merge**

7. **JaCoCo report**

8. **Upload coverage to Codecov**
   - Uses: `codecov/codecov-action@v5.5.2`
   - With:
     - `token`: `${{ secrets.CODECOV_TOKEN }}`
     - `files`: `build/reports/jacoco/aggregate/jacocoTestReport.xml`
     - `flags`: `smoke-tests`
     - `fail_ci_if_error`: `false`
     - `name`: `jreleaser-smoke-tests`
     - `verbose`: `true`

9. **Cleanup**
   - Condition: `always()`

</details>

### Sonar (`sonar`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

<details>
<summary>Steps (9)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `pattern`: `jacoco-*`
     - `merge-multiple`: `true`
     - `path`: `jacoco`

6. **JaCoCo merge**

7. **JaCoCo report**

8. **Sonar**

9. **Cleanup**
   - Condition: `always()`

</details>

[Back to contents](#contents)

# Trigger Early Access

**Triggers:** `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `trigger-early-access.yml` |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GIT_ACCESS_TOKEN` | `earlyaccess`: Release early-access artifacts (`token`) |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | `earlyaccess`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `earlyaccess`: Setup Java (`distribution`) |

## Jobs

### Trigger Early Access (`earlyaccess`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6.0.2**
   - With:
     - `persist-credentials`: `false`

2. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`

3. **Build**

4. **Rename artifacts**

5. **Release early-access artifacts**
   - Uses: `softprops/action-gh-release@v2.6.1`
   - With:
     - `generate_release_notes`: `false`
     - `tag_name`: `early-access`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `prerelease`: `true`
     - `name`: `JReleaser Early-Access`
     - `files`: `early-access/*`

</details>

[Back to contents](#contents)

# X-BachInfo

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-update-bach-info.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `project-version` | string | Yes | - | - |
| `project-tag` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `gh-access-token` | Yes | - |

## Permissions

- `actions`: `read`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `gh-access-token` | `update-bach-info`: Checkout (`token`) |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMIT_EMAIL` | `update-bach-info`: Commit (`COMMIT_EMAIL`) |

## Jobs

### Update bach-info (`update-bach-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (3)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`
     - `repository`: `jreleaser/bach-info`
     - `ref`: `main`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.gh-access-token }}`

2. **Download script**

3. **Commit**
   - Env:
     - `TAG`: `${{ inputs.project-tag }}`
     - `VERSION`: `${{ inputs.project-version }}`
     - `COMMIT_EMAIL`: `${{ vars.COMMIT_EMAIL }}`

</details>

[Back to contents](#contents)

# X-Jlink

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-jlink.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `project-version` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `gpg-passphrase` | Yes | - |
| `oci-compartment-id` | Yes | - |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `gpg-passphrase` | `jlink`: Decrypt secrets (`GPG_PASSPHRASE`) |
| `oci-compartment-id` | `jlink`: Build (`JRELEASER_OCI_COMPARTMENTID`) |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | `jlink`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `jlink`: Setup Java (`distribution`) |

## Jobs

### Jlink (`jlink`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `CI` | `true` |

<details>
<summary>Steps (15)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Decrypt secrets**
   - Env:
     - `GPG_PASSPHRASE`: `${{ secrets.gpg-passphrase }}`

3. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Version**
   - Condition: `${{ endsWith(inputs.project-version, '-SNAPSHOT') != true }}`
   - Env:
     - `PROJECT_VERSION`: `${{ inputs.project-version }}`

5. **Build**
   - Env:
     - `JRELEASER_OCI_COMPARTMENTID`: `${{ secrets.oci-compartment-id }}`

6. **Clear space**

7. **Java Archive**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler java-archive`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

8. **Jlink**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler jlink`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

9. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-jlink`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

10. **Dependencies**

11. **Upload artifacts**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `1`
     - `name`: `artifacts`
     - `path`: `plugins/jreleaser/build/libs/` ... (+4 more lines)

12. **Upload java-archive**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `1`
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive/*.zip out/jreleaser/assemble/jreleaser/java-archive/*.tar`

13. **Upload jlink**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `1`
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink/*.zip`

14. **Stop Gradle daemon**

15. **Delete JDK caches**

</details>

[Back to contents](#contents)

# X-JPackage

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-jpackage.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `project-version` | string | Yes | - | - |
| `project-effective-version` | string | Yes | - | - |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | `jpackage`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `jpackage`: Setup Java (`distribution`) |

## Jobs

### `jpackage`

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.runner }}` |
| Matrix | `job.runner`: ubuntu-latest, ubuntu-22.04-arm, macos-15-intel, macos-15, windows-latest; `job.platform`: linux-x86_64, linux-aarch_64, osx-x86_64, osx-aarch_64, windows-x86_64; `job.platformReplaced`: linux-x86_64, linux-aarch64, osx-x86_64, osx-aarch64, windows-x86_64; `job.jdkOs`: Linux, LinuxArm, Osx, OsxArm, Windows |

<details>
<summary>Steps (12)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Download artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download jlink**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

4. **Expand jlink**
   - Env:
     - `PROJECT_EFFECTIVE_VERSION`: `${{ inputs.project-effective-version }}`

5. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

6. **Version**
   - Condition: `${{ endsWith(inputs.project-version, '-SNAPSHOT') != true }}`
   - Env:
     - `PROJECT_VERSION`: `${{ inputs.project-version }}`

7. **Jdks**

8. **Add msbuild to PATH**
   - Uses: `microsoft/setup-msbuild@v3`
   - Condition: `runner.os == 'Windows'`

9. **Jpackage**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler jpackage --select-current-platform`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-jpackage-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

11. **Upload jpackage**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `1`
     - `name`: `jpackage-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage/*.pkg` ... (+4 more lines)

12. **Stop Gradle daemon**

</details>

[Back to contents](#contents)

# X-NativeImage

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-native-image.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `project-version` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `gh-access-token` | Yes | - |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | `native-image`: Setup Graal (`github-token`) |

**Variables:**

| Name | Used by |
|------|---------|
| `GRAAL_JAVA_VERSION` | `native-image`: Setup Graal (`java-version`) |
| `JAVA_VERSION` | `native-image`: Setup Java (`java-version`) |
| `JAVA_DISTRO` | `native-image`: Setup Java (`distribution`) |

## Jobs

### `native-image`

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.runner }}` |
| Matrix | `job.runner`: ubuntu-latest, ubuntu-22.04-arm, macos-15-intel, macos-15, windows-latest; `job.jdkOs`: Linux, LinuxArm, Osx, OsxArm, Windows |

<details>
<summary>Steps (10)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Download artifacts**
   - Uses: `actions/download-artifact@v8.0.1`
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@v1.5.0`
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
     - `distribution`: `graalvm`

4. **Setup Java**
   - Uses: `actions/setup-java@v5.2.0`
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Version**
   - Condition: `${{ endsWith(inputs.project-version, '-SNAPSHOT') != true }}`
   - Env:
     - `PROJECT_VERSION`: `${{ inputs.project-version }}`

6. **Jdks**

7. **NativeImage**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler native-image --select-current-platform`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

8. **JReleaser output**
   - Uses: `actions/upload-artifact@v7.0.0`
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-native-image-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

9. **Upload native-image**
   - Uses: `actions/upload-artifact@v7.0.0`
   - With:
     - `retention-days`: `1`
     - `name`: `native-image-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image/*.zip`

10. **Stop Gradle daemon**

</details>

[Back to contents](#contents)

# X-Precheck

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-precheck.yml` |

## Workflow call API

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `version` | version | `${{ jobs.precheck.outputs.version }}` |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `github-token` | Yes | - |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `github-token` | `precheck`: Cancel previous run (`access_token`) |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'jreleaser/jreleaser' && startsWith(github.event.head_commit.message, 'Releasing version') != true` |

<details>
<summary>Steps (3)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `false`

2. **Cancel previous run**
   - Uses: `styfle/cancel-workflow-action@v0.13.1`
   - With:
     - `access_token`: `${{ secrets.github-token }}`

3. **Version**
   - ID: `vars`

</details>

[Back to contents](#contents)

# X-UpdateWiki

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `step-update-wiki.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `project-version` | string | Yes | - | - |
| `project-tag` | string | Yes | - | - |
| `commit-email` | string | Yes | - | - |
| `template-params` | string | No | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `gh-access-token` | Yes | - |

## Permissions

- `actions`: `read`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `gh-access-token` | `update-wiki`: Checkout wiki (`token`), Generate wiki page (`JRELEASER_GITHUB_TOKEN`) |

## Jobs

### Update wiki for Release (project-tag) (`update-wiki`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (5)</summary>

1. **Checkout**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

2. **Checkout wiki**
   - Uses: `actions/checkout@v6.0.2`
   - With:
     - `persist-credentials`: `true`
     - `repository`: `jreleaser/jreleaser.wiki`
     - `path`: `wiki`
     - `token`: `${{ secrets.gh-access-token }}`

3. **Download checksums**
   - Env:
     - `JRELEASER_PROJECT_TAG`: `${{ inputs.project-tag }}`

4. **Generate wiki page**
   - Uses: `jreleaser/release-action@v2.5.0`
   - With:
     - `version`: `early-access`
     - `arguments`: `template eval --changelog --input-file src/jreleaser/templates/wiki-release-page.md.tpl --target-directory wiki/Releases ${TEMPLATE_PARAMS}`
   - Env:
     - `JRELEASER_GITHUB_TOKEN`: `${{ secrets.gh-access-token }}`
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`
     - `TEMPLATE_PARAMS`: `${{ inputs.template-params }}`

5. **Commit**
   - Env:
     - `TAG`: `${{ inputs.project-tag }}`
     - `VERSION`: `${{ inputs.project-version }}`
     - `COMMIT_EMAIL`: `${{ inputs.commit-email }}`

</details>

[Back to contents](#contents)

