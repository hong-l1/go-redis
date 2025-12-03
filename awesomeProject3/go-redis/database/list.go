//go:build ignore

package database

import (
	"awesomeProject3/go-redis/datastruct/list"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/utils"
	"awesomeProject3/go-redis/resp/reply"
	"strconv"
)

// getAsList 获取 List 实体，如果 key 存在但不是 List 类型，返回 ErrorReply
func (d *DB) getAsList(key string) (list.List, resp.Reply) {
	entity, exists := d.GetEntity(key)
	if !exists {
		return nil, nil
	}
	// 类型断言
	list, ok := entity.Data.(list.List)
	if !ok {
		return nil, reply.NewWrongTypeErrReply()
	}
	return list, nil
}

// execLPush LPUSH key value [value ...]
func execLPush(db *DB, args [][]byte) resp.Reply {
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
		l.LPush(string(v))
	}
	return reply.NewIntReply()
}

// execRPush RPUSH key value [value ...]
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
		l.RPush(string(v))
	}
	return reply.MakeIntReply(int64(l.Len()))
}

// execLPop LPOP key
func execLPop(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.MakeNullBulkReply()
	}

	val, _ := l.LPop()
	if val == nil {
		return reply.MakeNullBulkReply()
	}

	// 如果 List 空了，移除 Key
	if l.Len() == 0 {
		db.Remove(key)
	}

	return reply.MakeBulkReply(utils.ToCmdLine(val.(string)))
}

// execRPop RPOP key
func execRPop(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.MakeNullBulkReply()
	}

	val, _ := l.RPop()
	if val == nil {
		return reply.MakeNullBulkReply()
	}

	if l.Len() == 0 {
		db.Remove(key)
	}

	return reply.MakeBulkReply(utils.ToCmdLine(val.(string)))
}

// execLLen LLEN key
func execLLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.MakeIntReply(0)
	}

	return reply.MakeIntReply(int64(l.Len()))
}

// execLIndex LINDEX key index
func execLIndex(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	indexStr := string(args[1])
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}

	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.MakeNullBulkReply()
	}

	// 处理负数索引
	if index < 0 {
		index = l.Len() + index
	}

	val, _ := l.Index(index)
	if val == nil {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(utils.ToCmdLine(val.(string)))
}

// execLRange LRANGE key start stop
func execLRange(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	startStr := string(args[1])
	stopStr := string(args[2])

	start, err := strconv.Atoi(startStr)
	if err != nil {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(stopStr)
	if err != nil {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}

	l, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if l == nil {
		return reply.MakeEmptyMultiBulkReply()
	}

	// MyList 的 Range 应该处理好边界，但通常 DB 层也会处理
	// 这里假设底层 Range 实现完善
	vals := l.Range(start, stop)
	result := make([][]byte, len(vals))
	for i, v := range vals {
		result[i] = utils.ToCmdLine(v.(string))[0]
	}
	return reply.MakeMultiBulkReply(result)
}

// init 注册命令
func init() {
	RegisterCommand("LPush", execLPush, -3)
	RegisterCommand("RPush", execRPush, -3)
	RegisterCommand("LPop", execLPop, 2)
	RegisterCommand("RPop", execRPop, 2)
	RegisterCommand("LLen", execLLen, 2)
	RegisterCommand("LIndex", execLIndex, 3)
	RegisterCommand("LRange", execLRange, 4)
}
