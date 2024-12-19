package main

import (
	"fmt"
	"go_day_02/pkg/find"
	"os"
)

func main() {
	ent, directory, err := find.FlagParse()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	find.Init(directory, &ent)
}
