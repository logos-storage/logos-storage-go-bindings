package storage

import "testing"

func TestDownloadManifestPrivateWithoutMixFails(t *testing.T) {
	node := newStorageNode(t)
	cid, _ := uploadHelper(t, node)

	if _, err := node.DownloadManifest(cid, DownloadManifestOptions{Private: true}); err == nil {
		t.Fatal("expected an error when downloading privately without Mix configured")
	}
}
