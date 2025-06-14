package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	HandleTransactionWebSocket(c *gin.Context)
}

type Handler struct {
	sumupRepository  sumupRepo.RepositoryInterface
	sqliteRepository sqliteRepository
	purchaseService  purchaseService.Service
}

type sqliteRepository interface {
	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
}

func NewHandler(sqliteRepository sqliteRepository, sumupRepository sumupRepo.RepositoryInterface, purchaseService purchaseService.Service) *Handler {
	return &Handler{sqliteRepository: sqliteRepository, sumupRepository: sumupRepository, purchaseService: purchaseService}
}
