package repository // package declaration for the module

import ( // start import block
	"context" // import package
	"errors" // import package

	"github.com/google/uuid" // import package
	"github.com/jackc/pgx/v5" // import package
	"github.com/jackc/pgx/v5/pgconn" // import package
	"github.com/jackc/pgx/v5/pgxpool" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
) // end import block or block scope

type URLRepository interface { // declare custom type
	Create(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) // execute statement
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ShortURL, error) // execute statement
	FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (domain.ShortURL, error) // execute statement
	FindByShortCode(ctx context.Context, shortCode string) (domain.ShortURL, error) // execute statement
	Update(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) // execute statement
	Deactivate(ctx context.Context, id, userID uuid.UUID) error // execute statement
	ShortCodeExists(ctx context.Context, shortCode string) (bool, error) // execute statement
	CustomAliasExists(ctx context.Context, customAlias string) (bool, error) // execute statement
} // end block

type PostgresURLRepository struct { // declare struct type
	db *pgxpool.Pool // execute statement
} // end block

func NewPostgresURLRepository(db *pgxpool.Pool) *PostgresURLRepository { // declare function
	return &PostgresURLRepository{db: db} // return statement
} // end block

func (r *PostgresURLRepository) Create(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) { // declare method
	query := ` // declare and initialize variable
		INSERT INTO urls (user_id, original_url, short_code, custom_alias, title, is_active, expires_at) // execute statement
		VALUES ($1, $2, $3, $4, $5, $6, $7) // execute statement
		RETURNING id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at // execute statement
	` // execute statement

	err := r.db.QueryRow(ctx, query, // declare and initialize variable
		url.UserID, // execute statement
		url.OriginalURL, // execute statement
		url.ShortCode, // execute statement
		url.CustomAlias, // execute statement
		url.Title, // execute statement
		url.IsActive, // execute statement
		url.ExpiresAt, // execute statement
	).Scan( // execute statement
		&url.ID, // execute statement
		&url.UserID, // execute statement
		&url.OriginalURL, // execute statement
		&url.ShortCode, // execute statement
		&url.CustomAlias, // execute statement
		&url.Title, // execute statement
		&url.IsActive, // execute statement
		&url.ExpiresAt, // execute statement
		&url.CreatedAt, // execute statement
		&url.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		var pgErr *pgconn.PgError // execute statement
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // if condition
			return domain.ShortURL{}, ErrConflict // handle duplicate database entry
		} // end block
		return domain.ShortURL{}, err // return statement
	} // end block

	return url, nil // return statement
} // end block

func (r *PostgresURLRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ShortURL, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at // execute statement
		FROM urls // execute statement
		WHERE user_id = $1 // assign value
		ORDER BY created_at DESC // execute statement
	` // execute statement

	rows, err := r.db.Query(ctx, query, userID) // declare and initialize variable
	if err != nil { // if condition
		return nil, err // return statement
	} // end block
	defer rows.Close() // defer function call

	urls := make([]domain.ShortURL, 0) // declare and initialize variable
	for rows.Next() { // for loop
		var url domain.ShortURL // execute statement
		if err := rows.Scan( // if condition
			&url.ID, // execute statement
			&url.UserID, // execute statement
			&url.OriginalURL, // execute statement
			&url.ShortCode, // execute statement
			&url.CustomAlias, // execute statement
			&url.Title, // execute statement
			&url.IsActive, // execute statement
			&url.ExpiresAt, // execute statement
			&url.CreatedAt, // execute statement
			&url.UpdatedAt, // execute statement
		); err != nil { // assign value
			return nil, err // return statement
		} // end block
		urls = append(urls, url) // assign value
	} // end block

	return urls, rows.Err() // return statement
} // end block

func (r *PostgresURLRepository) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (domain.ShortURL, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at // execute statement
		FROM urls // execute statement
		WHERE id = $1 AND user_id = $2 // assign value
	` // execute statement

	return r.scanURL(ctx, query, id, userID) // return statement
} // end block

