# ghalint

GitHub Actions linter

## Policies

- `job_permissions`: All jobs should have `permissions` unless workflow's `permissions` is empty `{}`
  - Why: For least privilege
- `workflow_secrets`: Workflow should not set secrets to environment variables
  - How to fix: set secrets to jobs
  - Why: To limit the scope of secrets

### job_permissions

:x:

```yaml
jobs:
  hello:
    runs-on: ubuntu-latest
    # Without permissions
    steps:
      - run: echo hello
```

:o:

```yaml
jobs:
  hello:
    runs-on: ubuntu-latest
    permissions: {} # Set permissions
    steps:
      - run: echo hello
```

Or

```yaml
permissions: {} # Set permissions
jobs:
  hello:
    runs-on: ubuntu-latest
    steps:
      - run: echo hello
```

### workflow_secrets

:x:

```yaml
name: test
env:
  GITHUB_TOKEN: ${{github.token}}
  DATADOG_API_KEY: ${{secrets.DATADOG_API_KEY}}
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions: {}
    steps:
      - run: echo foo
  bar:
    runs-on: ubuntu-latest
    permissions: {}
    steps:
      - run: echo bar
```

:o:

```yaml
name: test
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions: {}
    env:
      GITHUB_TOKEN: ${{github.token}}
    steps:
      - run: echo foo
  bar:
    runs-on: ubuntu-latest
    permissions: {}
    env:
      DATADOG_API_KEY: ${{secrets.DATADOG_API_KEY}}
    steps:
      - run: echo bar
```

## How to install

- [Download a pre-built binary from GitHub Releases](https://github.com/suzuki-shunsuke/ghalint/releases) and locate an executable binary `ghalint` in `PATH`
- Homebrew: `brew install suzuki-shunsuke/ghalint/ghalint`
- [aqua](https://aquaproj.github.io/): `aqua g -i suzuki-shunsuke/ghalint`

## How to use

```console
$ ghalint run
```

```console
$ ghalint run
ERRO[0000] read a workflow file                          error="parse a workflow file as YAML: yaml: line 10: could not find expected ':'" program=ghalint version= workflow_file_path=.github/workflows/release.yaml
ERRO[0000] github.token should not be set to workflow's env  env_name=GITHUB_TOKEN policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
ERRO[0000] secret should not be set to workflow's env    env_name=DATADOG_API_KEY policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
```

## How does it works?

ghalint reads GitHub Actions Workflows `^\.github/workflows/.*\.ya?ml$` and validates them.
If there are violatation ghalint outputs error logs and fails.
If there is no violation ghalint succeeds.

## Why not `actionlint`?

We develop `ghalint` to support our policies that [actionlint](https://github.com/rhysd/actionlint) doesn't cover.
We don't aim to replace actionlint to ghalint. We use both `actionlint` and `ghalint`.

## LICENSE

[MIT](LICENSE)
