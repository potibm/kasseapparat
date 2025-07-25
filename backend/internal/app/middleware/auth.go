package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

var IdentityKey = "ID"

type login struct {
	Login    string `binding:"required" form:"login"    json:"login"`
	Password string `binding:"required" form:"password" json:"password"`
}

type loginResponse struct {
	Code        int     `json:"code"`
	Token       string  `json:"token"`
	Expire      string  `json:"expire"`
	Role        *string `json:"role"`
	Username    *string `json:"username"`
	GravatarUrl *string `json:"gravatarUrl"`
	Id          *uint   `json:"id"`
}

func HandlerMiddleWare(authMiddleware *ginjwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func RegisterRoute(r *gin.Engine, handle *ginjwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	auth := r.Group("/auth", handle.MiddlewareFunc())
	auth.GET("/refresh_token", handle.RefreshHandler)
}

func InitParams(repo *sqliteRepo.Repository, realm string, secret string, timeout int) *ginjwt.GinJWTMiddleware {
	if secret == "" {
		log.Println("JWT_SECRET is not set, using default value")

		secret = "secret"
	}

	return &ginjwt.GinJWTMiddleware{
		Realm:       realm,
		Key:         []byte(secret),
		Timeout:     time.Duration(timeout) * time.Minute,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(repo),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,

		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			user, err := c.Get(IdentityKey)
			var userObj *models.User = nil
			if err {
				userObj = user.(*models.User)
			}
			loginReponse(c, code, message, time, userObj)
		},
	}
}

func authenticator(repo *sqliteRepo.Repository) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", ginjwt.ErrMissingLoginValues
		}

		login := strings.TrimSpace(loginVals.Login)
		password := strings.TrimSpace(loginVals.Password)

		user, err := repo.GetUserByLoginAndPassword(login, password)
		if err == nil {
			c.Set(IdentityKey, user) // Set the user in the context

			return user, nil
		}

		return nil, ginjwt.ErrFailedAuthentication
	}
}

func payloadFunc() func(data interface{}) ginjwt.MapClaims {
	return func(data interface{}) ginjwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return ginjwt.MapClaims{
				IdentityKey: v.ID,
			}
		}

		return ginjwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := ginjwt.ExtractClaims(c)

		return &models.User{
			ID: uint(claims[IdentityKey].(float64)),
		}
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if _, ok := data.(*models.User); ok {
			return true
		}

		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func loginReponse(c *gin.Context, code int, token string, expire time.Time, user *models.User) {
	loginResponse := loginResponse{
		Code:   http.StatusOK,
		Token:  token,
		Expire: expire.Format(time.RFC3339),
	}

	if user != nil {
		role := user.Role()
		loginResponse.Role = &role

		username := user.Username
		loginResponse.Username = &username

		gravatarUrl := user.GravatarURL()
		loginResponse.GravatarUrl = &gravatarUrl

		id := user.ID
		loginResponse.Id = &id
	}

	c.JSON(code, loginResponse)
}
