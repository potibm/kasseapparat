package initializer

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/store"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	"github.com/potibm/kasseapparat/internal/app/models"
)

var newJwtFunc = jwt.New

type UserAuthenticator interface {
	GetUserByLoginAndPassword(login, password string) (*models.User, error)
}

func InitializeJwtMiddleware(
	repository UserAuthenticator,
	jwtConfig config.JwtConfig,
	redisConfig *config.RedisURL,
) *jwt.GinJWTMiddleware {
	const timeout = 10 // Duration that a JWT token is valid, in minutes

	var jwtRedisConfig *store.RedisConfig

	if redisConfig != nil {
		cfg := redisConfig.JwtConfig()
		jwtRedisConfig = &cfg
	}

	jwtMiddleware, err := newJwtFunc(
		middleware.InitParams(
			repository,
			jwtConfig.Realm,
			jwtConfig.Secret,
			timeout,
			jwtConfig.SecureCookie,
			jwtRedisConfig,
		),
	)
	if err != nil {
		panic("[Error] failed to initialize JWT middleware: " + err.Error())
	}

	return jwtMiddleware
}
