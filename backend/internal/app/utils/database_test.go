package utils

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestConnectToDatabaseInvalidFilename(t *testing.T) {
	assert.PanicsWithValue(t, "invalid database filename: \"invalid/name\"", func() {
		ConnectToDatabase("invalid/name")
	})

	assert.PanicsWithValue(t, "invalid database filename: \"name with spaces\"", func() {
		ConnectToDatabase("name with spaces")
	})

	assert.PanicsWithValue(t, "invalid database filename: \"../escaped\"", func() {
		ConnectToDatabase("../escaped")
	})
}

func TestConnectToDatabaseValidFilename(t *testing.T) {
	// Important: We need to ensure the "data" directory exists for the test, otherwise
	// ConnectToDatabase will panic when trying to create the SQLite file.
	err := os.MkdirAll("data", 0o755)
	require.NoError(t, err)

	defer os.RemoveAll("data") // Clean up after test

	var db *gorm.DB

	assert.NotPanics(t, func() {
		db = ConnectToDatabase("testdb_123")
	})
	assert.NotNil(t, db)

	assert.NotPanics(t, func() {
		db = ConnectToDatabase("")
	})
	assert.NotNil(t, db)
}

func TestConnectToLocalDatabase(t *testing.T) {
	var db *gorm.DB

	assert.NotPanics(t, func() {
		db = ConnectToLocalDatabase()
	})
	assert.NotNil(t, db)
}

func TestMigrateAndPurgeDatabase(t *testing.T) {
	db := ConnectToLocalDatabase()
	require.NotNil(t, db)

	assert.NotPanics(t, func() {
		MigrateDatabase(db)
	})

	assert.True(t, db.Migrator().HasTable(&models.User{}), "User table should exist after migration")
	assert.True(t, db.Migrator().HasTable(&models.Product{}), "Product table should exist after migration")

	// 2. Purge
	assert.NotPanics(t, func() {
		PurgeDatabase(db)
	})

	// Check if the tables were actually dropped
	assert.False(t, db.Migrator().HasTable(&models.User{}), "User table should be dropped after purge")
	assert.False(t, db.Migrator().HasTable(&models.Product{}), "Product table should be dropped after purge")
}

func TestSeedDatabase(t *testing.T) {
	db := ConnectToLocalDatabase()
	require.NotNil(t, db)

	MigrateDatabase(db)

	assert.NotPanics(t, func() {
		SeedDatabase(db, true) // Test with includeTestData = true
	})

	assert.NotPanics(t, func() {
		SeedDatabase(db, false) // Test with includeTestData = false
	})
}
