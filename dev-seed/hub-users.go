package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HubUser struct {
	Name                    string
	Handle                  string
	Email                   string
	ResidentCountry         string
	ResidentCity            string
	PreferredLanguage       string
	PreferredCompanyDomains []string
}

var hubUsers = []HubUser{
	// Primarily interested in Gryffindor
	{
		Name:              "Minerva McGonagall",
		Handle:            "minerva",
		Email:             "minerva@hub.example",
		ResidentCountry:   "IND",
		ResidentCity:      "Chennai",
		PreferredLanguage: "ta",
		PreferredCompanyDomains: []string{
			"gryffindor.example",
			"ravenclaw.example",
		},
	},
	{
		Name:              "Neville Longbottom",
		Handle:            "neville",
		Email:             "neville@hub.example",
		ResidentCountry:   "AUS",
		ResidentCity:      "Sydney",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"gryffindor.example",
			"hufflepuff.example",
		},
	},
	{
		Name:              "Rubeus Hagrid",
		Handle:            "hagrid",
		Email:             "hagrid@hub.example",
		ResidentCountry:   "GBR",
		ResidentCity:      "Old Trafford",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"gryffindor.example",
			"slytherin.example",
		},
	},

	// Primarily interested in Hufflepuff
	{
		Name:              "Cedric Diggory",
		Handle:            "cedric",
		Email:             "cedric@hub.example",
		ResidentCountry:   "NZL",
		ResidentCity:      "Wellington",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"hufflepuff.example",
			"gryffindor.example",
		},
	},
	{
		Name:              "Nymphadora Tonks",
		Handle:            "tonks",
		Email:             "tonks@hub.example",
		ResidentCountry:   "USA",
		ResidentCity:      "Provo",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"hufflepuff.example",
			"ravenclaw.example",
		},
	},
	{
		Name:              "Pomona Sprout",
		Handle:            "pomona",
		Email:             "pomona@hub.example",
		ResidentCountry:   "GER",
		ResidentCity:      "Nüremberg",
		PreferredLanguage: "de",
		PreferredCompanyDomains: []string{
			"hufflepuff.example",
			"slytherin.example",
		},
	},

	// Primarily interested in Slytherin
	{
		Name:              "Severus Snape",
		Handle:            "snape",
		Email:             "snape@hub.example",
		ResidentCountry:   "FRA",
		ResidentCity:      "Paris",
		PreferredLanguage: "fr",
		PreferredCompanyDomains: []string{
			"slytherin.example",
			"gryffindor.example",
		},
	},
	{
		Name:              "Draco Malfoy",
		Handle:            "draco",
		Email:             "draco@hub.example",
		ResidentCountry:   "USA",
		ResidentCity:      "Provo",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"slytherin.example",
			"hufflepuff.example",
		},
	},
	{
		Name:              "Tom Riddle",
		Handle:            "tom",
		Email:             "tom@hub.example",
		ResidentCountry:   "USA",
		ResidentCity:      "New York",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"slytherin.example",
			"ravenclaw.example",
		},
	},

	// Primarily interested in Ravenclaw
	{
		Name:              "Luna Lovegood",
		Handle:            "luna",
		Email:             "luna@hub.example",
		ResidentCountry:   "LKA",
		ResidentCity:      "Valvettithurai",
		PreferredLanguage: "ta",
		PreferredCompanyDomains: []string{
			"ravenclaw.example",
			"gryffindor.example",
		},
	},
	{
		Name:              "Cho Chang",
		Handle:            "cho",
		Email:             "cho@hub.example",
		ResidentCountry:   "CHN",
		ResidentCity:      "Shanghai",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"ravenclaw.example",
			"hufflepuff.example",
		},
	},
	{
		Name:              "Xenophilius Lovegood",
		Handle:            "xenophilius",
		Email:             "xenophilius@hub.example",
		ResidentCountry:   "IND",
		ResidentCity:      "Ulsoor",
		PreferredLanguage: "en",
		PreferredCompanyDomains: []string{
			"ravenclaw.example",
			"slytherin.example",
		},
	},
}

func initHubUsers(db *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	for i, user := range hubUsers {
		userID := fmt.Sprintf("12345678-0000-0000-0000-000000050%03d", i+1)
		_, err = tx.Exec(ctx, `
			INSERT INTO hub_users (
				id, full_name, handle, email, password_hash,
				state, resident_country_code, resident_city,
				preferred_language
			) VALUES (
				$1, $2, $3, $4,
				'$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
				'ACTIVE_HUB_USER', $5, $6, $7
			)
		`, userID, user.Name, user.Handle, user.Email,
			user.ResidentCountry, user.ResidentCity,
			user.PreferredLanguage)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.Name, err)
		}
		color.New(color.FgGreen).Printf("created hub user %s ", user.Name)
		color.New(color.FgCyan).Printf("<%s> ", user.Email)
		color.New(color.FgYellow).Printf("@%s\n", user.Handle)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func loginHubUsers() {
	for _, user := range hubUsers {
		go func(user HubUser) {
			color.New(color.FgGreen).Printf("Logging in %s\n", user.Email)
			hubLogin(user.Email, "NewPassword123$")
		}(user)
	}

	// Wait for 10 seconds to allow hubLogin to complete populating
	// the session tokens, as it needs to wait until TFA emails are sent
	<-time.After(10 * time.Second)
}
