package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func createColleagueConnections() {
	// Map to store users by employer
	employerToUsers := make(map[string][]struct {
		UserEmail string
		StartDate time.Time
		EndDate   *time.Time
	})

	// Group users by employer
	for _, user := range hubUsers {
		for _, workItem := range user.WorkHistoryItems {
			employerToUsers[workItem.EmployerID] = append(
				employerToUsers[workItem.EmployerID],
				struct {
					UserEmail string
					StartDate time.Time
					EndDate   *time.Time
				}{
					UserEmail: user.Email,
					StartDate: workItem.StartDate,
					EndDate:   workItem.EndDate,
				},
			)
		}
	}

	// Track connections to be created
	connectionsToCreate := make([]struct {
		HubUserEmail    string
		ColleagueHandle string
	}, 0)

	// Find users with overlapping work periods
	connectionCount := 0
	for employerID, users := range employerToUsers {
		// Skip employers with only one user
		if len(users) <= 1 {
			continue
		}

		employerName := "Unknown"
		for _, employer := range employers {
			if employer.Website == employerID {
				employerName = employer.Name
				break
			}
		}

		color.Blue(
			"Finding colleague connections at %s (%s)",
			employerName,
			employerID,
		)

		// Create connections between users with overlapping work periods
		for i := 0; i < len(users); i++ {
			user1 := users[i]

			// Find user1's handle
			var handle1 string
			for _, hubUser := range hubUsers {
				if hubUser.Email == user1.UserEmail {
					handle1 = hubUser.Handle
					break
				}
			}

			for j := i + 1; j < len(users); j++ {
				user2 := users[j]

				// Find user2's handle
				var handle2 string
				for _, hubUser := range hubUsers {
					if hubUser.Email == user2.UserEmail {
						handle2 = hubUser.Handle
						break
					}
				}

				// Check if they worked together (overlapping time periods)
				if periodsOverlap(
					user1.StartDate,
					user1.EndDate,
					user2.StartDate,
					user2.EndDate,
				) {
					color.Yellow(
						"Found connection: %s and %s worked together at %s",
						user1.UserEmail,
						user2.UserEmail,
						employerName,
					)

					// Add connections both ways
					connectionsToCreate = append(connectionsToCreate, struct {
						HubUserEmail    string
						ColleagueHandle string
					}{
						HubUserEmail:    user1.UserEmail,
						ColleagueHandle: handle2,
					})

					connectionsToCreate = append(connectionsToCreate, struct {
						HubUserEmail    string
						ColleagueHandle string
					}{
						HubUserEmail:    user2.UserEmail,
						ColleagueHandle: handle1,
					})

					connectionCount += 2
				}
			}
		}
	}

	// Limit the number of connections to prevent overwhelming the system
	maxConnections := 100
	if connectionCount > maxConnections {
		color.Magenta(
			"Limiting connections from %d to %d",
			connectionCount,
			maxConnections,
		)
		connectionsToCreate = connectionsToCreate[:maxConnections]
	}

	// Create the connection requests
	for _, connection := range connectionsToCreate {
		createConnectionRequest(
			connection.HubUserEmail,
			connection.ColleagueHandle,
		)
	}

	// Approve all connection requests
	for _, connection := range connectionsToCreate {
		// The approval is done by the person who received the request
		approveConnectionRequest(
			findEmailByHandle(connection.ColleagueHandle),
			findHandleByEmail(connection.HubUserEmail),
		)
	}

	color.Green(
		"Created and approved %d colleague connections",
		len(connectionsToCreate),
	)
}

// Check if two time periods overlap
func periodsOverlap(
	start1 time.Time,
	end1 *time.Time,
	start2 time.Time,
	end2 *time.Time,
) bool {
	// If end1 is nil, user1 still works there (current job)
	end1Time := time.Now()
	if end1 != nil {
		end1Time = *end1
	}

	// If end2 is nil, user2 still works there (current job)
	end2Time := time.Now()
	if end2 != nil {
		end2Time = *end2
	}

	// Check overlap: one range doesn't end before the other starts
	return !end1Time.Before(start2) && !end2Time.Before(start1)
}

// Helper function to find email by handle
func findEmailByHandle(handle string) string {
	for _, user := range hubUsers {
		if user.Handle == handle {
			return user.Email
		}
	}
	return ""
}

// Helper function to find handle by email
func findHandleByEmail(email string) string {
	for _, user := range hubUsers {
		if user.Email == email {
			return user.Handle
		}
	}
	return ""
}

func createConnectionRequest(hubUserEmail, colleagueHandle string) {
	tokenI, ok := hubSessionTokens.Load(hubUserEmail)
	if !ok {
		log.Printf(
			"No auth token found for %s, skipping connection request",
			hubUserEmail,
		)
		return
	}
	authToken := tokenI.(string)

	connectionReq := hub.ConnectColleagueRequest{
		Handle: common.Handle(colleagueHandle),
	}

	connectionReqJSON, err := json.Marshal(connectionReq)
	if err != nil {
		log.Printf("Failed to marshal connection request: %v", err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		serverURL+"/hub/connect-colleague",
		bytes.NewBuffer(connectionReqJSON),
	)
	if err != nil {
		log.Printf("Failed to create connection request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send connection request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send connection request: %v", resp.Status)
		return
	}

	color.Blue("Connect request from %s to %s", hubUserEmail, colleagueHandle)
}

func approveConnectionRequest(approverEmail, requestorHandle string) {
	if approverEmail == "" || requestorHandle == "" {
		log.Printf(
			"Invalid approver email or requestor handle, skipping approval",
		)
		return
	}

	tokenI, ok := hubSessionTokens.Load(approverEmail)
	if !ok {
		log.Printf(
			"No auth token found for %s, skipping approval",
			approverEmail,
		)
		return
	}
	authToken := tokenI.(string)

	approveColleagueReq := hub.ApproveColleagueRequest{
		Handle: common.Handle(requestorHandle),
	}

	approveColleagueReqJSON, err := json.Marshal(approveColleagueReq)
	if err != nil {
		log.Printf("Failed to marshal approve colleague request: %v", err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		serverURL+"/hub/approve-colleague",
		bytes.NewBuffer(approveColleagueReqJSON),
	)
	if err != nil {
		log.Printf("Failed to create approve colleague request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send approve colleague request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send approve colleague request: %v", resp.Status)
		return
	}

	color.Yellow(
		"Approved request from %s to %s",
		approverEmail,
		requestorHandle,
	)
}
