package llm

import (
	"context"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	doubao "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
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
	var M []*doubao.ChatCompletionMessage
	if len(sys) > 0 {
		M = append(M, &doubao.ChatCompletionMessage{
			Role: doubao.ChatMessageRoleSystem,
			Content: &doubao.ChatCompletionMessageContent{
				StringValue: volcengine.String(sys),
			},
		})
	}
	M = append(M, &doubao.ChatCompletionMessage{
		Role: doubao.ChatMessageRoleUser,
		Content: &doubao.ChatCompletionMessageContent{
			StringValue: volcengine.String(prompt),
		},
	})
	req := doubao.ChatCompletionRequest{
		MaxTokens:   4000,
		Model:       accessPoint,
		Messages:    M,
		Temperature: temperature,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil
}
