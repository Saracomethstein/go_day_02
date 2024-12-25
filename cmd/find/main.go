package main

import (
	"fmt"
	"go_day_02/pkg/find"
	"os"
)

func main() {
	settings, filename, err := find.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	find.Init(filename, settings)
}
