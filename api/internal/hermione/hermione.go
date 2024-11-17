package hermione

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/psankar/vetchi/api/internal/config"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Hermione struct {
	config *config.Hermione

	// These are initialized programmatically in New()
	hedwig hedwig.Hedwig
	pg     *postgres.PG
	log    util.Logger
	mw     *middleware.Middleware
	vator  *vetchi.Vator
}

func NewHermione() (*Hermione, error) {
	config, err := config.LoadHermioneConfig()
	if err != nil {
		return nil, err
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}

	logger := util.Logger{
		Log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	}

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.DB,
		pgPassword,
	)
	pg, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	vator, err := vetchi.InitValidator(logger)
	if err != nil {
		return nil, err
	}

	// Ensure that the db.DB interface is up to date with the postgres.PG
	// implementation. We somehow need to ensure that no new function is
	// added to the pg without it getting added to the db.DB interface first
	db := db.DB(pg)

	var hermione *Hermione

	hedwig, err := hedwig.NewHedwig(logger)
	if err != nil {
		return nil, fmt.Errorf("Hedwig initialisation failure: %w", err)
	}

	hermione = &Hermione{
		config: config,

		pg:  pg,
		log: logger,

		mw:    middleware.NewMiddleware(db, logger),
		vator: vator,

		hedwig: hedwig,
	}

	return hermione, nil
}

func (h *Hermione) Config() *config.Hermione {
	return h.config
}

func (h *Hermione) DB() *postgres.PG {
	return h.pg
}

func (h *Hermione) Vator() *vetchi.Vator {
	return h.vator
}

func (h *Hermione) Hedwig() hedwig.Hedwig {
	return h.hedwig
}

func (h *Hermione) Err(msg string, args ...any) {
	h.log.Err(msg, args...)
}

func (h *Hermione) Dbg(msg string, args ...any) {
	h.log.Dbg(msg, args...)
}

func (h *Hermione) Inf(msg string, args ...any) {
	h.log.Inf(msg, args...)
}
