package proloc

import (
	"os"
	"path"
	"strings"
	"sync"
)

func readFileStatsForDir(dirName string, ignore []string,
	wg *sync.WaitGroup, fsChan chan statsOrError,
	cMap *childrenMap, pMap *parentMap, ntMap *nodeTypeMap) {
	defer wg.Done()

	children, err := os.ReadDir(dirName)
	if err != nil {
		fsChan <- statsOrError{err: err}
	}

	for _, child := range children {
		childPath := path.Join(dirName, child.Name())
		if ignoreFile(childPath, ignore) || child.Type() == os.ModeSymlink {
			continue
		}
		cMap.addChild(dirName, childPath)
		pMap.addParent(childPath, dirName)

		if child.IsDir() {
			ntMap.addNodeType(childPath, dirType)
			wg.Add(1)
			go readFileStatsForDir(childPath, ignore, wg, fsChan, cMap, pMap, ntMap)
		} else {
			ntMap.addNodeType(childPath, fileType)
			wg.Add(1)
			go readStatsForFile(childPath, wg, fsChan)
		}
	}
}

func readStatsForFile(filename string, wg *sync.WaitGroup, fsChan chan statsOrError) {
	defer wg.Done()

	linesCount, err := calculateLinesCount(filename)
	if err != nil {
		fsChan <- statsOrError{err: err, fileStats: nil}
	}
	fsChan <- statsOrError{
		err:       nil,
		fileStats: &fileStats{linesCount: linesCount, filename: filename},
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
