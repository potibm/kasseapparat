package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
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

func (handler *Handler) GetPurchases(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "DESC")

	products, err := handler.repo.GetPurchases(end-start, start, sort, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := handler.repo.GetTotalPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, products)
}

func (handler *Handler) GetPurchaseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	purchase, err := handler.repo.GetPurchaseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (handler *Handler) GetPurchaseStats(c *gin.Context) {

	stats, err := handler.repo.GetPurchaseStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	// iterate over all stats and calculate the total quantity
	totalQuantity := 0
	for _, stat := range stats {
		totalQuantity += stat.Quantity
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"stats": stats, "totalQuantity": totalQuantity})
}
