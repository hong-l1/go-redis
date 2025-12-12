package hash

type Hash interface {
	HSet(field string, value []byte) (result int)
	HGet(field string) (value []byte, exists bool)
	HDel(field string) (result int)
	HExists(field string) bool
	HLen() int
	HKeys() []string
	HValues() [][]byte
	HGetAll() [][]byte
	ForEach(consumer func(field string, value []byte) bool)
}
