package database

import "awesomeProject3/go-redis/ineterface/resp"

func pingFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	return c.relay(c.self, connect, args)
}
