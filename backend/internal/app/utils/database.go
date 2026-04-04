package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func IsValidDatabaseFilename(filename string) bool {
	validName := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return validName.MatchString(filename)
}

func ConnectToDatabase(filename string) (*gorm.DB, error) {
	if filename == "" {
		filename = "kasseapparat"
	}

	if !IsValidDatabaseFilename(filename) {
		return nil, fmt.Errorf("invalid database filename: %q", filename)
	}

	dbPath := filepath.Join("data", filename+".db")

	dbDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create database directory: %w", err)
	}

	db, err := connectToSQLite(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func ConnectToLocalDatabase() (*gorm.DB, error) {
	db, err := connectToSQLite("file::memory:?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to local database: %w", err)
	}

	return db, nil
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

func PurgeDatabase(db *gorm.DB) error {
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
		return fmt.Errorf("failed to purge database: %w", err)
	}
	return nil
}

func MigrateDatabase(db *gorm.DB) error {
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
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}

func SeedDatabase(db *gorm.DB, includeTestData bool) {
	seed := NewDatabaseSeed(db)
	seed.Seed(includeTestData)
}
