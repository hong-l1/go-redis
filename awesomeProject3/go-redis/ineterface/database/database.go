package database

import "awesomeProject3/go-redis/ineterface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	AfterClientClose(c resp.Connection)
	Close()
}
type DataEntity struct {
	Data any
}
