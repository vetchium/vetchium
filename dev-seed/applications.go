package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

// Variables to store opening information
var openingsByCompany map[string][]OpeningInfo
var allOpenings []OpeningInfo

type OpeningInfo struct {
	CompanyDomain string
	OpeningID     string
	Title         string
	JobLevel      string
	Description   string
	Skills        []string
}

func createApplications() {
	// Fetch all openings from the server
	fetchAllOpenings()

	color.Cyan("Found %d openings across %d companies",
		len(allOpenings), len(openingsByCompany))

	// Track applications count
	applicationsCreated := 0
	maxApplications := 100 // Limit the total number of applications

	// Create applications for each user
	for _, user := range hubUsers {
		// Skip if we've reached the max applications limit
		if applicationsCreated >= maxApplications {
			break
		}

		// Find suitable openings for this user based on their experience
		suitableOpenings := findSuitableOpenings(user)

		// Limit applications per user to 3 maximum
		if len(suitableOpenings) > 3 {
			// Randomly select 3 openings
			rand.Shuffle(len(suitableOpenings), func(i, j int) {
				suitableOpenings[i], suitableOpenings[j] = suitableOpenings[j], suitableOpenings[i]
			})
			suitableOpenings = suitableOpenings[:3]
		}

		if len(suitableOpenings) == 0 {
			continue
		}

		color.Blue("Found %d suitable openings for %s (%s)",
			len(suitableOpenings), user.Name, user.ShortBio)

		// Generate a PDF resume for this user
		resumeFilename, resumeData, err := generateResumePDF(user)
		if err != nil {
			log.Printf("Failed to generate resume for %s: %v", user.Email, err)
			continue
		}

		// Apply to each suitable opening
		for _, opening := range suitableOpenings {
			// Find endorsers who worked at the same company
			endorsers := findEndorsers(user, opening.CompanyDomain)

			// Create the application
			err := createApplicationForOpening(
				user,
				opening.CompanyDomain,
				opening.OpeningID,
				resumeFilename,
				resumeData,
				endorsers,
			)

			if err != nil {
				log.Printf("Failed to create application for %s to %s/%s: %v",
					user.Email, opening.CompanyDomain, opening.OpeningID, err)
				continue
			}

			applicationsCreated++
		}
	}

	color.Green("Created %d applications for job openings", applicationsCreated)
}

// findSuitableOpenings finds openings that match the user's profile
func findSuitableOpenings(user HubSeedUser) []OpeningInfo {
	var suitable []OpeningInfo

	// Get the user's current job title and history
	currentJob := user.Jobs[len(user.Jobs)-1]

	// Extract keywords from the job title
	jobKeywords := extractKeywords(currentJob.Title)

	// For each opening, check if it's suitable
	for _, opening := range allOpenings {
		score := 0

		// Match by job title keywords
		openingKeywords := extractKeywords(opening.Title)
		for _, keyword := range jobKeywords {
			for _, openingKeyword := range openingKeywords {
				if strings.Contains(
					strings.ToLower(openingKeyword),
					strings.ToLower(keyword),
				) ||
					strings.Contains(
						strings.ToLower(keyword),
						strings.ToLower(openingKeyword),
					) {
					score += 2
				}
			}
		}

		// Check if user previously worked at this company
		for _, job := range user.Jobs {
			if job.Website == opening.CompanyDomain {
				score += 3 // Higher score for previous experience at the company
			}
		}

		// Check for career path progression
		// If the opening title contains words like "Senior", "Lead", "Manager" and
		// the user's current title doesn't, consider it a career advancement
		if (strings.Contains(strings.ToLower(opening.Title), "senior") ||
			strings.Contains(strings.ToLower(opening.Title), "lead") ||
			strings.Contains(strings.ToLower(opening.Title), "manager")) &&
			!strings.Contains(strings.ToLower(currentJob.Title), "senior") &&
			!strings.Contains(strings.ToLower(currentJob.Title), "lead") &&
			!strings.Contains(strings.ToLower(currentJob.Title), "manager") {
			score += 2
		}

		// If score is high enough, consider it suitable
		if score >= 2 {
			suitable = append(suitable, opening)
		}
	}

	return suitable
}

// findEndorsers finds users who can endorse an application to a specific company
func findEndorsers(user HubSeedUser, companyDomain string) []common.Handle {
	var endorsers []common.Handle

	// Find users who work or worked at the target company
	for _, otherUser := range hubUsers {
		// Skip the applicant
		if otherUser.Email == user.Email {
			continue
		}

		// Check if this user worked at the target company
		workedAtCompany := false
		for _, job := range otherUser.Jobs {
			if job.Website == companyDomain {
				workedAtCompany = true
				break
			}
		}

		if workedAtCompany {
			// Add as potential endorser (with 50% probability to simulate real-world scenarios)
			if rand.Intn(2) == 0 {
				endorsers = append(endorsers, common.Handle(otherUser.Handle))
			}
		}
	}

	// Limit to 3 endorsers maximum
	if len(endorsers) > 3 {
		endorsers = endorsers[:3]
	}

	return endorsers
}

