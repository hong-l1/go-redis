package es

import "github.com/elastic/go-elasticsearch/v8"

func InitEs() *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://43.142.57.35:9200"},
	})
	if err != nil {
		panic(err)
	}
	return client
}
