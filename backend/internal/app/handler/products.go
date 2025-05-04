package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
	response "github.com/potibm/kasseapparat/internal/app/response"
	"github.com/shopspring/decimal"
)

type ProductRequestCreate struct {
	Name      string          `binding:"required"         form:"name"      json:"name"`
	NetPrice  decimal.Decimal `binding:"required"         form:"netPrice"  json:"netPrice"`
	VATRate   decimal.Decimal `binding:"required"         form:"vatRate"   json:"vatRate"`
	WrapAfter bool            `form:"wrapAfter"           json:"wrapAfter"`
	Pos       int             `binding:"numeric,required" form:"pos"       json:"pos"`
	Hidden    bool            `binding:"boolean"          form:"hidden"    json:"hidden"`
}

type ProductRequestUpdate struct {
	Name       string          `binding:"required"         form:"name"       json:"name"`
	NetPrice   decimal.Decimal `binding:"required"         form:"netPrice"   json:"netPrice"`
	VATRate    decimal.Decimal `binding:"required"         form:"vatRate"    json:"vatRate"`
	WrapAfter  bool            `form:"wrapAfter"           json:"wrapAfter"`
	Pos        int             `binding:"numeric,required" form:"pos"        json:"pos"`
	ApiExport  bool            `binding:"boolean"          form:"apiExport"  json:"apiExport"`
	Hidden     bool            `binding:"boolean"          form:"hidden"     json:"hidden"`
	SoldOut    bool            `binding:"boolean"          form:"soldOut"    json:"soldOut"`
	TotalStock int             `binding:"numeric"          form:"totalStock" json:"totalStock"`
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
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	if filterHidden == "true" {
		products = filterHiddenProducts(products)
	}

	total, err := handler.repo.GetTotalProducts()
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	productsResponse := createExtendedProductResponse(handler.repo, products, handler.decimalPlaces)

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, productsResponse)
}

// filterHiddenProducts removes hidden products from the list while preserving
// wrap-after formatting by transferring the wrap-after property to the previous
// visible product when a hidden product with wrap-after is encountered.
func filterHiddenProducts(products []models.Product) []models.Product {
	var filteredProducts []models.Product

	for _, product := range products {
		if product.Hidden && product.WrapAfter {
			if len(filteredProducts) > 0 {
				filteredProducts[len(filteredProducts)-1].WrapAfter = true
			}
		}

		if !product.Hidden {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts
}

func createExtendedProductResponse(repo *repository.Repository, products []models.Product, decimalPlaces int32) []response.ExtendedProductResponse {
	var productsResponse = []response.ExtendedProductResponse{}

	for _, product := range products {
		unitsSold, _ := repo.GetPurchasedQuantitiesByProductID(product.ID)
		soldOutRequestCount, _ := repo.GetProductInterestCountByProductID(product.ID)

		productResponse := response.ToExtendedProductResponse(
			product,
			unitsSold,
			soldOutRequestCount,
			decimalPlaces,
		)

		productsResponse = append(productsResponse, productResponse)
	}

	return productsResponse
}

func (handler *Handler) GetProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	unitsSold, _ := handler.repo.GetPurchasedQuantitiesByProductID(product.ID)
	soldOutRequestCount, _ := handler.repo.GetProductInterestCountByProductID(product.ID)

	productResponse := response.ToExtendedProductResponse(
		*product,
		unitsSold,
		soldOutRequestCount,
		handler.decimalPlaces,
	)

	c.JSON(http.StatusOK, productResponse)
}

func (handler *Handler) UpdateProductByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	var productRequest ProductRequestUpdate
	if err := c.ShouldBind(&productRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	product.Name = productRequest.Name
	product.NetPrice = productRequest.NetPrice
	product.VATRate = productRequest.VATRate
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.ApiExport = productRequest.ApiExport
	product.Hidden = productRequest.Hidden
	product.UpdatedByID = &executingUserObj.ID
	product.SoldOut = productRequest.SoldOut
	product.TotalStock = productRequest.TotalStock

	product, err = handler.repo.UpdateProductByID(id, *product)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) CreateProduct(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	var product models.Product

	var productRequest ProductRequestCreate
	if err := c.ShouldBind(&productRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	product.Name = productRequest.Name
	product.NetPrice = productRequest.NetPrice
	product.VATRate = productRequest.VATRate
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.Hidden = productRequest.Hidden
	product.CreatedByID = &executingUserObj.ID

	product, err = handler.repo.CreateProduct(product)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteProductByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetProductByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	if !executingUserObj.Admin {
		_ = c.Error(Forbidden)

		return
	}

	handler.repo.DeleteProduct(*product, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
