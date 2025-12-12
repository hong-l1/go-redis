package zset

type ScoreBorder struct {
	Value   float64
	Exclude bool
}
type ZSet interface {
	Add(member string, score float64) bool
	Len() int64
	Get(member string) (Element, bool)
	Remove(member string) bool
	GetRank(member string) int64
	GetRevRank(member string) int64
	Range(start int64, stop int64) []Element
	RangeByScore(min *ScoreBorder, max *ScoreBorder, offset int64, limit int64, withScores bool) []Element
	RemoveRangeByScore(min *ScoreBorder, max *ScoreBorder) int64
	RemoveRangeByRank(start int64, stop int64) int64
}
