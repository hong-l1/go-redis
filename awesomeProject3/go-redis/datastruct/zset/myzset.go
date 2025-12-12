package zset

func (m *MyZSet) Add(member string, score float64) bool {
	node := m.sl.insert(member, score)
	if node == nil {
		return false
	}
	return true
}
func (m *MyZSet) Len() int64 {
	length := m.sl.length
	return int64(length)
}
func (m *MyZSet) Get(member string) (Element, bool) {
	val, ok := m.dict[member]
	if !ok {
		return Element{}, false
	}
	return Element{
		Member: val.data.Member,
		Score:  val.data.Score,
	}, true
}
func (m *MyZSet) Remove(member string) bool {
	node := m.dict[member]
	return m.sl.delete(member, node.data.Score)
}
func (m *MyZSet) GetRank(member string) int64 {
	node := m.dict[member]
	rank, ok := m.sl.getRank(member, node.data.Score)
	if ok {
		return int64(rank)
	}
	return -1
}
func (m *MyZSet) GetRevRank(member string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	length := m.Len()
	m.
}

func (m *MyZSet) Range(start int64, stop int64) []Element {
	//TODO implement me
	panic("implement me")
}

func (m *MyZSet) RangeByScore(min *ScoreBorder, max *ScoreBorder, offset int64, limit int64, withScores bool) []Element {
	//TODO implement me
	panic("implement me")
}

func (m *MyZSet) RemoveRangeByScore(min *ScoreBorder, max *ScoreBorder) int64 {
	//TODO implement me
	panic("implement me")
}

func (m *MyZSet) RemoveRangeByRank(start int64, stop int64) int64 {
	//TODO implement me
	panic("implement me")
}
