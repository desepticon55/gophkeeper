package storage

import (
	"context"
	"github.com/desepticon55/gophkeeper/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	pool, cleanup := utils.InitPostgresIntegrationTest(t, ctx, logger)

	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatalf("failed to cleanup test database: %s", err)
		}
	})

	userRepository := NewUserRepository(pool)

	t.Run("ExistUser", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := userRepository.CreateUser(ctx, "testUser", "testPassword")
		assert.NoError(t, err)

		result, err := userRepository.ExistUser(ctx, "testUser")
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("CreateUser", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := userRepository.CreateUser(ctx, "testUser", "testPassword")
		assert.NoError(t, err)

		var countUsers int
		err = pool.QueryRow(ctx, `select count(1) from gophkeeper.user where username = $1`, "testUser").Scan(&countUsers)
		assert.NoError(t, err)
		assert.Equal(t, 1, countUsers)
	})

	t.Run("FindUser", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})
		err := userRepository.CreateUser(ctx, "testUser", "testPassword")
		assert.NoError(t, err)

		result, err := userRepository.FindUser(ctx, "testUser")
		assert.NoError(t, err)
		assert.Equal(t, "testUser", result.Username)
		assert.Equal(t, "testPassword", result.Password)
	})
}
