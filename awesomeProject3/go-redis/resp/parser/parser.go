package parser

import (
	"awesomeProject3/go-redis/ineterface/resp"
	reply2 "awesomeProject3/go-redis/resp/reply"
	"bufio"
	"errors"
	"io"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err  error
}
type readState struct {
	readingMultiLine bool
	expectedArgsCnt  int
	msgType          byte
	args             [][]byte
	bulkLen          int
}

func (r *readState) IsFinished() bool {
	return r.expectedArgsCnt > 0 && r.expectedArgsCnt == len(r.args)
}

// ParseStream 给redis内核返回一个channel,内核可以不断从channel中拿到解析器解析到的数据
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

// for 循环一直执行,不断使用readline来读取客户端发来的指令
// readline每次读取一行；buldlen=0就是初始化的；buildlen就是需要读取多次。
// pase...和readBody都是通过修改readState的参数来控制readline读取的方式
// parseSingleLineReply 简单命令可以直接返回解析后的语句给redis内核
func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(debug.Stack())
		}
	}()
	var msg []byte
	var err error
	var state readState
	bufReader := bufio.NewReader(reader)
	for {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				ch <- &Payload{Err: err}
				close(ch)
				return
			} else {
				ch <- &Payload{Err: err}
				state = readState{}
				continue
			}
		}
		if !state.readingMultiLine {
			if msg[0] == '*' {
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{Err: errors.New("protocol error: " + string(msg))}
					state = readState{}
					continue
				}
				if state.expectedArgsCnt == 0 {
					ch <- &Payload{Data: &reply2.EmptyMultiBulkReply{}}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{Err: errors.New("protocol error: " + string(msg))}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &Payload{Data: &reply2.NilBulkReply{}}
					state = readState{}
					continue
				}
				continue
			} else {
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
				continue
			}
		} else {
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{Err: errors.New("protocol error: " + string(msg))}
			}
			if state.IsFinished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply2.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply2.NewBulkReply(state.args[0])
				}
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
			}
		}
	}
}
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol err" + string(msg))
		}
	} else {
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol err" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// 解析数组，设置state通过state来
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	state.expectedArgsCnt, err = strconv.Atoi(string(msg[1 : len(msg)-2]))
	if err != nil {
		return errors.New("protocol err" + string(msg))
	}
	if state.expectedArgsCnt == 0 {
		state.expectedArgsCnt = 0
		return nil
	} else if state.expectedArgsCnt > 0 {
		state.readingMultiLine = true
		state.msgType = msg[0]
		state.args = make([][]byte, 0, state.expectedArgsCnt)
		return nil
	} else {
		return errors.New("protocol err" + string(msg))
	}
}

// 解析字符串
func parseBulkHeader(msg []byte, state *readState) error {
	bulkLen, err := strconv.Atoi(string(msg[1 : len(msg)-2]))
	if err != nil {
		return errors.New("protocol err" + string(msg))
	}
	if bulkLen == -1 {
		state.bulkLen = -1
		return nil
	}
	// bulkLen >= 0: prepare to read body as a single bulk string
	state.bulkLen = bulkLen
	state.readingMultiLine = true
	state.msgType = msg[0]
	state.expectedArgsCnt = 1
	state.args = make([][]byte, 0, 1)
	return nil
}
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply2.NewStatusReply(str[1:])
	case '-':
		result = reply2.NewErrReply(str[1:])
	case ':':
		num, err := strconv.Atoi(str[1:])
		if err != nil {
			return nil, errors.New("protocol err" + string(msg))
		}
		result = reply2.NewIntReply(num)
	}
	return result, nil
}
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		state.bulkLen, err = strconv.Atoi(string(line[1:]))
		if err != nil {
			return errors.New("protocol err" + string(line))
		}
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
