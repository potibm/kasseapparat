package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductInterrestCreateRequest struct {
	ProductID uint `form:"productId" json:"productId" binding:"required"`
}

func (handler *Handler) GetProductInterrests(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	ids := queryArrayInt(c, "id")

	lists, err := handler.repo.GetProductInterrests(end-start, start, ids)
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

func (handler *Handler) DeleteProductInterrestByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	productInterrest, err := handler.repo.GetProductInterrestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	handler.repo.DeleteProductInterrest(*productInterrest, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) CreateProductInterrest(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var productInterrest models.ProductInterrest
	var productInterrestRequest ProductInterrestCreateRequest
	if c.ShouldBind(&productInterrestRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	productInterrest.ProductID = productInterrestRequest.ProductID
	product, err := handler.repo.GetProductByID(int(productInterrest.ProductID)) // check if product exists
	if product == nil || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	productInterrest, err = handler.repo.CreateProductInterrest(productInterrest, *executingUserObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, productInterrest)
}
