package guestlist

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/server/http/errors"
	"github.com/potibm/kasseapparat/internal/app/service"
)

type httpHandler struct {
	guestlistService service.GuestlistService
}

func (h httpHandler) findByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := h.guestlistService.FindByID(ctx, id)
	if err != nil {
		_ = c.Error(errors.ExtendHttpErrorWithDetails(errors.NotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, list)
}

func newHTTPHandler(guestlistService service.GuestlistService) *httpHandler {
	return &httpHandler{
		guestlistService: guestlistService,
	}
}
