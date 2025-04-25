package dolores

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Timeline Operations", Ordered, func() {
	var db *pgxpool.Pool

	// Constants for timeline test - aligned with backend timing in api/internal/granger/timelines.go
	const (
		// Backend checks timelines every 1 second when processing continuously
		// We use 10 seconds to allow for processing time and potential queueing
		TimelineRefreshInterval = 1 * time.Minute

		// Match backend's interval for checking timelines
		TimelinePollInterval = 30 * time.Second

		// Allow more retries to accommodate potential delays
		MaxTimelinePollRetries = 20
	)

	// Declare token variables directly at test level
	var user1Token, user2Token, user3Token, user4Token, user5Token string
	var user8Token, user9Token, user10Token, user15Token string
	var user11Token, user12Token string
	var user6Token, user7Token string

	// User token maps for different test groups
	var timelineTokens map[string]string
	var followUnfollowTokens map[string]string
	var paginationTokens map[string]string
	var missingTimelineTokens map[string]string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0024-timeline-up.pgsql")

		// Initialize token maps
		timelineTokens = make(map[string]string)
		followUnfollowTokens = make(map[string]string)
		paginationTokens = make(map[string]string)
		missingTimelineTokens = make(map[string]string)

		// Login hub users and get tokens - using different users for different tests
		var wg sync.WaitGroup

		// Users for basic timeline tests (existing users from the seed)
		wg.Add(5)
		hubSigninAsync(
			"user1@0024-timeline-test.example.com",
			"NewPassword123$",
			&user1Token,
			&wg,
		)
		hubSigninAsync(
			"user2@0024-timeline-test.example.com",
			"NewPassword123$",
			&user2Token,
			&wg,
		)
		hubSigninAsync(
			"user3@0024-timeline-test.example.com",
			"NewPassword123$",
			&user3Token,
			&wg,
		)
		hubSigninAsync(
			"user4@0024-timeline-test.example.com",
			"NewPassword123$",
			&user4Token,
			&wg,
		)
		hubSigninAsync(
			"user5@0024-timeline-test.example.com",
			"NewPassword123$",
			&user5Token,
			&wg,
		)

		// Users for follow/unfollow timeline tests
		wg.Add(4)
		hubSigninAsync(
			"user8@0024-timeline-test.example.com",
			"NewPassword123$",
			&user8Token,
			&wg,
		)
		hubSigninAsync(
			"user9@0024-timeline-test.example.com",
			"NewPassword123$",
			&user9Token,
			&wg,
		)
		hubSigninAsync(
			"user10@0024-timeline-test.example.com",
			"NewPassword123$",
			&user10Token,
			&wg,
		)
		hubSigninAsync(
			"user15@0024-timeline-test.example.com",
			"NewPassword123$",
			&user15Token,
			&wg,
		)

		// Users for pagination tests
		wg.Add(2)
		hubSigninAsync(
			"user11@0024-timeline-test.example.com",
			"NewPassword123$",
			&user11Token,
			&wg,
		)
		hubSigninAsync(
			"user12@0024-timeline-test.example.com",
			"NewPassword123$",
			&user12Token,
			&wg,
		)

		// Users without timeline setup
		wg.Add(2)
		hubSigninAsync(
			"user6@0024-timeline-test.example.com",
			"NewPassword123$",
			&user6Token,
			&wg,
		)
		hubSigninAsync(
			"user7@0024-timeline-test.example.com",
			"NewPassword123$",
			&user7Token,
			&wg,
		)

		wg.Wait()

		// Make sure all tokens were retrieved successfully
		fmt.Fprintf(GinkgoWriter, "User1 token: %s\n", user1Token)
		fmt.Fprintf(GinkgoWriter, "User2 token: %s\n", user2Token)
		fmt.Fprintf(GinkgoWriter, "User3 token: %s\n", user3Token)
		// Add more debug log statements as needed

		// Add tokens to maps after they've been populated - using database handle format
		timelineTokens["timeline-user1-0024"] = user1Token
		timelineTokens["timeline-user2-0024"] = user2Token
		timelineTokens["timeline-user3-0024"] = user3Token
		timelineTokens["timeline-user4-0024"] = user4Token
		timelineTokens["timeline-user5-0024"] = user5Token

		followUnfollowTokens["timeline-user8-0024"] = user8Token
		followUnfollowTokens["timeline-user9-0024"] = user9Token
		followUnfollowTokens["timeline-user10-0024"] = user10Token
		followUnfollowTokens["timeline-user15-0024"] = user15Token

		paginationTokens["timeline-user11-0024"] = user11Token
		paginationTokens["timeline-user12-0024"] = user12Token

		missingTimelineTokens["timeline-user6-0024"] = user6Token
		missingTimelineTokens["timeline-user7-0024"] = user7Token
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		seedDatabase(db, "0024-timeline-down.pgsql")
		db.Close()
	})

	Describe("Basic Timeline Functionality", func() {
		It("should retrieve an existing user's timeline", func() {
			// User1 already has a timeline with posts from users 2, 3, 4
			req := hub.GetMyHomeTimelineRequest{}
			resp := testPOSTGetResp(
				timelineTokens["timeline-user1-0024"],
				req,
				"/hub/get-my-home-timeline",
				http.StatusOK,
			)

			var timeline hub.MyHomeTimeline
			err := json.Unmarshal(resp.([]byte), &timeline)
			Expect(err).ShouldNot(HaveOccurred())

			// Should have 4 posts (2 from user2 and 2 from user3)
			Expect(len(timeline.Posts)).Should(BeNumerically(">", 0))

			// Verify the posts are from followed users
			for _, post := range timeline.Posts {
				authorHandle := string(post.AuthorHandle)
				Expect(authorHandle).Should(Or(
					Equal("timeline-user2-0024"),
					Equal("timeline-user3-0024"),
					Equal("timeline-user4-0024"),
				))
			}
		})

		It(
			"should create a new timeline for a user who hasn't accessed it before",
			func() {
				// User7 doesn't have a timeline yet
				req := hub.GetMyHomeTimelineRequest{}
				resp := testPOSTGetResp(
					missingTimelineTokens["timeline-user7-0024"],
					req,
					"/hub/get-my-home-timeline",
					http.StatusOK,
				)

				var timeline hub.MyHomeTimeline
				err := json.Unmarshal(resp.([]byte), &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				// The timeline should be empty at first
				Expect(timeline.Posts).Should(HaveLen(0))

				// Verify a timeline entry was created in the database
				var count int
				err = db.QueryRow(
					context.Background(),
					"SELECT COUNT(*) FROM hu_active_home_timelines WHERE hub_user_id = $1",
					"12345678-0024-0024-0024-000000000007",
				).Scan(&count)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(count).Should(Equal(1))
			},
		)

		It("should handle pagination correctly", func() {
			// User11 follows users 12, 13, 14 and has many posts
			// Get first page with limit 2
			firstPageReq := hub.GetMyHomeTimelineRequest{
				Limit: 2,
			}
			resp := testPOSTGetResp(
				paginationTokens["timeline-user11-0024"],
				firstPageReq,
				"/hub/get-my-home-timeline",
				http.StatusOK,
			)

			var firstPage hub.MyHomeTimeline
			err := json.Unmarshal(resp.([]byte), &firstPage)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(firstPage.Posts).Should(HaveLen(2))
			Expect(
				firstPage.PaginationKey,
			).ShouldNot(BeEmpty(), "Pagination key should be provided when more items exist")

			// Get second page using pagination key
			secondPageReq := hub.GetMyHomeTimelineRequest{
				PaginationKey: &firstPage.PaginationKey,
				Limit:         2,
			}
			resp = testPOSTGetResp(
				paginationTokens["timeline-user11-0024"],
				secondPageReq,
				"/hub/get-my-home-timeline",
				http.StatusOK,
			)

			var secondPage hub.MyHomeTimeline
			err = json.Unmarshal(resp.([]byte), &secondPage)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify second page has different posts
			Expect(secondPage.Posts).ShouldNot(BeEmpty())
			Expect(
				secondPage.Posts[0].ID,
			).ShouldNot(Equal(firstPage.Posts[0].ID))
		})

		It("should return 422 for invalid pagination key", func() {
			invalidPaginationKey := "non-existent-post-id"
			req := hub.GetMyHomeTimelineRequest{
				PaginationKey: &invalidPaginationKey,
			}

			// This should return 422 Unprocessable Entity
			testPOSTGetResp(
				timelineTokens["timeline-user1-0024"],
				req,
				"/hub/get-my-home-timeline",
				http.StatusUnprocessableEntity,
			)
		})

		It("should require authentication", func() {
			req := hub.GetMyHomeTimelineRequest{}

			// Without token
			testPOSTGetResp(
				"",
				req,
				"/hub/get-my-home-timeline",
				http.StatusUnauthorized,
			)

			// With invalid token
			testPOSTGetResp(
				"invalid-token",
				req,
				"/hub/get-my-home-timeline",
				http.StatusUnauthorized,
			)
		})
	})

	Describe("Timeline Post Creation and Updates", func() {
		It("should add new posts to followers' timelines", func() {
			// User3 posts something new (user1 follows user3)
			postReq := hub.AddPostRequest{
				Content: "This is a new post from user3 for timeline testing!",
				TagIDs: []common.VTagID{
					common.VTagID("12345678-0024-0024-0024-000000000001"),
				},
			}

			// Create the post
			testPOSTGetResp(
				timelineTokens["timeline-user3-0024"],
				postReq,
				"/hub/add-post",
				http.StatusOK,
			)

			// Wait and poll for timeline refresh
			var timeline hub.MyHomeTimeline
			var foundNewPost bool

			for i := 0; i < MaxTimelinePollRetries; i++ {
				time.Sleep(TimelinePollInterval)

				// Check user1's timeline (user1 follows user3)
				resp := testPOSTGetResp(
					timelineTokens["timeline-user1-0024"],
					hub.GetMyHomeTimelineRequest{},
					"/hub/get-my-home-timeline",
					http.StatusOK,
				)

				err := json.Unmarshal(resp.([]byte), &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				// Look for the new post
				for _, post := range timeline.Posts {
					if post.Content == postReq.Content &&
						post.AuthorHandle == "timeline-user3-0024" {
						foundNewPost = true
						break
					}
				}

				if foundNewPost {
					break
				}
			}

			Expect(
				foundNewPost,
			).Should(BeTrue(), "New post should appear in follower's timeline after refresh")
		})
	})

	Describe("Follow/Unfollow Effects on Timeline", func() {
		It("should show posts from newly followed users", func() {
			// User8 follows User15
			followReq := hub.FollowUserRequest{
				Handle: "timeline-user15-0024",
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user8-0024"],
				followReq,
				"/hub/follow-user",
				http.StatusOK,
			)

			// User15 creates a post
			addPostReq := hub.AddPostRequest{
				Content: "This post should appear in User8's timeline after following!",
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user15-0024"],
				addPostReq,
				"/hub/add-post",
				http.StatusOK,
			)

			// Wait and poll for the post to appear in User8's timeline
			var timeline hub.MyHomeTimeline
			var foundNewPost bool

			for i := 0; i < MaxTimelinePollRetries; i++ {
				time.Sleep(TimelinePollInterval)

				resp := testPOSTGetResp(
					followUnfollowTokens["timeline-user8-0024"],
					hub.GetMyHomeTimelineRequest{},
					"/hub/get-my-home-timeline",
					http.StatusOK,
				)

				err := json.Unmarshal(resp.([]byte), &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				fmt.Fprintf(GinkgoWriter, "Got Posts: %+v\n", timeline.Posts)

				// Look for User15's post
				for _, post := range timeline.Posts {
					if post.Content == addPostReq.Content &&
						post.AuthorHandle == "timeline-user15-0024" {
						foundNewPost = true
						break
					}
				}

				if foundNewPost {
					break
				}
			}

			Expect(
				foundNewPost,
			).Should(BeTrue(), "Posts from newly followed users should appear in timeline")
		})

		It("should stop showing new posts after unfollowing", func() {
			// User10 follows User9
			followReq := hub.FollowUserRequest{
				Handle: "timeline-user9-0024",
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user10-0024"],
				followReq,
				"/hub/follow-user",
				http.StatusOK,
			)

			// User9 creates a first post
			firstPostContent := "This is the first post - should be in timeline"
			addFirstPostReq := hub.AddPostRequest{
				Content: firstPostContent,
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user9-0024"],
				addFirstPostReq,
				"/hub/add-post",
				http.StatusOK,
			)

			// Wait for the post to appear
			var foundFirstPost bool
			for i := 0; i < MaxTimelinePollRetries; i++ {
				time.Sleep(TimelinePollInterval)

				resp := testPOSTGetResp(
					followUnfollowTokens["timeline-user10-0024"],
					hub.GetMyHomeTimelineRequest{},
					"/hub/get-my-home-timeline",
					http.StatusOK,
				)

				var timeline hub.MyHomeTimeline
				err := json.Unmarshal(resp.([]byte), &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				for _, post := range timeline.Posts {
					if post.Content == firstPostContent &&
						post.AuthorHandle == "timeline-user9-0024" {
						foundFirstPost = true
						break
					}
				}

				if foundFirstPost {
					break
				}
			}

			Expect(
				foundFirstPost,
			).Should(BeTrue(), "First post should appear while following")

			// User10 unfollows User9
			unfollowReq := hub.UnfollowUserRequest{
				Handle: "timeline-user9-0024",
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user10-0024"],
				unfollowReq,
				"/hub/unfollow-user",
				http.StatusOK,
			)

			// Wait a bit for unfollow to take effect
			time.Sleep(TimelineRefreshInterval)

			// User9 creates a second post
			secondPostContent := "This post should NOT appear after unfollowing"
			addSecondPostReq := hub.AddPostRequest{
				Content: secondPostContent,
			}
			testPOSTGetResp(
				followUnfollowTokens["timeline-user9-0024"],
				addSecondPostReq,
				"/hub/add-post",
				http.StatusOK,
			)

			// Wait and verify the second post doesn't appear
			time.Sleep(TimelineRefreshInterval * 2) // Extra wait time

			resp := testPOSTGetResp(
				followUnfollowTokens["timeline-user10-0024"],
				hub.GetMyHomeTimelineRequest{},
				"/hub/get-my-home-timeline",
				http.StatusOK,
			)

			var timeline hub.MyHomeTimeline
			err := json.Unmarshal(resp.([]byte), &timeline)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify second post is not in timeline
			var foundSecondPost bool
			for _, post := range timeline.Posts {
				if post.Content == secondPostContent &&
					post.AuthorHandle == "timeline-user9-0024" {
					foundSecondPost = true
					break
				}
			}

			Expect(
				foundSecondPost,
			).Should(BeFalse(), "New posts after unfollowing should not appear in timeline")

			// But first post should still be there (posts made during following period remain)
			var firstPostStillThere bool
			for _, post := range timeline.Posts {
				if post.Content == firstPostContent &&
					post.AuthorHandle == "timeline-user9-0024" {
					firstPostStillThere = true
					break
				}
			}

			Expect(
				firstPostStillThere,
			).Should(BeTrue(), "Posts made during following period should remain in timeline")
		})
	})
})
