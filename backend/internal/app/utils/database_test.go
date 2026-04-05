package utils

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabaseInvalidFilename(t *testing.T) {
	invalidFilenames := []string{
		"invalid/name",
		"name with spaces",
		"name/with/slash",
		"../traversal",
		"invalid\\name",
		"invalid:name",
	}

	for _, filename := range invalidFilenames {
		t.Run(filename, func(t *testing.T) {
			db, err := ConnectToDatabase(filename)
			assert.Nil(t, db)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "invalid database filename")
		})
	}
}

func TestIsValidDatabaseFilename(t *testing.T) {
	validFilenames := []string{
		"validname",
		"valid_name",
		"valid-name",
		"valid.name",
		"12345",
	}

	for _, filename := range validFilenames {
		t.Run(filename, func(t *testing.T) {
			assert.True(t, IsValidDatabaseFilename(filename))
		})
	}

	invalidFilenames := []string{
		"invalid/name",
		"name with spaces",
		"name/with/slash",
		"../traversal",
		"invalid\\name",
		"invalid:name",
	}

	for _, filename := range invalidFilenames {
		t.Run(filename, func(t *testing.T) {
			assert.False(t, IsValidDatabaseFilename(filename))
		})
	}
}

func TestConnectToDatabaseValidFilename(t *testing.T) {
	// Important: We need to ensure the "data" directory exists for the test, otherwise
	// ConnectToDatabase will panic when trying to create the SQLite file.
	err := os.MkdirAll("data", 0o755)
	require.NoError(t, err)

	defer os.RemoveAll("data") // Clean up after test

	db, err := ConnectToDatabase("testdb_123")
	assert.NotNil(t, db)
	assert.NoError(t, err)

	db, err = ConnectToDatabase("")
	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestConnectToLocalDatabase(t *testing.T) {
	db, err := ConnectToLocalDatabase()

	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestMigrateAndPurgeDatabase(t *testing.T) {
	db, err := ConnectToLocalDatabase()
	require.NoError(t, err)
	require.NotNil(t, db)

	err = MigrateDatabase(db)
	assert.NoError(t, err)

	assert.True(t, db.Migrator().HasTable(&models.User{}), "User table should exist after migration")
	assert.True(t, db.Migrator().HasTable(&models.Product{}), "Product table should exist after migration")

	// 2. Purge
	err = PurgeDatabase(db)
	assert.NoError(t, err)

	// Check if the tables were actually dropped
	assert.False(t, db.Migrator().HasTable(&models.User{}), "User table should be dropped after purge")
	assert.False(t, db.Migrator().HasTable(&models.Product{}), "Product table should be dropped after purge")
}

func TestSeedDatabase(t *testing.T) {
	db, err := ConnectToLocalDatabase()
	require.NoError(t, err)
	require.NotNil(t, db)

	err = MigrateDatabase(db)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		SeedDatabase(db, true) // Test with includeTestData = true
	})

	assert.NotPanics(t, func() {
		SeedDatabase(db, false) // Test with includeTestData = false
	})
}

func TestCloseDatabase(t *testing.T) {
	db, err := ConnectToLocalDatabase()
	require.NoError(t, err)

	err = CloseDatabase(db)
	assert.NoError(t, err)

	err = db.Exec("SELECT 1").Error
	assert.Error(t, err, "Operations should fail after closing the database")
}

func TestConnectToDatabaseDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	defer func() {
		_ = os.Chdir(originalWd)
	}()

	db, err := ConnectToDatabase("new_db")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	info, err := os.Stat("data")
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	_ = CloseDatabase(db)
}
