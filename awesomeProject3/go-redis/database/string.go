package database

import (
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
)

func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	data, ok := db.GetEntity(key)
	if !ok {
		return reply.NewNilBulkReply()
	}
	bytes := data.Data.([]byte)
	return reply.NewBulkReply(bytes)
}
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &database.DataEntity{Data: args[1]}
	_ = db.PutEntity(key, entity)
	return reply.NewOkReply()
}
func execSetNx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &database.DataEntity{Data: args[1]}
	result := db.PutIfAbsent(key, entity)
	return reply.NewIntReply(result)
}
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &database.DataEntity{Data: args[1]}
	data, ok := db.GetEntity(key)
	if !ok {
		return reply.NewNilBulkReply()
	}
	db.PutEntity(key, entity)
	return reply.NewBulkReply(data.Data.([]byte))

}
func execSreLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	data, ok := db.GetEntity(key)
	if !ok {
		return reply.NewNilBulkReply()
	}
	bytes := data.Data.([]byte)
	return reply.NewIntReply(len(bytes))
}
func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("setnx", execSetNx, 3)
	RegisterCommand("getset", execGetSet, 3)
	RegisterCommand("strlen", execSreLen, 2)
}
