// 分析文本，把文本进行分割
// 给分割后的文本生成提示词，这里应该是个列表
// 把提示词发给文生图 or 图生图
// 需要有个重新绘制的按钮
// 需要有个保存图片的按钮

// 最重要的是，提取出文本中的角色 保证角色的一致性

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"novel2video/backend/text_handler"
	"novel2video/backend/util"
)

func main() {
	r := gin.Default()

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/get/novel/fragments", text_handler.GetNovelFragments)       // 分割文本
	r.POST("/api/save/novel/fragments", text_handler.SaveCombinedFragments) // 合并文本
	r.GET("/api/get/novel/prompts", text_handler.ExtractPromptFromTexts)    // 提取文生图prompt
	r.POST("/api/novel/images", text_handler.GenerateImage)                 // 一键生成
	r.GET("/api/novel/images", text_handler.GetLocalImages)                 // 刷新图片
	r.GET("/api/novel/characters", text_handler.GetCharacters)              // 提取角色
	r.PATCH("/api/novel/characters", text_handler.PatchCharacters)          // 修改角色

	r.Static("/images", util.ImageDir)
	err := r.Run("localhost:1198")
	if err != nil {
		return
	}
}
