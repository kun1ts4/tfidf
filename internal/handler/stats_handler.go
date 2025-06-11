package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tfidf/internal/parser"
	"tfidf/internal/service"
)

// GetCollectionStats returns the top 50 words with the highest IDF for a collection of documents.
// @Summary Get collection statistics
// @Description Returns the top 50 words with the highest IDF for documents in the collection
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Param collection_id path string true "Collection ID"
// @Success 200 {object} map[string][]model.Word
// @Failure 500 {object} model.ErrorResponse
// @Router /collection/{collection_id}/statistics [get]
func (h *Handler) GetCollectionStats(c *gin.Context) {
	fileIDs, err := h.repo.GetFilesByCollectionId(c, c.Param("collection_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messages": "не удалось получить файлы из коллекции"})
		return
	}

	documents := make([][]string, 0)
	for _, id := range fileIDs {
		document, err := service.GetFile(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"messages": "ошибка получения документа"})
		}
		words := parser.ExtractWords(document)
		documents = append(documents, words)
	}

	tfidf, _ := service.CalculateTFIDF(documents, 0)
	tfidfTop := service.TopIDFRange(tfidf, 0, 50)
	c.JSON(http.StatusOK, gin.H{
		"statistics": tfidfTop,
	})
}

// GetDocumentStats returns the top 50 words with the highest TF-IDF for a specific document.
// @Summary Get document statistics
// @Description Returns the top 50 words with the highest TF-IDF for the specified document
// @Tags files
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string][]model.Word
// @Failure 500 {object} model.ErrorResponse
// @Router /documents/{id}/statistics [get]
func (h *Handler) GetDocumentStats(c *gin.Context) {
	fileID := c.Param("id")
	collections, err := h.repo.GetFileCollections(c, fileID)
	if err != nil {
		log.Printf("ошибка получения коллекций %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"messages": "ошибка получения коллекций"})
		return
	}

	collection := collections[0]

	fileIDs, err := h.repo.GetFilesByCollectionId(c, collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messages": "не удалось получить файлы из коллекции"})
		return
	}

	curFileIndex := 0
	for i, id := range fileIDs {
		if fileID == id {
			curFileIndex = i
		}
	}

	documents := make([][]string, 0)
	for _, id := range fileIDs {
		document, err := service.GetFile(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"messages": "ошибка получения документа"})
		}
		words := parser.ExtractWords(document)
		documents = append(documents, words)
	}

	_, tfidf := service.CalculateTFIDF(documents, curFileIndex)
	tfidfTop := service.TopIDFRange(tfidf, 0, 50)
	c.JSON(http.StatusOK, gin.H{
		"statistics": tfidfTop,
	})
}
