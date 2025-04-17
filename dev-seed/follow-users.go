package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/fatih/color"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func followUsers() {
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Fatalf("no auth token found for %s", user.Email)
		}
		authToken := tokenI.(string)

		usersToFollow := make(map[string]bool)
		for j := 0; j < 10; j++ {
			for {
				userNumToFollow := rand.Intn(len(hubUsers))
				userToFollow := hubUsers[userNumToFollow].Handle
				if userNumToFollow != i && !usersToFollow[userToFollow] {
					usersToFollow[userToFollow] = true
					break
				}
			}
		}

		for userIter := range usersToFollow {
			followUser(userIter, authToken)
		}

		color.Green("%s followed %+v", user.Handle, usersToFollow)
	}
}

func followUser(handle string, authToken string) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(hub.FollowUserRequest{
		Handle: common.Handle(handle),
	})
	if err != nil {
		log.Fatalf("failed to encode follow user request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/follow-user",
		&body,
	)
	if err != nil {
		log.Fatalf("failed to create follow user request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send follow user request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to follow user: %v", resp.Status)
	}
}
