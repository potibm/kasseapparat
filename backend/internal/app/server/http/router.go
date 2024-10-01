package http

import (
	"github.com/potibm/kasseapparat/internal/app/server/http/handler/guestlist"
)

func (s *Server) registerRoutes() {

	guestlistGroup := s.r.Group("/api/v1/guestlist")
	guestlist.RegisterRoutes(guestlistGroup, s.services.GuestlistService)
}
