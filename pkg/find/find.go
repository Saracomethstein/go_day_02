package find

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Entries struct {
	Link bool
	Dir  bool
	File bool
	Ext  string
}

func Init(root string, settings Entries) {
	var wg sync.WaitGroup

	wg.Add(1)
	go findRecursive(root, &settings, &wg)

	wg.Wait()
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

func findRecursive(root string, settings *Entries, wg *sync.WaitGroup) {
	defer wg.Done()

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", root, err)
		return
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		info, err := os.Lstat(path)
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			if settings.Dir {
				fmt.Println(path)
			}

			wg.Add(1)
			go findRecursive(path, settings, wg)
		} else if info.Mode()&os.ModeSymlink != 0 && settings.Link {
			processSymlink(path)
		} else if settings.File && matchesExtension(path, settings.Ext) {
			fmt.Println(path)
		}
	}
}

func processSymlink(path string) {
	target, err := os.Readlink(path)
	if err != nil {
		fmt.Printf("Error reading symlink %s: %v\n", path, err)
		return
	}

	absoluteTarget := filepath.Join(filepath.Dir(path), target)
	if _, err := os.Stat(absoluteTarget); err == nil {
		fmt.Printf("%s -> %s\n", path, target)
	} else if os.IsNotExist(err) {
		fmt.Printf("%s -> [broken]\n", path)
	} else {
		fmt.Printf("Error resolving symlink %s: %v\n", path, err)
	}
}

func matchesExtension(path, ext string) bool {
	return ext == "" || strings.EqualFold(filepath.Ext(path), "."+ext)
}
