package main

import (
	"context"
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

func initLocations() {

}

func initCostCenters() {

}
