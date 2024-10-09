package text_handler

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"novel2video/backend"
	"novel2video/backend/util"
)

func SaveCombinedFragments(c *gin.Context) {
	var fragments []string
	if err := c.ShouldBindJSON(&fragments); err != nil {
		backend.HandleError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}
	err := os.RemoveAll(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.NovelFragmentsDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	err = saveListToFiles(fragments, util.NovelFragmentsDir+"/", 0)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to save", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fragments saved successfully"})
}

func GetNovelFragments(c *gin.Context) {
	err := os.RemoveAll(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll(util.NovelFragmentsDir, os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	// 读取文件并保存每一行到单独的文件
	err = saveLinesToFiles("a.txt")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to process file", err)
		return
	}
	// 从目录中读取所有文件并返回内容
	lines, err := util.ReadLinesFromDirectory(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}

	c.JSON(http.StatusOK, lines)
}

func saveListToFiles(in []string, path string, offset int) error {
	for i, line := range in {
		filePath := fmt.Sprintf(path+"%d.txt", i+offset)
		err := os.WriteFile(filePath, []byte(line), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveLinesToFiles(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			filePath := fmt.Sprintf(util.NovelFragmentsDir+"/%d.txt", lineNumber)
			err := os.WriteFile(filePath, []byte(line), 0644)
			if err != nil {
				return err
			}
			lineNumber++
		}
	}
	return scanner.Err()
}

func GetInitial(c *gin.Context) {
	type InitialData struct {
		Fragments []string `json:"fragments"`
		Images    []string `json:"images"`
		Prompts   []string `json:"prompts"`
		PromptsEn []string `json:"promptsEn"`
	}
	novels, err := util.ReadLinesFromDirectory(util.NovelFragmentsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	prompts, err := util.ReadLinesFromDirectory(util.PromptsDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read prompts", err)
		return
	}
	promptsEn, err := util.ReadLinesFromDirectory(util.PromptsEnDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read prompts", err)
		return
	}
	files, err := util.ReadFilesFromDirectory(util.ImageDir)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "failed to read images", err)
		return
	}
	var images []string
	now := time.Now().Unix()
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, filepath.Join("/images", file.Name())+fmt.Sprintf("?v=%d", now))
		}
	}
	data := InitialData{
		Fragments: novels,
		Images:    images,
		Prompts:   prompts,
		PromptsEn: promptsEn,
	}
	c.JSON(http.StatusOK, data)
}
