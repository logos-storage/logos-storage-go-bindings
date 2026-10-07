package storage

import (
	"encoding/json"
	"os"
	"testing"
)

func mixConfig(t *testing.T) Config {
	t.Helper()

	raw, err := os.ReadFile("testdata/mix_config.json")
	if err != nil {
		t.Fatalf("Failed to read mix config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(raw, &config); err != nil {
		t.Fatalf("Failed to parse mix config: %v", err)
	}

	return config
}

func TestTogglePrivateQueriesDisableWithoutMix(t *testing.T) {
	node := newStorageNode(t)

	previous, err := node.TogglePrivateQueries(false)
	if err != nil {
		t.Fatalf("Failed to toggle private queries: %v", err)
	}
	if previous {
		t.Fatal("expected private queries to be disabled by default")
	}
}

func TestTogglePrivateQueriesEnableWithoutMixFails(t *testing.T) {
	node := newStorageNode(t)

	if _, err := node.TogglePrivateQueries(true); err == nil {
		t.Fatal("expected an error when enabling private queries without Mix configured")
	}
}

func TestTogglePrivateQueriesWithMix(t *testing.T) {
	node := newStorageNode(t, mixConfig(t))

	// Mix auto-enables private queries at startup, so disabling returns true.
	previous, err := node.TogglePrivateQueries(false)
	if err != nil {
		t.Fatalf("Failed to disable private queries: %v", err)
	}
	if !previous {
		t.Fatal("expected private queries to be enabled at startup when Mix is configured")
	}

	// Re-enabling is allowed because Mix is configured.
	previous, err = node.TogglePrivateQueries(true)
	if err != nil {
		t.Fatalf("Failed to enable private queries: %v", err)
	}
	if previous {
		t.Fatal("expected private queries to be disabled after the previous toggle")
	}
}
