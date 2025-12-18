# Install

ghalint is written in Go. So you only have to install a binary in your `PATH`.

There are some ways to install ghalint.

1. [Homebrew](#homebrew)
1. [Mise](#mise)
1. [Scoop](#scoop)
1. [aqua](#aqua)
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

## Mise

You can install ghalint using [Mise](https://github.com/jdx/mise), you can install and made available the latest version globally by using a command like:

```sh
mise use -g ghalint@latest
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
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.2.0
asset=ghalint_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/ghalint "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/ghalint \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.2.0
asset=ghalint_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/ghalint "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ghalint \
  --source-tag "$version"
```

### 3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
version=v1.2.0
checksum_file="ghalint_${version#v}_checksums.txt"
asset=ghalint_darwin_arm64.tar.gz
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
cat "$checksum_file" | sha256sum -c --ignore-missing
```
