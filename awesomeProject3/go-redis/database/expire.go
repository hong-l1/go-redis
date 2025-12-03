package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"fmt"
	"strconv"
	"time"
)

func execExpire(db *DB, args [][]byte) resp.Reply {
	fmt.Println(len(args))
	if len(args) != 2 {
		return reply.NewArgErrReply("expire")
	}
	key := string(args[0])
	secondsArg, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.NewErrReply("ERR value is not an integer or out of range")
	}
	if secondsArg <= 0 {
		return reply.NewErrReply("ERR invalid expire time in 'expire' command")
	}
	_, exists := db.GetEntity(key)
	if !exists {
		return reply.NewIntReply(0)
	}
	expireAt := time.Now().Add(time.Duration(secondsArg) * time.Second)
	db.Expire(key, expireAt)
	//db.addAof(aof.MakeExpireCmd(key, expireAt).Args)
	return reply.NewIntReply(1)
}
func init() {
	RegisterCommand("expire", execExpire, 3)
}
