package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manasdixit/url-shortener/internal/domain"
)

type URLRepository interface {
	Create(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ShortURL, error)
	FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (domain.ShortURL, error)
	FindByShortCode(ctx context.Context, shortCode string) (domain.ShortURL, error)
	Update(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error)
	Deactivate(ctx context.Context, id, userID uuid.UUID) error
	ShortCodeExists(ctx context.Context, shortCode string) (bool, error)
	CustomAliasExists(ctx context.Context, customAlias string) (bool, error)
}

type PostgresURLRepository struct {
	db *pgxpool.Pool
}

func NewPostgresURLRepository(db *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{db: db}
}

func (r *PostgresURLRepository) Create(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) {
	query := `
		INSERT INTO urls (user_id, original_url, short_code, custom_alias, title, is_active, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		url.UserID,
		url.OriginalURL,
		url.ShortCode,
		url.CustomAlias,
		url.Title,
		url.IsActive,
		url.ExpiresAt,
	).Scan(
		&url.ID,
		&url.UserID,
		&url.OriginalURL,
		&url.ShortCode,
		&url.CustomAlias,
		&url.Title,
		&url.IsActive,
		&url.ExpiresAt,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ShortURL{}, ErrConflict
		}
		return domain.ShortURL{}, err
	}

	return url, nil
}

func (r *PostgresURLRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ShortURL, error) {
	query := `
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at
		FROM urls
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]domain.ShortURL, 0)
	for rows.Next() {
		var url domain.ShortURL
		if err := rows.Scan(
			&url.ID,
			&url.UserID,
			&url.OriginalURL,
			&url.ShortCode,
			&url.CustomAlias,
			&url.Title,
			&url.IsActive,
			&url.ExpiresAt,
			&url.CreatedAt,
			&url.UpdatedAt,
		); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, rows.Err()
}

func (r *PostgresURLRepository) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (domain.ShortURL, error) {
	query := `
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at
		FROM urls
		WHERE id = $1 AND user_id = $2
	`

	return r.scanURL(ctx, query, id, userID)
}

func (r *PostgresURLRepository) FindByShortCode(ctx context.Context, shortCode string) (domain.ShortURL, error) {
	query := `
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at
		FROM urls
		WHERE short_code = $1
	`

	return r.scanURL(ctx, query, shortCode)
}

func (r *PostgresURLRepository) Update(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) {
	query := `
		UPDATE urls
		SET original_url = $1,
		    title = $2,
		    is_active = $3,
		    expires_at = $4,
		    updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		url.OriginalURL,
		url.Title,
		url.IsActive,
		url.ExpiresAt,
		url.ID,
		url.UserID,
	).Scan(
		&url.ID,
		&url.UserID,
		&url.OriginalURL,
		&url.ShortCode,
		&url.CustomAlias,
		&url.Title,
		&url.IsActive,
		&url.ExpiresAt,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ShortURL{}, ErrNotFound
		}
		return domain.ShortURL{}, err
	}

	return url, nil
}

func (r *PostgresURLRepository) Deactivate(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE urls
		SET is_active = FALSE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresURLRepository) ShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, shortCode).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresURLRepository) CustomAliasExists(ctx context.Context, customAlias string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE custom_alias = $1 OR short_code = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, customAlias).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresURLRepository) scanURL(ctx context.Context, query string, args ...any) (domain.ShortURL, error) {
	var url domain.ShortURL
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&url.ID,
		&url.UserID,
		&url.OriginalURL,
		&url.ShortCode,
		&url.CustomAlias,
		&url.Title,
		&url.IsActive,
		&url.ExpiresAt,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ShortURL{}, ErrNotFound
		}
		return domain.ShortURL{}, err
	}
	return url, nil
}
