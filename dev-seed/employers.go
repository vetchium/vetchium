package main

import (
	"context"
	"fmt"
	"log"
)

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
