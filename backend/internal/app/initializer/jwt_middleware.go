package initializer

import (
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/middleware"
)

func InitializeAuthMiddleware() gin.HandlerFunc {
	return middleware.HandlerMiddleWare()
}
