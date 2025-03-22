package proloc

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

type ProlocConfig struct {
	Project string
	Ignore  []string
}

type fileStats struct {
	parentDir  string
	filename   string
	linesCount uint64
}

type dirStats struct {
	dirName   string
	parentDir *string
}

type statsOrError struct {
	err       error
	fileStats *fileStats
	dirStats  *dirStats
}

func CountLines(config ProlocConfig) error {
	fsChan := make(chan statsOrError)

	// mapper routine
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go readFileStatsForDir(config.Project, nil, config.Ignore, wg, fsChan)
		wg.Wait()
		close(fsChan)
	}()

	// reducer
	var count uint64 = 0
	totalFiles := 0
	dirsLineCount := map[string]uint64{}
	for fs := range fsChan {
		if fs.err != nil {
			return fs.err
		}
		if fs.fileStats != nil {
			if _, ok := dirsLineCount[fs.fileStats.parentDir]; !ok {
				dirsLineCount[fs.fileStats.parentDir] = 0
			}
			count += fs.fileStats.linesCount
			dirsLineCount[fs.fileStats.parentDir] += fs.fileStats.linesCount
			totalFiles += 1
		}
	}
	fmt.Printf("total dirs = %s | total files = %s | total lines count = %s\n",
		formatNumber(uint64(len(dirsLineCount))), formatNumber(uint64(totalFiles)), formatNumber(count))
	return nil
}

func readFileStatsForDir(dirName string, parentDir *string, ignore []string, wg *sync.WaitGroup, fsChan chan statsOrError) {
	defer wg.Done()

	fsChan <- statsOrError{
		dirStats: &dirStats{
			dirName:   dirName,
			parentDir: parentDir,
		},
	}

	dir, err := os.ReadDir(dirName)
	if err != nil {
		fsChan <- statsOrError{err: err}
	}

	for _, child := range dir {
		childPath := path.Join(dirName, child.Name())
		if ignoreFile(childPath, ignore) {
			continue
		}
		if child.IsDir() {
			wg.Add(1)
			go readFileStatsForDir(childPath, parentDir, ignore, wg, fsChan)
		} else {
			if child.Type() != os.ModeSymlink {
				wg.Add(1)
				go readStatsForFile(childPath, dirName, wg, fsChan)
			}
		}
	}
}

func readStatsForFile(filename, parentDir string, wg *sync.WaitGroup, fsChan chan statsOrError) {
	defer wg.Done()

	linesCount := calculateLinesCount(filename)
	fsChan <- statsOrError{
		err: nil,
		fileStats: &fileStats{
			linesCount: linesCount,
			filename:   filename,
			parentDir:  parentDir,
		},
	}
}

func ignoreFile(path string, ignore []string) bool {
	for _, ignorePattern := range ignore {
		if strings.Contains(path, ignorePattern) {
			return true
		}
	}
	return false
}
