package repository

import (
	"context"
	"tfidf/internal/model"
)

func (r *Repository) CreateUserTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL)
    `

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	var user model.User
	err := r.pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *Repository) CheckUserPassword(ctx context.Context, user model.User) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2`
	var count int
	err := r.pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ChangeUserPassword(ctx context.Context, username, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE username = $2`
	_, err := r.pool.Exec(ctx, query, newPassword, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = $1`
	_, err := r.pool.Exec(ctx, query, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetNumberOfUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
