package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUpdatePasswordRequest struct {
	UserId   int    `form:"userId" json:"userId" binding:"required"`
	Token    string `form:"token" json:"token" binding:"required,len=32"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
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
	Login string `form:"login" json:"login" binding:"required"`
}

func (handler *Handler) RequestChangePasswordToken(c *gin.Context) {
	var request RequestChangePasswordTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(ExtendHttpErrorWithDetails(InvalidRequest, err.Error()))
		return
	}

	user, err := handler.repo.GetUserByUserameOrEmail(request.Login)
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
