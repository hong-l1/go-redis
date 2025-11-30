package reply

type UnknownErrReply struct {
}

var unknownErrReply = []byte("-unknown error\r\n")

func NewUnKnownErrReply() *UnknownErrReply {
	return &UnknownErrReply{}
}
func (u *UnknownErrReply) Error() string {
	return "unknown error"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unknownErrReply
}

type ArgErrReply struct {
	cmd string
}

func NewArgErrReply(cmd string) *ArgErrReply {
	return &ArgErrReply{cmd}
}
func (a *ArgErrReply) Error() string {
	return "err wrong number of arguments for '" + a.cmd + "'\r\n"
}

func (a *ArgErrReply) ToBytes() []byte {
	return []byte("- err wrong number of arguments for '" + a.cmd + "'\r\n")
}

type SyntaxErrReply struct{}

var syntaxErrReply = []byte("-syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "syntax error"
}
func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrReply
}

type WrongTypeErrReply struct {
}

var wrongTypeErrReply = []byte("- WRONGTYPE Operation against a key holding the wrong kind of value \r\n")

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrReply
}

type ProtocolErrReply struct {
	msg string
}

func (p *ProtocolErrReply) Error() string {
	return "protocol error'" + p.msg + "'\r\n"
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-protocol error'" + p.msg + "'\r\n")
}
