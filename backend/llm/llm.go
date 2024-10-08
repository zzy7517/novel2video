package llm

import (
	"context"
)

func QueryLLM(ctx context.Context, input, systemContent string, modelName string, temperature float32, maxOutputTokens int32) (res string, err error) {
	// return QueryGemini(ctx, input, systemContent, "gemini-1.5-pro-002", temperature, maxOutputTokens)
	return querySambaNova(ctx, input, systemContent, modelName, temperature)
	//switch modelName {
	//case "doubao":
	//	return queryVolEngine(ctx, input, systemContent, 0.98)
	//case "gemini-1.5-pro-002":
	//	return QueryGemini(ctx, input, systemContent, modelName, temperature, maxOutputTokens)
	//}
	//return "", fmt.Errorf("wrong llm")
}

var transLateSys = `把输入完全翻译成英文，不要输出翻译文本以外的内容，只需要输出翻译后的文本。
如果包含翻译之外的内容，则重新输出`

func LLMTranslate(ctx context.Context, input string) (output string, err error) {
	return querySiliconFlow(ctx, input, transLateSys, 0.01)
}
