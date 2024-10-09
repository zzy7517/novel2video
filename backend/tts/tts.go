package tts

import "github.com/gin-gonic/gin"

func GenerateAudioFiles(c *gin.Context) {
	byEdgeTTS()
}
