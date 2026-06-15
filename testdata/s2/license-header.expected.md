# Licensed CI

**Triggers:** `push`

| Property | Value |
|----------|-------|
| File | `license-header.yml` |

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

