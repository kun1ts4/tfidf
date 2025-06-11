package repository

import (
	"context"
	"log"
	"strings"
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
        INSERT INTO documents (id, file_name, author_id, collections, time_processed)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.pool.Exec(ctx, query, doc.Id, doc.Name, doc.AuthorId, doc.Collections, doc.TimeProcessed)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetFilesByAuthorId(ctx context.Context, authorId int) ([]model.Document, error) {
	query := `SELECT id, file_name, author_id, collections, time_processed FROM documents WHERE author_id = $1`
	rows, err := r.pool.Query(ctx, query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []model.Document
	for rows.Next() {
		var document model.Document
		err := rows.Scan(&document.Id, &document.Name, &document.AuthorId, &document.Collections, &document.TimeProcessed)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) DeleteDocument(ctx context.Context, id string) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFileCollections(ctx context.Context, id string) ([]string, error) {
	query := `SELECT array_to_string(collections, ',') FROM documents WHERE id = $1`
	var s string
	err := r.pool.QueryRow(ctx, query, id).Scan(&s)
	if err != nil {
		return nil, err
	}

	return strings.Split(s, ","), nil
}
