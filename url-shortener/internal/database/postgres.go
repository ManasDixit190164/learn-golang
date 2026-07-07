package database // package declaration for the module

import ( // start import block
	"context" // import package
	"time" // import package

	"github.com/jackc/pgx/v5/pgxpool" // import package
) // end import block or block scope

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) { // declare function
	config, err := pgxpool.ParseConfig(databaseURL) // declare and initialize variable
	if err != nil { // if condition
		return nil, err // return statement
	} // end block

	config.MaxConns = 10 // assign value
	config.MinConns = 2 // assign value
	config.MaxConnLifetime = time.Hour // assign value
	config.MaxConnIdleTime = 30 * time.Minute // assign value
	config.HealthCheckPeriod = time.Minute // assign value

	pool, err := pgxpool.NewWithConfig(ctx, config) // declare and initialize variable
	if err != nil { // if condition
		return nil, err // return statement
	} // end block

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // create a context with timeout
	defer cancel() // defer function call

	if err := pool.Ping(pingCtx); err != nil { // if condition
		pool.Close() // execute statement
		return nil, err // return statement
	} // end block

	return pool, nil // return statement
} // end block
