package main

import (
	database3 "awesomeProject3/go-redis/cluster"
	config2 "awesomeProject3/go-redis/config"
	"awesomeProject3/go-redis/database"
	database2 "awesomeProject3/go-redis/ineterface/database"
	consisithash "awesomeProject3/go-redis/pkg/consisit_hash"
	"awesomeProject3/go-redis/resp/handler"
	"awesomeProject3/go-redis/tcp"
	"fmt"
	"log"
)

func main() {
	nodemap := consisithash.NewNodeMap()
	var tempDB database2.Database
	if config2.RedisConfig.Self != "" && len(config2.RedisConfig.Peers) > 0 {
		tempDB = database3.NewClusterDataBase(nodemap)
	} else {
		tempDB = database.NewDatabase()
	}
	handle := handler.NewRespHandler(tempDB)
	fmt.Println("listen on:", config2.RedisConfig.Self)
	err := tcp.ListenAndServerWithSignal(config2.RedisConfig.Self, handle)
	if err != nil {
		log.Println(err)
	}
}
