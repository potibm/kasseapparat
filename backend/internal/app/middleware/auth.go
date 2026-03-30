package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v3"
	ginjwtCore "github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/potibm/kasseapparat/internal/app/exitcode"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	RefreshTokenLifetime = 7 * 24 * time.Hour
	loginEndpoint        = "/auth/login"
	refreshEndpoint      = "/auth/refresh"
	logoutEndpoint       = "/auth/logout"
)

var IdentityKey = "ID"
var meter = otel.Meter("kasseapparat-auth")

var (
	authEventsCounter, _ = meter.Int64Counter("kasseapparat_auth_events_total",
		metric.WithDescription("Number of authentication events"))
)

type login struct {
	Login    string `json:"login"    form:"login"    binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type loginResponseDTO struct {
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
			slog.Error("Error initializing auth middleware", "error", errInit)
			os.Exit(int(exitcode.Software))
		}
	}
}

func RegisterRoute(r *gin.RouterGroup, handle *ginjwt.GinJWTMiddleware) {
	r.POST(loginEndpoint, handle.LoginHandler)
	r.POST(refreshEndpoint, func(c *gin.Context) {
		handle.RefreshHandler(c)

		if c.Writer.Status() == http.StatusOK {
			authEventsCounter.Add(c.Request.Context(), 1,
				metric.WithAttributes(
					attribute.String("event_type", "refresh"),
					attribute.String("status", "success"),
				),
			)
		}
	})
	r.POST(logoutEndpoint, func(c *gin.Context) {
		authEventsCounter.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("event_type", "logout"),
				attribute.String("status", "success"),
			),
		)

		handle.LogoutHandler(c)
	})
}

func InitParams(
	repo *sqliteRepo.Repository,
	realm string,
	secret string,
	timeout int,
	secureCookie bool,
) *ginjwt.GinJWTMiddleware {
	if secret == "" {
		slog.Warn("JWT_SECRET is not set, using default value")

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
		LoginResponse:   loginResponse,
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
		eventType := "request"

		if strings.Contains(c.Request.URL.Path, loginEndpoint) {
			eventType = "login"
		} else if strings.Contains(c.Request.URL.Path, refreshEndpoint) {
			eventType = "refresh"
		}

		authEventsCounter.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("event_type", eventType),
				attribute.String("status", "failure"),
				attribute.Int("code", code),
			),
		)

		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func loginResponse(c *gin.Context, token *ginjwtCore.Token) {
	authEventsCounter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("event_type", "login"),
			attribute.String("status", "success"),
		),
	)

	user, exists := c.Get(IdentityKey)

	var userObj *models.User = nil
	if exists {
		userObj = user.(*models.User)
	}

	loginResponse := loginResponseDTO{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   token.ExpiresIn(),
	}

	if userObj != nil {
		role := userObj.Role()
		loginResponse.Role = &role

		username := userObj.Username
		loginResponse.Username = &username

		gravatarUrl := userObj.GravatarURL()
		loginResponse.GravatarUrl = &gravatarUrl

		id := userObj.ID
		loginResponse.Id = &id
	}

	c.JSON(http.StatusOK, loginResponse)
}
