package repository

import (
	"context"
	"log"
	"tfidf/internal/model"
)

func (r *Repository) CreateFileTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS documents(
		    id TEXT PRIMARY KEY,
		    file_name TEXT NOT NULL,
		    author_id INT NOT NULL,
		    collections TEXT[],
		    time_processed NUMERIC(10,3),
		    upload_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("не удалось создать таблицу file_uploads: %v", err)
		return err
	}
	return nil
}

func (r *Repository) SaveFileInfo(ctx context.Context, doc model.Document) error {
	query := `
        INSERT INTO documents (file_name, author_id, collections, time_processed)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.pool.Exec(ctx, query, doc.Name, doc.AuthorId, doc.Collections, doc.TimeProcessed)
	if err != nil {
		return err
	}

	return nil
}

//save document as a file
