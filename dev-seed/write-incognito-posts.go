package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

// Lorem ipsum phrases for generating varied reply content
var loremIpsumPhrases = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse.",
	"Excepteur sint occaecat cupidatat non proident, sunt in culpa.",
	"Nulla pariatur sed ut perspiciatis unde omnis iste natus error.",
	"Accusantium doloremque laudantium, totam rem aperiam eaque ipsa.",
	"Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet.",
	"Consectetur, adipisci velit, sed quia non numquam eius modi.",
	"Tempora incidunt ut labore et dolore magnam aliquam quaerat.",
	"Voluptatem accusantium doloremque laudantium totam rem aperiam.",
	"Eaque ipsa quae ab illo inventore veritatis et quasi architecto.",
	"Beatae vitae dicta sunt explicabo nemo enim ipsam voluptatem.",
	"Quia voluptas sit aspernatur aut odit aut fugit sed quia.",
	"Consequuntur magni dolores eos qui ratione voluptatem sequi.",
	"At vero eos et accusamus et iusto odio dignissimos ducimus.",
	"Qui blanditiis praesentium voluptatum deleniti atque corrupti.",
	"Quos dolores et quas molestias excepturi sint occaecati.",
	"Cupiditate non provident similique sunt in culpa qui officia.",
	"Deserunt mollitia animi id est laborum et dolorum fuga.",
}

// Reply templates for more natural-sounding responses
var replyTemplates = []string{
	"Great point! %s",
	"I think %s",
	"That's interesting. %s",
	"I disagree though. %s",
	"Adding to this: %s",
	"From my experience, %s",
	"Actually, %s",
	"This reminds me: %s",
	"I've seen this too. %s",
	"Totally agree! %s",
	"Not sure about that. %s",
	"Following up on this: %s",
	"Building on your point: %s",
	"Counter-argument: %s",
	"Related thought: %s",
}

// Generate varied lorem ipsum content for replies
func generateLoremReply() string {
	template := replyTemplates[rand.Intn(len(replyTemplates))]

	// Combine 1-3 lorem phrases
	numPhrases := 1 + rand.Intn(3)
	var phrases []string
	for i := 0; i < numPhrases; i++ {
		phrases = append(
			phrases,
			loremIpsumPhrases[rand.Intn(len(loremIpsumPhrases))],
		)
	}

	loremContent := ""
	for i, phrase := range phrases {
		if i > 0 {
			loremContent += " "
		}
		loremContent += phrase
	}

	// Sometimes just return the lorem content directly
	if rand.Float32() < 0.3 {
		return loremContent
	}

	// Otherwise use template
	return fmt.Sprintf(template, loremContent)
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

	// Second pass: Add comments to posts with enhanced volume
	writeIncognitoCommentsEnhanced(createdPosts)

	// Third pass: Add votes to posts and comments
	writeIncognitoVotes(createdPosts)
}

func writeIncognitoCommentsEnhanced(postIDs []string) {
	var wg sync.WaitGroup
	var allComments []struct {
		postID    string
		commentID string
		depth     int
	}
	var commentMutex sync.Mutex

	// Phase 1: Initial comments (significantly increased volume)
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]

		if i%4 == 0 {
			color.Magenta("incognito comments phase 1 thread %d waiting", i)
			wg.Wait()
			color.Magenta("incognito comments phase 1 thread %d resumes", i)
		}
		wg.Add(1)

		go func(user HubSeedUser, i int) {
			defer wg.Done()

			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			// Increased from 2-5 to 6-12 comments per user for 3x volume
			numComments := 6 + rand.Intn(7)
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
				allComments = append(allComments, struct {
					postID    string
					commentID string
					depth     int
				}{
					postID:    response.IncognitoPostID,
					commentID: response.CommentID,
					depth:     0, // Initial comments are depth 0
				})
				commentMutex.Unlock()

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Created %d initial incognito comments", len(allComments))

	// Phase 2: Multi-level nested replies
	writeMultiLevelReplies(allComments)
}

