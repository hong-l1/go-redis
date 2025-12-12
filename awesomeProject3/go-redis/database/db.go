package database

import (
	"awesomeProject3/go-redis/datastruct/dict"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/timewheel"
	"awesomeProject3/go-redis/resp/reply"
	"log"
	"strings"
	"time"
)

type DB struct {
	Index  int
	dict   dict.Dict
	ttl    dict.Dict
	addAof func(cmd [][]byte)
}
type ExecFunc func(db *DB, args [][]byte) resp.Reply

func NewDB() *DB {
	db := &DB{
		dict: dict.NewDict(),
		ttl:  dict.NewDict(),
		addAof: func(cmd [][]byte) {
		},
	}
	return db
}
func (d *DB) Exec(conn resp.Connection, args [][]byte) resp.Reply {
	name := strings.ToLower(string(args[0]))
	command, ok := CmdMap[name]
	if !ok {
		return reply.NewErrReply("invalid command" + name)
	}
	if !validateArity(command.arity, args) {
		return reply.NewArgErrReply(name)
	}
	return command.exec(d, args[1:])
}
func validateArity(arity int, args [][]byte) bool {
	argNum := len(args)
	if arity > 0 {
		return arity == argNum
	}
	return argNum >= -arity
}
func (d *DB) GetEntity(key string) (*database.DataEntity, bool) {
	if d.IsExpired(key) {
		return nil, false
	}
	val, exists := d.dict.Get(key)
	if !exists {
		return nil, exists
	}
	return val.(*database.DataEntity), exists
}
func (d *DB) PutEntity(key string, val *database.DataEntity) int {
	d.IsExpired(key)
	return d.dict.Put(key, val)
}
func (d *DB) PutIfExists(key string, val *database.DataEntity) int {
	d.IsExpired(key)
	return d.dict.PutIfExists(key, val)
}
func (d *DB) PutIfAbsent(key string, val *database.DataEntity) int {
	d.IsExpired(key)
	return d.dict.PutIfAbsent(key, val)
}
func (d *DB) Remove(key string) {
	d.dict.Remove(key)
	d.ttl.Remove(key)
	timewheel.Cancel(genExpireTask(key))
}
func (d *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		if _, ok := d.dict.Get(key); ok {
			d.Remove(key)
			deleted++
		}
	}
	return deleted
}
func (d *DB) Flush() {
	d.dict.Clear()
}
func genExpireTask(key string) string {
	return "expire:" + key
}
func (d *DB) Expire(key string, expireTime time.Time) {
	d.ttl.Put(key, expireTime)
	taskKey := genExpireTask(key)
	timewheel.At(expireTime, taskKey, func() {
		log.Println("expire " + key)
		rawExpireTime, ok := d.ttl.Get(key)
		if !ok {
			return
		}
		expireTime, _ := rawExpireTime.(time.Time)
		expired := time.Now().After(expireTime)
		if expired {
			d.Remove(key)
		}
	})
}
func (d *DB) Persist(key string) {
	d.ttl.Remove(key)
	taskKey := genExpireTask(key)
	timewheel.Cancel(taskKey)
}
func (d *DB) IsExpired(key string) bool {
	rawExpireTime, ok := d.ttl.Get(key)
	if !ok {
		return false
	}
	expireTime, _ := rawExpireTime.(time.Time)
	expired := time.Now().After(expireTime)
	if expired {
		d.Remove(key)
	}
	return expired
}
