# Contents

- [Build](#build)
- [Clear cache](#clear-cache)
- [CodeQL](#codeql)
- [EarlyAccess](#earlyaccess)
- [Lint](#lint)
- [OpenSSF Scorecard](#openssf-scorecard)
- [Release](#release)
- [SmokeTests](#smoketests)
- [X-Jlink](#x-jlink)
- [X-JPackage](#x-jpackage)
- [X-NativeImage](#x-nativeimage)
- [X-Precheck](#x-precheck)
- [X-BachInfo](#x-bachinfo)
- [X-UpdateWiki](#x-updatewiki)
- [Trigger Early Access](#trigger-early-access)

# Build

| Property | Value |
|----------|-------|
| File | `build.yml` |
| Triggers | `pull_request` |

## Permissions

- `contents`: `read`

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`, cancel-in-progress: `true`

## Jobs

### Build (`build`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `21`
     - `distribution`: `zulu`
     - `cache`: `gradle`

3. **Build**

# Clear cache

| Property | Value |
|----------|-------|
| File | `cache.yml` |
| Triggers | `schedule`, `workflow_dispatch` |

## Schedule

- `0 3 * * *`

## Jobs

### Delete all caches (`clear`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Clear caches**
   - Uses: `easimon/wipe-cache@e7ab82e64c328fd39c2e96933d426cd72ac2beba`

# CodeQL

| Property | Value |
|----------|-------|
| File | `codeql.yml` |
| Triggers | `workflow_dispatch`, `push`, `pull_request` |

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

```
codeql.yml [workflow_dispatch, push, pull_request]
+-- precheck (uses jreleaser/jreleaser/.github/workflows/step-precheck.yml@main)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `precheck` secrets `github-token`; job `codeql` step `Cancel previous run` with `access_token` |

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

#### Steps

1. **Checkout repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Cancel previous run**
   - Uses: `styfle/cancel-workflow-action@d07a454dad7609a92316b57b23c9ccfd4f59af66` (v0.13.1)
   - With:
     - `access_token`: `${{ secrets.GITHUB_TOKEN }}`

3. **Initialize CodeQL**
   - Uses: `github/codeql-action/init@c10b8064de6f491fea524254123dbe5e09572f13` (v4.35.1)
   - With:
     - `languages`: `java`
     - `build-mode`: `manual`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `21`
     - `distribution`: `zulu`
     - `cache`: `gradle`

5. **Build**

6. **Perform CodeQL Analysis**
   - Uses: `github/codeql-action/analyze@c10b8064de6f491fea524254123dbe5e09572f13` (v4.35.1)
   - With:
     - `category`: `/language:java`

# EarlyAccess

| Property | Value |
|----------|-------|
| File | `early-access.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Permissions

- `actions`: `write`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Call graph (rooted at this workflow)

```
early-access.yml [push]
+-- precheck (uses jreleaser/jreleaser/.github/workflows/step-precheck.yml@main)
+-- jlink (uses jreleaser/jreleaser/.github/workflows/step-jlink.yml@main)
+-- jpackage (uses jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main)
+-- native-image (uses jreleaser/jreleaser/.github/workflows/step-native-image.yml@main)
+-- provenance (uses slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0)
+-- update-wiki (uses jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `gh-access-token`, `github-token`, `gpg-passphrase`, `oci-compartment-id`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`, `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`, `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`, `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`, `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`, `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `precheck` secrets `github-token` |
| `GPG_PASSPHRASE` | job `jlink` secrets `gpg-passphrase`; job `release` step `Release` env `JRELEASER_GPG_PASSPHRASE` |
| `JRELEASER_OCI_COMPARTMENTID` | job `jlink` secrets `oci-compartment-id` |
| `GIT_ACCESS_TOKEN` | job `native-image` secrets `gh-access-token`; job `release` step `Release` env `JRELEASER_GITHUB_TOKEN`; job `update-wiki` secrets `gh-access-token` |
| `GPG_PUBLIC_KEY` | job `release` step `Release` env `JRELEASER_GPG_PUBLIC_KEY` |
| `GPG_SECRET_KEY` | job `release` step `Release` env `JRELEASER_GPG_SECRET_KEY` |
| `JRELEASER_DOCKER_PASSWORD` | job `release` step `Release` env `JRELEASER_DOCKER_DEFAULT_PASSWORD` |

**Variables:**

| Name | Used by |
|------|---------|
| `GH_BOT_EMAIL` | job `update-wiki` with `commit-email` |

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

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Download artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download java-archive**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive`

4. **Download jlink**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

5. **Download jpackage**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `jpackage-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage`

6. **Download native-image**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `native-image-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image`

7. **Release**
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-release`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

9. **SLSA**
   - ID: `slsa`

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

# Lint

| Property | Value |
|----------|-------|
| File | `lint.yml` |
| Triggers | `push` |

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

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **actionlint**
   - ID: `actionlint`
   - Uses: `raven-actions/actionlint@205b530c5d9fa8f44ae9ed59f341a0db994aa6f8` (v2.1.2)

# OpenSSF Scorecard

| Property | Value |
|----------|-------|
| File | `openssf-scorecard.yml` |
| Triggers | `branch_protection_rule`, `schedule`, `push`, `workflow_dispatch` |

## Schedule

- `30 1 * * 6`

## Event filters

- **push**
  - branches: `main`

## Permissions

All scopes: `read-all`.

## Call graph (rooted at this workflow)

```
openssf-scorecard.yml [branch_protection_rule, schedule, push, workflow_dispatch]
+-- precheck (uses jreleaser/jreleaser/.github/workflows/step-precheck.yml@main)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | job `precheck` secrets `github-token` |

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

- `security-events`: `write` - Needed to upload the results to code-scanning dashboard.
- `id-token`: `write` (OIDC) - Used to receive a badge. (Upcoming feature)
- `actions`: `read`
- `contents`: `read`

#### Steps

1. **Checkout code**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Run analysis**
   - Uses: `ossf/scorecard-action@4eaacf0543bb3f2c246792bd56e8cdeffafb205a` (v2.4.3)
   - With:
     - `results_file`: `results.sarif`
     - `results_format`: `sarif`
     - `publish_results`: `true`

3. **Upload artifact**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `name`: `SARIF file`
     - `path`: `results.sarif`
     - `retention-days`: `5`

4. **Upload to code-scanning**
   - Uses: `github/codeql-action/upload-sarif@c10b8064de6f491fea524254123dbe5e09572f13` (v4.35.1)
   - With:
     - `sarif_file`: `results.sarif`

# Release

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Triggers | `workflow_dispatch` |

## Permissions

- `actions`: `write`
- `id-token`: `write` (OIDC)
- `contents`: `write`

## Call graph (rooted at this workflow)

```
release.yml [workflow_dispatch]
+-- jlink (uses jreleaser/jreleaser/.github/workflows/step-jlink.yml@main)
+-- jpackage (uses jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main)
+-- native-image (uses jreleaser/jreleaser/.github/workflows/step-native-image.yml@main)
+-- provenance (uses slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0)
+-- update-wiki (uses jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `gh-access-token`, `gpg-passphrase`, `oci-compartment-id`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-jlink.yml@main`, `jreleaser/jreleaser/.github/workflows/step-jpackage.yml@main`, `jreleaser/jreleaser/.github/workflows/step-native-image.yml@main`, `jreleaser/jreleaser/.github/workflows/step-update-wiki.yml@main`, `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `COMMIT_EMAIL` | job `precheck` step `Commit version` (run); job `release` step `Bump version` (run) |
| `GPG_PASSPHRASE` | job `jlink` secrets `gpg-passphrase`; job `release` step `Release` env `JRELEASER_GPG_PASSPHRASE` |
| `JRELEASER_OCI_COMPARTMENTID` | job `jlink` secrets `oci-compartment-id` |
| `GIT_ACCESS_TOKEN` | job `native-image` secrets `gh-access-token`; job `release` step `Release` env `JRELEASER_GITHUB_TOKEN`; job `update-wiki` secrets `gh-access-token`; job `update-website` step `Checkout` with `token` |
| `GRADLE_PUBLISH_KEY` | job `release` step `Deploy` env `GRADLE_PUBLISH_KEY` |
| `GRADLE_PUBLISH_SECRET` | job `release` step `Deploy` env `GRADLE_PUBLISH_SECRET` |
| `GPG_PUBLIC_KEY` | job `release` step `Release` env `JRELEASER_GPG_PUBLIC_KEY` |
| `GPG_SECRET_KEY` | job `release` step `Release` env `JRELEASER_GPG_SECRET_KEY` |
| `JRELEASER_DOCKER_PASSWORD` | job `release` step `Release` env `JRELEASER_DOCKER_DEFAULT_PASSWORD` |
| `SDKMAN_CONSUMER_KEY` | job `release` step `Release` env `JRELEASER_SDKMAN_CONSUMER_KEY` |
| `SDKMAN_CONSUMER_TOKEN` | job `release` step `Release` env `JRELEASER_SDKMAN_CONSUMER_TOKEN` |
| `MASTODON_ACCESS_TOKEN` | job `release` step `Release` env `JRELEASER_MASTODON_ACCESS_TOKEN` |
| `SONATYPE_USERNAME` | job `release` step `Release` env `JRELEASER_MAVENCENTRAL_USERNAME` |
| `SONATYPE_PASSWORD` | job `release` step `Release` env `JRELEASER_MAVENCENTRAL_PASSWORD` |
| `NOTICEABLE_APIKEY` | job `release` step `Release` env `JRELEASER_HTTP_NOTICEABLE_PASSWORD` |
| `OPENCOLLECTIVE_TOKEN` | job `release` step `Release` env `JRELEASER_OPENCOLLECTIVE_TOKEN` |
| `BLUESKY_HOST` | job `release` step `Release` env `JRELEASER_BLUESKY_HOST` |
| `BLUESKY_HANDLE` | job `release` step `Release` env `JRELEASER_BLUESKY_HANDLE` |
| `BLUESKY_PASSWORD` | job `release` step `Release` env `JRELEASER_BLUESKY_PASSWORD` |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | job `release` step `Setup Java` with `java-version`; job `update-website` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `release` step `Setup Java` with `distribution`; job `update-website` step `Setup Java` with `distribution` |
| `GH_BOT_EMAIL` | job `update-wiki` with `commit-email` |
| `COMMIT_EMAIL` | job `update-website` step `Commit` env `COMMIT_EMAIL` |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`

2. **Version**
   - ID: `version`

3. **Commit version**

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
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `jlink`, `jpackage`, `native-image` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`
     - `ref`: `main`
     - `fetch-depth`: `0`

2. **Download artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download java-archive**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive`

4. **Download jlink**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

5. **Download jpackage**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `jpackage-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage`

6. **Download native-image**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `native-image-*`
     - `merge-multiple`: `true`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image`

7. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

8. **Deploy**
   - Env:
     - `GRADLE_PUBLISH_KEY`: `${{ secrets.GRADLE_PUBLISH_KEY }}`
     - `GRADLE_PUBLISH_SECRET`: `${{ secrets.GRADLE_PUBLISH_SECRET }}`

9. **Upload deploy artifacts**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `7`
     - `name`: `deploy`
     - `path`: `build/repos/local/release/`

10. **Release**
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-release`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

12. **SLSA**
   - ID: `slsa`

13. **Bump version**

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
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `release` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`
     - `repository`: `jreleaser/jreleaser.github.io`
     - `ref`: `main`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

2. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
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

# SmokeTests

| Property | Value |
|----------|-------|
| File | `smoke-tests.yml` |
| Triggers | `push` |

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

```
smoke-tests.yml [push]
+-- precheck (uses jreleaser/jreleaser/.github/workflows/step-precheck.yml@main)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `github-token`

External workflows referenced: `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GPG_PASSPHRASE` | workflow env `GPG_PASSPHRASE`; job `build-cli` step `JReleaser` env `JRELEASER_GPG_PASSPHRASE`; job `build-tool` step `JReleaser` env `JRELEASER_GPG_PASSPHRASE`; job `build-ant` step `JReleaser` env `JRELEASER_GPG_PASSPHRASE`; job `build-gradle` step `JReleaser` env `JRELEASER_GPG_PASSPHRASE`; job `build-gradle` step `Clean` env `JRELEASER_GPG_PASSPHRASE`; job `build-maven` step `JReleaser` env `JRELEASER_GPG_PASSPHRASE` |
| `JRELEASER_OCI_COMPARTMENTID` | workflow env `JRELEASER_OCI_COMPARTMENTID` |
| `GITHUB_TOKEN` | job `precheck` secrets `github-token` |
| `GIT_ACCESS_TOKEN` | job `build-cli` step `Setup Graal` with `github-token`; job `build-cli` step `Checkout smoketests repository` with `token`; job `build-tool` step `Setup Graal` with `github-token`; job `build-tool` step `Checkout smoketests repository` with `token`; job `build-ant` step `Setup Graal` with `github-token`; job `build-ant` step `Checkout smoketests repository` with `token`; job `build-gradle` step `Setup Graal` with `github-token`; job `build-gradle` step `Checkout smoketests repository` with `token`; job `build-maven` step `Setup Graal` with `github-token`; job `build-maven` step `Checkout smoketests repository` with `token` |
| `GIT_PAT_TOKEN` | job `build-cli` step `JReleaser` env `JRELEASER_GITHUB_TOKEN`; job `build-tool` step `JReleaser` env `JRELEASER_GITHUB_TOKEN`; job `build-ant` step `JReleaser` env `JRELEASER_GITHUB_TOKEN`; job `build-gradle` step `JReleaser` env `JRELEASER_GITHUB_TOKEN`; job `build-gradle` step `Clean` env `JRELEASER_GITHUB_TOKEN`; job `build-maven` step `JReleaser` env `JRELEASER_GITHUB_TOKEN` |
| `GPG_PUBLIC_KEY` | job `build-cli` step `JReleaser` env `JRELEASER_GPG_PUBLIC_KEY`; job `build-tool` step `JReleaser` env `JRELEASER_GPG_PUBLIC_KEY`; job `build-ant` step `JReleaser` env `JRELEASER_GPG_PUBLIC_KEY`; job `build-gradle` step `JReleaser` env `JRELEASER_GPG_PUBLIC_KEY`; job `build-gradle` step `Clean` env `JRELEASER_GPG_PUBLIC_KEY`; job `build-maven` step `JReleaser` env `JRELEASER_GPG_PUBLIC_KEY` |
| `GPG_SECRET_KEY` | job `build-cli` step `JReleaser` env `JRELEASER_GPG_SECRET_KEY`; job `build-tool` step `JReleaser` env `JRELEASER_GPG_SECRET_KEY`; job `build-ant` step `JReleaser` env `JRELEASER_GPG_SECRET_KEY`; job `build-gradle` step `JReleaser` env `JRELEASER_GPG_SECRET_KEY`; job `build-gradle` step `Clean` env `JRELEASER_GPG_SECRET_KEY`; job `build-maven` step `JReleaser` env `JRELEASER_GPG_SECRET_KEY` |
| `COVERALLS_TOKEN` | job `coveralls` step `Upload coverage to Coveralls` env `COVERALLS_REPO_TOKEN` |
| `CODECOV_TOKEN` | job `codecov` step `Upload coverage to Codecov` with `token` |
| `SONARCLOUD_TOKEN` | job `sonar` step `Sonar` (run) |

**Variables:**

| Name | Used by |
|------|---------|
| `GRAAL_JAVA_VERSION` | job `build-cli` step `Setup Graal` with `java-version`; job `build-tool` step `Setup Graal` with `java-version`; job `build-ant` step `Setup Graal` with `java-version`; job `build-gradle` step `Setup Graal` with `java-version`; job `build-maven` step `Setup Graal` with `java-version` |
| `JAVA_VERSION` | job `build-cli` step `Setup Java` with `java-version`; job `build-tool` step `Setup Java` with `java-version`; job `build-ant` step `Setup Java` with `java-version`; job `build-gradle` step `Setup Java` with `java-version`; job `build-maven` step `Setup Java` with `java-version`; job `unit-tests` step `Setup Java` with `java-version`; job `coveralls` step `Setup Java` with `java-version`; job `codecov` step `Setup Java` with `java-version`; job `sonar` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `build-cli` step `Setup Java` with `distribution`; job `build-tool` step `Setup Java` with `distribution`; job `build-ant` step `Setup Java` with `distribution`; job `build-gradle` step `Setup Java` with `distribution`; job `build-maven` step `Setup Java` with `distribution`; job `unit-tests` step `Setup Java` with `distribution`; job `coveralls` step `Setup Java` with `distribution`; job `codecov` step `Setup Java` with `distribution`; job `sonar` step `Setup Java` with `distribution` |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Uses workflow | `jreleaser/jreleaser/.github/workflows/step-precheck.yml@main` (external) |

#### Secrets forwarded

- `github-token`: `${{ secrets.GITHUB_TOKEN }}`

### CLI macos-15-intel, ubuntu-latest, windows-latest (`build-cli`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@cdf6c1fa76f9f475f3d7449005a359c84ca0f306` (v5.0.3)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-cli-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/out/jreleaser/trace.log smoketests-jreleaser/out/jreleaser/output.properties smoketests-jreleaser/out/jreleaser/release/CHANGELOG.md smoketests-jreleaser/out/jreleaser/prepare`

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-cli-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

### Tool macos-15-intel, ubuntu-latest, windows-latest (`build-tool`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@cdf6c1fa76f9f475f3d7449005a359c84ca0f306` (v5.0.3)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-tool-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/out/jreleaser/trace.log smoketests-jreleaser/out/jreleaser/output.properties smoketests-jreleaser/out/jreleaser/release/CHANGELOG.md smoketests-jreleaser/out/jreleaser/prepare`

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-tool-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

### Ant macos-15-intel, ubuntu-latest, windows-latest (`build-ant`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@cdf6c1fa76f9f475f3d7449005a359c84ca0f306` (v5.0.3)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-ant-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/build/jreleaser/trace.log smoketests-jreleaser/build/jreleaser/output.properties smoketests-jreleaser/build/jreleaser/release/CHANGELOG.md smoketests-jreleaser/build/jreleaser/prepare`

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-ant-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

### Gradle macos-15-intel, ubuntu-latest, windows-latest (`build-gradle`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@cdf6c1fa76f9f475f3d7449005a359c84ca0f306` (v5.0.3)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-gradle-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/build/jreleaser/trace.log smoketests-jreleaser/build/jreleaser/output.properties smoketests-jreleaser/build/jreleaser/release/CHANGELOG.md smoketests-jreleaser/build/jreleaser/prepare`

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
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

### Maven macos-15-intel, ubuntu-latest, windows-latest (`build-maven`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `distribution`: `graalvm-community`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

5. **Build**

6. **Checkout smoketests repository**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `repository`: `jreleaser/smoketests-jreleaser`
     - `path`: `smoketests-jreleaser`
     - `fetch-depth`: `0`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`

7. **Cache Maven packages**
   - Uses: `actions/cache@cdf6c1fa76f9f475f3d7449005a359c84ca0f306` (v5.0.3)
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
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `7`
     - `name`: `jreleaser-maven-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/target/jreleaser/trace.log smoketests-jreleaser/target/jreleaser/output.properties smoketests-jreleaser/target/jreleaser/release/CHANGELOG.md smoketests-jreleaser/target/jreleaser/prepare`

11. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-maven-${{ runner.os }}`
     - `path`: `smoketests-jreleaser/*.exec`

12. **Cleanup**
   - Condition: `always()`

### Unit Test ubuntu-latest, macos-15-intel, windows-latest (`unit-tests`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.os }}` |
| Depends on | `precheck` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Test**

5. **Rename JaCoCo execution data**

6. **JaCoCo upload**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `retention-days`: `1`
     - `name`: `jacoco-${{ runner.os }}`
     - `path`: `*.exec`

7. **Cleanup**
   - Condition: `always()`

### Coveralls (`coveralls`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
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

### Codecov (`codecov`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `jacoco-*`
     - `merge-multiple`: `true`
     - `path`: `jacoco`

6. **JaCoCo merge**

7. **JaCoCo report**

8. **Upload coverage to Codecov**
   - Uses: `codecov/codecov-action@671740ac38dd9b0130fbe1cec585b89eea48d3de` (v5.5.2)
   - With:
     - `token`: `${{ secrets.CODECOV_TOKEN }}`
     - `files`: `build/reports/jacoco/aggregate/jacocoTestReport.xml`
     - `flags`: `smoke-tests`
     - `fail_ci_if_error`: `false`
     - `name`: `jreleaser-smoke-tests`
     - `verbose`: `true`

9. **Cleanup**
   - Condition: `always()`

### Sonar (`sonar`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `precheck`, `build-cli`, `build-tool`, `build-ant`, `build-gradle`, `build-maven`, `unit-tests` |
| Condition | `${{ endsWith(needs.precheck.outputs.version, '-SNAPSHOT') }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `fetch-depth`: `0`

2. **Decrypt secrets**

3. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`
     - `cache`: `gradle`

4. **Build**

5. **Download JaCoCo execution data**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `pattern`: `jacoco-*`
     - `merge-multiple`: `true`
     - `path`: `jacoco`

6. **JaCoCo merge**

7. **JaCoCo report**

8. **Sonar**

9. **Cleanup**
   - Condition: `always()`

# X-Jlink

| Property | Value |
|----------|-------|
| File | `step-jlink.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `gpg-passphrase` | job `jlink` step `Decrypt secrets` env `GPG_PASSPHRASE` |
| `oci-compartment-id` | job `jlink` step `Build` env `JRELEASER_OCI_COMPARTMENTID` |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | job `jlink` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `jlink` step `Setup Java` with `distribution` |

## Jobs

### Jlink (`jlink`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `CI` | `true` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Decrypt secrets**
   - Env:
     - `GPG_PASSPHRASE`: `${{ secrets.gpg-passphrase }}`

3. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
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
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler java-archive`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

8. **Jlink**
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler jlink`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

9. **JReleaser output**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-jlink`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

10. **Dependencies**

11. **Upload artifacts**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `1`
     - `name`: `artifacts`
     - `path`: `plugins/jreleaser/build/libs/ plugins/jreleaser/build/dependencies/ plugins/jreleaser/build/distributions/ plugins/jreleaser-tool-provider/build/libs/*.jar plugins/jreleaser-ant-tasks/build/distributions/*.zip`

12. **Upload java-archive**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `1`
     - `name`: `java-archive`
     - `path`: `out/jreleaser/assemble/jreleaser/java-archive/*.zip out/jreleaser/assemble/jreleaser/java-archive/*.tar`

13. **Upload jlink**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `1`
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink/*.zip`

14. **Stop Gradle daemon**

15. **Delete JDK caches**

# X-JPackage

| Property | Value |
|----------|-------|
| File | `step-jpackage.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `JAVA_VERSION` | job `jpackage` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `jpackage` step `Setup Java` with `distribution` |

## Jobs

### Linux, LinuxArm, Osx, OsxArm, Windows (`jpackage`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.runner }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Download artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Download jlink**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `jlink`
     - `path`: `out/jreleaser/assemble/jreleaser-standalone/jlink`

4. **Expand jlink**
   - Env:
     - `PROJECT_EFFECTIVE_VERSION`: `${{ inputs.project-effective-version }}`

5. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
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
   - Uses: `microsoft/setup-msbuild@30375c66a4eea26614e0d39710365f22f8b0af57` (v3)
   - Condition: `runner.os == 'Windows'`

9. **Jpackage**
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler jpackage --select-current-platform`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

10. **JReleaser output**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-jpackage-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

11. **Upload jpackage**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `1`
     - `name`: `jpackage-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/assemble/jreleaser-installer/jpackage/*.pkg out/jreleaser/assemble/jreleaser-installer/jpackage/*.msi out/jreleaser/assemble/jreleaser-installer/jpackage/*.exe out/jreleaser/assemble/jreleaser-installer/jpackage/*.deb out/jreleaser/assemble/jreleaser-installer/jpackage/*.rpm`

12. **Stop Gradle daemon**

# X-NativeImage

| Property | Value |
|----------|-------|
| File | `step-native-image.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `GITHUB_TOKEN` | job `native-image` step `Setup Graal` with `github-token` |

**Variables:**

| Name | Used by |
|------|---------|
| `GRAAL_JAVA_VERSION` | job `native-image` step `Setup Graal` with `java-version` |
| `JAVA_VERSION` | job `native-image` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `native-image` step `Setup Java` with `distribution` |

## Jobs

### Linux, LinuxArm, Osx, OsxArm, Windows (`native-image`)

| Property | Value |
|----------|-------|
| Runs on | `${{ matrix.job.runner }}` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`
     - `ref`: `main`

2. **Download artifacts**
   - Uses: `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (v8.0.1)
   - With:
     - `name`: `artifacts`
     - `path`: `plugins`

3. **Setup Graal**
   - Uses: `graalvm/setup-graalvm@f744c72a42b1995d7b0cbc314bde4bace7ac1fe1` (v1.5.0)
   - With:
     - `java-version`: `${{ vars.GRAAL_JAVA_VERSION }}`
     - `github-token`: `${{ secrets.GITHUB_TOKEN }}`
     - `distribution`: `graalvm`

4. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
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
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
   - With:
     - `version`: `early-access`
     - `arguments`: `assemble --assembler native-image --select-current-platform`
     - `setup-java`: `false`
   - Env:
     - `JRELEASER_PROJECT_VERSION`: `${{ inputs.project-version }}`

8. **JReleaser output**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - Condition: `always()`
   - With:
     - `name`: `jreleaser-native-image-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/trace.log out/jreleaser/output.properties`

9. **Upload native-image**
   - Uses: `actions/upload-artifact@bbbca2ddaa5d8feaa63e36b76fdaad77386f024f` (v7.0.0)
   - With:
     - `retention-days`: `1`
     - `name`: `native-image-${{ runner.os }}-${{ runner.arch }}`
     - `path`: `out/jreleaser/assemble/jreleaser-native/native-image/*.zip`

10. **Stop Gradle daemon**

# X-Precheck

| Property | Value |
|----------|-------|
| File | `step-precheck.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `github-token` | job `precheck` step `Cancel previous run` with `access_token` |

## Jobs

### Precheck (`precheck`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.repository == 'jreleaser/jreleaser' && startsWith(github.event.head_commit.message, 'Releasing version') != true` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Cancel previous run**
   - Uses: `styfle/cancel-workflow-action@d07a454dad7609a92316b57b23c9ccfd4f59af66` (v0.13.1)
   - With:
     - `access_token`: `${{ secrets.github-token }}`

3. **Version**
   - ID: `vars`

# X-BachInfo

| Property | Value |
|----------|-------|
| File | `step-update-bach-info.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `gh-access-token` | job `update-bach-info` step `Checkout` with `token` |

**Variables:**

| Name | Used by |
|------|---------|
| `COMMIT_EMAIL` | job `update-bach-info` step `Commit` env `COMMIT_EMAIL` |

## Jobs

### Update bach-info (`update-bach-info`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
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

# X-UpdateWiki

| Property | Value |
|----------|-------|
| File | `step-update-wiki.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

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
| `gh-access-token` | job `update-wiki` step `Checkout wiki` with `token`; job `update-wiki` step `Generate wiki page` env `JRELEASER_GITHUB_TOKEN` |

## Jobs

### Update wiki for Release ${{ inputs.project-tag }} (`update-wiki`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`
     - `fetch-depth`: `0`

2. **Checkout wiki**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `true`
     - `repository`: `jreleaser/jreleaser.wiki`
     - `path`: `wiki`
     - `token`: `${{ secrets.gh-access-token }}`

3. **Download checksums**
   - Env:
     - `JRELEASER_PROJECT_TAG`: `${{ inputs.project-tag }}`

4. **Generate wiki page**
   - Uses: `jreleaser/release-action@90ac653bb9c79d11179e65d81499f3f34527dcd5` (v2.5.0)
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

# Trigger Early Access

| Property | Value |
|----------|-------|
| File | `trigger-early-access.yml` |
| Triggers | `workflow_dispatch` |

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GIT_ACCESS_TOKEN` | job `earlyaccess` step `Release early-access artifacts` with `token` |

**Variables:**

| Name | Used by |
|------|---------|
| `JAVA_VERSION` | job `earlyaccess` step `Setup Java` with `java-version` |
| `JAVA_DISTRO` | job `earlyaccess` step `Setup Java` with `distribution` |

## Jobs

### Trigger Early Access (`earlyaccess`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **actions/checkout@v6.0.2**
   - Uses: `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2)
   - With:
     - `persist-credentials`: `false`

2. **Setup Java**
   - Uses: `actions/setup-java@be666c2fcd27ec809703dec50e508c2fdc7f6654` (v5.2.0)
   - With:
     - `java-version`: `${{ vars.JAVA_VERSION }}`
     - `distribution`: `${{ vars.JAVA_DISTRO }}`

3. **Build**

4. **Rename artifacts**

5. **Release early-access artifacts**
   - Uses: `softprops/action-gh-release@153bb8e04406b158c6c84fc1615b65b24149a1fe` (v2.6.1)
   - With:
     - `generate_release_notes`: `false`
     - `tag_name`: `early-access`
     - `token`: `${{ secrets.GIT_ACCESS_TOKEN }}`
     - `prerelease`: `true`
     - `name`: `JReleaser Early-Access`
     - `files`: `early-access/*`

