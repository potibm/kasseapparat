package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductInterestCreateRequest struct {
	ProductID uint `form:"productId" json:"productId" binding:"required"`
}

func (handler *Handler) GetProductInterests(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	ids := queryArrayInt(c, "id")

	lists, err := handler.repo.GetProductInterests(end-start, start, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := handler.repo.GetTotalLists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) DeleteProductInterestByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	productInterest, err := handler.repo.GetProductInterestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	handler.repo.DeleteProductInterest(*productInterest, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) CreateProductInterest(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var productInterest models.ProductInterest
	var productInterestRequest ProductInterestCreateRequest
	if c.ShouldBind(&productInterestRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	productInterest.ProductID = productInterestRequest.ProductID
	product, err := handler.repo.GetProductByID(int(productInterest.ProductID)) // check if product exists
	if product == nil || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	productInterest, err = handler.repo.CreateProductInterest(productInterest, *executingUserObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, productInterest)
}
