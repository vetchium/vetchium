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
	color.Cyan(
		"Finding potential colleague connections based on work history overlap",
	)

	// Track connections to create and avoid duplicates
	processedPairs := make(map[string]bool)
	connectionsToCreate := make([]struct {
		HubUserEmail    string
		ColleagueHandle string
	}, 0)

	// Find overlapping work periods between users
	for email1, workHistory1 := range WorkHistoryMap {
		for email2, workHistory2 := range WorkHistoryMap {
			// Skip self-connections
			if email1 == email2 {
				continue
			}

			// Create a unique key for this user pair to avoid duplicate processing
			pairKey := getPairKey(email1, email2)
			if processedPairs[pairKey] {
				continue
			}

			processedPairs[pairKey] = true

			// Find overlapping work periods
			for _, item1 := range workHistory1 {
				for _, item2 := range workHistory2 {
					// Check if they worked at the same employer
					if item1.EmployerDomain == item2.EmployerDomain {
						// Check if time periods overlap
						if periodsOverlap(
							item1.StartDate,
							item1.EndDate,
							item2.StartDate,
							item2.EndDate,
						) {
							handle2 := findHandleByEmail(email2)

							color.Yellow(
								"Found connection: %s and %s worked together at %s",
								email1,
								email2,
								item1.EmployerDomain,
							)

							// Add connection from email1 to handle2
							connectionsToCreate = append(
								connectionsToCreate,
								struct {
									HubUserEmail    string
									ColleagueHandle string
								}{
									HubUserEmail:    email1,
									ColleagueHandle: handle2,
								},
							)

							break // No need to check other work items for this pair
						}
					}
				}
			}
		}
	}

	// Limit the number of connections to prevent overwhelming the system
	maxConnections := 100
	if len(connectionsToCreate) > maxConnections {
		color.Magenta(
			"Limiting connections from %d to %d",
			len(connectionsToCreate),
			maxConnections,
		)
		connectionsToCreate = connectionsToCreate[:maxConnections]
	}

	// Create the connection requests
	successCount := 0
	for _, connection := range connectionsToCreate {
		success := createConnectionRequest(
			connection.HubUserEmail,
			connection.ColleagueHandle,
		)
		if success {
			successCount++

			// Check to see if we need to approve the request
			approved := checkAndApproveRequest(
				findEmailByHandle(connection.ColleagueHandle),
				findHandleByEmail(connection.HubUserEmail),
			)
			if approved {
				successCount++
			} else {
				color.Red(
					"Failed to approve connection from %s to %s",
					connection.HubUserEmail,
					connection.ColleagueHandle,
				)
			}
		} else {
			color.Red(
				"Failed to create connection from %s to %s",
				connection.HubUserEmail,
				connection.ColleagueHandle,
			)
		}
	}

	color.Green(
		"Created %d colleague connections",
		successCount,
	)
}

// Create a unique key for a pair of users
func getPairKey(email1, email2 string) string {
	if email1 < email2 {
		return email1 + "|" + email2
	}
	return email2 + "|" + email1
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
func findEmailByEmail(email string) string {
	for _, user := range hubUsers {
		if user.Email == email {
			return user.Email
		}
	}
	return ""
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

// Create connection request between two users
func createConnectionRequest(hubUserEmail, colleagueHandle string) bool {
	tokenI, ok := hubSessionTokens.Load(hubUserEmail)
	if !ok {
		log.Printf(
			"No auth token found for %s, skipping connection request",
			hubUserEmail,
		)
		return false
	}
	authToken := tokenI.(string)

	connectionReq := hub.ConnectColleagueRequest{
		Handle: common.Handle(colleagueHandle),
	}

	connectionReqJSON, err := json.Marshal(connectionReq)
	if err != nil {
		log.Printf("Failed to marshal connection request: %v", err)
		return false
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/connect-colleague",
		bytes.NewBuffer(connectionReqJSON),
	)
	if err != nil {
		log.Printf("Failed to create connection request: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send connection request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf(
			"Failed to create connection from %s to %s: %s",
			hubUserEmail,
			colleagueHandle,
			resp.Status,
		)
		return false
	}

	color.Blue("Connect request from %s to %s", hubUserEmail, colleagueHandle)
	return true
}

// Check for and approve connection requests
func checkAndApproveRequest(approverEmail, requestorHandle string) bool {
	if approverEmail == "" || requestorHandle == "" {
		log.Printf(
			"Invalid approver email or requestor handle, skipping approval",
		)
		return false
	}

	tokenI, ok := hubSessionTokens.Load(approverEmail)
	if !ok {
		log.Printf(
			"No auth token found for %s, skipping approval",
			approverEmail,
		)
		return false
	}
	authToken := tokenI.(string)

	// Get pending approvals
	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(
		hub.MyColleagueApprovalsRequest{},
	); err != nil {
		log.Printf("Failed to encode request body: %v", err)
		return false
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/my-colleague-approvals",
		&reqBody,
	)
	if err != nil {
		log.Printf("Failed to create approvals request: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to get approvals: %v", err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get approvals: %s", resp.Status)
		return false
	}

	var approvals hub.MyColleagueApprovals
	if err := json.NewDecoder(resp.Body).Decode(&approvals); err != nil {
		resp.Body.Close()
		log.Printf("Failed to decode approvals: %v", err)
		return false
	}
	resp.Body.Close()

	// Check if the requestor is in the pending approvals
	found := false
	for _, approval := range approvals.Approvals {
		if approval.Handle == common.Handle(requestorHandle) {
			found = true
			break
		}
	}

	if !found {
		color.Red(
			"Requestor %s not found in pending approvals for %s",
			requestorHandle,
			approverEmail,
		)
		return false
	}

	// Approve the request
	approveReq := hub.ApproveColleagueRequest{
		Handle: common.Handle(requestorHandle),
	}

	approveReqJSON, err := json.Marshal(approveReq)
	if err != nil {
		log.Printf("Failed to marshal approve request: %v", err)
		return false
	}

	req, err = http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/approve-colleague",
		bytes.NewBuffer(approveReqJSON),
	)
	if err != nil {
		log.Printf("Failed to create approve request: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send approve request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to approve connection request: %s", resp.Status)
		return false
	}

	color.Green(
		"Approved connection from %s to %s",
		requestorHandle,
		approverEmail,
	)
	return true
}
