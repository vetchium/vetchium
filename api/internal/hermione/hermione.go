package hermione

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Config struct {
	Employer struct {
		TFATokLife string `json:"tfa_tok_life" validate:"required"`

		SessionTokLife string `json:"session_tok_life" validate:"required"`
		LTSTokLife     string `json:"lts_tok_life" validate:"required"`

		InviteTokLife string `json:"employee_invite_tok_life" validate:"required"`
	} `json:"employer" validate:"required"`

	Postgres struct {
		Host string `json:"host" validate:"required,min=1"`
		Port string `json:"port" validate:"required,min=1"`
		User string `json:"user" validate:"required,min=1"`
		DB   string `json:"db" validate:"required,min=1"`
	} `json:"postgres" validate:"required"`

	Port string `json:"port" validate:"required,min=1,number"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("/etc/hermione-config/config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Whatever can be validated by the struct tags, is done here. More
	// validations continue to happen in the New() function
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// We can store the token lifetimes as float64 instead of time.Duration and
// may be able to save some time avoiding the Mins() call everytime we need to
// refer to these values. But a pretty-printed time.Duration is easier to debug
// than a random float64 literal.
type employer struct {
	// TFA Token is sent as response to the signin request and should be used
	// in the subsequent tfa request, to get one of the session tokens.
	tfaTokLife time.Duration

	// One of the below Session Tokens is sent as response to the tfa request
	// and should be used in subsequent /employer/* requests.
	sessionTokLife         time.Duration
	longTermSessionTokLife time.Duration

	// Employee invite token life - Should be used for /employer/add-org-user
	employeeInviteTokLife time.Duration
}

type Hermione struct {
	// These are initialized from configmap
	employer employer
	port     string

	// These are initialized programmatically in New()
	hedwig hedwig.Hedwig
	pg     *postgres.PG
	log    *slog.Logger
	mw     *middleware.Middleware
	vator  *vetchi.Vator
}

func NewHermione() (*Hermione, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

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

	ec := config.Employer
	tfaTokLife, err := time.ParseDuration(ec.TFATokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.TGTLife: %w", err)
	}
	sessionTokLife, err := time.ParseDuration(ec.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.SessionTokLife: %w", err)
	}
	ltsTokLife, err := time.ParseDuration(ec.LTSTokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.LTSessionTokLife: %w", err)
	}
	employeeInviteTokLife, err := time.ParseDuration(ec.InviteTokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.EmployeeInviteTokLife: %w", err)
	}

	// Ensure that the db.DB interface is up to date with the postgres.PG
	// implementation. We somehow need to ensure that no new function is
	// added to the pg without it getting added to the db.DB interface first
	db := db.DB(pg)

	var hermione *Hermione

	hedwig, err := hedwig.NewHedwig(hermione)
	if err != nil {
		return nil, fmt.Errorf("Hedwig initialisation failure: %w", err)
	}

	hermione = &Hermione{
		pg:   pg,
		port: fmt.Sprintf(":%s", config.Port),
		log:  logger,

		mw:    middleware.NewMiddleware(db, logger),
		vator: vator,
		employer: employer{
			tfaTokLife:             tfaTokLife,
			sessionTokLife:         sessionTokLife,
			longTermSessionTokLife: ltsTokLife,
			employeeInviteTokLife:  employeeInviteTokLife,
		},

		hedwig: hedwig,
	}

	return hermione, nil
}

func (h *Hermione) DB() *postgres.PG {
	return h.pg
}

func (h *Hermione) Err(msg string, args ...any) {
	h.log.Error(msg, args...)
}

func (h *Hermione) Dbg(msg string, args ...any) {
	h.log.Debug(msg, args...)
}

func (h *Hermione) Inf(msg string, args ...any) {
	h.log.Info(msg, args...)
}

func (h *Hermione) Vator() *vetchi.Vator {
	return h.vator
}

func (h *Hermione) ConfigDuration(key db.TokenType) (time.Duration, error) {
	switch key {
	case db.EmployerSessionToken:
		return h.employer.sessionTokLife, nil
	case db.EmployerLTSToken:
		return h.employer.longTermSessionTokLife, nil
	case db.EmployerTFAToken:
		return h.employer.tfaTokLife, nil
	case db.EmployerInviteToken:
		return h.employer.employeeInviteTokLife, nil
	default:
		return 0, fmt.Errorf("unknown key: %s", key)
	}
}

func (h *Hermione) Hedwig() hedwig.Hedwig {
	return h.hedwig
}
