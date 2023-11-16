# ghalint

GitHub Actions linter for security best practices.

[Blog post](https://dev.to/suzukishunsuke/minimize-the-scope-of-secrets-and-permissions-in-github-actions-444b)

## Policies

- [job_permissions](docs/policies/001.md): All jobs should have `permissions`
- [deny_read_all_permission](docs/policies/002.md): `read-all` permission should not be used
- [deny_write_all_permission](docs/policies/003.md): `write-all` permission should not be used
- [deny_inherit_secrets](docs/policies/004.md): `secrets: inherit` should not be used
- [workflow_secrets](docs/policies/005.md): Workflow should not set secrets to environment variables
- [job_secrets](docs/policies/006.md): Job should not set secrets to environment variables
- [deny_job_container_latest_image](docs/policies/007.md)

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
