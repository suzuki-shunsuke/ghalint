# Install

ghalint is written in Go. So you only have to install a binary in your `PATH`.

There are some ways to install ghalint.

1. [Homebrew](#homebrew)
1. [Scoop](#scoop)
1. [aqua](#aqua)
1. [mise](#mise)
1. [GitHub Releases](#github-releases)
1. [Build an executable binary from source code yourself using Go](#build-an-executable-binary-from-source-code-yourself-using-go)

## Homebrew

You can install ghalint using [Homebrew](https://brew.sh/).

```sh
brew install ghalint
```

Or

```sh
brew install suzuki-shunsuke/ghalint/ghalint
```

## Scoop

You can install ghalint using [Scoop](https://scoop.sh/).

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install ghalint
```

## aqua

You can install ghalint using [aqua](https://aquaproj.github.io/).

```sh
aqua g -i suzuki-shunsuke/ghalint
```

## mise

You can install ghalint using [mise](https://github.com/jdx/mise).

```sh
mise use -g ghalint@latest
```

## Build an executable binary from source code yourself using Go

```sh
go install github.com/suzuki-shunsuke/ghalint/cmd/ghalint@latest
```

## GitHub Releases

You can download an asset from [GitHub Releases](https://github.com/suzuki-shunsuke/ghalint/releases).
Please unarchive it and install a pre built binary into `$PATH`. 

### Verify downloaded assets from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [Cosign](https://github.com/sigstore/cosign)
1. [Packslip](https://packslip.dev/) `ghalint >= v1.5.7`
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier) `ghalint <= v1.5.6`

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.5.7
asset="ghalint_${version#v}_darwin_arm64.tar.gz"
gh release download -R suzuki-shunsuke/ghalint "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/ghalint \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

From ghalint v1.5.7, the bundle file is `ghalint_<version>_checksums.txt.sigstore.json`.
Cosign v2.4.2 or later is required.

```sh
version=v1.5.7
checksum_file="ghalint_${version#v}_checksums.txt"
asset="ghalint_${version#v}_darwin_arm64.tar.gz"
gh release download "$version" \
  -R suzuki-shunsuke/ghalint \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.sigstore.json"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.sigstore.json" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

For ghalint v1.5.6 or older, `*.pem` and `*.sig` are available instead.

```sh
version=v1.5.6
checksum_file="ghalint_${version#v}_checksums.txt"
asset="ghalint_${version#v}_darwin_arm64.tar.gz"
gh release download "$version" \
  -R suzuki-shunsuke/ghalint \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

### 3. Packslip

`ghalint >= v1.5.7`

You can install packslip by aqua.

```sh
aqua g -i jdx/packslip
```

```sh
version=v1.5.7
asset="ghalint_${version#v}_darwin_arm64.tar.gz"
gh release download -R suzuki-shunsuke/ghalint "$version" -p "$asset" -p packslip.sigstore.json
packslip verify packslip.sigstore.json \
    --pin ps1_uk43kothwftpzcoebqslnwvl6m \
    --identity-prefix https://github.com/suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@ \
    --artifact "$asset"
```

- `--pin` requires the release to be signed from the repository `suzuki-shunsuke/ghalint`, identified by its repository ID rather than its name
- `--identity-prefix` requires the release to be signed by the release workflow of [suzuki-shunsuke/go-release-workflow](https://github.com/suzuki-shunsuke/go-release-workflow)

On success, packslip prints a line starting with `ok:`.
Please make sure it shows the project and version you meant, e.g. `github.com/suzuki-shunsuke/ghalint 1.5.7`.

### 4. slsa-verifier

> [!WARNING]
> As of ghalint v1.5.7, `multiple.intoto.jsonl` hasn't been published anymore.

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.5.6
asset="ghalint_${version#v}_darwin_arm64.tar.gz"
gh release download -R suzuki-shunsuke/ghalint "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ghalint \
  --source-tag "$version"
```
