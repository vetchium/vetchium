package hermione

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Config struct {
	Port                     string
	PostgresHost             string
	PostgresPort             string
	PostgresUser             string
	PostgresDB               string
	PostgresPassword         string
	SessionTokenValidMins    string
	LongTermSessionValidMins string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:                     os.Getenv("PORT"),
		PostgresHost:             os.Getenv("POSTGRES_HOST"),
		PostgresPort:             os.Getenv("POSTGRES_PORT"),
		PostgresUser:             os.Getenv("POSTGRES_USER"),
		PostgresDB:               os.Getenv("POSTGRES_DB"),
		PostgresPassword:         os.Getenv("POSTGRES_PASSWORD"),
		SessionTokenValidMins:    os.Getenv("SESSION_TOKEN_VALID_MINS"),
		LongTermSessionValidMins: os.Getenv("LONG_TERM_SESSION_VALID_MINS"),
	}

	// Validate required fields
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.Port == "" {
		return fmt.Errorf("PORT environment variable not set")
	}
	if config.PostgresHost == "" {
		return fmt.Errorf("POSTGRES_HOST environment variable not set")
	}
	if config.PostgresPort == "" {
		return fmt.Errorf("POSTGRES_PORT environment variable not set")
	}
	if config.PostgresUser == "" {
		return fmt.Errorf("POSTGRES_USER environment variable not set")
	}
	if config.PostgresDB == "" {
		return fmt.Errorf("POSTGRES_DB environment variable not set")
	}
	if config.PostgresPassword == "" {
		return fmt.Errorf("POSTGRES_PASSWORD environment variable not set")
	}
	if config.SessionTokenValidMins == "" {
		_, err := time.ParseDuration(config.SessionTokenValidMins)
		if err != nil {
			return fmt.Errorf("SESSION_TOKEN_VALID_MINS invalid duration")
		}
	}
	if config.LongTermSessionValidMins == "" {
		_, err := time.ParseDuration(config.LongTermSessionValidMins)
		if err != nil {
			return fmt.Errorf("LONG_TERM_SESSION_VALID_MINS invalid duration")
		}
	}

	return nil
}

type Hermione struct {
	// This is initialized from config
	port string

	// These are initialized from config
	longTermSessionValidMins float64
	sessionTokenValidMins    float64

	// These are initialized programmatically in New()
	db    db.DB
	log   *slog.Logger
	vator *vetchi.Vator
}

func New() (*Hermione, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresDB,
		config.PostgresPassword,
	)
	db, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	vator, err := vetchi.InitValidator(logger)
	if err != nil {
		return nil, err
	}

	sessionTokenValidity, err := time.ParseDuration(
		config.SessionTokenValidMins,
	)
	if err != nil {
		// This is unlikely to happen because we have already validated
		// the config value in LoadConfig()
		return nil, fmt.Errorf("SESSION_TOKEN_VALID_MINS invalid duration")
	}
	longTermSessionValidity, err := time.ParseDuration(
		config.LongTermSessionValidMins,
	)
	if err != nil {
		// This is unlikely to happen because we have already validated
		// the config value in LoadConfig()
		return nil, fmt.Errorf("LONG_TERM_SESSION_VALID_MINS invalid duration")
	}

	return &Hermione{
		db:                       db,
		port:                     fmt.Sprintf(":%s", config.Port),
		log:                      logger,
		vator:                    vator,
		sessionTokenValidMins:    sessionTokenValidity.Minutes(),
		longTermSessionValidMins: longTermSessionValidity.Minutes(),
	}, nil
}

func (h *Hermione) Run() error {
	http.HandleFunc("/employer/get-onboard-status", h.getOnboardStatus)
	http.HandleFunc("/employer/set-onboard-password", h.setOnboardPassword)
	http.HandleFunc("/employer/signin", h.employerSignin)

	return http.ListenAndServe(h.port, nil)
}
