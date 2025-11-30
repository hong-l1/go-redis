package resp

type Connection interface {
	GetDBIndex() int
	SelectDB(int)
	Write([]byte) error
}
