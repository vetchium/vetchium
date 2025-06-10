package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

// Sample incognito post content organized by categories
var incognitoPostContent = map[string][]struct {
	content string
	tags    []common.VTagID
}{
	"career-advice": {
		{
			content: "Has anyone successfully transitioned from a non-tech background to software engineering? Looking for advice on the best learning path.",
			tags:    []common.VTagID{"careers", "software-engineering"},
		},
		{
			content: "Feeling stuck in my current role. Management doesn't seem to value my contributions and I'm considering looking elsewhere. How do you know when it's time to leave?",
			tags:    []common.VTagID{"careers", "leadership"},
		},
		{
			content: "Is it worth getting an MBA later in your career? I'm 8 years into my professional journey and wondering if it would open new doors.",
			tags:    []common.VTagID{"careers", "entrepreneurship"},
		},
		{
			content: "Negotiating salary for the first time and I'm terrified. Any tips on how to approach this conversation with your manager?",
			tags:    []common.VTagID{"careers", "human-resources"},
		},
	},
	"workplace-issues": {
		{
			content: "Working with a colleague who takes credit for my ideas. HR seems unwilling to help. Anyone dealt with something similar?",
			tags:    []common.VTagID{"human-resources", "leadership"},
		},
		{
			content: "My company says they value work-life balance but expects 60+ hour weeks. The disconnect is getting to me. How do you handle this?",
			tags:    []common.VTagID{"leadership", "human-resources"},
		},
		{
			content: "Micromanaging boss is making my job impossible. Every decision needs approval and it's killing my productivity. Advice?",
			tags:    []common.VTagID{"leadership", "human-resources"},
		},
		{
			content: "Company culture has gotten toxic since the acquisition. Many good people are leaving. Should I stick it out or jump ship?",
			tags:    []common.VTagID{"leadership", "careers"},
		},
	},
	"tech-discussions": {
		{
			content: "What's your take on the current AI hype? Are we in another bubble or is this the real transformation everyone claims it is?",
			tags:    []common.VTagID{"technology", "artificial-intelligence"},
		},
		{
			content: "Remote work productivity tips? I'm struggling to stay focused working from home and my output has definitely decreased.",
			tags:    []common.VTagID{"technology", "productivity"},
		},
		{
			content: "Anyone else concerned about the direction of their tech stack? We're using legacy systems and management won't invest in modernization.",
			tags:    []common.VTagID{"technology", "software-engineering"},
		},
		{
			content: "Imposter syndrome is hitting hard lately. 5 years of experience but still feel like I don't know what I'm doing half the time.",
			tags:    []common.VTagID{"technology", "careers"},
		},
	},
	"industry-insights": {
		{
			content: "The job market feels incredibly unstable right now. Layoffs everywhere but also claims of talent shortages. What's really going on?",
			tags:    []common.VTagID{"careers", "entrepreneurship"},
		},
		{
			content: "Startup equity: is it still worth it? Hearing horror stories about valuations plummeting and equity becoming worthless.",
			tags:    []common.VTagID{"entrepreneurship", "finance"},
		},
		{
			content: "Why do companies still require office presence when productivity data shows remote work is just as effective?",
			tags:    []common.VTagID{"leadership", "human-resources"},
		},
		{
			content: "The interview process at most companies is broken. Why are we still doing whiteboard coding for senior positions?",
			tags:    []common.VTagID{"human-resources", "software-engineering"},
		},
	},
	"personal-growth": {
		{
			content: "Burnout is real and I'm experiencing it now. How do you recover while still meeting work expectations? Any strategies that worked?",
			tags:    []common.VTagID{"careers", "leadership"},
		},
		{
			content: "Finding it hard to build meaningful professional relationships when everything is remote. How do you network effectively now?",
			tags:    []common.VTagID{"careers", "human-resources"},
		},
		{
			content: "Considering a career pivot to consulting. Anyone made this transition? What should I expect in terms of lifestyle changes?",
			tags:    []common.VTagID{"careers", "entrepreneurship"},
		},
		{
			content: "Learning new skills while working full-time is exhausting. How do you maintain motivation for continuous learning?",
			tags:    []common.VTagID{"careers", "technology"},
		},
	},
}

