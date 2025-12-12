package database

import (
	"awesomeProject3/go-redis/aof"
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/resp/reply"
	"log"
	"strconv"
	"strings"
)

var DefaultDBSize = 16

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{}
	dbset := make([]*DB, DefaultDBSize)
	for i := range dbset {
		db := NewDB()
		db.Index = i
		dbset[i] = db
	}
	database.dbSet = dbset
	aofHandler, _ := aof.NewAofHandler(database)
	database.aofHandler = aofHandler
	for i, db := range database.dbSet {
		localIndex := i
		db.addAof = func(cmd [][]byte) {
			database.aofHandler.AddAof(localIndex, cmd)
		}
	}
	return database
}
func execSelect(conn resp.Connection, database *StandaloneDatabase, args [][]byte, ) resp.Reply {
	num, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("invalid db index")
	}
	if num < 0 || num >= len(database.dbSet) {
		return reply.NewErrReply(" db index is out of range")
	}
	conn.SelectDB(num)
	log.Printf("Selected DB %d", num)
	return reply.NewOkReply()
}
func (d *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	if len(args) == 0 {
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
	log.Printf("Exec command %s on DB %d", cmd, db)
	return d.dbSet[db].Exec(client, args)
}

func (d *StandaloneDatabase) AfterClientClose(c resp.Connection) {
}

func (d *StandaloneDatabase) Close() {
}
