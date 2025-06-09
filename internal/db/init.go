package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tfidf/internal/config"
	"tfidf/internal/repository"
	"tfidf/pkg/postgres"
)

func Initialize(ctx context.Context, cfg *config.AppConfig) (*pgxpool.Pool, error) {
	pool, err := postgres.NewConnection(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	return pool, nil
}

func InitializeTables(repo *repository.Repository) error {
	ctx := context.Background()
	if err := repo.CreateMetricsTables(ctx); err != nil {
		return fmt.Errorf("не удалось создать таблицы метрик: %v", err)
	}
	if err := repo.CreateUserTable(ctx); err != nil {
		return fmt.Errorf("не удалось создать таблицу пользователей: %v", err)
	}
	if err := repo.CreateFileTable(ctx); err != nil {
		return fmt.Errorf("не удалось создать таблицу документов: %v", err)
	}
	if err := repo.CreateCollectionsTable(ctx); err != nil {
		return fmt.Errorf("не удалось создать таблицу коллекций: %v", err)
	}
	return nil
}
