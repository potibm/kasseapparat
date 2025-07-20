package initializer

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

func InitializeJwtMiddleware(repository *sqliteRepo.Repository, jwtConfig config.JwtConfig) *jwt.GinJWTMiddleware {
	const timeout = 10 // Duration that a jwt token is valid, in minutes

	jwtMiddleware, err := jwt.New(middleware.InitParams(repository, jwtConfig.Realm, jwtConfig.Secret, timeout))
	if err != nil {
		panic("[Error] failed to initialize JWT middleware: " + err.Error())
	}

	return jwtMiddleware
}
