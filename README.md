# ghalint

GitHub Actions linter

[Blog post](https://dev.to/suzukishunsuke/minimize-the-scope-of-secrets-and-permissions-in-github-actions-444b)

## Policies

- `job_permissions`: All jobs should have `permissions`
  - Why: For least privilege
  - Exceptions
    - workflow's `permissions` is empty `{}`
    - workflow has only one job and the workflow has `permissions`
- `deny_read_all_permission`: `read-all` permission should not be used
  - Why: For least privilege
- `deny_write_all_permission`: `write-all` permission should not be used
  - Why: For least privilege
- `workflow_secrets`: Workflow should not set secrets to environment variables
  - How to fix: set secrets to jobs
  - Why: To limit the scope of secrets
  - Exceptions
    - workflow has only one job
- `job_secrets`: Job should not set secrets to environment variables
  - How to fix: set secrets to steps
  - Why: To limit the scope of secrets
  - Exceptions
    - job has only one step

### job_permissions

:x:

```yaml
permissions:
  contents: read
jobs:
  foo:
    runs-on: ubuntu-latest
    # Without permissions
    steps:
      - run: echo hello
  bar:
    runs-on: ubuntu-latest
    # Without permissions
    steps:
      - uses: actions/checkout@v3
```

:o:

```yaml
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions: {} # Set permissions
    steps:
      - run: echo hello
  bar:
    runs-on: ubuntu-latest
    permissions: # Set permissions
      contents: read
    steps:
      - uses: actions/checkout@v3
```

Or

```yaml
permissions: {} # empty permissions
jobs:
  foo:
    runs-on: ubuntu-latest
    steps:
      - run: echo hello
  bar:
    runs-on: ubuntu-latest
    permissions: # Set permissions
      contents: read
    steps:
      - uses: actions/checkout@v3
```

### deny_read_all_permission

:x:

```yaml
name: test
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions: read-all # Don't use read-all
    steps:
      - run: echo foo
```

:o:

```yaml
name: test
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - run: echo foo
```

### deny_write_all_permission

Same with `deny_read_all_permission`.

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

### job_secrets

:x:

```yaml
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions:
      issues: write
    env:
      GITHUB_TOKEN: ${{github.token}} # secret is set in job
    steps:
      - run: echo foo
      - run: gh label create bug
```

:o:

```yaml
jobs:
  foo:
    runs-on: ubuntu-latest
    permissions:
      issues: write
    steps:
      - run: echo foo
      - run: gh label create bug
        env:
          GITHUB_TOKEN: ${{github.token}} # secret is set in step
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

## Configuration file

Configuration file path: `^\.?ghalint\.ya?ml$`

You can exclude the policy `job_secrets`.

e.g.

```yaml
excludes:
  - policy_name: job_secrets
    workflow_file_path: .github/workflows/actionlint.yaml
    job_name: actionlint
```

* policy_name: Only `job_secrets` is supported

## Environment variables

* `GHALINT_LOG_COLOR`: Configure log color. One of `auto` (default), `always`, and `never`.

💡 If you want to enable log color in GitHub Actions, please try `GHALINT_LOG_COLOR=always` 

```yaml
env:
  GHALINT_LOG_COLOR: always
```

AS IS

<img width="986" alt="image" src="https://user-images.githubusercontent.com/13323303/216190768-cb09597f-5669-4907-b443-78d96b4491ab.png">

TO BE

<img width="1023" alt="image" src="https://user-images.githubusercontent.com/13323303/216190842-0c015088-dda2-4e6f-8dbe-2db89cfbf438.png">

## How does it works?

ghalint reads GitHub Actions Workflows `^\.github/workflows/.*\.ya?ml$` and validates them.
If there are violatation ghalint outputs error logs and fails.
If there is no violation ghalint succeeds.

## Why not `actionlint`?

We develop `ghalint` to support our policies that [actionlint](https://github.com/rhysd/actionlint) doesn't cover.
We don't aim to replace actionlint to ghalint. We use both `actionlint` and `ghalint`.

## LICENSE

[MIT](LICENSE)
