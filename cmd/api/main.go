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

	"tfidf/internal/service"
)

// @title         TF-IDF API
// @description   Это простой сервис, который помогает анализировать текстовые документы с помощью метода TF-IDF
// @BasePath      /
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

	r.Static("/styles", "./web/templates")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/status", handler.Status)
	r.GET("/version", handler.Version)
	r.GET("/", handler.ShowForm)
	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.Login)

	r.GET("/metrics", h.GetMetrics)

	authorized := r.Group("/")
	authorized.Use(h.Auth)
	{
		authorized.POST("/upload", h.UploadFile)
	}

	return r
}
