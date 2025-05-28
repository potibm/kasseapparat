package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/repository/sumup"
)

type ReaderReponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	DeviceIdentifier string    `json:"deviceIdentifier"`
	DeviceModel      string    `json:"deviceModel"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (handler *Handler) GetSumupReaders(c *gin.Context) {

	readers, _ := handler.sumupRepository.GetReaders()

	c.Header("X-Total-Count", strconv.Itoa(len(readers)))
	c.JSON(http.StatusOK, toReaderResponses(readers))
}

func (handler *Handler) GetSumupReaderByID(c *gin.Context) {
	id := c.Param("id")

	reader, err := handler.sumupRepository.GetReader(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reader"})
		return
	}

	if reader == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	c.JSON(http.StatusOK, toReaderResponse(*reader))
}

func toReaderResponse(c sumup.Reader) ReaderReponse {
	return ReaderReponse{
		ID:               c.ID,
		Name:             c.Name,
		Status:           c.Status,
		DeviceIdentifier: c.DeviceIdentifier,
		DeviceModel:      c.DeviceModel,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func toReaderResponses(readers []sumup.Reader) []ReaderReponse {
	responses := make([]ReaderReponse, len(readers))
	for i, reader := range readers {
		responses[i] = toReaderResponse(reader)
	}
	return responses
}
