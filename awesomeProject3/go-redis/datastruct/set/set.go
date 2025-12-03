package set

type Set interface {
	Add(val string) int
	Remove(val string) int
	Has(val string) bool
	Len() int
	Members() []string
	Intersect(another Set) []string
	Union(another Set) []string
	Diff(another Set) []string
	ForEach(consumer func(member string) bool)
	RandomMembers(limit int) []string
	RandomDistinctMembers(limit int) []string
}
