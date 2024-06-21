package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func (handler *Handler) GetUsers(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "pos")
	order := c.DefaultQuery("_order", "ASC")

	products, err := handler.repo.GetUsers(end-start, start, sort, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := handler.repo.GetTotalUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, products)
}

func (handler *Handler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

type UserCreateRequest struct {
	Usermame string `form:"username"  json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Admin	bool   `form:"admin" json:"admin" binding:""`
}

type UserUpdateRequest struct {
	Username string `form:"username"  json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:""`
	Admin	bool   `form:"admin" json:"admin" binding:""`
}

func (handler *Handler) UpdateUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := handler.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var userRequest UserUpdateRequest
	if err := c.ShouldBind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	user.Username = userRequest.Username
	user.Password = ""

	// an admin may change the password of another user
	// a user may change his own password
	if user.Admin || int(user.ID) == id {
		log.Println("password " + userRequest.Password)
		if userRequest.Password != "" {
			log.Println("password " + userRequest.Password)
			user.Password = userRequest.Password
		}
	}

	// only an admin may change the role of a user
	if user.Admin {
		user.Admin = userRequest.Admin
	}

	user, err = handler.repo.UpdateUserByID(id, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler *Handler) CreateUser(c *gin.Context) {
	var user models.User

	var userRequest UserCreateRequest
	if c.ShouldBind(&userRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user.Username = userRequest.Usermame
	user.Password = userRequest.Password
	
	// only an admin may change the role of a user
	if user.Admin {
		user.Admin = userRequest.Admin
	} else {
		user.Admin = false
	}

	product, err := handler.repo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (handler *Handler) DeleteUserByID(c *gin.Context) { 
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve the executing user"})
		return
	}
	
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := handler.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// only admins are allowed to delete users, and not themselves
	if !executingUserObj.Admin || executingUserObj.ID == user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	handler.repo.DeleteUser(*user)

	c.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) getUserFromContext(c *gin.Context) (*models.User, error) {
	user, _ := c.Get(middleware.IdentityKey)
	sparseUserObjFromJwt, _ := user.(*models.User)
	
	userObj, err := handler.repo.GetUserByID(int(sparseUserObjFromJwt.ID))
	if err != nil {
		return nil, errors.New("User not found")
	}

	return userObj, nil
}