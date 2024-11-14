package auth

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
)

type userService interface {
	CreateUser(ctx context.Context, user model.User) error

	FindUser(ctx context.Context, user model.User) (model.User, error)
}
