package repository // package declaration for the module

import ( // start import block
	"context" // import package
	"errors" // import package
	"time" // import package

	"github.com/jackc/pgx/v5" // import package
	"github.com/jackc/pgx/v5/pgxpool" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
) // end import block or block scope

type RefreshTokenRepository interface { // declare custom type
	Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error) // execute statement
	GetValidByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) // execute statement
	RevokeByHash(ctx context.Context, tokenHash string) error // execute statement
} // end block

type PostgresRefreshTokenRepository struct { // declare struct type
	db *pgxpool.Pool // execute statement
} // end block

func NewPostgresRefreshTokenRepository(db *pgxpool.Pool) *PostgresRefreshTokenRepository { // declare function
	return &PostgresRefreshTokenRepository{db: db} // return statement
} // end block

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error) { // declare method
	query := ` // declare and initialize variable
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at) // execute statement
		VALUES ($1, $2, $3) // execute statement
		RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at // execute statement
	` // execute statement

	err := r.db.QueryRow(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).Scan( // declare and initialize variable
		&token.ID, // execute statement
		&token.UserID, // execute statement
		&token.TokenHash, // execute statement
		&token.ExpiresAt, // execute statement
		&token.RevokedAt, // execute statement
		&token.CreatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		return domain.RefreshToken{}, err // return statement
	} // end block

	return token, nil // return statement
} // end block

func (r *PostgresRefreshTokenRepository) GetValidByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at // execute statement
		FROM refresh_tokens // execute statement
		WHERE token_hash = $1 // assign value
		  AND revoked_at IS NULL // execute statement
		  AND expires_at > NOW() // execute statement
	` // execute statement

	var token domain.RefreshToken // execute statement
	err := r.db.QueryRow(ctx, query, tokenHash).Scan( // declare and initialize variable
		&token.ID, // execute statement
		&token.UserID, // execute statement
		&token.TokenHash, // execute statement
		&token.ExpiresAt, // execute statement
		&token.RevokedAt, // execute statement
		&token.CreatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		if errors.Is(err, pgx.ErrNoRows) { // handle no result from database
			return domain.RefreshToken{}, ErrNotFound // handle missing database record
		} // end block
		return domain.RefreshToken{}, err // return statement
	} // end block

	return token, nil // return statement
} // end block

func (r *PostgresRefreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error { // declare method
	query := ` // declare and initialize variable
		UPDATE refresh_tokens // execute statement
		SET revoked_at = $1 // assign value
		WHERE token_hash = $2 AND revoked_at IS NULL // assign value
	` // execute statement

	result, err := r.db.Exec(ctx, query, time.Now(), tokenHash) // declare and initialize variable
	if err != nil { // if condition
		return err // return statement
	} // end block
	if result.RowsAffected() == 0 { // if condition
		return ErrNotFound // handle missing database record
	} // end block
	return nil // return statement
} // end block
