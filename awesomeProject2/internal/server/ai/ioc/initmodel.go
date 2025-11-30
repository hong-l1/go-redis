package ioc

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"log"
	"os"
)

var topP float32 = 0.32

func InitModel() *ark.ChatModel {
	model, err := ark.NewChatModel(context.Background(), &ark.ChatModelConfig{
		TopP:   &topP,
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  "doubao-seed-1-6-lite-251015",
	})
	if err != nil {
		log.Fatal("init model err: ", err)
	}
	return model
}
