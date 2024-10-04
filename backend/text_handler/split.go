package text_handler

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"novel2video/backend"
)

func GetCombinedNovelFragments(c *gin.Context) {

}

func GetNovelFragments(c *gin.Context) {
	err := os.RemoveAll("temp/fragments")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll("temp/fragments", os.ModePerm)
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
	lines, err := readLinesFromDirectory("temp/fragments")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}

	c.JSON(http.StatusOK, lines)
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
			filePath := fmt.Sprintf("temp/fragments/%d.txt", lineNumber)
			err := os.WriteFile(filePath, []byte(line), 0644)
			if err != nil {
				return err
			}
			lineNumber++
		}
	}
	return scanner.Err()
}

func readLinesFromDirectory(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			content, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
			if err != nil {
				return nil, err
			}
			lines = append(lines, string(content))
		}
	}
	return lines, nil
}
