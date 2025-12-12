package database

import (
	"awesomeProject3/go-redis/datastruct/hash"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/utils"
	"awesomeProject3/go-redis/resp/reply"
)

func (d *DB) GetAsHash(key string) (hash.Hash, resp.Reply) {
	entity, ok := d.GetEntity(key)
	if !ok {
		return nil, nil
	}
	data, ok := entity.Data.(hash.Hash)
	if !ok {
		return nil, reply.NewErrReply("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return data, nil
}
func execHSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	value := args[2]
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		h = hash.NewMyHash()
		db.PutEntity(key, &database.DataEntity{Data: h})
	}
	result := h.HSet(field, value)
	return reply.NewIntReply(result)
}
func execHGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewNilBulkReply()
	}
	val, ok := h.HGet(field)
	if !ok {
		return reply.NewNilBulkReply()
	}
	return reply.NewBulkReply(val)
}
func execHDel(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	fields := args[1:]
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewStatusReply("key does not exist")
	}
	deleted := 0
	for _, filed := range fields {
		deleted += h.HDel(string(filed))
	}
	if h.HLen() == 0 {
		db.Remove(key)
	}
	return reply.NewIntReply(deleted)
}
func execHGetAll(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewEmptyMultiBulkReply()
	}
	data := h.HGetAll()
	return reply.NewMultiBulkReply(data)
}
func execHExists(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewIntReply(0)
	}
	ok := h.HExists(field)
	if ok {
		return reply.NewIntReply(1)
	}
	return reply.NewIntReply(0)
}
func execHKeys(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewEmptyMultiBulkReply()
	}
	res := h.HKeys()
	return reply.NewMultiBulkReply(utils.ToCmdLine(res...))
}
func execHLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewIntReply(0)
	}
	return reply.NewIntReply(h.HLen())
}
func execHValues(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		return reply.NewEmptyMultiBulkReply()
	}
	data := h.HValues()
	return reply.NewMultiBulkReply(data)
}
func execHMGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	fields := args[1:]
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		// If hash doesn't exist, return array of nils
		ans := make([][]byte, len(fields))
		return reply.NewMultiBulkReply(ans)
	}
	ans := make([][]byte, 0, len(fields))
	for _, field := range fields {
		val, ok := h.HGet(string(field))
		if !ok {
			ans = append(ans, nil)
		} else {
			ans = append(ans, val)
		}
	}
	return reply.NewMultiBulkReply(ans)
}
func execHMSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	data := args[1:]
	if len(data)&1 != 0 {
		return reply.NewErrReply("ERR wrong number of arguments for 'hmset' command")
	}
	h, err := db.GetAsHash(key)
	if err != nil {
		return err
	}
	if h == nil {
		h = hash.NewMyHash()
		db.PutEntity(key, &database.DataEntity{Data: h})
	}
	result := 0
	for i := 0; i < len(data); i += 2 {
		result += h.HSet(string(data[i]), data[i+1])
	}
	return reply.NewIntReply(result)
}
func init() {
	RegisterCommand("hset", execHSet, 4)
	RegisterCommand("hget", execHGet, 3)
	RegisterCommand("hdel", execHDel, -3)
	RegisterCommand("hgetall", execHGetAll, 2)
	RegisterCommand("hexists", execHExists, 3)
	RegisterCommand("hkeys", execHKeys, 2)
	RegisterCommand("hvalues", execHValues, 2)
	RegisterCommand("hlen", execHLen, 2)
	RegisterCommand("hmget", execHMGet, -3)
	RegisterCommand("hmset", execHMSet, -3)
}
