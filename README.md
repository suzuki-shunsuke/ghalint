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

We develop `ghalint` to support our policies that [actionlint](https://github.com/rhysd/actionlint) doesn't cover.
We don't aim to replace actionlint to ghalint. We use both `actionlint` and `ghalint`.

## LICENSE

[MIT](LICENSE)
