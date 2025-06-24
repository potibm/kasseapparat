package http

import (
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
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

func NewHandler(config HandlerConfig) *Handler {
	return &Handler{
		repo:            config.Repo,
		sumupRepository: config.SumupRepository,
		purchaseService: config.PurchaseService,
		monitor:         config.Monitor,
		statusPublisher: config.StatusPublisher,
		mailer:          config.Mailer,
		config:          config.AppConfig,
		decimalPlaces:   int32(config.AppConfig.FormatConfig.FractionDigitsMax),
	}
}
