package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func (handler *Handler) GetLists(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")

	lists, err := handler.repo.GetLists(end-start, start, sort, order)
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

type ListRequest struct {
	Name      string  `form:"name"  json:"name" binding:"required"`
	TypeCode bool    `form:"typeCode" json:"typeCode" binding:"boolean"`
}

func (handler *Handler) UpdateListByID(c *gin.Context) {
	user, _ := c.Get(middleware.IdentityKey)
	userObj, _ := user.(*models.User)

	id, _ := strconv.Atoi(c.Param("id"))
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var listRequest ListRequest
	if err := c.ShouldBind(&listRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	list.UpdatedByID = &userObj.ID

	list, err = handler.repo.UpdateListByID(id, *list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (handler *Handler) CreateList(c *gin.Context) {
	userObj := handler.GetUserFromContext(c)

	var list models.List
	var listRequest ListRequest
	if c.ShouldBind(&listRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	list.Name = listRequest.Name
	list.TypeCode = listRequest.TypeCode
	list.CreatedByID = &userObj.ID

	product, err := handler.repo.CreateList(list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteListByID(c *gin.Context) {
	log.Println("lets go")
	userObj := handler.GetUserFromContext(c)
	log.Println(userObj)

	id, _ := strconv.Atoi(c.Param("id"))
	log.Println(id)
	list, err := handler.repo.GetListByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	log.Println(list)
	log.Println(userObj)
	if userObj.Admin {
		log.Println("user is admin")
	}
	if !userObj.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	handler.repo.DeleteList(*list, userObj)

	c.JSON(http.StatusOK, gin.H{})
}
