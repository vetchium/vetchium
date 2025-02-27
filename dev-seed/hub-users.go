package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/typespec/common"
)

type HubUser struct {
	Name                   string
	Handle                 string
	Email                  string
	ResidentCountry        string
	ResidentCity           string
	PreferredLanguage      string
	ShortBio               string
	LongBio                string
	ProfilePictureFilename string

	ApplyToCompanyDomains []string
	Endorsers             []common.Handle

	WorkHistoryDomains []string
}

var hubUsers = []HubUser{
	{
		Name:              "User One",
		Handle:            "user1",
		Email:             "user1@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "New York",
		PreferredLanguage: "en",
		ShortBio:          "Software Engineer",
		LongBio:           "Experienced software engineer with focus on backend development",
		WorkHistoryDomains: []string{
			"sunvaja.example",
			"decdpd.example",
			"nokiabricks.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
	},
	{
		Name:               "User Two",
		Handle:             "user2",
		Email:              "user2@example.com",
		ResidentCountry:    "GBR",
		ResidentCity:       "London",
		PreferredLanguage:  "en",
		ShortBio:           "Product Manager",
		LongBio:            "Product manager with experience in tech industry",
		WorkHistoryDomains: []string{"novelltenware.example", "decdpd.example"},
	},
	{
		Name:              "User Three",
		Handle:            "user3",
		Email:             "user3@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Berlin",
		PreferredLanguage: "de",
		ShortBio:          "Frontend Developer",
		LongBio:           "Frontend developer specializing in React and TypeScript",
		WorkHistoryDomains: []string{
			"decdpd.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:               "User Four",
		Handle:             "user4",
		Email:              "user4@example.com",
		ResidentCountry:    "FRA",
		ResidentCity:       "Paris",
		PreferredLanguage:  "fr",
		ShortBio:           "DevOps Engineer",
		LongBio:            "DevOps engineer with cloud expertise",
		WorkHistoryDomains: []string{"nokiabricks.example", "sunvaja.example"},
	},
	{
		Name:              "User Five",
		Handle:            "user5",
		Email:             "user5@example.com",
		ResidentCountry:   "IND",
		ResidentCity:      "Bangalore",
		PreferredLanguage: "en",
		ShortBio:          "Data Scientist",
		LongBio:           "Data scientist with ML expertise",
		WorkHistoryDomains: []string{
			"novelltenware.example",
			"decdpd.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Six",
		Handle:            "user6",
		Email:             "user6@example.com",
		ResidentCountry:   "CAN",
		ResidentCity:      "Toronto",
		PreferredLanguage: "en",
		ShortBio:          "UX Designer",
		LongBio:           "UX designer focused on user-centered design",
		WorkHistoryDomains: []string{
			"decdpd.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:               "User Seven",
		Handle:             "user7",
		Email:              "user7@example.com",
		ResidentCountry:    "AUS",
		ResidentCity:       "Sydney",
		PreferredLanguage:  "en",
		ShortBio:           "System Architect",
		LongBio:            "System architect with distributed systems experience",
		WorkHistoryDomains: []string{"sunvaja.example", "nokiabricks.example"},
	},
	{
		Name:              "User Eight",
		Handle:            "user8",
		Email:             "user8@example.com",
		ResidentCountry:   "JPN",
		ResidentCity:      "Tokyo",
		PreferredLanguage: "ja",
		ShortBio:          "Mobile Developer",
		LongBio:           "Mobile developer specializing in iOS",
		WorkHistoryDomains: []string{
			"novelltenware.example",
			"decdpd.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Nine",
		Handle:            "user9",
		Email:             "user9@example.com",
		ResidentCountry:   "SGP",
		ResidentCity:      "Singapore",
		PreferredLanguage: "en",
		ShortBio:          "QA Engineer",
		LongBio:           "QA engineer with automation expertise",
		WorkHistoryDomains: []string{
			"decdpd.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
	},
	{
		Name:               "User Ten",
		Handle:             "user10",
		Email:              "user10@example.com",
		ResidentCountry:    "NLD",
		ResidentCity:       "Amsterdam",
		PreferredLanguage:  "nl",
		ShortBio:           "Security Engineer",
		LongBio:            "Security engineer focused on application security",
		WorkHistoryDomains: []string{"nokiabricks.example", "decdpd.example"},
	},
	{
		Name:              "User Eleven",
		Handle:            "user11",
		Email:             "user11@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "Boston",
		PreferredLanguage: "en",
		ShortBio:          "Backend Developer",
		LongBio:           "Backend developer with Java expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"decdpd.example",
			"nokiabricks.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user12"),
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Twelve",
		Handle:            "user12",
		Email:             "user12@example.com",
		ResidentCountry:   "GBR",
		ResidentCity:      "Manchester",
		PreferredLanguage: "en",
		ShortBio:          "Full Stack Developer",
		LongBio:           "Full stack developer with MEAN stack expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"sunvaja.example",
			"novelltenware.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Thirteen",
		Handle:            "user13",
		Email:             "user13@example.com",
		ResidentCountry:   "IND",
		ResidentCity:      "Mumbai",
		PreferredLanguage: "en",
		ShortBio:          "Backend Developer",
		LongBio:           "Backend developer with Go expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"nokiabricks.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Fourteen",
		Handle:            "user14",
		Email:             "user14@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "San Francisco",
		PreferredLanguage: "en",
		ShortBio:          "ML Engineer",
		LongBio:           "Machine learning engineer specializing in NLP",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"sunvaja.example",
			"novelltenware.example",
			"decdpd.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user11"),
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Fifteen",
		Handle:            "user15",
		Email:             "user15@example.com",
		ResidentCountry:   "CAN",
		ResidentCity:      "Vancouver",
		PreferredLanguage: "en",
		ShortBio:          "DevOps Engineer",
		LongBio:           "DevOps engineer with Kubernetes expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user11"),
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Sixteen",
		Handle:            "user16",
		Email:             "user16@example.com",
		ResidentCountry:   "AUS",
		ResidentCity:      "Melbourne",
		PreferredLanguage: "en",
		ShortBio:          "Frontend Developer",
		LongBio:           "Frontend developer with Vue.js expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"decdpd.example",
			"novelltenware.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user11"),
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Seventeen",
		Handle:            "user17",
		Email:             "user17@example.com",
		ResidentCountry:   "GBR",
		ResidentCity:      "Edinburgh",
		PreferredLanguage: "en",
		ShortBio:          "Data Engineer",
		LongBio:           "Data engineer with Apache Spark expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"nokiabricks.example",
			"sunvaja.example",
			"decdpd.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user11"),
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Eighteen",
		Handle:            "user18",
		Email:             "user18@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Hamburg",
		PreferredLanguage: "de",
		ShortBio:          "Security Engineer",
		LongBio:           "Security engineer with pentesting expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"novelltenware.example",
			"decdpd.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
		Endorsers: []common.Handle{
			common.Handle("user13"),
		},
	},
	{
		Name:              "User Nineteen",
		Handle:            "user19",
		Email:             "user19@example.com",
		ResidentCountry:   "FRA",
		ResidentCity:      "Marseille",
		PreferredLanguage: "fr",
		ShortBio:          "System Architect",
		LongBio:           "System architect with microservices expertise",
		WorkHistoryDomains: []string{
			"gryffindor.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
		ApplyToCompanyDomains: []string{
			"gryffindor.example",
		},
	},
	{
		Name:              "User Twenty",
		Handle:            "user20",
		Email:             "user20@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "Seattle",
		PreferredLanguage: "en",
		ShortBio:          "Cloud Architect",
		LongBio:           "Cloud architect with AWS expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"nokiabricks.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Twenty One",
		Handle:            "user21",
		Email:             "user21@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "Austin",
		PreferredLanguage: "en",
		ShortBio:          "Mobile Developer",
		LongBio:           "Mobile developer with Android expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"decdpd.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Twenty Two",
		Handle:            "user22",
		Email:             "user22@example.com",
		ResidentCountry:   "CAN",
		ResidentCity:      "Montreal",
		PreferredLanguage: "fr",
		ShortBio:          "UI Designer",
		LongBio:           "UI designer with focus on mobile apps",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"nokiabricks.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Twenty Three",
		Handle:            "user23",
		Email:             "user23@example.com",
		ResidentCountry:   "GBR",
		ResidentCity:      "Bristol",
		PreferredLanguage: "en",
		ShortBio:          "Backend Developer",
		LongBio:           "Backend developer with Python expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"sunvaja.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Twenty Four",
		Handle:            "user24",
		Email:             "user24@example.com",
		ResidentCountry:   "IND",
		ResidentCity:      "Hyderabad",
		PreferredLanguage: "en",
		ShortBio:          "Full Stack Developer",
		LongBio:           "Full stack developer with MERN stack expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"nokiabricks.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Twenty Five",
		Handle:            "user25",
		Email:             "user25@example.com",
		ResidentCountry:   "SGP",
		ResidentCity:      "Singapore",
		PreferredLanguage: "en",
		ShortBio:          "DevOps Engineer",
		LongBio:           "DevOps engineer with AWS expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"decdpd.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Twenty Six",
		Handle:            "user26",
		Email:             "user26@example.com",
		ResidentCountry:   "AUS",
		ResidentCity:      "Brisbane",
		PreferredLanguage: "en",
		ShortBio:          "Data Scientist",
		LongBio:           "Data scientist with deep learning expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Twenty Seven",
		Handle:            "user27",
		Email:             "user27@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Frankfurt",
		PreferredLanguage: "de",
		ShortBio:          "Frontend Developer",
		LongBio:           "Frontend developer with Angular expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Twenty Eight",
		Handle:            "user28",
		Email:             "user28@example.com",
		ResidentCountry:   "JPN",
		ResidentCity:      "Osaka",
		PreferredLanguage: "ja",
		ShortBio:          "Security Engineer",
		LongBio:           "Security engineer with cloud security expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"decdpd.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Twenty Nine",
		Handle:            "user29",
		Email:             "user29@example.com",
		ResidentCountry:   "FRA",
		ResidentCity:      "Nice",
		PreferredLanguage: "fr",
		ShortBio:          "System Architect",
		LongBio:           "System architect with cloud native expertise",
		WorkHistoryDomains: []string{
			"hufflepuff.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Thirty",
		Handle:            "user30",
		Email:             "user30@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Munich",
		PreferredLanguage: "de",
		ShortBio:          "ML Engineer",
		LongBio:           "Machine learning engineer with deep learning expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"sunvaja.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Thirty One",
		Handle:            "user31",
		Email:             "user31@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "Chicago",
		PreferredLanguage: "en",
		ShortBio:          "Backend Developer",
		LongBio:           "Backend developer with Node.js expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"decdpd.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Thirty Two",
		Handle:            "user32",
		Email:             "user32@example.com",
		ResidentCountry:   "IND",
		ResidentCity:      "Chennai",
		PreferredLanguage: "en",
		ShortBio:          "Mobile Developer",
		LongBio:           "Mobile developer with React Native expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"sunvaja.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Thirty Three",
		Handle:            "user33",
		Email:             "user33@example.com",
		ResidentCountry:   "GBR",
		ResidentCity:      "Leeds",
		PreferredLanguage: "en",
		ShortBio:          "DevOps Engineer",
		LongBio:           "DevOps engineer with GCP expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Thirty Four",
		Handle:            "user34",
		Email:             "user34@example.com",
		ResidentCountry:   "CAN",
		ResidentCity:      "Ottawa",
		PreferredLanguage: "en",
		ShortBio:          "Data Engineer",
		LongBio:           "Data engineer with Apache Kafka expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"decdpd.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Thirty Five",
		Handle:            "user35",
		Email:             "user35@example.com",
		ResidentCountry:   "AUS",
		ResidentCity:      "Perth",
		PreferredLanguage: "en",
		ShortBio:          "Frontend Developer",
		LongBio:           "Frontend developer with Svelte expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"nokiabricks.example",
			"sunvaja.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Thirty Six",
		Handle:            "user36",
		Email:             "user36@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Cologne",
		PreferredLanguage: "de",
		ShortBio:          "UX Designer",
		LongBio:           "UX designer with focus on web applications",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Thirty Seven",
		Handle:            "user37",
		Email:             "user37@example.com",
		ResidentCountry:   "SGP",
		ResidentCity:      "Singapore",
		PreferredLanguage: "en",
		ShortBio:          "System Architect",
		LongBio:           "System architect with serverless expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Thirty Eight",
		Handle:            "user38",
		Email:             "user38@example.com",
		ResidentCountry:   "JPN",
		ResidentCity:      "Kyoto",
		PreferredLanguage: "ja",
		ShortBio:          "ML Engineer",
		LongBio:           "Machine learning engineer with computer vision expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"decdpd.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Thirty Nine",
		Handle:            "user39",
		Email:             "user39@example.com",
		ResidentCountry:   "FRA",
		ResidentCity:      "Bordeaux",
		PreferredLanguage: "fr",
		ShortBio:          "Security Engineer",
		LongBio:           "Security engineer with DevSecOps expertise",
		WorkHistoryDomains: []string{
			"slytherin.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Forty",
		Handle:            "user40",
		Email:             "user40@example.com",
		ResidentCountry:   "FRA",
		ResidentCity:      "Lyon",
		PreferredLanguage: "fr",
		ShortBio:          "Data Engineer",
		LongBio:           "Data engineer with big data expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"nokiabricks.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Forty One",
		Handle:            "user41",
		Email:             "user41@example.com",
		ResidentCountry:   "USA",
		ResidentCity:      "Portland",
		PreferredLanguage: "en",
		ShortBio:          "Backend Developer",
		LongBio:           "Backend developer with Ruby expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"decdpd.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Forty Two",
		Handle:            "user42",
		Email:             "user42@example.com",
		ResidentCountry:   "IND",
		ResidentCity:      "Pune",
		PreferredLanguage: "en",
		ShortBio:          "Full Stack Developer",
		LongBio:           "Full stack developer with Django expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"sunvaja.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Forty Three",
		Handle:            "user43",
		Email:             "user43@example.com",
		ResidentCountry:   "GBR",
		ResidentCity:      "Glasgow",
		PreferredLanguage: "en",
		ShortBio:          "DevOps Engineer",
		LongBio:           "DevOps engineer with Azure expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
	{
		Name:              "User Forty Four",
		Handle:            "user44",
		Email:             "user44@example.com",
		ResidentCountry:   "CAN",
		ResidentCity:      "Calgary",
		PreferredLanguage: "en",
		ShortBio:          "Mobile Developer",
		LongBio:           "Mobile developer with Flutter expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"decdpd.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Forty Five",
		Handle:            "user45",
		Email:             "user45@example.com",
		ResidentCountry:   "AUS",
		ResidentCity:      "Adelaide",
		PreferredLanguage: "en",
		ShortBio:          "Data Scientist",
		LongBio:           "Data scientist with time series analysis expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"nokiabricks.example",
			"sunvaja.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Forty Six",
		Handle:            "user46",
		Email:             "user46@example.com",
		ResidentCountry:   "DEU",
		ResidentCity:      "Stuttgart",
		PreferredLanguage: "de",
		ShortBio:          "Frontend Developer",
		LongBio:           "Frontend developer with WebGL expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"novelltenware.example",
			"decdpd.example",
		},
	},
	{
		Name:              "User Forty Seven",
		Handle:            "user47",
		Email:             "user47@example.com",
		ResidentCountry:   "SGP",
		ResidentCity:      "Singapore",
		PreferredLanguage: "en",
		ShortBio:          "System Architect",
		LongBio:           "System architect with event-driven architecture expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"sunvaja.example",
			"nokiabricks.example",
		},
	},
	{
		Name:              "User Forty Eight",
		Handle:            "user48",
		Email:             "user48@example.com",
		ResidentCountry:   "JPN",
		ResidentCity:      "Sapporo",
		PreferredLanguage: "ja",
		ShortBio:          "Security Engineer",
		LongBio:           "Security engineer with blockchain security expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"decdpd.example",
			"novelltenware.example",
		},
	},
	{
		Name:              "User Forty Nine",
		Handle:            "user49",
		Email:             "user49@example.com",
		ResidentCountry:   "FRA",
		ResidentCity:      "Toulouse",
		PreferredLanguage: "fr",
		ShortBio:          "ML Engineer",
		LongBio:           "Machine learning engineer with reinforcement learning expertise",
		WorkHistoryDomains: []string{
			"ravenclaw.example",
			"nokiabricks.example",
			"sunvaja.example",
		},
	},
}

func initHubUsers(db *pgxpool.Pool) {
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
				preferred_language, short_bio, long_bio
			) VALUES (
				$1, $2, $3, $4,
				'$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
				'ACTIVE_HUB_USER', $5, $6, $7, $8, $9
			)
		`, userID, user.Name, user.Handle, user.Email,
			user.ResidentCountry, user.ResidentCity,
			user.PreferredLanguage, user.ShortBio, user.LongBio)
		if err != nil {
			log.Fatalf("failed to create hub user %s: %v", user.Name, err)
		}

		// Print with color directly
		color.New(color.FgGreen).Printf("created hub user %s ", user.Name)
		color.New(color.FgCyan).Printf("<%s> ", user.Email)
		color.New(color.FgYellow).Printf("@%s\n", user.Handle)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func loginHubUsers() {
	var wg sync.WaitGroup
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubUser) {
			color.Green("Logging in %s", user.Email)
			hubLogin(user.Email, "NewPassword123$", &wg)
		}(user)
	}

	// Wait for all hubLogin to complete
	wg.Wait()
}
