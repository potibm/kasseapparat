package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
)

type StatusPublisher interface {
	PushUpdate(purchaseID uuid.UUID, status models.PurchaseStatus)
}

var _ HandlerInterface = (*Handler)(nil)

type HandlerInterface interface {
	HandleTransactionWebSocket(c *gin.Context)
}

type WebsocketPublisher struct{}

func (w *WebsocketPublisher) PushUpdate(purchaseID uuid.UUID, status models.PurchaseStatus) {
	PushUpdate(purchaseID, status)
}

type Handler struct {
	sumupRepository  sumupRepo.RepositoryInterface
	sqliteRepository sqliteRepository
	purchaseService  purchaseService.Service
	upgrader         websocket.Upgrader
}

type sqliteRepository interface {
	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
}

func NewHandler(sqliteRepository sqliteRepository, sumupRepository sumupRepo.RepositoryInterface, purchaseService purchaseService.Service, corsAllowOrigins *config.CorsAllowOriginsConfig) *Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: makeCheckOrigin(corsAllowOrigins),
	}

	return &Handler{sqliteRepository: sqliteRepository, sumupRepository: sumupRepository, purchaseService: purchaseService, upgrader: upgrader}
}

func makeCheckOrigin(allowedOrigins *config.CorsAllowOriginsConfig) func(r *http.Request) bool {
	allowed := make(map[string]struct{}, len(*allowedOrigins))
	for _, o := range *allowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("WebSocket connection attempt from origin: %s", origin)

		if origin == "" {
			return false
		}

		_, ok := allowed[origin]

		return ok
	}
}
