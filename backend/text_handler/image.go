package text_handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"novel2video/backend"
	"novel2video/backend/image"
	"novel2video/backend/util"
)

func GenerateImage(c *gin.Context) {
	err := os.RemoveAll(util.ImageDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.ImageDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	lines, err := readLinesFromDirectory("temp/prompts")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	go func() {
		for i, p := range lines {
			pi := p
			err := image.GenerateImage(pi, 114514191981, 540, 960, i)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}()
}

func GetLocalImages(c *gin.Context) {
	files, err := os.ReadDir(util.ImageDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取图像目录失败"})
		return
	}
	imageMap := make(map[string]string)
	re := regexp.MustCompile(`(\d+)\.png`) // 从文件名中提取数字
	now := time.Now().Unix()
	for _, file := range files {
		if !file.IsDir() {
			matches := re.FindStringSubmatch(file.Name())
			if len(matches) > 1 {
				key := matches[1]
				absPath := filepath.Join("/images", file.Name())
				imageMap[key] = absPath + fmt.Sprintf("?v=%d", now)
			}
		}
	}

	c.JSON(http.StatusOK, imageMap)
}
