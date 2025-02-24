package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HubUser struct {
	Name                   string
	Handle                 string
	Email                  string
	ResidentCountry        string
	ResidentCity           string
	PreferredLanguage      string
	ShortBio               string
	LongBio                string
	ProfilePictureFilename string

	ApplyToCompanyDomains []string

	WorkHistoryDomains []string
}

var hubUsers = []HubUser{}

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
	var wg sync.WaitGroup
	green := color.New(color.FgGreen).SprintFunc()
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubUser) {
			fmt.Printf("%s\n", green(fmt.Sprintf("Logging in %s", user.Email)))
			hubLogin(user.Email, "NewPassword123$", &wg)
		}(user)
	}

	// Wait for all hubLogin to complete
	wg.Wait()
}
