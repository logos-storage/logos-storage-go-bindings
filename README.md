# Logos Storage Go Bindings

This repository provides Go bindings for the Logos Stroage library, enabling seamless integration with Go projects.

## Usage

Include in your Go project:

```sh
go get github.com/logos-storage/logos-storage-go-bindings
```

Then the easiest way is to download our prebuilt artifacts and configure your project.
You can use this `Makefile` (or integrates the commands in your build process):

```makefile
# Path configuration
LIBS_DIR := $(abspath ./libs)
CGO_CFLAGS  := -I$(LIBS_DIR)
CGO_LDFLAGS := -L$(LIBS_DIR) -lstorage -Wl,-rpath,$(LIBS_DIR)

# Fetch configuration
OS ?= "linux"
ARCH ?= "amd64"
VERSION ?= "v0.0.21"
DOWNLOAD_URL := "https://github.com/logos-storage/logos-storage-go-bindings/releases/download/$(VERSION)/storage-${OS}-${ARCH}.zip"

# Edit your binary name here
ifeq ($(OS),Windows_NT)
  BIN_NAME := example.exe
else
  BIN_NAME := example
endif

fetch:
	@echo "Fetching libstorage from GitHub Actions from: ${DOWNLOAD_URL}"
	@curl -fSL --create-dirs -o $(LIBS_DIR)/storage-${OS}-${ARCH}.zip ${DOWNLOAD_URL}
	@unzip -o -qq $(LIBS_DIR)/storage-${OS}-${ARCH}.zip -d $(LIBS_DIR)
	@rm -f $(LIBS_DIR)/*.zip

build:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o $(BIN_NAME) main.go

clean:
	rm -f $(BIN_NAME)
	rm -Rf $(LIBS_DIR)/*
```

First you need to `fetch` the artefacts for your `OS` and `ARCH`:

```sh
OS=macos ARCH=arm64 make fetch
```

Then you can build your project using:

```sh
make build
```

That's it!

