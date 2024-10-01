package guestlist

import (
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/service"
)

func RegisterRoutes(rg *gin.RouterGroup,
	guestlistService service.GuestlistService) {

	handler := newHTTPHandler(guestlistService) // handlers creating
	//rg.GET("", handler.GetProducts)
	rg.GET("/:id", handler.findByID)
	//rg.GET("/:id/listEntries", handler.GetListEntriesByProductID)
	//rg.PUT("/:id", handler.UpdateProductByID)
	//rg.DELETE("/:id", handler.DeleteProductByID)
	//rg.POST("", handler.CreateProduct)
}
