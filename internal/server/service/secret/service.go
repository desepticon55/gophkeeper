package secret

import (
	"context"
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/internal/server"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type SecretService struct {
	logger     *zap.Logger
	repository secretRepository
}

func NewSecretService(l *zap.Logger, r secretRepository) *SecretService {
	return &SecretService{logger: l, repository: r}
}

func (s *SecretService) CreateSecret(ctx context.Context, secret model.Secret) error {
	currentUserName := fmt.Sprintf("%v", ctx.Value(server.UserNameContextKey))
	if secret.Name == "" {
		return model.ErrSecretNameIsEmpty
	}

	existSecret, err := s.repository.ExistSecret(ctx, currentUserName, secret.Name)
	if err != nil {
		s.logger.Error("Error during create secret", zap.String("name", secret.Name), zap.String("userName", currentUserName), zap.Error(err))
		return err
	}

	if existSecret {
		return model.ErrSecretExistToCurrentUser
	}

	err = s.repository.CreateSecret(ctx, currentUserName, secret)
	if err != nil {
		s.logger.Error("Error during create secret", zap.String("name", secret.Name), zap.String("userName", currentUserName), zap.Error(err))
		return err
	}

	return nil
}

func (s *SecretService) FindSecret(ctx context.Context, secretName string) (model.Secret, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(server.UserNameContextKey))
	secret, err := s.repository.FindSecret(ctx, currentUserName, secretName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("Secret was not found", zap.String("name", secretName), zap.String("userName", currentUserName))
			return model.Secret{}, model.ErrSecretWasNotFound
		}

		s.logger.Error("Error during find secret", zap.String("name", secretName), zap.String("userName", currentUserName), zap.Error(err))
		return model.Secret{}, err
	}
	return secret, nil
}

func (s *SecretService) FindAllSecrets(ctx context.Context) ([]model.Secret, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(server.UserNameContextKey))
	secrets, err := s.repository.FindAllSecrets(ctx, currentUserName)
	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Warn("Secrets was not found", zap.String("userName", currentUserName))
		return nil, model.ErrSecretsWasNotFound
	}
	return secrets, nil
}

func (s *SecretService) DeleteSecret(ctx context.Context, secretName string) error {
	currentUserName := fmt.Sprintf("%v", ctx.Value(server.UserNameContextKey))
	err := s.repository.DeleteSecret(ctx, currentUserName, secretName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("Secret was not found", zap.String("name", secretName), zap.String("userName", currentUserName))
			return model.ErrSecretWasNotFound
		}
		return err
	}
	return nil
}
