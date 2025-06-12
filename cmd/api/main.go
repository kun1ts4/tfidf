package main

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"tfidf/internal/config"
	"tfidf/internal/db"
	"tfidf/internal/handler"
	"tfidf/internal/repository"
	"time"

	"tfidf/internal/service"
)

// @title         TF-IDF API
// @description   Это простой сервис, который помогает анализировать текстовые документы с помощью метода TF-IDF
// @host		  37.9.53.154:8080
// @BasePath      /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
	cfg, err := config.Initialize()
	if err != nil {
		log.Fatalf("не удалось инициализировать конфигурацию: %v", err)
	}

	pgxPool, err := db.Initialize(context.Background(), cfg)
	if err != nil {
		log.Fatalf("не удалось инициализировать базу данных: %v", err)
	}
	defer pgxPool.Close()

	//ожидание создания таблиц
	time.Sleep(2 * time.Second)

	repo := repository.NewRepository(pgxPool)
	if err := db.InitializeTables(repo); err != nil {
		log.Fatalf("не удалось инициализировать таблицы: %v", err)
	}

	tokenService := service.NewTokenService("secret")
	h := handler.NewHandler(repo, tokenService)

	r := SetupRouter(h)

	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/status", handler.Status)
	r.GET("/version", handler.Version)

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
		//TODO authorized.GET("/documents/:id/huffman", h.GetHuffman)

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
