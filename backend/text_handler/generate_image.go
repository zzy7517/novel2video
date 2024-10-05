package text_handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"novel2video/backend"
	"novel2video/backend/image"
)

func GenerateImage(c *gin.Context) {
	err := os.RemoveAll("temp/image")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to remove directory", err)
		return
	}
	err = os.MkdirAll("temp/image", os.ModePerm)
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to create directory", err)
		return
	}
	lines, err := readLinesFromDirectory("temp/prompts")
	if err != nil {
		backend.HandleError(c, http.StatusInternalServerError, "Failed to read fragments", err)
		return
	}
	for i, p := range lines {
		err := image.GenerateImage(p, 114514191981, 540, 960, i)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
