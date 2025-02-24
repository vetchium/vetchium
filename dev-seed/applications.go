package main

import (
	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/hub"
)

func createApplications() {
	for _, user := range hubUsers {
		for _, companyDomain := range user.ApplyToCompanyDomains {
			firstOpeningID := activeOpenings[companyDomain][0]
			createApplicationForOpening(user, companyDomain, firstOpeningID)

			secondOpeningID := activeOpenings[companyDomain][1]
			createApplicationForOpening(user, companyDomain, secondOpeningID)
		}
	}
}

func createApplicationForOpening(
	user HubUser,
	company string,
	openingID string,
) {
	color.Green(
		"Creating application for %q for %s/%s",
		user.Name,
		company,
		openingID,
	)

	// Get the user's session token
	tokenVal, ok := hubSessionTokens.Load(user.Email)
	if !ok {
		color.Red("Failed to get session token for %s", user.Email)
		return
	}
	token := tokenVal.(string)

	// Create the application request
	req := hub.ApplyForOpeningRequest{
		OpeningIDWithinCompany: openingID,
		CompanyDomain:          company,
		Resume:                 sampleResumePDF,
		CoverLetter:            "I am excited to apply for this position...",
		Filename:               user.Handle + "-resume.pdf",
	}

	// Make the API request and get the response
	var resp hub.ApplyForOpeningResponse
	makeRequest("POST", "/hub/apply-for-opening", token, req, &resp)
	color.Magenta("Successfully created application: %s", resp.ApplicationID)
}
