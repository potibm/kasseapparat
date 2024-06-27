package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
)


type ListCreateRequest struct {
	Name      string  `form:"name"  json:"name" binding:"required"`
	TypeCode bool    `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint  `form:"productId" json:"productId" binding:"required"`
}

type ListUpdateRequest struct {
	Name      string  `form:"name"  json:"name" binding:"required"`
	TypeCode bool    `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint  `form:"productId" json:"productId" binding:"required"`
}

func (handler *Handler) GetLists(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	ids := queryArrayInt(c, "id");
	
	lists, err := handler.repo.GetLists(end-start, start, sort, order, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
		
	total, err := handler.repo.GetTotalLists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetListByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}


func (handler *Handler) UpdateListByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var listRequest ListUpdateRequest
	if err := c.ShouldBind(&listRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	if listRequest.ProductID > 0 {
		list.ProductID = listRequest.ProductID
	}
	list.UpdatedByID = &executingUserObj.ID

	list, err = handler.repo.UpdateListByID(id, *list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) CreateList(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	var list models.List
	var listRequest ListCreateRequest
	if c.ShouldBind(&listRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	list.ProductID = listRequest.ProductID
	list.CreatedByID = &executingUserObj.ID

	product, err := handler.repo.CreateList(list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteListByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	if !executingUserObj.Admin && *list.CreatedByID != executingUserObj.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	handler.repo.DeleteList(*list, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
