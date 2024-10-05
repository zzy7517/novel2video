package text_handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend"
	"novel2video/backend/llm"
	"novel2video/backend/util"
)

var extractCharacterSys = `
	#Task: #
	Extract characters from the novel fragment
	
	#Rule#
	the extracted names should be in English
	
	#Output Format:#
	name1, name2, name3
`

func GetCharacters(c *gin.Context) {
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
			builder.WriteString(strconv.Itoa(j) + ". ")
			builder.WriteString(fmt.Sprintf("%v", v))
			if j != len(lines[i:end])-1 {
				builder.WriteString("\n")
			}
		}
		prompt := builder.String()
		// logrus.Infof("prompt is %v", prompt)
		prompts = append(prompts, prompt)
	}

	for _, p := range prompts {
		res, err := llm.QueryGemini(c.Request.Context(), p, extractCharacterSys, "gemini-1.5-pro-002", 0.01, 8192)
		if err != nil {
			logrus.Errorf("query gemini failed, err is %v", err)
			continue
		}
		for _, ch := range strings.Split(res, ",") {
			characterMap[ch] = ch
		}
	}
	c.JSON(http.StatusOK, characterMap)
}

func PatchCharacters(c *gin.Context) {

}
