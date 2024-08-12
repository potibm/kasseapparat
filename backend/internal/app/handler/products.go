package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)

type ProductRequestCreate struct {
	Name      string  `form:"name"  json:"name" binding:"required"`
	Price     float64 `form:"price" json:"price" binding:"numeric"`
	WrapAfter bool    `form:"wrapAfter" json:"wrapAfter"`
	Pos       int     `form:"pos" json:"pos" binding:"numeric,required"`
	Hidden    bool    `form:"hidden" json:"hidden" binding:"boolean"`
}

type ProductRequestUpdate struct {
	Name       string  `form:"name"  json:"name" binding:"required"`
	Price      float64 `form:"price" json:"price" binding:"numeric"`
	WrapAfter  bool    `form:"wrapAfter" json:"wrapAfter"`
	Pos        int     `form:"pos" json:"pos" binding:"numeric,required"`
	ApiExport  bool    `form:"apiExport" json:"apiExport" binding:"boolean"`
	Hidden     bool    `form:"hidden" json:"hidden" binding:"boolean"`
	SoldOut    bool    `form:"soldOut" json:"soldOut" binding:"boolean"`
	TotalStock int     `form:"totalStock" json:"totalStock" binding:"numeric"`
}

func (handler *Handler) GetProducts(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")
	filterHidden := c.DefaultQuery("_filter_hidden", "false")
	ids := queryArrayInt(c, "id")

	products, err := handler.repo.GetProducts(end-start, start, sort, order, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if filterHidden == "true" {
		var filteredProducts []models.ProductWithSalesAndInterrest
		for _, product := range products {
			if product.Hidden && product.WrapAfter {
				// find last product in filteredProducts and set wrapAfter to true
				if len(filteredProducts) > 0 {
					filteredProducts[len(filteredProducts)-1].WrapAfter = true
				}
			}
			if !product.Hidden {
				filteredProducts = append(filteredProducts, product)
			}
		}
		products = filteredProducts
	}

	total, err := handler.repo.GetTotalProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// iterate over all products and set the SoldOutRequestCount and UnitsSold
	for i := range products {
		products[i].UnitsSold, _ = handler.repo.GetPurchasedQuantitiesByProductID(products[i].ID)
		products[i].SoldOutRequestCount, _ = handler.repo.GetProductInterestCountByProductID(products[i].ID)
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, products)
}

func (handler *Handler) GetProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByIDWithSalesAndInterrest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	product.UnitsSold, _ = handler.repo.GetPurchasedQuantitiesByProductID(product.ID)
	product.SoldOutRequestCount, _ = handler.repo.GetProductInterestCountByProductID(product.ID)

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) UpdateProductByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var productRequest ProductRequestUpdate
	if err := c.ShouldBind(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	product.Name = productRequest.Name
	product.Price = productRequest.Price
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.ApiExport = productRequest.ApiExport
	product.Hidden = productRequest.Hidden
	product.UpdatedByID = &executingUserObj.ID
	product.SoldOut = productRequest.SoldOut
	product.TotalStock = productRequest.TotalStock

	product, err = handler.repo.UpdateProductByID(id, *product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) CreateProduct(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var product models.Product
	var productRequest ProductRequestCreate
	if c.ShouldBind(&productRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product.Name = productRequest.Name
	product.Price = productRequest.Price
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.Hidden = productRequest.Hidden
	product.CreatedByID = &executingUserObj.ID

	product, err = handler.repo.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteProductByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !executingUserObj.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	handler.repo.DeleteProduct(*product, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
