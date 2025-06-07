package handler

import (
	"fmt"
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

	authorID, err := GetUserID(c)
	if err != nil {
		log.Printf("не удалось получить UserID %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
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

func (h *Handler) GetUserDocuments(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		log.Printf("не удалось получить UserID %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
	}
	documents, err := h.repo.GetFilesByAuthorId(c.Request.Context(), userID)
	if err != nil {
		log.Printf("ошибка получения информации о файлах %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	ids := make([]string, len(documents))
	for i, doc := range documents {
		ids[i] = doc.Id
	}

	c.JSON(http.StatusOK, ids)
}

func GetUserID(c *gin.Context) (int, error) {
	idValue, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("user ID не существует")
	}
	authorID, ok := idValue.(int)
	if !ok {
		return 0, fmt.Errorf("user ID не существует")
	}
	return authorID, nil
}

func (h *Handler) GetDocumentById(c *gin.Context) {
	fileID := c.Param("id")
	file, err := service.GetFile(fileID)
	if err != nil {
		log.Printf("ошибка получения файла %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte(file))
}

//func (h *Handler) GetStats(c *gin.Context) {
//
//}
