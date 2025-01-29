package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store session tokens for each user in a thread-safe map
var employerSessionTokens sync.Map
var hubSessionTokens sync.Map

func main() {
	log.SetFlags(log.Lshortfile)

	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Fatal("POSTGRES_URI environment variable is required")
	}
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Directly write to the database

	color.Cyan("Initializing employers and domains")
	initEmployersAndDomains(db)
	color.Cyan("Initializing org users")
	initOrgUsers(db)
	color.Cyan("Initializing hub users")
	initHubUsers(db)

	// Use APIs to write to the database
	color.Cyan("Signing in admins")
	signinAdmins()
	color.Cyan("Initializing locations")
	createLocations()
	color.Cyan("Initializing cost centers")
	createCostCenters()
	color.Cyan("Create Openings")
	createOpenings()

	color.Cyan("Signing in hub users")
	signinHubUsers()
}