// extractKeywords splits a job title into individual words for matching
func extractKeywords(title string) []string {
	// Remove common words that aren't useful for matching
	commonWords := map[string]bool{
		"the": true, "and": true, "of": true, "in": true, "at": true,
		"for": true, "with": true, "to": true, "a": true, "an": true,
	}

	// Split and filter
	words := strings.Fields(title)
	var keywords []string

	for _, word := range words {
		word = strings.ToLower(word)
		if !commonWords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// fetchAllOpenings retrieves all active openings from the server
func fetchAllOpenings() {
	// Initialize maps
	openingsByCompany = make(map[string][]OpeningInfo)
	allOpenings = []OpeningInfo{}

	// Copy from activeOpenings (which is loaded by createOpenings() function)
	for companyDomain, openingIDs := range activeOpenings {
		for _, openingID := range openingIDs {
			// Create an opening info object with basic information
			// In a real implementation, we would fetch more details from the API
			opening := OpeningInfo{
				CompanyDomain: companyDomain,
				OpeningID:     openingID,
				Title:         getRandomJobTitle(companyDomain),
				JobLevel:      getRandomJobLevel(),
				Description:   "This is a job opening at " + companyDomain,
				Skills:        []string{},
			}

			// Add to our collections
			openingsByCompany[companyDomain] = append(
				openingsByCompany[companyDomain],
				opening,
			)
			allOpenings = append(allOpenings, opening)
		}
	}
}

// Helper function to get a random job title for an opening
func getRandomJobTitle(companyDomain string) string {
	// Find a job title that matches the company's domain
	for _, employer := range employers {
		if employer.Website == companyDomain {
			// Get one of the career paths matching this employer's tags
			var matchingPaths []CareerPath
			for _, careerPath := range careerPaths {
				for _, tag := range employer.Tags {
					if careerPath.Tag == tag {
						matchingPaths = append(matchingPaths, careerPath)
						break
					}
				}
			}

			if len(matchingPaths) > 0 {
				// Choose a random career path
				path := matchingPaths[rand.Intn(len(matchingPaths))]
				// Choose a random job title from this path
				return path.Steps[rand.Intn(len(path.Steps))]
			}
		}
	}

	// Fallback to generic titles
	titles := []string{
		"Software Engineer", "Product Manager", "Data Scientist",
		"Marketing Specialist", "HR Manager", "Operations Manager",
	}
	return titles[rand.Intn(len(titles))]
}

// Helper function to get a random job level
func getRandomJobLevel() string {
	levels := []string{"Entry", "Mid", "Senior", "Lead", "Manager", "Director"}
	return levels[rand.Intn(len(levels))]
}

// createApplicationForOpening creates an application for a specific opening
func createApplicationForOpening(
	user HubSeedUser,
	companyDomain string,
	openingID string,
	resumeFilename string,
	resumeData []byte,
	endorsers []common.Handle,
) error {
	color.Green(
		"Creating application for %q for %s/%s",
		user.Name,
		companyDomain,
		openingID,
	)

	// Get the user's session token
	tokenVal, ok := hubSessionTokens.Load(user.Email)
	if !ok {
		return fmt.Errorf("failed to get session token for %s", user.Email)
	}
	token := tokenVal.(string)

	// Generate a custom cover letter based on the user's experience
	coverLetter := generateCoverLetter(user, companyDomain, openingID)

	// Create the application request
	req := hub.ApplyForOpeningRequest{
		OpeningIDWithinCompany: openingID,
		CompanyDomain:          companyDomain,
		Resume:                 string(resumeData),
		CoverLetter:            coverLetter,
		Filename:               filepath.Base(resumeFilename),
		EndorserHandles:        endorsers,
	}

	// Make the API request and get the response
	var resp hub.ApplyForOpeningResponse
	err := sendRequest("POST", "/hub/apply-for-opening", token, req, &resp)
	if err != nil {
		return err
	}

	color.Magenta("Successfully created application: %s with %d endorsers",
		resp.ApplicationID, len(endorsers))

	return nil
}

// generateCoverLetter creates a custom cover letter for an application
func generateCoverLetter(
	user HubSeedUser,
	companyDomain string,
	openingID string,
) string {
	// Get company name from domain
	companyName := companyDomain
	for _, employer := range employers {
		if employer.Website == companyDomain {
			companyName = employer.Name
			break
		}
	}

	// Get user's most recent job
	currentJob := user.Jobs[len(user.Jobs)-1]

	// Check if user worked at this company before
	workedAtCompany := false
	for _, job := range user.Jobs {
		if job.Website == companyDomain {
			workedAtCompany = true
			break
		}
	}

	// Generate the cover letter
	var coverLetter strings.Builder

	coverLetter.WriteString(
		fmt.Sprintf("Dear %s Hiring Team,\n\n", companyName),
	)

	coverLetter.WriteString(
		fmt.Sprintf(
			"I am excited to apply for the opening at %s. ",
			companyName,
		),
	)

	if workedAtCompany {
		coverLetter.WriteString(
			fmt.Sprintf(
				"Having previously worked at %s, I am particularly drawn to the opportunity to rejoin the team and contribute my skills once again. ",
				companyName,
			),
		)
	} else {
		coverLetter.WriteString(fmt.Sprintf("As a %s, I believe my skills and experience align well with the requirements of this position. ", currentJob.Title))
	}

	coverLetter.WriteString(fmt.Sprintf("\n\n%s\n\n", user.LongBio))

	coverLetter.WriteString(
		"I am particularly interested in this role because it allows me to leverage my experience while taking on new challenges. ",
	)

	coverLetter.WriteString(
		fmt.Sprintf(
			"\n\nThank you for considering my application. I look forward to the opportunity to discuss how I can contribute to %s's continued success.\n\n",
			companyName,
		),
	)

	coverLetter.WriteString(fmt.Sprintf("Sincerely,\n%s", user.Name))

	return coverLetter.String()
}

// makeRequest is a helper function to make API requests
func sendRequest(
	method, endpoint, token string,
	reqBody interface{},
	respBody interface{},
) error {
	// Marshal the request body
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create the request
	req, err := http.NewRequest(
		method,
		serverURL+endpoint,
		bytes.NewBuffer(reqData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(
			"request failed with status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	// Parse the response
	if respBody != nil {
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		err = json.Unmarshal(respData, respBody)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %v", err)
		}
	}

	return nil
}
