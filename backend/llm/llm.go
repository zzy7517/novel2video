package llm

import (
	"context"
	"fmt"
)

func QueryLLM(ctx context.Context, input, systemContent string, modelName string, temperature float32, maxOutputTokens int32) (res string, err error) {
	// return QueryGemini(ctx, input, systemContent, "gemini-1.5-pro-002", temperature, maxOutputTokens)
	return querySambaNova(ctx, input, systemContent, modelName, temperature)
	switch modelName {
	case "doubao":
		return queryVolEngine(ctx, input, systemContent, 0.98)
	case "gemini-1.5-pro-002":
		return QueryGemini(ctx, input, systemContent, modelName, temperature, maxOutputTokens)
	}
	return "", fmt.Errorf("wrong llm")
}
