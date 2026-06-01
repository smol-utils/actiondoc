# Release

Orchestrates the full release: build, publish, and notify.

| Property | Value |
|----------|-------|
| File | `release.yml` |
| Triggers | `workflow_dispatch` |

## Manual trigger inputs

Inputs for the `workflow_dispatch` event.

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | The version to release |

## Call graph (rooted at this workflow)

```
release.yml [workflow_dispatch]
+-- publish (uses build_and_publish.yml)
|   +-- build-matrix (uses build.yml)
|   +-- publish / Set up toolchain (uses ./actions/setup)
+-- notify (uses external-org/notifications/.github/workflows/notify.yml@v2)
```

## Transitive requirements (from full call graph)

Secrets referenced (literal names): `GPG_KEY`, `SIGNING_KEY`

External workflows referenced: `external-org/notifications/.github/workflows/notify.yml@v2`

## Jobs

### `publish`

Build and publish all release artifacts.

| Property | Value |
|----------|-------|
| Uses workflow | [Build and Publish](#build-and-publish) |

#### Inputs forwarded

- `version`: `${{ inputs.version }}`
- `channel`: `stable`

#### Secrets forwarded

- `GPG_KEY`: `${{ secrets.RELEASE_GPG_KEY }}`

### `notify`

Announce the release in external channels.

| Property | Value |
|----------|-------|
| Uses workflow | `external-org/notifications/.github/workflows/notify.yml@v2` (external) |
| Depends on | `publish` |
| Condition | `success()` |

#### Secrets forwarded

- `secrets: inherit` (all caller secrets are passed to the callee)

### `tag`

Tag the release commit.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `publish` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Tag**

# Build and Publish

Builds the artifacts for every architecture and publishes them.

| Property | Value |
|----------|-------|
| File | `build_and_publish.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | - |
| `channel` | string | No | - | - |

## Called by

```
build_and_publish.yml
+-- release.yml (job: publish)  <- entry point
```

## Jobs

### `build-matrix`

Compile the artifacts for each target architecture.

| Property | Value |
|----------|-------|
| Uses workflow | [Build](#build) |

#### Inputs forwarded

- `arch`: `amd64,arm64`

### `publish`

Upload the built artifacts to the package registry.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |
| Depends on | `build-matrix` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Set up toolchain**
   - Uses: `./actions/setup`

3. **Publish**

# Build

Compiles the project for a set of architectures.

| Property | Value |
|----------|-------|
| File | `build.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `arch` | string | Yes | - | - |

## Called by

```
build.yml
+-- build_and_publish.yml (job: build-matrix)
    +-- release.yml (job: publish)  <- entry point
```

## Secrets

| Name | Type | Description |
|------|------|-------------|
| `SIGNING_KEY` | - | Key used to sign the compiled artifacts |

## Jobs

### `compile`

Compile and sign the artifacts.

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Checkout**
   - Uses: `actions/checkout@v4`

2. **Compile**

3. **Sign**

# Setup Toolchain

Installs the compilers and signing tools.

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `composite` |

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `cache` | Whether to cache the toolchain | No | `true` |

