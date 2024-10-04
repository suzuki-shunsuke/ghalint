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
1. [job_timeout_minutes_is_required](docs/policies/012.md): All jobs should set [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes)

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

<details>
<summary>Verify downloaded assets from GitHub Releases</summary>

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

### 1. GitHub CLI

ghalint >= v1.0.0

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
gh release download -R suzuki-shunsuke/ghalint v1.0.0 -p ghalint_1.0.0_darwin_arm64.tar.gz
gh attestation verify ghalint_1.0.0_darwin_arm64.tar.gz \
  -R suzuki-shunsuke/ghalint \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

Output:

```
Loaded digest sha256:3e3fda71ffae83cf713295df2bef09fc268811deab11dea58d8caa287642c9dc for file://ghalint_1.0.0_darwin_arm64.tar.gz
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:3e3fda71ffae83cf713295df2bef09fc268811deab11dea58d8caa287642c9dc was attested by:
REPO                                 PREDICATE_TYPE                  WORKFLOW
suzuki-shunsuke/go-release-workflow  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@7f97a226912ee2978126019b1e95311d7d15c97a
```

### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
gh release download -R suzuki-shunsuke/ghalint v1.0.0
slsa-verifier verify-artifact ghalint_1.0.0_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ghalint \
  --source-tag v1.0.0
```

Output:

```
Verified signature against tlog entry index 137012838 at URL: https://rekor.sigstore.dev/api/v1/log/entries/108e9186e8c5677a89619c7db02cfb94d2609666f60a8a48d41ee49b2e6553195f36fce510626ca7
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v2.0.0" at commit 292bc11372c8b0dc8dc23476bde9bab19c8d663b
Verifying artifact ghalint_1.0.0_darwin_arm64.tar.gz: PASSED

PASSED: SLSA verification passed
```

### 3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
gh release download -R suzuki-shunsuke/ghalint v1.0.0
cosign verify-blob \
  --signature ghalint_1.0.0_checksums.txt.sig \
  --certificate ghalint_1.0.0_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghalint_1.0.0_checksums.txt
```

Output:

```
Verified OK
```

After verifying the checksum, verify the artifact.

```sh
cat ghalint_1.0.0_checksums.txt | sha256sum -c --ignore-missing
```

</details>

5. go install

```sh
go install github.com/suzuki-shunsuke/ghalint/cmd/ghalint@latest
```

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
