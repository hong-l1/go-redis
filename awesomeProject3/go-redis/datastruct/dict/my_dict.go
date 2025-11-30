package dict

import (
	"github.com/cespare/xxhash/v2"
	"math/rand"
	"sync"
)

const DefaultShardCount = 256

type MyDict struct {
	shards    []*Shard
	shardMask uint64
}

//通过业务的key计算hash,来分配执行当前业务的shard
func NewDict() *MyDict {
	d := &MyDict{
		shards:    make([]*Shard, DefaultShardCount),
		shardMask: uint64(DefaultShardCount - 1),
	}
	for i := 0; i < DefaultShardCount; i++ {
		d.shards[i] = NewShard()
	}
	return d
}
func (d *MyDict) getShard(key string) *Shard {
	hash := getHash(key)
	return d.shards[hash&d.shardMask]
}

//getHash xxhash 性能最好；分布足够均匀；无GC开销
func getHash(key string) uint64 {
	return xxhash.Sum64String(key)
}

type Shard struct {
	sync.RWMutex
	ht        [2]*hTable
	rehashIdx int64
}

func NewShard() *Shard {
	s := &Shard{
		rehashIdx: -1,
	}
	s.ht[0] = &hTable{}
	s.ht[0].init(4)
	return s
}
func (s *Shard) expand() {
	if s.rehashIdx != -1 {
		return
	}
	if s.ht[0].used >= s.ht[0].size {
		newSize := s.ht[0].size * 2
		s.ht[1] = &hTable{}
		s.ht[1].init(newSize)
		s.rehashIdx = 0
	}
}
func (s *Shard) rehashStep() {
	if s.rehashIdx == -1 {
		return
	}
	idx := s.rehashIdx
	entry := s.ht[0].buckets[idx]
	for entry != nil {
		nextEntry := entry.Next
		hash := getHash(entry.Key)
		idx1 := hash & s.ht[1].mask
		entry.Next = s.ht[1].buckets[idx1]
		s.ht[1].buckets[idx1] = entry
		s.ht[0].used--
		s.ht[1].used++
		entry = nextEntry
	}
	s.ht[0].buckets[idx] = nil
	s.rehashIdx++
	if s.rehashIdx >= int64(s.ht[0].size) {
		s.ht[0] = s.ht[1]
		s.ht[1] = nil
		s.rehashIdx = -1
	}
}
func (s *Shard) getRandomKeyFromShard() (string, bool) {
	var htIdx int
	if s.rehashIdx >= 0 {
		if rand.Intn(2) == 0 {
			htIdx = 1
		} else {
			htIdx = 0
		}
	} else {
		htIdx = 0
	}
	for i := 0; i < 8; i++ {
		idx := rand.Uint64() & s.ht[htIdx].mask
		entry := s.ht[htIdx].buckets[idx]
		if entry == nil {
			continue
		}
		cnt := 0
		for head := entry; head != nil; head = head.Next {
			cnt++
			if rand.Intn(cnt) == 0 {
				return head.Key, true
			}
		}
	}
	return "", false
}

type hTable struct {
	buckets []*Entry
	mask    uint64
	size    uint64
	used    uint64
}

func (h *hTable) lookup(key string, hash uint64) (interface{}, bool) {
	if h.used == 0 {
		return nil, false
	}
	idx := hash & h.mask
	entry := h.buckets[idx]
	for entry != nil {
		if entry.Key == key {
			return entry.Val, true
		}
		entry = entry.Next
	}
	return nil, false
}
func (h *hTable) update(key string, val interface{}, hash uint64) bool {
	if h.used == 0 {
		return false
	}
	idx := hash & h.mask
	entry := h.buckets[idx]
	for entry != nil {
		if entry.Key == key {
			entry.Val = val
			return true
		}
		entry = entry.Next
	}
	return false
}
func (h *hTable) del(key string, hash uint64) (any, bool) {
	if h.used == 0 {
		return nil, false
	}
	idx := hash & h.mask
	var prev *Entry
	entry := h.buckets[idx]
	for entry != nil {
		if entry.Key == key {
			if prev == nil {
				h.buckets[idx] = entry.Next
			} else {
				prev.Next = entry.Next
			}
			h.used--
			return entry.Val, true
		}
		prev = entry
		entry = entry.Next
	}
	return nil, false
}
func (h *hTable) init(size uint64) {
	h.buckets = make([]*Entry, size)
	h.size = size
	h.mask = size - 1
	h.used = 0
}

type Entry struct {
	Key  string
	Val  any
	Next *Entry
}
