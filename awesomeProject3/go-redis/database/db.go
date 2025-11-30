package database

import (
	"awesomeProject3/go-redis/datastruct/dict"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"strings"
)

type DB struct {
	Index int
	dict  dict.Dict
}
type ExecFunc func(db *DB, args [][]byte) resp.Reply

func NewDB() *DB {
	return &DB{
		dict: dict.NewDict(),
	}
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
	val, exists := d.dict.Get(key)
	if !exists {
		return nil, exists
	}
	return val.(*database.DataEntity), exists
}
func (d *DB) PutEntity(key string, val *database.DataEntity) int {
	return d.dict.Put(key, val)
}
func (d *DB) PutIfExists(key string, val *database.DataEntity) int {
	return d.dict.PutIfExists(key, val)
}
func (d *DB) PutIfAbsent(key string, val *database.DataEntity) int {
	return d.dict.PutIfAbsent(key, val)
}
func (d *DB) Remove(key string) {
	d.dict.Remove(key)
}
func (d *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		if _, ok := d.dict.Get(key); ok {
			d.dict.Remove(key)
			deleted++
		}
	}
	return deleted
}
func (d *DB) Flush() {
	d.dict.Clear()
}