For an example on how to use this package, please take a look at our [example-go-bindings](https://github.com/logos-storage/example-logos-storage-go-bindings) repo.

If you want to build the library yourself, you need to clone this repo and follow the instructions
of the next step.

## Development

To build the required dependencies for this module, the `make` command needs to be executed.
If you are integrating this module into another project via `go get`, ensure that you navigate
to the `logos-storage-go-bindings` module directory and run the `make` commands.

### Steps to install

Follow these steps to install and set up the module:

1. Make sure your system has the [prerequisites](https://github.com/logos-storage/logos-storage-nim) to run a local Logos Storage node.

2. Fetch the dependencies:
   ```sh
   make update
   ```

3. Build the library:
   ```sh
   make libstorage
   ```

You can pass flags to the Logos Storage building step by using `STORAGE_LIB_PARAMS`. For example,
if you want to enable debug API for peers, you can build the library using:

```sh
STORAGE_LIB_PARAMS="-d:storage_enable_api_debug_peers=true" make libstorage
```

or you can use a convenience `libstorage-with-debug-api` make target:

```sh
make libstorage-with-debug-api
```

To run the test, you have to make sure you have `gotestsum` installed on your system, e.g.:

```sh
go install gotest.tools/gotestsum@latest
```

Then you can run the tests as follows.

To run all the tests:

```sh
make test
```

To run selected test only:

```sh
make test "TestDownloadManifest$"
```

> We use `$` to make sure we run only the `TestDownloadManifest` test.
> Without `$` we would run all the tests starting with `TestDownloadManifest` and
> so also `TestDownloadManifestWithNotExistingCid`
>

If you need to pass more arguments to the underlying `go test` (`gotestsum` passes
everything after `--` to `go test`), you can use: `test-with-params` make target, e.g.:

```sh
make test-with-params ARGS='-run "TestDownloadManifest$$" -count=2'
```

> Here, we use double escape `$$` instead of just `$`, otherwise make
> will interpret `$` as a make variable inside `ARGS`. 

Now the module is ready for use in your project.

The release process is defined [here](./RELEASE.md).

## API

### Init

First you need to create a Logos Storage node:

```go
dataDir := "..."
node, err := StorageNew(StorageConfig{
   DataDir:        dataDir,
   BlockRetries:   10,
})
/// ....
err := node.Destroy()
```

The `StorageConfig` object provides several options to configure your node. You should at least
adjust the `DataDir` folder and the `BlockRetries` setting to avoid long retrieval times when
the data is unavailable.

When you are done with your node, you **have to** call `Destroy` method to free resources.

### Start / Stop

use `Start` method to start your node. You **have to** call `Stop` before `Destroy` when you are done
with your node.

```go
err := node.Start()
err := node.Stop()
err := node.Destroy()
```

### Info

You can get the version and revision without starting the node:

```go
version, err := node.Version()
revision, err := node.Revision()
```

Other information is available after the node is started:

```go
version, err := node.Version()
spr, err := node.Spr()
peerId, err := node.PeerId()
```

### Upload

There are 3 strategies for uploading: `reader`, `file` or `chunks`. Each one requires its own upload session.

#### reader

The `reader` strategy is the easiest option when you already have a Go `Reader`.
It handles creating the upload session and cancels it if an error occurs.

The `filepath` should contain the data’s name with its extension, because Logos Storage uses that to
infer the MIME type.

An `onProgress` callback is available to receive progress updates and notify the user.
The total size of the reader is determined via `stat` when it wraps a file, or from the buffer length otherwise.
From there, the callback can compute and report the percentage complete.

The `UploadReader` returns the cid of the content uploaded.

```go
buf := bytes.NewBuffer([]byte("Hello World!"))
onProgress := func(read, total int, percent float64, err error) {
   // Do something with the data
}
ctx := context.Background()
cid, err := storage.UploadReader(ctx, UploadOptions{filepath: "hello.txt", onProgress: onProgress}, buf)
```

#### file

The `file` strategy allows you to upload a file on Logos Storage using the path.
It handles creating the upload session and cancels it if an error occurs.

The `onProgress` callback is the same as for `reader` strategy.

The `UploadFile` returns the cid of the content uploaded.

```go
onProgress := func(read, total int, percent float64, err error) {
   // Do something with the data
}
ctx := context.Background()
cid, err := storage.UploadFile(ctx, UploadOptions{filepath: "./testdata/hello.txt", onProgress: onProgress})
```

#### chunks

The `chunks` strategy allows you to manage the upload by yourself. It requires more code
but provides more flexibility. You have to create the upload session, send the chunks
and then finalize to get the cid.

```go
sessionId, err := storage.UploadInit(&UploadOptions{filepath: "hello.txt"})

err = storage.UploadChunk(sessionId, []byte("Hello "))

err = storage.UploadChunk(sessionId, []byte("World!"))

cid, err := storage.UploadFinalize(sessionId)
```

Using this strategy, you can handle resumable uploads and cancel the upload
whenever you want!

### Download

When you receive a cid, you can download the `Manifest` to get information about the data:

```go
manifest, err := storage.DownloadManifest(cid)
```

It is not mandatory for downloading the data but it is really useful.

There are 2 strategies for downloading: `stream` and `chunks`.

#### stream

The `stream` strategy is the easiest to use.

It provides an `onProgress` callback to receive progress updates and notify the user.
The percentage is calculated from the `datasetSize` (taken from the manifest).
If you don’t provide it, you can enable `datasetSizeAuto` so `DownloadStream` fetches the
manifest first and uses its `datasetSize`.

You can pass a `writer` and/or a `filepath` as destinations. They are not mutually exclusive,
letting you write the content to two places for the same download.

```go
opt := DownloadStreamOptions{
   writer:      f,
   datasetSize: len,
   filepath:    "testdata/hello.downloaded.writer.txt",
   onProgress: func(read, total int, percent float64, err error) {
      // Handle progress
   },
}
ctx := context.Background()
err := storage.DownloadStream(ctx, cid, opt)
```

#### chunks

The `chunks` strategy allows to manage the download by yourself. It requires more code
but provide more flexibility.

This strategy **assumes you already know the total size to download** (from the manifest).
After you believe all chunks have been retrieved, you **must** call `DownloadCancel`
to terminate the download session.

```go
cid := "..."
err := storage.DownloadInit(cid, DownloadInitOptions{})
chunk, err := storage.DownloadChunk(cid)
err := storage.DownloadCancel(cid)
```

Using this strategy, you can handle resumable downloads and cancel the download
whenever you want !

### Storage

Several methods are available to manage the data on your node:

```go
manifests, err := node.Manifests()
space, err := node.Space()

cid := "..."
err := node.Delete(cid)
err := node.Fetch(cid)
```

The `Fetch` method downloads remote data into your local node.

### P2P

You can connect to a node using the `peerId` or the `listenAddresses`:

```go
peerId := "..."
addrs := ["..."]
err := node.Connect(peerId, addrs)
```

### Debug

Several methods are available to debug your node:

```go
// Get node info
info, err := node.Debug()

// Update the chronicles level log on runtime
err := node.UpdateLogLevel("DEBUG")

peerId := "..."
record, err := node.StoragePeerDebug(peerId)
```

`StoragePeerDebug` is only available if you built with `-d:STORAGE_enable_api_debug_peers=true` flag.

### Context and cancellation

Go contexts are exposed only on the long-running operations as `UploadReader`, `UploadFile`, and `DownloadFile`. If the
context is cancelled, those methods cancel the active upload or download. Short lived API calls don’t take a context
because they usually finish before a cancellation signal could matter.