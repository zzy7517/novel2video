package backend

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, gin.H{"error": message})
}
