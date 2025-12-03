package database

import (
	"awesomeProject3/go-redis/aof"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"strconv"
	"strings"
)

var DefaultDBSize = 16

type Database struct {
	dbSet      []*DB
	aofHandler aof.AofHandler
}

func NewDatabase() *Database {
	database := new(Database)
	aofHandler, _ := aof.NewAofHandler(database)
	dbset := make([]*DB, DefaultDBSize)
	for k := range dbset {
		db := NewDB()
		db.addAof = func(cmd [][]byte) {
			aofHandler.AddAof(k, cmd)
		}
		db.Index = k
		dbset[k] = db
	}
	return &Database{dbSet: dbset}
}
func execSelect(conn resp.Connection, database *Database, args [][]byte, ) resp.Reply {
	num, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("invalid db index")
	}
	if num < 0 || num >= len(database.dbSet) {
		return reply.NewErrReply(" db index is out of range")
	}
	conn.SelectDB(num)
	return reply.NewOkReply()
}
func (d *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	if args == nil {
		return reply.NewErrReply("empty command")
	}
	cmd := strings.ToLower(string(args[0]))
	if cmd == "select" {
		if len(args) != 2 {
			return reply.NewArgErrReply("invalid args")
		}
		return execSelect(client, d, args[1:])
	}
	db := client.GetDBIndex()
	return d.dbSet[db].Exec(client, args)
}

func (d *Database) AfterClientClose(c resp.Connection) {
}

func (d *Database) Close() {
}
