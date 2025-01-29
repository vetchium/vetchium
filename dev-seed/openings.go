package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

// Track openings per company using company domain as key, openingID as value
var companyOpenings = make(map[string][]string)

// Track active openings per company using company domain as key, openingID as value
var activeOpenings = make(map[string][]string)

func createOpening(token string, req employer.CreateOpeningRequest) string {
	var resp employer.CreateOpeningResponse
	makeRequest("POST", "/employer/create-opening", token, req, &resp)
	return resp.OpeningID
}

func changeOpeningState(
	token string,
	openingID string,
	fromState, toState common.OpeningState,
) {
	req := employer.ChangeOpeningStateRequest{
		OpeningID: openingID,
		FromState: fromState,
		ToState:   toState,
	}
	makeRequest("POST", "/employer/change-opening-state", token, req, nil)
}

func createOpenings() {
	// Get tokens from the global map
	gryffindorVal, ok := employerSessionTokens.Load("admin@gryffindor.example")
	if !ok {
		log.Fatal("failed to get gryffindor token")
	}
	gryffindorToken := gryffindorVal.(string)

	hufflepuffVal, ok := employerSessionTokens.Load("admin@hufflepuff.example")
	if !ok {
		log.Fatal("failed to get hufflepuff token")
	}
	hufflepuffToken := hufflepuffVal.(string)

	ravenclawVal, ok := employerSessionTokens.Load("admin@ravenclaw.example")
	if !ok {
		log.Fatal("failed to get ravenclaw token")
	}
	ravenclawToken := ravenclawVal.(string)

	slytherinVal, ok := employerSessionTokens.Load("admin@slytherin.example")
	if !ok {
		log.Fatal("failed to get slytherin token")
	}
	slytherinToken := slytherinVal.(string)

	openings := []struct {
		domain string
		token  string
		req    employer.CreateOpeningRequest
	}{
		// Gryffindor openings
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Senior Backend Engineer",
				Positions:         2,
				JD:                "Looking for experienced backend engineers to join our UK team. Must have strong Go experience.",
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "UK Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Product Manager",
				Positions:         1,
				JD:                "Seeking an experienced product manager to lead our Irish expansion.",
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "ron@gryffindor.example",
				CostCenterName:    "Ireland Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Diagon"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "DevOps Engineer",
				Positions:         2,
				JD:                "Looking for DevOps engineers to support our global infrastructure.",
				Recruiter:         "ron@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "APAC Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Frontend Developer",
				Positions:         3,
				JD:                "Seeking React developers for our Canadian office.",
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "ron@gryffindor.example",
				CostCenterName:    "Canada Business",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            6,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Marketing Lead",
				Positions:         1,
				JD:                "Looking for a marketing lead to oversee global campaigns.",
				Recruiter:         "ron@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "Global Marketing",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            10,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Diagon"},
			},
		},

		// Hufflepuff openings
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Software Architect",
				Positions:         1,
				JD:                "Seeking a software architect for our Benelux operations.",
				Recruiter:         "cedric@hufflepuff.example",
				HiringManager:     "newt@hufflepuff.example",
				CostCenterName:    "Benelux Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Data Scientist",
				Positions:         2,
				JD:                "Looking for data scientists to join our Nordic team.",
				Recruiter:         "nymphadora@hufflepuff.example",
				HiringManager:     "cedric@hufflepuff.example",
				CostCenterName:    "Nordic Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Marketing Manager",
				Positions:         1,
				JD:                "Seeking a marketing manager for EU operations.",
				Recruiter:         "newt@hufflepuff.example",
				HiringManager:     "nymphadora@hufflepuff.example",
				CostCenterName:    "EU Marketing",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Research Engineer",
				Positions:         3,
				JD:                "Join our R&D team working on cutting-edge technology.",
				Recruiter:         "cedric@hufflepuff.example",
				HiringManager:     "newt@hufflepuff.example",
				CostCenterName:    "R&D Labs",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            4,
				YoeMax:            12,
				MinEducationLevel: common.DoctorateEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Operations Manager",
				Positions:         1,
				JD:                "Looking for an operations manager for EU administration.",
				Recruiter:         "nymphadora@hufflepuff.example",
				HiringManager:     "cedric@hufflepuff.example",
				CostCenterName:    "EU Admin",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            6,
				YoeMax:            12,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
			},
		},

		// Ravenclaw openings
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Technical Lead",
				Positions:         1,
				JD:                "Seeking a technical lead for our APAC headquarters.",
				Recruiter:         "luna@ravenclaw.example",
				HiringManager:     "filius@ravenclaw.example",
				CostCenterName:    "APAC Headquarters",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Flourish"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Mobile Developer",
				Positions:         2,
				JD:                "Looking for iOS/Android developers for our Japan team.",
				Recruiter:         "cho@ravenclaw.example",
				HiringManager:     "luna@ravenclaw.example",
				CostCenterName:    "Japan Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "QA Engineer",
				Positions:         2,
				JD:                "Join our Korean QA team ensuring product quality.",
				Recruiter:         "filius@ravenclaw.example",
				HiringManager:     "cho@ravenclaw.example",
				CostCenterName:    "Korea Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            6,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Full Stack Developer",
				Positions:         3,
				JD:                "Seeking full stack developers for our India operations.",
				Recruiter:         "luna@ravenclaw.example",
				HiringManager:     "filius@ravenclaw.example",
				CostCenterName:    "India Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            4,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Solutions Architect",
				Positions:         1,
				JD:                "Looking for a solutions architect for Middle East expansion.",
				Recruiter:         "cho@ravenclaw.example",
				HiringManager:     "luna@ravenclaw.example",
				CostCenterName:    "Middle East Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            10,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Flourish"},
			},
		},

		// Slytherin openings
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Engineering Manager",
				Positions:         1,
				JD:                "Seeking an engineering manager for DACH region.",
				Recruiter:         "draco@slytherin.example",
				HiringManager:     "severus@slytherin.example",
				CostCenterName:    "DACH Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Security Engineer",
				Positions:         2,
				JD:                "Join our French security team.",
				Recruiter:         "severus@slytherin.example",
				HiringManager:     "horace@slytherin.example",
				CostCenterName:    "France Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Cloud Engineer",
				Positions:         2,
				JD:                "Looking for cloud engineers for Southern Europe operations.",
				Recruiter:         "horace@slytherin.example",
				HiringManager:     "draco@slytherin.example",
				CostCenterName:    "Southern Europe",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Project Manager",
				Positions:         1,
				JD:                "Seeking a project manager for EU special projects.",
				Recruiter:         "draco@slytherin.example",
				HiringManager:     "severus@slytherin.example",
				CostCenterName:    "EU Projects",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            6,
				YoeMax:            12,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Research Scientist",
				Positions:         2,
				JD:                "Join our continental R&D team.",
				Recruiter:         "severus@slytherin.example",
				HiringManager:     "horace@slytherin.example",
				CostCenterName:    "Continental R&D",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            12,
				MinEducationLevel: common.DoctorateEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
			},
		},
	}

	for _, opening := range openings {
		openingID := createOpening(opening.token, opening.req)
		color.Green("Created opening %s for %s", openingID, opening.domain)
		// Track openings by domain
		companyOpenings[opening.domain] = append(
			companyOpenings[opening.domain],
			openingID,
		)
	}

	// Publish first two openings for each company
	for domain, openings := range companyOpenings {
		employerTokenRaw, ok := employerSessionTokens.Load("admin@" + domain)
		if !ok {
			log.Fatalf("failed to get employer token for %s", domain)
		}
		employerToken, ok := employerTokenRaw.(string)
		if !ok {
			log.Fatalf("failed to cast employer token for %s", domain)
		}

		for i := 0; i < 2 && i < len(openings); i++ {
			changeOpeningState(
				employerToken,
				openings[i],
				common.DraftOpening,
				common.ActiveOpening,
			)
			color.Green("Published opening %s for %s", openings[i], domain)
			activeOpenings[domain] = append(activeOpenings[domain], openings[i])
		}

	}
}
