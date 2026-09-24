package sqlite

import (
	"context"
	"errors"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := utils.ConnectToLocalDatabase()
	require.NoError(t, err)

	return NewRepository(db, 2)
}

func TestNewRepository(t *testing.T) {
	db, err := utils.ConnectToLocalDatabase()
	require.NoError(t, err)

	repo := NewRepository(db, 3)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
	assert.Equal(t, int32(3), repo.decimalPlaces)
}

func TestGetDB(t *testing.T) {
	repo := setupTestRepository(t)

	db := repo.GetDB()

	assert.NotNil(t, db)
	assert.Equal(t, repo.db, db)
}

func TestPing_Success(t *testing.T) {
	repo := setupTestRepository(t)

	err := repo.Ping()

	assert.NoError(t, err)
}

func TestPing_DatabaseClosed(t *testing.T) {
	repo := setupTestRepository(t)

	sqlDB, err := repo.db.DB()
	require.NoError(t, err)

	err = sqlDB.Close()
	require.NoError(t, err)

	err = repo.Ping()
	assert.Error(t, err)
}

func TestWithTransaction_Commit(t *testing.T) {
	repo := setupTestRepository(t)

	err := repo.WithTransaction(context.Background(), func(txRepo RepositoryInterface) error {
		assert.NotNil(t, txRepo)
		assert.NotNil(t, txRepo.GetDB())

		return nil
	})

	assert.NoError(t, err)
}

func TestWithTransaction_Rollback(t *testing.T) {
	repo := setupTestRepository(t)

	expectedErr := errors.New("transaction failed")

	err := repo.WithTransaction(context.Background(), func(txRepo RepositoryInterface) error {
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestWithTransaction_ContextCancellation(t *testing.T) {
	repo := setupTestRepository(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.WithTransaction(ctx, func(txRepo RepositoryInterface) error {
		return nil
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestCloneWithDB(t *testing.T) {
	repo := setupTestRepository(t)

	newDB, err := utils.ConnectToLocalDatabase()
	require.NoError(t, err)

	cloned := repo.cloneWithDB(newDB)

	assert.NotNil(t, cloned)
	assert.Equal(t, newDB, cloned.db)
	assert.Equal(t, repo.decimalPlaces, cloned.decimalPlaces)
	assert.NotSame(t, repo, cloned)
}
