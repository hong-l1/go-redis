package client

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/wait"
	"awesomeProject3/go-redis/resp/parser"
	"awesomeProject3/go-redis/resp/reply"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	chanSize = 256
	maxWait  = 3 * time.Second
)
const (
	created = iota
	running
	closed
)

type Client struct {
	conn        net.Conn
	sendingReqs chan *request
	waitingReqs chan *request
	ticker      *time.Ticker
	addr        string
	status      int32
	working     *sync.WaitGroup
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr:        addr,
		conn:        conn,
		status:      created,
		sendingReqs: make(chan *request, chanSize),
		waitingReqs: make(chan *request, chanSize),
		working:     &sync.WaitGroup{},
	}, nil
}

type request struct {
	args      [][]byte
	reply     resp.Reply
	heartbeat bool
	waiting   *wait.Wait
	err       error
}

func (client *Client) Start() {
	client.ticker = time.NewTicker(10 * time.Second)
	atomic.StoreInt32(&client.status, running)
	go client.handleWrite()
	go func() {
		client.handleRead()
	}()
	go client.heartbeat()
}
func (client *Client) Close() {
	atomic.StoreInt32(&client.status, closed)
	client.ticker.Stop()
	close(client.sendingReqs)
	client.working.Wait()
	_ = client.conn.Close()
	close(client.waitingReqs)
}
func (client *Client) Send(args [][]byte) resp.Reply {
	if atomic.LoadInt32(&client.status) != running {
		return reply.NewErrReply("client closed")
	}
	res := &request{
		args:      args,
		heartbeat: false,
		waiting:   &wait.Wait{},
	}
	res.waiting.Add(1)
	client.working.Add(1)
	defer client.working.Done()
	client.sendingReqs <- res
	timeout := res.waiting.WaitWithTimeout(maxWait)
	if timeout {
		return reply.NewErrReply("server time out")
	}
	if res.err != nil {
		return reply.NewErrReply("request failed: " + res.err.Error())
	}
	return res.reply
}
func (client *Client) handleWrite() {
	for req := range client.sendingReqs {
		client.doRequest(req)
	}
}
func (client *Client) doRequest(req *request) {
	if req == nil || len(req.args) == 0 {
		return
	}
	re := reply.NewMultiBulkReply(req.args)
	bytes := re.ToBytes()
	_, err := client.conn.Write(bytes)
	i := 0
	for err != nil && i < 3 {
		_, err = client.conn.Write(bytes)
		if err == nil {
			break
		}
		i++
	}
	if err == nil {
		client.waitingReqs <- req
	} else {
		req.err = err
		req.waiting.Done()
	}
}
func (client *Client) finishRequest(reply resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			log.Println(err)
		}
	}()
	res := <-client.waitingReqs
	if res == nil {
		return
	}
	res.reply = reply
	if res.waiting != nil {
		res.waiting.Done()
	}
}
func (client *Client) handleRead() {
	ch := parser.ParseStream(client.conn)
	for payload := range ch {
		if payload.Err != nil {
			client.finishRequest(reply.NewErrReply(payload.Err.Error()))
			continue
		}
		client.finishRequest(payload.Data)
	}
}
func (client *Client) doHeartbeat() {
	res := &request{
		args:      [][]byte{[]byte("PING")},
		heartbeat: true,
		waiting:   &wait.Wait{},
	}
	res.waiting.Add(1)
	client.working.Add(1)
	defer client.working.Done()
	client.sendingReqs <- res
	res.waiting.WaitWithTimeout(maxWait)
}
func (client *Client) heartbeat() {
	for range client.ticker.C {
		client.doHeartbeat()
	}
}
func (client *Client) RemoteAddress() string {
	return client.addr
}
