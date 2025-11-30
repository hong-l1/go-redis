package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/wildcard"
	"awesomeProject3/go-redis/resp/reply"
)

func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, key := range args {
		keys[i] = string(key)
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(deleted)
}
func execExists(db *DB, args [][]byte) resp.Reply {
	num := 0
	for _, key := range args {
		_, ok := db.GetEntity(string(key))
		if ok {
			num++
		}
	}
	return reply.NewIntReply(num)
}
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.NewOkReply()
}
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	data, ok := db.GetEntity(key)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	switch data.Data.(type) {
	case []byte:
		return reply.NewBulkReply([]byte("string"))
	}
	return reply.NewUnKnownErrReply()
}
func execRename(db *DB, args [][]byte) resp.Reply {
	src := args[0]
	dst := args[1]
	data, ok := db.GetEntity(string(src))
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(string(dst), data)
	db.Remove(string(src))
	return reply.NewOkReply()
}
func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := args[0]
	dst := args[1]
	_, ok := db.GetEntity(string(dst))
	if ok {
		return reply.NewIntReply(0)
	}
	data, ok := db.GetEntity(string(src))
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(string(dst), data)
	db.Remove(string(src))
	return reply.NewIntReply(1)
}
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.dict.ForEach(func(k string, v interface{}) bool {
		if pattern.IsMatch(k) {
			result = append(result, []byte(k))
		}
		return true
	})
	return reply.NewMultiBulkReply(result)
}
func init() {
	RegisterCommand("delete", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushDB", execFlushDB, -1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renameNx", execRenameNx, 3)
	RegisterCommand("keys", execKeys, 2)
}
