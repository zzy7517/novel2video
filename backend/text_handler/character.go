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
	1. 提取出所有的人名，把人名翻译成英文输出，只需要输出英文
	2. 所有的人名，别名，称呼，包括对话中引用到的名字都需要提取
	3. 不要输出中文
	
	#Output Format:#
	name1, name2, name3, name4...
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
	lines, err := readLinesFromDirectory("temp/fragments")
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

func GetRandomAppearance(c *gin.Context) {
	prompt := "随机生成一个二次元角色的外形描述，包括年龄发色眼睛穿着等等，使用英文输出"
	appearance, err := llm.QueryLLM(c.Request.Context(), prompt, "", "doubao", 1, 100)
	if err != nil {
		logrus.Errorf("get random appearance from llm failed, err %v", err)
		backend.HandleError(c, http.StatusBadRequest, "get random appearance from llm failed", nil)
		return
	}
	c.JSON(http.StatusOK, appearance)
}
