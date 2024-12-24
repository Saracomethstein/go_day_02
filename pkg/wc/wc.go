package wc

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"unicode"
)

type WCSettings struct {
	Line   bool
	Symbol bool
	Words  bool
}

func Init(files []string, settings WCSettings) {
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go wcRecursive(file, settings, &wg)
	}

	wg.Wait()
}

func FlagParce() (WCSettings, []string, error) {
	lineFlag := flag.Bool("l", false, "output line count")
	symbolFlag := flag.Bool("m", false, "output symbol count")
	wordsFlag := flag.Bool("w", false, "output word count")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		return WCSettings{}, []string{}, fmt.Errorf("Usage: ./wc [flags] <directory>")
	}

	settings := WCSettings{
		Line:   *lineFlag,
		Symbol: *symbolFlag,
		Words:  *wordsFlag,
	}

	if !*lineFlag && !*symbolFlag && !*wordsFlag {
		settings = WCSettings{
			Line:   false,
			Symbol: false,
			Words:  true,
		}
	}

	return settings, args[0:], nil
}

func wcRecursive(file string, settings WCSettings, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	defer f.Close()

	var cLine, cSymbol, cWord int
	inWord := false

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		char := scanner.Text()
		cSymbol++

		if char == "\n" {
			cLine++
		}

		if unicode.IsSpace(rune(char[0])) {
			if inWord {
				cWord++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
		cWord++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}

	var output string
	if settings.Line {
		output += fmt.Sprintf("%d", cLine)
	} else if settings.Symbol {
		output += fmt.Sprintf("%d", cSymbol)
	} else if settings.Words {
		output += fmt.Sprintf("%d", cWord)
	}

	fmt.Println(output, "\t", file)
}
