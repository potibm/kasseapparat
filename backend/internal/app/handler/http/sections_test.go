package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	apperrors "github.com/potibm/kasseapparat/internal/app/errors"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withUser(c *gin.Context, role string) {
	c.Set(middleware.AuthUserKey, models.AuthUser{Username: "tester", Role: role})
}

func lastErrorStatus(t *testing.T, c *gin.Context) int {
	t.Helper()

	require.NotEmpty(t, c.Errors, "expected the handler to report an error")

	var basic *apperrors.BasicError
	require.ErrorAs(t, c.Errors.Last().Err, &basic)

	return basic.StatusCode()
}

func TestCreateSection_AsAdmin(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, w := ginContextWithParam(t, http.MethodPost, "/")
	withUser(c, adminRole)
	jsonBody(t, c, SectionRequest{Name: "Tickets", Pos: 1})

	handler.CreateSection(c)

	require.Empty(t, c.Errors)
	require.Equal(t, http.StatusCreated, w.Code)

	var section models.Section
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &section))
	assert.NotZero(t, section.ID)
	assert.Equal(t, "Tickets", section.Name)
	assert.Equal(t, 1, section.Pos)
}

func TestCreateSection_AllowsAnEmptyName(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, w := ginContextWithParam(t, http.MethodPost, "/")
	withUser(c, adminRole)
	jsonBody(t, c, SectionRequest{Name: "", Pos: 2})

	handler.CreateSection(c)

	require.Empty(t, c.Errors)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSectionWrites_OpenToEveryUser(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, w := ginContextWithParam(t, http.MethodPost, "/")
	withUser(c, "user")
	jsonBody(t, c, SectionRequest{Name: "Tickets", Pos: 1})

	handler.CreateSection(c)

	require.Empty(t, c.Errors)
	require.Equal(t, http.StatusCreated, w.Code)

	var created models.Section
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))

	c, w = ginContextWithParam(t, http.MethodPut, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(created.ID)}}
	withUser(c, "user")
	jsonBody(t, c, SectionRequest{Name: "Entry", Pos: 2})

	handler.UpdateSectionByID(c)

	require.Empty(t, c.Errors)
	assert.Equal(t, http.StatusOK, w.Code)

	reloaded, err := handler.repo.GetSectionByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Entry", reloaded.Name)
}

func TestDeleteSectionByID_RequiresAdmin(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	id := strconv.Itoa(section.ID)

	t.Run("as user", func(t *testing.T) {
		c, _ := ginContextWithParam(t, http.MethodDelete, "/")
		c.Params = gin.Params{{Key: "id", Value: id}}
		withUser(c, "user")

		handler.DeleteSectionByID(c)

		assert.Equal(t, http.StatusForbidden, lastErrorStatus(t, c))
	})

	t.Run("without user", func(t *testing.T) {
		c, _ := ginContextWithParam(t, http.MethodDelete, "/")
		c.Params = gin.Params{{Key: "id", Value: id}}

		handler.DeleteSectionByID(c)

		assert.Equal(t, http.StatusUnauthorized, lastErrorStatus(t, c))
	})

	_, err = handler.repo.GetSectionByID(section.ID)
	assert.NoError(t, err, "the section must still exist")
}

func TestGetSections_ListsByPosWithTotalCount(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	_, err := handler.repo.CreateSection(models.Section{Name: "Merch", Pos: 2})
	require.NoError(t, err)

	_, err = handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodGet, "/?_start=0&_end=10")

	handler.GetSections(c)

	require.Empty(t, c.Errors)
	assert.Equal(t, "2", w.Header().Get("X-Total-Count"))

	var sections []models.Section
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &sections))
	require.Len(t, sections, 2)
	assert.Equal(t, "Tickets", sections[0].Name)
	assert.Equal(t, "Merch", sections[1].Name)
}

func TestGetSectionByID(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodGet, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(section.ID)}}

	handler.GetSectionByID(c)

	require.Empty(t, c.Errors)
	assert.Contains(t, w.Body.String(), `"name":"Tickets"`)

	c, _ = ginContextWithParam(t, http.MethodGet, "/")
	c.Params = gin.Params{{Key: "id", Value: "4711"}}

	handler.GetSectionByID(c)

	assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
}

func TestUpdateSectionByID_AsAdmin(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodPut, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(section.ID)}}
	withUser(c, adminRole)
	jsonBody(t, c, SectionRequest{Name: "Entry", Pos: 4})

	handler.UpdateSectionByID(c)

	require.Empty(t, c.Errors)
	assert.Equal(t, http.StatusOK, w.Code)

	reloaded, err := handler.repo.GetSectionByID(section.ID)
	require.NoError(t, err)
	assert.Equal(t, "Entry", reloaded.Name)
	assert.Equal(t, 4, reloaded.Pos)
}

func TestUpdateSectionByID_NotFound(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, _ := ginContextWithParam(t, http.MethodPut, "/")
	c.Params = gin.Params{{Key: "id", Value: "4711"}}
	withUser(c, adminRole)
	jsonBody(t, c, SectionRequest{Name: "Entry", Pos: 4})

	handler.UpdateSectionByID(c)

	assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
}

func TestDeleteSectionByID_Empty(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Empty", Pos: 1})
	require.NoError(t, err)

	c, _ := ginContextWithParam(t, http.MethodDelete, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(section.ID)}}
	withUser(c, adminRole)

	handler.DeleteSectionByID(c)

	require.Empty(t, c.Errors)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
}

