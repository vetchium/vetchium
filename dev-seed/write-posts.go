package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/fatih/color"
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

			var posts []hub.AddPostRequest
			for j := 0; j < numPosts; j++ {
				post := postUniverse[rand.Intn(len(postUniverse))]
				posts = append(posts, hub.AddPostRequest{
					Content: post.Content,
					NewTags: post.NewTags,
				})
			}

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			for k, post := range posts {
				var body bytes.Buffer
				err := json.NewEncoder(&body).Encode(post)
				if err != nil {
					log.Fatalf("failed to encode post: %v", err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/add-post",
					&body,
				)
				if err != nil {
					log.Fatalf("failed to create request: %v", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				color.Yellow("thread %d in action for post %d", i, k)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalf("failed to send request: %v", err)
				}

				if resp.StatusCode != http.StatusOK {
					log.Fatalf("failed to add post: %v", resp.Status)
				}

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
}
