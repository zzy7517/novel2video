package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

var apikey = "sk-Weqj90RMpPbfss51QpqYFvYUwGfC8i1Q1iSs6J45d7DT6nd0"
var baseUrl = "https://api.chatanywhere.tech"
var model = openai.GPT3Dot5Turbo

var groq_llama_base_url = "https://api.groq.com/openai/v1"
var groq_llama_api_key = "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"
var groq_llama_3_1_70b = "llama-3.1-70b-versatile"

var client *openai.Client

func getOpenAiClient() *openai.Client {
	if client == nil {
		config := openai.DefaultConfig(groq_llama_api_key)
		config.BaseURL = groq_llama_base_url
		client = openai.NewClientWithConfig(config)
	}
	return client
}

func queryOpenai(ctx context.Context, input, sys string, modelName string, temperature float32, maxOutputTokens int32) (res string, err error) {
	var M []openai.ChatCompletionMessage
	if len(sys) > 0 {
		M = append(M, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: sys,
		})
	}
	M = append(M, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	resp, err := getOpenAiClient().CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       groq_llama_3_1_70b,
			Messages:    M,
			Temperature: temperature,
			MaxTokens:   int(maxOutputTokens),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
