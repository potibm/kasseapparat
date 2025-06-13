package websocket

import (
	"github.com/gin-gonic/gin"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
)

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	HandleTransactionWebSocket(c *gin.Context)
}

type Handler struct {
	sumupRepository sumupRepo.RepositoryInterface
	purchaseService purchaseService.Service	
}

func NewHandler(sumupRepository sumupRepo.RepositoryInterface, purchaseService purchaseService.Service) *Handler {
	return &Handler{sumupRepository: sumupRepository, purchaseService: purchaseService}
}
