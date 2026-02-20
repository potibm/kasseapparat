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
	err := db.Migrator().
		DropTable(
			&models.Product{},
			&models.Purchase{},
			&models.PurchaseItem{},
			&models.User{},
			&models.Guestlist{},
			&models.Guest{},
			&models.ProductInterest{},
		)
	if err != nil {
		panic(err)
	}
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Product{},
		&models.Purchase{},
		&models.PurchaseItem{},
		&models.User{},
		models.Guestlist{},
		models.Guest{},
		models.ProductInterest{},
	)
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB, includeTestData bool) {
	seed := NewDatabaseSeed(db)
	seed.Seed(includeTestData)
}
