package main

import (
	"context"
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

	r := handler.SetupRouter(h)

	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
