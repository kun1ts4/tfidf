package handler

import (
	"github.com/gin-gonic/gin"
	"tfidf/internal/service"
)

// GetHuffman godoc
// @Summary Get Huffman encoding for a document
// @Description Returns the Huffman code and code tree for the specified document
// @Tags documents
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} model.ErrorResponse
// @Router /documents/{id}/huffman [get]
func (h *Handler) GetHuffman(c *gin.Context) {
	documentID := c.Param("id")
	documentString, err := service.GetFile(documentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "не найден документ"})
		return
	}
	document := []byte(documentString)

	code, codeTree, err := service.HuffmanEncode(document)
	if err != nil {
		c.JSON(500, gin.H{"error": "ошибка при кодировании Хаффмана"})
		return
	}

	c.JSON(200, gin.H{
		"code": code,
		"tree": codeTree,
	})
}