func writeMultiLevelReplies(initialComments []struct {
	postID    string
	commentID string
	depth     int
}) {
	if len(initialComments) == 0 {
		return
	}

	var allComments []struct {
		postID    string
		commentID string
		depth     int
	}

	// Start with initial comments
	allComments = append(allComments, initialComments...)

	var commentMutex sync.Mutex
	maxDepth := 4 // Allow up to 4 levels of nesting

	// Generate replies for each depth level
	for currentDepth := 0; currentDepth < maxDepth; currentDepth++ {
		var wg sync.WaitGroup

		// Get comments at current depth to reply to
		var commentsAtDepth []struct {
			postID    string
			commentID string
			depth     int
		}

		for _, comment := range allComments {
			if comment.depth == currentDepth {
				commentsAtDepth = append(commentsAtDepth, comment)
			}
		}

		if len(commentsAtDepth) == 0 {
			continue
		}

		// Calculate reply ratio - more replies for shallower comments
		var replyRatio float32
		switch currentDepth {
		case 0:
			replyRatio = 0.7 // 70% of top-level comments get replies
		case 1:
			replyRatio = 0.5 // 50% of depth-1 comments get replies
		case 2:
			replyRatio = 0.3 // 30% of depth-2 comments get replies
		case 3:
			replyRatio = 0.15 // 15% of depth-3 comments get replies
		}

		targetReplies := int(float32(len(commentsAtDepth)) * replyRatio)
		color.Cyan(
			"Creating %d replies at depth %d",
			targetReplies,
			currentDepth+1,
		)

		// Create replies in batches to avoid overwhelming the server
		batchSize := 20
		for batch := 0; batch < targetReplies; batch += batchSize {
			batchEnd := batch + batchSize
			if batchEnd > targetReplies {
				batchEnd = targetReplies
			}

			for i := batch; i < batchEnd; i++ {
				wg.Add(1)

				go func(replyIndex int) {
					defer wg.Done()

					// Pick a random user
					user := hubUsers[rand.Intn(len(hubUsers))]

					tokenI, ok := hubSessionTokens.Load(user.Email)
					if !ok {
						log.Printf("no auth token found for %s", user.Email)
						return
					}
					authToken := tokenI.(string)

					// Pick a random comment at current depth
					parentComment := commentsAtDepth[rand.Intn(len(commentsAtDepth))]

					// Generate lorem ipsum content for replies
					content := generateLoremReply()

					replyRequest := hub.AddIncognitoPostCommentRequest{
						IncognitoPostID: parentComment.postID,
						Content:         content,
						InReplyTo:       &parentComment.commentID,
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
						log.Printf(
							"failed to create incognito reply request: %v",
							err,
						)
						return
					}

					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authToken)

					color.Blue(
						"adding reply at depth %d (index %d)",
						currentDepth+1,
						replyIndex,
					)

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Printf(
							"failed to send incognito reply request: %v",
							err,
						)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							color.Red(
								"failed to read reply error response: %v",
								err,
							)
							return
						}
						color.Red(
							"failed to add incognito reply: %v",
							string(body),
						)
						return
					}

					var response hub.AddIncognitoPostCommentResponse
					err = json.NewDecoder(resp.Body).Decode(&response)
					if err != nil {
						log.Printf(
							"failed to decode incognito reply response: %v",
							err,
						)
						return
					}

					// Add the new reply to our tracking
					commentMutex.Lock()
					allComments = append(allComments, struct {
						postID    string
						commentID string
						depth     int
					}{
						postID:    response.IncognitoPostID,
						commentID: response.CommentID,
						depth:     currentDepth + 1,
					})
					commentMutex.Unlock()
				}(i)
			}

			// Wait for this batch to complete before starting the next
			wg.Wait()
		}

		color.Green("Completed depth %d replies", currentDepth+1)
	}

	// Count total comments by depth
	depthCounts := make(map[int]int)
	for _, comment := range allComments {
		depthCounts[comment.depth]++
	}

	color.Green("Final comment distribution:")
	totalComments := 0
	for depth := 0; depth <= maxDepth; depth++ {
		if count, exists := depthCounts[depth]; exists {
			color.Green("  Depth %d: %d comments", depth, count)
			totalComments += count
		}
	}
	color.Green("Total comments created: %d", totalComments)
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

				// 422 will be returned on self-post voting
				if resp.StatusCode != http.StatusOK &&
					resp.StatusCode != http.StatusUnprocessableEntity {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						color.Red("vote error: %v %v", err, resp.StatusCode)
					} else {
						color.Red("error: %v %v", string(body), resp.StatusCode)
						log.Fatalf("failed to vote on post: %v %v", string(body), resp.StatusCode)
					}
				}

				resp.Body.Close()
			}
		}(user, i)
	}

	wg.Wait()
	color.Green("Completed incognito post voting")
}

