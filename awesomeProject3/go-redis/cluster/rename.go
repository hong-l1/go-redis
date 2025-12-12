package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/utils"
)

func renameFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	key1 := string(args[1])
	key2 := string(args[2])
	src := c.peerPicker.PickNode(key1)
	dst := c.peerPicker.PickNode(key2)
	if src != dst {
		utils.ToCmdLine("")
	}
	return c.relay(src, connect, args)
}
