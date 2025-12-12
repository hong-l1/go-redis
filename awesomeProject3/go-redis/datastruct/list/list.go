package list

type List interface {
	RPop() (val interface{}, exists bool)
	RPush(val interface{}) (result int)
	Index(index int) (val interface{}, exists bool)
	Range(start, stop int) []interface{}
	Len() int
}
