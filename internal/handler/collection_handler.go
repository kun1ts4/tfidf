package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tfidf/internal/service"
)

// ListCollections godoc
// @Summary List all collections
// @Description Returns a list of all collections
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string][]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection [get]
func (h *Handler) ListCollections(c *gin.Context) {
	collections, err := h.repo.GetAllCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить коллекции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"collections": collections})
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Creates a new collection and returns its ID
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection [post]
func (h *Handler) CreateCollection(c *gin.Context) {
	id := service.GenerateUUID()

	err := h.repo.CreateCollection(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать коллекцию"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "коллекция успешно создана", "id": id})
}

// AddDocumentToCollection godoc
// @Summary Add a document to a collection
// @Description Adds a document to the specified collection
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Param collection_id path string true "Collection ID"
// @Param document_id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection/{collection_id}/{document_id} [post]
func (h *Handler) AddDocumentToCollection(c *gin.Context) {
	collectionID := c.Param("collection_id")
	documentID := c.Param("document_id")

	err := h.repo.AddFileToCollection(c, collectionID, documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось добавить документ в коллекцию"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "документ успешно добавлен в коллекцию"})
}

// DeleteCollection godoc
// @Summary Delete a collection
// @Description Deletes the specified collection
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Param collection_id path string true "Collection ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection/{collection_id} [delete]
func (h *Handler) DeleteCollection(c *gin.Context) {
	err := h.repo.DeleteCollection(c, c.Param("collection_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить коллекцию"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "коллекция удалена"})
}

// ListCollectionDocuments godoc
// @Summary List documents in a collection
// @Description Returns a list of documents in the specified collection
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Param collection_id path string true "Collection ID"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection/{collection_id} [get]
func (h *Handler) ListCollectionDocuments(c *gin.Context) {
	files, err := h.repo.GetFilesByCollectionId(c, c.Param("collection_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить файлы из коллекции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// DeleteDocumentFromCollection godoc
// @Summary Remove a document from a collection
// @Description Removes a document from the specified collection
// @Tags collections
// @Security ApiKeyAuth
// @Produce json
// @Param collection_id path string true "Collection ID"
// @Param document_id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} model.ErrorResponse
// @Router /collection/{collection_id}/{document_id} [delete]
func (h *Handler) DeleteDocumentFromCollection(c *gin.Context) {
	collectionID := c.Param("collection_id")
	documentID := c.Param("document_id")

	err := h.repo.RemoveFileFromCollection(c, collectionID, documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить документ из коллекции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "документ успешно удален из коллекции"})
}
