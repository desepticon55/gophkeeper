package storage

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) ExistUser(ctx context.Context, userName string) (bool, error) {
	var exist bool
	query := "select exists(select 1 from gophkeeper.user where username = $1)"
	err := r.pool.QueryRow(ctx, query, userName).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, userName string, password string) error {
	query := "insert into gophkeeper.user(username, password) values ($1, $2)"
	_, err := r.pool.Exec(ctx, query, userName, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindUser(ctx context.Context, userName string) (model.User, error) {
	var user model.User
	query := "select username, password from gophkeeper.user where username = $1"
	err := r.pool.QueryRow(ctx, query, userName).Scan(&user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
