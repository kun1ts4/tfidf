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

// UploadFile godoc
// @Summary Upload file
// @Description Uploads a file to the server
// @Tags files
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /upload [post]
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

	timeProcessed := time.Since(startUploadFileTime).Seconds()

	fileName := fileHeader.Filename

	authorID, err := GetUserID(c)
	if err != nil {
		log.Printf("не удалось получить UserID %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "для загрузки файлов необходимо войти в аккаунт",
		})
	}

	document := model.Document{
		Id:            generateUUID(),
		Name:          fileName,
		AuthorId:      authorID,
		TimeProcessed: timeProcessed,
	}

	err = service.SaveFile(content, document.Id)
	if err != nil {
		log.Printf("ошибка сохранения файла на диск %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	err = h.repo.SaveFileInfo(c.Request.Context(), document)

	if err != nil {
		log.Printf("ошибка записи информации о загрузке файла %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "файл загружен успешно",
	})
}

// generateUUID creates a new unique identifier
func generateUUID() string {
	id := uuid.New().String()
	return id
}

// GetUserDocuments godoc
// @Summary Get user documents
// @Description Returns list of document IDs belonging to the authenticated user
// @Tags files
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} string "List of document IDs"
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /documents [get]
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

// GetDocumentById godoc
// @Summary Get document by ID
// @Description Returns file content by its ID
// @Tags files
// @Security ApiKeyAuth
// @Produce text/plain
// @Param id path string true "Document ID"
// @Success 200 {string} string "File content"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /documents/{id} [get]
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

// GetStats godoc
// @Summary Get statistics
// @Description Returns system statistics (to be implemented)
// @Tags admin
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} model.Stats
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Router /stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	// TODO: Implement statistics endpoint
}
