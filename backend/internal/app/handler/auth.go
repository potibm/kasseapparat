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
	err := c.ShouldBind(&userPasswordChangeRequest)
	if err != nil {
		log.Println("Invalid request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := handler.repo.GetUserByID(userPasswordChangeRequest.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	if !user.ChangePasswordTokenIsValid(userPasswordChangeRequest.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or has expired."})
		return
	}

	user.Password = userPasswordChangeRequest.Password
	user.ChangePasswordToken = nil
	user.ChangePasswordTokenExpiry = nil

	user, err = handler.repo.UpdateUserByID(int(user.ID), *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type RequestChangePasswordTokenRequest struct {
	Login string `form:"login" json:"login" binding:"required"`
}

func (handler *Handler) RequestChangePasswordToken(c *gin.Context) {
	var request RequestChangePasswordTokenRequest
	if c.ShouldBind(&request) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = handler.mailer.SendChangePasswordTokenMail(
		user.Email, user.ID, user.Username, *user.ChangePasswordToken)

	if err != nil {
		log.Println("Error sending email", err)
	}

	c.JSON(http.StatusOK, "OK")
}
