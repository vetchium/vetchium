package dolores

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB() *pgxpool.Pool {
	connStr := "host=localhost port=5432 user=user dbname=vdb password=pass sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return pool
}
