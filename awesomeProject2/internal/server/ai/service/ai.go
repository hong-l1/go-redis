package service

import (
	"awesomeProject2/internal/api/gen/ai"
	"awesomeProject2/internal/api/gen/order"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
)

type ChatService struct {
	model       *ark.ChatModel
	OrderClient order.OrderServiceClient
	ai.UnimplementedAiAgentServiceServer
}
type ChatServiceImpl interface {
	Chat(ctx context.Context, req *ai.ChatReq) (*ai.ChatResp, error)
}
type Order struct {
	ID     int64  `json:"id"`
	Option string `json:"option"`
}

func NewChatService(model *ark.ChatModel) *ChatService {
	return &ChatService{model: model}
}

func (c *ChatService) GetOrderInfo(ctx context.Context, input Order) (output string, err error) {
	//res, err := c.OrderClient.ListOrder(ctx, &order.ListOrderReq{
	//	UserId: input.ID,
	//})
	//if err != nil {
	//	return "", err
	//}
	fmt.Println(input)
	data, err := json.Marshal(input.ID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(data), nil
}
func (c *ChatService) Chat(ctx context.Context, req *ai.ChatReq) (*ai.ChatResp, error) {
	tool1 := utils.NewTool(&schema.ToolInfo{
		Name: "find_order_service",
		Desc: "用户用来查询订单服务",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"id": {
				Type:     schema.Integer,
				Desc:     "用户id",
				Required: true,
			},
			"option": {
				Type: schema.String,
				Desc: "用户查询的条件",
			},
		}),
	}, c.GetOrderInfo)
	toolinfo, _ := tool1.Info(ctx)
	chain := compose.NewChain[string, *schema.Message]()
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		template := prompt.FromMessages(schema.FormatType(0),
			schema.SystemMessage("你是一个AI助手，请你从用户的输入中提取出id信息和option信息以便后续查询订单"),
			schema.UserMessage("我的id是3，我要查询我未支付的订单"),
			schema.SystemMessage(`Order{
				ID:     3,
				Option: "未支付",
			}`),
		)
		msg, _ := template.Format(ctx, nil)
		res, _ := c.model.Generate(ctx, msg)
		fmt.Printf("%#v", res)
		return res.Content, nil
	})).AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (output *schema.Message, err error) {
		template := prompt.FromMessages(schema.FormatType(0),
			schema.SystemMessage("你是一个AI助手,已经有了订单的id和option,请你调用工具函数完成查询."),
			schema.UserMessage("{text}"),
		)
		msg, err := template.Format(ctx, map[string]interface{}{
			"text": input,
		})
		if err != nil {
			log.Println(err)
			return
		}
		res, err := c.model.Generate(ctx, msg, model.WithTools([]*schema.ToolInfo{toolinfo}))
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(res.Content)
		return res, nil
	}))
	runnable, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}
	res, _ := runnable.Invoke(ctx, req.TextInput)
	return &ai.ChatResp{TextReply: res.Content}, nil
}
