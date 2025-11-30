package handler

import (
	databaseface "awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/resp/connection"
	"awesomeProject3/go-redis/resp/parser"
	"awesomeProject3/go-redis/resp/reply"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type RespHandler struct {
	activeConn  sync.Map
	db          databaseface.Database
	closingConn atomic.Bool
}

func NewRespHandler(db databaseface.Database) *RespHandler {
	return &RespHandler{
		db: db,
	}
}
func (r *RespHandler) CloseHandleClient(client *connection.Connect) error {
	client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
	return nil
}
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closingConn.Load() {
		conn.Close()
		return
	}
	client := connection.NewConnect(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		fmt.Println(payload)
		if payload.Err != nil {
			if errors.Is(payload.Err, io.EOF) || errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "usr of closed network connection") {
				r.CloseHandleClient(client)
				return
			}
			reply := reply.NewErrReply(payload.Err.Error())
			err := client.Write(reply.ToBytes())
			if err != nil {
				r.CloseHandleClient(client)
				return
			}
			continue
		}
		if payload.Data == nil {
			continue
		}
		val, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			continue
		}
		fmt.Println("val :", val.ToBytes())
		result := r.db.Exec(client, val.Content)
		if result != nil {
			client.Write(result.ToBytes())
		} else {
			client.Write(reply.NewUnKnownErrReply().ToBytes())
		}
	}
}

func (r *RespHandler) Close() error {
	log.Println("closing resp handler")
	r.closingConn.Store(true)
	r.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*connection.Connect)
		client.Close()
		return true
	})
	r.db.Close()
	return nil
}
