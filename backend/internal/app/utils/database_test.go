package utils

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestConnectToDatabaseBoundsConnectionPool(t *testing.T) {
	err := os.MkdirAll("data", 0o755)
	require.NoError(t, err)

	defer os.RemoveAll("data")

	db, err := ConnectToDatabase("testdb_pool")
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	assert.Equal(t, 1, sqlDB.Stats().MaxOpenConnections)
}

// The single-connection bound deadlocks if a transaction callback reaches for the
// outer repository instead of the transaction-scoped one, so pin that behaviour.
func TestConnectToDatabaseTransactionDoesNotDeadlock(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("testdb_tx")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := &models.Product{Name: "Pool Test", NetPrice: decimal.NewFromInt(10)}

		return tx.Create(product).Error
	})

	require.NoError(t, err)
}

// The seed is the one transaction that does not go through WithTransaction, and it
// runs on the bounded pool in CI (`mise run e2e:setup` -> `database reset
// --test-data`). A reach-back there deadlocks, so run it on the real pool under a
// watchdog: a regression fails the test instead of hanging the job until the
// workflow timeout.
func TestSeedDatabaseOnBoundedPoolDoesNotDeadlock(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("testdb_seed")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.Equal(t, 1, sqlDB.Stats().MaxOpenConnections, "seed must run on the bounded pool")

	done := make(chan error, 1)

	go func() {
		done <- SeedDatabase(db, true)
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
		t.Fatal("SeedDatabase did not finish: a transaction callback waits for a connection the pool cannot give it")
	}

	// A silent rollback would leave the fixture database empty, so assert the
	// transaction actually committed its purchases.
	var purchaseCount int64

	require.NoError(t, db.Model(&models.Purchase{}).Count(&purchaseCount).Error)
	assert.Positive(t, purchaseCount, "seeded purchases should have been committed")
}

// Reaching for the outer handle from inside a transaction can only ever end in a
// context error, never in a connection. Bound it so the failure is an assertion
// rather than a hang.
func TestBoundedPoolTransactionCallbackUsingOuterHandleErrors(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("testdb_outer")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
		product := &models.Product{Name: "Outer Handle", NetPrice: decimal.NewFromInt(10)}

		return db.WithContext(ctx).Create(product).Error
	})

	require.ErrorIs(t, err, context.DeadlineExceeded)
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

	assert.True(t, db.Migrator().HasTable(&models.Product{}), "Product table should exist after migration")

	// 2. Purge
	err = PurgeDatabase(db)
	assert.NoError(t, err)

	// Check if the tables were actually dropped
	assert.False(t, db.Migrator().HasTable(&models.Product{}), "Product table should be dropped after purge")
}

func TestSeedDatabase(t *testing.T) {
	db, err := ConnectToLocalDatabase()
	require.NoError(t, err)
	require.NotNil(t, db)

	require.NoError(t, SeedDatabase(db, true), "seeding with test data should succeed")

	var purchaseCount int64

	require.NoError(t, db.Model(&models.Purchase{}).Count(&purchaseCount).Error)
	assert.Positive(t, purchaseCount, "test data should seed purchases")

	require.NoError(t, SeedDatabase(db, false), "seeding without test data should succeed")
}

// A failed seed used to be swallowed, so `database reset` reported success over a
// half-written fixture database. Drop the table the first seeding step writes to
// and pin that the failure reaches the caller.
func TestSeedDatabaseReportsFailure(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("testdb_seedfail")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	require.NoError(t, db.Migrator().DropTable(&models.Product{}))

	err = SeedDatabase(db, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to seed database")
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
