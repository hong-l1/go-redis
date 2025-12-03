package list

type List interface {
	LPop() (val interface{}, exists bool)
	RPop() (val interface{}, exists bool)
	LPush(val interface{}) (result int)
	RPush(val interface{}) (result int)
	Index(index int) (val interface{}, exists bool)
	Insert(index int, val interface{}) (result int)
	Range(start, stop int) []interface{}
	Trim(start, stop int)
}
