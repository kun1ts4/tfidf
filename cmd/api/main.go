package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"tfidf/internal/config"
	"tfidf/internal/handler"
	"tfidf/internal/repository"
	"tfidf/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	pgxConn, err := db.NewConnection(context.Background(), cfg)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err.Error())
	}
	defer pgxConn.Close()

	repo := repository.NewRepository(pgxConn)
	err = repo.CreateMetricsTables(context.Background())
	if err != nil {
		log.Fatalf("не удалось создать таблицы: %v", err)
	}

	r := gin.Default()

	r.Static("/styles", "./web/templates")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/status", handler.Status)
	r.GET("/version", handler.Version)
	r.GET("/", handler.ShowForm)

	uploadHandler := handler.NewHandler(repo)
	r.POST("/upload", uploadHandler.UploadFile)

	metricsHandler := handler.NewHandler(repo)
	r.GET("/metrics", metricsHandler.GetMetrics)

	err = r.Run(":" + cfg.App.Port)
	if err != nil {
		return
	}
}
