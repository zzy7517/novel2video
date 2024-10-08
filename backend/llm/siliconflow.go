package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var (
	// https://siliconflow.cn/zh-cn/models
	siliconFlowFreeModels = []string{
		"Qwen/Qwen2-7B-Instruct",
		"Qwen/Qwen2.5-7B-Instruct",
		"THUDM/glm-4-9b-chat",
		"01-ai/Yi-1.5-9B-Chat-16K",
		"internlm/internlm2_5-7b-chat",
		"google/gemma-2-9b-it",
		"meta-llama/Meta-Llama-3-8B-Instruct",
		"meta-llama/Meta-Llama-3.1-8B-Instruct",
	}
	modelIndex int32
)

func getNextModel() string {
	index := atomic.AddInt32(&modelIndex, 1) % int32(len(siliconFlowFreeModels))
	return siliconFlowFreeModels[index]
}

func querySiliconFlow(ctx context.Context, input, sys string, temperature float32) (res string, err error) {
	url := "https://api.siliconflow.cn/v1/chat/completions"
	var messages []map[string]string
	if len(sys) > 0 {
		messages = append(messages, map[string]string{"role": "system", "content": sys})
	}
	messages = append(messages, map[string]string{"role": "user", "content": input})
	sFModel := getNextModel()
	requestBody := map[string]interface{}{
		"temperature": temperature,
		"messages":    messages,
		"model":       sFModel,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+siliconFLowApiKey)
	snClient := &http.Client{}
	resp, err := snClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Errorf("query siliconFlow error %v, modelName %v", err, sFModel)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("query siliconFlow Unexpected response status: %v, modelName %v", body, sFModel)
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	} else {
		return "", fmt.Errorf("No choices found in response.")
	}
}
