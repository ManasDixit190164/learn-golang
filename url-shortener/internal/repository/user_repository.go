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

type UserRepository interface { // declare custom type
	Create(ctx context.Context, user domain.User) (domain.User, error) // execute statement
	GetByEmail(ctx context.Context, email string) (domain.User, error) // execute statement
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) // execute statement
} // end block

type PostgresUserRepository struct { // declare struct type
	db *pgxpool.Pool // execute statement
} // end block

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository { // declare function
	return &PostgresUserRepository{db: db} // return statement
} // end block

func (r *PostgresUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) { // declare method
	query := ` // declare and initialize variable
		INSERT INTO users (name, email, password_hash) // execute statement
		VALUES ($1, $2, $3) // execute statement
		RETURNING id, name, email, password_hash, is_email_verified, created_at, updated_at // execute statement
	` // execute statement

	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash).Scan( // declare and initialize variable
		&user.ID, // execute statement
		&user.Name, // execute statement
		&user.Email, // execute statement
		&user.PasswordHash, // execute statement
		&user.IsEmailVerified, // execute statement
		&user.CreatedAt, // execute statement
		&user.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		var pgErr *pgconn.PgError // execute statement
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // if condition
			return domain.User{}, ErrConflict // handle duplicate database entry
		} // end block
		return domain.User{}, err // return statement
	} // end block

	return user, nil // return statement
} // end block

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, name, email, password_hash, is_email_verified, created_at, updated_at // execute statement
		FROM users // execute statement
		WHERE email = $1 // assign value
	` // execute statement

	var user domain.User // execute statement
	err := r.db.QueryRow(ctx, query, email).Scan( // declare and initialize variable
		&user.ID, // execute statement
		&user.Name, // execute statement
		&user.Email, // execute statement
		&user.PasswordHash, // execute statement
		&user.IsEmailVerified, // execute statement
		&user.CreatedAt, // execute statement
		&user.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		if errors.Is(err, pgx.ErrNoRows) { // handle no result from database
			return domain.User{}, ErrNotFound // handle missing database record
		} // end block
		return domain.User{}, err // return statement
	} // end block

	return user, nil // return statement
} // end block

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) { // declare method
	query := ` // declare and initialize variable
		SELECT id, name, email, password_hash, is_email_verified, created_at, updated_at // execute statement
		FROM users // execute statement
		WHERE id = $1 // assign value
	` // execute statement

	var user domain.User // execute statement
	err := r.db.QueryRow(ctx, query, id).Scan( // declare and initialize variable
		&user.ID, // execute statement
		&user.Name, // execute statement
		&user.Email, // execute statement
		&user.PasswordHash, // execute statement
		&user.IsEmailVerified, // execute statement
		&user.CreatedAt, // execute statement
		&user.UpdatedAt, // execute statement
	) // execute statement
	if err != nil { // if condition
		if errors.Is(err, pgx.ErrNoRows) { // handle no result from database
			return domain.User{}, ErrNotFound // handle missing database record
		} // end block
		return domain.User{}, err // return statement
	} // end block

	return user, nil // return statement
} // end block
