package secret

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
)

type secretRepository interface {
	CreateSecret(ctx context.Context, userName string, secret model.Secret) error

	FindSecret(ctx context.Context, userName string, secretName string) (model.Secret, error)

	ExistSecret(ctx context.Context, userName string, secretName string) (bool, error)

	FindAllSecrets(ctx context.Context, userName string) ([]model.Secret, error)

	DeleteSecret(ctx context.Context, userName string, secretName string) error
}
