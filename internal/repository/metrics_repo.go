package repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"tfidf/internal/model"
)

func (r *Repository) CreateMetricsTables(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS word_frequencies(
		    id SERIAL PRIMARY KEY,
		    word TEXT NOT NULL,
		    freq DOUBLE PRECISION NOT NULL
		);
	`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("не удалось создать таблицу word_frequencies: %v", err)
		return err
	}

	query = `
		CREATE TABLE IF NOT EXISTS file_uploads(
		    id SERIAL PRIMARY KEY,
		    file_name TEXT NOT NULL,
		    time_processed NUMERIC(10,3),
		    upload_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`

	_, err = r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("не удалось создать таблицу file_uploads: %v", err)
		return err
	}

	return nil
}

func (r *Repository) RecordWordFrequency(ctx context.Context, word string, freq int) error {
	query := `
		INSERT INTO word_frequencies (word, freq)
		VALUES ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, word, freq)
	if err != nil {
		log.Printf("не удалось сохранить частоту слова %s: %v", word, err)
		return err
	}
	return nil
}

func (r *Repository) RecordFileUpload(ctx context.Context, fileName string, timeProcessed float64) error {
	query := `
		INSERT INTO file_uploads (file_name, time_processed)
		VALUES ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, fileName, timeProcessed)
	if err != nil {
		log.Printf("не удалось сохранить информацию о загрузке файла %s: %v", fileName, err)
		return err
	}
	return nil
}

func (r *Repository) GetPeakUploadTime(ctx context.Context) (string, error) {
	query := `
		SELECT TO_CHAR(upload_time, 'HH24:00') AS upload_hour 
		FROM file_uploads 
		GROUP BY upload_hour 
		ORDER BY COUNT(*) DESC 
		LIMIT 1`

	var uploadHour string

	err := r.pool.QueryRow(ctx, query).Scan(&uploadHour)
	if err != nil {
		return "", fmt.Errorf("не удалось получить время загрузки: %v", err)
	}

	return uploadHour, nil
}

func (r *Repository) GetTopFreqWords(ctx context.Context, limit int) ([]model.WordTFIDF, error) {
	query := `SELECT word, freq FROM word_frequencies ORDER BY freq DESC LIMIT $1`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("не удалось запросить самые популярные слова: %v", err)
	}
	defer rows.Close()

	var result []model.WordTFIDF

	for rows.Next() {
		var word string
		var freq int
		if err := rows.Scan(&word, &freq); err != nil {
			return nil, fmt.Errorf("не удалось отсканировать строку: %v", err)
		}
		result = append(result, model.WordTFIDF{
			Word: word,
			Freq: freq,
		})
	}

	return result, nil
}

func (r *Repository) GetFilesProcessed(ctx context.Context) (float64, error) {
	query := `
		SELECT COUNT(*) FROM file_uploads
	`

	var count float64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество обработанных файлов: %v", err)
	}

	return count, nil
}

func (r *Repository) GetMinTimeProcessed(ctx context.Context) (float64, error) {
	query := `
		SELECT MIN(time_processed) FROM file_uploads
	`

	var minTime float64
	err := r.pool.QueryRow(ctx, query).Scan(&minTime)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить минимальное время обработки: %v", err)
	}

	return minTime, nil
}

func (r *Repository) GetMaxTimeProcessed(ctx context.Context) (float64, error) {
	query := `
		SELECT MAX(time_processed) FROM file_uploads
	`

	var maxTime float64
	err := r.pool.QueryRow(ctx, query).Scan(&maxTime)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить максимальное время обработки: %v", err)
	}

	return maxTime, nil
}

func (r *Repository) GetAvgTimeProcessed(ctx context.Context) (float64, error) {
	query := `
		SELECT AVG(time_processed) FROM file_uploads
	`

	var avgTime float64
	err := r.pool.QueryRow(ctx, query).Scan(&avgTime)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить среднее время обработки: %v", err)
	}

	roundedAvgTime := math.Round(avgTime*1000) / 1000

	return roundedAvgTime, nil
}

func (r *Repository) GetLatestFileProcessedTimestamp(ctx context.Context) (string, error) {
	query := `
		SELECT TO_CHAR(upload_time, 'YYYY-MM-DD HH24:MI:SS') FROM file_uploads
		ORDER BY upload_time DESC LIMIT 1
	`

	var latestTimestamp string
	err := r.pool.QueryRow(ctx, query).Scan(&latestTimestamp)
	if err != nil {
		return "", fmt.Errorf("не удалось получить время последней обработки файла: %v", err)
	}

	return latestTimestamp, nil
}
