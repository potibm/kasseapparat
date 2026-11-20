package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
)

// listQueryError reports a failed list query. An unknown sort field comes from the
// client (react-admin sends _sort from the column headers), so it is a 400; anything
// else is a server error.
func listQueryError(c *gin.Context, err error) {
	if errors.Is(err, sqliteRepo.ErrInvalidSortField) {
		_ = c.Error(InvalidRequest.WithMsg("Invalid sort field").WithCause(err))

		return
	}

	_ = c.Error(InternalServerError.WithCauseMsg(err))
}
