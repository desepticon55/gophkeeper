package storage

import (
	"context"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SecretRepository struct {
	pool *pgxpool.Pool
}

func NewSecretRepository(pool *pgxpool.Pool) *SecretRepository {
	return &SecretRepository{
		pool: pool,
	}
}

func (r *SecretRepository) ExistSecret(ctx context.Context, userName string, secretName string) (bool, error) {
	var exist bool
	query := "select exists(select 1 from gophkeeper.secret where username = $1 and name = $2)"
	err := r.pool.QueryRow(ctx, query, userName, secretName).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r *SecretRepository) CreateSecret(ctx context.Context, userName string, secret model.Secret) error {
	query := "insert into gophkeeper.secret(name, username, content, type, opt_lock) values ($1, $2, $3, $4, $5)"
	_, err := r.pool.Exec(ctx, query, secret.Name, userName, secret.Content, secret.Type, 0)
	if err != nil {
		return err
	}

	return nil
}

func (r *SecretRepository) FindSecret(ctx context.Context, userName string, secretName string) (model.Secret, error) {
	var secret model.Secret
	query := "select name, username, content, type, opt_lock from gophkeeper.secret where username = $1 and name = $2"
	err := r.pool.QueryRow(ctx, query, userName, secretName).Scan(&secret.Name, &secret.Username, &secret.Content, &secret.Type, &secret.Version)
	if err != nil {
		return model.Secret{}, err
	}

	return secret, nil
}

func (r *SecretRepository) FindAllSecrets(ctx context.Context, userName string) ([]model.Secret, error) {
	query := `
		SELECT name, username, content, type, opt_lock
		FROM gophkeeper.secret
		WHERE username = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []model.Secret
	for rows.Next() {
		var secret model.Secret
		if err := rows.Scan(&secret.Name, &secret.Username, &secret.Content, &secret.Type, &secret.Version); err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(secrets) == 0 {
		return nil, pgx.ErrNoRows
	}

	return secrets, nil
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, userName string, secretName string) error {
	query := "DELETE FROM gophkeeper.secret WHERE username = $1 AND name = $2"
	result, err := r.pool.Exec(ctx, query, userName, secretName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
