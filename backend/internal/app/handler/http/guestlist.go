package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

type GuestlistCreateRequest struct {
	Name      string `json:"name"      form:"name"      binding:"required"`
	TypeCode  bool   `json:"typeCode"  form:"typeCode"  binding:"boolean"`
	ProductID int    `json:"productId" form:"productId" binding:"required"`
}

type GuestlistUpdateRequest struct {
	Name      string `json:"name"      form:"name"      binding:"required"`
	TypeCode  bool   `json:"typeCode"  form:"typeCode"  binding:"boolean"`
	ProductID int    `json:"productId" form:"productId" binding:"required"`
}

func (handler *Handler) GetGuestlists(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := sqliteRepo.GuestlistFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.IDs = queryArrayInt(c, "id")

	lists, err := handler.repo.GetGuestlists(end-start, start, sort, order, filters)
	if err != nil {
		listQueryError(c, err)

		return
	}

	total, err := handler.repo.GetTotalGuestlists()
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	ids := make([]int, 0, len(lists))

	for _, list := range lists {
		ids = append(ids, list.ID)
	}

	stats, err := handler.repo.GetGuestlistStats(ids)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	response := make([]GuestlistWithStats, 0, len(lists))

	for _, list := range lists {
		response = append(response, GuestlistWithStats{Guestlist: list, GuestlistStats: stats[list.ID]})
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, response)
}

// GuestlistWithStats is a guest list in the list endpoint, with how many guests are
// expected and how many have arrived.
type GuestlistWithStats struct {
	models.Guestlist
	sqliteRepo.GuestlistStats
}

func (handler *Handler) GetGuestlistByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	list, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) UpdateGuestlistByID(c *gin.Context) {
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	guestlist, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	var guestlistRequest GuestlistUpdateRequest
	if err := c.ShouldBind(&guestlistRequest); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	guestlist.Name = guestlistRequest.Name
	guestlist.TypeCode = guestlistRequest.TypeCode

	if guestlistRequest.ProductID > 0 {
		guestlist.ProductID = guestlistRequest.ProductID
	}

	guestlist, err = handler.repo.UpdateGuestlistByID(id, *guestlist)
	if err != nil {
		_ = c.Error(InternalServerError.WithCause(err))

		return
	}

	c.JSON(http.StatusOK, guestlist)
}

func (handler *Handler) CreateGuestlist(c *gin.Context) {
	c = handler.contextWithUser(c)

	var guestlist models.Guestlist

	var guestlistRequest GuestlistCreateRequest
	if err := c.ShouldBind(&guestlistRequest); err != nil {
		_ = c.Error(InvalidRequest.WithCauseMsg(err))

		return
	}

	guestlist.Name = guestlistRequest.Name
	guestlist.TypeCode = guestlistRequest.TypeCode
	guestlist.ProductID = guestlistRequest.ProductID

	newGuestlist, err := handler.repo.CreateGuestlist(guestlist)
	if err != nil {
		_ = c.Error(InternalServerError.WithCauseMsg(err))

		return
	}

	c.JSON(http.StatusCreated, newGuestlist)
}

func (handler *Handler) DeleteGuestlistByID(c *gin.Context) {
	c = handler.contextWithUser(c)

	id, _ := strconv.Atoi(c.Param("id"))

	guestlist, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(NotFound.WithCause(err))

		return
	}

	handler.repo.DeleteGuestlist(*guestlist)

	c.Status(http.StatusNoContent)
}
