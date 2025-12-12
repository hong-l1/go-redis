package aof

import (
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/pkg/utils"
	"awesomeProject3/go-redis/resp/connection"
	"awesomeProject3/go-redis/resp/parser"
	"awesomeProject3/go-redis/resp/reply"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var DefaultName = "aof.txt"

type CmdLine [][]byte

type payload struct {
	cmdLine CmdLine
	dbIndex int
}
type AofHandler struct {
	db          database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

func NewAofHandler(database database.Database) (*AofHandler, error) {
	var handlerAof = new(AofHandler)
	handlerAof.db = database
	handlerAof.aofFilename = DefaultName
	path := filepath.Join("./", handlerAof.aofFilename)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	handlerAof.aofFile = file
	handlerAof.LoadAof()
	handlerAof.aofChan = make(chan *payload, 1024)
	go func() {
		handlerAof.HandlerAof()
	}()
	return handlerAof, nil
}
func (h *AofHandler) AddAof(dbIndex int, cmdLine CmdLine) {
	h.aofChan <- &payload{
		cmdLine: cmdLine,
		dbIndex: dbIndex,
	}
}
func (h *AofHandler) HandlerAof() {
	for data := range h.aofChan {
		if data.dbIndex != h.currentDB {
			cmd := utils.ToCmdLine("select", strconv.Itoa(data.dbIndex))
			bytes := reply.NewMultiBulkReply(cmd).ToBytes()
			_, err := h.aofFile.Write(bytes)
			if err != nil {
				log.Println("Error writing select to aof file:", err)
				continue
			}
			h.currentDB = data.dbIndex
			log.Printf("Switched AOF DB to %d", h.currentDB)
		}
		bytes := reply.NewMultiBulkReply(data.cmdLine).ToBytes()
		_, err := h.aofFile.Write(bytes)
		if err != nil {
			log.Println("Error writing to aof file:", err)
			continue
		}
	}
}
func (h *AofHandler) LoadAof() {
	file, err := os.Open(h.aofFilename)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}
	ch := parser.ParseStream(file)
	fakeConn := &connection.Connect{}
	for data := range ch {
		if data.Err != nil {
			if data.Err == io.EOF {
				break
			}
			log.Println(data.Err)
			continue
		}
		r, ok := data.Data.(*reply.MultiBulkReply)
		if !ok {
			log.Println("need multi bulk reply")
			continue
		}
		resp := h.db.Exec(fakeConn, r.Content)
		if reply.IsErrorReply(resp) {
			log.Println(resp)
			continue
		}
	}
	h.currentDB = fakeConn.GetDBIndex()
}
