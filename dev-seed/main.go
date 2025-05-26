package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/typespec/hub"
)

// Store session tokens for each user in a thread-safe map
var employerSessionTokens sync.Map
var hubSessionTokens sync.Map

var serverURL = "http://localhost:8080"
var mailPitURL = "http://localhost:8025"

func showTierStatistics() {
	freeTierCount := 0
	paidTierCount := 0

	for _, user := range hubUsers {
		if user.Tier == hub.FreeHubUserTier {
			freeTierCount++
		} else {
			paidTierCount++
		}
	}

	total := len(hubUsers)
	freePercent := float64(freeTierCount) / float64(total) * 100
	paidPercent := float64(paidTierCount) / float64(total) * 100

	color.Green("Hub User Tier Statistics:")
	color.Green("  Free Tier: %d users (%.1f%%)", freeTierCount, freePercent)
	color.Green("  Paid Tier: %d users (%.1f%%)", paidTierCount, paidPercent)
	color.Green("  Total: %d users", total)
}

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

	if os.Getenv("SERVER_URL") != "" {
		serverURL = os.Getenv("SERVER_URL")
	}

	if os.Getenv("MAIL_PIT_URL") != "" {
		mailPitURL = os.Getenv("MAIL_PIT_URL")
	}

	// Directly write to the database
	color.Cyan("Initializing employers and domains")
	initEmployersAndDomains(db)
	color.Cyan("Initializing org users")
	initOrgUsers(db)
	color.Cyan("Initializing hub users")
	initHubUsers(db)
	showTierStatistics()

	color.Cyan("Signing in hub users")
	loginHubUsers()
	color.Cyan("Follow other users")
	followUsers()
	color.Cyan("Follow organizations")
	followOrgs()
	color.Cyan("Write posts")
	writePosts()

	color.Cyan("Create Achievements for Hub Users")
	createAchievements()
	color.Cyan("Creating work histories")
	createWorkHistories()
	color.Cyan("Add Official Emails to Hub Users")
	addOfficialEmails()
	color.Cyan("Uploading profile pictures")
	uploadHubUserProfilePictures()

	// Initialize the PDF directory for resumes
	color.Cyan("Initializing PDF directory for resumes")
	initResumePDFDirectory()

	// Generate PDF resumes for all users
	color.Cyan("Generating PDF resumes for all users")
	generateResumesForAllUsers()

	// Create colleague connections based on overlapping work history
	color.Cyan("Creating colleague connections")
	createColleagueConnections()

	// Use APIs to write to the database
	color.Cyan("Signing in admins")
	signinAdmins()
	color.Cyan("Write employer posts")
	writeEmployerPosts()
	color.Cyan("Initializing locations")
	createLocations()
	color.Cyan("Initializing cost centers")
	createCostCenters()
	color.Cyan("Create Openings")
	createOpenings()

	// Create applications with generated PDF resumes
	color.Cyan("Creating applications")
	createApplications()
}
