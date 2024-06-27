package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type PurchaseListItemRequest struct {
	ID       int `form:"ID" binding:"required"`
	AttendedGuests uint `form:"attendedGuests" binding:"required"`
}

type PurchaseCartRequest struct {
	ID       int `form:"ID" binding:"required"`
	Quantity int `form:"quantity" binding:"required"`
	ListItems	[]PurchaseListItemRequest `form:"listItems" binding:"required,dive"`
}

type PurchaseRequest struct {
	TotalPrice float64               `form:"totalPrice" binding:"numeric"`
	Cart       []PurchaseCartRequest `form:"cart" binding:"required,dive"`
}

func (handler *Handler) OptionsPurchases(c *gin.Context) {

	c.JSON(http.StatusOK, nil)
}

func (handler *Handler) DeletePurchase(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	handler.repo.DeletePurchaseByID(id, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted"})
}

func (handler *Handler) PostPurchases(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var purchase models.Purchase
	var updatedListEntries []models.ListEntry = make([]models.ListEntry, 0)

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
		purchase.CreatedByID = &executingUserObj.ID

		for j := 0; j < len(purchaseRequest.Cart[i].ListItems); j++ {
			var listEntry *models.ListEntry
			listEntry, err = handler.repo.GetFullListEntryByID(purchaseRequest.Cart[i].ListItems[j].ID)
			if err != nil || listEntry == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "List item not found."})
				return
			}
			
			if (listEntry.AttendedGuests != 0) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "List item has already been attended."})
				return
			}

			if (listEntry.AdditionalGuests+1 < purchaseRequest.Cart[i].ListItems[j].AttendedGuests) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Additional guests exceed available guests."})
				return
			}

			if (listEntry.List.ProductID != uint(id)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "List item does not belong to product."})
				return
			}

			listEntry.AttendedGuests = purchaseRequest.Cart[i].ListItems[j].AttendedGuests
			updatedListEntries = append(updatedListEntries, *listEntry)
		}

	}
	// check that total price is correct
	if calculatedTotalPrice != purchaseRequest.TotalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Total price does not match."})
		return
	}

	purchase.TotalPrice = calculatedTotalPrice

	// update the list of listEntries
	for i := 0; i < len(updatedListEntries); i++ {
		_, err := handler.repo.UpdateListEntryByID(int(updatedListEntries[i].ID), updatedListEntries[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
			return
		}
	}

	purchase, err = handler.repo.StorePurchases(purchase)
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
