package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// check jira

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usaage %s <dir>", os.Args[0])
	}

	dir := flag.Arg(0)
	readDir(dir)
}

func readDir(dir string) {
	dh, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Could not open %s: %s", dir, err.Error())
	}
	defer dh.Close()

	for {
		fis, err := dh.Readdir(10)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not read dir names in %s: %s", dir, err.Error())
		}

		for _, fi := range fis {
			fmt.Printf("%s/%s\n", dir, fi.Name())
			if fi.IsDir() {
				readDir(dir + "/" + fi.Name())
			}
		}
	}
}
