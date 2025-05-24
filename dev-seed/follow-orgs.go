package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/fatih/color"
	"github.com/vetchium/vetchium/typespec/hub"
)

func followOrgs() {
	// List of available employer domains
	employerDomains := []string{
		"gryffindor.example",
		"hufflepuff.example",
		"ravenclaw.example",
		"slytherin.example",
		"sunvaja.example",
		"novelltenware.example",
		"decdpd.example",
		"nokiabricks.example",
	}

	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Fatalf("no auth token found for %s", user.Email)
		}
		authToken := tokenI.(string)

		// Each user follows 2-4 random employers
		numOrgsToFollow := 2 + rand.Intn(3) // 2-4 orgs
		orgsToFollow := make(map[string]bool)

		for j := 0; j < numOrgsToFollow; j++ {
			for {
				orgToFollow := employerDomains[rand.Intn(len(employerDomains))]
				if !orgsToFollow[orgToFollow] {
					orgsToFollow[orgToFollow] = true
					break
				}
			}
		}

		for orgDomain := range orgsToFollow {
			followOrg(orgDomain, authToken)
		}

		var orgList []string
		for org := range orgsToFollow {
			orgList = append(orgList, org)
		}
		color.Green("%s followed orgs: %v", user.Handle, orgList)
	}
}

func followOrg(domain string, authToken string) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(hub.FollowOrgRequest{
		Domain: domain,
	})
	if err != nil {
		log.Fatalf("failed to encode follow org request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/follow-org",
		&body,
	)
	if err != nil {
		log.Fatalf("failed to create follow org request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send follow org request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to follow org %s: %v", domain, resp.Status)
	}
}
