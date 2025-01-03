package guestlist

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/server/http/errors"
	"github.com/potibm/kasseapparat/internal/app/server/http/helper"
	"github.com/potibm/kasseapparat/internal/app/service"
	"github.com/potibm/kasseapparat/internal/app/storage"
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

	listResponse := newGuestlistReponse(*list)
	c.JSON(http.StatusOK, listResponse)
}

func (h httpHandler) find(c *gin.Context) {
	ctx := c.Request.Context()

	queryOptions := helper.NewQueryOptions(c)

	filters := storage.GuestListFilters{}
	filters.Query = c.DefaultQuery("q", "")
	filters.IDs = helper.QueryArrayUint(c, "id")

	lists, err := h.guestlistService.FindAllWithParams(ctx, queryOptions, filters)
	if err != nil {
		_ = c.Error(errors.ExtendHttpErrorWithDetails(errors.InternalServerError, err.Error()))
	}

	total, err := h.guestlistService.GetTotalCount(ctx)
	if err != nil {
		_ = c.Error(errors.InternalServerError)
		return
	}

	reponseLists := make([]GuestlistResponse, 0)
	for _, list := range lists {
		reponseLists = append(reponseLists, newGuestlistReponse(*list))
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	c.JSON(http.StatusOK, reponseLists)
}

func newHTTPHandler(guestlistService service.GuestlistService) *httpHandler {
	return &httpHandler{
		guestlistService: guestlistService,
	}
}
