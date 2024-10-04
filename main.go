// 分析文本，把文本进行分割
// 给分割后的文本生成提示词，这里应该是个列表
// 把提示词发给文生图 or 图生图
// 需要有个重新绘制的按钮
// 需要有个保存图片的按钮

// 最重要的是，提取出文本中的角色 保证角色的一致性

package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetNovelFragments(c *gin.Context) {
	file, err := os.Open("a.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	c.JSON(http.StatusOK, lines)
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/get/novel/fragments", GetNovelFragments)
	r.Run("localhost:1198") // 监听并在 0.0.0.0:8080 上启动服务
}
