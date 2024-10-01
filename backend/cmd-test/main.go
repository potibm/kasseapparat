package main

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/server/http"
	"github.com/potibm/kasseapparat/internal/app/service/domain"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite"
	"github.com/potibm/kasseapparat/internal/app/storage/sqlite/guestlist"

	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()
	db, _ := gorm.Open(gormsqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(guestlist.GuestlistModel{})

	db.Create(&guestlist.GuestlistModel{Name: "Test Guestlist"})
	db.Create(&guestlist.GuestlistModel{Name: "Second Guestlist"})
	db.Create(&guestlist.GuestlistModel{Name: "Third Guestlist", TypeCode: true})

	repositories := sqlite.NewRepository(db)
	services := domain.NewService(repositories)

	httpServer := http.NewServer(services)

	err := httpServer.Serve(ctx)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
