package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
)

func flushDBFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	replys := c.broadcast(connect, args)
	for _, r := range replys {
		if reply.IsErrorReply(r) {
			return r
		}
	}
	return reply.NewOkReply()
}
