package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUpdatePasswordRequest struct {
	UserId   int    `binding:"required"        form:"userId"   json:"userId"`
	Token    string `binding:"required,len=32" form:"token"    json:"token"`
	Password string `binding:"required,min=8"  form:"password" json:"password"`
}

func (handler *Handler) UpdateUserPassword(c *gin.Context) {
	var userPasswordChangeRequest UserUpdatePasswordRequest
	if err := c.ShouldBind(&userPasswordChangeRequest); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	user, err := handler.repo.GetUserByID(userPasswordChangeRequest.UserId)
	if err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(BadRequest, "User not found"))

		return
	}

	if !user.ChangePasswordTokenIsValid(userPasswordChangeRequest.Token) {
		_ = c.Error(ExtendHttpErrorWithDetails(BadRequest, "Token is invalid or has expired"))

		return
	}

	user.Password = userPasswordChangeRequest.Password
	user.ChangePasswordToken = nil
	user.ChangePasswordTokenExpiry = nil

	user, err = handler.repo.UpdateUserByID(int(user.ID), *user)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	c.JSON(http.StatusOK, user)
}

type RequestChangePasswordTokenRequest struct {
	Login string `binding:"required" form:"login" json:"login"`
}

func (handler *Handler) RequestChangePasswordToken(c *gin.Context) {
	var request RequestChangePasswordTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))

		return
	}

	user, err := handler.repo.GetUserByUsernameOrEmail(request.Login)
	if err != nil {
		c.JSON(http.StatusOK, "OK")

		return
	}

	if user.ChangePasswordTokenExpiry != nil && !user.ChangePasswordTokenIsExpired() {
		c.JSON(http.StatusOK, "OK")

		return
	}

	user.GenerateChangePasswordToken(nil)

	user, err = handler.repo.UpdateUserByID(int(user.ID), *user)
	if err != nil {
		_ = c.Error(InternalServerError)

		return
	}

	err = handler.mailer.SendChangePasswordTokenMail(
		user.Email, user.ID, user.Username, *user.ChangePasswordToken)
	if err != nil {
		log.Println("Error sending email", err)
	}

	c.JSON(http.StatusOK, "OK")
}
