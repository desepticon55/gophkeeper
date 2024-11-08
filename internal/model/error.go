package model

import "errors"

var (
	ErrUserDataIsNotValid       = errors.New("user data is not valid")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrSecretWasNotFound        = errors.New("secret with specific name to current user was not found")
	ErrSecretsWasNotFound       = errors.New("secrets to current user was not found")
	ErrSecretNameIsEmpty        = errors.New("secret name is empty")
	ErrSecretExistToCurrentUser = errors.New("secret already exists to current user")
)
