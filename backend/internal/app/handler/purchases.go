package handler

import (
	"math"
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

func (handler *Handler) GetLastPurchases(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	purchases, err := handler.repo.GetLastPurchases(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	totalRows, err := handler.repo.GetTotalPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	var nextPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}
	var prevPage *int
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	type Pagination struct {
		TotalRecords int64 `json:"total_records"`
		CurrentPage  int   `json:"current_page"`
		TotalPages   int   `json:"total_pages"`
		NextPage     *int  `json:"next_page"`
		PrevPage     *int  `json:"prev_page"`
	}

	pagination := Pagination{
		TotalRecords: totalRows,
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     nextPage,
		PrevPage:     prevPage,
	}

	c.JSON(http.StatusOK, gin.H{"data": purchases, "pagination": pagination})
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
