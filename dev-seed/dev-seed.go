package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	initDB()
	initEmployersAndDomains()
	initOrgUsers()
	initLocations()
	initCostCenters()
	initHubUsers()
}

func initDB() {
	connStr := "host=localhost port=5432 user=user dbname=vdb password=pass sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db = pool
}

func initEmployersAndDomains() {
	employers := []struct {
		displayName string
		email       string
		domain      string
		shortDomain string
	}{
		{
			displayName: "Gryffindor",
			email:       "admin@gryffindor.example",
			domain:      "gryffindor.example",
			shortDomain: "g.ex",
		},
		{
			displayName: "Hufflepuff",
			email:       "admin@hufflepuff.example",
			domain:      "hufflepuff.example",
			shortDomain: "h.ex",
		},
		{
			displayName: "Ravenclaw",
			email:       "admin@ravenclaw.example",
			domain:      "ravenclaw.example",
			shortDomain: "r.ex",
		},
		{
			displayName: "Slytherin",
			email:       "admin@slytherin.example",
			domain:      "slytherin.example",
			shortDomain: "s.ex",
		},
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	for i, emp := range employers {
		// Create welcome email
		emailID := fmt.Sprintf("12345678-0000-0000-0000-00000000001%d", i+1)
		_, err := tx.Exec(ctx, `
			INSERT INTO emails (
				email_key, email_from, email_to, email_cc, email_bcc,
				email_subject, email_html_body, email_text_body,
				email_state, created_at, processed_at
			) VALUES (
				$1, 'no-reply@vetchi.org', $2, NULL, NULL,
				'Welcome to Vetchi', 'Welcome HTML', 'Welcome Text',
				'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now())
			)
		`, emailID, []string{emp.email})
		if err != nil {
			log.Fatalf(
				"failed to create email for %s: %v",
				emp.displayName,
				err,
			)
		}

		// Create employer
		employerID := fmt.Sprintf("12345678-0000-0000-0000-00000000020%d", i+1)
		_, err = tx.Exec(ctx, `
			INSERT INTO employers (
				id, client_id_type, employer_state, company_name,
				onboard_admin_email, onboard_secret_token, token_valid_till,
				onboard_email_id, created_at
			) VALUES (
				$1, 'DOMAIN', 'ONBOARDED', $2,
				$3, 'blah', timezone('UTC'::text, now()) + interval '1 day',
				$4, timezone('UTC'::text, now())
			)
		`, employerID, emp.displayName, emp.email, emailID)
		if err != nil {
			log.Fatalf("failed to create employer %s: %v", emp.displayName, err)
		}

		// Create primary domain
		primaryDomainID := fmt.Sprintf(
			"12345678-0000-0000-0000-00000000300%d",
			i*2+1,
		)
		_, err = tx.Exec(ctx, `
			INSERT INTO domains (
				id, domain_name, domain_state, employer_id, created_at
			) VALUES (
				$1, $2, 'VERIFIED', $3, timezone('UTC'::text, now())
			)
		`, primaryDomainID, emp.domain, employerID)
		if err != nil {
			log.Fatalf(
				"failed to create primary domain for %s: %v",
				emp.displayName,
				err,
			)
		}

		// Create short domain
		shortDomainID := fmt.Sprintf(
			"12345678-0000-0000-0000-00000000300%d",
			i*2+2,
		)
		_, err = tx.Exec(ctx, `
			INSERT INTO domains (
				id, domain_name, domain_state, employer_id, created_at
			) VALUES (
				$1, $2, 'VERIFIED', $3, timezone('UTC'::text, now())
			)
		`, shortDomainID, emp.shortDomain, employerID)
		if err != nil {
			log.Fatalf(
				"failed to create short domain for %s: %v",
				emp.displayName,
				err,
			)
		}

		// Set primary domain
		_, err = tx.Exec(ctx, `
			INSERT INTO employer_primary_domains (employer_id, domain_id)
			VALUES ($1, $2)
		`, employerID, primaryDomainID)
		if err != nil {
			log.Fatalf(
				"failed to set primary domain for %s: %v",
				emp.displayName,
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func initOrgUsers() {
	users := []struct {
		employerID   string
		employerName string
		adminEmail   string
		orgUsers     []struct {
			name  string
			email string
		}
	}{
		{
			employerID:   "12345678-0000-0000-0000-000000000201",
			employerName: "gryffindor.example",
			adminEmail:   "admin@gryffindor.example",
			orgUsers: []struct {
				name  string
				email string
			}{
				{name: "Harry Potter", email: "harry@gryffindor.example"},
				{
					name:  "Hermione Granger",
					email: "hermione@gryffindor.example",
				},
				{name: "Ron Weasley", email: "ron@gryffindor.example"},
			},
		},
		{
			employerID:   "12345678-0000-0000-0000-000000000202",
			employerName: "hufflepuff.example",
			adminEmail:   "admin@hufflepuff.example",
			orgUsers: []struct {
				name  string
				email string
			}{
				{name: "Cedric Diggory", email: "cedric@hufflepuff.example"},
				{
					name:  "Nymphadora Tonks",
					email: "nymphadora@hufflepuff.example",
				},
				{name: "Newt Scamander", email: "newt@hufflepuff.example"},
			},
		},
		{
			employerID:   "12345678-0000-0000-0000-000000000203",
			employerName: "ravenclaw.example",
			adminEmail:   "admin@ravenclaw.example",
			orgUsers: []struct {
				name  string
				email string
			}{
				{name: "Luna Lovegood", email: "luna@ravenclaw.example"},
				{name: "Cho Chang", email: "cho@ravenclaw.example"},
				{name: "Filius Flitwick", email: "filius@ravenclaw.example"},
			},
		},
		{
			employerID:   "12345678-0000-0000-0000-000000000204",
			employerName: "slytherin.example",
			adminEmail:   "admin@slytherin.example",
			orgUsers: []struct {
				name  string
				email string
			}{
				{name: "Draco Malfoy", email: "draco@slytherin.example"},
				{name: "Severus Snape", email: "severus@slytherin.example"},
				{name: "Horace Slughorn", email: "horace@slytherin.example"},
			},
		},
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	for i, emp := range users {
		// Create admin user first
		adminID := fmt.Sprintf(
			"12345678-0000-0000-0000-000000040%03d",
			(i+1)*100+1,
		)
		_, err = tx.Exec(ctx, `
			INSERT INTO org_users (
				id, email, name, password_hash,
				org_user_roles, org_user_state, employer_id, created_at
			) VALUES (
				$1, $2, 'Admin User', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
				ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', $3, timezone('UTC'::text, now())
			)
		`, adminID, emp.adminEmail, emp.employerID)
		if err != nil {
			log.Fatalf(
				"failed to create admin user for %s: %v",
				emp.employerName,
				err,
			)
		}

		// Create other org users
		for j, user := range emp.orgUsers {
			userID := fmt.Sprintf(
				"12345678-0000-0000-0000-000000040%03d",
				(i+1)*100+j+2,
			)
			_, err = tx.Exec(ctx, `
				INSERT INTO org_users (
					id, email, name, password_hash,
					org_user_roles, org_user_state, employer_id, created_at
				) VALUES (
					$1, $2, $3, '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
					ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', $4, timezone('UTC'::text, now())
				)
			`, userID, user.email, user.name, emp.employerID)
			if err != nil {
				log.Fatalf(
					"failed to create user %s for %s: %v",
					user.name,
					emp.employerName,
					err,
				)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func initLocations() {

}

func initCostCenters() {

}

func initHubUsers() {
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
