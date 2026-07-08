package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	gormaudit "github.com/potibm/kasseapparat/internal/app/store/gorm"
)

type StatusPublisher interface {
	PushUpdate(purchaseID uuid.UUID, status models.PurchaseStatus)
}

type Handler struct {
	repo            sqliteRepo.RepositoryInterface
	sumupRepository sumupRepo.RepositoryInterface
	purchaseService purchaseService.Service
	monitor         monitor.Poller
	statusPublisher StatusPublisher
	mailer          mailer.Mailer
	config          config.Config
	decimalPlaces   int32
}

type HandlerConfig struct {
	Repo            sqliteRepo.RepositoryInterface
	SumupRepository sumupRepo.RepositoryInterface
	PurchaseService purchaseService.Service
	Monitor         monitor.Poller
	StatusPublisher StatusPublisher
	Mailer          mailer.Mailer
	AppConfig       config.Config
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		repo:            cfg.Repo,
		sumupRepository: cfg.SumupRepository,
		purchaseService: cfg.PurchaseService,
		monitor:         cfg.Monitor,
		statusPublisher: cfg.StatusPublisher,
		mailer:          cfg.Mailer,
		config:          cfg.AppConfig,
		decimalPlaces:   cfg.AppConfig.Format.Currency.FractionDigitsMax,
	}
}

func (handler *Handler) getUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get(middleware.IdentityKey)
	if !exists {
		return "", errors.New("username not found in context")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", errors.New("username in context is not a string")
	}

	return usernameStr, nil
}

func (handler *Handler) contextWithUser(c *gin.Context) *gin.Context {
	username, err := handler.getUsernameFromContext(c)
	if err != nil {
		return c
	}

	ctx := gormaudit.WithUserID(c.Request.Context(), username)
	c.Request = c.Request.WithContext(ctx)

	return c
}