// createMegaIncognitoThread creates one massive incognito post with ~5000 comments
// across multiple nested levels with time-staggered top-level comments
func createMegaIncognitoThread() {
	color.Cyan("Creating mega incognito thread with ~5000 comments...")

	// Step 1: Create the main mega post using user0 as the author
	if len(hubUsers) == 0 {
		log.Fatalf("no hub users available for mega thread creation")
	}
	user := hubUsers[0] // Always use user0 as the mega thread author
	tokenI, ok := hubSessionTokens.Load(user.Email)
	if !ok {
		log.Fatalf("no auth token found for user0 (%s)", user.Email)
	}
	authToken := tokenI.(string)

	megaPostContent := `The future of remote work and career trajectories - let's discuss!

Remote work has fundamentally changed how we build careers and navigate professional development. This shift seems permanent, but what are the long-term implications?

Key areas to explore:
• Career Development: Building mentorship relationships remotely vs. in-person
• Company Culture: Can authentic culture exist in fully remote environments?
• Skills: What becomes more important in a remote-first world?
• Economics: Impact on salary negotiations and geographic pay equity
• Future: Permanent hybrid model or return to office mandates?

Whether you're early career or experienced, your perspective would be valuable. What has remote work taught you about your professional goals and development? How has it changed your approach to career growth?

Looking forward to hearing diverse experiences and predictions about where we're headed!`

	postRequest := hub.AddIncognitoPostRequest{
		Content: megaPostContent,
		TagIDs: []common.VTagID{
			"careers",
			"remote-work",
			"leadership",
		},
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(postRequest)
	if err != nil {
		log.Fatalf("failed to encode mega post: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/add-incognito-post",
		&body,
	)
	if err != nil {
		log.Fatalf("failed to create mega post request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	color.Cyan("Creating mega incognito post with user0 (%s)...", user.Email)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send mega post request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}
		log.Fatalf("failed to add mega post: %v", string(body))
	}

	var response hub.AddIncognitoPostResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("failed to decode mega post response: %v", err)
	}
	resp.Body.Close()

	megaPostID := response.IncognitoPostID
	color.Green("Created mega post with ID: %s", megaPostID)

	// Step 2: Create top-level comments with time staggering
	var topLevelComments []struct {
		postID    string
		commentID string
		depth     int
	}
	var commentMutex sync.Mutex

	// Create 25 top-level comments to keep total comment count reasonable (~50k instead of 200k)
	numTopLevelComments := 25
	color.Cyan(
		"Creating %d top-level comments with time staggering...",
		numTopLevelComments,
	)

	topLevelContents := []string{
		"This is such an important discussion. Remote work has completely changed my career trajectory.",
		"I've been working remotely for 5 years and the learning curve was steep initially.",
		"The mentorship aspect is crucial. I've found it much harder to build those relationships remotely.",
		"Company culture remotely feels different, but not necessarily worse in my experience.",
		"I think we're still in the experimental phase. The real impact won't be clear for another decade.",
		"Geographic pay equity is fascinating. I'm earning Silicon Valley wages while living in a small town.",
		"The skills that matter most now are written communication and self-motivation.",
		"I disagree with the premise. Office work had its own set of problems we're conveniently forgetting.",
		"Hybrid seems like the worst of both worlds to me. Either commit to remote or don't.",
		"The economic implications are huge. Entire industries are being disrupted.",
		"As someone early in their career, I worry about missing out on informal learning opportunities.",
		"Remote work has made me more intentional about professional development.",
		"The future is definitely hybrid, but the balance will vary by industry and role.",
		"I've seen both successful and failed attempts at remote culture building.",
		"The productivity gains are real, but so are the collaboration challenges.",
	}

	for i := 0; i < numTopLevelComments; i++ {
		// Use a different user for each comment
		user := hubUsers[i%len(hubUsers)]
		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Printf("no auth token found for %s", user.Email)
			continue
		}
		authToken := tokenI.(string)

		// Pick content (cycling through available ones and generating lorem for extras)
		var content string
		if i < len(topLevelContents) {
			content = topLevelContents[i]
		} else {
			content = generateLoremReply()
		}

		commentRequest := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: megaPostID,
			Content:         content,
		}

		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(commentRequest)
		if err != nil {
			log.Printf("failed to encode top-level comment %d: %v", i, err)
			continue
		}

		req, err := http.NewRequest(
			http.MethodPost,
			serverURL+"/hub/add-incognito-post-comment",
			&body,
		)
		if err != nil {
			log.Printf(
				"failed to create top-level comment request %d: %v",
				i,
				err,
			)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		color.Blue("Creating top-level comment %d/%d", i+1, numTopLevelComments)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf(
				"failed to send top-level comment request %d: %v",
				i,
				err,
			)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				color.Red("failed to read comment error response: %v", err)
			} else {
				color.Red("failed to add top-level comment: %v", string(body))
			}
			resp.Body.Close()
			continue
		}

		var response hub.AddIncognitoPostCommentResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Printf(
				"failed to decode top-level comment response %d: %v",
				i,
				err,
			)
			resp.Body.Close()
			continue
		}

		commentMutex.Lock()
		topLevelComments = append(topLevelComments, struct {
			postID    string
			commentID string
			depth     int
		}{
			postID:    response.IncognitoPostID,
			commentID: response.CommentID,
			depth:     0,
		})
		commentMutex.Unlock()

		resp.Body.Close()

		// Add time staggering - 1-3 seconds between some comments
		if i%10 == 0 && i > 0 {
			sleepTime := 1 + rand.Intn(3)
			color.Yellow(
				"Staggering comment creation - sleeping %d seconds",
				sleepTime,
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

		// Smaller delays between all comments to avoid overwhelming the server
		time.Sleep(100 * time.Millisecond)
	}

	color.Green("Created %d top-level comments", len(topLevelComments))

	// Step 3: Create massive nested thread structure and collect all comments
	allMegaComments := createMegaThreadReplies(topLevelComments, megaPostID)

	// Step 4: Add realistic voting patterns to the mega thread
	color.Cyan("Adding voting patterns to mega thread...")
	voteMegaThreadComments(allMegaComments)
}

