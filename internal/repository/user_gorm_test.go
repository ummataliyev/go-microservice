//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"
	"go-microservice/internal/repository"
	"go-microservice/internal/testutil"
)

func TestCreateUser_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupTestDB(t, db) })

	repo := repository.NewGORMUser(db)
	ctx := context.Background()

	user := &models.User{
		Email:          "create@example.com",
		HashedPassword: "$2a$10$somethinghashed",
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.Greater(t, user.ID, uint(0))
}

func TestGetByEmail_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupTestDB(t, db) })

	repo := repository.NewGORMUser(db)
	ctx := context.Background()

	user := &models.User{
		Email:          "lookup@example.com",
		HashedPassword: "$2a$10$somethinghashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.GetByEmail(ctx, "lookup@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, "lookup@example.com", found.Email)
	assert.Equal(t, user.HashedPassword, found.HashedPassword)
}

func TestGetAll_Pagination_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupTestDB(t, db) })

	repo := repository.NewGORMUser(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		u := &models.User{
			Email:          fmt.Sprintf("page%d@example.com", i),
			HashedPassword: "$2a$10$somethinghashed",
		}
		require.NoError(t, repo.Create(ctx, u))
	}

	users, err := repo.GetAll(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestSoftDelete_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupTestDB(t, db) })

	repo := repository.NewGORMUser(db)
	ctx := context.Background()

	user := &models.User{
		Email:          "delete@example.com",
		HashedPassword: "$2a$10$somethinghashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	// Soft-delete the user.
	err := repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// GetByID should return ErrNotFound for a soft-deleted user.
	_, err = repo.GetByID(ctx, user.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoerrors.ErrNotFound)

	// Restore the user.
	err = repo.Restore(ctx, user.ID)
	require.NoError(t, err)

	// GetByID should succeed after restore.
	restored, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, restored.ID)
}

func TestCount_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupTestDB(t, db) })

	repo := repository.NewGORMUser(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		u := &models.User{
			Email:          fmt.Sprintf("count%d@example.com", i),
			HashedPassword: "$2a$10$somethinghashed",
		}
		require.NoError(t, repo.Create(ctx, u))
	}

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}
