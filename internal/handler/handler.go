package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"tfidf/internal/service"
)

func ShowForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func UploadFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ошибка загрузки файла",
		})
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ошибка чтения файла",
		})
	}

	top50, err := service.ProcessFile(content, 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Words": top50,
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
