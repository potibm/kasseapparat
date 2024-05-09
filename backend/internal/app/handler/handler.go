package handler

import "github.com/potibm/kasseapparat/internal/app/repository"

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}
