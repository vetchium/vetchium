package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func writePosts() {
	var wg sync.WaitGroup
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%2 == 0 {
			color.Yellow("thread %d waiting", i)
			wg.Wait()
			color.Yellow("thread %d resumes", i)
		}
		wg.Add(1)

		// Parallelism to have more messages from many authors at similar times
		go func(user HubSeedUser, i int) {
			defer wg.Done()
			postUniverse := broadAreaPostsMap[user.BroadArea]
			numPosts := rand.Intn(len(postUniverse))

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			// Determine if user is free tier
			isFreeTier := user.Tier == hub.FreeHubUserTier

			for j := 0; j < numPosts; j++ {
				post := postUniverse[rand.Intn(len(postUniverse))]

				// Prepare request based on user tier
				var body bytes.Buffer
				var endpoint string
				var err error

				randomTagIDs := []common.VTagID{}
				for i := 0; i < rand.Intn(3); i++ {
					randomTagIDs = append(
						randomTagIDs,
						tagsList[rand.Intn(len(tagsList))],
					)
				}

				if isFreeTier {
					// Free tier users: use AddFTPostRequest (255 char limit, existing tags only)
					content := post.Content
					if len(content) > 255 {
						content = content[:252] + "..." // Truncate to fit limit
					}

					ftPost := hub.AddFTPostRequest{
						Content: content,
						TagIDs:  randomTagIDs,
					}
					err = json.NewEncoder(&body).Encode(ftPost)
					endpoint = serverURL + "/hub/add-ft-post"
				} else {
					// Paid tier users: use AddPostRequest (4096 char limit)
					paidPost := hub.AddPostRequest{
						Content: post.Content,
						TagIDs:  randomTagIDs,
					}
					err = json.NewEncoder(&body).Encode(paidPost)
					endpoint = serverURL + "/hub/add-post"
				}

				if err != nil {
					log.Fatalf("failed to encode post: %v", err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					endpoint,
					&body,
				)
				if err != nil {
					log.Fatalf("failed to create request: %v", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				tierStr := "paid"
				if isFreeTier {
					tierStr = "free"
				}
				color.Yellow(
					"thread %d (%s tier) in action for post %d",
					i,
					tierStr,
					j,
				)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalf("failed to send request: %v", err)
				}

				if resp.StatusCode != http.StatusOK {
					log.Fatalf(
						"failed to add post for %s tier user: %v",
						tierStr,
						resp.Status,
					)
				}

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
}
