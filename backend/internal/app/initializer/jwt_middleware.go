package initializer

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/middleware"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

func InitializeJwtMiddleware(repository *sqliteRepo.Repository, jwtConfig config.JwtConfig) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(middleware.InitParams(repository, jwtConfig.Realm, jwtConfig.Secret, 10))
	if err != nil {
		panic("[Error] failed to initialize JWT middleware: " + err.Error())
	}

	return jwtMiddleware
}
