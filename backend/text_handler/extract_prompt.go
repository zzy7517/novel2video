package text_handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
	"novel2video/backend/llm"
)

var sys = `
#Task: #
Extract scenes from the given inputs.

#Rules:#
1. 越简单，越具体越好，比如Mike is eating seafood或者sun rises
2. 关于地点和环境，可以详细一点
3. 每个输入必须要有一个输出
4. 如果没有出现人，可以只描述下场景
5. 如果不是很具体的内容，可以输出null
6. 如果某一行输入无法提取出内容，输出null，但不要不输出
6. 输出需要使用英文，一定不要出现中文，人名也要翻译成英文

#Input Format:#
1.text1
2.text2
3.text3
每个数字开头的行代表一个输入，每行输入必须对应一行输出

#Output Format:#
每行输入需要对应一行输出，每个输出用空行隔开
如果输出的人名是中文，改成英文

# 检查 #
输出的行数是否和输入一致，如果不一致，则重新生成输出内容
`

func ExtractPromptFromTexts(c *gin.Context) {
	lines, err := readLinesFromDirectory("temp/fragments")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	promptsMid := generateInputPrompts(lines, 50)
	var t2iPrompts []string
	for _, p := range promptsMid {
		res, err := llm.QueryLLM(c.Request.Context(), p, sys, "doubao", 0.01, 8192)
		if err != nil {
			logrus.Errorf("query gemini failed, err is %v", err)
			continue
		}
		lines := strings.Split(res, "\n")
		// 编译正则表达式，用于匹配行首的序号和点
		re := regexp.MustCompile(`^\d+\.\s*`)
		// 遍历每一行，去掉序号和点
		for _, line := range lines {
			line = re.ReplaceAllString(line, "")
			if line != "" {
				t2iPrompts = append(t2iPrompts, line)
			}
		}
	}
	err = os.RemoveAll("temp/prompts")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll("temp/prompts", os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	for i, element := range t2iPrompts {
		for key, value := range characterMap {
			if strings.Contains(element, key) {
				t2iPrompts[i] = strings.ReplaceAll(element, key, value)
			}
		}
	}
	err = saveListToFiles(t2iPrompts, "temp/prompts/")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "save list to file failed", err)
		return
	}
	// 从目录中读取所有文件并返回内容
	lines, err = readLinesFromDirectory("temp/prompts")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	logrus.Infof("extract prompts from novel fragments finished")
	c.JSON(http.StatusOK, lines)
}

// 每100个生成为一组，发给ai，提取提示词
func generateInputPrompts(list []string, step int) []string {
	var prompts []string
	for i := 0; i < len(list); i += step {
		end := i + step
		if end > len(list) {
			end = len(list)
		}
		var builder strings.Builder
		for j, v := range list[i:end] {
			builder.WriteString(strconv.Itoa(j) + ". ")
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
