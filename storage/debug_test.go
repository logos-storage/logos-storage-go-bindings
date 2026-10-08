package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	storage := newStorageNode(t)

	info, err := storage.Debug()
	if err != nil {
		t.Fatalf("Debug call failed: %v", err)
	}
	if info.ID == "" {
		t.Error("Debug info ID is empty")
	}
	if info.Spr == "" {
		t.Error("Debug info Spr is empty")
	}
}

func runAndReadLogHelper(t *testing.T, node *StorageNode, logFile string) string {
	t.Helper()

	if err := node.Start(); err != nil {
		t.Fatalf("Failed to start Logos Storage node: %v", err)
	}

	if err := node.Stop(); err != nil {
		t.Fatalf("Failed to stop Logos Storage node: %v", err)
	}

	if err := node.Destroy(); err != nil {
		t.Fatalf("Failed to destroy Logos Storage node: %v", err)
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	return string(content)
}

func TestLogLevel(t *testing.T) {
	config := defaultConfigHelper(t)
	config.LogLevel = "ERROR"
	config.LogFile = filepath.Join(t.TempDir(), "storage.log")

	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	content := runAndReadLogHelper(t, node, config.LogFile)

	if strings.Contains(content, "INF") {
		t.Errorf("Log file contains INFO statement with ERROR log level: %s", content)
	}
}

func TestUpdateLogLevel(t *testing.T) {
	config := defaultConfigHelper(t)
	config.LogLevel = "ERROR"
	config.LogFile = filepath.Join(t.TempDir(), "storage.log")

	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	if err := node.UpdateLogLevel("INFO"); err != nil {
		t.Fatalf("UpdateLogLevel call failed: %v", err)
	}

	content := runAndReadLogHelper(t, node, config.LogFile)

	if !strings.Contains(content, "INF") {
		t.Errorf("Log file does not contain INFO statement after log level update: %s", content)
	}
}
