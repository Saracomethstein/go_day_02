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

type WCInfo struct {
	cLine   int
	cSymbol int
	cWord   int
	File    string
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

func wcRecursive(filename string, settings WCSettings, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	defer file.Close()

	info := WCInfo{
		File: filename,
	}
	err = getCounts(file, &info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erorr: %s\n", err)
		return
	}
	printCounts(settings, info)
}

func getCounts(file *os.File, info *WCInfo) error {
	inWord := false
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		char := scanner.Text()
		info.cSymbol++

		if char == "\n" {
			info.cLine++
		}

		if unicode.IsSpace(rune(char[0])) {
			if inWord {
				info.cWord++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
		info.cWord++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func printCounts(settings WCSettings, info WCInfo) {
	var output string
	if settings.Line {
		output += fmt.Sprintf("%d", info.cLine)
	} else if settings.Symbol {
		output += fmt.Sprintf("%d", info.cSymbol)
	} else if settings.Words {
		output += fmt.Sprintf("%d", info.cWord)
	}
	fmt.Println(output, info.File)
}
