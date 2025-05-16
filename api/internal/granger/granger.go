package granger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	ristretto "github.com/dgraph-io/ristretto/v2"
	"github.com/go-playground/validator/v10"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/postgres"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
)

type Config struct {
	Env              string `json:"env"                validate:"required,min=1"`
	OnboardTokenLife string `json:"onboard_token_life" validate:"required,min=1"`
	Port             string `json:"port"               validate:"required,min=1,number"`
	EmployerBaseURL  string `json:"employer_base_url"  validate:"required"`
	HubBaseURL       string `json:"hub_base_url"       validate:"required"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("/etc/granger-config/config.json")
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

	if config.Env != vetchi.ProdEnv && config.Env != vetchi.DevEnv {
		return nil, fmt.Errorf(
			"%q is not one of [%q, %q]",
			config.Env,
			vetchi.ProdEnv,
			vetchi.DevEnv,
		)
	}

	return config, nil
}

type smtpCredentials struct {
	host     string
	port     int
	user     string
	password string
}

type Granger struct {
	// These are initialized from configmap
	env              string
	onboardTokenLife time.Duration
	port             string
	smtp             smtpCredentials

	employerBaseURL string
	hubBaseURL      string

	// These are initialized programatically in NewGranger()
	db  db.DB
	log util.Logger
	wg  sync.WaitGroup

	employerActiveJobCountCache *ristretto.Cache
	employerEmployeeCountCache  *ristretto.Cache
}

func NewGranger() (*Granger, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgConnStr := os.Getenv("POSTGRES_URI")
	if pgConnStr == "" {
		return nil, fmt.Errorf("POSTGRES_URI not set")
	}

	logger := util.Logger{
		Log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	}

	var sc smtpCredentials
	sc.host = os.Getenv("SMTP_HOST")
	if sc.host == "" {
		return nil, fmt.Errorf("SMTP_HOST not set")
	}

	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		sc.port = 587
	} else {
		sc.port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("SMTP_PORT is invalid: %w", err)
		}
	}

	sc.user = os.Getenv("SMTP_USER")
	if sc.user == "" {
		return nil, fmt.Errorf("SMTP_USER not set")
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD not set")
	}

	db, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	tokenDuration, err := time.ParseDuration(config.OnboardTokenLife)
	if err != nil {
		return nil, fmt.Errorf("OnboardTokenLife is invalid: %w", err)
	}

	g := &Granger{
		env:              config.Env,
		port:             fmt.Sprintf(":%s", config.Port),
		smtp:             sc,
		onboardTokenLife: tokenDuration,

		employerBaseURL: config.EmployerBaseURL,
		hubBaseURL:      config.HubBaseURL,

		db:  db,
		log: logger,

		employerActiveJobCountCache: ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e6,
			MaxCost:     1 << 20,
			BufferItems: 64,
		}),
		employerEmployeeCountCache: ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e6,
			MaxCost:     1 << 20,
			BufferItems: 64,
		}),
	}

	return g, nil
}

func (g *Granger) Run() error {
	g.wg.Add(1)
	pruneTokensQuit := make(chan struct{})
	go g.pruneTokens(pruneTokensQuit)

	// TODO: This should be moved to Hedwig templates
	g.wg.Add(1)
	createOnboardEmailsQuit := make(chan struct{})
	go g.createOnboardEmails(createOnboardEmailsQuit)

	g.wg.Add(1)
	pruneOfficialEmailCodesQuit := make(chan struct{})
	go g.pruneOfficialEmailCodes(pruneOfficialEmailCodesQuit)

	g.wg.Add(1)
	mailSenderQuit := make(chan struct{})
	go g.mailSender(mailSenderQuit)

	g.wg.Add(1)
	scoreApplicationsQuit := make(chan struct{})
	go g.scoreApplications(scoreApplicationsQuit)

	g.wg.Add(1)
	timelineRefresherQuit := make(chan struct{})
	go g.TimelineRefresher(timelineRefresherQuit)

	go func() {
		http.HandleFunc(
			"/internal/get-employer-active-job-count",
			g.getEmployerActiveJobCount,
		)

		http.HandleFunc(
			"/internal/get-employer-employee-count",
			g.getEmployerEmployeeCount,
		)

		err := http.ListenAndServe(g.port, nil)
		if err != nil {
			g.log.Err("Failed to start HTTP server", "error", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	go func() {
		<-signalChan
		g.log.Dbg("Received TERM signal, closing all channels")
		close(pruneTokensQuit)
		close(createOnboardEmailsQuit)
		close(pruneOfficialEmailCodesQuit)
		close(mailSenderQuit)
		close(scoreApplicationsQuit)
	}()

	g.wg.Wait()
	return nil
}
