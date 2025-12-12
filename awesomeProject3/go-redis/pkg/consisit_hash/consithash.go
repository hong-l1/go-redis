package consisithash

import (
	"github.com/cespare/xxhash/v2"
	"slices"
	"sort"
)

type NodeMap struct {
	nodesHash    []uint64
	nodesHashMap map[uint64]string
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		nodesHashMap: map[uint64]string{},
		nodesHash:    []uint64{},
	}
}
func (n *NodeMap) IsEmpty() bool {
	return len(n.nodesHash) == 0
}
func (n *NodeMap) AddNode(keys []string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := xxhash.Sum64String(key)
		n.nodesHashMap[hash] = key
		n.nodesHash = append(n.nodesHash, hash)
	}
	slices.Sort(n.nodesHash)
}
func (n *NodeMap) PickNode(keys string) string {
	if n.IsEmpty() {
		return ""
	}
	hash := xxhash.Sum64String(keys)
	index := sort.Search(len(n.nodesHash), func(i int) bool {
		return n.nodesHash[i] >= hash
	}) % len(n.nodesHash)
	return n.nodesHashMap[n.nodesHash[index]]
}
