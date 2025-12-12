package zset

import (
	"math/rand"
	"sync"
	"time"
)

const (
	maxLevel = 32
)

type Element struct {
	Member string
	Score  float64
}
type Node struct {
	data *Element
	prev *Node
	next []*Node
	span []uint32
}
type skipList struct {
	head   *Node //哨兵节点不参与
	tail   *Node //指向最后一个元素
	length uint32
	level  int8 //level-1是有效下标
}
type MyZSet struct {
	dict map[string]*Node
	sl   *skipList
	mu   sync.RWMutex
}

func NewMyZSet() *MyZSet {
	rand.Seed(time.Now().UnixNano())
	return &MyZSet{
		dict: make(map[string]*Node),
		sl:   newSkipList(),
	}
}
func randomLevel() int8 {
	level := 1
	for rand.Float64() < 0.5 && level < maxLevel {
		level++
	}
	return int8(level)
}
func newSkipList() *skipList {
	head := &Node{
		next: make([]*Node, maxLevel),
		span: make([]uint32, maxLevel),
	}
	return &skipList{
		head:  head,
		tail:  nil,
		level: 1,
	}
}
func (s *skipList) insert(member string, score float64) *Node {
	head := s.head
	update := make([]*Node, maxLevel)
	rank := make([]uint32, maxLevel)
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for head.next[i] != nil && (head.next[i].data.Score < score ||
			head.next[i].data.Score == score && head.next[i].data.Member < member) {
			rank[i] += head.span[i]
			head = head.next[i]
		}
		update[i] = head
	}
	level := randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			rank[i] = 0
			update[i] = s.head
			update[i].span[i] = s.length
		}
		s.level = level
	}
	newNode := &Node{
		data: &Element{Member: member, Score: score},
		next: make([]*Node, level),
		span: make([]uint32, level),
	}
	for i := 0; i < int(level); i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
		newNode.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].span[i] = rank[0] - rank[i] + 1
	}
	for j := level; j < s.level; j++ {
		update[j].span[j]++
	}
	if update[0] == s.head {
		newNode.next = nil
	} else {
		newNode.prev = update[0]
	}
	if newNode.next[0] != nil {
		newNode.next[0].prev = newNode
	} else {
		s.tail = newNode
	}
	s.length++
	return newNode
}
func (s *skipList) delete(member string, score float64) bool {
	head := s.head
	update := make([]*Node, maxLevel)
	for i := s.level - 1; i >= 0; i-- {
		for head.next[i] != nil && (head.next[i].data.Score < score ||
			head.next[i].data.Score == score && head.next[i].data.Member < member) {
			update[i] = head.next[i]
		}
		if head.next[i] != nil &&
			head.next[i].data.Score == score &&
			head.next[i].data.Member == member {
			s.deleteNode(head, update)
			return true
		}
	}
	return false
}
func (s *skipList) deleteNode(node *Node, update []*Node) {
	for i := int8(0); i < s.level; i++ {
		if update[i].next[i] == node {
			update[i].next[i] = node.next[i]
			update[i].span[i] += node.span[i] - 1
		} else {
			update[i].span[i]--
		}
	}
	if node.next[0] != nil {
		node.next[0].prev = node.prev
	} else {
		s.tail = node.prev
	}
	if s.level > 1 && s.head.next[s.level-1] == nil {
		s.level--
	}
	s.length--
}
func (s *skipList) getByRank(rank uint32) *Node {
	var step uint32
	head := s.head
	for i := s.level - 1; i >= 0; i-- {
		for head.next[i] != nil && (step+head.next[i].span[i] <= rank) {
			step += head.span[i]
			head = head.next[i]
		}
		if step == rank {
			return head.next[i]
		}
	}
	return nil
}
func (s *skipList) getRank(member string, score float64) (uint32, bool) {
	var step uint32
	head := s.head
	for i := s.level - 1; i >= 0; i-- {
		for head.next[i] != nil && (head.next[i].data.Score < score ||
			head.next[i].data.Score == score && head.next[i].data.Member < member) {
			step += head.span[i]
			head = head.next[i]
		}
		if head.next[i] != nil &&
			head.next[i].data.Score == score &&
			head.next[i].data.Member == member {
			return step + head.span[i], true
		}
	}
	return 0, false
}
func (s *skipList) getFirstInScoreRange(min *ScoreBorder) *Node {
	head := s.head
	for i := s.level - 1; i >= 0; i-- {
		for head.next[i] != nil {
			score := head.next[i].data.Score
			if score < min.Value || (min.Exclude && score == min.Value) {
				head = head.next[i]
			} else {
				break
			}
		}
	}
	return head.next[0]
}
