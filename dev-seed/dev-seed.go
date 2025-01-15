package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "host=localhost port=5432 user=user dbname=vdb password=pass sslmode=disable"
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Directly write to the database
	initEmployersAndDomains(db)
	initOrgUsers(db)
	initHubUsers(db)

	// Use APIs to write to the database
	initLocations()
	initCostCenters()
}
