package utils

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
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

	db, err := connectToSQLite(dbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	return db
}

func ConnectToLocalDatabase() *gorm.DB {
	db, err := connectToSQLite("file::memory:?cache=shared")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to local database: %v", err))
	}

	return db
}

func connectToSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
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