// Sample comment content for incognito posts
var incognitoCommentContent = []string{
	"I've been through something similar. Happy to share my experience.",
	"This resonates with me so much. Thanks for sharing.",
	"Have you considered talking to a mentor about this?",
	"I disagree with this approach. Here's why...",
	"Great insights! This gives me a lot to think about.",
	"Been there. It gets better, trust me.",
	"This is exactly what I needed to hear today.",
	"I had a different experience, but I understand your perspective.",
	"Have you tried approaching it from a different angle?",
	"This is more common than you think. You're not alone.",
	"Thanks for the honest post. More people need to talk about this.",
	"I learned this the hard way too. Wish someone had told me earlier.",
	"Solid advice. I've seen this work in practice.",
	"Interesting point. I hadn't thought about it that way.",
	"This is why I love this community. Real talk about real issues.",
	"Going through this right now. Would love to connect if possible.",
	"This post made my day. Thank you for sharing your story.",
	"I have some resources that might help. How can I share them?",
	"Your experience mirrors mine almost exactly. Crazy!",
	"This is the kind of discussion we need more of.",
}

func writeIncognitoPosts() {
	var wg sync.WaitGroup
	var createdPosts []string // Store post IDs for commenting later
	var postMutex sync.Mutex

	// First pass: Create incognito posts
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%3 == 0 {
			color.Cyan("incognito posts thread %d waiting", i)
			wg.Wait()
			color.Cyan("incognito posts thread %d resumes", i)
		}
		wg.Add(1)

		go func(user HubSeedUser, i int) {
			defer wg.Done()

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			// Create 1-4 incognito posts per user
			numPosts := 1 + rand.Intn(4)
			categories := []string{
				"career-advice",
				"workplace-issues",
				"tech-discussions",
				"industry-insights",
				"personal-growth",
			}

			for j := 0; j < numPosts; j++ {
				category := categories[rand.Intn(len(categories))]
				categoryPosts := incognitoPostContent[category]
				selectedPost := categoryPosts[rand.Intn(len(categoryPosts))]

				postRequest := hub.AddIncognitoPostRequest{
					Content: selectedPost.content,
					TagIDs:  selectedPost.tags,
				}

				var body bytes.Buffer
				err := json.NewEncoder(&body).Encode(postRequest)
				if err != nil {
					log.Fatalf("failed to encode incognito post: %v", err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/add-incognito-post",
					&body,
				)
				if err != nil {
					log.Fatalf(
						"failed to create incognito post request: %v",
						err,
					)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				color.Cyan("thread %d creating incognito post %d", i, j)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalf("failed to send incognito post request: %v", err)
				}

				if resp.StatusCode != http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Fatalf("failed to read response body: %v", err)
					}
					log.Fatalf("failed to add incognito post: %v", string(body))
				}

				var response hub.AddIncognitoPostResponse
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					log.Fatalf(
						"failed to decode incognito post response: %v",
						err,
					)
				}

				// Store post ID for later commenting
				postMutex.Lock()
				createdPosts = append(createdPosts, response.IncognitoPostID)
				postMutex.Unlock()

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Created %d incognito posts", len(createdPosts))

	// Second pass: Add comments to posts
	writeIncognitoComments(createdPosts)

	// Third pass: Add votes to posts and comments
	writeIncognitoVotes(createdPosts)
}

func writeIncognitoComments(postIDs []string) {
	var wg sync.WaitGroup
	var createdComments []struct {
		postID    string
		commentID string
	}
	var commentMutex sync.Mutex

	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%4 == 0 {
			color.Magenta("incognito comments thread %d waiting", i)
			wg.Wait()
			color.Magenta("incognito comments thread %d resumes", i)
		}
		wg.Add(1)

		go func(user HubSeedUser, i int) {
			defer wg.Done()

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			// Each user comments on 2-5 random posts
			numComments := 2 + rand.Intn(4)
			for j := 0; j < numComments; j++ {
				if len(postIDs) == 0 {
					continue
				}

				postID := postIDs[rand.Intn(len(postIDs))]
				content := incognitoCommentContent[rand.Intn(len(incognitoCommentContent))]

				commentRequest := hub.AddIncognitoPostCommentRequest{
					IncognitoPostID: postID,
					Content:         content,
				}

				var body bytes.Buffer
				err := json.NewEncoder(&body).Encode(commentRequest)
				if err != nil {
					log.Fatalf("failed to encode incognito comment: %v", err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/add-incognito-post-comment",
					&body,
				)
				if err != nil {
					log.Fatalf(
						"failed to create incognito comment request: %v",
						err,
					)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				color.Magenta("thread %d adding comment %d", i, j)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalf(
						"failed to send incognito comment request: %v",
						err,
					)
				}

				if resp.StatusCode != http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						color.Red(
							"failed to read comment error response: %v",
							err,
						)
						continue
					}
					color.Red(
						"failed to add incognito comment: %v",
						string(body),
					)
					resp.Body.Close()
					continue
				}

				var response hub.AddIncognitoPostCommentResponse
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					log.Fatalf(
						"failed to decode incognito comment response: %v",
						err,
					)
				}

				commentMutex.Lock()
				createdComments = append(createdComments, struct {
					postID    string
					commentID string
				}{
					postID:    response.IncognitoPostID,
					commentID: response.CommentID,
				})
				commentMutex.Unlock()

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Created %d incognito comments", len(createdComments))

	// Add some nested replies
	writeIncognitoReplies(createdComments)
}

