package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

// signinAdmins logs in all admin users and stores their tokens
func signinAdmins() {
	admins := []struct {
		email    string
		clientID string
	}{
		{
			email:    "admin@gryffindor.example",
			clientID: "gryffindor.example",
		},
		{
			email:    "admin@hufflepuff.example",
			clientID: "hufflepuff.example",
		},
		{
			email:    "admin@ravenclaw.example",
			clientID: "ravenclaw.example",
		},
		{
			email:    "admin@slytherin.example",
			clientID: "slytherin.example",
		},
		{
			email:    "admin@nithra.example",
			clientID: "nithra.example",
		},
	}

	for _, admin := range admins {
		employerSessionTokens.Store(
			admin.email,
			employerSignin(admin.email, "NewPassword123$", admin.clientID),
		)
	}
}

func initOrgUsers(db *pgxpool.Pool) {
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
				{
					name:  "Remus Lupin",
					email: "remus@gryffindor.example",
				},
				{
					name:  "Minerva McGonagall",
					email: "minerva@gryffindor.example",
				},
				{
					name:  "Rubeus Hagrid",
					email: "hagrid@gryffindor.example",
				},
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
				{
					name:  "Pomona Sprout",
					email: "pomona@hufflepuff.example",
				},
				{
					name:  "Silvanus Kettleburn",
					email: "silvanus@hufflepuff.example",
				},
				{
					name:  "Newt Scamander",
					email: "newt@hufflepuff.example",
				},
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
				{
					name:  "Filius Flitwick",
					email: "filius@ravenclaw.example",
				},
				{
					name:  "Rowena Ravenclaw",
					email: "rowena@ravenclaw.example",
				},
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
				{name: "Salazar Slytherin", email: "salazar@slytherin.example"},
				{name: "Severus Snape", email: "severus@slytherin.example"},
				{name: "Horace Slughorn", email: "horace@slytherin.example"},
			},
		},
		{
			employerID:   "12345678-0000-0000-0000-000000000205",
			employerName: "nithra.example",
			adminEmail:   "admin@nithra.example",
			orgUsers: []struct {
				name  string
				email string
			}{
				{name: "HR", email: "hr@nithra.example"},
				{name: "Doctor", email: "doctor@nithra.example"},
			},
		},
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

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
		color.New(color.FgGreen).
			Printf("created admin user %s ", emp.adminEmail)
		color.New(color.FgCyan).Printf("<%s> ", emp.adminEmail)
		color.New(color.FgYellow).Printf("for %s\n", emp.employerName)

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
			color.New(color.FgGreen).Printf("created user %s ", user.name)
			color.New(color.FgCyan).Printf("<%s> ", user.email)
			color.New(color.FgYellow).Printf("for %s\n", emp.employerName)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}
