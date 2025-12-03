package list

import "sync"

type listNode struct {
	value interface{}
	prev  *listNode
	next  *listNode
}
type MyList struct {
	head  *listNode
	tail  *listNode
	len   int
	mutex sync.RWMutex
}

func NewMyList() *MyList {
	return &MyList{
		head:  nil,
		tail:  nil,
		len:   0,
		mutex: sync.RWMutex{},
	}
}
func (m *MyList) LPop() (val interface{}, exists bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.len == 0 {
		return nil, false
	}
	node := m.head
	val = node.value
	if m.len == 1 {
		m.head = nil
		m.tail = nil
	} else {
		m.head = node.next
		m.head.prev = nil
	}
	m.len--
	return val, true
}

func (m *MyList) RPop() (val interface{}, exists bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.len == 0 {
		return nil, false
	}
	node := m.tail
	val = node.value
	if m.len == 1 {
		m.tail = nil
		m.head = nil
	} else {
		m.tail = node.prev
		m.tail.next = nil
	}
	m.len--
	return val, true
}

func (m *MyList) LPush(val interface{}) (result int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	node := &listNode{
		value: val,
	}
	if m.len == 0 {
		m.head = node
		m.tail = node
	} else {
		node.next = m.head
		m.head.prev = node
		m.head = node
	}
	m.len++
	return m.len
}

func (m *MyList) RPush(val interface{}) (result int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	node := &listNode{
		value: val,
	}
	if m.len == 0 {
		m.head = node
		m.tail = node
	} else {
		node.prev = m.tail
		m.tail.next = node
		m.tail = node
	}
	m.len++
	return m.len
}

func (m *MyList) findNode(index int) *listNode {
	var start *listNode
	if index < m.len/2 {
		start = m.head
		for i := 0; i < index; i++ {
			start = start.next
		}
	} else {
		start = m.tail
		for i := m.len - 1; i > index; i-- {
			start = start.prev
		}
	}
	return start
}

func (m *MyList) Index(index int) (val interface{}, exists bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.len == 0 || index < 0 || index >= m.len {
		return nil, false
	}
	node := m.findNode(index)
	return node.value, true
}

func (m *MyList) Insert(index int, val interface{}) (result int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if index < 0 || index > m.len {
		return 0
	}
	node := &listNode{value: val}
	if index == 0 {
		if m.len == 0 {
			m.head = node
			m.tail = node
		} else {
			node.next = m.head
			m.head.prev = node
			m.head = node
		}
		m.len++
		return 1
	}
	if index == m.len {
		node.prev = m.tail
		m.tail.next = node
		m.tail = node
		m.len++
		return 1
	}
	current := m.findNode(index)
	prev := current.prev
	node.next = current
	node.prev = prev
	prev.next = node
	current.prev = node
	m.len++
	return 1
}

func (m *MyList) Range(start, stop int) []interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.len == 0 || start < 0 || stop >= m.len || start > stop {
		return nil
	}
	limit := stop - start + 1
	ans := make([]any, 0, limit)
	node := m.findNode(start)
	for i := 0; i < limit; i++ {
		ans = append(ans, node.value)
		node = node.next
	}
	return ans
}

func (m *MyList) Trim(start, stop int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.len == 0 || start < 0 || stop >= m.len || start > stop {
		return
	}
	if start == stop {
		node := m.findNode(start)
		m.head = node
		m.tail = node
		return
	}
	newHead := m.findNode(start)
	newTail := m.findNode(stop)
	m.len = stop - start + 1
	m.head = newHead
	m.head.prev = nil
	m.tail = newTail
	m.tail.next = nil
}
