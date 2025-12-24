package storage

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"
)

const expectedCID = "zDvZRwzmAkhzDRPH5EW242gJBNZ2T7aoH2v1fVH66FxXL4kSbvyM"

func TestUploadReader(t *testing.T) {
	storage := newStorageNode(t)
	totalBytes := 0
	finalPercent := 0.0

	buf := bytes.NewBuffer([]byte("Hello World!"))
	len := buf.Len()
	cid, err := storage.UploadReader(context.Background(), UploadOptions{Filepath: "hello.txt", OnProgress: func(read, total int, percent float64, err error) {
		if err != nil {
			log.Fatalf("Error happened during upload: %v\n", err)
		}

		totalBytes = total
		finalPercent = percent
	}}, buf)

	if err != nil {
		t.Fatalf("UploadReader failed: %v", err)
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback read %d bytes but expected %d", totalBytes, len)
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}
}

func TestUploadReaderCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	storage := newStorageNode(t)
	buf := bytes.NewBuffer(make([]byte, 1024*1024*10))

	channelErr := make(chan error, 1)
	go func() {
		_, e := storage.UploadReader(ctx, UploadOptions{Filepath: "hello.txt"}, buf)
		channelErr <- e
	}()

	cancel()
	err := <-channelErr

	if err == nil {
		t.Fatal("UploadReader should have been canceled")
	}

	if err.Error() != context.Canceled.Error() {
		t.Fatalf("UploadReader returned unexpected error: %v expected %v", err, context.Canceled)
	}
}

func TestUploadFile(t *testing.T) {
	storage := newStorageNode(t)
	totalBytes := 0
	finalPercent := 0.0

	stat, err := os.Stat("./testdata/hello.txt")
	if err != nil {
		log.Fatalf("Error happened during file stat: %v\n", err)
	}

	options := UploadOptions{Filepath: "./testdata/hello.txt", OnProgress: func(read, total int, percent float64, err error) {
		if err != nil {
			log.Fatalf("Error happened during upload: %v\n", err)
		}

		totalBytes = total
		finalPercent = percent
	}}

	cid, err := storage.UploadFile(context.Background(), options)
	if err != nil {
		t.Fatalf("UploadReader failed: %v", err)
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}

	if totalBytes != int(stat.Size()) {
		t.Fatalf("UploadReader progress callback read %d bytes but expected %d", totalBytes, int(stat.Size()))
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}
}

func TestUploadFileCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	tmpFile, err := os.Create(os.TempDir() + "/large_file.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	largeContent := make([]byte, 1024*1024*50)
	if _, err := tmpFile.Write(largeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	storage := newStorageNode(t)

	channelError := make(chan error, 1)
	go func() {
		_, err := storage.UploadFile(ctx, UploadOptions{Filepath: tmpFile.Name()})
		channelError <- err
	}()

	cancel()
	err = <-channelError

	if err == nil {
		t.Fatal("UploadFile should have been canceled")
	}

	if err.Error() != context.Canceled.Error() {
		t.Fatalf("UploadFile returned unexpected error: %v", err)
	}
}

func TestUploadFileNoProgress(t *testing.T) {
	storage := newStorageNode(t)

	options := UploadOptions{Filepath: "./testdata/doesnt_exist.txt"}

	cid, err := storage.UploadFile(context.Background(), options)
	if err == nil {
		t.Fatalf("UploadReader should have failed")
	}

	if cid != "" {
		t.Fatalf("Cid should be empty but got %s", cid)
	}
}

func TestManualUpload(t *testing.T) {
	storage := newStorageNode(t)

	sessionId, err := storage.UploadInit(&UploadOptions{Filepath: "hello.txt"})
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	err = storage.UploadChunk(sessionId, []byte("Hello "))
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	err = storage.UploadChunk(sessionId, []byte("World!"))
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	cid, err := storage.UploadFinalize(sessionId)
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}
}
