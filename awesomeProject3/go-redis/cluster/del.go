package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
)

func delFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	replys := c.broadcast(connect, args)
	del := 0
	for _, r := range replys {
		if reply.IsErrorReply(r) {
			continue
		}
		delnum, ok := r.(*reply.IntReply)
		if !ok {
			return reply.NewErrReply("wrong reply type")
		}
		del += delnum.Num
	}
	return reply.NewIntReply(del)
}
