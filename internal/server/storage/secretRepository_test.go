package storage

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestSecretRepository(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	pool, cleanup := utils.InitPostgresIntegrationTest(t, ctx, logger)

	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatalf("failed to cleanup test database: %s", err)
		}
	})

	secretRepository := NewSecretRepository(pool)

	t.Run("ExistSecret", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := secretRepository.CreateSecret(ctx, "testUser", model.Secret{
			Name:    "testName",
			Type:    model.CredentialsSecretType,
			Content: []byte("Hello"),
		})
		assert.NoError(t, err)

		result, err := secretRepository.ExistSecret(ctx, "testUser", "testName")
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("CreateSecret", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := secretRepository.CreateSecret(ctx, "testUser", model.Secret{
			Name:    "testName",
			Type:    model.CredentialsSecretType,
			Content: []byte("Hello"),
		})
		assert.NoError(t, err)

		var countUsers int
		err = pool.QueryRow(ctx, `select count(1) from gophkeeper.secret where username = $1 and name = $2`, "testUser", "testName").Scan(&countUsers)
		assert.NoError(t, err)
		assert.Equal(t, 1, countUsers)
	})

	t.Run("FindSecret", func(t *testing.T) {
		t.Cleanup(func() {
			if err := utils.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})
		err := secretRepository.CreateSecret(ctx, "testUser", model.Secret{
			Name:    "testName",
			Type:    model.CredentialsSecretType,
			Content: []byte("Hello"),
		})
		assert.NoError(t, err)

		result, err := secretRepository.FindSecret(ctx, "testUser", "testName")
		assert.NoError(t, err)
		assert.Equal(t, "testName", result.Name)
		assert.Equal(t, "testUser", result.Username)
		assert.Equal(t, model.CredentialsSecretType, result.Type)
		assert.Equal(t, []byte("Hello"), result.Content)
	})
}
