package main

import (
	"awesomeProject2/internal/api/gen/ai"
	"awesomeProject2/internal/server/ai/ioc"
	"awesomeProject2/internal/server/ai/service"
	"context"
	"fmt"
	"log"
)

func main() {
	model := ioc.InitModel()
	svc := service.NewChatService(model)
	res, err := svc.Chat(context.Background(), &ai.ChatReq{TextInput: "我的id是3，我要查询我未支付的订单"})
	if err != nil {
		println("c")
		log.Println(err)
		return
	}
	fmt.Println("查询结果：", res)
}
