package handler

import (
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
