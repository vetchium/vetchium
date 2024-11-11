package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/api/internal/util"
)

type PG struct {
	pool *pgxpool.Pool
	log  util.Logger
}

func New(connStr string, logger util.Logger) (*PG, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	cdb := PG{pool: pool, log: logger}
	return &cdb, nil
}
