package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
)

func (handler *Handler) GetUsers(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	sort := c.DefaultQuery("_sort", "id")
	order := c.DefaultQuery("_order", "ASC")
	filters := repository.UserFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.IsAdmin = c.DefaultQuery("isAdmin", "false") == "true"

	products, err := handler.repo.GetUsers(end-start, start, sort, order, filters)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InternalServerError, err.Error()))

		return
	}

	total, err := handler.repo.GetTotalUsers(&filters)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, products)
}

func (handler *Handler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := handler.repo.GetUserByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	c.JSON(http.StatusOK, product)
}

type UserCreateRequest struct {
	Username string `binding:"required" form:"username" json:"username"`
	Email    string `binding:"required" form:"email"    json:"email"`
	Admin    bool   `binding:""         form:"admin"    json:"admin"`
}

type UserUpdateRequest struct {
	Username string `binding:"required" form:"username" json:"username"`
	Password string `binding:""         form:"password" json:"password"`
	Email    string `binding:"required" form:"email"    json:"email"`
	Admin    bool   `binding:""         form:"admin"    json:"admin"`
}

func (handler *Handler) UpdateUserByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	user, err := handler.repo.GetUserByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}

	var userRequest UserUpdateRequest
	if err := c.ShouldBind(&userRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	user.Username = userRequest.Username
	user.Password = ""
	user.Email = userRequest.Email

	// an admin may change the password of another user
	// a user may change his own password
	if executingUserObj.Admin || int(executingUserObj.ID) == id {
		if userRequest.Password != "" {
			user.Password = userRequest.Password
		}
	}

	// only an admin may change the role of a user
	if executingUserObj.Admin {
		user.Admin = userRequest.Admin
	}

	user, err = handler.repo.UpdateUserByID(id, *user)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler *Handler) CreateUser(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	var user models.User

	var userRequest UserCreateRequest
	if err := c.ShouldBind(&userRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	user.Username = userRequest.Username
	user.Email = userRequest.Email
	user.GenerateRandomPassword()

	validity := 3 * time.Hour
	user.GenerateChangePasswordToken(&validity)

	// only an admin may change the role of a user
	if executingUserObj.Admin {
		user.Admin = userRequest.Admin
	} else {
		user.Admin = false
	}

	user, err = handler.repo.CreateUser(user)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	err = handler.mailer.SendNewUserTokenMail(user.Email, user.ID, user.Username, *user.ChangePasswordToken)
	if err != nil {
		log.Println("Error sending email", err)
	}

	c.JSON(http.StatusCreated, user)
}

func (handler *Handler) DeleteUserByID(c *gin.Context) {
	executingUserObj, err := handler.getUserFromContext(c)
	if err != nil {
		_ = c.Error(UnableToRetrieveExecutingUser)

		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	user, err := handler.repo.GetUserByID(id)

	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(NotFound, err.Error()))

		return
	}
	// only admins are allowed to delete users, and not themselves
	if !executingUserObj.Admin || executingUserObj.ID == user.ID {
		_ = c.Error(Forbidden)

		return
	}

	handler.repo.DeleteUser(*user)

	c.Status(http.StatusNoContent)
}

func (handler *Handler) getUserFromContext(c *gin.Context) (*models.User, error) {
	user, _ := c.Get(middleware.IdentityKey)
	sparseUserObjFromJwt, _ := user.(*models.User)

	userObj, err := handler.repo.GetUserByID(int(sparseUserObjFromJwt.ID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	return userObj, nil
}
