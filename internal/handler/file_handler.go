package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"tfidf/internal/model"
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
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка чтения файла",
		})
		return
	}

	top50, err := service.ProcessFile(content, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
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

	timeProcessed := time.Since(startUploadFileTime).Seconds()

	fileName := fileHeader.Filename

	idValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "не удалось получить userID"})
		return
	}
	authorID, ok := idValue.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	document := model.Document{
		Id:            generateUUID(),
		Name:          fileName,
		AuthorId:      authorID,
		Collections:   nil, //TODO COLLECTIONS
		TimeProcessed: timeProcessed,
	}

	err = service.SaveFile(file, document.Id)
	if err != nil {
		log.Printf("ошибка сохранения файла на диск %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	err = h.repo.SaveFileInfo(ctx, document)

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

func generateUUID() string {
	id := uuid.New().String()
	return id
}
