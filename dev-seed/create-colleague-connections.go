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
		createColleagueConnection(
			connection.HubUserEmail,
			connection.ColleagueHandle,
		)
	}
}

func createColleagueConnection(hubUserEmail, colleagueHandle string) {
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
