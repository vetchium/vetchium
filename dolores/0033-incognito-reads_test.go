package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Incognito Posts Read APIs", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool

		// User tokens
		aliceToken   string
		bobToken     string
		charlieToken string
		eveToken     string
		frankToken   string
		graceToken   string

		// Test post IDs created via APIs
		alicePost1ID   string
		bobPost1ID     string
		charliePost1ID string
		deletedPostID  string

		// Test comment IDs created via APIs
		aliceComment1ID   string
		bobComment1ID     string
		charlieComment1ID string
		nestedCommentID   string
	)

	createTestPosts := func() {
		// Alice creates a technology post
		alicePost1 := hub.AddIncognitoPostRequest{
			Content: "Alice's first technology post about programming",
			TagIDs:  []common.VTagID{"technology"},
		}

		// Debug: Log before making the API call
		fmt.Printf(
			"DEBUG: About to create Alice's post with token: %s\n",
			aliceToken,
		)
		fmt.Printf("DEBUG: Request body: %+v\n", alicePost1)

		resp1 := testPOSTGetResp(
			aliceToken,
			alicePost1,
			"/hub/add-incognito-post",
			http.StatusOK,
		)

		// Debug: Log after successful API call
		fmt.Printf(
			"DEBUG: Successfully created Alice's post, response: %v\n",
			resp1,
		)

		var addResp1 hub.AddIncognitoPostResponse
		err := json.Unmarshal(resp1.([]byte), &addResp1)
		Expect(err).ShouldNot(HaveOccurred())
		alicePost1ID = addResp1.IncognitoPostID

		fmt.Printf("DEBUG: Alice's post ID: %s\n", alicePost1ID)

		// Alice creates a second post (different topic)
		alicePost2 := hub.AddIncognitoPostRequest{
			Content: "Alice's second post about startup entrepreneurship",
			TagIDs:  []common.VTagID{"startups"},
		}

		fmt.Printf(
			"DEBUG: About to create Alice's second post with token: %s\n",
			aliceToken,
		)

		resp2 := testPOSTGetResp(
			aliceToken,
			alicePost2,
			"/hub/add-incognito-post",
			http.StatusOK,
		)
		var addResp2 hub.AddIncognitoPostResponse
		err = json.Unmarshal(resp2.([]byte), &addResp2)
		Expect(err).ShouldNot(HaveOccurred())
		// We don't need to store this ID in a variable since it's not used elsewhere

		fmt.Printf(
			"DEBUG: Alice's second post ID: %s\n",
			addResp2.IncognitoPostID,
		)

		// Bob creates a mentorship post
		bobPost1 := hub.AddIncognitoPostRequest{
			Content: "Bob's guide to effective mentorship in tech industry",
			TagIDs:  []common.VTagID{"mentorship"},
		}

		fmt.Printf(
			"DEBUG: About to create Bob's post with token: %s\n",
			bobToken,
		)

		resp3 := testPOSTGetResp(
			bobToken,
			bobPost1,
			"/hub/add-incognito-post",
			http.StatusOK,
		)
		var addResp3 hub.AddIncognitoPostResponse
		err = json.Unmarshal(resp3.([]byte), &addResp3)
		Expect(err).ShouldNot(HaveOccurred())
		bobPost1ID = addResp3.IncognitoPostID

		// Charlie creates a technology post
		charliePost1 := hub.AddIncognitoPostRequest{
			Content: "Charlie's insights on modern software architecture patterns",
			TagIDs:  []common.VTagID{"technology"},
		}

		fmt.Printf(
			"DEBUG: About to create Charlie's post with token: %s\n",
			charlieToken,
		)

		resp4 := testPOSTGetResp(
			charlieToken,
			charliePost1,
			"/hub/add-incognito-post",
			http.StatusOK,
		)
		var addResp4 hub.AddIncognitoPostResponse
		err = json.Unmarshal(resp4.([]byte), &addResp4)
		Expect(err).ShouldNot(HaveOccurred())
		charliePost1ID = addResp4.IncognitoPostID

		// Frank creates a post that will be deleted
		deletedPost := hub.AddIncognitoPostRequest{
			Content: "Frank's post that will be deleted for testing",
			TagIDs:  []common.VTagID{"technology"},
		}

		fmt.Printf(
			"DEBUG: About to create Frank's post with token: %s\n",
			frankToken,
		)

		resp6 := testPOSTGetResp(
			frankToken,
			deletedPost,
			"/hub/add-incognito-post",
			http.StatusOK,
		)
		var addResp6 hub.AddIncognitoPostResponse
		err = json.Unmarshal(resp6.([]byte), &addResp6)
		Expect(err).ShouldNot(HaveOccurred())
		deletedPostID = addResp6.IncognitoPostID

		// Delete Frank's post
		deleteReq := hub.DeleteIncognitoPostRequest{
			IncognitoPostID: deletedPostID,
		}
		testPOST(
			frankToken,
			deleteReq,
			"/hub/delete-incognito-post",
			http.StatusOK,
		)
	}

	createTestComments := func() {
		// Alice comments on Charlie's tech post
		aliceComment1 := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: charliePost1ID,
			Content:         "Alice's insightful comment on Charlie's architecture post",
		}
		resp1 := testPOSTGetResp(
			aliceToken,
			aliceComment1,
			"/hub/add-incognito-post-comment",
			http.StatusOK,
		)
		var commentResp1 hub.AddIncognitoPostCommentResponse
		err := json.Unmarshal(resp1.([]byte), &commentResp1)
		Expect(err).ShouldNot(HaveOccurred())
		aliceComment1ID = commentResp1.CommentID

		// Bob comments on Alice's tech post
		bobComment1 := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: alicePost1ID,
			Content:         "Bob's helpful feedback on Alice's programming post",
		}
		resp2 := testPOSTGetResp(
			bobToken,
			bobComment1,
			"/hub/add-incognito-post-comment",
			http.StatusOK,
		)
		var commentResp2 hub.AddIncognitoPostCommentResponse
		err = json.Unmarshal(resp2.([]byte), &commentResp2)
		Expect(err).ShouldNot(HaveOccurred())
		bobComment1ID = commentResp2.CommentID

		// Charlie comments on Bob's mentorship post
		charlieComment1 := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: bobPost1ID,
			Content:         "Charlie's thoughts on effective mentorship strategies",
		}
		resp3 := testPOSTGetResp(
			charlieToken,
			charlieComment1,
			"/hub/add-incognito-post-comment",
			http.StatusOK,
		)
		var commentResp3 hub.AddIncognitoPostCommentResponse
		err = json.Unmarshal(resp3.([]byte), &commentResp3)
		Expect(err).ShouldNot(HaveOccurred())
		charlieComment1ID = commentResp3.CommentID

		// Alice replies to Bob's comment (nested comment)
		nestedComment := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: alicePost1ID,
			Content:         "Alice's nested reply to Bob's helpful feedback",
			InReplyTo:       &bobComment1ID,
		}
		resp4 := testPOSTGetResp(
			aliceToken,
			nestedComment,
			"/hub/add-incognito-post-comment",
			http.StatusOK,
		)
		var commentResp4 hub.AddIncognitoPostCommentResponse
		err = json.Unmarshal(resp4.([]byte), &commentResp4)
		Expect(err).ShouldNot(HaveOccurred())
		nestedCommentID = commentResp4.CommentID

		// Alice adds another comment on Bob's mentorship post
		aliceComment2 := hub.AddIncognitoPostCommentRequest{
			IncognitoPostID: bobPost1ID,
			Content:         "Alice's thoughts on mentorship best practices",
		}
		resp5 := testPOSTGetResp(
			aliceToken,
			aliceComment2,
			"/hub/add-incognito-post-comment",
			http.StatusOK,
		)
		var commentResp5 hub.AddIncognitoPostCommentResponse
		err = json.Unmarshal(resp5.([]byte), &commentResp5)
		Expect(err).ShouldNot(HaveOccurred())
		// We don't need to store this ID since it's not used elsewhere

		// Delete one of Alice's comments to test deleted comment inclusion
		deleteCommentReq := hub.DeleteIncognitoPostCommentRequest{
			IncognitoPostID: charliePost1ID,
			CommentID:       aliceComment1ID,
		}
		testPOST(
			aliceToken,
			deleteCommentReq,
			"/hub/delete-incognito-post-comment",
			http.StatusOK,
		)
	}

	performTestVotingScenarios := func() {
		// Scenario 1: Multiple users upvote Charlie's post
		upvoteCharlie := hub.UpvoteIncognitoPostRequest{
			IncognitoPostID: charliePost1ID,
		}
		testPOST(
			aliceToken,
			upvoteCharlie,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)
		testPOST(
			bobToken,
			upvoteCharlie,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)
		testPOST(
			eveToken,
			upvoteCharlie,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)

		// Scenario 2: Mixed voting on Alice's first post
		upvoteAlice1 := hub.UpvoteIncognitoPostRequest{
			IncognitoPostID: alicePost1ID,
		}
		downvoteAlice1 := hub.DownvoteIncognitoPostRequest{
			IncognitoPostID: alicePost1ID,
		}
		testPOST(
			charlieToken,
			upvoteAlice1,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)
		testPOST(
			eveToken,
			downvoteAlice1,
			"/hub/downvote-incognito-post",
			http.StatusOK,
		)

		// Scenario 3: Vote then unvote on Bob's post
		upvoteBob := hub.UpvoteIncognitoPostRequest{IncognitoPostID: bobPost1ID}
		unvoteBob := hub.UnvoteIncognitoPostRequest{IncognitoPostID: bobPost1ID}
		testPOST(
			graceToken,
			upvoteBob,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)
		testPOST(
			graceToken,
			unvoteBob,
			"/hub/unvote-incognito-post",
			http.StatusOK,
		)
		testPOST(
			frankToken,
			upvoteBob,
			"/hub/upvote-incognito-post",
			http.StatusOK,
		)

		// Scenario 4: Comment voting
		upvoteBobComment := hub.UpvoteIncognitoPostCommentRequest{
			IncognitoPostID: alicePost1ID,
			CommentID:       bobComment1ID,
		}
		testPOST(
			charlieToken,
			upvoteBobComment,
			"/hub/upvote-incognito-post-comment",
			http.StatusOK,
		)
		testPOST(
			eveToken,
			upvoteBobComment,
			"/hub/upvote-incognito-post-comment",
			http.StatusOK,
		)

		downvoteCharlieComment := hub.DownvoteIncognitoPostCommentRequest{
			IncognitoPostID: bobPost1ID,
			CommentID:       charlieComment1ID,
		}
		testPOST(
			graceToken,
			downvoteCharlieComment,
			"/hub/downvote-incognito-post-comment",
			http.StatusOK,
		)
	}

	BeforeAll(func() {
		fmt.Printf("DEBUG: Starting BeforeAll setup\n")

		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		fmt.Printf("DEBUG: Database setup complete\n")

		seedDatabase(pool, "0033-incognito-reads-up.pgsql")
		fmt.Printf(
			"DEBUG: Database seeded with 0033-incognito-reads-up.pgsql\n",
		)

		// Login hub users and get tokens using async signin
		// Note: Diana is excluded as she's DISABLED_HUB_USER and cannot log in
		fmt.Printf("DEBUG: Starting user authentication\n")
		var wg sync.WaitGroup
		wg.Add(6)
		hubSigninAsync(
			"alice@test0033.com",
			"NewPassword123$",
			&aliceToken,
			&wg,
		)
		hubSigninAsync("bob@test0033.com", "NewPassword123$", &bobToken, &wg)
		hubSigninAsync(
			"charlie@company0033.com",
			"NewPassword123$",
			&charlieToken,
			&wg,
		)
		hubSigninAsync("eve@test0033.com", "NewPassword123$", &eveToken, &wg)
		hubSigninAsync(
			"frank@test0033.com",
			"NewPassword123$",
			&frankToken,
			&wg,
		)
		hubSigninAsync(
			"grace@test0033.com",
			"NewPassword123$",
			&graceToken,
			&wg,
		)
		wg.Wait()

		fmt.Printf("DEBUG: All users authenticated successfully\n")
		fmt.Printf("DEBUG: Alice token: %s\n", aliceToken)
		fmt.Printf("DEBUG: Bob token: %s\n", bobToken)
		fmt.Printf("DEBUG: Charlie token: %s\n", charlieToken)
		fmt.Printf("DEBUG: Eve token: %s\n", eveToken)
		fmt.Printf("DEBUG: Frank token: %s\n", frankToken)
		fmt.Printf("DEBUG: Grace token: %s\n", graceToken)

		// Create test posts via APIs for consistent testing
		fmt.Printf("DEBUG: About to start creating test posts\n")
		createTestPosts()
		fmt.Printf("DEBUG: Test posts created successfully\n")

		createTestComments()
		fmt.Printf("DEBUG: Test comments created successfully\n")

		performTestVotingScenarios()
		fmt.Printf("DEBUG: Test voting scenarios completed\n")

		fmt.Printf("DEBUG: BeforeAll setup completed successfully\n")
	})

	AfterAll(func() {
		seedDatabase(pool, "0033-incognito-reads-down.pgsql")
		pool.Close()
	})

	Describe("GetIncognitoPosts", func() {
		Describe("Successful retrieval and filtering", func() {
			It(
				"should get technology posts with proper sorting by score",
				func() {
					reqBody := hub.GetIncognitoPostsRequest{
						TagID: common.VTagID("technology"),
						Limit: 25,
					}

					respData := testPOSTGetResp(
						aliceToken,
						reqBody,
						"/hub/get-incognito-posts",
						http.StatusOK,
					)

					var response hub.GetIncognitoPostsResponse
					err := json.Unmarshal(respData.([]byte), &response)
					Expect(err).ShouldNot(HaveOccurred())

					// Should exclude deleted posts
					Expect(
						len(response.Posts),
					).Should(Equal(2))
					// Charlie's and Alice's tech posts

					// Verify sorting by score (desc), then creation date (desc)
					// Charlie's post should have higher score due to multiple upvotes
					charliePostFound := false
					alicePostFound := false
					for _, post := range response.Posts {
						if post.IncognitoPostID == charliePost1ID {
							charliePostFound = true
							Expect(
								post.UpvotesCount,
							).Should(Equal(int32(3)))
							// Alice, Bob, Eve
							Expect(post.Score).Should(Equal(int32(3)))
							Expect(
								post.MeUpvoted,
							).Should(BeTrue())
							// Alice upvoted
							Expect(
								post.CanUpvote,
							).Should(BeFalse())
							// Already voted
						} else if post.IncognitoPostID == alicePost1ID {
							alicePostFound = true
							Expect(post.UpvotesCount).Should(Equal(int32(1)))   // Charlie
							Expect(post.DownvotesCount).Should(Equal(int32(1))) // Eve
							Expect(post.Score).Should(Equal(int32(0)))
							Expect(post.IsCreatedByMe).Should(BeTrue())
							Expect(post.CanUpvote).Should(BeFalse()) // Can't vote on own post
						}
					}
					Expect(charliePostFound).Should(BeTrue())
					Expect(alicePostFound).Should(BeTrue())
				},
			)

			It("should filter by mentorship tag", func() {
				reqBody := hub.GetIncognitoPostsRequest{
					TagID: common.VTagID("mentorship"),
					Limit: 25,
				}

				respData := testPOSTGetResp(
					bobToken,
					reqBody,
					"/hub/get-incognito-posts",
					http.StatusOK,
				)

				var response hub.GetIncognitoPostsResponse
				err := json.Unmarshal(respData.([]byte), &response)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(len(response.Posts)).Should(Equal(1))
				Expect(
					response.Posts[0].IncognitoPostID,
				).Should(Equal(bobPost1ID))
				Expect(response.Posts[0].IsCreatedByMe).Should(BeTrue())
				Expect(
					response.Posts[0].UpvotesCount,
				).Should(Equal(int32(1)))
				// Frank upvoted
			})

			It("should handle pagination correctly", func() {
				reqBody := hub.GetIncognitoPostsRequest{
					TagID: common.VTagID("technology"),
					Limit: 1,
				}

				// Get first page
				respData := testPOSTGetResp(
					eveToken,
					reqBody,
					"/hub/get-incognito-posts",
					http.StatusOK,
				)
				var response hub.GetIncognitoPostsResponse
				err := json.Unmarshal(respData.([]byte), &response)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(len(response.Posts)).Should(Equal(1))
				Expect(response.PaginationKey).ShouldNot(BeEmpty())
				firstPostID := response.Posts[0].IncognitoPostID

				// Get second page
				reqBody.PaginationKey = &response.PaginationKey
				respData2 := testPOSTGetResp(
					eveToken,
					reqBody,
					"/hub/get-incognito-posts",
					http.StatusOK,
				)
				var response2 hub.GetIncognitoPostsResponse
				err = json.Unmarshal(respData2.([]byte), &response2)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(len(response2.Posts)).Should(Equal(1))
				Expect(
					response2.Posts[0].IncognitoPostID,
				).ShouldNot(Equal(firstPostID))
			})
		})

		Describe("Error handling", func() {
			It("should fail without authentication", func() {
				reqBody := hub.GetIncognitoPostsRequest{
					TagID: common.VTagID("technology"),
					Limit: 25,
				}
				testPOST(
					"",
					reqBody,
					"/hub/get-incognito-posts",
					http.StatusUnauthorized,
				)
			})

			It("should fail with invalid parameters", func() {
				reqBody := hub.GetIncognitoPostsRequest{
					TagID: common.VTagID("technology"),
					Limit: 0, // Invalid
				}
				testPOST(
					aliceToken,
					reqBody,
					"/hub/get-incognito-posts",
					http.StatusBadRequest,
				)
			})
		})
	})

	Describe("GetMyIncognitoPosts", func() {
		It("should return user's own posts sorted by creation date", func() {
			reqBody := hub.GetMyIncognitoPostsRequest{Limit: 25}

			respData := testPOSTGetResp(
				aliceToken,
				reqBody,
				"/hub/get-my-incognito-posts",
				http.StatusOK,
			)
			var response hub.GetMyIncognitoPostsResponse
			err := json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())

			// Alice has 2 posts
			Expect(len(response.Posts)).Should(Equal(2))

			// All should be marked as created by me
			for _, post := range response.Posts {
				Expect(post.IsCreatedByMe).Should(BeTrue())
				Expect(
					post.CanUpvote,
				).Should(BeFalse())
				// Can't vote on own posts
			}

			// Verify vote counts reflect actual voting
			for _, post := range response.Posts {
				if post.IncognitoPostID == alicePost1ID {
					Expect(post.UpvotesCount).Should(Equal(int32(1)))
					Expect(post.DownvotesCount).Should(Equal(int32(1)))
				}
			}
		})

		It("should return empty for user with no posts", func() {
			reqBody := hub.GetMyIncognitoPostsRequest{Limit: 25}

			respData := testPOSTGetResp(
				graceToken,
				reqBody,
				"/hub/get-my-incognito-posts",
				http.StatusOK,
			)
			var response hub.GetMyIncognitoPostsResponse
			err := json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(len(response.Posts)).Should(Equal(0))
		})
	})

	Describe("GetMyIncognitoPostComments", func() {
		It("should return user's own comments with post context", func() {
			reqBody := hub.GetMyIncognitoPostCommentsRequest{Limit: 25}

			respData := testPOSTGetResp(
				aliceToken,
				reqBody,
				"/hub/get-my-incognito-post-comments",
				http.StatusOK,
			)
			var response hub.GetMyIncognitoPostCommentsResponse
			err := json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())

			// Alice has 3 comments (including deleted and nested)
			Expect(len(response.Comments)).Should(Equal(3))

			// Verify post context is provided
			for _, comment := range response.Comments {
				Expect(comment.IncognitoPostID).ShouldNot(BeEmpty())
				Expect(comment.PostContentPreview).ShouldNot(BeEmpty())
				Expect(len(comment.PostTags)).Should(BeNumerically(">", 0))
			}

			// Find the nested comment and verify depth
			nestedFound := false
			deletedFound := false
			for _, comment := range response.Comments {
				if comment.CommentID == nestedCommentID {
					nestedFound = true
					Expect(comment.Depth).Should(Equal(int32(1)))
					Expect(comment.InReplyTo).ShouldNot(BeNil())
				}
				if comment.IsDeleted {
					deletedFound = true
				}
			}
			Expect(nestedFound).Should(BeTrue())
			Expect(deletedFound).Should(BeTrue())
		})
	})

	Describe("Complex voting scenarios", func() {
		It("should handle vote changes correctly", func() {
			// Create a new post for testing vote changes
			testPost := hub.AddIncognitoPostRequest{
				Content: "Test post for vote change scenarios",
				TagIDs:  []common.VTagID{"technology"},
			}
			resp := testPOSTGetResp(
				graceToken,
				testPost,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(resp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())
			testPostID := addResp.IncognitoPostID

			// Alice upvotes
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: testPostID,
			}
			testPOST(
				aliceToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Bob downvotes
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: testPostID,
			}
			testPOST(
				bobToken,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusOK,
			)

			// Charlie upvotes then unvotes
			testPOST(
				charlieToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: testPostID,
			}
			testPOST(
				charlieToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusOK,
			)

			// Verify final vote state through GetIncognitoPost
			getReq := hub.GetIncognitoPostRequest{IncognitoPostID: testPostID}
			getResp := testPOSTGetResp(
				aliceToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(
				getResponse.UpvotesCount,
			).Should(Equal(int32(1)))
			// Only Alice
			Expect(
				getResponse.DownvotesCount,
			).Should(Equal(int32(1)))
			// Only Bob
			Expect(
				getResponse.Score,
			).Should(Equal(int32(0)))
			// 1 - 1 = 0
			Expect(
				getResponse.MeUpvoted,
			).Should(BeTrue())
			// Alice's perspective
		})

		It("should handle idempotent voting correctly", func() {
			// Test that multiple same votes don't change counts
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: charliePost1ID,
			}

			// Alice upvotes multiple times (already upvoted in setup)
			testPOST(
				aliceToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)
			testPOST(
				aliceToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Verify count hasn't changed
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: charliePost1ID,
			}
			getResp := testPOSTGetResp(
				aliceToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err := json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(
				getResponse.UpvotesCount,
			).Should(Equal(int32(3)))
			// Still 3, not more
		})
	})
})
