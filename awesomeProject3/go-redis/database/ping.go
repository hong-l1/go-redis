package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.NewPongReply()
}
func init() {
	RegisterCommand("ping", Ping, 1)
}
