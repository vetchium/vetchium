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
	ShortBio                string
	LongBio                 string
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
		ShortBio: "Minerva McGonagall is wise",
		LongBio:  "Minerva McGonagall was born in Scotland and finished education at Hogwarts and has 40 years as experience.",
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
		ShortBio: "Neville Longbottom is brave",
		LongBio:  "Neville Longbottom was born in England and finished education at Hogwarts and has 10 years as experience.",
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
		ShortBio: "Rubeus Hagrid is caring",
		LongBio:  "Rubeus Hagrid was born in England and finished education at Hogwarts and has 50 years as experience.",
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
		ShortBio: "Cedric Diggory is fair",
		LongBio:  "Cedric Diggory was born in England and finished education at Hogwarts and has 7 years as experience.",
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
		ShortBio: "Nymphadora Tonks is adaptable",
		LongBio:  "Nymphadora Tonks was born in England and finished education at Hogwarts and has 8 years as experience.",
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
		ShortBio: "Pomona Sprout is nurturing",
		LongBio:  "Pomona Sprout was born in Wales and finished education at Hogwarts and has 35 years as experience.",
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
		ShortBio: "Severus Snape is precise",
		LongBio:  "Severus Snape was born in England and finished education at Hogwarts and has 20 years as experience.",
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
		ShortBio: "Draco Malfoy is ambitious",
		LongBio:  "Draco Malfoy was born in England and finished education at Hogwarts and has 5 years as experience.",
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
		ShortBio: "Tom Riddle is determined",
		LongBio:  "Tom Riddle was born in England and finished education at Hogwarts and has 50 years as experience.",
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
		ShortBio: "Luna Lovegood is creative",
		LongBio:  "Luna Lovegood was born in England and finished education at Hogwarts and has 3 years as experience.",
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
		ShortBio: "Cho Chang is intelligent",
		LongBio:  "Cho Chang was born in Scotland and finished education at Hogwarts and has 4 years as experience.",
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
		ShortBio: "Xenophilius Lovegood is innovative",
		LongBio:  "Xenophilius Lovegood was born in England and finished education at Hogwarts and has 25 years as experience.",
	},
}

func initHubUsers(db *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for i, user := range hubUsers {
		userID := fmt.Sprintf("12345678-0000-0000-0000-000000050%03d", i+1)
		_, err = tx.Exec(ctx, `
			INSERT INTO hub_users (
				id, full_name, handle, email, password_hash,
				state, resident_country_code, resident_city,
				preferred_language, short_bio, long_bio
			) VALUES (
				$1, $2, $3, $4,
				'$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
				'ACTIVE_HUB_USER', $5, $6, $7, $8, $9
			)
		`, userID, user.Name, user.Handle, user.Email,
			user.ResidentCountry, user.ResidentCity,
			user.PreferredLanguage, user.ShortBio, user.LongBio)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.Name, err)
		}
		fmt.Printf("%s %s %s\n",
			green(fmt.Sprintf("created hub user %s", user.Name)),
			cyan(fmt.Sprintf("<%s>", user.Email)),
			yellow(fmt.Sprintf("@%s", user.Handle)),
		)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func loginHubUsers() {
	green := color.New(color.FgGreen).SprintFunc()
	for _, user := range hubUsers {
		go func(user HubUser) {
			fmt.Printf("%s\n", green(fmt.Sprintf("Logging in %s", user.Email)))
			hubLogin(user.Email, "NewPassword123$")
		}(user)
	}

	// Wait for 10 seconds to allow hubLogin to complete populating
	// the session tokens, as it needs to wait until TFA emails are sent
	<-time.After(10 * time.Second)
}
