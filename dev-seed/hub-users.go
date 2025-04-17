package main

import (
	"context"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

var hubUsers []HubSeedUser

func initHubUsers(db *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Should be at least 10 as each user follows 10 other users
	hubUsers = generateHubSeedUsers(200)

	for _, user := range hubUsers {
		query := `
INSERT INTO hub_users (
	full_name, handle, email, password_hash,
	state, tier, resident_country_code, resident_city,
	preferred_language, short_bio, long_bio
) VALUES (
	$1, $2, $3,
	'$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
	'ACTIVE_HUB_USER', $4, $5, $6, $7, $8, $9
)`
		_, err = tx.Exec(
			ctx,
			query,
			user.Name,
			user.Handle,
			user.Email,
			user.Tier,
			user.ResidentCountry,
			user.ResidentCity,
			user.PreferredLanguage,
			user.ShortBio,
			user.LongBio,
		)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.Name, err)
		}

		// Print with color directly
		color.New(color.FgGreen).Printf("created hub user %s ", user.Name)
		color.New(color.FgCyan).Printf("<%s> ", user.Email)
		color.New(color.FgYellow).Printf("@%s\n", user.Handle)
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func loginHubUsers() {
	var wg sync.WaitGroup
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubSeedUser) {
			color.Green("Logging in %s", user.Email)
			hubLogin(user.Email, "NewPassword123$", &wg)
		}(user)
	}

	// Wait for all hubLogin to complete
	wg.Wait()
}
