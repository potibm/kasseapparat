package helper

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/storage"
)

func QueryArrayUint(c *gin.Context, field string) []uint {

	idStrings := c.QueryArray(field)
	var ids []uint

	for _, s := range idStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("Error converting %s to int: %v", s, err)
		}
		ids = append(ids, uint(id))
	}

	return ids
}

func NewQueryOptions(ctx *gin.Context) storage.QueryOptions {
	queryOptions := storage.QueryOptions{}
	queryOptions.Offset, _ = strconv.Atoi(ctx.DefaultQuery("_start", "0"))
	queryOptions.Limit, _ = strconv.Atoi(ctx.DefaultQuery("_end", "10"))
	queryOptions.SortBy = ctx.DefaultQuery("_sort", "id")
	queryOptions.SortAsc = ctx.DefaultQuery("_order", "ASC") == "ASC"

	return queryOptions
}
