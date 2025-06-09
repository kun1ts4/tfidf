package repository

import (
	"context"
	"log"
)

func (r *Repository) CreateCollectionsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS collections (
			id TEXT NOT NULL PRIMARY KEY
		)
	`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		log.Printf("не удалось создать таблицу collections: %v", err)
		return err
	}
	return nil
}

func (r *Repository) CreateCollection(ctx context.Context, id string) error {
	query := `
		INSERT INTO collections (id)
		VALUES ($1)
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("не удалось создать коллекцию с id %s: %v", id, err)
		return err
	}
	return nil
}

func (r *Repository) GetAllCollections(ctx context.Context) ([]string, error) {
	query := `
		SELECT id FROM collections
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Printf("не удалось получить коллекции: %v", err)
		return nil, err
	}
	defer rows.Close()

	var collections []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("не удалось сканировать коллекцию: %v", err)
			return nil, err
		}
		collections = append(collections, id)
	}

	return collections, nil
}

func (r *Repository) AddFileToCollection(ctx context.Context, collectionId, fileId string) error {
	query := `
		UPDATE documents
		SET collections = array_append(collections, $1)
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, collectionId, fileId)
	if err != nil {
		log.Printf("не удалось добавить файл с id %s в коллекцию %s: %v", fileId, collectionId, err)
		return err
	}
	return nil
}

func (r *Repository) GetFilesByCollectionId(ctx context.Context, collectionId string) ([]string, error) {
	query := `
		SELECT id FROM documents
		WHERE $1 = ANY(collections)
	`

	rows, err := r.pool.Query(ctx, query, collectionId)
	if err != nil {
		log.Printf("не удалось получить файлы из коллекции %s: %v", collectionId, err)
		return nil, err
	}
	defer rows.Close()

	var fileIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("не удалось сканировать файл: %v", err)
			return nil, err
		}
		fileIds = append(fileIds, id)
	}

	return fileIds, nil
}

func (r *Repository) RemoveFileFromCollection(ctx context.Context, collectionId, fileId string) error {
	query := `
		UPDATE documents
		SET collections = array_remove(collections, $1)
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, collectionId, fileId)
	if err != nil {
		log.Printf("не удалось удалить файл с id %s из коллекции %s: %v", fileId, collectionId, err)
		return err
	}
	return nil
}

func (r *Repository) DeleteCollection(ctx context.Context, collectionId string) error {
	query := `
		DELETE FROM collections
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, collectionId)
	if err != nil {
		log.Printf("не удалось удалить коллекцию с id %s: %v", collectionId, err)
		return err
	}

	query = `		
		UPDATE documents
		SET collections = array_remove(collections, $1)
		WHERE $1 = ANY(collections)
	`

	_, err = r.pool.Exec(ctx, query, collectionId)
	if err != nil {
		log.Printf("не удалось удалить файлы из коллекции %s: %v", collectionId, err)
		return err
	}

	return nil
}
