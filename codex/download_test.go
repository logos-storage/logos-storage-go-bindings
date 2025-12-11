package codex

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDownloadStream(t *testing.T) {
	codex := newCodexNode(t)
	cid, len := uploadHelper(t, codex)

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

	if err := codex.DownloadStream(context.Background(), cid, opt); err != nil {
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
	codex := newCodexNode(t)
	cid, len := uploadHelper(t, codex)

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

	if err := codex.DownloadStream(context.Background(), cid, opt); err != nil {
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
	codex := newCodexNode(t, Config{BlockRetries: 1})

	opt := DownloadStreamOptions{}
	if err := codex.DownloadStream(context.Background(), "bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku", opt); err == nil {
		t.Fatal("Error expected when downloading non-existing cid")
	}
}

func TestDownloadStreamCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	codex := newCodexNode(t)
	cid, _ := uploadBigFileHelper(t, codex)

	channelError := make(chan error, 1)
	go func() {
		err := codex.DownloadStream(ctx, cid, DownloadStreamOptions{Local: true})
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
	ctx := context.Background()
	codex := newCodexNode(t)
	cid, _ := uploadHelper(t, codex)

	if err := codex.DownloadInit(ctx, cid, DownloadInitOptions{}); err != nil {
		t.Fatal("Error when initializing download:", err)
	}

	var b strings.Builder
	if chunk, err := codex.DownloadChunk(cid); err != nil {
		t.Fatal("Error when downloading chunk:", err)
	} else {
		b.Write(chunk)
	}

	data := b.String()
	if data != "Hello World!" {
		t.Fatalf("Expected data was \"Hello World!\" got %s", data)
	}

	if err := codex.DownloadCancel(cid); err != nil {
		t.Fatalf("Error when cancelling the download %s", err)
	}
}

func TestDownloadManifest(t *testing.T) {
	ctx := context.Background()
	codex := newCodexNode(t)
	cid, _ := uploadHelper(t, codex)

	manifest, err := codex.DownloadManifest(ctx, cid)
	if err != nil {
		t.Fatal("Error when downloading manifest:", err)
	}

	if manifest.Cid != cid {
		t.Errorf("expected cid %q, got %q", cid, manifest.Cid)
	}
}

func TestDownloadManifestWithNotExistingCid(t *testing.T) {
	ctx := context.Background()
	codex := newCodexNode(t, Config{BlockRetries: 1})

	manifest, err := codex.DownloadManifest(ctx, "bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku")
	if err == nil {
		t.Fatal("Error when downloading manifest:", err)
	}

	if manifest.Cid != "" {
		t.Errorf("expected empty cid, got %q", manifest.Cid)
	}
}

func TestDownloadInitWithNotExistingCid(t *testing.T) {
	ctx := context.Background()
	codex := newCodexNode(t, Config{BlockRetries: 1})

	if err := codex.DownloadInit(ctx, "bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku", DownloadInitOptions{}); err == nil {
		t.Fatal("expected error when initializing download for non-existent cid")
	}
}

func TestDownloadWithTwoNodes(t *testing.T) {
	var err error

	node1 := newCodexNode(t, Config{
		DiscoveryPort: 8100,
	})

	info1, err := node1.Debug()
	if err != nil {
		t.Fatal(err)
	}

	spr, err := node1.Spr()
	if err != nil {
		t.Fatalf("Failed to get bootstrap spr: %v", err)
	}

	t.Logf("spr: %s, info.spr: %s", spr, info1.Spr)

	bootstrapNodes := []string{info1.Spr}

	data := []byte("Hello World!")
	cid, _ := uploadData(t, node1, data)

	node2 := newCodexNode(t, Config{
		DiscoveryPort:  8101,
		BootstrapNodes: bootstrapNodes,
	})

	var buf bytes.Buffer
	options := DownloadStreamOptions{
		Writer: &buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = node2.DownloadStream(ctx, cid, options)

	if err != nil {
		t.Fatal("Error happened during download:", err.Error())
	}

	downloadedData := buf.Bytes()
	if string(downloadedData) != string(data) {
		t.Fatalf("Downloaded content does not match, expected %s got %s", data, downloadedData)
	}
}

func TestCancellingContextWhenNodeCannotResolveCID(t *testing.T) {
	var err error

	node1 := newCodexNode(t, Config{
		DiscoveryPort: 8100,
	})

	data := []byte("Hello World!")
	cid, _ := uploadData(t, node1, data)

	// Notice - no bootstrap nodes, so node2 cannot resolve the CID

	node2 := newCodexNode(t, Config{
		DiscoveryPort: 8101,
	})

	// Set a test-specific timeout to catch if DownloadStream hangs
	testTimeout := time.AfterFunc(10*time.Second, func() {
		panic("Test exceeded 10 second timeout - DownloadStream likely not respecting context cancellation")
	})
	defer testTimeout.Stop()

	var buf bytes.Buffer
	options := DownloadStreamOptions{
		Writer: &buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = node2.DownloadStream(ctx, cid, options)

	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}

	// The error should be context.DeadlineExceeded since node2 cannot resolve the CID
	if err != context.DeadlineExceeded && !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Got unexpected error): %v", err)
	}
}