func TestDeleteSectionByID_WithProductsIsAConflict(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	_, err = handler.repo.CreateProduct(models.Product{Name: "Regular", SectionID: section.ID})
	require.NoError(t, err)

	c, _ := ginContextWithParam(t, http.MethodDelete, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(section.ID)}}
	withUser(c, adminRole)

	handler.DeleteSectionByID(c)

	assert.Equal(t, http.StatusConflict, lastErrorStatus(t, c))

	_, err = handler.repo.GetSectionByID(section.ID)
	assert.NoError(t, err, "the section must still exist")
}

func TestDeleteSectionByID_NotFound(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, _ := ginContextWithParam(t, http.MethodDelete, "/")
	c.Params = gin.Params{{Key: "id", Value: "4711"}}
	withUser(c, adminRole)

	handler.DeleteSectionByID(c)

	assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
}

func productCreateRequest(sectionID int) ProductRequestCreate {
	return ProductRequestCreate{
		Name:      "Mate",
		NetPrice:  decimal.NewFromInt(2),
		VATRate:   decimal.NewFromInt(19),
		Pos:       1,
		SectionID: sectionID,
	}
}

func TestCreateProduct_WithoutSectionGoesIntoTheDefaultSection(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, w := ginContextWithParam(t, http.MethodPost, "/")
	jsonBody(t, c, productCreateRequest(0))

	handler.CreateProduct(c)

	require.Empty(t, c.Errors)
	require.Equal(t, http.StatusCreated, w.Code)

	var product models.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &product))

	defaultSection, err := handler.repo.GetOrCreateDefaultSection()
	require.NoError(t, err)
	assert.Equal(t, defaultSection.ID, product.SectionID)
}

func TestCreateProduct_InTheRequestedSection(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	_, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	merch, err := handler.repo.CreateSection(models.Section{Name: "Merch", Pos: 2})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodPost, "/")
	jsonBody(t, c, productCreateRequest(merch.ID))

	handler.CreateProduct(c)

	require.Empty(t, c.Errors)

	var product models.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &product))
	assert.Equal(t, merch.ID, product.SectionID)
}

func TestCreateProduct_WithUnknownSectionIsRejected(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	c, _ := ginContextWithParam(t, http.MethodPost, "/")
	jsonBody(t, c, productCreateRequest(4711))

	handler.CreateProduct(c)

	assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
}

func productUpdateRequest(sectionID int) ProductRequestUpdate {
	return ProductRequestUpdate{
		Name:      "Mate",
		NetPrice:  decimal.NewFromInt(2),
		VATRate:   decimal.NewFromInt(19),
		Pos:       1,
		SectionID: sectionID,
	}
}

func TestUpdateProduct_SectionHandling(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	tickets, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	merch, err := handler.repo.CreateSection(models.Section{Name: "Merch", Pos: 2})
	require.NoError(t, err)

	product, err := handler.repo.CreateProduct(models.Product{Name: "Mate", SectionID: merch.ID})
	require.NoError(t, err)

	update := func(t *testing.T, sectionID int) *gin.Context {
		t.Helper()

		c, _ := ginContextWithParam(t, http.MethodPut, "/")
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(product.ID)}}
		jsonBody(t, c, productUpdateRequest(sectionID))

		handler.UpdateProductByID(c)

		return c
	}

	t.Run("without a section the product keeps its section", func(t *testing.T) {
		c := update(t, 0)

		require.Empty(t, c.Errors)

		reloaded, err := handler.repo.GetProductByID(product.ID)
		require.NoError(t, err)
		assert.Equal(t, merch.ID, reloaded.SectionID)
	})

	t.Run("with a section the product moves there", func(t *testing.T) {
		c := update(t, tickets.ID)

		require.Empty(t, c.Errors)

		reloaded, err := handler.repo.GetProductByID(product.ID)
		require.NoError(t, err)
		assert.Equal(t, tickets.ID, reloaded.SectionID)
	})

	t.Run("an unknown section is rejected", func(t *testing.T) {
		c := update(t, 4711)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))

		reloaded, err := handler.repo.GetProductByID(product.ID)
		require.NoError(t, err)
		assert.Equal(t, tickets.ID, reloaded.SectionID)
	})
}

func TestUpdateProduct_WithoutSectionGetsTheDefault(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	// a product stored before sections existed
	product, err := handler.repo.CreateProduct(models.Product{Name: "Mate"})
	require.NoError(t, err)
	require.Zero(t, product.SectionID)

	c, _ := ginContextWithParam(t, http.MethodPut, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(product.ID)}}
	jsonBody(t, c, productUpdateRequest(0))

	handler.UpdateProductByID(c)

	require.Empty(t, c.Errors)

	reloaded, err := handler.repo.GetProductByID(product.ID)
	require.NoError(t, err)
	assert.NotZero(t, reloaded.SectionID)
}

func TestGetProductByID_ExposesTheSection(t *testing.T) {
	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	section, err := handler.repo.CreateSection(models.Section{Name: "Tickets", Pos: 1})
	require.NoError(t, err)

	product, err := handler.repo.CreateProduct(models.Product{Name: "Mate", SectionID: section.ID})
	require.NoError(t, err)

	c, w := ginContextWithParam(t, http.MethodGet, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(product.ID)}}

	handler.GetProductByID(c)

	require.Empty(t, c.Errors)
	assert.Contains(t, w.Body.String(), `"sectionId":`+strconv.Itoa(section.ID))
}
