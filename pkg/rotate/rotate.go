package rotate

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func Init(files []string, archiveDir string) {
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go ProcessLogFile(file, archiveDir, &wg)
	}
	wg.Wait()
}

func ParseFlag() ([]string, *string, error) {
	archiveDir := flag.String("a", "", "Directory to store archives (optional)")
	flag.Parse()

	if flag.NArg() < 1 {
		return []string{}, nil, fmt.Errorf("Usage: ./myRotate [-a archive_dir] log_file1 [log_file2 ...]")
	}
	logFiles := flag.Args()

	return logFiles, archiveDir, nil
}

func ProcessLogFile(logFile, archiveDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	fileInfo, err := os.Stat(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to stat file: %s", err)
		return
	}

	if !fileInfo.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "not a regular file: %s", logFile)
		return
	}

	modTime := fileInfo.ModTime().Unix()

	archiveName := fmt.Sprintf("%s_%d.tar.gz", filepath.Base(logFile), modTime)
	if archiveDir != "" {
		archiveName = filepath.Join(archiveDir, archiveName)
	}

	if err := CreateTarGz(archiveName, logFile, fileInfo); err != nil {
		fmt.Fprintf(os.Stderr, "unable to create archive: %s", err)
		return
	}

	fmt.Printf("Created archive: %s\n", archiveName)
}

func CreateTarGz(archiveName, logFile string, fileInfo os.FileInfo) error {
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("unable to create archive file: %w", err)
	}
	defer archiveFile.Close()

	gzipWriter := gzip.NewWriter(archiveFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	logFileReader, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("unable to open log file: %w", err)
	}
	defer logFileReader.Close()

	header := &tar.Header{
		Name:    filepath.Base(logFile),
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode().Perm()),
		ModTime: fileInfo.ModTime(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("unable to write tar header: %w", err)
	}

	if _, err := io.Copy(tarWriter, logFileReader); err != nil {
		return fmt.Errorf("unable to write log file to tar: %w", err)
	}

	return nil
}
