package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"tfidf/internal/model"
	"tfidf/internal/parser"
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

	words := parser.ExtractWords(content)
	stats := service.CalculateTFIDF(words)
	top50 := model.TopIDFRange(stats, 0, 50)

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Words": top50,
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
