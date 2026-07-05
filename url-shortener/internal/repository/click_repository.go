package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manasdixit/url-shortener/internal/domain"
)

type ClickRepository interface {
	Create(ctx context.Context, click domain.Click) error
	CountByURLID(ctx context.Context, urlID uuid.UUID) (int64, error)
	DailyStatsByURLID(ctx context.Context, urlID uuid.UUID, days int) ([]domain.DailyClickStat, error)
}

type PostgresClickRepository struct {
	db *pgxpool.Pool
}

func NewPostgresClickRepository(db *pgxpool.Pool) *PostgresClickRepository {
	return &PostgresClickRepository{db: db}
}

func (r *PostgresClickRepository) Create(ctx context.Context, click domain.Click) error {
	query := `
		INSERT INTO clicks (url_id, ip_address, user_agent, referrer)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, click.URLID, click.IPAddress, click.UserAgent, click.Referrer)
	return err
}

func (r *PostgresClickRepository) CountByURLID(ctx context.Context, urlID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM clicks WHERE url_id = $1`
	var count int64
	err := r.db.QueryRow(ctx, query, urlID).Scan(&count)
	return count, err
}

func (r *PostgresClickRepository) DailyStatsByURLID(ctx context.Context, urlID uuid.UUID, days int) ([]domain.DailyClickStat, error) {
	query := `
		SELECT TO_CHAR(DATE(clicked_at), 'YYYY-MM-DD') AS date, COUNT(*) AS count
		FROM clicks
		WHERE url_id = $1
		  AND clicked_at >= NOW() - ($2::int * INTERVAL '1 day')
		GROUP BY DATE(clicked_at)
		ORDER BY DATE(clicked_at)
	`

	rows, err := r.db.Query(ctx, query, urlID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]domain.DailyClickStat, 0)
	for rows.Next() {
		var stat domain.DailyClickStat
		if err := rows.Scan(&stat.Date, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}
