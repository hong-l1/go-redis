package database

import (
	"awesomeProject3/go-redis/ineterface/resp"
	"awesomeProject3/go-redis/pkg/utils"
	"awesomeProject3/go-redis/resp/client"
	"awesomeProject3/go-redis/resp/reply"
	"context"
	"errors"
	"strconv"
)

func (c *ClusterDataBase) getObjectPool(node string) (*client.Client, error) {
	objectPool, ok := c.peerConnection[node]
	if !ok {
		return nil, errors.New("no such node")
	}
	object, err := objectPool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	Client, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("error casting object to client.Client")
	}
	return Client, nil
}
func (c *ClusterDataBase) returnObjectPool(node string, client *client.Client) error {
	objectPool, ok := c.peerConnection[node]
	if !ok {
		return errors.New("no such node")
	}
	return objectPool.ReturnObject(context.Background(), client)
}
func (c *ClusterDataBase) relay(node string, connect resp.Connection, args [][]byte) resp.Reply {
	if node == c.self {
		return c.db.Exec(connect, args)
	}
	Client, err := c.getObjectPool(node)
	if err != nil {
		return reply.NewErrReply(err.Error())
	}
	defer c.returnObjectPool(node, Client)
	Client.Send(utils.ToCmdLine("select", strconv.Itoa(connect.GetDBIndex())))
	return Client.Send(args)
}
func (c *ClusterDataBase) broadcast(connect resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range c.nodes {
		result := c.relay(node, connect, args)
		results[node] = result
	}
	return results
}
