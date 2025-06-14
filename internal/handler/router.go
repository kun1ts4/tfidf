package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/status", Status)
	r.GET("/version", Version)

	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.Login)

	r.GET("/metrics", h.GetMetrics)

	authorized := r.Group("/")
	authorized.Use(h.Auth)
	{
		authorized.POST("/upload", h.UploadFile)
		authorized.GET("/documents", h.GetUserDocuments)
		authorized.GET("/documents/:id", h.GetDocumentById)
		authorized.GET("/documents/:id/statistics", h.GetDocumentStats)
		authorized.DELETE("/documents/:id", h.DeleteDocument)
		authorized.GET("/documents/:id/huffman", h.GetHuffman)

		authorized.GET("/logout", h.Logout)
		authorized.PATCH("/user", h.ChangeUserPassword)
		authorized.DELETE("/user", h.DeleteUser)

		authorized.POST("/collection", h.CreateCollection)
		authorized.POST("/collection/:collection_id/:document_id", h.AddDocumentToCollection)
		authorized.GET("/collection", h.ListCollections)
		authorized.GET("/collection/:collection_id", h.ListCollectionDocuments)
		authorized.DELETE("/collection/:collection_id/:document_id", h.DeleteDocumentFromCollection)
		authorized.DELETE("/collection/:collection_id", h.DeleteCollection)
		authorized.GET("/collection/:collection_id/statistics", h.GetCollectionStats)
	}

	return r
}
