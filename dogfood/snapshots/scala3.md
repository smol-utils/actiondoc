# scala3

11 workflows, 9 reusable workflows

## Contents

**Workflows**

- [Add to backporting project - push](#add-to-backporting-project)
- [Compile Full Standard Library - push, pull_request, workflow_call](#compile-full-standard-library)
- [Language reference documentation - push, pull_request, workflow_dispatch](#language-reference-documentation)
- [Nightly Release of Scala 3 - workflow_dispatch, schedule](#nightly-release-of-scala-3)
- [Official release of Scala - workflow_dispatch](#official-release-of-scala)
- [Scala 3 - push, pull_request, merge_group, workflow_dispatch](#scala-3)
- [Scala CLA - pull_request](#scala-cla)
- [scaladoc - push, pull_request, merge_group](#scaladoc)
- [Specification - push, pull_request, merge_group, workflow_dispatch](#specification)
- [Test CLI Launchers on all the platforms - pull_request, workflow_dispatch](#test-cli-launchers-on-all-the-platforms)
- [Update Dependency Graph - push](#update-dependency-graph)

**Reusable workflows**

- [Build 'scala' Chocolatey Package](#build-scala-chocolatey-package)
- [Build Scala Launchers](#build-scala-launchers)
- [Build the MSI Package](#build-the-msi-package)
- [Publish Scala to Chocolatey](#publish-scala-to-chocolatey)
- [Publish Scala to SDKMAN!](#publish-scala-to-sdkman)
- [Publish Scala to winget](#publish-scala-to-winget)
- [Release Artifacts to Maven](#release-artifacts-to-maven)
- [Test 'scala' Chocolatey Package](#test-scala-chocolatey-package)
- [Test 'scala' MSI Package](#test-scala-msi-package)

## Secrets and variables used across this repository

**Secrets:**

| Name | Used by |
|------|---------|
| `API-KEY` | [Publish Scala to Chocolatey](#publish-scala-to-chocolatey) |
| `CHOCOLATEY_KEY` | [Official release of Scala](#official-release-of-scala) |
| `CONSUMER-KEY` | [Publish Scala to SDKMAN!](#publish-scala-to-sdkman) |
| `CONSUMER-TOKEN` | [Publish Scala to SDKMAN!](#publish-scala-to-sdkman) |
| `DOCS_DEPLOY_KEY` | [Language reference documentation](#language-reference-documentation) |
| `DOCS_KEY` | [Language reference documentation](#language-reference-documentation) |
| `DOTTYBOT-TOKEN` | [Publish Scala to winget](#publish-scala-to-winget) |
| `DOTTYBOT_WINGET_TOKEN` | [Official release of Scala](#official-release-of-scala) |
| `MAVEN_REPOSITORY_TOKEN` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `MAVEN_REPOSITORY_USER` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `PGP_PW` | [Release Artifacts to Maven](#release-artifacts-to-maven), [Scala 3](#scala-3) |
| `PGP_SECRET` | [Release Artifacts to Maven](#release-artifacts-to-maven), [Scala 3](#scala-3) |
| `SCALA_APP_ID` | [Add to backporting project](#add-to-backporting-project), [Nightly Release of Scala 3](#nightly-release-of-scala-3) |
| `SCALA_APP_PRIVATE_KEY` | [Add to backporting project](#add-to-backporting-project), [Nightly Release of Scala 3](#nightly-release-of-scala-3) |
| `SCALA_PGP_KEY` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `SCALA_PGP_PASSPHRASE` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `SDKMAN_KEY` | [Official release of Scala](#official-release-of-scala) |
| `SDKMAN_TOKEN` | [Official release of Scala](#official-release-of-scala) |
| `SONATYPE_PW_ORGSCALALANG` | [Scala 3](#scala-3) |
| `SONATYPE_USER_ORGSCALALANG` | [Scala 3](#scala-3) |
| `SPEC_DEPLOY_HOST` | [Specification](#specification) |
| `SPEC_DEPLOY_KEY` | [Specification](#specification) |
| `SPEC_DEPLOY_PASS` | [Specification](#specification) |
| `SPEC_DEPLOY_PATH` | [Specification](#specification) |
| `SPEC_DEPLOY_USER` | [Specification](#specification) |

**Variables:**

| Name | Used by |
|------|---------|
| `MAVEN_REPOSITORY_HOST` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `MAVEN_REPOSITORY_REALM` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `MAVEN_REPOSITORY_URL` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `NEWNIGHTLY` | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `NIGHTLYBUILD` | [Nightly Release of Scala 3](#nightly-release-of-scala-3), [Release Artifacts to Maven](#release-artifacts-to-maven) |
| `SCALA_PGP_FINGERPRINT` | [Release Artifacts to Maven](#release-artifacts-to-maven) |

## Permissions across this repository

| Scope | Level |
|-------|-------|
| `contents` | `write` (also granted as `read` elsewhere) |

# Add to backporting project

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `lts-backport.yaml` |

## Event filters

- **push**
  - branches: `main`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SCALA_APP_ID` | `add-to-backporting-project`: Generate GitHub App Token (`app-id`) |
| `SCALA_APP_PRIVATE_KEY` | `add-to-backporting-project`: Generate GitHub App Token (`private-key`) |

## Jobs

### `add-to-backporting-project`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `!contains(github.event.push.head_commit.message, '[Next only]') && github.repository == 'scala/scala3'` |

<details>
<summary>Steps (5)</summary>

1. **Generate GitHub App Token**
   - ID: `app-token`
   - Uses: `actions/create-github-app-token@v3`
   - With:
     - `app-id`: `${{ secrets.SCALA_APP_ID }}`
     - `private-key`: `${{ secrets.SCALA_APP_PRIVATE_KEY }}`
     - `owner`: `scala`

2. **actions/checkout@v6**
   - With:
     - `fetch-depth`: `0`

3. **coursier/cache-action@v8**

4. **VirtusLab/scala-cli-setup@v1.14**

5. **scala-cli ./project/scripts/addToBackportingProject.scala...**
   - Env:
     - `GRAPHQL_API_TOKEN`: `${{ steps.app-token.outputs.token }}`

</details>

[Back to contents](#contents)

# Compile Full Standard Library

**Triggers:** `push`, `pull_request`, `workflow_call`

| Property | Value |
|----------|-------|
| File | `stdlib.yaml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:**

- [Non-Bootstrapped Library Unit Tests](#non-bootstrapped-library-unit-tests-test-scala-library-nonbootstrapped)
- [Bootstrapped Library Unit Tests](#bootstrapped-library-unit-tests-test-scala-library-bootstrapped)
- [`mima-scala-library-nonbootstrapped`](#mima-scala-library-nonbootstrapped)
- [`mima-scala3-interfaces`](#mima-scala3-interfaces)
- [`mima-tasty-core-nonbootstrapped`](#mima-tasty-core-nonbootstrapped)
- [`static-analysis-scala-library-bootstrapped`](#static-analysis-scala-library-bootstrapped)
- [`mima-tasty-core-bootstrapped`](#mima-tasty-core-bootstrapped)
- [`mima-scala-library-sjs`](#mima-scala-library-sjs)
- [`test-scala3-compiler-nonbootstrapped`](#test-scala3-compiler-nonbootstrapped)
- [`test-scala3-compiler-bootstrapped`](#test-scala3-compiler-bootstrapped)
- [`test-scala3-bootstrapped-compilation-coverage`](#test-scala3-bootstrapped-compilation-coverage)
- [`test-scala3-sbt-bridge-nonbootstrapped`](#test-scala3-sbt-bridge-nonbootstrapped)
- [`test-scala3-sbt-bridge-bootstrapped`](#test-scala3-sbt-bridge-bootstrapped)
- [`test-tasty-core-nonbootstrapped`](#test-tasty-core-nonbootstrapped)
- [`test-tasty-core-bootstrapped`](#test-tasty-core-bootstrapped)
- [`test-scala-js`](#test-scala-js)
- [`test-repl`](#test-repl)
- [`test-presentation-compiler`](#test-presentation-compiler)
- [`test-language-server`](#test-language-server)
- [`scripted-tests`](#scripted-tests)
- [`community_build_a`](#community_build_a)
- [`community_build_b`](#community_build_b)
- [`community_build_c`](#community_build_c)
- [`scala-library-docs`](#scala-library-docs)

## Event filters

- **push**
  - branches: `main`

## Permissions

- `contents`: `read`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `DOTTY_CI_RUN` | `true` |

## Called by

- [ci.yaml](#stdlib-tests-1) (job: `stdlib-tests`) - entry point
- [release-nightly.yml](#stdlib-tests) (job: `stdlib-tests`) - entry point

## Jobs

### Non-Bootstrapped Library Unit Tests (`test-scala-library-nonbootstrapped`)

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **./project/scripts/sbt scala-library-nonbootstrapped/test**

</details>

### Bootstrapped Library Unit Tests (`test-scala-library-bootstrapped`)

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **./project/scripts/sbt scala-library-bootstrapped/test**

</details>

### `mima-scala-library-nonbootstrapped`

| Property | Value |
|----------|-------|
| Condition | `false` |

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`scala-library-nonbootstrapped\`**

</details>

### `mima-scala3-interfaces`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`scala3-interfaces\`**

</details>

### `mima-tasty-core-nonbootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`tasty-core-nonbootstrapped\`**

</details>

### `static-analysis-scala-library-bootstrapped`

<details>
<summary>Steps (5)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`scala-library-bootstrapped\`**

5. **Report missingLink checks**

</details>

### `mima-tasty-core-bootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`tasty-core-bootstrapped\`**

</details>

### `mima-scala-library-sjs`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Report MiMa issues in \`scala-library-sjs\`**

</details>

### `test-scala3-compiler-nonbootstrapped`

<details>
<summary>Steps (5)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`scala3-compiler-nonbootstrapped\`**

5. **Cmd Tests**

</details>

### `test-scala3-compiler-bootstrapped`

<details>
<summary>Steps (5)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`scala3-compiler-bootstrapped\`**

5. **Cmd Tests**

</details>

### `test-scala3-bootstrapped-compilation-coverage`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Run \`scala3-bootstrapped/testCompilation --enable-coverage-phase\`**

</details>

### `test-scala3-sbt-bridge-nonbootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`scala3-sbt-bridge-nonbootstrapped\`**

</details>

### `test-scala3-sbt-bridge-bootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`scala3-sbt-bridge-bootstrapped\`**

</details>

### `test-tasty-core-nonbootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`tasty-core-nonbootstrapped\`**

</details>

### `test-tasty-core-bootstrapped`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test \`tasty-core-bootstrapped\`**

</details>

### `test-scala-js`

<details>
<summary>Steps (9)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **actions/setup-node@v6**
   - With:
     - `node-version`: `24.x`

4. **sbt/setup-sbt@v1**

5. **Scala.js compiler tests**

6. **Scala.js sandbox**

7. **Scala.js JUnit tests**

8. **Scala.js JUnit tests with latest ES version**

9. **Scala.js JUnit tests with WebAssembly**

</details>

### `test-repl`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test REPL**

</details>

### `test-presentation-compiler`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test Presentation Compiler**

</details>

### `test-language-server`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Test Language Server**

</details>

### `scripted-tests`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Run SBT scripted tests**

</details>

### `community_build_a`

<details>
<summary>Steps (4)</summary>

1. **Checkout cleanup script**
   - Uses: `actions/checkout@v6`
   - With:
     - `submodules`: `true`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Run Community Build A**

</details>

### `community_build_b`

<details>
<summary>Steps (4)</summary>

1. **Checkout cleanup script**
   - Uses: `actions/checkout@v6`
   - With:
     - `submodules`: `true`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Run Community Build B**

</details>

### `community_build_c`

<details>
<summary>Steps (4)</summary>

1. **Checkout cleanup script**
   - Uses: `actions/checkout@v6`
   - With:
     - `submodules`: `true`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Run Community Build C**

</details>

### `scala-library-docs`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Generate Documentation of the Standard Library**

</details>

[Back to contents](#contents)

# Language reference documentation

**Triggers:** `push`, `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `language-reference.yaml` |

## Event filters

- **push**
  - branches: `language-reference-stable`
- **pull_request**
  - branches: `language-reference-stable`

## Permissions

- `contents`: `read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOCS_KEY` | `build-and-push`: Git Checkout (`ssh-key`) |
| `DOCS_DEPLOY_KEY` | `build-and-push`: Push changes to scala3-reference-docs (`ssh-key`) |

## Jobs

### Build reference documentation and push it (`build-and-push`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

**Permissions:**

- `contents`: `write` - for Git to git push

<details>
<summary>Steps (7)</summary>

1. **Get current date**
   - ID: `date`

2. **Git Checkout**
   - Uses: `actions/checkout@v6`
   - With:
     - `path`: `dotty`
     - `fetch-depth`: `0`
     - `ssh-key`: `${{ secrets.DOCS_KEY }}`

3. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

4. **sbt/setup-sbt@v1**

5. **Generate reference documentation and test links**

6. **Push changes to scala3-reference-docs**
   - Uses: `actions/checkout@v6`
   - Condition: `github.event_name == 'push'`
   - With:
     - `repository`: `lampepfl/scala3-reference-docs`
     - `fetch-depth`: `0`
     - `submodules`: `true`
     - `ssh-key`: `${{ secrets.DOCS_DEPLOY_KEY }}`
     - `path`: `scala3-reference-docs`

7. **\cp -a dotty/scaladoc/output/reference/. scala3-reference...**
   - Condition: `github.event_name == 'push'`

</details>

[Back to contents](#contents)

# Nightly Release of Scala 3

**Triggers:** `workflow_dispatch`, `schedule`

| Property | Value |
|----------|-------|
| File | `release-nightly.yml` |

**Jobs:** [`stdlib-tests`](#stdlib-tests), [`release-maven-artifacts`](#release-maven-artifacts), [`release-documentation`](#release-documentation)

## Schedule

- `0 3 * * *` - Every day at 3 AM

## Call graph (rooted at this workflow)

- `stdlib-tests` uses [stdlib.yaml](#compile-full-standard-library)
- `release-maven-artifacts` uses [release-maven-artifacts.yml](#release-artifacts-to-maven)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SCALA_APP_ID` | `release-documentation`: Generate GitHub App Token (`app-id`) |
| `SCALA_APP_PRIVATE_KEY` | `release-documentation`: Generate GitHub App Token (`private-key`) |

**Variables:**

| Name | Used by |
|------|---------|
| `NIGHTLYBUILD` | `release-documentation`: (`NIGHTLYBUILD`) |

## Jobs

### `stdlib-tests`

| Property | Value |
|----------|-------|
| Uses workflow | [Compile Full Standard Library](#compile-full-standard-library) |

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `release-maven-artifacts`

| Property | Value |
|----------|-------|
| Uses workflow | [Release Artifacts to Maven](#release-artifacts-to-maven) |
| Depends on | `stdlib-tests` |

#### Inputs forwarded

- `environment`: `release-nightly`

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `release-documentation`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `release-maven-artifacts` |

**Deploys to environment:** `release-nightly` [gated]

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `NIGHTLYBUILD` | `${{ vars.NIGHTLYBUILD }}` |

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `adopt`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Generate The Website**

5. **Generate GitHub App Token**
   - ID: `app-token`
   - Uses: `actions/create-github-app-token@v3`
   - With:
     - `app-id`: `${{ secrets.SCALA_APP_ID }}`
     - `private-key`: `${{ secrets.SCALA_APP_PRIVATE_KEY }}`
     - `owner`: `scala`
     - `repositories`: `nightly.scala-lang.org`

6. **Deploy Website to https://nightly.scala-lang.org**
   - Uses: `peaceiris/actions-gh-pages@v4`
   - With:
     - `personal_token`: `${{ steps.app-token.outputs.token }}`
     - `publish_dir`: `scaladoc/output/scala3`
     - `external_repository`: `scala/nightly.scala-lang.org`
     - `publish_branch`: `main`

</details>

[Back to contents](#contents)

# Official release of Scala

**Triggers:** `workflow_dispatch`

OFFICIAL RELEASE WORKFLOW HOW TO USE: - THIS WORKFLOW WILL NEED TO BE TRIGGERED MANUALLY NOTE: - THIS WORKFLOW SHOULD ONLY BE RUN ON STABLE RELEASES - IT ASSUMES THAT THE PRE-RELEASE WORKFLOW WAS PREVIOUSLY EXECUTED

| Property | Value |
|----------|-------|
| File | `releases.yml` |

**Jobs:** [`publish-sdkman`](#publish-sdkman), [`publish-winget`](#publish-winget), [`compute-digest`](#compute-digest), [`build-chocolatey`](#build-chocolatey), [`test-chocolatey`](#test-chocolatey), [`publish-chocolatey`](#publish-chocolatey)

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | The version to officially release |

## Call graph (rooted at this workflow)

- `publish-sdkman` uses [publish-sdkman.yml](#publish-scala-to-sdkman)
- `publish-winget` uses [publish-winget.yml](#publish-scala-to-winget)
- `build-chocolatey` uses [build-chocolatey.yml](#build-scala-chocolatey-package)
- `test-chocolatey` uses [test-chocolatey.yml](#test-scala-chocolatey-package)
- `publish-chocolatey` uses [publish-chocolatey.yml](#publish-scala-to-chocolatey)

## Transitive requirements (from full call graph)

Secrets required (declared/forwarded names): `API-KEY`, `CONSUMER-KEY`, `CONSUMER-TOKEN`, `DOTTYBOT-TOKEN`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SDKMAN_KEY` | `publish-sdkman`: (`CONSUMER-KEY`) |
| `SDKMAN_TOKEN` | `publish-sdkman`: (`CONSUMER-TOKEN`) |
| `DOTTYBOT_WINGET_TOKEN` | `publish-winget`: (`DOTTYBOT-TOKEN`) |
| `CHOCOLATEY_KEY` | `publish-chocolatey`: (`API-KEY`) |

## Jobs

### `publish-sdkman`

| Property | Value |
|----------|-------|
| Uses workflow | [Publish Scala to SDKMAN!](#publish-scala-to-sdkman) |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`

#### Secrets forwarded

- `CONSUMER-KEY`: `${{ secrets.SDKMAN_KEY }}`
- `CONSUMER-TOKEN`: `${{ secrets.SDKMAN_TOKEN }}`

### `publish-winget`

| Property | Value |
|----------|-------|
| Uses workflow | [Publish Scala to winget](#publish-scala-to-winget) |
| Condition | `false` |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`

#### Secrets forwarded

- `DOTTYBOT-TOKEN`: `${{ secrets.DOTTYBOT_WINGET_TOKEN }}`

### `compute-digest`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (1)</summary>

1. **Compute the SHA256 of scala3-${{ inputs.version }}-x86\_64-pc-win32.zip in GitHub Release**
   - ID: `digest`
   - Env:
     - `VERSION`: `${{ inputs.version }}`

</details>

### `build-chocolatey`

| Property | Value |
|----------|-------|
| Uses workflow | [Build 'scala' Chocolatey Package](#build-scala-chocolatey-package) |
| Depends on | `compute-digest` |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`
- `url`: `https://github.com/scala/scala3/releases/download/${{ inputs.version }}/scala3-${{ inputs.version }}-x86_64-pc-win32.zip`
- `digest`: `${{ needs.compute-digest.outputs.digest }}`

### `test-chocolatey`

| Property | Value |
|----------|-------|
| Uses workflow | [Test 'scala' Chocolatey Package](#test-scala-chocolatey-package) |
| Depends on | `build-chocolatey` |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`
- `java-version`: `17`

### `publish-chocolatey`

| Property | Value |
|----------|-------|
| Uses workflow | [Publish Scala to Chocolatey](#publish-scala-to-chocolatey) |
| Depends on | `build-chocolatey`, `test-chocolatey` |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`

#### Secrets forwarded

- `API-KEY`: `${{ secrets.CHOCOLATEY_KEY }}`

[Back to contents](#contents)

# Scala 3

**Triggers:** `push`, `pull_request`, `merge_group`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `ci.yaml` |

**Jobs:**

- [`stdlib-tests`](#stdlib-tests-1)
- [`test_windows_full`](#test_windows_full)
- [`publish_release`](#publish_release)
- [`build-msi-package`](#build-msi-package)
- [`test-msi-package`](#test-msi-package)
- [`build-sdk-package`](#build-sdk-package)
- [`build-chocolatey-package`](#build-chocolatey-package)
- [`test-chocolatey-package`](#test-chocolatey-package)
- [`scalafmt`](#scalafmt)

## Event filters

- **push**
  - tags: `*`
  - branches-ignore: `gh-readonly-queue/**`, `release-**`, `lts-**`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `DOTTY_CI_RUN` | `true` |

**Concurrency:** group `${{ github.workflow }}-${{ github.ref }}`, cancel-in-progress: `${{ github.ref != 'refs/heads/main' }}`

## Call graph (rooted at this workflow)

- `stdlib-tests` uses [stdlib.yaml](#compile-full-standard-library)
- `build-msi-package` uses [build-msi.yml](#build-the-msi-package)
- `test-msi-package` uses [test-msi.yml](#test-scala-msi-package)
- `build-sdk-package` uses [build-sdk.yml](#build-scala-launchers)
- `build-chocolatey-package` uses [build-chocolatey.yml](#build-scala-chocolatey-package)
- `test-chocolatey-package` uses [test-chocolatey.yml](#test-scala-chocolatey-package)

## Transitive requirements (from full call graph)

Permissions declared across the chain: `contents: read`, `contents: write`

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `PGP_PW` | `publish_release`: (`PGP_PW`) |
| `PGP_SECRET` | `publish_release`: (`PGP_SECRET`) |
| `SONATYPE_PW_ORGSCALALANG` | `publish_release`: (`SONATYPE_PW`) |
| `SONATYPE_USER_ORGSCALALANG` | `publish_release`: (`SONATYPE_USER`) |
| `GITHUB_TOKEN` | `publish_release`: Create GitHub Release (`GITHUB_TOKEN`) |

## Jobs

### `stdlib-tests`

| Property | Value |
|----------|-------|
| Uses workflow | [Compile Full Standard Library](#compile-full-standard-library) |
| Condition | `github.event_name == 'push' && startsWith(github.event.ref, 'refs/tags/')` |

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `test_windows_full`

| Property | Value |
|----------|-------|
| Runs on | `self-hosted, Windows` |
| Condition | `github.event_name == 'schedule' && github.repository == 'scala/scala3' \|\| github.event_name == 'push' \|\| ( github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') && contains(github.event.pull_request.body, '[test_windows_full]') )` |

<details>
<summary>Steps (3)</summary>

1. **Reset existing repo**

2. **Git Checkout**
   - Uses: `actions/checkout@v6`

3. **Test**

</details>

### `publish_release`

| Property | Value |
|----------|-------|
| Runs on | `self-hosted, Linux` |
| Depends on | `stdlib-tests` |
| Condition | `github.event_name == 'push' && startsWith(github.event.ref, 'refs/tags/')` |

**Permissions:**

- `contents`: `write` - for GH CLI to create a release

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `RELEASEBUILD` | `yes` |
| `PGP_PW` | `${{ secrets.PGP_PW }}` |
| `PGP_SECRET` | `${{ secrets.PGP_SECRET }}` |
| `SONATYPE_PW` | `${{ secrets.SONATYPE_PW_ORGSCALALANG }}` |
| `SONATYPE_USER` | `${{ secrets.SONATYPE_USER_ORGSCALALANG }}` |

<details>
<summary>Steps (15)</summary>

1. **Set JDK 17 as default**

2. **Reset existing repo**

3. **Checkout cleanup script**
   - Uses: `actions/checkout@v6`

4. **Cleanup**

5. **Git Checkout**
   - Uses: `actions/checkout@v6`

6. **Add SBT proxy repositories**

7. **Extract the release tag**

8. **Check compiler version**

9. **Prepare the SDKs**

10. **Download MSI package**
   - Uses: `actions/download-artifact@v8`
   - With:
     - `name`: `scala.msi`
     - `path`: `.`

11. **Prepare MSI package**

12. **Install GH CLI**
   - Uses: `dev-hanz-ops/install-gh-cli-action@v0.2.1`
   - With:
     - `gh-cli-version`: `2.59.0`

13. **Create GitHub Release**
   - Env:
     - `GITHUB_TOKEN`: `${{ secrets.GITHUB_TOKEN }}`

14. **Publish Release (org.scala-lang)**

15. **Publish Release (org.scala-js)**

</details>

### `build-msi-package`

| Property | Value |
|----------|-------|
| Uses workflow | [Build the MSI Package](#build-the-msi-package) |
| Condition | `(github.event_name == 'pull_request' && contains(github.event.pull_request.body, '[test_msi]')) \|\| (github.event_name == 'push' && startsWith(github.event.ref, 'refs/tags/'))` |

### `test-msi-package`

| Property | Value |
|----------|-------|
| Uses workflow | [Test 'scala' MSI Package](#test-scala-msi-package) |
| Depends on | `build-msi-package` |

#### Inputs forwarded

- `version`: `${{ needs.build-msi-package.outputs.version }}`
- `java-version`: `17`

### `build-sdk-package`

| Property | Value |
|----------|-------|
| Uses workflow | [Build Scala Launchers](#build-scala-launchers) |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]')) \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3') \|\| (github.event_name == 'schedule' && github.repository == 'scala/scala3') \|\| github.event_name == 'push' \|\| github.event_name == 'merge_group'` |

#### Inputs forwarded

- `java-version`: `17`

### `build-chocolatey-package`

| Property | Value |
|----------|-------|
| Uses workflow | [Build 'scala' Chocolatey Package](#build-scala-chocolatey-package) |
| Depends on | `build-sdk-package` |

#### Inputs forwarded

- `version`: `3.6.0-SNAPSHOT`
- `url`: `https://api.github.com/repos/scala/scala3/actions/artifacts/${{ needs.build-sdk-package.outputs.win-x86_64-id }}/zip`
- `digest`: `${{ needs.build-sdk-package.outputs.win-x86_64-digest }}`

### `test-chocolatey-package`

| Property | Value |
|----------|-------|
| Uses workflow | [Test 'scala' Chocolatey Package](#test-scala-chocolatey-package) |
| Depends on | `build-chocolatey-package` |
| Condition | `github.event_name == 'pull_request' && contains(github.event.pull_request.body, '[test_chocolatey]')` |

#### Inputs forwarded

- `version`: `3.6.0-SNAPSHOT`
- `java-version`: `17`

### `scalafmt`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**
   - With:
     - `fetch-depth`: `0`

2. **coursier/cache-action@v8**

3. **VirtusLab/scala-cli-setup@v1.14**

4. **scala-cli format --check**

</details>

[Back to contents](#contents)

# Scala CLA

**Triggers:** `pull_request`

| Property | Value |
|----------|-------|
| File | `cla.yml` |

## Event filters

- **pull_request**
  - branches-ignore: `language-reference-stable`

## Jobs

### `check`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `github.event.pull_request.user.login != 'dependabot'` |

<details>
<summary>Steps (1)</summary>

1. **Verify CLA**
   - Uses: `scala/cla-checker@v1`
   - With:
     - `author`: `${{ github.event.pull_request.user.login }}`

</details>

[Back to contents](#contents)

# scaladoc

**Triggers:** `push`, `pull_request`, `merge_group`

| Property | Value |
|----------|-------|
| File | `scaladoc.yaml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`build`](#build), [`validate-docs`](#validate-docs), [`validate-generated-docs`](#validate-generated-docs), [`stdlib-sourcelinks-test`](#stdlib-sourcelinks-test), [`check-error-code-snippets`](#check-error-code-snippets)

## Event filters

- **push**
  - branches-ignore: `language-reference-stable`, `gh-readonly-queue/**`
- **pull_request**
  - branches-ignore: `language-reference-stable`

## Permissions

- `contents`: `read`

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Condition | `github.event_name == 'merge_group' \|\| (    github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') && !contains(github.event.pull_request.body, '[skip docs]') ) \|\| contains(github.event.ref, 'scaladoc') \|\| contains(github.event.ref, 'main')` |

<details>
<summary>Steps (11)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Compile and test scala3doc-js**

5. **Compile and test**

6. **Locally publish self**

7. **Generate self documentation**

8. **Generate testcases documentation**

9. **Generate reference documentation**

10. **Generate Scala 3 documentation**

11. **Upload scaladoc output**
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `name`: `scaladoc-output`
     - `path`: `scaladoc/output/scala3`

</details>

### `validate-docs`

<details>
<summary>Steps (4)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **coursier/cache-action@v8**

3. **VirtusLab/scala-cli-setup@v1.14**

4. **Validate docs sidebars**

</details>

### `validate-generated-docs`

| Property | Value |
|----------|-------|
| Depends on | `build` |

<details>
<summary>Steps (6)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`
   - With:
     - `fetch-depth`: `0`

2. **Download scaladoc output**
   - Uses: `actions/download-artifact@v8`
   - With:
     - `name`: `scaladoc-output`
     - `path`: `./scaladoc-output/`

3. **Set up Ruby**
   - Uses: `ruby/setup-ruby@v1`
   - With:
     - `ruby-version`: `3.3`

4. **Install html-proofer**

5. **Check if docs/ was modified**
   - ID: `docs-changed`

6. **Validate documentation links**

</details>

### `stdlib-sourcelinks-test`

| Property | Value |
|----------|-------|
| Condition | `false && ((    github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') && !contains(github.event.pull_request.body, '[skip docs]') ) \|\| contains(github.event.ref, 'scaladoc') \|\| contains(github.event.ref, 'scala3doc') \|\| contains(github.event.ref, 'main'))` |

<details>
<summary>Steps (3)</summary>

1. **Git Checkout**
   - Uses: `actions/checkout@v6`

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `distribution`: `temurin`
     - `java-version`: `17`

3. **Test sourcelinks to stdlib**

</details>

### `check-error-code-snippets`

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6**

2. **coursier/cache-action@v8**

3. **VirtusLab/scala-cli-setup@v1.14**
   - With:
     - `jvm`: `temurin:17`
     - `apps`: `sbt`

4. **Publish compiler locally**

5. **Test error code snippets**

6. **[On failure] Print reproduction/fix steps**
   - Condition: `failure()`

</details>

[Back to contents](#contents)

# Specification

**Triggers:** `push`, `pull_request`, `merge_group`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `spec.yml` |

## Event filters

- **push**
  - tags: `*`
  - branches-ignore: `gh-readonly-queue/**`

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `DOTTY_CI_RUN` | `true` |

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SPEC_DEPLOY_PATH` | `specification`: Deployment (`remote_path`) |
| `SPEC_DEPLOY_HOST` | `specification`: Deployment (`remote_host`) |
| `SPEC_DEPLOY_USER` | `specification`: Deployment (`remote_user`), Deployment (`USER_FOR_TEST`) |
| `SPEC_DEPLOY_KEY` | `specification`: Deployment (`remote_key`) |
| `SPEC_DEPLOY_PASS` | `specification`: Deployment (`remote_key_pass`) |

## Jobs

### `specification`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]'))  \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3') \|\| github.event_name == 'push' \|\| github.event_name == 'merge_group'` |

**Defaults:** working-directory `./docs/_spec`

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6**

2. **ruby/setup-ruby@v1**
   - With:
     - `ruby-version`: `2.7`

3. **Install required gems**

4. **Build the specification**

5. **Deployment**
   - Uses: `burnett01/rsync-deployments@8.0.5`
   - Condition: `${{ env.USER_FOR_TEST != '' }}`
   - With:
     - `switches`: `-rzv`
     - `path`: `docs/_spec/_site/`
     - `remote_path`: `${{ secrets.SPEC_DEPLOY_PATH }}`
     - `remote_host`: `${{ secrets.SPEC_DEPLOY_HOST }}`
     - `remote_user`: `${{ secrets.SPEC_DEPLOY_USER }}`
     - `remote_key`: `${{ secrets.SPEC_DEPLOY_KEY }}`
     - `remote_key_pass`: `${{ secrets.SPEC_DEPLOY_PASS }}`
   - Env:
     - `USER_FOR_TEST`: `${{ secrets.SPEC_DEPLOY_USER }}`

</details>

[Back to contents](#contents)

# Test CLI Launchers on all the platforms

**Triggers:** `pull_request`, `workflow_dispatch`

| Property | Value |
|----------|-------|
| File | `test-launchers.yml` |

**Jobs:** [Deploy and Test on Linux x64 architecture](#deploy-and-test-on-linux-x64-architecture-linux-x86_64), [Deploy and Test on Linux ARM64 architecture](#deploy-and-test-on-linux-arm64-architecture-linux-aarch64), [Deploy and Test on Mac x64 architecture](#deploy-and-test-on-mac-x64-architecture-mac-x86_64), [Deploy and Test on Mac ARM64 architecture](#deploy-and-test-on-mac-arm64-architecture-mac-aarch64), [Deploy and Test on Windows x64 architecture](#deploy-and-test-on-windows-x64-architecture-win-x86_64)

## Jobs

### Deploy and Test on Linux x64 architecture (`linux-x86_64`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') ) \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3' )` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `java-version`: `17`
     - `distribution`: `temurin`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build and test launcher command**
   - Env:
     - `LAUNCHER_EXPECTED_PROJECT`: `dist-linux-x86_64`

</details>

### Deploy and Test on Linux ARM64 architecture (`linux-aarch64`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-24.04-arm` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `java-version`: `17`
     - `distribution`: `temurin`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build and test launcher command**
   - Env:
     - `LAUNCHER_EXPECTED_PROJECT`: `dist-linux-aarch64`

</details>

### Deploy and Test on Mac x64 architecture (`mac-x86_64`)

| Property | Value |
|----------|-------|
| Runs on | `macos-15-intel` |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') ) \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3' )` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `java-version`: `17`
     - `distribution`: `temurin`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build and test launcher command**
   - Env:
     - `LAUNCHER_EXPECTED_PROJECT`: `dist-mac-x86_64`

</details>

### Deploy and Test on Mac ARM64 architecture (`mac-aarch64`)

| Property | Value |
|----------|-------|
| Runs on | `macos-latest` |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') ) \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3' )` |

<details>
<summary>Steps (4)</summary>

1. **actions/checkout@v6**

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `java-version`: `17`
     - `distribution`: `temurin`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build and test launcher command**
   - Env:
     - `LAUNCHER_EXPECTED_PROJECT`: `dist-mac-aarch64`

</details>

### Deploy and Test on Windows x64 architecture (`win-x86_64`)

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |
| Condition | `(github.event_name == 'pull_request' && !contains(github.event.pull_request.body, '[skip ci]') ) \|\| (github.event_name == 'workflow_dispatch' && github.repository == 'scala/scala3' )` |

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6**

2. **Set up JDK 17**
   - Uses: `actions/setup-java@v5`
   - With:
     - `java-version`: `17`
     - `distribution`: `temurin`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build the launcher command**

5. **Run the launcher command tests**

</details>

[Back to contents](#contents)

# Update Dependency Graph

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `dependency-graph.yml` |

## Event filters

- **push**
  - branches: `main`

## Jobs

### Update Dependency Graph (`dependency-graph`)

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (3)</summary>

1. **actions/checkout@v6**

2. **sbt/setup-sbt@v1**

3. **scalacenter/sbt-dependency-submission@v3**

</details>

[Back to contents](#contents)

# Build 'scala' Chocolatey Package

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO BUILD SCALA WITH CHOCOLATEY HOW TO USE: NOTE:

| Property | Value |
|----------|-------|
| File | `build-chocolatey.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |
| `url` | string | Yes | - | - |
| `digest` | string | Yes | - | - |

## Called by

- [ci.yaml](#build-chocolatey-package) (job: `build-chocolatey-package`) - entry point
- [releases.yml](#build-chocolatey) (job: `build-chocolatey`) - entry point

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

<details>
<summary>Steps (6)</summary>

1. **actions/checkout@v6**

2. **Replace the version placeholder**
   - Uses: `richardrigutins/replace-in-files@v3`
   - With:
     - `files`: `./pkgs/chocolatey/scala.nuspec`
     - `search-text`: `@LAUNCHER_VERSION@`
     - `replacement-text`: `${{ inputs.version }}`

3. **Replace the URL placeholder**
   - Uses: `richardrigutins/replace-in-files@v3`
   - With:
     - `files`: `./pkgs/chocolatey/tools/chocolateyInstall.ps1`
     - `search-text`: `@LAUNCHER_URL@`
     - `replacement-text`: `${{ inputs.url }}`

4. **Replace the CHECKSUM placeholder**
   - Uses: `richardrigutins/replace-in-files@v3`
   - With:
     - `files`: `./pkgs/chocolatey/tools/chocolateyInstall.ps1`
     - `search-text`: `@LAUNCHER_SHA256@`
     - `replacement-text`: `${{ inputs.digest }}`

5. **Build the Chocolatey package (.nupkg)**

6. **Upload the Chocolatey package to GitHub**
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `name`: `scala.nupkg`
     - `path`: `./pkgs/chocolatey/scala.${{ inputs.version }}.nupkg`
     - `if-no-files-found`: `error`

</details>

[Back to contents](#contents)

# Build Scala Launchers

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO BUILD THE SCALA LAUNCHERS HOW TO USE: - THSI WORKFLOW WILL PACKAGE THE ALL THE LAUNCHERS AND UPLOAD THEM TO GITHUB ARTIFACTS NOTE: - SEE THE WORFLOW FOR THE NAMES OF THE ARTIFACTS

| Property | Value |
|----------|-------|
| File | `build-sdk.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `java-version` | string | Yes | - | - |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `universal-id` | ID of the `universal` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.universal-id }}` |
| `linux-x86_64-id` | ID of the `linux x86-64` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.linux-x86_64-id }}` |
| `linux-aarch64-id` | ID of the `linux aarch64` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.linux-aarch64-id }}` |
| `mac-x86_64-id` | ID of the `mac x86-64` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.mac-x86_64-id }}` |
| `mac-aarch64-id` | ID of the `mac aarch64` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.mac-aarch64-id }}` |
| `win-x86_64-id` | ID of the `win x86-64` package from GitHub Artifacts (Authentication Required) | `${{ jobs.build.outputs.win-x86_64-id }}` |
| `win-x86_64-digest` | The SHA256 of the uploaded artifact (`win x86-64`) | `${{ jobs.build.outputs.win-x86_64-digest }}` |

## Called by

- [ci.yaml](#build-sdk-package) (job: `build-sdk-package`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | `build`: Compute SHA256 of the uploaded artifact (win x86-64) (`run`) |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

<details>
<summary>Steps (16)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `${{ inputs.java-version }}`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Build and pack the SDK (universal)**

5. **Build and pack the SDK (linux x86-64)**

6. **Build and pack the SDK (linux aarch64)**

7. **Build and pack the SDK (mac x86-64)**

8. **Build and pack the SDK (mac aarch64)**

9. **Build and pack the SDK (win x86-64)**

10. **Upload zip archive to GitHub Artifact (universal)**
   - ID: `universal`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/target/universal/stage`
     - `name`: `scala3-universal`

11. **Upload zip archive to GitHub Artifact (linux x86-64)**
   - ID: `linux-x86_64`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/linux-x86_64/target/universal/stage`
     - `name`: `scala3-x86_64-pc-linux`

12. **Upload zip archive to GitHub Artifact (linux aarch64)**
   - ID: `linux-aarch64`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/linux-aarch64/target/universal/stage`
     - `name`: `scala3-aarch64-pc-linux`

13. **Upload zip archive to GitHub Artifact (mac x86-64)**
   - ID: `mac-x86_64`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/mac-x86_64/target/universal/stage`
     - `name`: `scala3-x86_64-apple-darwin`

14. **Upload zip archive to GitHub Artifact (mac aarch64)**
   - ID: `mac-aarch64`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/mac-aarch64/target/universal/stage`
     - `name`: `scala3-aarch64-apple-darwin`

15. **Upload zip archive to GitHub Artifact (win x86-64)**
   - ID: `win-x86_64`
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `path`: `./dist/win-x86_64/target/universal/stage`
     - `name`: `scala3-x86_64-pc-win32`

16. **Compute SHA256 of the uploaded artifact (win x86-64)**
   - ID: `win-x86_64-digest`

</details>

[Back to contents](#contents)

# Build the MSI Package

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO BUILD SCALA MSI HOW TO USE: - THE RELEASE WORKFLOW SHOULD CALL THIS WORKFLOW - IT WILL UPLOAD TO GITHUB THE MSI FILE FOR SCALA UNDER THE 'scala.msi' NAME NOTE: - WE SHOULD BUILD SCALA USING JAVA 8

| Property | Value |
|----------|-------|
| File | `build-msi.yml` |

## Workflow call API

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `version` | The developed Scala version (major.minor.patch) extracted from project/Build.scala | `${{ jobs.build.outputs.version }}` |

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `RELEASEBUILD` | `${{ startsWith(github.event.ref, 'refs/tags/') && 'yes' \|\| 'no' }}` |

## Called by

- [ci.yaml](#build-msi-package) (job: `build-msi-package`) - entry point

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

<details>
<summary>Steps (5)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `adopt`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **Extract Scala version**
   - ID: `extract-version`

4. **Build MSI package**

5. **Upload MSI Artifact**
   - Uses: `actions/upload-artifact@v7`
   - With:
     - `name`: `scala.msi`
     - `path`: `./dist/win-x86_64/target/windows/scala.msi`

</details>

[Back to contents](#contents)

# Publish Scala to Chocolatey

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO PUBLISH SCALA TO CHOCOLATEY HOW TO USE: - THE RELEASE WORKFLOW SHOULD CALL THIS WORKFLOW - IT WILL PUBLISH TO CHOCOLATEY THE MSI NOTE: - WE SHOULD KEEP IN SYNC THE NAME OF THE MSI WITH THE ACTUAL BUILD - WE SHOULD KEEP IN SYNC THE URL OF THE RELEASE - IT ASSUMES THAT THE `build-chocolatey` WORKFLOW WAS EXECUTED BEFORE

| Property | Value |
|----------|-------|
| File | `publish-chocolatey.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `API-KEY` | Yes | - |

## Called by

- [releases.yml](#publish-chocolatey) (job: `publish-chocolatey`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `API-KEY` | `publish`: Publish the package to Chocolatey (`KEY`) |

## Jobs

### `publish`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

<details>
<summary>Steps (2)</summary>

1. **Fetch the Chocolatey package from GitHub**
   - Uses: `actions/download-artifact@v8`
   - With:
     - `name`: `scala.nupkg`

2. **Publish the package to Chocolatey**
   - Env:
     - `VERSION`: `${{ inputs.version }}`
     - `KEY`: `${{ secrets.API-KEY }}`

</details>

[Back to contents](#contents)

# Publish Scala to SDKMAN!

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO PUBLISH SCALA TO SDKMAN! HOW TO USE: - THE RELEASE WORKFLOW SHOULD CALL THIS WORKFLOW - IT WILL PUBLISH TO SDKMAN! THE BINARIES TO EACH SUPPORTED PLATFORM AND A UNIVERSAL JAR - IT CHANGES THE DEFAULT VERSION IN SDKMAN! NOTE: - WE SHOULD KEEP IN SYNC THE NAME OF THE ARCHIVES WITH THE ACTUAL BUILD - WE SHOULD KEEP IN SYNC THE URL OF THE RELEASE

| Property | Value |
|----------|-------|
| File | `publish-sdkman.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`publish`](#publish-1), [`default`](#default)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `CONSUMER-KEY` | Yes | - |
| `CONSUMER-TOKEN` | Yes | - |

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `RELEASE-URL` | `https://github.com/scala/scala3/releases/download/${{ inputs.version }}` |

## Called by

- [releases.yml](#publish-sdkman) (job: `publish-sdkman`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `CONSUMER-KEY` | `publish`: sdkman/sdkman-release-action (`CONSUMER-KEY`)<br>`default`: sdkman/sdkman-default-action (`CONSUMER-KEY`) |
| `CONSUMER-TOKEN` | `publish`: sdkman/sdkman-release-action (`CONSUMER-TOKEN`)<br>`default`: sdkman/sdkman-default-action (`CONSUMER-TOKEN`) |

## Jobs

### `publish`

| Property | Value |
|----------|-------|
| Matrix | `platform`: LINUX_64, LINUX_ARM64, MAC_OSX, MAC_ARM64, WINDOWS_64; `archive`: scala3-${{ inputs.version }}-x86_64-pc-linux.zip, scala3-${{ inputs.version }}-aarch64-pc-linux.zip, scala3-${{ inputs.version }}-x86_64-apple-darwin.zip, scala3-${{ inputs.version }}-aarch64-apple-darwin.zip, scala3-${{ inputs.version }}-x86_64-pc-win32.zip (combinations adjusted by include/exclude) |

<details>
<summary>Steps (1)</summary>

1. **sdkman/sdkman-release-action**
   - Uses: `sdkman/sdkman-release-action@f93e93c50d5c60902d0bf87b58aaa58c027094da`
   - With:
     - `CONSUMER-KEY`: `${{ secrets.CONSUMER-KEY }}`
     - `CONSUMER-TOKEN`: `${{ secrets.CONSUMER-TOKEN }}`
     - `CANDIDATE`: `scala`
     - `VERSION`: `${{ inputs.version }}`
     - `URL`: `${{ env.RELEASE-URL }}/${{ matrix.archive }}`
     - `PLATFORM`: `${{ matrix.platform }}`

</details>

### `default`

| Property | Value |
|----------|-------|
| Depends on | `publish` |

<details>
<summary>Steps (1)</summary>

1. **sdkman/sdkman-default-action**
   - Uses: `sdkman/sdkman-default-action@b3f991bd109e40155af1b13a4c6fc8e8ccada65e`
   - With:
     - `CONSUMER-KEY`: `${{ secrets.CONSUMER-KEY }}`
     - `CONSUMER-TOKEN`: `${{ secrets.CONSUMER-TOKEN }}`
     - `CANDIDATE`: `scala`
     - `VERSION`: `${{ inputs.version }}`

</details>

[Back to contents](#contents)

# Publish Scala to winget

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO PUBLISH SCALA TO WINGET HOW TO USE: - THE RELEASE WORKFLOW SHOULD CALL THIS WORKFLOW - IT WILL PUBLISH THE MSI TO WINGET NOTE: - WE SHOULD KEEP IN SYNC THE https://github.com/dottybot/winget-pkgs REPOSITORY

| Property | Value |
|----------|-------|
| File | `publish-winget.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `DOTTYBOT-TOKEN` | Yes | - |

## Called by

- [releases.yml](#publish-winget) (job: `publish-winget`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `DOTTYBOT-TOKEN` | `publish`: vedantmgoyal9/winget-releaser (`token`) |

## Jobs

### `publish`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

<details>
<summary>Steps (1)</summary>

1. **vedantmgoyal9/winget-releaser**
   - Uses: `vedantmgoyal9/winget-releaser@4ffc7888bffd451b357355dc214d43bb9f23917e`
   - With:
     - `identifier`: `Scala.Scala.3`
     - `version`: `${{ inputs.version }}`
     - `installers-regex`: `\.msi$`
     - `release-tag`: `${{ inputs.version }}`
     - `fork-user`: `dottybot`
     - `token`: `${{ secrets.DOTTYBOT-TOKEN }}`

</details>

[Back to contents](#contents)

# Release Artifacts to Maven

**Triggers:** `workflow_call`

| Property | Value |
|----------|-------|
| File | `release-maven-artifacts.yml` |
| Default runs-on | `ubuntu-latest` |

**Jobs:** [`release-maven-artifacts`](#release-maven-artifacts-1), [`release-maven-lts-artifacts`](#release-maven-lts-artifacts)

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `environment` | string | Yes | - | - |

## Called by

- [release-nightly.yml](#release-maven-artifacts) (job: `release-maven-artifacts`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `SCALA_PGP_KEY` | `release-maven-artifacts`: Import Scala PGP Key (`gpg_private_key`) |
| `SCALA_PGP_PASSPHRASE` | `release-maven-artifacts`: Import Scala PGP Key (`passphrase`) |
| `MAVEN_REPOSITORY_USER` | `release-maven-artifacts`: Publish Artifacts to the Maven Repository (`MAVEN_REPOSITORY_USER`)<br>`release-maven-lts-artifacts`: Publish Artifacts to the Maven Repository (`SONATYPE_USER`) |
| `MAVEN_REPOSITORY_TOKEN` | `release-maven-artifacts`: Publish Artifacts to the Maven Repository (`MAVEN_REPOSITORY_TOKEN`)<br>`release-maven-lts-artifacts`: Publish Artifacts to the Maven Repository (`SONATYPE_PW`) |
| `PGP_SECRET` | `release-maven-lts-artifacts`: Setup PGP Key (`PGP_SECRET`), Setup SBT PGP (`PGP_SECRET`), Publish Artifacts to the Maven Repository (`PGP_SECRET`) |
| `PGP_PW` | `release-maven-lts-artifacts`: Publish Artifacts to the Maven Repository (`PGP_PW`) |

**Variables:**

| Name | Used by |
|------|---------|
| `MAVEN_REPOSITORY_HOST` | `release-maven-artifacts`: (`MAVEN_REPOSITORY_HOST`)<br>`release-maven-lts-artifacts`: (`MAVEN_REPOSITORY_HOST`) |
| `MAVEN_REPOSITORY_REALM` | `release-maven-artifacts`: (`MAVEN_REPOSITORY_REALM`)<br>`release-maven-lts-artifacts`: (`MAVEN_REPOSITORY_REALM`) |
| `MAVEN_REPOSITORY_URL` | `release-maven-artifacts`: (`MAVEN_REPOSITORY_URL`)<br>`release-maven-lts-artifacts`: (`MAVEN_REPOSITORY_URL`) |
| `NEWNIGHTLY` | `release-maven-artifacts`: (`NEWNIGHTLY`)<br>`release-maven-lts-artifacts`: (`NEWNIGHTLY`) |
| `NIGHTLYBUILD` | `release-maven-artifacts`: (`NIGHTLYBUILD`)<br>`release-maven-lts-artifacts`: (`NIGHTLYBUILD`) |
| `SCALA_PGP_FINGERPRINT` | `release-maven-artifacts`: Import Scala PGP Key (`fingerprint`) |

## Jobs

### `release-maven-artifacts`

**Deploys to environment:** `${{ inputs.environment }}` [gated]

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `MAVEN_REPOSITORY_HOST` | `${{ vars.MAVEN_REPOSITORY_HOST }}` |
| `MAVEN_REPOSITORY_REALM` | `${{ vars.MAVEN_REPOSITORY_REALM }}` |
| `MAVEN_REPOSITORY_URL` | `${{ vars.MAVEN_REPOSITORY_URL }}` |
| `NEWNIGHTLY` | `${{ vars.NEWNIGHTLY }}` |
| `NIGHTLYBUILD` | `${{ vars.NIGHTLYBUILD }}` |

<details>
<summary>Steps (7)</summary>

1. **actions/checkout@v6**

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `adopt`
     - `java-version`: `17`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Import Scala PGP Key**
   - Uses: `crazy-max/ghaction-import-gpg@v7`
   - With:
     - `gpg_private_key`: `${{ secrets.SCALA_PGP_KEY }}`
     - `passphrase`: `${{ secrets.SCALA_PGP_PASSPHRASE }}`
     - `fingerprint`: `${{ vars.SCALA_PGP_FINGERPRINT }}`

5. **Get version string for this build**

6. **Check whether not yet published** `[continue-on-error]`
   - ID: `not_yet_published`

7. **Publish Artifacts to the Maven Repository**
   - Condition: `steps.not_yet_published.outcome == 'success'`
   - Env:
     - `MAVEN_REPOSITORY_USER`: `${{ secrets.MAVEN_REPOSITORY_USER }}`
     - `MAVEN_REPOSITORY_TOKEN`: `${{ secrets.MAVEN_REPOSITORY_TOKEN }}`

</details>

### `release-maven-lts-artifacts`

**Deploys to environment:** `${{ inputs.environment }}` [gated]

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `MAVEN_REPOSITORY_HOST` | `${{ vars.MAVEN_REPOSITORY_HOST }}` |
| `MAVEN_REPOSITORY_REALM` | `${{ vars.MAVEN_REPOSITORY_REALM }}` |
| `MAVEN_REPOSITORY_URL` | `${{ vars.MAVEN_REPOSITORY_URL }}` |
| `NEWNIGHTLY` | `${{ vars.NEWNIGHTLY }}` |
| `NIGHTLYBUILD` | `${{ vars.NIGHTLYBUILD }}` |

<details>
<summary>Steps (9)</summary>

1. **actions/checkout@v6**
   - With:
     - `repository`: `scala/scala3-lts`
     - `ref`: `lts-3.3`

2. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `8`
     - `cache`: `sbt`

3. **sbt/setup-sbt@v1**

4. **Setup PGP Key**
   - Env:
     - `PGP_SECRET`: `${{ secrets.PGP_SECRET }}`

5. **Setup SBT PGP directory**

6. **Setup SBT PGP**
   - Env:
     - `PGP_SECRET`: `${{ secrets.PGP_SECRET }}`

7. **Get version string for this build**

8. **Check whether not yet published** `[continue-on-error]`
   - ID: `not_yet_published`

9. **Publish Artifacts to the Maven Repository**
   - Condition: `steps.not_yet_published.outcome == 'success'`
   - Env:
     - `SONATYPE_USER`: `${{ secrets.MAVEN_REPOSITORY_USER }}`
     - `SONATYPE_PW`: `${{ secrets.MAVEN_REPOSITORY_TOKEN }}`
     - `PGP_PW`: `${{ secrets.PGP_PW }}`
     - `PGP_SECRET`: `${{ secrets.PGP_SECRET }}`

</details>

[Back to contents](#contents)

# Test 'scala' Chocolatey Package

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO TEST SCALA WITH CHOCOLATEY HOW TO USE: NOTE:

| Property | Value |
|----------|-------|
| File | `test-chocolatey.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |
| `java-version` | string | Yes | - | - |

## Environment (`env`)

| Variable | Value |
|----------|-------|
| `CHOCOLATEY-REPOSITORY` | `chocolatey-pkgs` |
| `DOTTY_CI_INSTALLATION` | `${{ endsWith(inputs.version, '-SNAPSHOT') && secrets.GITHUB_TOKEN \|\| '' }}` |

## Called by

- [ci.yaml](#test-chocolatey-package) (job: `test-chocolatey-package`) - entry point
- [releases.yml](#test-chocolatey) (job: `test-chocolatey`) - entry point

## Referenced secrets and variables

**Secrets:**

| Name | Used by |
|------|---------|
| `GITHUB_TOKEN` | workflow env: (`DOTTY_CI_INSTALLATION`) |

## Jobs

### `test`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

<details>
<summary>Steps (7)</summary>

1. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `${{ inputs.java-version }}`

2. **Download the 'nupkg' from GitHub Artifacts**
   - Uses: `actions/download-artifact@v8`
   - With:
     - `name`: `scala.nupkg`
     - `path`: `${{ env.CHOCOLATEY-REPOSITORY }}`

3. **Install the \`scala\` package with Chocolatey**

4. **Test the \`scala\` command**

5. **Test the \`scalac\` command**

6. **Test the \`scaladoc\` command**

7. **Uninstall the \`scala\` package**

</details>

[Back to contents](#contents)

# Test 'scala' MSI Package

**Triggers:** `workflow_call`

THIS IS A REUSABLE WORKFLOW TO TEST SCALA WITH MSI RUNNER HOW TO USE: Provide optional `version` to test if installed binaries are installed with correct Scala version.

<details>
<summary>more</summary>

NOTE: Requires `scala.msi` artifact uploaded within the same run

</details>

| Property | Value |
|----------|-------|
| File | `test-msi.yml` |

## Workflow call API

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |
| `java-version` | string | Yes | - | - |

## Called by

- [ci.yaml](#test-msi-package) (job: `test-msi-package`) - entry point

## Jobs

### `test`

| Property | Value |
|----------|-------|
| Runs on | `windows-latest` |

**Environment (`env`):**

| Variable | Value |
|----------|-------|
| `VERSION` | `${{ inputs.version }}` |

<details>
<summary>Steps (9)</summary>

1. **actions/setup-java@v5**
   - With:
     - `distribution`: `temurin`
     - `java-version`: `${{ inputs.java-version }}`

2. **Download MSI artifact**
   - Uses: `actions/download-artifact@v8`
   - With:
     - `name`: `scala.msi`
     - `path`: `.`

3. **Install Scala Runner**

4. **Verify installation layout**

5. **Test Scala Runner**

6. **Test the \`scalac\` command**

7. **Test the \`scaladoc\` command**

8. **Smoke test - compile and run**

9. **Uninstall the \`scala\` package**

</details>

[Back to contents](#contents)