// createMegaThreadReplies creates nested replies up to maximum depth and returns all comments
func createMegaThreadReplies(initialComments []struct {
	postID    string
	commentID string
	depth     int
}, postID string) []struct {
	postID    string
	commentID string
	depth     int
} {
	if len(initialComments) == 0 {
		return nil
	}

	var allComments []struct {
		postID    string
		commentID string
		depth     int
	}

	// Start with initial comments
	allComments = append(allComments, initialComments...)
	var commentMutex sync.Mutex

	maxDepth := 5 // Maximum depth is 5 levels (0-5) as per API specification

	color.Cyan(
		"Creating nested mega thread to maximum depth %d (natural comment distribution)",
		maxDepth,
	)

	// Generate replies for each depth level
	for currentDepth := 0; currentDepth < maxDepth; currentDepth++ {
		// Get comments at current depth to reply to
		var commentsAtDepth []struct {
			postID    string
			commentID string
			depth     int
		}

		commentMutex.Lock()
		for _, comment := range allComments {
			if comment.depth == currentDepth {
				commentsAtDepth = append(commentsAtDepth, comment)
			}
		}
		currentTotal := len(allComments)
		commentMutex.Unlock()

		if len(commentsAtDepth) == 0 {
			color.Yellow("No comments at depth %d, stopping", currentDepth)
			break
		}

		color.Blue(
			"Currently at depth %d with %d total comments",
			currentDepth,
			currentTotal,
		)

		// Calculate reply distribution for this depth (optimized for ~50k total comments)
		var repliesPerComment int
		switch currentDepth {
		case 0:
			repliesPerComment = 6 // Each top-level comment gets 6 replies
		case 1:
			repliesPerComment = 4 // Each depth-1 comment gets 4 replies
		case 2:
			repliesPerComment = 3 // Each depth-2 comment gets 3 replies
		case 3:
			repliesPerComment = 2 // Each depth-3 comment gets 2 replies
		case 4:
			repliesPerComment = 1 // Each depth-4 comment gets 1 reply
		}

		targetReplies := len(commentsAtDepth) * repliesPerComment

		color.Cyan(
			"Creating %d replies at depth %d (from %d parent comments)",
			targetReplies,
			currentDepth+1,
			len(commentsAtDepth),
		)

		// Create replies in batches to control server load
		batchSize := 25
		createdReplies := 0

		for batch := 0; batch < targetReplies; batch += batchSize {
			batchEnd := batch + batchSize
			if batchEnd > targetReplies {
				batchEnd = targetReplies
			}

			batchWg := sync.WaitGroup{}
			for i := batch; i < batchEnd; i++ {
				batchWg.Add(1)

				go func(replyIndex int) {
					defer batchWg.Done()

					// Pick a random user
					user := hubUsers[rand.Intn(len(hubUsers))]
					tokenI, ok := hubSessionTokens.Load(user.Email)
					if !ok {
						log.Printf("no auth token found for %s", user.Email)
						return
					}
					authToken := tokenI.(string)

					// Pick a random parent comment at current depth
					parentComment := commentsAtDepth[rand.Intn(len(commentsAtDepth))]

					// Generate content for this reply
					content := generateLoremReply()

					replyRequest := hub.AddIncognitoPostCommentRequest{
						IncognitoPostID: parentComment.postID,
						Content:         content,
						InReplyTo:       &parentComment.commentID,
					}

					var body bytes.Buffer
					err := json.NewEncoder(&body).Encode(replyRequest)
					if err != nil {
						log.Printf("failed to encode mega reply: %v", err)
						return
					}

					req, err := http.NewRequest(
						http.MethodPost,
						serverURL+"/hub/add-incognito-post-comment",
						&body,
					)
					if err != nil {
						log.Printf(
							"failed to create mega reply request: %v",
							err,
						)
						return
					}

					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authToken)

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Printf("failed to send mega reply request: %v", err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							color.Red(
								"failed to read mega reply error response: %v",
								err,
							)
						} else {
							color.Red("failed to add mega reply: %v", string(body))
						}
						return
					}

					var response hub.AddIncognitoPostCommentResponse
					err = json.NewDecoder(resp.Body).Decode(&response)
					if err != nil {
						log.Printf(
							"failed to decode mega reply response: %v",
							err,
						)
						return
					}

					// Add the new reply to our tracking
					commentMutex.Lock()
					allComments = append(allComments, struct {
						postID    string
						commentID string
						depth     int
					}{
						postID:    response.IncognitoPostID,
						commentID: response.CommentID,
						depth:     currentDepth + 1,
					})
					commentMutex.Unlock()
				}(i)
			}

			// Wait for this batch to complete
			batchWg.Wait()
			createdReplies += (batchEnd - batch)

			// Progress update
			color.Blue(
				"Completed batch: %d/%d replies at depth %d",
				createdReplies,
				targetReplies,
				currentDepth+1,
			)

			// Small delay between batches
			time.Sleep(200 * time.Millisecond)
		}

		commentMutex.Lock()
		totalSoFar := len(allComments)
		commentMutex.Unlock()

		color.Green(
			"Completed depth %d: created %d replies (total comments: %d)",
			currentDepth+1,
			createdReplies,
			totalSoFar,
		)

		// Continue to next depth level
	}

	// Final statistics
	depthCounts := make(map[int]int)
	commentMutex.Lock()
	for _, comment := range allComments {
		depthCounts[comment.depth]++
	}
	totalComments := len(allComments)
	commentMutex.Unlock()

	color.Green("Mega thread completed! Final statistics:")
	for depth := 0; depth <= maxDepth; depth++ {
		if count, exists := depthCounts[depth]; exists {
			color.Green("  Depth %d: %d comments", depth, count)
		}
	}
	color.Green("Total comments in mega thread: %d", totalComments)

	return allComments
}

