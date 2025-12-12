package database

import (
	"awesomeProject3/go-redis/resp/client"
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectFactory struct {
	peer string
}

func (c *connectFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	newClient, err := client.NewClient(c.peer)
	if err != nil {
		return nil, err
	}
	newClient.Start()
	return pool.NewPooledObject(newClient), nil
}
func (c *connectFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	Client, ok := object.Object.(client.Client)
	if !ok {
		return errors.New("object does not implement client.Client")
	}
	Client.Close()
	return nil
}
func (c *connectFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (c *connectFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c *connectFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
