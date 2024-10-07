package llm

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const (
	projectID = "p-sgx9cobj"
	location  = "asia-east1"
)

var geminiClient *genai.Client

func QueryGemini(ctx context.Context, input, systemContent string, modelName string, temperature float32, maxOutputTokens int32) (res string, err error) {
	client := GetClient()
	if client == nil {
		return res, fmt.Errorf("error getting gemini client")
	}
	gemini := client.GenerativeModel(modelName)
	gemini.SetTemperature(temperature)
	if maxOutputTokens != 0 {
		gemini.SetMaxOutputTokens(maxOutputTokens)
	}
	gemini.SafetySettings = GetSafetySetting()
	if len(systemContent) > 0 {
		var sysPart []genai.Part
		sysPart = append(sysPart, genai.Text(systemContent))
		gemini.SystemInstruction = &genai.Content{Parts: sysPart}
	}
	var prompt []genai.Part
	prompt = append(prompt, genai.Text(input))
	resp, err := gemini.GenerateContent(ctx, prompt...)
	if err != nil {
		return res, err
	}
	if len(resp.Candidates) == 0 {
		return res, fmt.Errorf("get result from gemini failed, nil response")
	}
	if resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return res, fmt.Errorf("get result from gemini failed, nil response")
	}
	t, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if ok {
		return string(t), nil
	}
	return res, fmt.Errorf("get result from gemini failed")
}

func InitClient() {
	var err error
	options := getGrpcOptions()
	geminiClient, err = genai.NewClient(context.Background(), projectID, location, options...)
	if err != nil {
		log.Fatalf("error creating gemini client: %+v,location:%s", err, location)
	}
}

func GetClient() *genai.Client {
	if geminiClient == nil {
		InitClient()
	}
	return geminiClient
}

func getGrpcOptions() []option.ClientOption {
	credentialsFile := "/Users/zhongyuanzhang/Desktop/gemini-credential.json"
	return []option.ClientOption{
		option.WithCredentialsFile(credentialsFile),
		option.WithGRPCDialOption(grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)),
		// option.WithEndpoint("asia-east1-aiplatform.googleapis.com:443"),
		// option.WithGRPCConnectionPool(500),
	}
}

func GetSafetySetting() []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryUnspecified,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}
}