func (r *PostgresURLRepository) FindByShortCode(ctx context.Context, shortCode string) (domain.ShortURL, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at // execute statement
		FROM urls // execute statement
		WHERE short_code = $1 // assign value
	` // execute statement

	return r.scanURL(ctx, query, shortCode) // return statement
} // end block

func (r *PostgresURLRepository) Update(ctx context.Context, url domain.ShortURL) (domain.ShortURL, error) { // declare method
	query := ` // declare and initialize variable
		UPDATE urls // execute statement
		SET original_url = $1, // assign value
		    title = $2, // assign value
		    is_active = $3, // assign value
		    expires_at = $4, // assign value
		    updated_at = NOW() // assign value
		WHERE id = $5 AND user_id = $6 // assign value
		RETURNING id, user_id, original_url, short_code, custom_alias, title, is_active, expires_at, created_at, updated_at // execute statement
	` // execute statement

	err := r.db.QueryRow(ctx, query, // declare and initialize variable
		url.OriginalURL, // execute statement
		url.Title, // execute statement
		url.IsActive, // execute statement
		url.ExpiresAt, // execute statement
		url.ID, // execute statement
		url.UserID, // execute statement
	).Scan( // execute statement
		&url.ID, // execute statement
		&url.UserID, // execute statement
		&url.OriginalURL, // execute statement
		&url.ShortCode, // execute statement
		&url.CustomAlias, // execute statement
		&url.Title, // execute statement
		&url.IsActive, // execute statement
		&url.ExpiresAt, // execute statement
		&url.CreatedAt, // execute statement
		&url.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		if errors.Is(err, pgx.ErrNoRows) { // handle no result from database
			return domain.ShortURL{}, ErrNotFound // handle missing database record
		} // end block
		return domain.ShortURL{}, err // return statement
	} // end block

	return url, nil // return statement
} // end block

func (r *PostgresURLRepository) Deactivate(ctx context.Context, id, userID uuid.UUID) error { // declare method
	query := ` // declare and initialize variable
		UPDATE urls // execute statement
		SET is_active = FALSE, updated_at = NOW() // assign value
		WHERE id = $1 AND user_id = $2 // assign value
	` // execute statement

	result, err := r.db.Exec(ctx, query, id, userID) // declare and initialize variable
	if err != nil { // if condition
		return err // return statement
	} // end block
	if result.RowsAffected() == 0 { // if condition
		return ErrNotFound // handle missing database record
	} // end block
	return nil // return statement
} // end block

func (r *PostgresURLRepository) ShortCodeExists(ctx context.Context, shortCode string) (bool, error) { // declare method
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)` // declare and initialize variable
	var exists bool // execute statement
	if err := r.db.QueryRow(ctx, query, shortCode).Scan(&exists); err != nil { // if condition
		return false, err // return statement
	} // end block
	return exists, nil // return statement
} // end block

func (r *PostgresURLRepository) CustomAliasExists(ctx context.Context, customAlias string) (bool, error) { // declare method
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE custom_alias = $1 OR short_code = $1)` // declare and initialize variable
	var exists bool // execute statement
	if err := r.db.QueryRow(ctx, query, customAlias).Scan(&exists); err != nil { // if condition
		return false, err // return statement
	} // end block
	return exists, nil // return statement
} // end block

func (r *PostgresURLRepository) scanURL(ctx context.Context, query string, args ...any) (domain.ShortURL, error) { // declare method
	var url domain.ShortURL // execute statement
	err := r.db.QueryRow(ctx, query, args...).Scan( // declare and initialize variable
		&url.ID, // execute statement
		&url.UserID, // execute statement
		&url.OriginalURL, // execute statement
		&url.ShortCode, // execute statement
		&url.CustomAlias, // execute statement
		&url.Title, // execute statement
		&url.IsActive, // execute statement
		&url.ExpiresAt, // execute statement
		&url.CreatedAt, // execute statement
		&url.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		if errors.Is(err, pgx.ErrNoRows) { // handle no result from database
			return domain.ShortURL{}, ErrNotFound // handle missing database record
		} // end block
		return domain.ShortURL{}, err // return statement
	} // end block
	return url, nil // return statement
} // end block
