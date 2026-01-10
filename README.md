# ghalint

[![DeepWiki](https://img.shields.io/badge/Ask_DeepWiki-000000.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAyCAYAAAAnWDnqAAAAAXNSR0IArs4c6QAAA05JREFUaEPtmUtyEzEQhtWTQyQLHNak2AB7ZnyXZMEjXMGeK/AIi+QuHrMnbChYY7MIh8g01fJoopFb0uhhEqqcbWTp06/uv1saEDv4O3n3dV60RfP947Mm9/SQc0ICFQgzfc4CYZoTPAswgSJCCUJUnAAoRHOAUOcATwbmVLWdGoH//PB8mnKqScAhsD0kYP3j/Yt5LPQe2KvcXmGvRHcDnpxfL2zOYJ1mFwrryWTz0advv1Ut4CJgf5uhDuDj5eUcAUoahrdY/56ebRWeraTjMt/00Sh3UDtjgHtQNHwcRGOC98BJEAEymycmYcWwOprTgcB6VZ5JK5TAJ+fXGLBm3FDAmn6oPPjR4rKCAoJCal2eAiQp2x0vxTPB3ALO2CRkwmDy5WohzBDwSEFKRwPbknEggCPB/imwrycgxX2NzoMCHhPkDwqYMr9tRcP5qNrMZHkVnOjRMWwLCcr8ohBVb1OMjxLwGCvjTikrsBOiA6fNyCrm8V1rP93iVPpwaE+gO0SsWmPiXB+jikdf6SizrT5qKasx5j8ABbHpFTx+vFXp9EnYQmLx02h1QTTrl6eDqxLnGjporxl3NL3agEvXdT0WmEost648sQOYAeJS9Q7bfUVoMGnjo4AZdUMQku50McDcMWcBPvr0SzbTAFDfvJqwLzgxwATnCgnp4wDl6Aa+Ax283gghmj+vj7feE2KBBRMW3FzOpLOADl0Isb5587h/U4gGvkt5v60Z1VLG8BhYjbzRwyQZemwAd6cCR5/XFWLYZRIMpX39AR0tjaGGiGzLVyhse5C9RKC6ai42ppWPKiBagOvaYk8lO7DajerabOZP46Lby5wKjw1HCRx7p9sVMOWGzb/vA1hwiWc6jm3MvQDTogQkiqIhJV0nBQBTU+3okKCFDy9WwferkHjtxib7t3xIUQtHxnIwtx4mpg26/HfwVNVDb4oI9RHmx5WGelRVlrtiw43zboCLaxv46AZeB3IlTkwouebTr1y2NjSpHz68WNFjHvupy3q8TFn3Hos2IAk4Ju5dCo8B3wP7VPr/FGaKiG+T+v+TQqIrOqMTL1VdWV1DdmcbO8KXBz6esmYWYKPwDL5b5FA1a0hwapHiom0r/cKaoqr+27/XcrS5UwSMbQAAAABJRU5ErkJggg==)](https://deepwiki.com/suzuki-shunsuke/ghalint)
[Install](docs/install.md) | [Policies](#policies) | [How to use](#how-to-use) | [Configuration](#configuration)

GitHub Actions linter for security best practices.

```console
$ ghalint run
ERRO[0000] read a workflow file                          error="parse a workflow file as YAML: yaml: line 10: could not find expected ':'" program=ghalint version= workflow_file_path=.github/workflows/release.yaml
ERRO[0000] github.token should not be set to workflow's env  env_name=GITHUB_TOKEN policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
ERRO[0000] secret should not be set to workflow's env    env_name=DATADOG_API_KEY policy_name=workflow_secrets program=ghalint version= workflow_file_path=.github/workflows/test.yaml
```

ghalint is a command line tool to check GitHub Actions Workflows and action.yaml for security policy compliance.

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
1. [github_app_should_limit_repositories](docs/policies/009.md): GitHub Actions issuing GitHub Access tokens from GitHub Apps should limit repositories
1. [github_app_should_limit_permissions](docs/policies/010.md): GitHub Actions issuing GitHub Access tokens from GitHub Apps should limit permissions
1. [job_timeout_minutes_is_required](docs/policies/012.md): All jobs should set [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes)
1. [checkout_persist_credentials_should_be_false](docs/policies/013.md): [actions/checkout](https://github.com/actions/checkout)'s input `persist-credentials` should be `false`

### 2. Action Policies

1. [action_ref_should_be_full_length_commit_sha](docs/policies/008.md): action's ref should be full length commit SHA
1. [github_app_should_limit_repositories](docs/policies/009.md): GitHub Actions issuing GitHub Access tokens from GitHub Apps should limit repositories
1. [github_app_should_limit_permissions](docs/policies/010.md): GitHub Actions issuing GitHub Access tokens from GitHub Apps should limit permissions
1. [action_shell_is_required](docs/policies/011.md): `shell` is required if `run` is set
1. [checkout_persist_credentials_should_be_false](docs/policies/013.md): [actions/checkout](https://github.com/actions/checkout)'s input `persist-credentials` should be `false`

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

Then ghalint validates action files `^([^/]+/){0,3}action\.ya?ml$` on the current directory.
You can also specify file paths.

```sh
ghalint act foo/action.yaml bar/action.yml
```

## Configuration file

Configuration file path: `^(\.|\.github/)?ghalint\.ya?ml$`

You can specify the configuration file with the command line option `-config (-c)` or the environment variable `GHALINT_CONFIG`.

```sh
ghalint -c foo.yaml run
```

### JSON Schema

- [ghalint.json](json-schema/ghalint.json)
- https://raw.githubusercontent.com/suzuki-shunsuke/ghalint/refs/heads/main/json-schema/ghalint.json

If you look for a CLI tool to validate configuration with JSON Schema, [ajv-cli](https://ajv.js.org/packages/ajv-cli.html) is useful.

```sh
ajv --spec=draft2020 -s json-schema/ghalint.json -d ghalint.yaml
```

#### Input Complementation by YAML Language Server

[Please see the comment too.](https://github.com/szksh-lab/.github/issues/67#issuecomment-2564960491)

Version: `main`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghalint/main/json-schema/ghalint.json
```

Or pinning version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghalint/v1.2.1/json-schema/ghalint.json
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
- `GHALINT_LOG_LEVEL`: Log level One of `error`, `warn`, `info` (default), `debug`
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

## Experimental Features

> [!WARNING]
> These features are experimental, meaning they are unstable and may be changed or removed at minor or patch versions.

### Validate inputs of actions and reusable workflows

[#904](https://github.com/suzuki-shunsuke/ghalint/pull/904)

```console
$ ghalint exp validate-input
ERRO[0000] invalid input key                             action=suzuki-shunsuke/actionlint-action@c8d3c0dcc9152f1d1c7d4a38cbf4953c3a55953d input_key=actionlint_option job_key=actionlint program=ghalint required_inputs= valid_inputs="sparse-checkout, actionlint_options" version=v1.0.0-local workflow_file_path=.github/workflows/actionlint.yaml
```

`ghalint exp validate-input` command validates inputs of actions and reusable workflows.
It fails if required inputs aren't given or unknown inputs are passed.

> [!WARNING]
> [Actions using `required: true` will not automatically return an error if the input is not specified.](https://docs.github.com/en/actions/sharing-automations/creating-actions/metadata-syntax-for-github-actions#inputs)
> This means if `ghalint exp validate-input` fails as required inputs aren't given, the action may work without any problem.
> Now `ghalint exp validate-input` can't ignore those errors.
> Ideally, actions should be fixed.

By default, the following files are validated.

```
.github/workflows/*.yaml
.github/workflows/*.yml
action.yaml
action.yml
*/action.yaml
*/action.yml
*/*/action.yaml
*/*/action.yml
*/*/*/action.yaml
*/*/*/action.yml
```

This command uses a GitHub access token with `contents:read` permission to download actions and reusable workflows.
It downloads them into `XDG_DATA_HOME/ghalint`.
You can pass a GitHub access token by environment variables `GITHUB_TOKEN` or `GHALINT_GITHUB_TOKEN`.
You can also manage it by secret stores such as GNOME Keyring, Windows Credential Manager, and macOS Keychain.

```sh
ghalint exp token set [-stdin]
```

```sh
ghalint exp token rm # Remove a token from secret store
```

## LICENSE

[MIT](LICENSE)