// voteMegaThreadComments adds realistic voting patterns to the mega thread comments
func voteMegaThreadComments(allComments []struct {
	postID    string
	commentID string
	depth     int
}) {
	if len(allComments) == 0 {
		color.Yellow("No comments to vote on")
		return
	}

	color.Cyan(
		"Adding realistic voting patterns to %d mega thread comments...",
		len(allComments),
	)

	var wg sync.WaitGroup

	// Create voting patterns:
	// - 20% get significant upvotes (5-20 upvotes)
	// - 5% get significant downvotes (5-15 downvotes)
	// - 15% get mixed moderate voting (1-4 votes either way)
	// - 60% get no votes (realistic - most comments don't get voted on)

	popularComments := allComments[:int(float64(len(allComments))*0.20)]                                          // Top 20% get lots of upvotes
	controversialComments := allComments[int(float64(len(allComments))*0.20):int(float64(len(allComments))*0.25)] // Next 5% get downvotes
	moderateComments := allComments[int(float64(len(allComments))*0.25):int(float64(len(allComments))*0.40)]      // Next 15% get moderate voting
	// Remaining 60% get no votes

	// Vote on popular comments (lots of upvotes)
	for i, comment := range popularComments {
		if i%10 == 0 {
			color.Green("Voting on popular comments batch %d", i/10)
			wg.Wait()
		}

		// Each popular comment gets 5-20 upvotes from different users
		numUpvotes := 5 + rand.Intn(16)
		for j := 0; j < numUpvotes && j < len(hubUsers); j++ {
			wg.Add(1)
			go func(commentData struct {
				postID    string
				commentID string
				depth     int
			}, userIndex int) {
				defer wg.Done()

				user := hubUsers[userIndex]
				tokenI, ok := hubSessionTokens.Load(user.Email)
				if !ok {
					return
				}
				authToken := tokenI.(string)

				voteRequest := hub.UpvoteIncognitoPostCommentRequest{
					IncognitoPostID: commentData.postID,
					CommentID:       commentData.commentID,
				}

				var body bytes.Buffer
				json.NewEncoder(&body).Encode(voteRequest)

				req, _ := http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/upvote-incognito-post-comment",
					&body,
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					resp.Body.Close()
				}
			}(comment, (i*7+j)%len(hubUsers)) // Spread across different users
		}
	}

	wg.Wait()
	color.Green("Completed upvoting popular comments")

	// Vote on controversial comments (lots of downvotes)
	for i, comment := range controversialComments {
		if i%5 == 0 {
			color.Red("Voting on controversial comments batch %d", i/5)
			wg.Wait()
		}

		// Each controversial comment gets 5-15 downvotes
		numDownvotes := 5 + rand.Intn(11)
		for j := 0; j < numDownvotes && j < len(hubUsers); j++ {
			wg.Add(1)
			go func(commentData struct {
				postID    string
				commentID string
				depth     int
			}, userIndex int) {
				defer wg.Done()

				user := hubUsers[userIndex]
				tokenI, ok := hubSessionTokens.Load(user.Email)
				if !ok {
					return
				}
				authToken := tokenI.(string)

				voteRequest := hub.DownvoteIncognitoPostCommentRequest{
					IncognitoPostID: commentData.postID,
					CommentID:       commentData.commentID,
				}

				var body bytes.Buffer
				json.NewEncoder(&body).Encode(voteRequest)

				req, _ := http.NewRequest(
					http.MethodPost,
					serverURL+"/hub/downvote-incognito-post-comment",
					&body,
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					resp.Body.Close()
				}
			}(comment, (i*11+j)%len(hubUsers)) // Different user pattern
		}
	}

	wg.Wait()
	color.Green("Completed downvoting controversial comments")

	// Vote on moderate comments (mixed low-level voting)
	for i, comment := range moderateComments {
		if i%15 == 0 {
			color.Blue("Voting on moderate comments batch %d", i/15)
			wg.Wait()
		}

		// Each moderate comment gets 1-4 votes (mix of up/down)
		numVotes := 1 + rand.Intn(4)
		for j := 0; j < numVotes && j < len(hubUsers); j++ {
			wg.Add(1)
			go func(commentData struct {
				postID    string
				commentID string
				depth     int
			}, userIndex int, isUpvote bool) {
				defer wg.Done()

				user := hubUsers[userIndex]
				tokenI, ok := hubSessionTokens.Load(user.Email)
				if !ok {
					return
				}
				authToken := tokenI.(string)

				var endpoint string
				var voteRequest interface{}

				if isUpvote {
					endpoint = "/hub/upvote-incognito-post-comment"
					voteRequest = hub.UpvoteIncognitoPostCommentRequest{
						IncognitoPostID: commentData.postID,
						CommentID:       commentData.commentID,
					}
				} else {
					endpoint = "/hub/downvote-incognito-post-comment"
					voteRequest = hub.DownvoteIncognitoPostCommentRequest{
						IncognitoPostID: commentData.postID,
						CommentID:       commentData.commentID,
					}
				}

				var body bytes.Buffer
				json.NewEncoder(&body).Encode(voteRequest)

				req, _ := http.NewRequest(
					http.MethodPost,
					serverURL+endpoint,
					&body,
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+authToken)

				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					resp.Body.Close()
				}
			}(comment, (i*13+j)%len(hubUsers), rand.Float32() < 0.6) // 60% upvotes, 40% downvotes
		}
	}

	wg.Wait()
	color.Green("Completed voting on moderate comments")

	// Final statistics
	color.Green("Mega thread voting completed:")
	color.Green("  Popular comments (high upvotes): %d", len(popularComments))
	color.Green(
		"  Controversial comments (downvotes): %d",
		len(controversialComments),
	)
	color.Green("  Moderate comments (mixed votes): %d", len(moderateComments))
	color.Green(
		"  Unvoted comments (realistic): %d",
		len(
			allComments,
		)-len(
			popularComments,
		)-len(
			controversialComments,
		)-len(
			moderateComments,
		),
	)
}
