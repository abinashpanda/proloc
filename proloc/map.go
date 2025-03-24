package proloc

import "sync"

type childrenMap struct {
	sync.Mutex
	children map[string][]string
}

func initChildrenMap() *childrenMap {
	return &childrenMap{children: make(map[string][]string)}
}

func (cMap *childrenMap) addChild(parent, child string) {
	cMap.Lock()
	if _, ok := cMap.children[parent]; !ok {
		cMap.children[parent] = []string{child}
	} else {
		cMap.children[parent] = append(cMap.children[parent], child)
	}
	cMap.Unlock()
}

func (cMap *childrenMap) getChildren(parent string) []string {
	cMap.Lock()
	defer cMap.Unlock()
	return cMap.children[parent]
}

type parentMap struct {
	sync.Mutex
	parent map[string]string
}

func initParentMap() *parentMap {
	return &parentMap{parent: make(map[string]string)}
}

func (pMap *parentMap) addParent(child, parent string) {
	pMap.Lock()
	pMap.parent[child] = parent
	pMap.Unlock()
}

func (pMap *parentMap) getParent(child string) string {
	pMap.Lock()
	defer pMap.Unlock()
	return pMap.parent[child]
}

type lineCountMap struct {
	sync.Mutex
	lineCount map[string]uint64
}

func initLineCountMap() *lineCountMap {
	return &lineCountMap{lineCount: make(map[string]uint64)}
}

func (lcMap *lineCountMap) addLineCount(node string, lineCount uint64) {
	lcMap.Lock()
	if _, ok := lcMap.lineCount[node]; !ok {
		lcMap.lineCount[node] = lineCount
	} else {
		lcMap.lineCount[node] += lineCount
	}
	lcMap.Unlock()
}

func (lcMap *lineCountMap) getLineCount(node string) uint64 {
	lcMap.Lock()
	defer lcMap.Unlock()
	return lcMap.lineCount[node]
}

type nodeType uint8

const (
	fileType = iota
	dirType
)

type nodeTypeMap struct {
	sync.Mutex
	nodeType map[string]nodeType
}

func initNodeTypeMap() *nodeTypeMap {
	return &nodeTypeMap{nodeType: make(map[string]nodeType)}
}

func (ntMap *nodeTypeMap) addNodeType(node string, nt nodeType) {
	ntMap.Lock()
	ntMap.nodeType[node] = nt
	ntMap.Unlock()
}

func (ntMap *nodeTypeMap) getNodeType(node string) nodeType {
	ntMap.Lock()
	defer ntMap.Unlock()
	return ntMap.nodeType[node]
}
