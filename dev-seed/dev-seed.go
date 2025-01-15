package main

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store session tokens for each user in a thread-safe map
var sessionTokens sync.Map

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	connStr := "host=localhost port=5432 user=user dbname=vdb password=pass sslmode=disable"
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Directly write to the database
	log.Println("Initializing employers and domains")
	initEmployersAndDomains(db)
	log.Println("Initializing org users")
	initOrgUsers(db)
	log.Println("Initializing hub users")
	initHubUsers(db)

	// Use APIs to write to the database
	log.Println("Signing in admins")
	signinAdmins()
	log.Println("Initializing locations")
	initLocations()
	log.Println("Initializing cost centers")
	initCostCenters()
}
