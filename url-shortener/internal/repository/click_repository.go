package repository // package declaration for the module

import ( // start import block
	"context" // import package

	"github.com/google/uuid" // import package
	"github.com/jackc/pgx/v5/pgxpool" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
) // end import block or block scope

type ClickRepository interface { // declare custom type
	Create(ctx context.Context, click domain.Click) error // execute statement
	CountByURLID(ctx context.Context, urlID uuid.UUID) (int64, error) // execute statement
	DailyStatsByURLID(ctx context.Context, urlID uuid.UUID, days int) ([]domain.DailyClickStat, error) // execute statement
} // end block

type PostgresClickRepository struct { // declare struct type
	db *pgxpool.Pool // execute statement
} // end block

func NewPostgresClickRepository(db *pgxpool.Pool) *PostgresClickRepository { // declare function
	return &PostgresClickRepository{db: db} // return statement
} // end block

func (r *PostgresClickRepository) Create(ctx context.Context, click domain.Click) error { // declare method
	query := ` // declare and initialize variable
		INSERT INTO clicks (url_id, ip_address, user_agent, referrer) // execute statement
		VALUES ($1, $2, $3, $4) // execute statement
	` // execute statement
	_, err := r.db.Exec(ctx, query, click.URLID, click.IPAddress, click.UserAgent, click.Referrer) // declare and initialize variable
	return err // return statement
} // end block

func (r *PostgresClickRepository) CountByURLID(ctx context.Context, urlID uuid.UUID) (int64, error) { // declare method
	query := `SELECT COUNT(*) FROM clicks WHERE url_id = $1` // declare and initialize variable
	var count int64 // execute statement
	err := r.db.QueryRow(ctx, query, urlID).Scan(&count) // declare and initialize variable
	return count, err // return statement
} // end block

func (r *PostgresClickRepository) DailyStatsByURLID(ctx context.Context, urlID uuid.UUID, days int) ([]domain.DailyClickStat, error) { // declare method
	query := ` // declare and initialize variable
		SELECT TO_CHAR(DATE(clicked_at), 'YYYY-MM-DD') AS date, COUNT(*) AS count // execute statement
		FROM clicks // execute statement
		WHERE url_id = $1 // assign value
		  AND clicked_at >= NOW() - ($2::int * INTERVAL '1 day') // assign value
		GROUP BY DATE(clicked_at) // execute statement
		ORDER BY DATE(clicked_at) // execute statement
	` // execute statement

	rows, err := r.db.Query(ctx, query, urlID, days) // declare and initialize variable
	if err != nil { // if condition
		return nil, err // return statement
	} // end block
	defer rows.Close() // defer function call

	stats := make([]domain.DailyClickStat, 0) // declare and initialize variable
	for rows.Next() { // for loop
		var stat domain.DailyClickStat // execute statement
		if err := rows.Scan(&stat.Date, &stat.Count); err != nil { // if condition
			return nil, err // return statement
		} // end block
		stats = append(stats, stat) // assign value
	} // end block

	return stats, rows.Err() // return statement
} // end block
