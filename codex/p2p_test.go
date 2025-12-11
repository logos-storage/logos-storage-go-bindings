package codex

import (
	"log"
	"testing"
)

func TestConnectWithAddress(t *testing.T) {
	var err error

	node1 := newCodexNode(t, Config{
		DiscoveryPort: 8090,
	})

	node2 := newCodexNode(t, Config{
		DiscoveryPort: 8091,
	})

	info2, err := node2.Debug()
	if err != nil {
		t.Fatal(err)
	}

	if err := node1.Connect(info2.ID, info2.Addrs); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
}

func TestCodexWithPeerId(t *testing.T) {
	var bootstrap, node1, node2 *CodexNode
	var err error

	bootstrap = newCodexNode(t, Config{
		DiscoveryPort: 8092,
	})

	spr, err := bootstrap.Spr()
	if err != nil {
		t.Fatalf("Failed to get bootstrap spr: %v", err)
	}

	bootstrapNodes := []string{spr}

	node1 = newCodexNode(t, Config{
		DiscoveryPort:  8090,
		BootstrapNodes: bootstrapNodes,
	})

	node2 = newCodexNode(t, Config{
		DiscoveryPort:  8091,
		BootstrapNodes: bootstrapNodes,
	})

	peerId, err := node2.PeerId()
	if err != nil {
		t.Fatal(err)
	}

	if err := node1.Connect(peerId, []string{}); err != nil {
		log.Println(err)
	}
}
