package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.Run(":9001")

}

func uploadFile(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	filePath := filepath.Join("uploads", f.Filename)

	if err := c.SaveUploadedFile(f, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to save file"})
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "saved"})
}
