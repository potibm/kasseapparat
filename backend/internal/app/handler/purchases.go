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

	var purchaseRequest PurchaseRequest
	if err := c.ShouldBind(&purchaseRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	purchase, updatedListEntries, err := handler.createPurchase(purchaseRequest, *executingUserObj)
	if err != nil {
		_ = c.Error(err)
		return
	}

	purchase, err = handler.repo.StorePurchases(purchase)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	// update the list of listEntries
	handler.updateGuestEntries(c, updatedListEntries, purchase.ID)
	handler.sendNotifications(updatedListEntries)

	purchasesResponse := response.ToPurchaseResponse(purchase, handler.decimalPlaces)

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase successful", "purchase": purchasesResponse})
}

func (handler *Handler) createPurchase(purchaseRequest PurchaseRequest, executingUserObj models.User) (models.Purchase, []models.Guest, error) {
	var purchase models.Purchase

	updatedListEntries := make([]models.Guest, 0)

	calculatedTotalNetPrice := decimal.NewFromInt(0)
	calculatedTotalGrossPrice := decimal.NewFromInt(0)

	for _, cartItem := range purchaseRequest.Cart {
		product, err := handler.repo.GetProductByID(cartItem.ID)
		if err != nil {
			return purchase, nil, ExtendHttpErrorWithDetails(InvalidRequest, "Product not found")
		}

		calculatedTotalNetPrice = calculatedTotalNetPrice.Add(product.NetPrice.Mul(decimal.NewFromInt(int64(cartItem.Quantity))))
		calculatedTotalGrossPrice = calculatedTotalGrossPrice.Add(product.GrossPrice(handler.decimalPlaces).Mul(decimal.NewFromInt(int64(cartItem.Quantity))))

		purchase.PurchaseItems = append(purchase.PurchaseItems, models.PurchaseItem{
			Product:  *product,
			Quantity: cartItem.Quantity,
			NetPrice: product.NetPrice,
			VATRate:  product.VATRate,
		})

		listEntries, err := handler.validateAndProcessListEntries(cartItem, product.ID)
		if err != nil {
			return purchase, nil, err
		}

		updatedListEntries = append(updatedListEntries, listEntries...)
	}

	if !calculatedTotalNetPrice.Equal(purchaseRequest.TotalNetPrice) {
		return purchase, nil, ExtendHttpErrorWithDetails(InvalidRequest, "Total net price does not match")
	}

	if !calculatedTotalGrossPrice.Equal(purchaseRequest.TotalGrossPrice) {
		return purchase, nil, ExtendHttpErrorWithDetails(InvalidRequest, "Total gross price does not match")
	}

	purchase.TotalNetPrice = calculatedTotalNetPrice
	purchase.TotalGrossPrice = calculatedTotalGrossPrice
	purchase.CreatedByID = &executingUserObj.ID

	return purchase, updatedListEntries, nil
}

func (handler *Handler) validateAndProcessListEntries(cartItem PurchaseCartRequest, productID uint) ([]models.Guest, error) {
	var updatedListEntries []models.Guest

	for _, listItem := range cartItem.ListItems {
		listEntry, err := handler.repo.GetFullGuestByID(listItem.ID)
		if err != nil || listEntry == nil {
			return nil, ExtendHttpErrorWithDetails(InvalidRequest, "List item not found")
		}

		if listEntry.AttendedGuests != 0 {
			return nil, ExtendHttpErrorWithDetails(InvalidRequest, "List item has already been attended")
		}

		if listEntry.AdditionalGuests+1 < listItem.AttendedGuests {
			return nil, ExtendHttpErrorWithDetails(InvalidRequest, "Additional guests exceed available guests")
		}

		if listEntry.Guestlist.ProductID != productID {
			return nil, ExtendHttpErrorWithDetails(InvalidRequest, "List item does not belong to product")
		}

		listEntry.AttendedGuests = listItem.AttendedGuests
		listEntry.MarkAsArrived()
		updatedListEntries = append(updatedListEntries, *listEntry)
	}

	return updatedListEntries, nil
}

func (handler *Handler) updateGuestEntries(c *gin.Context, updatedListEntries []models.Guest, purchaseID uint) {
	for i := range updatedListEntries {
		updatedListEntries[i].PurchaseID = &purchaseID

		_, err := handler.repo.UpdateGuestByID(int(updatedListEntries[i].ID), updatedListEntries[i])
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		}
	}
}

func (handler *Handler) sendNotifications(updatedListEntries []models.Guest) {
	for _, entry := range updatedListEntries {
		if entry.NotifyOnArrivalEmail != nil {
			_ = handler.mailer.SendNotificationOnArrival(*entry.NotifyOnArrivalEmail, entry.Name)
		}
	}
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

	purchasesResponse := response.ToPurchasesResponse(purchases, handler.decimalPlaces)

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, purchasesResponse)
}

func (handler *Handler) GetPurchaseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	purchase, err := handler.repo.GetPurchaseByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	purchaseResponse := response.ToPurchaseResponse(*purchase, handler.decimalPlaces)

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
