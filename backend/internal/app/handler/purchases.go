package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/die-kassa/internal/app/models"
)

type PurchaseCartRequest struct {
	ID       int `form:"ID" binding:"required"`
	Quantity int `form:"quantity" binding:"required"`
}

type PurchaseRequest struct {
	TotalPrice float64               `form:"totalPrice" binding:"numeric"`
	Cart       []PurchaseCartRequest `form:"cart" binding:"required,dive"`
}

func (handler *Handler) OptionsPurchases(c *gin.Context) {

	c.JSON(http.StatusOK, nil)
}

func (handler *Handler) DeletePurchases(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	handler.repo.DeletePurchases(id)

	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted"})
}

func (handler *Handler) PostPurchases(c *gin.Context) {

	var purchase models.Purchase

	var purchaseRequest PurchaseRequest
	if c.ShouldBind(&purchaseRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	calculatedTotalPrice := 0.0
	for i := 0; i < len(purchaseRequest.Cart); i++ {
		id := purchaseRequest.Cart[i].ID
		quantity := purchaseRequest.Cart[i].Quantity

		product, err := handler.repo.GetProductByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Product not found."})
			return
		}
		calculatedPurchaseItemPrice := product.Price * float64(quantity)
		calculatedTotalPrice += calculatedPurchaseItemPrice

		purchaseItem := models.PurchaseItem{
			Product:    *product,
			Quantity:   purchaseRequest.Cart[i].Quantity,
			Price:      product.Price,
			TotalPrice: calculatedPurchaseItemPrice,
		}
		purchase.PurchaseItems = append(purchase.PurchaseItems, purchaseItem)
	}
	// check that total price is correct
	if calculatedTotalPrice != purchaseRequest.TotalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Total price does not match."})
		return
	}

	purchase.TotalPrice = calculatedTotalPrice

	purchase, err := handler.repo.StorePurchases(purchase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase successful", "purchase": purchase})
}
