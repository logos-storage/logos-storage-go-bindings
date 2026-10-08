package storage

import (
	"os"
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
	if len(info.ProviderAddresses) == 0 {
		t.Error("Debug info ProviderAddresses is empty")
	}
}

func TestUpdateLogLevel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "storage-log-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	node := newStorageNode(t, Config{
		LogLevel:  "INFO",
		LogFile:   tmpFile.Name(),
		LogFormat: LogFormatNoColors,
	})

	content, err := os.ReadFile(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "INF") {
		t.Errorf("Log file does not contain INFO statement %s", string(content))
	}

	err = node.UpdateLogLevel("ERROR")
	if err != nil {
		t.Fatalf("UpdateLogLevel call failed: %v", err)
	}

	if err := node.Stop(); err != nil {
		t.Fatalf("Failed to stop Logos Storage node: %v", err)
	}

	// Clear the file
	if err := os.WriteFile(tmpFile.Name(), []byte{}, 0644); err != nil {
		t.Fatalf("Failed to clear log file: %v", err)
	}

	err = node.Start()
	if err != nil {
		t.Fatalf("Failed to start Logos Storage node: %v", err)
	}

	content, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if strings.Contains(string(content), "INF") {
		t.Errorf("Log file contains INFO statement after log level update: %s", string(content))
	}
}
