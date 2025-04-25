package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Posts", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool

		// User tokens
		addUserToken        string
		authTestUserToken   string
		getUser1Token       string
		getUser2Token       string
		getDetailsUserToken string

		// Post IDs
		getDetailsPostID string
		user1PostIDs     []string
		user2PostID      string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0022-posts-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(5) // Now 5 hub users to sign in
		hubSigninAsync(
			"add-user@0022-posts.example.com",
			"NewPassword123$",
			&addUserToken,
			&wg,
		)
		hubSigninAsync(
			"auth-user@0022-posts.example.com",
			"NewPassword123$",
			&authTestUserToken, // Using this token just to have another valid one if needed
			&wg,
		)
		hubSigninAsync(
			"get-user1@0022-posts.example.com",
			"NewPassword123$",
			&getUser1Token,
			&wg,
		)
		hubSigninAsync(
			"get-user2@0022-posts.example.com",
			"NewPassword123$",
			&getUser2Token,
			&wg,
		)
		hubSigninAsync(
			"get-details@0022-posts.example.com",
			"NewPassword123$",
			&getDetailsUserToken,
			&wg,
		)
		wg.Wait()

		// Create posts for get-user1
		// Post 1 with no tags
		post1Resp := testPOSTGetResp(
			getUser1Token,
			hub.AddPostRequest{
				Content: "First post by get-user1",
				NewTags: nil,
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post1AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post1Resp, &post1AddResp)).To(Succeed())
		user1PostIDs = append(user1PostIDs, post1AddResp.PostID)

		// Post 2 with tags
		post2Resp := testPOSTGetResp(
			getUser1Token,
			hub.AddPostRequest{
				Content: "Second post by get-user1, with tags",
				NewTags: []common.VTagName{"golang", "testing"},
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post2AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post2Resp, &post2AddResp)).To(Succeed())
		user1PostIDs = append(user1PostIDs, post2AddResp.PostID)

		// Post 3 with tag
		post3Resp := testPOSTGetResp(
			getUser1Token,
			hub.AddPostRequest{
				Content: "Third post, updated recently",
				NewTags: []common.VTagName{"pagination"},
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post3AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post3Resp, &post3AddResp)).To(Succeed())
		user1PostIDs = append(user1PostIDs, post3AddResp.PostID)

		// Post 4 with no tags
		post4Resp := testPOSTGetResp(
			getUser1Token,
			hub.AddPostRequest{
				Content: "Fourth post, newest",
				NewTags: nil,
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post4AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post4Resp, &post4AddResp)).To(Succeed())
		user1PostIDs = append(user1PostIDs, post4AddResp.PostID)

		// Create post for get-user2
		addPostReq := hub.AddPostRequest{
			Content: "First post by get-user2",
			NewTags: []common.VTagName{"specific-test"},
		}
		respBytes := testPOSTGetResp(
			getUser2Token,
			addPostReq,
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var addPostResp hub.AddPostResponse
		err := json.Unmarshal(respBytes, &addPostResp)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(addPostResp.PostID).ShouldNot(BeEmpty())
		user2PostID = addPostResp.PostID

		// Create the post for GetDetails tests
		addPostReq = hub.AddPostRequest{
			Content: "Post for GetDetails test (created via API)",
			NewTags: []common.VTagName{"details-tag1", "details-tag2"},
		}
		respBytes = testPOSTGetResp(
			getDetailsUserToken,
			addPostReq,
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		err = json.Unmarshal(respBytes, &addPostResp)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(addPostResp.PostID).ShouldNot(BeEmpty())
		getDetailsPostID = addPostResp.PostID // Store the ID

		fmt.Fprintf(GinkgoWriter, "### User1PostIDs: %v\n", user1PostIDs)
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		// Assumes 0022-posts-down.pgsql handles cleanup of users, posts, and tags created here.
		seedDatabase(pool, "0022-posts-down.pgsql")
		pool.Close()
	})

	Describe("Add Post", func() {
		type addPostTestCase struct {
			description string
			token       string
			request     hub.AddPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various add post test cases correctly", func() {
			testCases := []addPostTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddPostRequest{
						Content: "This post should not be added.",
						NewTags: []common.VTagName{"test"},
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddPostRequest{
						Content: "Another post that should not be added.",
						NewTags: []common.VTagName{"fail"},
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid post with content only",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "This is my first post!",
						NewTags: []common.VTagName{}, // Empty tags
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							response.PostID,
						).ShouldNot(BeEmpty(), "PostID should be returned")

						// TODO: Optionally verify the post was actually created in the DB
					},
				},
				{
					description: "add valid post with content and tags",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "Exploring the world of Go testing.",
						NewTags: []common.VTagName{
							"golang",
							"testing",
							"backend",
						},
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.PostID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add post with missing content",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "", // Missing content
						NewTags: []common.VTagName{"empty"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add post with content exceeding max length",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: strings.Repeat("x", 4097), // MaxLength is 4096
						NewTags: []common.VTagName{"long"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add post with exactly max content length",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: strings.Repeat("y", 4096), // Exactly MaxLength
						NewTags: []common.VTagName{"maxlength"},
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.PostID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add post with too many tags",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "Trying to add too many tags.",
						NewTags: []common.VTagName{
							"one",
							"two",
							"three",
							"four",
						}, // MaxItems is 3
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add post with exactly max tags",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "Testing maximum number of tags.",
						NewTags: []common.VTagName{
							"tag1",
							"tag2",
							"tag3",
						}, // Exactly MaxItems
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.PostID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add post with null tags (should be treated as empty)",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "Testing with null tags.",
						NewTags: nil,
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.PostID).ShouldNot(BeEmpty())
					},
				},
				// Add more edge cases if needed, e.g., tags with special characters, unicode content etc.
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing AddPost: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/add-post",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Get User Posts", func() {
		type getUserPostsTestCase struct {
			description string
			token       string                  // Auth token for the request
			request     hub.GetUserPostsRequest // The request body
			wantStatus  int
			validate    func([]byte) // Function to validate the response body
		}

		// Helper function to create a string pointer
		strPtr := func(s string) *string { return &s }
		hdlPtr := func(s string) *common.Handle { h := common.Handle(s); return &h }

		It("should handle various get user posts scenarios", func() {
			testCases := []getUserPostsTestCase{
				{
					description: "without authentication",
					token:       "",
					request:     hub.GetUserPostsRequest{}, // Fetch own posts (default)
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.GetUserPostsRequest{
						Handle: hdlPtr("get-user1"),
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "fetch own posts (user1) - default limit",
					token:       getUser1Token,
					request:     hub.GetUserPostsRequest{}, // No handle, no limit -> fetch self, default limit (10)
					wantStatus:  http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							resp.Posts,
						).Should(HaveLen(4))
						// User1 has 4 posts, default limit is 10
						Expect(
							resp.PaginationKey,
						).Should(BeEmpty())
						// No more posts
						// Check order (newest first based on created_at)
						// Since we created posts in order, the last one should be first
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user1PostIDs[3])) // post4ID
						Expect(
							resp.Posts[1].ID,
						).Should(Equal(user1PostIDs[2]))
						// post3ID
						Expect(
							resp.Posts[2].ID,
						).Should(Equal(user1PostIDs[1]))
						// post2ID
						Expect(
							resp.Posts[3].ID,
						).Should(Equal(user1PostIDs[0]))
						// post1ID
						// Check author details
						Expect(
							resp.Posts[0].AuthorHandle,
						).Should(Equal(common.Handle("get-user1")))
						Expect(
							resp.Posts[0].AuthorName,
						).Should(Equal("Get Posts User One"))
						// Check tags for the second post (should have golang and testing tags)
						Expect(
							resp.Posts[2].Tags,
						).Should(ConsistOf("golang", "testing"))

						// Validate voting and authorship fields for own posts
						for _, post := range resp.Posts {
							Expect(
								post.AmIAuthor,
							).Should(BeTrue(), "Should be marked as author for own posts")
							Expect(
								post.CanUpvote,
							).Should(BeFalse(), "Should not be able to upvote own posts")
							Expect(
								post.CanDownvote,
							).Should(BeFalse(), "Should not be able to downvote own posts")
							Expect(
								post.MeUpvoted,
							).Should(BeFalse(), "Should not have upvoted own posts")
							Expect(
								post.MeDownvoted,
							).Should(BeFalse(), "Should not have downvoted own posts")
						}
					},
				},
				{
					description: "fetch user1 posts by handle (from user2)",
					token:       getUser2Token, // Authenticated as user2
					request: hub.GetUserPostsRequest{
						Handle: hdlPtr("get-user1"),
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							resp.Posts,
						).Should(HaveLen(4))
						// User1 has 4 posts
						Expect(resp.PaginationKey).Should(BeEmpty())
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user1PostIDs[3])) // Newest post
						// Compare with string
						Expect(
							resp.Posts[0].AuthorHandle,
						).Should(Equal(common.Handle("get-user1")))

						// Validate voting and authorship fields when viewing other user's posts
						for _, post := range resp.Posts {
							Expect(
								post.AmIAuthor,
							).Should(BeFalse(), "Should not be marked as author for other's posts")
							Expect(
								post.CanUpvote,
							).Should(BeTrue(), "Should be able to upvote other's posts")
							Expect(
								post.CanDownvote,
							).Should(BeTrue(), "Should be able to downvote other's posts")
							Expect(
								post.MeUpvoted,
							).Should(BeFalse(), "Should not have upvoted other's posts yet")
							Expect(
								post.MeDownvoted,
							).Should(BeFalse(), "Should not have downvoted other's posts yet")
						}
					},
				},
				{
					description: "fetch user2 posts by handle (from user1)",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Handle: hdlPtr("get-user2"),
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							resp.Posts,
						).Should(HaveLen(1))
						// User2 has 1 post
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user2PostID))
						// Compare with string
						Expect(
							resp.Posts[0].AuthorHandle,
						).Should(Equal(common.Handle("get-user2")))
						Expect(
							resp.Posts[0].Tags,
						).Should(ConsistOf("specific-test"))
						Expect(resp.PaginationKey).Should(BeEmpty())
					},
				},
				{
					description: "fetch posts with invalid handle",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Handle: hdlPtr("non-existent-user"),
					},
					wantStatus: http.StatusNotFound, // Expecting 404 Not Found due to db.ErrNoHubUser
				},
				{
					description: "fetch own posts (user1) - limit 2",
					token:       getUser1Token,
					request:     hub.GetUserPostsRequest{Limit: 2},
					wantStatus:  http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(2))
						// Should get the two newest posts (last two in user1PostIDs)
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user1PostIDs[3])) // Newest post
						Expect(
							resp.Posts[1].ID,
						).Should(Equal(user1PostIDs[2]))
						// Second newest
						// Pagination key should be the ID of the last post returned
						Expect(
							resp.PaginationKey,
						).Should(Equal(user1PostIDs[2]))
						// The ID of the last post returned
					},
				},
				{
					description: "fetch own posts (user1) - limit 2, page 2",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit:         2,
						PaginationKey: strPtr(user1PostIDs[2]),
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(2))
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user1PostIDs[1]))
						Expect(resp.Posts[1].ID).Should(Equal(user1PostIDs[0]))
						Expect(
							resp.PaginationKey,
						).Should(Equal(user1PostIDs[0]))
						// Last post ID returned
					},
				},
				{
					description: "fetch own posts (user1) - limit 2, page 3 (empty)",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit:         2,
						PaginationKey: strPtr(user1PostIDs[0]),
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(BeEmpty())
						Expect(resp.PaginationKey).Should(BeEmpty())
					},
				},
				{
					description: "fetch posts with invalid pagination key",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit:         5,
						PaginationKey: strPtr("invalid-post-id"),
					},
					wantStatus: http.StatusOK, // Should return the first page if pagination key is invalid
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							resp.Posts,
						).Should(HaveLen(4))
						// Should get all 4 posts (limit 5)
						Expect(
							resp.Posts[0].ID,
						).Should(Equal(user1PostIDs[3])) // Newest post
						// Compare with string
						Expect(
							resp.PaginationKey,
						).Should(BeEmpty())
						// No more posts
					},
				},
				{
					description: "fetch posts with zero limit (should use default)",
					token:       getUser1Token,
					request:     hub.GetUserPostsRequest{Limit: 0},
					wantStatus:  http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(4))
						// Default limit is 10, gets all 4
						Expect(resp.PaginationKey).Should(BeEmpty())
					},
				},
				{
					description: "fetch posts with negative limit (should fail validation)",
					token:       getUser1Token,
					request:     hub.GetUserPostsRequest{Limit: -1},
					wantStatus:  http.StatusBadRequest,
				},
				{
					description: "fetch posts with limit exceeding max (should fail validation)",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit: 41,
					}, // Max is 40
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing GetUserPosts: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/get-user-posts",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Get Post Details", func() {
		type getPostDetailsTestCase struct {
			description string
			token       string                    // Auth token for the request
			request     hub.GetPostDetailsRequest // The request body
			wantStatus  int
			validate    func([]byte) // Function to validate the response body (common.Post)
		}

		It("should handle various get post details scenarios", func() {
			testCases := []getPostDetailsTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.GetPostDetailsRequest{
						PostID: getDetailsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.GetPostDetailsRequest{
						PostID: getDetailsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "fetch valid post details",
					token:       getDetailsUserToken,
					request: hub.GetPostDetailsRequest{
						PostID: getDetailsPostID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.Post
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.ID).Should(Equal(getDetailsPostID))
						Expect(
							resp.Content,
						).Should(Equal("Post for GetDetails test (created via API)"))
						Expect(
							resp.AuthorHandle,
						).Should(Equal(common.Handle("get-details-user")))
						Expect(
							resp.AuthorName,
						).Should(Equal("Get Details User"))
						Expect(
							resp.Tags,
						).Should(ConsistOf("details-tag1", "details-tag2"))
						Expect(resp.CreatedAt).ShouldNot(BeZero())

						// Validate voting and authorship fields for post author
						Expect(
							resp.AmIAuthor,
						).Should(BeTrue(), "Author should see AmIAuthor as true")
						Expect(
							resp.CanUpvote,
						).Should(BeFalse(), "Author should not be able to upvote own post")
						Expect(
							resp.CanDownvote,
						).Should(BeFalse(), "Author should not be able to downvote own post")
						Expect(
							resp.MeUpvoted,
						).Should(BeFalse(), "Author should not have upvoted their own post")
						Expect(
							resp.MeDownvoted,
						).Should(BeFalse(), "Author should not have downvoted their own post")
					},
				},
				{
					description: "fetch non-existent post",
					token:       getDetailsUserToken,
					request: hub.GetPostDetailsRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "fetch post details as non-author",
					token:       getUser1Token, // Using a different user's token
					request: hub.GetPostDetailsRequest{
						PostID: getDetailsPostID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.Post
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())

						// Basic post details should match
						Expect(resp.ID).Should(Equal(getDetailsPostID))
						Expect(
							resp.AuthorHandle,
						).Should(Equal(common.Handle("get-details-user")))

						// Validate voting and authorship fields for non-author
						Expect(
							resp.AmIAuthor,
						).Should(BeFalse(), "Non-author should see AmIAuthor as false")
						Expect(
							resp.CanUpvote,
						).Should(BeTrue(), "Non-author should be able to upvote post")
						Expect(
							resp.CanDownvote,
						).Should(BeTrue(), "Non-author should be able to downvote post")
						Expect(
							resp.MeUpvoted,
						).Should(BeFalse(), "Non-author should not have upvoted post yet")
						Expect(
							resp.MeDownvoted,
						).Should(BeFalse(), "Non-author should not have downvoted post yet")
					},
				},
				{
					description: "fetch with empty post ID (should fail validation)",
					token:       getDetailsUserToken,
					request:     hub.GetPostDetailsRequest{PostID: ""},
					wantStatus:  http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing GetPostDetails: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/get-post-details", // Ensure this matches your actual endpoint path
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
