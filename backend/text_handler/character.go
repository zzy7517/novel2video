package text_handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
	"novel2video/backend/llm"
	"novel2video/backend/util"
)

var characterMap = make(map[string]string)

var extractCharacterSys = `
	#Task: #
	Extract characters from the novel fragment
	
	#Rule#
	1. 提取出所有的人名
	2. 所有的人名，别名，称呼，包括对话中引用到的名字都需要提取
    3. 所有出现过的和人有关的称呼都需要提取
	
	#Output Format:#
	名字1, 名字2, 名字3, ...
`

func GetCharacters(c *gin.Context) {
	if len(characterMap) > 0 {
		c.JSON(http.StatusOK, characterMap)
		return
	}
	err := os.RemoveAll(util.CharacterDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.CharacterDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	lines, err := readLinesFromDirectory(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	// 每500行发送给ai一次
	var prompts []string
	characterMap := make(map[string]string)
	for i := 0; i < len(lines); i += 500 {
		end := i + 100
		if end > len(lines) {
			end = len(lines)
		}
		var builder strings.Builder
		for j, v := range lines[i:end] {
			builder.WriteString(fmt.Sprintf("%v", v))
			if j != len(lines[i:end])-1 {
				builder.WriteString("\n")
			}
		}
		prompt := builder.String()
		prompts = append(prompts, prompt)
	}

	for _, p := range prompts {
		res, err := llm.QueryLLM(c.Request.Context(), p, extractCharacterSys, "doubao", 0.01, 8192)
		if err != nil {
			logrus.Errorf("query doubao failed, err is %v", err)
			continue
		}
		for _, ch := range strings.Split(res, ",") {
			characterMap[ch] = ch
		}
	}
	c.JSON(http.StatusOK, characterMap)
}

func PutCharacters(c *gin.Context) {
	var descriptions map[string]string
	if err := c.ShouldBindJSON(&descriptions); err != nil {
		backend.HandleError(c, http.StatusBadRequest, `"error":"Invalid JSON"`, err)
		return
	}

	if len(descriptions) <= 0 {
		backend.HandleError(c, http.StatusBadRequest, `"error":"find no description`, nil)
	}
	for k, v := range descriptions {
		characterMap[k] = v
	}
	c.JSON(http.StatusOK, gin.H{"message": "Descriptions updated successfully"})
}

var appearancePrompt = `
随机生成动漫角色的外形描述，输出简练，以一组描述词的形式输出，每个描述用逗号隔开
数量：一个
包含：性别，年龄，衣着，脸型，眼睛，发色，发型
使用英文输出`

// todo 感觉这个api需要适配一下topp & topk
func GetRandomAppearance(c *gin.Context) {
	prompt := appearancePrompt
	appearance, err := llm.QueryLLM(c.Request.Context(), prompt, "", "doubao", 1, 100)
	if err != nil {
		logrus.Errorf("get random appearance from llm failed, err %v", err)
		backend.HandleError(c, http.StatusBadRequest, "get random appearance from llm failed", nil)
		return
	}
	c.JSON(http.StatusOK, appearance)
}
