package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

type GuestlistCreateRequest struct {
	Name      string `form:"name"  json:"name" binding:"required"`
	TypeCode  bool   `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint   `form:"productId" json:"productId" binding:"required"`
}

type GuestlistUpdateRequest struct {
	Name      string `form:"name"  json:"name" binding:"required"`
	TypeCode  bool   `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint   `form:"productId" json:"productId" binding:"required"`
}

func (handler *Handler) GetGuestlists(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.GuestlistFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.IDs = queryArrayInt(c, "id")

	lists, err := handler.repo.GetGuestlists(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	total, err := handler.repo.GetTotalGuestlists()
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetGuestlistByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) UpdateGuestlistByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	var listRequest GuestlistUpdateRequest
	if err := c.ShouldBind(&listRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	if listRequest.ProductID > 0 {
		list.ProductID = listRequest.ProductID
	}
	list.UpdatedByID = &executingUserObj.ID

	list, err = handler.repo.UpdateGuestlistByID(id, *list)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) CreateGuestlist(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	var guestlist models.Guestlist
	var listRequest GuestlistCreateRequest
	if err := c.ShouldBind(&listRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	guestlist.Name = listRequest.Name
	guestlist.TypeCode = listRequest.TypeCode
	guestlist.ProductID = listRequest.ProductID
	guestlist.CreatedByID = &executingUserObj.ID

	product, err := handler.repo.CreateGuestlist(guestlist)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteGuestlistByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetGuestlistByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	if !executingUserObj.Admin && *list.CreatedByID != executingUserObj.ID {
		_ = c.Error(Forbidden)
		return
	}

	handler.repo.DeleteGuestlist(*list, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
