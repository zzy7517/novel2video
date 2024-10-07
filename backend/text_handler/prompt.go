package text_handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
	"novel2video/backend/llm"
	"novel2video/backend/util"
)

var sys = `
#Task: #
从输入中提取画面信息

#Rules:#
1. 越简单，越具体越好，比如Mike is eating seafood或者sun rises
2. 关于文本中描述的场景信息，可以详细一点
3. 每个输入必须要有一个输出
4. 如果出现了人物，需要具体到名字
5. 如果没有出现人，可以只描述下场景
6. 如果不是很具体的内容，可以挑选文本中某个词语发散一下
7. 如果某一行输入无法提取出内容，输出“无”，但不要不输出


#Input Format:#
0.输入0
1.输入1
2.输入2
每个数字开头的行代表一个输入，每行输入必须对应一行输出

#Output Format:#
每行输入需要对应一行输出，每个输出用空行隔开
输入和输出的序号需要对应
不要输出多余内容
这是一个输出示例
0. 输出0
1. 输出1
2. 输出2

# 检查 #
输出的行数是否和输入一致，如果不一致，则重新生成输出内容
`

var fragmentsLen = 30

func ExtractSceneFromTexts(c *gin.Context) {
	lines, err := readLinesFromDirectory(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	err = os.RemoveAll(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.PromptsDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	var offset int
	promptsMid := generateInputPrompts(lines, fragmentsLen)
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
		var t2iPrompts []string
		for _, line := range lines {
			line = re.ReplaceAllString(line, "")
			if len(strings.TrimSpace(line)) > 0 {
				t2iPrompts = append(t2iPrompts, line)
				offset++
			}
		}

		err = saveListToFiles(t2iPrompts, util.PromptsDir+"/", offset-len(t2iPrompts))
		if err != nil {
			backend.HandleError(c, http.StatusInternalServerError, "save list to file failed", err)
			return
		}
	}
	// 从目录中读取所有文件并返回内容
	lines, err = readLinesFromDirectory(util.PromptsDir)
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

func GetPromptsEn(c *gin.Context) {
	err := os.RemoveAll(util.PromptsEnDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.PromptsEnDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	lines, err := readLinesFromDirectory(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	for i, element := range lines {
		for key, value := range characterMap {
			if strings.Contains(element, key) {
				lines[i] = strings.ReplaceAll(element, key, value)
			}
		}
	}
	logrus.Infof("translate prompts to English finished")
	c.JSON(http.StatusOK, lines)
}

func SavePromptEn(c *gin.Context) {
	type SaveRequest struct {
		Index   int    `json:"index"`
		Content string `json:"content"`
	}
	var req SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		backend.HandleError(c, http.StatusBadRequest, "parse request body failed", err)
		return
	}

	filePath := filepath.Join(util.PromptsEnDir, fmt.Sprintf("%d.txt", req.Index))
	if err := os.MkdirAll(util.PromptsEnDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attachment saved successfully"})
}
