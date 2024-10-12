package granger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

type Config struct {
	Port             string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDB       string
	PostgresPassword string
	SMTPHost         string
	SMTPPort         string
	SMTPUser         string
	SMTPPassword     string
	Env              string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:             os.Getenv("PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         os.Getenv("SMTP_PORT"),
		SMTPUser:         os.Getenv("SMTP_USER"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		Env:              os.Getenv("ENV"),
	}

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

	if config.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST environment variable not set")
	}

	if config.SMTPPort == "" {
		return fmt.Errorf("SMTP_PORT environment variable not set")
	}

	if config.SMTPUser == "" {
		return fmt.Errorf("SMTP_USER environment variable not set")
	}

	if config.SMTPPassword == "" {
		return fmt.Errorf("SMTP_PASSWORD environment variable not set")
	}

	_, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		return fmt.Errorf(
			"SMTP_PORT environment variable is not a valid integer: %w",
			err,
		)
	}

	if config.Env != libvetchi.ProdEnv && config.Env != libvetchi.DevEnv &&
		config.Env != libvetchi.TestEnv {
		return fmt.Errorf("ENV environment variable is not valid")
	}

	return nil
}

type Granger struct {
	port         string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	env          string

	db   db.DB
	log  *slog.Logger
	wg   sync.WaitGroup
	quit chan struct{}
}

func NewGranger() (*Granger, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.PostgresHost, config.PostgresPort, config.PostgresUser,
		config.PostgresDB, config.PostgresPassword)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	pg, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	smtpPort, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		// This is unlikely to happen as we already validated the
		// SMTP_PORT environment variable earlier, but we'll check anyway.
		return nil, fmt.Errorf(
			"SMTP_PORT environment variable is not a valid integer: %w",
			err,
		)
	}

	return &Granger{
		port:         fmt.Sprintf(":%s", config.Port),
		SMTPHost:     config.SMTPHost,
		SMTPPort:     smtpPort,
		SMTPUser:     config.SMTPUser,
		SMTPPassword: config.SMTPPassword,
		env:          config.Env,

		db:  pg,
		log: logger,
	}, nil
}

func (g *Granger) Run() error {
	g.wg.Add(1)
	go g.createOnboardEmails()

	g.wg.Add(1)
	mailSenderQuit := make(chan struct{})
	go g.mailSender(mailSenderQuit)

	return http.ListenAndServe(g.port, nil)
}
