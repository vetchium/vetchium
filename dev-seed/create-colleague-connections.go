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
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user12",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user13",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user14",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user15",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user16",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user17",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user18",
		},
		{
			HubUserEmail:    "user11@example.com",
			ColleagueHandle: "user19",
		},
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
		{
			ApproverEmail:   "user12@example.com",
			RequestorHandle: "user11",
		},
		{
			ApproverEmail:   "user13@example.com",
			RequestorHandle: "user11",
		},
		{
			ApproverEmail:   "user14@example.com",
			RequestorHandle: "user11",
		},
		{
			ApproverEmail:   "user15@example.com",
			RequestorHandle: "user11",
		},
		{
			ApproverEmail:   "user16@example.com",
			RequestorHandle: "user11",
		},
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
