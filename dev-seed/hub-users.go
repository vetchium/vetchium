package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

type HubUser struct {
	Name                   string
	Handle                 string
	Email                  string
	Tier                   hub.HubUserTier
	ResidentCountry        string
	ResidentCity           string
	PreferredLanguage      string
	ShortBio               string
	LongBio                string
	ProfilePictureFilename string

	ApplyToCompanyDomains []string
	Endorsers             []common.Handle

	WorkHistoryDomains []string
}

func initHubUsers(db *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	hubUsers := generateHubSeedUsers(50)

	for i, user := range hubUsers {
		userID := fmt.Sprintf("12345678-0000-0000-0000-000000050%03d", i+1)
		_, err = tx.Exec(ctx, `
			INSERT INTO hub_users (
				id, full_name, handle, email, password_hash,
				state, tier, resident_country_code, resident_city,
				preferred_language, short_bio, long_bio
			) VALUES (
				$1, $2, $3, $4,
				'$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
				'ACTIVE_HUB_USER', $5, $6, $7, $8, $9, $10
			)
		`, userID, user.Name, user.Handle, user.Email,
			user.Tier, user.ResidentCountry, user.ResidentCity,
			user.PreferredLanguage, user.ShortBio, user.LongBio)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.Name, err)
		}

		// Print with color directly
		color.New(color.FgGreen).Printf("created hub user %s ", user.Name)
		color.New(color.FgCyan).Printf("<%s> ", user.Email)
		color.New(color.FgYellow).Printf("@%s\n", user.Handle)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func loginHubUsers() {
	var wg sync.WaitGroup
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubUser) {
			color.Green("Logging in %s", user.Email)
			hubLogin(user.Email, "NewPassword123$", &wg)
		}(user)
	}

	// Wait for all hubLogin to complete
	wg.Wait()
}
