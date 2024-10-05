package text_handler

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"novel2video/backend"
)

func SaveCombinedFragments(c *gin.Context) {
	var fragments []string
	if err := c.ShouldBindJSON(&fragments); err != nil {
		backend.HandleError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}
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
	err = saveListToFiles(fragments, "temp/fragments/")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to save", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fragments saved successfully"})
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

func saveListToFiles(in []string, path string) error {
	for i, line := range in {
		filePath := fmt.Sprintf(path+"%d.txt", i)
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
	// 按照数字顺序从小到大读取
	// 正则表达式用于提取文件名中的数字
	re := regexp.MustCompile(`\d+`)

	// 创建一个切片来存储文件名和对应的数字
	type fileWithNumber struct {
		name   string
		number int
	}

	var fileList []fileWithNumber

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			matches := re.FindStringSubmatch(file.Name())
			if len(matches) > 0 {
				number, err := strconv.Atoi(matches[0])
				if err == nil {
					fileList = append(fileList, fileWithNumber{name: file.Name(), number: number})
				}
			}
		}
	}

	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].number < fileList[j].number
	})

	var lines []string
	for _, file := range fileList {
		content, err := os.ReadFile(filepath.Join(dir, file.name))
		if err != nil {
			// 打印错误并继续处理其他文件
			fmt.Printf("Error reading file %s: %v\n", file.name, err)
			continue
		}
		lines = append(lines, string(content))
	}
	return lines, nil
}
