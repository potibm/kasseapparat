package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

type GuestCreateRequest struct {
	GuestlistID          uint    `form:"guestlistId"  json:"guestlistId" binding:"required"`
	Name                 string  `form:"name"  json:"name" binding:"required"`
	Code                 string  `form:"code"  json:"code"`
	AdditionalGuests     uint    `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests       uint    `form:"attendedGuests"  json:"attendedGuests"`
	ArrivalNote          *string `form:"arrivalNote" json:"arrivalNote"`
	NotifyOnArrivalEmail *string `form:"notifyOnArrivalEmail" json:"notifyOnArrivalEmail"`
}

type GuestUpdateRequest struct {
	GuestlistID          uint       `form:"guestlistId"  json:"guestlistId"`
	Name                 string     `form:"name"  json:"name" binding:"required"`
	Code                 string     `form:"code"  json:"code"`
	AdditionalGuests     uint       `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests       uint       `form:"attendedGuests"  json:"attendedGuests"`
	ArrivedAt            *time.Time `form:"arrivedAt" json:"arrivedAt"`
	ArrivalNote          *string    `form:"arrivalNote" json:"arrivalNote"`
	NotifyOnArrivalEmail *string    `form:"notifyOnArrivalEmail" json:"notifyOnArrivalEmail"`
}

func (handler *Handler) GetGuests(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.GuestFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.GuestlistID, _ = strconv.Atoi(c.DefaultQuery("guestlist_id", "0"))
	filters.Present = c.DefaultQuery("isPresent", "false") == "true"
	filters.NotPresent = c.DefaultQuery("isNotPresent", "false") == "true"
	filters.IDs = queryArrayInt(c, "id")

	guests, err := handler.repo.GetGuests(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	total, err := handler.repo.GetTotalGuests(&filters)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, guests)
}

func (handler *Handler) GetGuestsByProductID(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	query := strings.TrimSpace(c.DefaultQuery("q", ""))

	guests, err := handler.repo.GetUnattendedGuestsByProductID(productID, query)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	if query != "" {
		guests.SortByQuery(query)
	}

	c.JSON(http.StatusOK, guests)
}

func (handler *Handler) GetGuestByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	guest, err := handler.repo.GetGuestByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, guest)
}

func (handler *Handler) UpdateGuestByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	guest, err := handler.repo.GetGuestByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	var guestRequest GuestUpdateRequest
	if err := c.ShouldBind(&guestRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	guest.Name = guestRequest.Name
	if guestRequest.Code != "" {
		guest.Code = &guestRequest.Code
	} else {
		guest.Code = nil
	}
	if guestRequest.GuestlistID > 0 {
		guest.GuestlistID = guestRequest.GuestlistID
	}
	guest.AdditionalGuests = guestRequest.AdditionalGuests
	guest.AttendedGuests = guestRequest.AttendedGuests
	guest.UpdatedByID = &executingUserObj.ID
	guest.ArrivedAt = guestRequest.ArrivedAt
	guest.ArrivalNote = guestRequest.ArrivalNote
	guest.NotifyOnArrivalEmail = guestRequest.NotifyOnArrivalEmail

	guest, err = handler.repo.UpdateGuestByID(id, *guest)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, guest)
}

func (handler *Handler) CreateGuest(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	var guest models.Guest
	var guestRequest GuestCreateRequest
	if err := c.ShouldBind(&guestRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	guest.Name = guestRequest.Name
	guest.GuestlistID = guestRequest.GuestlistID
	if guestRequest.Code != "" {
		guest.Code = &guestRequest.Code
	} else {
		guest.Code = nil
	}
	guest.AdditionalGuests = guestRequest.AdditionalGuests
	guest.AttendedGuests = guestRequest.AttendedGuests
	guest.CreatedByID = &executingUserObj.ID
	guest.ArrivalNote = guestRequest.ArrivalNote
	guest.NotifyOnArrivalEmail = guestRequest.NotifyOnArrivalEmail

	product, err := handler.repo.CreateGuest(guest)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteGuestByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	guest, err := handler.repo.GetGuestByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	if !executingUserObj.Admin && *guest.CreatedByID != executingUserObj.ID {
		_ = c.Error(Forbidden)
		return
	}

	handler.repo.DeleteGuest(*guest, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
