package rotate_test

import (
	"archive/tar"
	"compress/gzip"
	"go_day_02/pkg/rotate"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return path
}

func TestProcessLogFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := createTestFile(t, tempDir, "test.log", "test content")
	archiveDir := t.TempDir()

	var wg sync.WaitGroup
	wg.Add(1)
	rotate.ProcessLogFile(logFile, archiveDir, &wg)
	wg.Wait()

	files, err := filepath.Glob(filepath.Join(archiveDir, "*.tar.gz"))
	if err != nil {
		t.Fatalf("failed to search for archives: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("no archive created in directory: %s", archiveDir)
	}

	checkTarGz(t, files[0], logFile, "test content")
}

func TestInit(t *testing.T) {
	tempDir := t.TempDir()
	logFile1 := createTestFile(t, tempDir, "log1.log", "log1 content")
	logFile2 := createTestFile(t, tempDir, "log2.log", "log2 content")
	archiveDir := t.TempDir()

	rotate.Init([]string{logFile1, logFile2}, archiveDir)

	files, err := filepath.Glob(filepath.Join(archiveDir, "*.tar.gz"))
	if err != nil {
		t.Fatalf("failed to search for archives: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 archives, got %d", len(files))
	}
}

func TestParseFlag(t *testing.T) {
	os.Args = []string{"./myRotate", "-a", "/archive", "file1.log", "file2.log"}
	files, archiveDir, err := rotate.ParseFlag()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	if *archiveDir != "/archive" {
		t.Fatalf("expected archive directory '/archive', got '%s'", *archiveDir)
	}
}

func checkTarGz(t *testing.T, archivePath, expectedFile, expectedContent string) {
	t.Helper()
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("failed to open archive: %v", err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	header, err := tarReader.Next()
	if err != nil {
		t.Fatalf("failed to read tar header: %v", err)
	}

	if header.Name != filepath.Base(expectedFile) {
		t.Fatalf("expected file name '%s', got '%s'", filepath.Base(expectedFile), header.Name)
	}

	content, err := io.ReadAll(tarReader)
	if err != nil {
		t.Fatalf("failed to read file content from tar: %v", err)
	}

	if string(content) != expectedContent {
		t.Fatalf("expected content '%s', got '%s'", expectedContent, content)
	}
}
