package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

// ProductGroupRequest is the body for creating and updating a product group.
// SectionID: on create, 0 puts the group into the default section; on update, 0 keeps
// the current section. Moving a group moves its products with it (ADR 006).
type ProductGroupRequest struct {
	Name      string `json:"name"      form:"name"      binding:"required"`
	SectionID int    `json:"sectionId" form:"sectionId"`
	Pos       int    `json:"pos"       form:"pos"       binding:"numeric"`
}

func (handler *Handler) GetProductGroups(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")

	filters := sqliteRepo.ProductFilters{
		IDs:        queryArrayInt(c, "id"),
		SectionIDs: queryArrayInt(c, "sectionId"),
	}

	groups, err := handler.repo.GetProductGroups(end-start, start, sort, order, filters)
	if err != nil {
		listQueryError(c, err)

		return
	}

	total, err := handler.repo.GetTotalProductGroups(filters)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, groups)
}

func (handler *Handler) GetProductGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	group, err := handler.repo.GetProductGroupByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, group)
}

func (handler *Handler) CreateProductGroup(c *gin.Context) {
	c = handler.contextWithUser(c)

	var request ProductGroupRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	sectionID, err := handler.resolveSectionID(request.SectionID)
	if err != nil {
		sectionResolveError(c, request.SectionID, err)

		return
	}

	group, err := handler.repo.CreateProductGroup(models.ProductGroup{
		SectionID: sectionID,
		Name:      request.Name,
		Pos:       request.Pos,
	})
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	c.JSON(http.StatusCreated, group)
}

func (handler *Handler) UpdateProductGroupByID(c *gin.Context) {
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	group, err := handler.repo.GetProductGroupByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	var request ProductGroupRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	if request.SectionID > 0 {
		if _, err := handler.repo.GetSectionByID(request.SectionID); err != nil {
			_ = c.Error(InvalidRequest.WithMsg("Unknown section").WithCause(err))

			return
		}

		group.SectionID = request.SectionID
	}

	group.Name = request.Name
	group.Pos = request.Pos

	updated, err := handler.repo.UpdateProductGroupByID(id, *group)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, updated)
}

func (handler *Handler) DeleteProductGroupByID(c *gin.Context) {
	if !requireAdmin(c, "Only admins can delete product groups") {
		return
	}

	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	group, err := handler.repo.GetProductGroupByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	if err := handler.repo.DeleteProductGroup(*group); err != nil {
		if errors.Is(err, sqliteRepo.ErrProductGroupNotEmpty) {
			_ = c.Error(Conflict.WithMsg("The product group still contains products").WithCause(err))

			return
		}

		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.Status(http.StatusNoContent)
}

// assignProductLayout places a created or edited product (ADR 006):
//
//   - groupID > 0: into that group, taking the group's section; sectionID is ignored.
//   - groupID 0: out of any group, into sectionID (0 keeps the current section).
//   - groupID nil: the group stays as it is. A grouped product keeps its section, a
//     standalone one goes into sectionID (0 keeps the current or, for a new product,
//     the default section).
//
// It reports an error itself and returns false on failure.
func (handler *Handler) assignProductLayout(
	c *gin.Context,
	product *models.Product,
	sectionID int,
	groupID *int,
) bool {
	if groupID != nil && *groupID > 0 {
		group, err := handler.repo.GetProductGroupByID(*groupID)
		if err != nil {
			_ = c.Error(InvalidRequest.WithMsg("Unknown product group").WithCause(err))

			return false
		}

		product.ProductGroupID = &group.ID
		product.SectionID = group.SectionID

		return true
	}

	if groupID != nil {
		product.ProductGroupID = nil
	}

	if product.ProductGroupID != nil {
		return true
	}

	return handler.assignSectionOnUpdate(c, product, sectionID)
}
