package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var llama_405b = "Meta-Llama-3.1-405B-Instruct"

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

// free llama api
func querySambaNova(ctx context.Context, input, sys string, modelName string, temperature float32) (res string, err error) {
	url := "https://api.sambanova.ai/v1/chat/completions"
	var messages []map[string]string
	if len(sys) > 0 {
		messages = append(messages, map[string]string{"role": "system", "content": sys})
	}
	messages = append(messages, map[string]string{"role": "user", "content": input})
	requestBody := map[string]interface{}{
		"temperature": temperature,
		"messages":    messages,
		"model":       llama_405b,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+SambaNovaApiKey)
	snClient := &http.Client{}
	resp, err := snClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected response status: %s\n", resp.Status)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	} else {
		return "", fmt.Errorf("No choices found in response.")
	}
}
