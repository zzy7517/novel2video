package llm

import (
	"context"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type Client struct {
	apiClient *arkruntime.Client
}

var arkClient *arkruntime.Client

func GetDoubaoClient() *arkruntime.Client {
	if arkClient == nil {
		arkClient = arkruntime.NewClientWithApiKey(APIKEY)
	}
	return arkClient
}

func queryVolEngine(ctx context.Context, prompt, sys string, temperature float32) (string, error) {
	client := GetDoubaoClient()
	req := model.ChatCompletionRequest{
		MaxTokens: 4000,
		Model:     accessPoint,

		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(sys),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(prompt),
				},
			},
		},
		Temperature: temperature,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil
}