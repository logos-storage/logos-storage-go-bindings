package storage

import (
	"bytes"
	"context"
	"testing"
)

func defaultConfigHelper(t *testing.T) Config {
	t.Helper()

	return Config{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		BlockRetries:   3000,
		Nat:            "none",
	}
}

func newStorageNode(t *testing.T, opts ...Config) *StorageNode {
	config := defaultConfigHelper(t)

	if len(opts) > 0 {
		defaults := config
		config = opts[0]

		if config.DataDir == "" {
			config.DataDir = defaults.DataDir
		}

		if config.LogFormat == "" {
			config.LogFormat = defaults.LogFormat
		}

		if config.BlockRetries == 0 {
			config.BlockRetries = defaults.BlockRetries
		}

		if config.Nat == "" {
			config.Nat = defaults.Nat
		}
	}

	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	err = node.Start()
	if err != nil {
		t.Fatalf("Failed to start Logos Storage node: %v", err)
	}

	t.Cleanup(func() {
		if err := node.Stop(); err != nil {
			t.Logf("cleanup storage: %v", err)
		}

		if err := node.Destroy(); err != nil {
			t.Logf("cleanup storage: %v", err)
		}
	})

	return node
}

func uploadHelper(t *testing.T, storage *StorageNode) (string, int) {
	t.Helper()

	buf := bytes.NewBuffer([]byte("Hello World!"))
	len := buf.Len()
	cid, err := storage.UploadReader(context.Background(), UploadOptions{Filepath: "hello.txt"}, buf)
	if err != nil {
		t.Fatalf("Error happened during upload: %v\n", err)
	}

	return cid, len
}

func uploadBigFileHelper(t *testing.T, storage *StorageNode) (string, int) {
	t.Helper()

	len := 1024 * 1024 * 50
	buf := bytes.NewBuffer(make([]byte, len))

	cid, err := storage.UploadReader(context.Background(), UploadOptions{Filepath: "hello.txt"}, buf)
	if err != nil {
		t.Fatalf("Error happened during upload: %v\n", err)
	}

	return cid, len
}
