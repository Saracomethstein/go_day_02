package wc_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWCUtility(t *testing.T) {
	testDir, err := ioutil.TempDir("", "wc_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	testFiles := map[string]string{
		"file1.txt": "Hello world\nThis is a test file\n",
		"file2.txt": "Another test file\nWith multiple lines\nAnd more words\n",
		"empty.txt": "",
	}

	for name, content := range testFiles {
		filePath := filepath.Join(testDir, name)
		if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write to %s: %v", filePath, err)
		}
	}

	testCases := []struct {
		args []string
	}{
		{args: []string{"-l", filepath.Join(testDir, "file1.txt")}},
		{args: []string{"-w", filepath.Join(testDir, "file2.txt")}},
		{args: []string{"-m", filepath.Join(testDir, "file1.txt")}},
		{args: []string{"-l", filepath.Join(testDir, "empty.txt")}},
	}

	for _, tc := range testCases {
		customWCOut := runCustomWC(t, tc.args)
		originalWCOut := runOriginalWC(t, tc.args)

		if customWCOut != originalWCOut {
			t.Errorf("Test failed for args %v: expected %q, got %q", tc.args, originalWCOut, customWCOut)
		} else {
			t.Logf("Success for args %v", tc.args)
		}
	}
}

func runCustomWC(t *testing.T, args []string) string {
	cmd := exec.Command("./../../build/wc", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("Custom wc failed for %v: %v", args, err)
	}

	return strings.TrimSpace(out.String())
}

func runOriginalWC(t *testing.T, args []string) string {
	cmd := exec.Command("wc", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("Original wc failed for %v: %v", args, err)
	}

	return strings.TrimSpace(out.String())
}
