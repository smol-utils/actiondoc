# Licensed CI

| Property | Value |
|----------|-------|
| File | `license-header.yml` |
| Triggers | `push` |

## Event filters

- **push**
  - branches: `main`

## Jobs

### `build`

| Property | Value |
|----------|-------|
| Runs on | `ubuntu-latest` |

#### Steps

1. **make build**

