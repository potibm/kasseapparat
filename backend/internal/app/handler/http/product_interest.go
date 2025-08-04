package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductInterestCreateRequest struct {
	ProductID uint `binding:"required" form:"productId" json:"productId"`
}

func (handler *Handler) GetProductInterests(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	ids := queryArrayInt(c, "id")

	lists, err := handler.repo.GetProductInterests(end-start, start, ids)
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	total, err := handler.repo.GetTotalProductInterests()
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) DeleteProductInterestByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser.WithCause(err))

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	productInterest, err := handler.repo.GetProductInterestByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	handler.repo.DeleteProductInterest(*productInterest, *executingUserObj)

	c.Status(http.StatusNoContent)
}

func (handler *Handler) CreateProductInterest(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser.WithCause(err))

		return
	}

	var productInterest models.ProductInterest

	var productInterestRequest ProductInterestCreateRequest
	if err := c.ShouldBind(&productInterestRequest); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	productInterest.ProductID = productInterestRequest.ProductID

	product, err := handler.repo.GetProductByID(int(productInterest.ProductID)) // check if product exists
	if product == nil || err != nil {
		_ = c.Error(BadRequest.WithMsg("Product not found").WithCause(err))

		return
	}

	productInterest, err = handler.repo.CreateProductInterest(productInterest, *executingUserObj)
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	c.JSON(http.StatusCreated, productInterest)
}
