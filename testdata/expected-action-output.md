# Deploy Action

Deploy an application to the target environment.

| Property | Value |
|----------|-------|
| File | `action.yml` |
| Runs with | `node20` |
| Since | v1.0.0 |

**See also:** https://docs.example.com/actions/deploy

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|--------|
| `environment` | Target deployment environment | Yes | - |
| `version` | Version or commit SHA to deploy | Yes | - |
| `dry-run` | Run without making changes | No | `false` |

## Outputs

| Name | Description |
|------|-------------|
| `deploy-url` | URL of the deployed application |
| `deploy-id` | Unique identifier for this deployment |

## Secrets

| Name | Type | Description |
|------|------|-------------|
| `DEPLOY_TOKEN` | - | Token for authenticating with the deployment target |

## Environment Variables

| Name | Type | Description |
|------|------|-------------|
| `CLUSTER_NAME` | - | Kubernetes cluster name (must be set in repository variables) |

## Example

```
  - uses: my-org/deploy-action@v1
    with:
      environment: staging
      version: ${{ github.sha }}
```

