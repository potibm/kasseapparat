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
	// Your own implementation of purging the database
	db.Migrator().DropTable(&models.Product{})
}

func MigrateDatabase(db *gorm.DB) {
	// Your own implementation of migrating the database
	db.AutoMigrate(&models.Product{})
}

func SeedDatabase(db *gorm.DB) {
	// Your own implementation of seeding the database
	db.Create(&models.Product{Name: "🎟️ Regular", Price: 40})
	db.Create(&models.Product{Name: "🎟️ Reduced", Price: 20})
	db.Create(&models.Product{Name: "🎟️ Free", Price: 0})
	db.Create(&models.Product{Name: "👕 T-Shirt Male S", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Male M", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Male L", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Male XL", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female S", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female M", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female L", Price: 20})
	db.Create(&models.Product{Name: "👕 T-Shirt Female XL", Price: 20})
	db.Create(&models.Product{Name: "☕ Coffee Mug", Price: 1})
}
