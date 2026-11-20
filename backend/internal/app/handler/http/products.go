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

// ProductRequestCreate is the body for a new product. SectionID is optional: without
// it the product goes into the default section (ADR 006).
type ProductRequestCreate struct {
	Name      string          `json:"name"      form:"name"      binding:"required"`
	NetPrice  decimal.Decimal `json:"netPrice"  form:"netPrice"  binding:"required"`
	VATRate   decimal.Decimal `json:"vatRate"   form:"vatRate"   binding:"required"`
	WrapAfter bool            `json:"wrapAfter" form:"wrapAfter"`
	Pos       int             `json:"pos"       form:"pos"       binding:"numeric,required"`
	Hidden    bool            `json:"hidden"    form:"hidden"    binding:"boolean"`
	SectionID int             `json:"sectionId" form:"sectionId"`
}

// ProductRequestUpdate is the body for editing a product. SectionID is optional: when
// it is omitted or 0 the product keeps its current section. 0 does not mean "the
// default section" – to move a product there, send that section's id.
type ProductRequestUpdate struct {
	Name       string          `json:"name"       form:"name"       binding:"required"`
	NetPrice   decimal.Decimal `json:"netPrice"   form:"netPrice"   binding:"required"`
	VATRate    decimal.Decimal `json:"vatRate"    form:"vatRate"    binding:"required"`
	WrapAfter  bool            `json:"wrapAfter"  form:"wrapAfter"`
	Pos        int             `json:"pos"        form:"pos"        binding:"numeric,required"`
	APIExport  bool            `json:"apiExport"  form:"apiExport"  binding:"boolean"`
	Hidden     bool            `json:"hidden"     form:"hidden"     binding:"boolean"`
	SoldOut    bool            `json:"soldOut"    form:"soldOut"    binding:"boolean"`
	TotalStock int             `json:"totalStock" form:"totalStock" binding:"numeric"`
	SectionID  int             `json:"sectionId"  form:"sectionId"`
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
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.APIExport = productRequest.APIExport
	product.Hidden = productRequest.Hidden
	product.SoldOut = productRequest.SoldOut
	product.TotalStock = productRequest.TotalStock

	if !handler.assignSectionOnUpdate(c, product, productRequest.SectionID) {
		return
	}

	// A mode-owned flag is read-only for products with a stock limit: the form
	// still sends its current value, so ignore it rather than rejecting.
	if product.TotalStock > 0 {
		behavior := config.OutOfStockIgnore
		if ps, ok := handler.purchaseService.(*purchaseService.PurchaseService); ok {
			behavior = ps.OutOfStockBehavior
		}

		switch behavior {
		case config.OutOfStockAutoSoldOut:
			stored, _ := handler.repo.GetProductByID(id)
			if stored != nil {
				product.SoldOut = stored.SoldOut
			}
		case config.OutOfStockAutoHide:
			stored, _ := handler.repo.GetProductByID(id)
			if stored != nil {
				product.Hidden = stored.Hidden
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
	product.WrapAfter = productRequest.WrapAfter
	product.Pos = productRequest.Pos
	product.Hidden = productRequest.Hidden

	sectionID, err := handler.resolveSectionID(productRequest.SectionID)
	if err != nil {
		sectionResolveError(c, productRequest.SectionID, err)

		return
	}

	product.SectionID = sectionID

	product, err = handler.repo.CreateProduct(product)
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
