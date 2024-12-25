package xargs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Init() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: <program> <command> [args...]")
		os.Exit(1)
	}

	command := os.Args[1]
	cmdArgs := os.Args[2:]

	scanner := bufio.NewScanner(os.Stdin)
	var execErrors []string

	for scanner.Scan() {
		arg := strings.TrimSpace(scanner.Text())
		if arg == "" {
			continue
		}

		if err := executeCommand(command, cmdArgs, arg); err != nil {
			execErrors = append(execErrors, fmt.Sprintf("arg '%s': %v", arg, err))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	if len(execErrors) > 0 {
		fmt.Fprintln(os.Stderr, "Errors occurred:")
		for _, err := range execErrors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func executeCommand(command string, cmdArgs []string, arg string) error {
	cmd := exec.Command(command, append(cmdArgs, arg)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
