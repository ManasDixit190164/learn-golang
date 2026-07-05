package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manasdixit/url-shortener/internal/domain"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error)
	GetValidByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
}

type PostgresRefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRefreshTokenRepository(db *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at
	`

	err := r.db.QueryRow(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return token, nil
}

func (r *PostgresRefreshTokenRepository) GetValidByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`

	var token domain.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, ErrNotFound
		}
		return domain.RefreshToken{}, err
	}

	return token, nil
}

func (r *PostgresRefreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2 AND revoked_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, time.Now(), tokenHash)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
