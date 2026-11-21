package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// groupFixture is a shirts section with a "Male" group and a merch section.
type groupFixture struct {
	handler *Handler
	shirts  models.Section
	merch   models.Section
	male    models.ProductGroup
}

func setupGroupFixture(t *testing.T) groupFixture {
	t.Helper()

	handler := setupStockFlagHandler(t, config.OutOfStockIgnore)

	shirts, err := handler.repo.CreateSection(models.Section{Name: "T-Shirts", Pos: 1})
	require.NoError(t, err)

	merch, err := handler.repo.CreateSection(models.Section{Name: "Merch", Pos: 2})
	require.NoError(t, err)

	male, err := handler.repo.CreateProductGroup(models.ProductGroup{SectionID: shirts.ID, Name: "Male", Pos: 10})
	require.NoError(t, err)

	return groupFixture{handler: handler, shirts: shirts, merch: merch, male: male}
}

func groupContext(t *testing.T, method string, id int, payload any) *gin.Context {
	t.Helper()

	c, _ := ginContextWithParam(t, method, "/")
	if id > 0 {
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(id)}}
	}

	if payload != nil {
		jsonBody(t, c, payload)
	}

	return c
}

func TestCreateProductGroup(t *testing.T) {
	f := setupGroupFixture(t)

	t.Run("in the requested section", func(t *testing.T) {
		c, w := ginContextWithParam(t, http.MethodPost, "/")
		withUser(c, "user")
		jsonBody(t, c, ProductGroupRequest{Name: "Female", SectionID: f.shirts.ID, Pos: 20})

		f.handler.CreateProductGroup(c)

		require.Empty(t, c.Errors)
		require.Equal(t, http.StatusCreated, w.Code)

		var created models.ProductGroup
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
		assert.Equal(t, "Female", created.Name)
		assert.Equal(t, f.shirts.ID, created.SectionID)
	})

	t.Run("without a section it goes into the default section", func(t *testing.T) {
		c, w := ginContextWithParam(t, http.MethodPost, "/")
		jsonBody(t, c, ProductGroupRequest{Name: "Hoodie"})

		f.handler.CreateProductGroup(c)

		require.Empty(t, c.Errors)

		var created models.ProductGroup
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
		assert.Equal(t, f.shirts.ID, created.SectionID)
	})

	t.Run("an unknown section is rejected", func(t *testing.T) {
		c := groupContext(t, http.MethodPost, 0, ProductGroupRequest{Name: "Ghost", SectionID: 4711})

		f.handler.CreateProductGroup(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})

	t.Run("a name is required", func(t *testing.T) {
		c := groupContext(t, http.MethodPost, 0, ProductGroupRequest{SectionID: f.shirts.ID})

		f.handler.CreateProductGroup(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})
}

func TestGetProductGroups(t *testing.T) {
	f := setupGroupFixture(t)

	_, err := f.handler.repo.CreateProductGroup(models.ProductGroup{SectionID: f.merch.ID, Name: "Hoodie", Pos: 1})
	require.NoError(t, err)

	t.Run("filtered by section", func(t *testing.T) {
		c, w := ginContextWithParam(t, http.MethodGet, "/?_start=0&_end=10&sectionId="+strconv.Itoa(f.shirts.ID))

		f.handler.GetProductGroups(c)

		require.Empty(t, c.Errors)
		assert.Equal(t, "1", w.Header().Get("X-Total-Count"))

		var groups []models.ProductGroup
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &groups))
		require.Len(t, groups, 1)
		assert.Equal(t, "Male", groups[0].Name)
	})

	t.Run("an unknown sort field is a bad request", func(t *testing.T) {
		c, _ := ginContextWithParam(t, http.MethodGet, "/?_sort=nope")

		f.handler.GetProductGroups(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})
}

func TestGetProductGroupByID(t *testing.T) {
	f := setupGroupFixture(t)

	c, w := ginContextWithParam(t, http.MethodGet, "/")
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(f.male.ID)}}

	f.handler.GetProductGroupByID(c)

	require.Empty(t, c.Errors)
	assert.Contains(t, w.Body.String(), `"name":"Male"`)

	c = groupContext(t, http.MethodGet, 4711, nil)

	f.handler.GetProductGroupByID(c)

	assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
}

