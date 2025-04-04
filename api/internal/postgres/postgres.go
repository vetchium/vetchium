package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/api/internal/util"
)

type PG struct {
	pool *pgxpool.Pool
	log  util.Logger
}

// registerCustomTypes registers custom type mappings with pgx
func registerCustomTypes(conn *pgx.Conn) error {
	// Load the endorsement_states enum type from the database
	endorsementType, err := conn.LoadType(
		context.Background(),
		"endorsement_states",
	)
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterType(endorsementType)

	// Load the array type for endorsement_states
	endorsementArrayType, err := conn.LoadType(
		context.Background(),
		"_endorsement_states",
	)
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterType(endorsementArrayType)

	return nil
}

func New(connStr string, logger util.Logger) (*PG, error) {
	// Create a configuration
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Set up the AfterConnect hook to register custom types
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logger.Dbg("registering custom type mappings for EndorsementState")
		err := registerCustomTypes(conn)
		if err != nil {
			logger.Err("failed to register custom types", "error", err)
		}
		return nil // Don't fail connection on type registration error
	}

	// Create the connection pool with the modified config
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	cdb := PG{pool: pool, log: logger}
	return &cdb, nil
}
