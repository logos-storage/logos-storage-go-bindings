package codex

import (
	"testing"
)

func TestCodexVersion(t *testing.T) {
	config := defaultConfigHelper(t)
	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}

	version, err := node.Version()
	if err != nil {
		t.Fatalf("Failed to get Codex version: %v", err)
	}
	if version == "" {
		t.Fatal("Codex version is empty")
	}

	t.Logf("Codex version: %s", version)
}

func TestCodexRevision(t *testing.T) {
	config := defaultConfigHelper(t)
	node, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}

	revision, err := node.Revision()
	if err != nil {
		t.Fatalf("Failed to get Codex revision: %v", err)
	}
	if revision == "" {
		t.Fatal("Codex revision is empty")
	}

	t.Logf("Codex revision: %s", revision)
}

func TestCodexRepo(t *testing.T) {
	node := newCodexNode(t)

	repo, err := node.Repo()
	if err != nil {
		t.Fatalf("Failed to get Codex repo: %v", err)
	}
	if repo == "" {
		t.Fatal("Codex repo is empty")
	}

	t.Logf("Codex repo: %s", repo)
}

func TestSpr(t *testing.T) {
	node := newCodexNode(t)

	spr, err := node.Spr()
	if err != nil {
		t.Fatalf("Failed to get Codex SPR: %v", err)
	}
	if spr == "" {
		t.Fatal("Codex SPR is empty")
	}

	t.Logf("Codex SPR: %s", spr)
}

func TestPeerId(t *testing.T) {
	node := newCodexNode(t)

	peerId, err := node.PeerId()
	if err != nil {
		t.Fatalf("Failed to get Codex PeerId: %v", err)
	}
	if peerId == "" {
		t.Fatal("Codex PeerId is empty")
	}

	t.Logf("Codex PeerId: %s", peerId)
}

func TestStorageQuota(t *testing.T) {
	node := newCodexNode(t, Config{
		StorageQuota: 1024 * 1024 * 1024, // 1GB
	})

	if node == nil {
		t.Fatal("expected codex node to be created")
	}
}

func TestCreateAndDestroyMultipleInstancesWithSameDatadir(t *testing.T) {
	datadir := t.TempDir()

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
			t.Fatalf("Failed to create Codex node: %v", err)
		}

		if err := node.Start(); err != nil {
			t.Fatalf("Failed to start Codex node: %v", err)
		}

		if err := node.Stop(); err != nil {
			t.Fatalf("Failed to stop Codex node: %v", err)
		}

		if err := node.Destroy(); err != nil {
			t.Fatalf("Failed to stop Codex node after restart: %v", err)
		}
	}
}

func TestNumThreads(t *testing.T) {
	node := newCodexNode(t, Config{
		NumThreads: 1,
	})

	if node == nil {
		t.Fatal("expected codex node to be created")
	}
}

func TestBlockTtl(t *testing.T) {
	node := newCodexNode(t, Config{
		BlockTtl: "10H",
	})

	if node == nil {
		t.Fatal("expected codex node to be created")
	}
}

func TestBlockMaintenanceInterval(t *testing.T) {
	node := newCodexNode(t, Config{
		BlockMaintenanceInterval: "10H",
	})

	if node == nil {
		t.Fatal("expected codex node to be created")
	}
}
