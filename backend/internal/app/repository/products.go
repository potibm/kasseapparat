package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/potibm/die-kassa/internal/app/models"
	"github.com/potibm/die-kassa/internal/app/utils"
)

func GetProducts(c *gin.Context) {
	db := utils.ConnectToDatabase()

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	db := utils.ConnectToDatabase()

	var product models.Product
	if err := db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, product)
}
