package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/typespec/hub"
)

func writePosts() {
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%10 == 0 {
			// Just to reduce the parallelism for a while
			// so we do not run out of sockets
			<-time.After(3 * time.Second)
		}

		// Parallelism to have more messages from many authors at similar times
		go func(user HubSeedUser, i int) {
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

			for _, post := range posts {
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

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalf("failed to send request: %v", err)
				}

				if resp.StatusCode != http.StatusOK {
					log.Fatalf("failed to add post: %v", resp.Status)
				}
			}
		}(user, i)
	}
}
