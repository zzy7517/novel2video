package text_handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
)

var sys = `
#Task: #
Extract stable diffusion prompts from the given inputs.

#Rules:#
1.The extracted prompts should be as simple and concrete as possible. Focus mainly on someone doing something with someone at somewhere.
2.You can describe the environment in detail.
3.The prompt should consist of detailed words, not sentences.
4.Each input should have one output
6.I need detailed people name in the prompt

#Input Format:#
1.text1
2.text2
3.text3

#Output Format:#
each line in input should have one output, and the outputs should be separated by a new line
`

func ExtractPromptFromTexts(c *gin.Context) {
	//lines, err := readLinesFromDirectory("temp/fragments")
	//if err != nil {
	//	backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
	//	return
	//}
	//promptsMid := generateInputPrompts(lines)
	//var t2iPrompts []string
	//for _, p := range promptsMid {
	//	res, err := llm.QueryGemini(c.Request.Context(), p, sys, "gemini-1.5-pro-002", 0.01, 8192)
	//	if err != nil {
	//		logrus.Errorf("query gemini failed, err is %v", err)
	//		continue
	//	}
	//	lines := strings.Split(res, "\n")
	//	// 编译正则表达式，用于匹配行首的序号和点
	//	re := regexp.MustCompile(`^\d+\.\s*`)
	//	// 遍历每一行，去掉序号和点
	//	for _, line := range lines {
	//		line = re.ReplaceAllString(line, "")
	//		if line != "" {
	//			t2iPrompts = append(t2iPrompts, line)
	//		}
	//	}
	//}
	//err = os.RemoveAll("temp/prompts")
	//if err != nil {
	//	backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
	//	return
	//}
	//err = os.MkdirAll("temp/prompts", os.ModePerm)
	//if err != nil {
	//	backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
	//	return
	//}
	//err = saveListToFiles(t2iPrompts, "temp/prompts/")
	//if err != nil {
	//	backend.HandleError(c, http.StatusInternalServerError, "save list to file failed", err)
	//	return
	//}
	// todo 直接读取完了返回
	// 从目录中读取所有文件并返回内容
	lines, err := readLinesFromDirectory("temp/prompts")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	logrus.Infof("extract prompts from novel fragments finished")
	c.JSON(http.StatusOK, lines)
}

// 每100个生成为一组，发给ai，提取提示词
func generateInputPrompts(list []string) []string {
	var prompts []string
	for i := 0; i < len(list); i += 100 {
		end := i + 100
		if end > len(list) {
			end = len(list)
		}
		var builder strings.Builder
		for j, v := range list[i:end] {
			builder.WriteString(strconv.Itoa(j+1) + ". ")
			builder.WriteString(fmt.Sprintf("%v", v))
			if j != len(list[i:end])-1 {
				builder.WriteString("\n")
			}
		}
		prompt := builder.String()
		logrus.Infof("prompt is %v", prompt)
		prompts = append(prompts, prompt)
	}
	return prompts
}
