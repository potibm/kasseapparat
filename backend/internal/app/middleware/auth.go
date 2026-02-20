package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v3"
	ginjwtCore "github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

const RefreshTokenLifetime = 7 * time.Hour * 24

var IdentityKey = "ID"

type login struct {
	Login    string `json:"login"    form:"login"    binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	Role         *string `json:"role"`
	Username     *string `json:"username"`
	GravatarUrl  *string `json:"gravatarUrl"`
	Id           *uint   `json:"id"`
}

func HandlerMiddleWare(authMiddleware *ginjwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func RegisterRoute(r *gin.RouterGroup, handle *ginjwt.GinJWTMiddleware) {
	r.POST("/auth/login", handle.LoginHandler)
	r.POST("/auth/refresh", handle.RefreshHandler)
	r.POST("/auth/logout", handle.LogoutHandler)
}

func InitParams(
	repo *sqliteRepo.Repository,
	realm string,
	secret string,
	timeout int,
	secureCookie bool,
) *ginjwt.GinJWTMiddleware {
	if secret == "" {
		log.Println("JWT_SECRET is not set, using default value")

		secret = "secret"
	}

	return &ginjwt.GinJWTMiddleware{
		Realm:      realm,
		Key:        []byte(secret),
		Timeout:    time.Minute * time.Duration(timeout), // Short-lived access tokens
		MaxRefresh: RefreshTokenLifetime,

		SecureCookie:   secureCookie,            // HTTPS only
		CookieHTTPOnly: true,                    // Prevent XSS
		CookieSameSite: http.SameSiteStrictMode, // CSRF protection
		SendCookie:     true,                    // Enable secure cookies

		IdentityKey:     IdentityKey,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(repo),
		Authorizer:      authorizer(),
		Unauthorized:    unauthorized(),

		LoginResponse: func(c *gin.Context, token *ginjwtCore.Token) {
			user, err := c.Get(IdentityKey)

			var userObj *models.User = nil
			if err {
				userObj = user.(*models.User)
			}

			loginReponse(c, token, userObj)
		},
	}
}

func authenticator(repo *sqliteRepo.Repository) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
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

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return jwt.MapClaims{
				IdentityKey: v.ID,
			}
		}

		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := ginjwt.ExtractClaims(c)

		return &models.User{
			ID: uint(claims[IdentityKey].(float64)),
		}
	}
}

func authorizer() func(c *gin.Context, data any) bool {
	return func(c *gin.Context, data any) bool {
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

func loginReponse(c *gin.Context, token *ginjwtCore.Token, user *models.User) {
	loginResponse := loginResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   token.ExpiresIn(),
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

	c.JSON(http.StatusOK, loginResponse)
}
