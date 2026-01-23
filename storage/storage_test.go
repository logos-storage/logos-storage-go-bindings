package storage

import "testing"

func TestManifests(t *testing.T) {
	storage := newStorageNode(t)

	manifests, err := storage.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) != 0 {
		t.Fatal("expected manifests to be empty")
	}

	cid, _ := uploadHelper(t, storage)

	manifests, err = storage.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) == 0 {
		t.Fatal("expected manifests to be non-empty")
	}

	for _, m := range manifests {
		if m.Cid != cid {
			t.Errorf("expected cid %q, got %q", cid, m.Cid)
		}
	}
}

func TestSpace(t *testing.T) {
	storage := newStorageNode(t)

	space, err := storage.Space()
	if err != nil {
		t.Fatal(err)
	}

	if space.TotalBlocks != 0 {
		t.Fatal("expected total blocks to be non-zero")
	}

	if space.QuotaMaxBytes == 0 {
		t.Fatal("expected quota max bytes to be non-zero")
	}

	if space.QuotaUsedBytes != 0 {
		t.Fatal("expected quota used bytes to be non-zero")
	}

	if space.QuotaReservedBytes != 0 {
		t.Fatal("expected quota reserved bytes to be non-zero")
	}

	uploadHelper(t, storage)

	space, err = storage.Space()
	if err != nil {
		t.Fatal(err)
	}

	if space.TotalBlocks == 0 {
		t.Fatal("expected total blocks to be non-zero after upload")
	}

	if space.QuotaUsedBytes == 0 {
		t.Fatal("expected quota used bytes to be non-zero after upload")
	}
}

func TestFetch(t *testing.T) {
	storage := newStorageNode(t)

	cid, _ := uploadHelper(t, storage)

	_, err := storage.Fetch(cid)
	if err != nil {
		t.Fatal("expected error when fetching non-existent manifest")
	}
}

func TestFetchCidDoesNotExist(t *testing.T) {
	storage := newStorageNode(t, Config{BlockRetries: 1})

	_, err := storage.Fetch("bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku")
	if err == nil {
		t.Fatal("expected error when fetching non-existent manifest")
	}
}

func TestDelete(t *testing.T) {
	storage := newStorageNode(t)

	cid, _ := uploadHelper(t, storage)

	manifests, err := storage.Manifests()
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 1 {
		t.Fatal("expected manifests to be empty after deletion")
	}

	err = storage.Delete(cid)
	if err != nil {
		t.Fatal(err)
	}

	manifests, err = storage.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) != 0 {
		t.Fatal("expected manifests to be empty after deletion")
	}
}

func TestExists(t *testing.T) {
	storage := newStorageNode(t)

	cid, _ := uploadHelper(t, storage)

	exists, err := storage.Exists(cid)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected cid to exist")
	}

	err = storage.Delete(cid)
	if err != nil {
		t.Fatal(err)
	}

	exists, err = storage.Exists(cid)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("expected cid to not exist after deletion")
	}
}
