package utils

import (
	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./data/kasseapparat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func ConnectToLocalDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func PurgeDatabase(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{})
	if err != nil {
		panic(err)
	}
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Product{}, &models.Purchase{}, &models.PurchaseItem{}, &models.User{})
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB) {
	// Your own implementation of seeding the database
	db.Create(&models.Product{Name: "🎟️ Regular", Price: 40, Pos: 1, ApiExport: true})
	db.Create(&models.Product{Name: "🎟️ Reduced", Price: 20, Pos: 2, ApiExport: true})
	db.Create(&models.Product{Name: "🎟️ Free", Price: 0, Pos: 3, WrapAfter: true, ApiExport: true})
	db.Create(&models.Product{Name: "👕 T-Shirt Male S", Price: 20, Pos: 10})
	db.Create(&models.Product{Name: "👕 T-Shirt Male M", Price: 20, Pos: 11})
	db.Create(&models.Product{Name: "👕 T-Shirt Male L", Price: 20, Pos: 12})
	db.Create(&models.Product{Name: "👕 T-Shirt Male XL", Price: 20, Pos: 13, WrapAfter: true})
	db.Create(&models.Product{Name: "👕 T-Shirt Female S", Price: 20, Pos: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female M", Price: 20, Pos: 21})
	db.Create(&models.Product{Name: "👕 T-Shirt Female L", Price: 20, Pos: 22})
	db.Create(&models.Product{Name: "👕 T-Shirt Female XL", Price: 20, Pos: 23, WrapAfter: true})
	db.Create(&models.Product{Name: "☕ Coffee Mug", Price: 1, Pos: 30})
	db.Create(&models.User{Username: "admin", Password: "admin", Admin: true})
	db.Create(&models.User{Username: "demo", Password: "demo", Admin: false})
}
