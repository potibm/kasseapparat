package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

type PurchaseListItemRequest struct {
	ID             int  `binding:"required" form:"ID"`
	AttendedGuests uint `binding:"required" form:"attendedGuests"`
}

type PurchaseCartRequest struct {
	ID        int                       `binding:"required"      form:"ID"`
	Quantity  int                       `binding:"required"      form:"quantity"`
	ListItems []PurchaseListItemRequest `binding:"required,dive" form:"listItems"`
}

type PurchaseRequest struct {
	TotalNetPrice   decimal.Decimal       `binding:"required"      form:"totalNetPrice"`
	TotalGrossPrice decimal.Decimal       `binding:"required"      form:"totalGrossPrice"`
	Cart            []PurchaseCartRequest `binding:"required,dive" form:"cart"`
}

func (handler *Handler) DeletePurchase(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	handler.repo.DeletePurchaseByID(id, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted"})
}

func (handler *Handler) PostPurchases(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	var purchase models.Purchase

	updatedListEntries := make([]models.Guest, 0)

	var purchaseRequest PurchaseRequest
	if err := c.ShouldBind(&purchaseRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	calculatedTotalNetPrice := decimal.NewFromInt(0)
	calculatedTotalGrossPrice := decimal.NewFromInt(0)

	for i := range len(purchaseRequest.Cart) {
		id := purchaseRequest.Cart[i].ID
		quantity := purchaseRequest.Cart[i].Quantity

		product, err := handler.repo.GetProductByID(id)
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Product not found"))

			return
		}

		calculatedPurchaseItemNetPrice := product.NetPrice.Mul(decimal.NewFromInt(int64(quantity)))
		calculatedTotalNetPrice = calculatedTotalNetPrice.Add(calculatedPurchaseItemNetPrice)

		calculatedPurchaseItemGrossPrice := product.GrossPrice().Mul(decimal.NewFromInt(int64(quantity)))
		calculatedTotalGrossPrice = calculatedTotalGrossPrice.Add(calculatedPurchaseItemGrossPrice)

		purchaseItem := models.PurchaseItem{
			Product:  *product,
			Quantity: purchaseRequest.Cart[i].Quantity,
			NetPrice: product.NetPrice,
			VATRate:  product.VATRate,
		}
		purchase.PurchaseItems = append(purchase.PurchaseItems, purchaseItem)
		purchase.CreatedByID = &executingUserObj.ID

		for j := range len(purchaseRequest.Cart[i].ListItems) {
			var listEntry *models.Guest

			listEntry, err = handler.repo.GetFullGuestByID(purchaseRequest.Cart[i].ListItems[j].ID)
			if err != nil || listEntry == nil {
				_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "List item not found"))

				return
			}

			if listEntry.AttendedGuests != 0 {
				_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "List item has already been attended"))

				return
			}

			if listEntry.AdditionalGuests+1 < purchaseRequest.Cart[i].ListItems[j].AttendedGuests {
				_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Additional guests exceed available guests"))

				return
			}

			if listEntry.Guestlist.ProductID != uint(id) {
				_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "List item does not belong to product"))

				return
			}

			listEntry.AttendedGuests = purchaseRequest.Cart[i].ListItems[j].AttendedGuests
			listEntry.MarkAsArrived()

			updatedListEntries = append(updatedListEntries, *listEntry)
		}
	}
	// check that total price is correct
	if !calculatedTotalNetPrice.Equal(purchaseRequest.TotalNetPrice) {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Total net price does not match"))

		return
	}

	if !calculatedTotalGrossPrice.Equal(purchaseRequest.TotalGrossPrice) {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Total gross price does not match"))

		return
	}

	purchase.TotalNetPrice = calculatedTotalNetPrice
	purchase.TotalGrossPrice = calculatedTotalGrossPrice

	purchase, err = handler.repo.StorePurchases(purchase)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	// update the list of listEntries
	for i := range updatedListEntries {
		updatedListEntry := updatedListEntries[i]
		updatedListEntry.PurchaseID = &purchase.ID

		_, err := handler.repo.UpdateGuestByID(int(updatedListEntry.ID), updatedListEntry)
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

			return
		}
	}

	for i := range updatedListEntries {
		updatedListEntry := updatedListEntries[i]
		if updatedListEntry.NotifyOnArrivalEmail != nil {
			_ = handler.mailer.SendNotificationOnArrival(*updatedListEntry.NotifyOnArrivalEmail, updatedListEntry.Name)
		}
	}

	purchasesResponse := response.ToPurchaseResponse(purchase)

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase successful", "purchase": purchasesResponse})
}

func (handler *Handler) GetPurchases(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "DESC")

	purchases, err := handler.repo.GetPurchases(end-start, start, sort, order)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	total, err := handler.repo.GetTotalPurchases()
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	purchasesRespone := response.ToPurchasesResponse(purchases)

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, purchasesRespone)
}

func (handler *Handler) GetPurchaseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	purchase, err := handler.repo.GetPurchaseByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	purchaseResponse := response.ToPurchaseResponse(*purchase)

	c.JSON(http.StatusOK, purchaseResponse)
}

func (handler *Handler) GetPurchaseStats(c *gin.Context) {
	stats, err := handler.repo.GetPurchaseStats()
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

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
