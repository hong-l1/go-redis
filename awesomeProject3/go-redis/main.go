package main

import (
	"awesomeProject3/go-redis/database"
	"awesomeProject3/go-redis/resp/handler"
	"awesomeProject3/go-redis/tcp"
	"log"
)

func main() {
	tempDB := database.NewDatabase()
	handle := handler.NewRespHandler(tempDB)
	err := tcp.ListenAndServerWithSignal(tcp.Config{Addr: "127.0.0.1:9090"}, handle)
	if err != nil {
		log.Println(err)
	}
}
