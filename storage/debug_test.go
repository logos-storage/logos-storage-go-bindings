package storage

import (
	"os"
	"strings"
	"testing"
	"time"
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
	if len(info.AnnounceAddresses) == 0 {
		t.Error("Debug info AnnounceAddresses is empty")
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

func TestStoragePeerDebug(t *testing.T) {
	var bootstrap, node1, node2 *StorageNode
	var err error

	bootstrap = newStorageNode(t, Config{
		DiscoveryPort: 8092,
	})

	spr, err := bootstrap.Spr()
	if err != nil {
		t.Fatalf("Failed to get bootstrap spr: %v", err)
	}

	bootstrapNodes := []string{spr}

	node1 = newStorageNode(t, Config{
		DiscoveryPort:  8090,
		BootstrapNodes: bootstrapNodes,
	})

	node2 = newStorageNode(t, Config{
		DiscoveryPort:  8091,
		BootstrapNodes: bootstrapNodes,
	})

	peerId, err := node2.PeerId()
	if err != nil {
		t.Fatal(err)
	}

	var record PeerRecord
	for range 10 {
		record, err = node1.StoragePeerDebug(peerId)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		t.Fatalf("Logos StoragePeerDebug call failed: %v", err)
	}

	if record.PeerId == "" {
		t.Fatalf("Logos StoragePeerDebug call failed: %v", err)
	}
	if record.PeerId == "" {
		t.Error("Logos StoragePeerDebug info PeerId is empty")
	}
	if record.SeqNo == 0 {
		t.Error("Logos StoragePeerDebug info SeqNo is empty")
	}
	if len(record.Addresses) == 0 {
		t.Error("Logos StoragePeerDebug info Addresses is empty")
	}
	if record.PeerId != peerId {
		t.Errorf("Logos StoragePeerDebug info PeerId (%s) does not match requested PeerId (%s)", record.PeerId, peerId)
	}
}
