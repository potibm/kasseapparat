package handler

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

type Handler struct {
	repo    *repository.Repository
	mailer  mailer.Mailer
	version string
}

func NewHandler(repo *repository.Repository, mailer mailer.Mailer, version string) *Handler {
	return &Handler{repo: repo, mailer: mailer, version: version}
}

func queryArrayInt(c *gin.Context, field string) []int {
	idStrings := c.QueryArray(field)

	var ids []int

	for _, s := range idStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("Error converting %s to int: %v", s, err)
		}

		ids = append(ids, id)
	}

	return ids
}
