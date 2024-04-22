package utils

import (
	"github.com/potibm/die-kassa/internal/app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	// Your own implementation of connecting to the database
	db, err := gorm.Open(sqlite.Open("./data/kassa.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func PurgeDatabase(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Product{})
	if err != nil {
		panic(err)
	}
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB) {
	// Your own implementation of seeding the database
	db.Create(&models.Product{Name: "🎟️ Regular", Price: 40, Pos: 1, WrapAfter: false})
	db.Create(&models.Product{Name: "🎟️ Reduced", Price: 20, Pos: 2, WrapAfter: false})
	db.Create(&models.Product{Name: "🎟️ Free", Price: 0, Pos: 3, WrapAfter: true})
	db.Create(&models.Product{Name: "👕 T-Shirt Male S", Price: 20, Pos: 10, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Male M", Price: 20, Pos: 11, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Male L", Price: 20, Pos: 12, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Male XL", Price: 20, Pos: 13, WrapAfter: true})
	db.Create(&models.Product{Name: "👕 T-Shirt Female S", Price: 20, Pos: 20, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Female M", Price: 20, Pos: 21, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Female L", Price: 20, Pos: 22, WrapAfter: false})
	db.Create(&models.Product{Name: "👕 T-Shirt Female XL", Price: 20, Pos: 23, WrapAfter: true})
	db.Create(&models.Product{Name: "☕ Coffee Mug", Price: 1, Pos: 30, WrapAfter: false})
}
