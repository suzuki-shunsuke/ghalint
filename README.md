# ghalint

GitHub Actions linter

## Policies

- `job_permissions`: All jobs should have `permissions`
  - Why: For least privilege
- `workflow_secrets`: Workflow should not set secrets to environment variables
  - How to fix: set secrets to jobs
  - Why: To limit the scope of secrets

## How to use

```console
$ ghalint run
```

## Why not `actionlint`?

- We don't think ghalint can replace `actionlint`
- We use both `actionlint` and `ghalint`
- `ghalint` doesn't support features that `actionlint` supports
- We develop `ghalint` to support our policies that `actionlint` doesn't cover

## LICENSE

[MIT](LICENSE)
