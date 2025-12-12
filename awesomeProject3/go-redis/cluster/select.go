package database

import "awesomeProject3/go-redis/ineterface/resp"

func selectFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	return c.db.Exec(connect, args)
}
