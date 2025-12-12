package database

import (
	"awesomeProject3/go-redis/config"
	database2 "awesomeProject3/go-redis/database"
	"awesomeProject3/go-redis/ineterface/database"
	"awesomeProject3/go-redis/ineterface/resp"
	consisithash "awesomeProject3/go-redis/pkg/consisit_hash"
	"awesomeProject3/go-redis/resp/reply"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

type ClusterDataBase struct {
	self           string
	nodes          []string
	peerPicker     *consisithash.NodeMap
	peerConnection map[string]*pool.ObjectPool
	db             database.Database
}

func NewClusterDataBase(peerPicker *consisithash.NodeMap) *ClusterDataBase {
	cluster := &ClusterDataBase{
		self:           config.RedisConfig.Self,
		db:             database2.NewDatabase(),
		peerPicker:     peerPicker,
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.RedisConfig.Peers)+1)
	for _, node := range config.RedisConfig.Peers {
		nodes = append(nodes, node)
	}
	nodes = append(nodes, config.RedisConfig.Self)
	cluster.nodes = nodes
	peerPicker.AddNode(nodes)
	ctx := context.Background()
	for _, node := range nodes {
		cluster.peerConnection[node] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectFactory{node})
	}
	return cluster
}

var cmdMap = NewRouter()

type cmdline func(c *ClusterDataBase, args [][]byte, connect resp.Connection) resp.Reply

func (c *ClusterDataBase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	cmd := strings.ToLower(string(args[0]))
	exec, ok := cmdMap[cmd]
	if !ok {
		return reply.NewErrReply("cmd not found")
	}
	return exec(c, args, client)
}

func (c *ClusterDataBase) AfterClientClose(conn resp.Connection) {
	c.db.AfterClientClose(conn)
}

func (c *ClusterDataBase) Close() {
	c.db.Close()
}
