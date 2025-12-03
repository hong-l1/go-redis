package connection

import (
	"awesomeProject3/go-redis/pkg/wait"
	"net"
	"sync"
	"time"
)

type Connect struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	RedisDB      int
}

func NewConnect(conn net.Conn) *Connect {
	return &Connect{
		conn: conn,
	}
}
func (c *Connect) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
func (c *Connect) GetDBIndex() int {
	return c.RedisDB
}

func (c *Connect) SelectDB(i int) {
	c.RedisDB = i
}

func (c *Connect) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.waitingReply.Add(1)
	defer c.waitingReply.Done()
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Connect) Close() error {
	c.waitingReply.WaitWithTimeout(time.Second * 3)
	c.conn.Close()
	return nil
}
