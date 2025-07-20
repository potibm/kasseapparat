package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	response "github.com/potibm/kasseapparat/internal/app/response"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

const invalidPurchaseIDMsg = "Invalid purchase ID"

func (handler *Handler) DeletePurchase(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, invalidPurchaseIDMsg))
		return
	}

	handler.repo.DeletePurchaseByID(id, *executingUserObj)

	_ = handler.repo.RollbackVisitedGuestsByPurchaseID(id)

	c.Status(http.StatusNoContent)
}

func (handler *Handler) RefundPurchase(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, invalidPurchaseIDMsg))
		return
	}

	purchase, err := handler.repo.GetPurchaseByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, "Purchase not found"))
		return
	}

	isCreator := purchase.CreatedByID != nil && *purchase.CreatedByID == uint(executingUserObj.ID)
	if !executingUserObj.Admin && !isCreator {
		_ = c.Error(ExtendHttpErrorWithDetails(Forbidden, "You are not allowed to refund this purchase"))
		return
	}

	if !executingUserObj.Admin && time.Since(purchase.CreatedAt) > 15*time.Minute {
		_ = c.Error(ExtendHttpErrorWithDetails(Forbidden, "You can only refund purchases within 15 minutes of creation"))
		return
	}

	purchase, err = handler.purchaseService.RefundPurchase(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	purchaseResponse := response.ToPurchaseResponse(*purchase, handler.decimalPlaces)

	c.JSON(http.StatusOK, purchaseResponse)
}

func (handler *Handler) PostPurchases(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	var req PurchaseRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	err = handler.ValidatePaymentMethodPayload(req.PaymentMethod, req.SumupReaderID)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	input := req.ToInput()

	// create a variable purchase that is a nil pointer to models.Purchase
	var purchase *models.Purchase

	if req.PaymentMethod == models.PaymentMethodSumUp {
		purchase, err = handler.purchaseService.CreatePendingPurchase(c.Request.Context(), input, int(executingUserObj.ID))
	} else {
		purchase, err = handler.purchaseService.CreateConfirmedPurchase(c.Request.Context(), input, int(executingUserObj.ID))
	}

	if err != nil {
		switch err {
		case purchaseService.ErrInvalidProductPrice,
			purchaseService.ErrInvalidTotalGrossPrice,
			purchaseService.ErrInvalidTotalNetPrice,
			purchaseService.ErrProductNotFound,
			purchaseService.ErrGuestNotFound,
			purchaseService.ErrGuestAlreadyAttended,
			purchaseService.ErrTooManyAdditionalGuests,
			purchaseService.ErrListItemWrongProduct:
			_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, utils.CapitalizeFirstRune(err.Error())))
			return
		default:
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
			return
		}
	}

	reloadedPurchase, err := handler.repo.GetPurchaseByID(purchase.ID)
	if err != nil {
		reloadedPurchase = purchase // Fallback to the created purchase if reloading fails
	}

	if req.PaymentMethod == models.PaymentMethodSumUp {
		clientTransactionId, err := handler.sumupRepository.CreateReaderCheckout(
			req.SumupReaderID,
			purchase.TotalGrossPrice,
			"Purchase from Kasseapparat",
			purchase.ID.String(),
			handler.sumupRepository.GetWebhookUrl(),
		)
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to create SumUp reader checkout: "+err.Error()))

			log.Printf("Error creating SumUp reader checkout: %v", err)

			return
		}

		log.Printf("Created SumUp reader checkout: %s", *clientTransactionId)

		_, err = handler.repo.UpdatePurchaseSumupClientTransactionIDByID(reloadedPurchase.ID, *clientTransactionId)
		if err != nil {
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, "Failed to update purchase with SumUp transaction ID"))
			return
		}

		log.Printf("Updated purchase %s with SumUp client transaction ID %s", reloadedPurchase.ID, *clientTransactionId)
		log.Printf("Monitor: %+v", handler.monitor)

		handler.monitor.Start(reloadedPurchase.ID)
		log.Printf("Started monitoring for purchase %s", reloadedPurchase.ID)
	}

	purchaseResponse := response.ToPurchaseResponse(*reloadedPurchase, handler.decimalPlaces)

	c.JSON(http.StatusCreated, purchaseResponse)
}

func (handler *Handler) GetPurchases(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "createdAt")
	order := c.DefaultQuery("_order", "DESC")

	filters := sqliteRepo.PurchaseFilters{}
	filters.PaymentMethods = queryPaymentMethods(c, "paymentMethod", handler.config.PaymentMethods)
	filters.CreatedByID, _ = strconv.Atoi(c.DefaultQuery("createdById", "0"))
	filters.TotalGrossPriceGte = queryDecimal(c, "totalGrossPrice_gte")
	filters.TotalGrossPriceLte = queryDecimal(c, "totalGrossPrice_lte")
	filters.IDs = queryArrayInt(c, "id")
	filters.Status = queryPurchaseStatus(c, "status")

	purchases, err := handler.repo.GetPurchases(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	total, err := handler.repo.GetTotalPurchases(filters)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	purchasesResponse := response.ToPurchasesResponse(purchases, handler.decimalPlaces)

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, purchasesResponse)
}

func (handler *Handler) GetPurchaseByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, invalidPurchaseIDMsg))
		return
	}

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
