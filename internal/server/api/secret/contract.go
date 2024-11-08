package secret

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
)

type secretService interface {
	CreateSecret(ctx context.Context, secret model.Secret) error

	FindSecret(ctx context.Context, name string) (model.Secret, error)

	FindAllSecrets(ctx context.Context) ([]model.Secret, error)

	DeleteSecret(ctx context.Context, name string) error
}