func writeIncognitoReplies(comments []struct {
	postID    string
	commentID string
}) {
	if len(comments) == 0 {
		return
	}

	var wg sync.WaitGroup

	// Create some replies to existing comments
	replyCount := len(comments) / 3 // About 1/3 of comments get replies
	if replyCount > 20 {
		replyCount = 20 // Cap at 20 replies
	}

	for i := 0; i < replyCount; i++ {
		user := hubUsers[rand.Intn(len(hubUsers))]
		wg.Add(1)

		go func(user HubSeedUser, i int) {
			defer wg.Done()

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Printf("no auth token found for %s", user.Email)
				return
			}
			authToken := tokenI.(string)

			comment := comments[rand.Intn(len(comments))]
			content := incognitoCommentContent[rand.Intn(len(incognitoCommentContent))]

			replyRequest := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: comment.postID,
				Content:         content,
				InReplyTo:       &comment.commentID,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(replyRequest)
			if err != nil {
				log.Printf("failed to encode incognito reply: %v", err)
				return
			}

			req, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/hub/add-incognito-post-comment",
				&body,
			)
			if err != nil {
				log.Printf("failed to create incognito reply request: %v", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+authToken)

			color.Blue("adding reply %d", i)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("failed to send incognito reply request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					color.Red("failed to read reply error response: %v", err)
					return
				}
				color.Red("failed to add incognito reply: %v", string(body))
				return
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Created %d incognito replies", replyCount)
}

func writeIncognitoVotes(postIDs []string) {
	if len(postIDs) == 0 {
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%5 == 0 {
			color.Yellow("incognito votes thread %d waiting", i)
			wg.Wait()
			color.Yellow("incognito votes thread %d resumes", i)
		}
		wg.Add(1)

		go func(user HubSeedUser, i int) {
			defer wg.Done()

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			// Each user votes on 3-8 random posts
			numVotes := 3 + rand.Intn(6)
			for j := 0; j < numVotes; j++ {
				postID := postIDs[rand.Intn(len(postIDs))]

				// 70% upvotes, 30% downvotes
				isUpvote := rand.Float32() < 0.7

				var endpoint string
				var voteRequest interface{}

				if isUpvote {
					endpoint = "/hub/upvote-incognito-post"
					voteRequest = hub.UpvoteIncognitoPostRequest{
						IncognitoPostID: postID,
					}
				} else {
					endpoint = "/hub/downvote-incognito-post"
					voteRequest = hub.DownvoteIncognitoPostRequest{
						IncognitoPostID: postID,
					}
				}

				var body bytes.Buffer
				err := json.NewEncoder(&body).Encode(voteRequest)
				if err != nil {
					log.Printf("failed to encode vote request: %v", err)
					continue
				}

				req, err := http.NewRequest(
					http.MethodPost,
					serverURL+endpoint,
					&body,
				)
				if err != nil {
					log.Printf("failed to create vote request: %v", err)
					continue
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				voteType := "upvote"
				if !isUpvote {
					voteType = "downvote"
				}
				color.Yellow("thread %d adding %s %d", i, voteType, j)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("failed to send vote request: %v", err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						color.Red("failed to read vote error response: %v", err)
					} else {
						color.Red("failed to vote on post: %v", string(body))
					}
				}

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Completed incognito post voting")
}
