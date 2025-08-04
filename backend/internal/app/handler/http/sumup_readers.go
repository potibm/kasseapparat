package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
)

type SumupReaderReponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	DeviceIdentifier string    `json:"deviceIdentifier"`
	DeviceModel      string    `json:"deviceModel"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (handler *Handler) GetSumupReaders(c *gin.Context) {
	readers, err := handler.sumupRepository.GetReaders()
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("Failed to retrieve readers").WithCause(err))

		return
	}

	c.Header("X-Total-Count", strconv.Itoa(len(readers)))
	c.JSON(http.StatusOK, toSumupReaderResponses(readers))
}

func (handler *Handler) GetSumupReaderByID(c *gin.Context) {
	id := c.Param("id")

	reader, err := handler.sumupRepository.GetReader(id)
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("Failed to retrieve reader").WithCause(err))

		return
	}

	if reader == nil {
		_ = c.Error(NotFound.WithMsg("Reader not found"))

		return
	}

	c.JSON(http.StatusOK, toSumupReaderResponse(*reader))
}

func (handler *Handler) CreateSumupReader(c *gin.Context) {
	var request struct {
		PairingCode string `binding:"required" json:"pairingCode"`
		Name        string `json:"name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(InvalidRequest.WithCause(err))

		return
	}

	reader, err := handler.sumupRepository.CreateReader(strings.ToUpper(request.PairingCode), request.Name)
	if err != nil {
		_ = c.Error(InternalServerError.WithMsg("Failed to create reader").WithCause(err))

		return
	}

	c.JSON(http.StatusCreated, toSumupReaderResponse(*reader))
}

func (handler *Handler) DeleteSumupReader(c *gin.Context) {
	id := c.Param("id")

	if err := handler.sumupRepository.DeleteReader(id); err != nil {
		_ = c.Error(InternalServerError.WithMsg("Failed to delete reader").WithCause(err))

		return
	}

	c.Status(http.StatusNoContent)
}

func toSumupReaderResponse(c sumup.Reader) SumupReaderReponse {
	return SumupReaderReponse{
		ID:               c.ID,
		Name:             c.Name,
		Status:           c.Status,
		DeviceIdentifier: c.DeviceIdentifier,
		DeviceModel:      c.DeviceModel,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func toSumupReaderResponses(readers []sumup.Reader) []SumupReaderReponse {
	responses := make([]SumupReaderReponse, len(readers))
	for i, reader := range readers {
		responses[i] = toSumupReaderResponse(reader)
	}

	return responses
}
