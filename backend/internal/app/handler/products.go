package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func (handler *Handler) GetProducts(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")

	products, err := handler.repo.GetProducts(end-start, start, sort, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := handler.repo.GetTotalProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, products)
}

func (handler *Handler) GetProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

type ProductRequest struct {
	Name      string  `form:"name" binding:"required"`
	Price     float64 `form:"price" binding:"numeric,required"`
	WrapAfter bool    `form:"wrapAfter"`
	Pos       int     `form:"pos" binding:"numeric,required"`
	ApiExport bool    `form:"apiExport"`
}

func (handler *Handler) UpdateProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var productRequest ProductRequest
	if c.ShouldBind(&productRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product.Name = productRequest.Name
	product.Price = productRequest.Price
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.ApiExport = productRequest.ApiExport

	product, err = handler.repo.UpdateProductByID(id, *product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) CreateProduct(c *gin.Context) {
	var product models.Product
	var productRequest ProductRequest
	if c.ShouldBind(&productRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product.Name = productRequest.Name
	product.Price = productRequest.Price
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.ApiExport = productRequest.ApiExport

	product, err := handler.repo.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	handler.repo.DeleteProduct(*product)

	c.JSON(http.StatusOK, gin.H{})
}
