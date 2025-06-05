package handler

import (
	_ "tfidf/cmd/api/docs"
	"tfidf/internal/repository"
	"tfidf/internal/service"
)

type Handler struct {
	repo         *repository.Repository
	tokenService *service.TokenService
}

func NewHandler(repo *repository.Repository, tokenService *service.TokenService) *Handler {
	return &Handler{
		repo:         repo,
		tokenService: tokenService,
	}
}
