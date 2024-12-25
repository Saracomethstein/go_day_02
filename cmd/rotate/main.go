package main

import (
	"fmt"
	"go_day_02/pkg/rotate"
	"os"
)

func main() {
	files, archDir, err := rotate.ParseFlag()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	rotate.Init(files, *archDir)
}
