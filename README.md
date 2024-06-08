# ghalint

[Install](#how-to-install) | [Policies](#policies) | [How to use](#how-to-use) | [Configuration](#configuration)

GitHub Actions linter for security best practices.

```console
$ ghalint run
ERRO[0000] read a workflow file                          error="parse a workflow file as YAML: yaml: line 10: could not find expected ':'" program=ghalint version= workflow_file_path=.github/workflows/release.yaml
ERRO[0000] github.token should not be set to workflow's env  env_name=GITHUB_TOKEN policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
ERRO[0000] secret should not be set to workflow's env    env_name=DATADOG_API_KEY policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
```

ghalint is a command line tool to check GitHub Actions Workflows anc action.yaml for security policy compliance.

## :bulb: We've ported ghalint to lintnet module

- https://lintnet.github.io/
- https://github.com/lintnet-modules/ghalint

lintnet is a general purpose linter powered by Jsonnet.
We've ported ghalint to [the lintnet module](https://github.com/lintnet-modules/ghalint), so you can migrate ghalint to lintnet!

## Policies

### 1. Workflow Policies

1. [job_permissions](docs/policies/001.md): All jobs should have `permissions`
1. [deny_read_all_permission](docs/policies/002.md): `read-all` permission should not be used
1. [deny_write_all_permission](docs/policies/003.md): `write-all` permission should not be used
1. [deny_inherit_secrets](docs/policies/004.md): `secrets: inherit` should not be used
1. [workflow_secrets](docs/policies/005.md): Workflow should not set secrets to environment variables
1. [job_secrets](docs/policies/006.md): Job should not set secrets to environment variables
1. [deny_job_container_latest_image](docs/policies/007.md): Job's container image tag should not be `latest`
1. [action_ref_should_be_full_length_commit_sha](docs/policies/008.md): action's ref should be full length commit SHA
1. [github_app_should_limit_repositories](docs/policies/009.md): GitHub Actions issueing GitHub Access tokens from GitHub Apps should limit repositories
1. [github_app_should_limit_permissions](docs/policies/010.md): GitHub Actions issueing GitHub Access tokens from GitHub Apps should limit permissions

### 2. Action Policies

1. [action_ref_should_be_full_length_commit_sha](docs/policies/008.md): action's ref should be full length commit SHA
1. [github_app_should_limit_repositories](docs/policies/009.md): GitHub Actions issueing GitHub Access tokens from GitHub Apps should limit repositories
1. [github_app_should_limit_permissions](docs/policies/010.md): GitHub Actions issueing GitHub Access tokens from GitHub Apps should limit permissions
1. [action_shell_is_required](docs/policies/011.md): `shell` is required if `run` is set

## How to install

1. Homebrew:

```sh
brew install suzuki-shunsuke/ghalint/ghalint
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install ghalint
```

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/ghalint
```

4. [Download a pre-built binary from GitHub Releases](https://github.com/suzuki-shunsuke/ghalint/releases) and locate an executable binary `ghalint` in `PATH`

## How to use

### 1. Validate workflows

Run the command `ghalint run` on the repository root directory.

```sh
ghalint run
```

Then ghalint validates workflow files `^\.github/workflows/.*\.ya?ml$`.

### 2. Validate action.yaml

Run the command `ghalint run-action`.

```sh
ghalint run-action
```

The alias `act` is available.

```sh
ghalint act
```

Then ghalint validates action files `^action\.ya?ml$` on the current directory.
You can also specify file paths.

```sh
ghalint act foo/action.yaml bar/action.yml
```

## Configuration file

Configuration file path: `^\.?ghalint\.ya?ml$`

You can specify the configuration file with the command line option `-config (-c)` or the environment variable `GHALINT_CONFIG`.

```sh
ghalint -c foo.yaml run
```

### Disable policies

You can disable the following policies.

- [deny_inherit_secrets](docs/policies/004.md)
- [job_secrets](docs/policies/006.md)
- [action_ref_should_be_full_length_commit_sha](docs/policies/008.md)
- [github_app_should_limit_repositories](docs/policies/009.md)

e.g.

```yaml
excludes:
  - policy_name: deny_inherit_secrets
    workflow_file_path: .github/workflows/actionlint.yaml
    job_name: actionlint
  - policy_name: job_secrets
    workflow_file_path: .github/workflows/actionlint.yaml
    job_name: actionlint
  - policy_name: action_ref_should_be_full_length_commit_sha
    action_name: slsa-framework/slsa-github-generator
  - policy_name: github_app_should_limit_repositories
    workflow_file_path: .github/workflows/test.yaml
    job_name: test
    step_id: create_token
```

## Environment variables

- `GHALINT_CONFIG`: Configuration file path
- `GHALINT_LOG_LEVEL`: Log level One of `panic`, `fatal`, `error`, `warn`, `warning`, `info` (default), `debug`, `trace`
- `GHALINT_LOG_COLOR`: Configure log color. One of `auto` (default), `always`, and `never`.

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

## LICENSE

[MIT](LICENSE)
