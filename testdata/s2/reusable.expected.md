# Reusable Build

| Property | Value |
|----------|-------|
| File | `reusable.yml` |
| Triggers | `workflow_call` |

## Workflow call API

This workflow is reusable via `workflow_call`.

**Inputs:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `version` | string | Yes | - | Semantic version to build. |
| `publish` | boolean | No | `false` | Whether to publish artifacts. |

**Outputs:**

| Name | Description | Value |
|------|-------------|-------|
| `artifact-url` | URL of the produced artifact. | `${{ jobs.build.outputs.url }}` |

**Secrets:**

| Name | Required | Description |
|------|----------|-------------|
| `CONSUMER-KEY` | Yes | SDKMAN consumer key for publishing. |
| `CONSUMER-TOKEN` | No | - |

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **Step 1**

