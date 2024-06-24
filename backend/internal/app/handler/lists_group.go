package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

func (handler *Handler) GetListGroups(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.ListGroupFilters{}
	filters.IDs = queryArrayInt(c, "id");
	filters.ListID, _ = strconv.Atoi(c.DefaultQuery("list", "0"))
	
	lists, err := handler.repo.GetListsGroups(end-start, start, sort, order, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
		
	total, err := handler.repo.GetTotalListGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetListGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListGroupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

type ListGroupRequest struct {
	Name      string  `form:"name"  json:"name" binding:"required"`
	ListID	uint    `form:"listId" json:"listId" binding:"required"`
}

func (handler *Handler) UpdateListGroupByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listGroup, err := handler.repo.GetListGroupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var listGroupRequest ListGroupRequest
	if err := c.ShouldBind(&listGroupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	listGroup.Name = listGroupRequest.Name
	listGroup.ListID = listGroupRequest.ListID
	listGroup.UpdatedByID = &executingUserObj.ID

	listGroup, err = handler.repo.UpdateListGroupByID(id, *listGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, listGroup)
}

func (handler *Handler) CreateListGroup(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var listGroup models.ListGroup
	var listGroupRequest ListGroupRequest
	if c.ShouldBind(&listGroupRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	listGroup.Name = listGroupRequest.Name
	listGroup.ListID = listGroupRequest.ListID
	listGroup.CreatedByID = &executingUserObj.ID

	listGroup, err = handler.repo.CreateListGroup(listGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, listGroup)
}

func (handler *Handler) DeleteListGroupByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	listGroup, err := handler.repo.GetListGroupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	if !executingUserObj.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	handler.repo.DeleteListGroup(*listGroup, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