func TestUpdateProductGroupByID(t *testing.T) {
	f := setupGroupFixture(t)

	product, err := f.handler.repo.CreateProduct(models.Product{
		Name: "Male S", SectionID: f.shirts.ID, ProductGroupID: &f.male.ID, Variant: "S",
	})
	require.NoError(t, err)

	reloadGroup := func(t *testing.T) *models.ProductGroup {
		t.Helper()

		group, err := f.handler.repo.GetProductGroupByID(f.male.ID)
		require.NoError(t, err)

		return group
	}

	t.Run("without a section the group keeps its section", func(t *testing.T) {
		c := groupContext(t, http.MethodPut, f.male.ID, ProductGroupRequest{Name: "Men", Pos: 5})

		f.handler.UpdateProductGroupByID(c)

		require.Empty(t, c.Errors)

		group := reloadGroup(t)
		assert.Equal(t, "Men", group.Name)
		assert.Equal(t, 5, group.Pos)
		assert.Equal(t, f.shirts.ID, group.SectionID)
	})

	t.Run("moving the group moves its products", func(t *testing.T) {
		c := groupContext(t, http.MethodPut, f.male.ID, ProductGroupRequest{Name: "Men", SectionID: f.merch.ID})

		f.handler.UpdateProductGroupByID(c)

		require.Empty(t, c.Errors)
		assert.Equal(t, f.merch.ID, reloadGroup(t).SectionID)

		reloaded, err := f.handler.repo.GetProductByID(product.ID)
		require.NoError(t, err)
		assert.Equal(t, f.merch.ID, reloaded.SectionID)
	})

	t.Run("an unknown section is rejected", func(t *testing.T) {
		c := groupContext(t, http.MethodPut, f.male.ID, ProductGroupRequest{Name: "Men", SectionID: 4711})

		f.handler.UpdateProductGroupByID(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})

	t.Run("a name is required", func(t *testing.T) {
		c := groupContext(t, http.MethodPut, f.male.ID, ProductGroupRequest{})

		f.handler.UpdateProductGroupByID(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})

	t.Run("an unknown group is not found", func(t *testing.T) {
		c := groupContext(t, http.MethodPut, 4711, ProductGroupRequest{Name: "Ghost"})

		f.handler.UpdateProductGroupByID(c)

		assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
	})
}

func TestDeleteProductGroupByID(t *testing.T) {
	f := setupGroupFixture(t)

	empty, err := f.handler.repo.CreateProductGroup(models.ProductGroup{SectionID: f.shirts.ID, Name: "Empty"})
	require.NoError(t, err)

	_, err = f.handler.repo.CreateProduct(models.Product{
		Name: "Male S", SectionID: f.shirts.ID, ProductGroupID: &f.male.ID,
	})
	require.NoError(t, err)

	t.Run("as user", func(t *testing.T) {
		c := groupContext(t, http.MethodDelete, empty.ID, nil)
		withUser(c, "user")

		f.handler.DeleteProductGroupByID(c)

		assert.Equal(t, http.StatusForbidden, lastErrorStatus(t, c))
	})

	t.Run("without user", func(t *testing.T) {
		c := groupContext(t, http.MethodDelete, empty.ID, nil)

		f.handler.DeleteProductGroupByID(c)

		assert.Equal(t, http.StatusUnauthorized, lastErrorStatus(t, c))
	})

	t.Run("a group with products is a conflict", func(t *testing.T) {
		c := groupContext(t, http.MethodDelete, f.male.ID, nil)
		withUser(c, adminRole)

		f.handler.DeleteProductGroupByID(c)

		assert.Equal(t, http.StatusConflict, lastErrorStatus(t, c))
	})

	t.Run("an unknown group is not found", func(t *testing.T) {
		c := groupContext(t, http.MethodDelete, 4711, nil)
		withUser(c, adminRole)

		f.handler.DeleteProductGroupByID(c)

		assert.Equal(t, http.StatusNotFound, lastErrorStatus(t, c))
	})

	t.Run("an empty group is deleted", func(t *testing.T) {
		c := groupContext(t, http.MethodDelete, empty.ID, nil)
		withUser(c, adminRole)

		f.handler.DeleteProductGroupByID(c)

		require.Empty(t, c.Errors)
		assert.Equal(t, http.StatusNoContent, c.Writer.Status())

		_, err := f.handler.repo.GetProductGroupByID(empty.ID)
		assert.Error(t, err)
	})
}

func TestDeleteSectionByID_WithProductGroupsIsAConflict(t *testing.T) {
	f := setupGroupFixture(t)

	c := groupContext(t, http.MethodDelete, f.shirts.ID, nil)
	withUser(c, adminRole)

	f.handler.DeleteSectionByID(c)

	assert.Equal(t, http.StatusConflict, lastErrorStatus(t, c))
}

func intPtr(v int) *int {
	return &v
}

func TestCreateProduct_InAGroup(t *testing.T) {
	f := setupGroupFixture(t)

	t.Run("takes the section of the group", func(t *testing.T) {
		request := productCreateRequest(f.merch.ID)
		request.ProductGroupID = intPtr(f.male.ID)
		request.Variant = "M"

		c, w := ginContextWithParam(t, http.MethodPost, "/")
		jsonBody(t, c, request)

		f.handler.CreateProduct(c)

		require.Empty(t, c.Errors)
		require.Equal(t, http.StatusCreated, w.Code)

		var created models.Product
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
		require.NotNil(t, created.ProductGroupID)
		assert.Equal(t, f.male.ID, *created.ProductGroupID)
		assert.Equal(t, f.shirts.ID, created.SectionID, "the group's section wins")
		assert.Equal(t, "M", created.Variant)
	})

	t.Run("an unknown group is rejected", func(t *testing.T) {
		request := productCreateRequest(0)
		request.ProductGroupID = intPtr(4711)

		c := groupContext(t, http.MethodPost, 0, request)

		f.handler.CreateProduct(c)

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
	})
}

func TestUpdateProduct_GroupHandling(t *testing.T) {
	f := setupGroupFixture(t)

	product, err := f.handler.repo.CreateProduct(models.Product{Name: "Mate", SectionID: f.merch.ID})
	require.NoError(t, err)

	update := func(t *testing.T, sectionID int, groupID *int) *gin.Context {
		t.Helper()

		request := productUpdateRequest(sectionID)
		request.ProductGroupID = groupID
		request.Variant = "L"

		c := groupContext(t, http.MethodPut, product.ID, request)

		f.handler.UpdateProductByID(c)

		return c
	}

	reload := func(t *testing.T) *models.Product {
		t.Helper()

		reloaded, err := f.handler.repo.GetProductByID(product.ID)
		require.NoError(t, err)

		return reloaded
	}

	t.Run("a group moves the product into the group's section", func(t *testing.T) {
		c := update(t, f.merch.ID, intPtr(f.male.ID))

		require.Empty(t, c.Errors)

		reloaded := reload(t)
		require.NotNil(t, reloaded.ProductGroupID)
		assert.Equal(t, f.male.ID, *reloaded.ProductGroupID)
		assert.Equal(t, f.shirts.ID, reloaded.SectionID)
		assert.Equal(t, "L", reloaded.Variant)
	})

	t.Run("without a group the product stays in its group", func(t *testing.T) {
		c := update(t, f.merch.ID, nil)

		require.Empty(t, c.Errors)

		reloaded := reload(t)
		require.NotNil(t, reloaded.ProductGroupID)
		assert.Equal(t, f.shirts.ID, reloaded.SectionID)
	})

	t.Run("an unknown group is rejected", func(t *testing.T) {
		c := update(t, 0, intPtr(4711))

		assert.Equal(t, http.StatusBadRequest, lastErrorStatus(t, c))
		assert.NotNil(t, reload(t).ProductGroupID)
	})

	t.Run("group 0 takes the product out and keeps its section", func(t *testing.T) {
		c := update(t, 0, intPtr(0))

		require.Empty(t, c.Errors)

		reloaded := reload(t)
		assert.Nil(t, reloaded.ProductGroupID)
		assert.Equal(t, f.shirts.ID, reloaded.SectionID)
	})

	t.Run("a standalone product moves to the requested section", func(t *testing.T) {
		c := update(t, f.merch.ID, nil)

		require.Empty(t, c.Errors)

		reloaded := reload(t)
		assert.Nil(t, reloaded.ProductGroupID)
		assert.Equal(t, f.merch.ID, reloaded.SectionID)
	})
}
