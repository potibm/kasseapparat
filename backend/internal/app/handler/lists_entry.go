package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

type ListEntryCreateRequest struct {
	ListID           uint   `form:"listId"  json:"listId" binding:"required"`
	Name             string `form:"name"  json:"name" binding:"required"`
	Code             string `form:"code"  json:"code"`
	AdditionalGuests uint   `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests   uint   `form:"attendedGuests"  json:"attendedGuests"`
}

type ListEntryUpdateRequest struct {
	ListID           uint   `form:"listId"  json:"listId"`
	Name             string `form:"name"  json:"name" binding:"required"`
	Code             string `form:"code"  json:"code"`
	AdditionalGuests uint   `form:"additionalGuests"  json:"additionalGuests"`
	AttendedGuests   uint   `form:"attendedGuests"  json:"attendedGuests"`
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := handler.repo.GetTotalListEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetListEntryByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) UpdateListEntryByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listEntry, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var listEntryRequest ListEntryUpdateRequest
	if err := c.ShouldBind(&listEntryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	listEntry.Name = listEntryRequest.Name
	if listEntryRequest.Code != "" {
		listEntry.Code = &listEntryRequest.Code
	} else {
		listEntry.Code = nil
	}
	if listEntryRequest.ListID > 0 {
		listEntry.ListID = listEntryRequest.ListID
	}
	listEntry.AdditionalGuests = listEntryRequest.AdditionalGuests
	listEntry.AttendedGuests = listEntryRequest.AttendedGuests
	listEntry.UpdatedByID = &executingUserObj.ID

	listEntry, err = handler.repo.UpdateListEntryByID(id, *listEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, listEntry)
}

func (handler *Handler) CreateListEntry(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var listEntry models.ListEntry
	var listEntryRequest ListEntryCreateRequest
	if c.ShouldBind(&listEntryRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	listEntry.Name = listEntryRequest.Name
	listEntry.ListID = listEntryRequest.ListID
	if listEntryRequest.Code != "" {
		listEntry.Code = &listEntryRequest.Code
	} else {
		listEntry.Code = nil
	}
	listEntry.AdditionalGuests = listEntryRequest.AdditionalGuests
	listEntry.AttendedGuests = listEntryRequest.AttendedGuests
	listEntry.CreatedByID = &executingUserObj.ID

	product, err := handler.repo.CreateListEntry(listEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteListEntryByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listEntry, err := handler.repo.GetListEntryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !executingUserObj.Admin && *listEntry.CreatedByID != executingUserObj.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if query != "" {
		listEntries.SortByQuery(query)
	}

	c.JSON(http.StatusOK, listEntries)
}
