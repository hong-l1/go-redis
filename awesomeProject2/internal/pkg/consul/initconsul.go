package consul

import "github.com/hashicorp/consul/api"

func InitConsul() *api.Client {
	client, err := api.NewClient(&api.Config{
		Address: "43.142.57.35:8500",
	})
	if err != nil {
		panic(err)
	}
	return client
}
