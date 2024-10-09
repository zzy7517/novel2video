package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ReadLinesFromDirectory(dir string) ([]string, error) {
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

func ReadFilesFromDirectory(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// Create a slice to hold the files along with their extracted numbers
	type fileWithNumber struct {
		entry  os.DirEntry
		number int
	}
	var fileList []fileWithNumber
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			// Extract the numeric part from the file name
			numberStr := strings.TrimSuffix(name, filepath.Ext(name))
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				fmt.Printf("failed to convert %s to number: %s\n", numberStr, err)
				continue
			}
			fileList = append(fileList, fileWithNumber{entry: file, number: number})
		}
	}
	// Sort the files based on the extracted number
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].number < fileList[j].number
	})
	// Extract the sorted os.DirEntry from the sorted fileWithNumber slice
	sortedFiles := make([]os.DirEntry, len(fileList))
	for i, file := range fileList {
		sortedFiles[i] = file.entry
	}
	return sortedFiles, nil
}
