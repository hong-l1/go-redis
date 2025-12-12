package hash

import "awesomeProject3/go-redis/datastruct/dict"

type MyHash struct {
	dict *dict.MyDict
}

func (m *MyHash) HSet(field string, value []byte) (result int) {
	return m.dict.Put(field, value)
}
func (m *MyHash) HGet(field string) (value []byte, exists bool) {
	val, ok := m.dict.Get(field)
	if !ok {
		return nil, ok
	}
	return val.([]byte), ok
}

func (m *MyHash) HDel(field string) (result int) {
	_, result = m.dict.Remove(field)
	return result
}

func (m *MyHash) HExists(field string) bool {
	_, ok := m.dict.Get(field)
	return ok
}

func (m *MyHash) HLen() int {
	return m.dict.Len()
}

func (m *MyHash) HKeys() []string {
	return m.dict.Keys()
}

func (m *MyHash) HValues() [][]byte {
	res := make([][]byte, 0, m.HLen())
	m.ForEach(func(field string, value []byte) bool {
		res = append(res, value)
		return true
	})
	return res
}

func (m *MyHash) HGetAll() [][]byte {
	res := make([][]byte, 0, m.HLen()*2)
	m.ForEach(func(field string, value []byte) bool {
		res = append(res, []byte(field))
		res = append(res, value)
		return true
	})
	return res
}
func (m *MyHash) ForEach(consumer func(field string, value []byte) bool) {
	m.dict.ForEach(func(key string, val interface{}) bool {
		return consumer(key, val.([]byte))
	})
}

func NewMyHash() *MyHash {
	return &MyHash{
		dict: dict.NewDict(),
	}
}
