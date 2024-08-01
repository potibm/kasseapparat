package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetProductStats(c *gin.Context) {
	products, err := handler.repo.GetProductStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(len(products)))
	c.JSON(http.StatusOK, products)
}
