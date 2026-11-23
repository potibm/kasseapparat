package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	response "github.com/potibm/kasseapparat/internal/app/response"
	purchaseService "github.com/potibm/kasseapparat/internal/app/service/purchase"
	"github.com/shopspring/decimal"
)

// ProductRequestCreate is the body for a new product. SectionID and ProductGroupID are
// optional: without them the product goes into the default section. A product in a
// group always takes the group's section (ADR 006).
type ProductRequestCreate struct {
	Name           string          `json:"name"           form:"name"           binding:"required"`
	NetPrice       decimal.Decimal `json:"netPrice"       form:"netPrice"       binding:"required"`
	VATRate        decimal.Decimal `json:"vatRate"        form:"vatRate"        binding:"required"`
	Pos            int             `json:"pos"            form:"pos"            binding:"numeric,required"`
	Hidden         bool            `json:"hidden"         form:"hidden"         binding:"boolean"`
	SectionID      int             `json:"sectionId"      form:"sectionId"`
	ProductGroupID *int            `json:"productGroupId" form:"productGroupId"`
	Variant        string          `json:"variant"        form:"variant"`
}

// ProductRequestUpdate is the body for editing a product.
//
// SectionID: omitted or 0 keeps the current section. 0 does not mean "the default
// section" – to move a product there, send that section's id.
//
// ProductGroupID: omitted or null keeps the current group, 0 takes the product out of
// its group, any other id moves it into that group. A product in a group always takes
// the group's section; a SectionID sent along is ignored then (ADR 006).
type ProductRequestUpdate struct {
	Name           string          `json:"name"           form:"name"           binding:"required"`
	NetPrice       decimal.Decimal `json:"netPrice"       form:"netPrice"       binding:"required"`
	VATRate        decimal.Decimal `json:"vatRate"        form:"vatRate"        binding:"required"`
	Pos            int             `json:"pos"            form:"pos"            binding:"numeric,required"`
	APIExport      bool            `json:"apiExport"      form:"apiExport"      binding:"boolean"`
	Hidden         bool            `json:"hidden"         form:"hidden"         binding:"boolean"`
	SoldOut        bool            `json:"soldOut"        form:"soldOut"        binding:"boolean"`
	TotalStock     int             `json:"totalStock"     form:"totalStock"     binding:"numeric"`
	SectionID      int             `json:"sectionId"      form:"sectionId"`
	ProductGroupID *int            `json:"productGroupId" form:"productGroupId"`
	Variant        string          `json:"variant"        form:"variant"`
}

func (handler *Handler) GetProducts(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")
	filterHidden := c.DefaultQuery("_filter_hidden", "false")
	filters := sqliteRepo.ProductFilters{
		IDs:        queryArrayInt(c, "id"),
		SectionIDs: queryArrayInt(c, "sectionId"),
	}

	products, err := handler.repo.GetProducts(end-start, start, sort, order, filters)
	if err != nil {
		listQueryError(c, err)

		return
	}

	if filterHidden == "true" {
		products = filterHiddenProducts(products)
	}

	total, err := handler.repo.GetTotalProducts(filters)
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	productsResponse := createExtendedProductResponse(handler.repo, products, handler.decimalPlaces)

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, productsResponse)
}

// filterHiddenProducts removes hidden products from the list.
func filterHiddenProducts(products []models.Product) []models.Product {
	var filteredProducts []models.Product

	for _, product := range products {
		if !product.Hidden {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts
}

func createExtendedProductResponse(
	repo sqliteRepo.RepositoryInterface,
	products []models.Product,
	decimalPlaces int32,
) []response.ExtendedProductResponse {
	productsResponse := make([]response.ExtendedProductResponse, 0, len(products))

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
		_ = c.Error(NotFound.WithCauseMsg(err))

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
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	var productRequest ProductRequestUpdate
	if err := c.ShouldBind(&productRequest); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	product.Name = productRequest.Name
	product.NetPrice = productRequest.NetPrice
	product.VATRate = productRequest.VATRate
	product.Pos = productRequest.Pos
	product.APIExport = productRequest.APIExport
	product.Hidden = productRequest.Hidden
	product.SoldOut = productRequest.SoldOut
	product.TotalStock = productRequest.TotalStock
	product.Variant = productRequest.Variant

	if !handler.assignProductLayout(c, product, productRequest.SectionID, productRequest.ProductGroupID) {
		return
	}

	// A mode-owned flag is read-only for products with a stock limit: the form
	// still sends its current value, so ignore it rather than rejecting.
	if product.TotalStock > 0 {
		behavior := config.OutOfStockIgnore
		if ps, ok := handler.purchaseService.(*purchaseService.PurchaseService); ok {
			behavior = ps.OutOfStockBehavior
		}

		if behavior == config.OutOfStockAutoSoldOut {
			stored, _ := handler.repo.GetProductByID(id)
			if stored != nil {
				product.SoldOut = stored.SoldOut
			}
		}
	}

	product, err = handler.repo.UpdateProductByID(id, *product)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	// The stock limit is part of what the out-of-stock behaviour derives the product
	// flags from, so an edit can make a flag — or its release — due.
	handler.reconcileStockFlags(c)

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) CreateProduct(c *gin.Context) {
	c = handler.contextWithUser(c)

	var product models.Product

	var productRequest ProductRequestCreate
	if err := c.ShouldBind(&productRequest); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	product.Name = productRequest.Name
	product.NetPrice = productRequest.NetPrice
	product.VATRate = productRequest.VATRate
	product.Pos = productRequest.Pos
	product.Hidden = productRequest.Hidden
	product.Variant = productRequest.Variant

	if !handler.assignProductLayout(c, &product, productRequest.SectionID, productRequest.ProductGroupID) {
		return
	}

	product, err := handler.repo.CreateProduct(product)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteProductByID(c *gin.Context) {
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	product, err := handler.repo.GetProductByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	handler.repo.DeleteProduct(*product)

	c.Status(http.StatusNoContent)
}
