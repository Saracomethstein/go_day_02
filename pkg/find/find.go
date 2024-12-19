package find

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Entries struct {
	Link bool
	Dir  bool
	File bool
	Ext  string
}

func Init(root string, ent *Entries) {
	var wg sync.WaitGroup

	wg.Add(1)
	go find(root, ent, &wg)

	wg.Wait()
}

func FlagParse() (Entries, string, error) {
	symlinkFlag := flag.Bool("sl", false, "output symlinks")
	dirFlag := flag.Bool("d", false, "output dirs")
	fileFlag := flag.Bool("f", false, "output files")
	extFlag := flag.String("ext", "", "accept file extension")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || (*extFlag != "" && !(*fileFlag)) {
		return Entries{}, "", fmt.Errorf("Usage: ./find <flags: -sl | -d | -f | -ext (reqiered f)> <directory>")
	}

	directory := args[0]
	var ent Entries
	if !(*symlinkFlag) && !(*dirFlag) && !(*fileFlag) {
		ent = Entries{
			Dir:  true,
			File: true,
			Link: true,
			Ext:  "",
		}
	} else {
		ent = Entries{
			Dir:  *dirFlag,
			File: *fileFlag,
			Link: *symlinkFlag,
			Ext:  *extFlag,
		}
	}
	return ent, directory, nil
}

func find(root string, ent *Entries, wg *sync.WaitGroup) {
	defer wg.Done()

	entries, err := os.ReadDir(root)

	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		info, err := os.Lstat(path)

		if err != nil {
			continue
		}

		if info.IsDir() {
			if ent.Dir {
				fmt.Println(path)
			}
			wg.Add(1)
			go find(path, ent, wg)
		} else if info.Mode()&os.ModeSymlink != 0 && ent.Link {
			target, err := os.Readlink(path)

			if err != nil {
				continue
			}

			if _, err := os.Stat(target); err == nil {
				fmt.Println(path, " -> ", target)
			} else if os.IsNotExist(err) {
				fmt.Println(path, " -> ", " [broken] ")
			} else {
				continue
			}
		} else {
			if ent.File && (ent.Ext == "" || ent.Ext != "" && ("."+ent.Ext) == filepath.Ext(path)) {
				fmt.Println(path)
				continue
			}
		}
	}
}
