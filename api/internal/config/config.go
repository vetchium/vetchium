package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type HermioneConfigOnDisk struct {
	Employer struct {
		TFATokLife     string `json:"tfa_tok_life" validate:"required"`
		SessionTokLife string `json:"session_tok_life" validate:"required"`
		LTSTokLife     string `json:"lts_tok_life" validate:"required"`
		InviteTokLife  string `json:"employee_invite_tok_life" validate:"required"`
	} `json:"employer" validate:"required"`

	Hub struct {
		TFATokLife     string `json:"tfa_tok_life" validate:"required"`
		SessionTokLife string `json:"session_tok_life" validate:"required"`
		LTSTokLife     string `json:"lts_tok_life" validate:"required"`
		InviteTokLife  string `json:"hub_user_invite_tok_life" validate:"required"`
	} `json:"hub" validate:"required"`

	Postgres struct {
		Host string `json:"host" validate:"required,min=1"`
		Port string `json:"port" validate:"required,min=1"`
		User string `json:"user" validate:"required,min=1"`
		DB   string `json:"db" validate:"required,min=1"`
	} `json:"postgres" validate:"required"`

	Port string `json:"port" validate:"required,min=1,number"`
}

type Hermione struct {
	Employer struct {
		TFATokLife     time.Duration
		SessionTokLife time.Duration
		LTSTokLife     time.Duration
		InviteTokLife  time.Duration
	}

	Hub struct {
		TFATokLife     time.Duration
		SessionTokLife time.Duration
		LTSTokLife     time.Duration
		InviteTokLife  time.Duration
	}

	Postgres struct {
		Host string
		Port string
		User string
		DB   string
	}

	Port int
}

func LoadHermioneConfig() (*Hermione, error) {
	data, err := os.ReadFile("/etc/hermione-config/config.json")
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	onDiskConfig := &HermioneConfigOnDisk{}
	if err := json.Unmarshal(data, onDiskConfig); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Whatever can be validated by the struct tags, is done here. More
	// validations continue to happen below
	validate := validator.New()
	if err := validate.Struct(onDiskConfig); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	hc := &Hermione{}

	hc.Postgres.Host = onDiskConfig.Postgres.Host
	hc.Postgres.Port = onDiskConfig.Postgres.Port
	hc.Postgres.User = onDiskConfig.Postgres.User
	hc.Postgres.DB = onDiskConfig.Postgres.DB

	hc.Port, err = strconv.Atoi(onDiskConfig.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	emp := onDiskConfig.Employer
	hc.Employer.TFATokLife, err = time.ParseDuration(emp.TFATokLife)
	if err != nil {
		return nil, fmt.Errorf("employer tfa token life: %w", err)
	}
	hc.Employer.SessionTokLife, err = time.ParseDuration(emp.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer session token life: %w", err)
	}
	hc.Employer.LTSTokLife, err = time.ParseDuration(emp.LTSTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer lts token life: %w", err)
	}
	hc.Employer.InviteTokLife, err = time.ParseDuration(emp.InviteTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer invite token life: %w", err)
	}

	hub := onDiskConfig.Hub
	hc.Hub.TFATokLife, err = time.ParseDuration(hub.TFATokLife)
	if err != nil {
		return nil, fmt.Errorf("hub tfa token life: %w", err)
	}
	hc.Hub.SessionTokLife, err = time.ParseDuration(hub.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub session token life: %w", err)
	}
	hc.Hub.LTSTokLife, err = time.ParseDuration(hub.LTSTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub lts token life: %w", err)
	}
	hc.Hub.InviteTokLife, err = time.ParseDuration(hub.InviteTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub invite token life: %w", err)
	}

	return hc, nil
}
