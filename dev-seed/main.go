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

	color.Cyan("Signing in hub users")
	loginHubUsers()
	color.Cyan("Creating work histories")
	createWorkHistories()
	color.Cyan("Add Official Emails to Hub Users")
	addOfficialEmails()
	color.Cyan("Uploading profile pictures")
	uploadHubUserProfilePictures()

	// Use APIs to write to the database
	color.Cyan("Signing in admins")
	signinAdmins()
	color.Cyan("Initializing locations")
	createLocations()
	color.Cyan("Initializing cost centers")
	createCostCenters()
	color.Cyan("Create Openings")
	createOpenings()

	// Initialize the PDF directory for resumes
	color.Cyan("Initializing PDF directory for resumes")
	initResumePDFDirectory()

	// Generate PDF resumes for all users
	color.Cyan("Generating PDF resumes for all users")
	generateResumesForAllUsers()

	// Create colleague connections based on overlapping work history
	color.Cyan("Creating colleague connections")
	createColleagueConnections()

	// Create applications with generated PDF resumes
	color.Cyan("Creating applications")
	createApplications()
}
