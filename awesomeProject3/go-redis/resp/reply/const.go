package reply

type PongReply struct {
}

var pongReply = []byte("+pong\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongReply
}
func NewPongReply() *PongReply {
	return &PongReply{}
}

type OkReply struct {
}

var okReply = []byte("+OK\r\n")

func (p *OkReply) ToBytes() []byte {
	return okReply
}
func NewOkReply() *OkReply {
	return &OkReply{}
}

type NilBulkReply struct {
}

func NewNilBulkReply() *NilBulkReply {
	return &NilBulkReply{}
}

var nilBulkReply = []byte("$-1\r\n")

func (p *NilBulkReply) ToBytes() []byte {
	return nilBulkReply
}

type EmptyMultiBulkReply struct {
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

var emptyMultiBulkReply = []byte("*0\r\n")

func (p *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkReply
}

type NOBytes struct {
}

var noBytes = []byte("")

func (p *NOBytes) ToBytes() []byte {
	return noBytes
}
