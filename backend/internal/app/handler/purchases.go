package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potibm/kasseapparat/internal/app/repository"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/potibm/kasseapparat/internal/app/service"
	"github.com/potibm/kasseapparat/internal/app/utils"
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
	PaymentMethod   string                `binding:"required"      form:"paymentMethod"`
}

func (req PurchaseRequest) ToInput() service.PurchaseInput {
	input := service.PurchaseInput{
		PaymentMethod:   req.PaymentMethod,
		TotalNetPrice:   req.TotalNetPrice,
		TotalGrossPrice: req.TotalGrossPrice,
	}

	for _, cart := range req.Cart {
		item := service.PurchaseCartItem{
			ID:       cart.ID,
			Quantity: cart.Quantity,
		}
		for _, li := range cart.ListItems {
			item.ListItems = append(item.ListItems, service.ListItemInput{
				ID:             li.ID,
				AttendedGuests: int(li.AttendedGuests),
			})
		}

		input.Cart = append(input.Cart, item)
	}

	return input
}

func (handler *Handler) DeletePurchase(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id := c.Param("id")
	if uuid.Validate(id) != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Invalid purchase ID"))
		return
	}

	handler.repo.DeletePurchaseByID(id, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted"})
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

	if !handler.IsValidPaymentMethod(req.PaymentMethod) {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Invalid payment method"))
		return
	}

	input := req.ToInput()
	purchaseService := service.NewPurchaseService(handler.repo, &handler.mailer, handler.decimalPlaces)

	purchase, err := purchaseService.CreatePurchase(c.Request.Context(), input, int(executingUserObj.ID))
	if err != nil {
		switch err {
		case service.ErrInvalidProductPrice,
			service.ErrInvalidTotalGrossPrice,
			service.ErrInvalidTotalNetPrice,
			service.ErrProductNotFound,
			service.ErrGuestNotFound,
			service.ErrGuestAlreadyAttended,
			service.ErrTooManyAdditionalGuests,
			service.ErrListItemWrongProduct:
			_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, utils.CapitalizeFirstRune(err.Error())))
			return
		default:
			_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
			return
		}
	}

	purchaseResponse := response.ToPurchaseResponse(*purchase, handler.decimalPlaces)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Purchase successful",
		"purchase": purchaseResponse,
	})
}

func (handler *Handler) GetPurchases(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "DESC")
	filters := repository.PurchaseFilters{}

	filters.PaymentMethods = queryPaymentMethods(c, "paymentMethod", handler.paymentMethods)
	filters.CreatedByID, _ = strconv.Atoi(c.DefaultQuery("createdById", "0"))
	filters.TotalGrossPriceGte = queryDecimal(c, "totalGrossPrice_gte")
	filters.TotalGrossPriceLte = queryDecimal(c, "totalGrossPrice_lte")
	filters.IDs = queryArrayInt(c, "id")

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
	id := c.Param("id")

	// validate that id is a uuid
	if uuid.Validate(id) != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, "Invalid purchase ID"))
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
