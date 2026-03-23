package utils

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/potibm/kasseapparat/internal/app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDatabase(filename string) *gorm.DB {
	if filename == "" {
		filename = "kasseapparat"
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validName.MatchString(filename) {
		panic(fmt.Sprintf("invalid database filename: %q", filename))
	}

	dbPath := filepath.Join("data", filename+".db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
		&models.Guestlist{},
		&models.Guest{},
		&models.ProductInterest{},
	)
	if err != nil {
		panic(err)
	}
}

func SeedDatabase(db *gorm.DB, includeTestData bool) {
	seed := NewDatabaseSeed(db)
	seed.Seed(includeTestData)
}
