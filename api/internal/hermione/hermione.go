package hermione

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Port             string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDB       string
	PostgresPassword string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:             os.Getenv("PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
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
	return nil
}

type Hermione struct {
	Port   string
	DB     *pgxpool.Pool
	logger *slog.Logger
}

func New() (*Hermione, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.PostgresHost, config.PostgresPort, config.PostgresUser,
		config.PostgresDB, config.PostgresPassword)
	log.Println("pgConnStr", pgConnStr)

	dbpool, err := pgxpool.New(context.Background(), pgConnStr)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Database connection pool established")

	return &Hermione{
		Port:   fmt.Sprintf(":%s", config.Port),
		DB:     dbpool,
		logger: logger,
	}, nil
}

func (h *Hermione) Run() error {
	return http.ListenAndServe(h.Port, nil)
}
