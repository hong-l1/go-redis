package reply

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"bytes"
	"strconv"
)

const (
	CRLF = "\r\n"
)

type BulkReply struct {
	Content []byte
}

func NewBulkReply(content []byte) *BulkReply {
	return &BulkReply{
		content,
	}
}
func (p *BulkReply) ToBytes() []byte {
	if p.Content == nil {
		return []byte("$-1" + CRLF)
	}
	return []byte("$" + strconv.Itoa(len(p.Content)) + CRLF + string(p.Content) + CRLF)
}

type MultiBulkReply struct {
	Content [][]byte
}

func NewMultiBulkReply(content [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		content,
	}
}
func (p *MultiBulkReply) ToBytes() []byte {
	n := len(p.Content)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(n) + CRLF)
	for _, v := range p.Content {
		if v == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(v)) + CRLF + string(v) + CRLF)
		}
	}
	return buf.Bytes()
}

type StatusReply struct {
	Status string
}

func NewStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

type IntReply struct {
	Num int
}

func (p *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.Itoa(p.Num) + CRLF)
}
func NewIntReply(num int) *IntReply {
	return &IntReply{Num: num}
}

type StandardErrReply struct {
	status string
}

func NewErrReply(status string) *StandardErrReply {
	return &StandardErrReply{status: status}
}
func (p *StandardErrReply) ToBytes() []byte {
	return []byte("-" + p.status + CRLF)
}
func (p *StandardErrReply) Error() string {
	return p.status
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
