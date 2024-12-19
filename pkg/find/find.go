package find

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Entries struct {
	Link bool
	Dir  bool
	File bool
	Ext  string
}

func Init(root string, settings Entries) {
	ch := make(chan struct{}, 10)
	defer close(ch)

	done := make(chan struct{})
	go func() {
		findRecursive(root, settings, ch)
		close(done)
	}()

	<-done
}

func ParseFlags() (Entries, string, error) {
	symlinkFlag := flag.Bool("sl", false, "output symlinks")
	dirFlag := flag.Bool("d", false, "output dirs")
	fileFlag := flag.Bool("f", false, "output files")
	extFlag := flag.String("ext", "", "accept file extension")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		return Entries{}, "", fmt.Errorf("Usage: ./find [flags] <directory>")
	}
	if *extFlag != "" && !*fileFlag {
		return Entries{}, "", fmt.Errorf("-ext flag requires -f flag")
	}

	entries := Entries{
		Link: *symlinkFlag,
		Dir:  *dirFlag,
		File: *fileFlag,
		Ext:  strings.TrimPrefix(*extFlag, "."),
	}

	if !*symlinkFlag && !*dirFlag && !*fileFlag {
		entries = Entries{
			Link: true,
			Dir:  true,
			File: true,
		}
	}

	return entries, args[0], nil
}

func findRecursive(root string, settings Entries, ch chan struct{}) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("Failed to read directory %s: %v", root, err)
		return
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil {
			log.Printf("Failed to get info for %s: %v", path, err)
			continue
		}

		if info.IsDir() {
			if settings.Dir {
				fmt.Println(path)
			}

			ch <- struct{}{}
			go func(subDir string) {
				defer func() { <-ch }()
				findRecursive(subDir, settings, ch)
			}(path)
		} else if settings.Link && (info.Mode()&os.ModeSymlink != 0) {
			processSymlink(path)
		} else if settings.File && matchesExtension(path, settings.Ext) {
			fmt.Println(path)
		}
	}
}

func processSymlink(path string) {
	target, err := os.Readlink(path)
	if err != nil {
		log.Printf("Failed to read symlink %s: %v", path, err)
		return
	}

	if _, err := os.Stat(target); err == nil {
		fmt.Printf("%s -> %s\n", path, target)
	} else if os.IsNotExist(err) {
		fmt.Printf("%s -> [broken]\n", path)
	} else {
		log.Printf("Failed to stat target %s: %v", target, err)
	}
}

func matchesExtension(path, ext string) bool {
	return ext == "" || strings.EqualFold(filepath.Ext(path), "."+ext)
}
