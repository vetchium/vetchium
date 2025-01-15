package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initHubUsers(db *pgxpool.Pool) {
	hubUsers := []struct {
		name              string
		handle            string
		email             string
		residentCountry   string
		residentCity      string
		preferredLanguage string
	}{
		{
			name:              "Minerva McGonagall",
			handle:            "minerva",
			email:             "minerva@hub.example",
			residentCountry:   "IRL",
			residentCity:      "Dublin",
			preferredLanguage: "en",
		},
		{
			name:              "Pomona Sprout",
			handle:            "pomona",
			email:             "pomona@hub.example",
			residentCountry:   "NZL",
			residentCity:      "Wellington",
			preferredLanguage: "en",
		},
		{
			name:              "Rubeus Hagrid",
			handle:            "hagrid",
			email:             "hagrid@hub.example",
			residentCountry:   "NOR",
			residentCity:      "Bergen",
			preferredLanguage: "en",
		},
		{
			name:              "Sybill Trelawney",
			handle:            "sybill",
			email:             "sybill@hub.example",
			residentCountry:   "ISL",
			residentCity:      "Reykjavik",
			preferredLanguage: "en",
		},
	}

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
		`, userID, user.name, user.handle, user.email,
			user.residentCountry, user.residentCity,
			user.preferredLanguage)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}
