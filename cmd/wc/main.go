package main

import (
	"fmt"
	"go_day_02/pkg/wc"
	"os"
)

func main() {
	settings, files, err := wc.FlagParce()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	wc.Init(files, settings)
}
