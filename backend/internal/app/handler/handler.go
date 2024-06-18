package handler

import "github.com/potibm/kasseapparat/internal/app/repository"

type Handler struct {
	repo *repository.Repository
	version string
}

func NewHandler(repo *repository.Repository, version string) *Handler {
	return &Handler{repo: repo, version: version}
}
