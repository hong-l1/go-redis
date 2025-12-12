package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
)

func NewRouter() map[string]cmdline {
	routerMap := make(map[string]cmdline)
	routerMap["get"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["ping"] = pingFunc
	routerMap["rename"] = renameFunc
	routerMap["renamenx"] = renameFunc
	routerMap["flushdb"] = flushDBFunc
	routerMap["select"] = selectFunc
	return routerMap
}
func defaultFunc(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply {
	key := string(args[1])
	peer := c.peerPicker.PickNode(key)
	return reply.NewBulkReply([]byte(peer))
}
