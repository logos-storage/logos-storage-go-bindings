# Release process

This document describes the release process for the Go bindings.

## Description

1. Ensure the main branch is up-to-date and all tests are passing.

2. Update the CHANGELOG.md with the description of the changes

3. Create a new tag, example:

```sh
git tag v0.0.15
git push --tags
```

4. The CI job will build the artifacts and create a draft release with the artifacts uploaded.

5. Copy the description added in the `CHANGELOG.md` file to the release description.

6. Publish it.

Once published, the artifacts can be downloaded using  the `version`, example: 

`https://github.com/logos-storage/logos-storage-nim/releases/download/v0.3.2/libstorage-linux-amd64-v0.3.2.zip`

It is not recommended to use the `latest` URL because you may face cache issues.

## Integration

Once released, you can integrate it into your Go project using:

```bash
go get github.com/logos-storage/logos-storage-go-bindings@v0.0.26
```

Then you can use the following `Makefile` command to fetch the artifact:

```bash
LIBS_DIR := $(abspath ./libs)
STORAGE_OS := linux
STORAGE_ARCH := amd64
STORAGE_VERSION := v0.3.2
STORAGE_DOWNLOAD_URL := "https://github.com/logos-storage/logos-storage-nim/releases/download/$(STORAGE_VERSION)/storage-${STORAGE_OS}-${STORAGE_ARCH}-${STORAGE_VERSION}.zip"

fetch-libstorage:
    mkdir -p $(LIBS_DIR); \
    curl -fSL --create-dirs -o $(LIBS_DIR)/storage-${STORAGE_OS}-${STORAGE_ARCH}-${STORAGE_VESION}.zip ${STORAGE_DOWNLOAD_URL}; \
    unzip -o -qq $(LIBS_DIR)/libstorage-${STORAGE_OS}-${STORAGE_ARCH}-${STORAGE_VESION}.zip -d $(LIBS_DIR); \
    rm -f $(LIBS_DIR)/libstorage*.zip;
```

`STORAGE_VERSION` is the release version from Logos Storage Nim.

### Nix

If you use Nix in a sandboxed environment, you cannot use curl to download the artifacts, so you have to prefetch them using the artifacts `SHA-256` hash. To generate the hash, you can use the following command: 

```bash
nix store prefetch-file --json --unpack https://github.com/logos-storage/logos-storage-nim/releases/download/v0.3.2/libstorage-darwin-arm64-v0.3.2.zip | jq -r .hash

# [10.4 MiB DL] sha256-3CHIWoSjo0plsYqzXQWm1EtY1STcljV4yfXTPon90uE=
```

Then include this hash in your Nix configuration. For example:

```nix
let
  optionalString = pkgs.lib.optionalString;
  storageVersion = "v0.0.26";
  arch =
    if stdenv.hostPlatform.isx86_64 then "amd64"
    else if stdenv.hostPlatform.isAarch64 then "arm64"
    else stdenv.hostPlatform.arch;
  os = if stdenv.isDarwin then "macos" else "Linux";
  hash =
    if stdenv.hostPlatform.isDarwin
    # nix store prefetch-file --json --unpack https://github.com/logos-storage/logos-storage-nim/releases/download/v0.3.2/libstorage-darwin-arm64-v0.3.2.zip | jq -r .hash
    then "sha256-3CHIWoSjo0plsYqzXQWm1EtY1STcljV4yfXTPon90uE="
    # nix store prefetch-file --json --unpack https://github.com/logos-storage/logos-storage-nim/releases/download/v0.3.2/libstorage-linux-amd64-v0.3.2.zip | jq -r .hash
    else "sha256-YxW2vFZlcLrOx1PYgWW4MIstH/oFBRF0ooS0sl3v6ig=";

  # Pre-fetch libstorage to avoid network during build
  storageLib = pkgs.fetchzip {
    url = "https://github.com/logos-storage/logos-storage-nim/releases/download/${storageVersion}/storage-${os}-${arch}-${storageVersion}.zip";
    hash = hash;
    stripRoot = false;
  };

  preBuild = ''
    export LIBS_DIR="${storageLib}"
    # Build something cool with Logos Storage
  '';
```

