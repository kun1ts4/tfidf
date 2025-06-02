package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"tfidf/internal/service"
	"time"
)

func ShowForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) UploadFile(c *gin.Context) {
	startUploadFileTime := time.Now()
	file, fileHeader, err := c.Request.FormFile("file")
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

	ctx := c.Request.Context()
	for _, word := range top50 {
		err = h.repo.RecordWordFrequency(ctx, word.Word, word.Freq)
		if err != nil {
			log.Printf("ошибка сохранения данных в базе %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "внутренняя ошибка сервера",
			})
			return
		}
	}

	uploadTime := time.Since(startUploadFileTime).Seconds()

	fileName := fileHeader.Filename
	err = h.repo.RecordFileUpload(ctx, fileName, uploadTime)
	if err != nil {
		log.Printf("ошибка записи информации о загрузке файла %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Words": top50,
	})
}
