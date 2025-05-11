package hashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type HashRing struct {
	mu       sync.RWMutex
	nodes    []uint32
	nodeMap  map[uint32]string
	replicas int
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		nodes:    []uint32{},
		nodeMap:  make(map[uint32]string),
		replicas: replicas,
	}
}

func Hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (hr *HashRing) AddNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for i := 0; i < hr.replicas; i++ {
		hash := Hash(fmt.Sprintf("%s-%d", node, i))
		hr.nodes = append(hr.nodes, hash)
		hr.nodeMap[hash] = node
	}
	sort.Slice(hr.nodes, func(i, j int) bool { return hr.nodes[i] < hr.nodes[j] })
}

func (hr *HashRing) GetNode(key string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if len(hr.nodes) == 0 {
		return ""
	}

	hsh := Hash(key)
	idx := sort.Search(len(hr.nodes),
		func(i int) bool { return hr.nodes[i] >= hsh })
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.nodeMap[hr.nodes[idx]]
}
