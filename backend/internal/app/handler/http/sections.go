package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

// SectionRequest is the body for creating and updating a section. The name may be
// empty: an unnamed section is shown as a plain row break (ADR 006).
type SectionRequest struct {
	Name string `json:"name" form:"name"`
	Pos  int    `json:"pos"  form:"pos"  binding:"numeric"`
}

const adminRole = "admin"

// requireAdmin stops the request unless the executing user is an admin. Like products,
// sections can be created and changed by every user, but only admins may delete them
// (ADR 006).
func requireAdmin(c *gin.Context) bool {
	authUser, exists := middleware.GetAuthUser(c)
	if !exists {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return false
	}

	if authUser.Role != adminRole {
		_ = c.Error(Forbidden.WithMsg("Only admins can delete sections"))

		return false
	}

	return true
}

func (handler *Handler) GetSections(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")
	ids := queryArrayInt(c, "id")

	sections, err := handler.repo.GetSections(end-start, start, sort, order, ids)
	if err != nil {
		listQueryError(c, err)

		return
	}

	total, err := handler.repo.GetTotalSections()
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, sections)
}

func (handler *Handler) GetSectionByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	section, err := handler.repo.GetSectionByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, section)
}

func (handler *Handler) CreateSection(c *gin.Context) {
	c = handler.contextWithUser(c)

	var request SectionRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	section, err := handler.repo.CreateSection(models.Section{
		Name: request.Name,
		Pos:  request.Pos,
	})
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	c.JSON(http.StatusCreated, section)
}

func (handler *Handler) UpdateSectionByID(c *gin.Context) {
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	if _, err := handler.repo.GetSectionByID(id); err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	var request SectionRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	section, err := handler.repo.UpdateSectionByID(id, models.Section{
		Name: request.Name,
		Pos:  request.Pos,
	})
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, section)
}

func (handler *Handler) DeleteSectionByID(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	section, err := handler.repo.GetSectionByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	if err := handler.repo.DeleteSection(*section); err != nil {
		if errors.Is(err, sqliteRepo.ErrSectionNotEmpty) {
			_ = c.Error(Conflict.WithMsg("The section still contains products").WithCause(err))

			return
		}

		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.Status(http.StatusNoContent)
}

// resolveSectionID returns the section a product is stored in: the requested one if
// it exists, or the default section when none was requested.
func (handler *Handler) resolveSectionID(requested int) (int, error) {
	if requested > 0 {
		if _, err := handler.repo.GetSectionByID(requested); err != nil {
			return 0, err
		}

		return requested, nil
	}

	section, err := handler.repo.GetOrCreateDefaultSection()
	if err != nil {
		return 0, err
	}

	return section.ID, nil
}

// sectionResolveError reports a failed resolveSectionID: an unknown requested section
// is the client's mistake, a failure to find or create the default one is ours.
func sectionResolveError(c *gin.Context, requested int, err error) {
	if requested > 0 {
		_ = c.Error(InvalidRequest.WithMsg("Unknown section").WithCause(err))

		return
	}

	_ = c.Error(InternalServerError.WithCause(err))
}

// assignSectionOnUpdate sets the section of an edited product: the requested one, or
// the default section for a product that has none yet. A requested 0 means
// "unchanged", not "default section". It reports the error itself and returns false
// on failure.
func (handler *Handler) assignSectionOnUpdate(c *gin.Context, product *models.Product, requested int) bool {
	if requested == 0 && product.SectionID != 0 {
		return true
	}

	sectionID, err := handler.resolveSectionID(requested)
	if err != nil {
		sectionResolveError(c, requested, err)

		return false
	}

	product.SectionID = sectionID

	return true
}
