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

	"novel2video/backend/text_handler"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/get/novel/fragments", text_handler.GetNovelFragments)
	r.POST("/api/save/novel/fragments", text_handler.SaveCombinedFragments)
	r.Run("localhost:1198")
}
