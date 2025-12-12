package database

import (
	"awesomeProject3/go-redis/datastruct/list"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"strconv"
)

func (d *DB) getAsList(key string) (list.List, resp.Reply) {
	entity, exists := d.GetEntity(key)
	if !exists {
		return nil, nil
	}
	list, ok := entity.Data.(list.List)
	if !ok {
		return nil, reply.NewErrReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return list, nil
}
func execRPush(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	values := args[1:]
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		l = list.NewMyList()
		db.PutEntity(key, &database.DataEntity{Data: l})
	}
	for _, v := range values {
		l.RPush(v)
	}
	return reply.NewIntReply(l.Len())
}
func execRPop(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.NewNilBulkReply()
	}
	val, exists := l.RPop()
	if !exists {
		return reply.NewNilBulkReply()
	}
	return reply.NewBulkReply(val.([]byte))
}
func execLLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.NewIntReply(0)
	}
	return reply.NewIntReply(l.Len())
}
func execLIndex(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	indexStr := string(args[1])
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.NewNilBulkReply()
	}
	val, exists := l.Index(index)
	if !exists {
		return reply.NewNilBulkReply()
	}
	return reply.NewBulkReply(val.([]byte))
}
func execLRange(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	startStr := string(args[1])
	stopStr := string(args[2])
	start, err1 := strconv.Atoi(startStr)
	stop, err2 := strconv.Atoi(stopStr)
	if err1 != nil || err2 != nil {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}
	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.NewEmptyMultiBulkReply()
	}
	vals := l.Range(start, stop)
	result := make([][]byte, len(vals))
	for i, v := range vals {
		result[i] = v.([]byte)
	}
	return reply.NewMultiBulkReply(result)
}
func init() {
	RegisterCommand("rpush", execRPush, -3)
	RegisterCommand("rpop", execRPop, 2)
	RegisterCommand("llen", execLLen, 2)
	RegisterCommand("lindex", execLIndex, 3)
	RegisterCommand("lrange", execLRange, 4)
}
