package storage

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestDownloadStream(t *testing.T) {
	storage := newStorageNode(t)
	cid, len := uploadHelper(t, storage)

	f, err := os.Create("testdata/hello.downloaded.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	totalBytes := 0
	finalPercent := 0.0
	opt := DownloadStreamOptions{
		Writer:      f,
		DatasetSize: len,
		Filepath:    "testdata/hello.downloaded.writer.txt",
		OnProgress: func(read, total int, percent float64, err error) {
			if err != nil {
				t.Fatalf("Error happening during download: %v\n", err)
			}

			totalBytes = total
			finalPercent = percent
		},
	}

	if err := storage.DownloadStream(context.Background(), cid, opt); err != nil {
		t.Fatal("Error happened:", err.Error())
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback total bytes %d but expected %d", totalBytes, len)
	}

	data, err := os.ReadFile("testdata/hello.downloaded.writer.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hello World!" {
		t.Fatalf("Downloaded content does not match, expected Hello World! got %s", data)
	}
}

func TestDownloadStreamWithAutosize(t *testing.T) {
	storage := newStorageNode(t)
	cid, len := uploadHelper(t, storage)

	totalBytes := 0
	finalPercent := 0.0
	opt := DownloadStreamOptions{
		DatasetSizeAuto: true,
		OnProgress: func(read, total int, percent float64, err error) {
			if err != nil {
				t.Fatalf("Error happening during download: %v\n", err)
			}

			totalBytes = total
			finalPercent = percent
		},
	}

	if err := storage.DownloadStream(context.Background(), cid, opt); err != nil {
		t.Fatal("Error happened:", err.Error())
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback total bytes %d but expected %d", totalBytes, len)
	}
}

func TestDownloadStreamWithNotExisting(t *testing.T) {
	storage := newStorageNode(t, Config{BlockRetries: 1})

	opt := DownloadStreamOptions{}
	if err := storage.DownloadStream(context.Background(), "bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku", opt); err == nil {
		t.Fatal("Error expected when downloading non-existing cid")
	}
}

func TestDownloadStreamCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	storage := newStorageNode(t)
	cid, _ := uploadBigFileHelper(t, storage)

	channelError := make(chan error, 1)
	go func() {
		err := storage.DownloadStream(ctx, cid, DownloadStreamOptions{Local: true})
		channelError <- err
	}()

	cancel()
	err := <-channelError

	if err == nil {
		t.Fatal("DownloadStream should have been canceled")
	}

	if err.Error() != context.Canceled.Error() {
		t.Fatalf("DownloadStream returned unexpected error: %v", err)
	}
}

func TestDownloadManual(t *testing.T) {
	storage := newStorageNode(t)
	cid, _ := uploadHelper(t, storage)

	if err := storage.DownloadInit(cid, DownloadInitOptions{}); err != nil {
		t.Fatal("Error when initializing download:", err)
	}

	var b strings.Builder
	if chunk, err := storage.DownloadChunk(cid); err != nil {
		t.Fatal("Error when downloading chunk:", err)
	} else {
		b.Write(chunk)
	}

	data := b.String()
	if data != "Hello World!" {
		t.Fatalf("Expected data was \"Hello World!\" got %s", data)
	}

	if err := storage.DownloadCancel(cid); err != nil {
		t.Fatalf("Error when cancelling the download %s", err)
	}
}

func TestDownloadManifest(t *testing.T) {
	storage := newStorageNode(t)
	cid, _ := uploadHelper(t, storage)

	manifest, err := storage.DownloadManifest(cid)
	if err != nil {
		t.Fatal("Error when downloading manifest:", err)
	}

	if manifest.Cid != cid {
		t.Errorf("expected cid %q, got %q", cid, manifest.Cid)
	}
}

func TestDownloadManifestWithNotExistingCid(t *testing.T) {
	storage := newStorageNode(t, Config{BlockRetries: 1})

	manifest, err := storage.DownloadManifest("bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku")
	if err == nil {
		t.Fatal("Error when downloading manifest:", err)
	}

	if manifest.Cid != "" {
		t.Errorf("expected empty cid, got %q", manifest.Cid)
	}
}

func TestDownloadInitWithNotExistingCid(t *testing.T) {
	storage := newStorageNode(t, Config{BlockRetries: 1})

	if err := storage.DownloadInit("bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku", DownloadInitOptions{}); err == nil {
		t.Fatal("expected error when initializing download for non-existent cid")
	}
}
