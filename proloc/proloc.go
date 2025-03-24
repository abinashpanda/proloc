package proloc

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
)

type ProlocConfig struct {
	Project  string
	Ignore   []string
	MaxDepth uint64
}

type fileStats struct {
	filename   string
	linesCount uint64
}

type statsOrError struct {
	err       error
	fileStats *fileStats
}

func CountLines(config ProlocConfig) error {
	fsChan := make(chan statsOrError)
	cMap := initChildrenMap()
	pMap := initParentMap()

	ntMap := initNodeTypeMap()
	ntMap.addNodeType(config.Project, dirType)

	// mapper routine
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go readFileStatsForDir(config.Project, config.Ignore, wg, fsChan, cMap, pMap, ntMap)
		wg.Wait()
		close(fsChan)
	}()

	// reducer
	lcMap := initLineCountMap()
	for fs := range fsChan {
		if fs.err != nil {
			return fs.err
		}
		if fs.fileStats != nil {
			lineCount := fs.fileStats.linesCount
			fileName := fs.fileStats.filename
			lcMap.addLineCount(fileName, lineCount)
			parent := pMap.getParent(fileName)
			for parent != "" {
				lcMap.addLineCount(parent, lineCount)
				parent = pMap.getParent(parent)
			}
		}
	}

	totalFiles := 0
	totalDirs := 0
	fmt.Println("Line Count\t\tName")
	fmt.Println("----------\t\t----")
	formatLineCount(config.Project, "", 0, config.MaxDepth, cMap, lcMap, ntMap, &totalFiles, &totalDirs)
	fmt.Printf("lines count = %s, files = %s, dirs = %s\n",
		humanize.Comma(int64(lcMap.getLineCount(config.Project))),
		humanize.Comma(int64(totalFiles)), humanize.Comma(int64(totalDirs)))

	return nil
}

func formatLineCount(node string, parent string,
	depth int, maxDepth uint64,
	cMap *childrenMap, lcMap *lineCountMap, ntMap *nodeTypeMap,
	totalFiles, totalDirs *int) {
	if maxDepth != 0 && depth > int(maxDepth) {
		return
	}

	if ntMap.getNodeType(node) == fileType {
		*totalFiles += 1
	} else if ntMap.getNodeType(node) == dirType {
		*totalDirs += 1
	}

	nodeName := node
	if parent != "" {
		nodeName = strings.Replace(nodeName, fmt.Sprintf("%s/", parent), "", 1)
	}
	fmt.Printf("%s\t\t\t%s%s\n", humanize.Comma(int64(lcMap.getLineCount(node))), strings.Repeat("  ", depth), nodeName)

	if ntMap.getNodeType(node) == fileType {
		return
	}

	for _, child := range cMap.getChildren(node) {
		formatLineCount(child, node, depth+1, maxDepth, cMap, lcMap, ntMap, totalFiles, totalDirs)
	}
}
