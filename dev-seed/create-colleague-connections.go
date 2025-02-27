package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func createColleagueConnections() {
	createConnectionRequests()
	approveConnectionRequests()
}

func createConnectionRequests() {
	var connections = []struct {
		HubUserEmail    string
		ColleagueHandle string
	}{
		// user11 connects to all other 1 series users
		{"user11@example.com", "user12"},
		{"user11@example.com", "user13"},
		{"user11@example.com", "user14"},
		{"user11@example.com", "user15"},
		{"user11@example.com", "user16"},
		{"user11@example.com", "user17"},
		{"user11@example.com", "user18"},
		{"user11@example.com", "user19"},

		// user13 is also connected to all 1 series users
		{"user13@example.com", "user12"},
		{"user13@example.com", "user14"},
		{"user13@example.com", "user15"},
		{"user13@example.com", "user16"},
		{"user13@example.com", "user17"},
		{"user13@example.com", "user18"},
		{"user13@example.com", "user19"},
	}
	for _, connection := range connections {
		createConnectionRequest(
			connection.HubUserEmail,
			connection.ColleagueHandle,
		)
	}
}

func createConnectionRequest(hubUserEmail, colleagueHandle string) {
	tokenI, ok := hubSessionTokens.Load(hubUserEmail)
	if !ok {
		log.Fatalf("no auth token found for %s", hubUserEmail)
	}
	authToken := tokenI.(string)

	connectionReq := hub.ConnectColleagueRequest{
		Handle: common.Handle(colleagueHandle),
	}

	connectionReqJSON, err := json.Marshal(connectionReq)
	if err != nil {
		log.Fatalf("failed to marshal connection request: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		serverURL+"/hub/connect-colleague",
		bytes.NewBuffer(connectionReqJSON),
	)
	if err != nil {
		log.Fatalf("failed to create connection request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send connection request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to send connection request: %v", resp.Status)
	}

	color.Blue("connect request from %s to %s", hubUserEmail, colleagueHandle)
}

func approveConnectionRequests() {
	var connectionRequests = []struct {
		ApproverEmail   string
		RequestorHandle string
	}{
		{"user12@example.com", "user11"},
		{"user13@example.com", "user11"},
		{"user14@example.com", "user11"},
		{"user15@example.com", "user11"},
		{"user16@example.com", "user11"},
		{"user17@example.com", "user11"},

		{"user12@example.com", "user13"},
		{"user14@example.com", "user13"},
		{"user15@example.com", "user13"},
		{"user16@example.com", "user13"},
		{"user17@example.com", "user13"},
		{"user18@example.com", "user13"},
		{"user19@example.com", "user13"},
	}
	for _, connectionRequest := range connectionRequests {
		approveConnectionRequest(
			connectionRequest.ApproverEmail,
			connectionRequest.RequestorHandle,
		)
	}
}

func approveConnectionRequest(approverEmail, requestorHandle string) {
	tokenI, ok := hubSessionTokens.Load(approverEmail)
	if !ok {
		log.Fatalf("no auth token found for %s", approverEmail)
	}
	authToken := tokenI.(string)

	approveColleagueReq := hub.ApproveColleagueRequest{
		Handle: common.Handle(requestorHandle),
	}

	approveColleagueReqJSON, err := json.Marshal(approveColleagueReq)
	if err != nil {
		log.Fatalf("failed to marshal approve colleague request: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		serverURL+"/hub/approve-colleague",
		bytes.NewBuffer(approveColleagueReqJSON),
	)
	if err != nil {
		log.Fatalf("failed to create approve colleague request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send approve colleague request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to send approve colleague request: %v", resp.Status)
	}

	color.Yellow("Approved req from %s to %s", approverEmail, requestorHandle)
}
