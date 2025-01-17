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

type ListEntryCreateRequest struct {
	GuestlistID          uint    `form:"guestlistId"  json:"guestlistId" binding:"required"`
	Name                 string  `form:"name"  json:"name" binding:"required"`
	Code                 string  `form:"code"  json:"code"`
	AdditionalGuests     uint    `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests       uint    `form:"attendedGuests"  json:"attendedGuests"`
	ArrivalNote          *string `form:"arrivalNote" json:"arrivalNote"`
	NotifyOnArrivalEmail *string `form:"notifyOnArrivalEmail" json:"notifyOnArrivalEmail"`
}

type ListEntryUpdateRequest struct {
	GuestlistID          uint       `form:"guestlistId"  json:"guestlistId"`
	Name                 string     `form:"name"  json:"name" binding:"required"`
	Code                 string     `form:"code"  json:"code"`
	AdditionalGuests     uint       `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests       uint       `form:"attendedGuests"  json:"attendedGuests"`
	ArrivedAt            *time.Time `form:"arrivedAt" json:"arrivedAt"`
	ArrivalNote          *string    `form:"arrivalNote" json:"arrivalNote"`
	NotifyOnArrivalEmail *string    `form:"notifyOnArrivalEmail" json:"notifyOnArrivalEmail"`
}

func (handler *Handler) GetListEntries(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.ListEntryFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.ListID, _ = strconv.Atoi(c.DefaultQuery("list", "0"))
	filters.Present = c.DefaultQuery("isPresent", "false") == "true"
	filters.NotPresent = c.DefaultQuery("isNotPresent", "false") == "true"
	filters.IDs = queryArrayInt(c, "id")

	lists, err := handler.repo.GetListEntries(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	total, err := handler.repo.GetTotalListEntries(&filters)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetListEntryByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) UpdateListEntryByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listEntry, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	var listEntryRequest ListEntryUpdateRequest
	if err := c.ShouldBind(&listEntryRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	listEntry.Name = listEntryRequest.Name
	if listEntryRequest.Code != "" {
		listEntry.Code = &listEntryRequest.Code
	} else {
		listEntry.Code = nil
	}
	if listEntryRequest.GuestlistID > 0 {
		listEntry.GuestlistID = listEntryRequest.GuestlistID
	}
	listEntry.AdditionalGuests = listEntryRequest.AdditionalGuests
	listEntry.AttendedGuests = listEntryRequest.AttendedGuests
	listEntry.UpdatedByID = &executingUserObj.ID
	listEntry.ArrivedAt = listEntryRequest.ArrivedAt
	listEntry.ArrivalNote = listEntryRequest.ArrivalNote
	listEntry.NotifyOnArrivalEmail = listEntryRequest.NotifyOnArrivalEmail

	listEntry, err = handler.repo.UpdateListEntryByID(id, *listEntry)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, listEntry)
}

func (handler *Handler) CreateListEntry(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	var listEntry models.ListEntry
	var listEntryRequest ListEntryCreateRequest
	if err := c.ShouldBind(&listEntryRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	listEntry.Name = listEntryRequest.Name
	listEntry.GuestlistID = listEntryRequest.GuestlistID
	if listEntryRequest.Code != "" {
		listEntry.Code = &listEntryRequest.Code
	} else {
		listEntry.Code = nil
	}
	listEntry.AdditionalGuests = listEntryRequest.AdditionalGuests
	listEntry.AttendedGuests = listEntryRequest.AttendedGuests
	listEntry.CreatedByID = &executingUserObj.ID
	listEntry.ArrivalNote = listEntryRequest.ArrivalNote
	listEntry.NotifyOnArrivalEmail = listEntryRequest.NotifyOnArrivalEmail

	product, err := handler.repo.CreateListEntry(listEntry)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteListEntryByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listEntry, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	if !executingUserObj.Admin && *listEntry.CreatedByID != executingUserObj.ID {
		_ = c.Error(Forbidden)
		return
	}

	handler.repo.DeleteListEntry(*listEntry, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) GetListEntriesByProductID(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	query := strings.TrimSpace(c.DefaultQuery("q", ""))

	listEntries, err := handler.repo.GetUnattendedListEntriesByProductID(productID, query)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	if query != "" {
		listEntries.SortByQuery(query)
	}

	c.JSON(http.StatusOK, listEntries)
}
