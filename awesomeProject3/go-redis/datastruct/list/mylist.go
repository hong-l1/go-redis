package list

import (
	"container/list"
	"encoding/binary"
	"sync"
)

const bufSize = 8 * 1024

type MyList struct {
	data *list.List
	size int
	mu   sync.RWMutex
}
type listPack struct {
	buf []byte
	cnt int
}

func NewMyList() *MyList {
	return &MyList{
		data: list.New(),
		size: 0,
	}
}
func newListPack() *listPack {
	lp := &listPack{
		buf: make([]byte, 4, bufSize),
		cnt: 0,
	}
	binary.LittleEndian.PutUint32(lp.buf[0:4], 4)
	return lp
}
func (l *listPack) totalBytes() int {
	return len(l.buf)
}
func (l *listPack) insert(val []byte) {
	entryLen := 2 + len(val)
	entry := make([]byte, entryLen)
	binary.LittleEndian.PutUint16(entry[0:], uint16(len(val)))
	copy(entry[2:], val)
	l.buf = append(l.buf, entry...)
	l.cnt++
	binary.LittleEndian.PutUint32(l.buf[0:4], uint32(len(l.buf)))
}
func (l *listPack) remove() []byte {
	pos := 4
	for i := 0; i < l.cnt-1; i++ {
		temp := int(binary.LittleEndian.Uint16(l.buf[pos:]))
		pos += 2 + temp
	}
	valLen := int(binary.LittleEndian.Uint16(l.buf[pos:]))
	val := make([]byte, valLen)
	copy(val, l.buf[pos+2:pos+valLen+2])
	l.buf = l.buf[:pos]
	l.cnt--
	binary.LittleEndian.PutUint32(l.buf[0:4], uint32(len(l.buf)))
	return val
}
func (l *listPack) get(idx int) []byte {
	if idx < 0 || idx >= l.cnt {
		return nil
	}
	pos := 4
	for i := 0; i < idx; i++ {
		length := int(binary.LittleEndian.Uint16(l.buf[pos:]))
		pos += 2 + length
	}
	length := int(binary.LittleEndian.Uint16(l.buf[pos:]))
	val := make([]byte, length)
	copy(val, l.buf[pos+2:pos+2+length])
	return val
}
func (m *MyList) RPush(val interface{}) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := val.([]byte)
	if !ok {
		_, ok1 := val.(string)
		if !ok1 {
			return 0
		}
		v = []byte(val.(string))
	}
	if m.data.Len() == 0 {
		m.data.PushBack(newListPack())
	}
	tail := m.data.Back()
	lp := tail.Value.(*listPack)
	if lp.totalBytes()+len(v)+2 > bufSize && lp.cnt > 0 {
		lp = newListPack()
		m.data.PushBack(lp)
	}
	lp.insert(v)
	m.size++
	return 1
}
func (m *MyList) RPop() (val interface{}, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.size == 0 {
		return nil, false
	}
	tail := m.data.Back()
	lp := tail.Value.(*listPack)
	valBytes := lp.remove()
	m.size--
	if lp.cnt == 0 {
		m.data.Remove(tail)
	}
	return valBytes, true
}
func (m *MyList) Index(index int) (val interface{}, exists bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if index < 0 || index >= m.size {
		return nil, false
	}
	var e *list.Element
	var offset int
	if index < m.size/2 {
		e = m.data.Front()
		offset = index
		for e != nil {
			lp := e.Value.(*listPack)
			if offset < lp.cnt {
				return lp.get(offset), true
			}
			offset -= lp.cnt
			e = e.Next()
		}
	} else {
		e = m.data.Back()
		offset = m.size - 1 - index
		for e != nil {
			lp := e.Value.(*listPack)
			if offset < lp.cnt {
				return lp.get(lp.cnt - 1 - offset), true
			}
			offset -= lp.cnt
			e = e.Prev()
		}
	}
	return nil, false
}
func (m *MyList) Range(start, stop int) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if start < 0 || stop >= m.size || start > stop {
		return nil
	}
	limit := stop - start + 1
	res := make([]interface{}, 0, limit)
	e := m.data.Front()
	currentIdx := 0
	for e != nil {
		lp := e.Value.(*listPack)
		if currentIdx+lp.cnt > start {
			break
		}
		currentIdx += lp.cnt
		e = e.Next()
	}
	offset := start - currentIdx
	for e != nil && len(res) < limit {
		lp := e.Value.(*listPack)
		pos := 4
		for i := 0; i < offset; i++ {
			l := int(binary.LittleEndian.Uint16(lp.buf[pos:]))
			pos += 2 + l
		}
		for i := offset; i < lp.cnt && len(res) < limit; i++ {
			l := int(binary.LittleEndian.Uint16(lp.buf[pos:]))
			val := make([]byte, l)
			copy(val, lp.buf[pos+2:pos+2+l])
			res = append(res, val)
			pos += 2 + l
		}
		offset = 0
		e = e.Next()
	}
	return res
}
func (m *MyList) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.size
}
