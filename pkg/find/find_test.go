package find_test

import (
	"bytes"
	"go_day_02/pkg/find"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestFindProgram(t *testing.T) {
	rootDir := "../../test_dir"

	tests := []struct {
		name     string
		settings find.Entries
		expected []string
	}{
		{
			name: "find all files and directories",
			settings: find.Entries{
				Dir:  true,
				File: true,
				Link: true,
			},
			expected: []string{
				"../../test_dir/dir1",
				"../../test_dir/dir1/file1.txt",
				"../../test_dir/dir1/file2.log",
				"../../test_dir/dir2",
				"../../test_dir/dir2/file3.txt",
				"../../test_dir/dir2/subdir",
				"../../test_dir/dir2/subdir/file4.md",
				"../../test_dir/symlink1 -> dir1/file1.txt",
				"../../test_dir/broken_symlink -> [broken]",
			},
		},
		{
			name: "find only directories",
			settings: find.Entries{
				Dir: true,
			},
			expected: []string{
				"../../test_dir/dir1",
				"../../test_dir/dir2",
				"../../test_dir/dir2/subdir",
			},
		},
		{
			name: "find only .txt files",
			settings: find.Entries{
				File: true,
				Ext:  "txt",
			},
			expected: []string{
				"../../test_dir/dir2/file3.txt",
				"../../test_dir/dir1/file1.txt",
			},
		},
		{
			name: "find symlinks",
			settings: find.Entries{
				Link: true,
			},
			expected: []string{
				"../../test_dir/symlink1 -> dir1/file1.txt",
				"../../test_dir/broken_symlink -> [broken]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}

			oldStdout := os.Stdout
			os.Stdout = w

			done := make(chan struct{})
			go func() {
				defer close(done)
				find.Init(rootDir, tt.settings)
				w.Close()
			}()

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("failed to read from pipe: %v", err)
			}
			r.Close()

			os.Stdout = oldStdout

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			sort.Strings(lines)

			sort.Strings(tt.expected)
			if !equal(lines, tt.expected) {
				t.Errorf("unexpected output: got %v, want %v", lines, tt.expected)
			}
		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
