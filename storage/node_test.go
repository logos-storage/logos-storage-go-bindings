package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestStorageVersion(t *testing.T) {
	config := defaultConfigHelper(t)
	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	version, err := node.Version()
	if err != nil {
		t.Fatalf("Failed to get Logos Storage version: %v", err)
	}
	if version == "" {
		t.Fatal("Logos Storage version is empty")
	}

	t.Logf("Logos Storage version: %s", version)
}

func TestStorageRevision(t *testing.T) {
	config := defaultConfigHelper(t)
	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	revision, err := node.Revision()
	if err != nil {
		t.Fatalf("Failed to get Logos Storage revision: %v", err)
	}
	if revision == "" {
		t.Fatal("Logos Storage revision is empty")
	}

	t.Logf("Logos Storage revision: %s", revision)
}

func TestStorageRepo(t *testing.T) {
	node := newStorageNode(t)

	repo, err := node.Repo()
	if err != nil {
		t.Fatalf("Failed to get Logos Storage repo: %v", err)
	}
	if repo == "" {
		t.Fatal("Logos Storage repo is empty")
	}

	t.Logf("Logos Storage repo: %s", repo)
}

func TestSpr(t *testing.T) {
	node := newStorageNode(t)

	spr, err := node.Spr()
	if err != nil {
		t.Fatalf("Failed to get Logos Storage SPR: %v", err)
	}
	if spr == "" {
		t.Fatal("Logos Storage SPR is empty")
	}

	t.Logf("Logos Storage SPR: %s", spr)
}

func TestPeerId(t *testing.T) {
	node := newStorageNode(t)

	peerId, err := node.PeerId()
	if err != nil {
		t.Fatalf("Failed to get Logos Storage PeerId: %v", err)
	}
	if peerId == "" {
		t.Fatal("Logos Storage PeerId is empty")
	}

	t.Logf("Logos Storage PeerId: %s", peerId)
}

func TestStorageQuota(t *testing.T) {
	node := newStorageNode(t, Config{
		StorageQuota: 1024 * 1024 * 1024, // 1GB
	})

	if node == nil {
		t.Fatal("expected Logos Storage node to be created")
	}
}

func TestCreateAndDestroyMultipleInstancesWithSameDatadir(t *testing.T) {
	t.Skip("Enable when the PR https://github.com/logos-storage/logos-storage-nim/pull/1364 is merged into master.")

	datadir := fmt.Sprintf("%s/special-test", t.TempDir())

	config := Config{
		DataDir:        datadir,
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		BlockRetries:   5,
		Nat:            "none",
	}

	for range 2 {
		node, err := New(config)
		if err != nil {
			t.Fatalf("Failed to create Logos Storage node: %v", err)
		}

		if err := node.Start(); err != nil {
			t.Fatalf("Failed to start Logos Storage node: %v", err)
		}

		if err := node.Stop(); err != nil {
			t.Fatalf("Failed to stop Logos Storage node: %v", err)
		}

		if err := node.Destroy(); err != nil {
			t.Fatalf("Failed to stop Logos Storage node after restart: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestNumThreads(t *testing.T) {
	node := newStorageNode(t, Config{
		NumThreads: 1,
	})

	if node == nil {
		t.Fatal("expected Logos Storage node to be created")
	}
}

func TestBlockTtl(t *testing.T) {
	node := newStorageNode(t, Config{
		BlockTtl: "10H",
	})

	if node == nil {
		t.Fatal("expected Logos Storage node to be created")
	}
}

func TestBlockMaintenanceInterval(t *testing.T) {
	node := newStorageNode(t, Config{
		BlockMaintenanceInterval: "10H",
	})

	if node == nil {
		t.Fatal("expected Logos Storage node to be created")
	}
}
