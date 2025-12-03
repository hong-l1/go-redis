package set

import "awesomeProject3/go-redis/datastruct/dict"

type MySet struct {
	dict dict.Dict
}

func NewMySet() *MySet {
	return &MySet{
		dict: dict.NewDict(),
	}
}
func (m *MySet) Add(val string) int {
	return m.dict.Put(val, struct{}{})
}

func (m *MySet) Remove(val string) int {
	_, result := m.dict.Remove(val)
	return result
}

func (m *MySet) Has(val string) bool {
	_, ok := m.dict.Get(val)
	return ok
}

func (m *MySet) Len() int {
	return m.dict.Len()
}

func (m *MySet) Members() []string {
	return m.dict.Keys()
}
func (m *MySet) ForEach(consumer func(member string) bool) {
	m.dict.ForEach(func(key string, val any) bool {
		return consumer(key)
	})
}
func (m *MySet) Intersect(another Set) []string {
	if another == nil {
		return nil
	}
	result := make([]string, 0)
	m.ForEach(func(member string) bool {
		if another.Has(member) {
			result = append(result, member)
		}
		return true
	})
	return result
}

func (m *MySet) Union(another Set) []string {
	if another == nil {
		return m.Members()
	}
	uniqueMap := make(map[string]struct{})
	m.ForEach(func(member string) bool {
		uniqueMap[member] = struct{}{}
		return true
	})
	another.ForEach(func(member string) bool {
		uniqueMap[member] = struct{}{}
		return true
	})
	result := make([]string, 0, len(uniqueMap))
	for member := range uniqueMap {
		result = append(result, member)
	}
	return result
}

func (m *MySet) Diff(another Set) []string {
	if another == nil {
		return m.Members()
	}
	result := make([]string, 0)
	m.ForEach(func(member string) bool {
		if !another.Has(member) {
			result = append(result, member)
		}
		return true
	})
	return result
}

func (m *MySet) RandomMembers(limit int) []string {
	return m.dict.RandomKeys(limit)
}

func (m *MySet) RandomDistinctMembers(limit int) []string {
	return m.dict.RandomDistinctKeys(limit)
}
