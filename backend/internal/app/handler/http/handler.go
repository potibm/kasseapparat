package http

import (
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/monitor"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	sumupRepo "github.com/potibm/kasseapparat/internal/app/repository/sumup"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
)

type Handler struct {
	repo            *sqliteRepo.Repository
	sumupRepository sumupRepo.RepositoryInterface
	purchaseService purchaseService.Service
	monitor         monitor.Poller
	mailer          mailer.Mailer
	version         string
	decimalPlaces   int32
	paymentMethods  map[string]string
}

func NewHandler(repo *sqliteRepo.Repository, sumupRepository sumupRepo.RepositoryInterface, purchaseService purchaseService.Service, monitor monitor.Poller, mailer mailer.Mailer, version string, decimalPlaces int32, paymentMethods map[string]string) *Handler {
	return &Handler{repo: repo, sumupRepository: sumupRepository, purchaseService: purchaseService, monitor: monitor, mailer: mailer, version: version, decimalPlaces: decimalPlaces, paymentMethods: paymentMethods}
}
