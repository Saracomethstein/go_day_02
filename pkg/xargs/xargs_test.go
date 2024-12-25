package xargs_test

import (
	"bytes"
	"os/exec"
	"testing"
)

func runTest(t *testing.T, description, commandOrg string, commandXargs string) {
	t.Helper()
	t.Logf("Running test: %s", description)

	originalCmd := exec.Command(commandOrg)
	var originalOut bytes.Buffer
	originalCmd.Stdout = &originalOut
	originalCmd.Stderr = &originalOut
	originalErr := originalCmd.Run()

	xargsCmd := exec.Command(commandXargs)
	var xargsOut bytes.Buffer
	xargsCmd.Stdout = &xargsOut
	xargsCmd.Stderr = &xargsOut
	xargsErr := xargsCmd.Run()

	if originalErr != nil && xargsErr == nil || originalErr == nil && xargsErr != nil {
		t.Fatalf("Status mismatch:\nOriginal error: %v\nXargs error: %v", originalErr, xargsErr)
	}
	if originalOut.String() != xargsOut.String() {
		t.Fatalf("Output mismatch:\nOriginal:\n%s\nXargs:\n%s", originalOut.String(), xargsOut.String())
	}
	t.Logf("Test passed: %s", description)
}

func TestXargs(t *testing.T) {
	tests := []struct {
		description  string
		commandOrg   string
		commandXargs string
	}{
		{
			description:  "List directories with ls",
			commandOrg:   "ls -l ./../../build ./../../cmd ./../../tests",
			commandXargs: "echo -e \"./../../build\n./../../cmd\n./../../tests\" | ./../../build/xargs ls -l",
		},
		{
			description:  "Count lines in file with wc",
			commandOrg:   "wc -l ./../../Makefile ./../../README.md",
			commandXargs: "echo -e \"./../../Makefile\n./../../README.md\" | ./../../build/xargs wc -l",
		},
		{
			description:  "Concatenate files with cat",
			commandOrg:   "cat ./../../Makefile ./../../README.md",
			commandXargs: "echo -e \"./../../Makefile\n./../../README.md\" | ./../../build/xargs cat",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			runTest(t, test.description, test.commandOrg, test.commandXargs)
		})
	}
}
