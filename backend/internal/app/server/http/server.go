package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/server/http/middleware"
	"github.com/potibm/kasseapparat/internal/app/service"
)

type Server struct {
	services *service.Service
	r        *gin.Engine
}

func NewServer(services *service.Service) *Server {
	return &Server{
		services: services,
		r:        gin.Default(),
	}
}

func (s *Server) Serve(ctx context.Context) error {
	s.r.Use(middleware.ErrorHandlingMiddleware())

	s.registerRoutes()

	// Start the server
	return s.r.Run()
}
