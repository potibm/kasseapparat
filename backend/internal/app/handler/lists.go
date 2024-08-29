package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

type ListCreateRequest struct {
	Name      string `form:"name"  json:"name" binding:"required"`
	TypeCode  bool   `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint   `form:"productId" json:"productId" binding:"required"`
}

type ListUpdateRequest struct {
	Name      string `form:"name"  json:"name" binding:"required"`
	TypeCode  bool   `form:"typeCode" json:"typeCode" binding:"boolean"`
	ProductID uint   `form:"productId" json:"productId" binding:"required"`
}

func (handler *Handler) GetLists(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.ListFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.IDs = queryArrayInt(c, "id")

	lists, err := handler.repo.GetLists(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))
		return
	}

	total, err := handler.repo.GetTotalLists()
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, lists)
}

func (handler *Handler) GetListByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) UpdateListByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	var listRequest ListUpdateRequest
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

	list, err = handler.repo.UpdateListByID(id, *list)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) CreateList(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	var list models.List
	var listRequest ListCreateRequest
	if err := c.ShouldBind(&listRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	list.ProductID = listRequest.ProductID
	list.CreatedByID = &executingUserObj.ID

	product, err := handler.repo.CreateList(list)
	if err != nil {
		_ = c.Error(InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (handler *Handler) DeleteListByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))
		return
	}

	if !executingUserObj.Admin && *list.CreatedByID != executingUserObj.ID {
		_ = c.Error(Forbidden)
		return
	}

	handler.repo.DeleteList(*list, *executingUserObj)

	c.JSON(http.StatusOK, gin.H{})
}
