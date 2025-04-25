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
	var db *pgxpool.Pool
	var addUserToken, authTestUserToken, getUser1Token, getUser2Token string
	var getDetailsUserToken string
	var getDetailsPostID string // Variable to store the dynamically created post ID

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0022-posts-up.pgsql")

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

		// Create the post for GetDetails tests via API
		addPostReq := hub.AddPostRequest{
			Content: "Post for GetDetails test (created via API)",
			NewTags: []common.VTagName{"details-tag1", "details-tag2"},
		}
		respBytes := testPOSTGetResp(
			getDetailsUserToken,
			addPostReq,
			"/hub/add-post",
			http.StatusOK,
		).([]byte)

		var addPostResp hub.AddPostResponse
		err := json.Unmarshal(respBytes, &addPostResp)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(addPostResp.PostID).ShouldNot(BeEmpty())
		getDetailsPostID = addPostResp.PostID // Store the ID
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		// Assumes 0022-posts-down.pgsql handles cleanup of users, posts, and tags created here.
		seedDatabase(db, "0022-posts-down.pgsql")
		db.Close()
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
						// Check order (newest first based on updated_at)
						Expect(
							resp.Posts[0].ID,
						).Should(Equal("post-g1-04"))
						// Compare with string
						Expect(resp.Posts[1].ID).Should(Equal("post-g1-03"))
						Expect(resp.Posts[2].ID).Should(Equal("post-g1-02"))
						Expect(resp.Posts[3].ID).Should(Equal("post-g1-01"))
						// Check author details
						Expect(
							resp.Posts[0].AuthorHandle,
						).Should(Equal(common.Handle("get-user1")))
						Expect(
							resp.Posts[0].AuthorName,
						).Should(Equal("Get Posts User One"))
						// Check tags for post-g1-02
						Expect(
							resp.Posts[2].Tags,
						).Should(ConsistOf("golang", "testing"))
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
						).Should(Equal("post-g1-04"))
						// Compare with string
						Expect(
							resp.Posts[0].AuthorHandle,
						).Should(Equal(common.Handle("get-user1")))
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
						).Should(Equal("post-g2-01"))
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
						Expect(
							resp.Posts[0].ID,
						).Should(Equal("post-g1-04"))
						// Compare with string
						Expect(resp.Posts[1].ID).Should(Equal("post-g1-03"))
						Expect(
							resp.PaginationKey,
						).Should(Equal("post-g1-03"))
						// The ID of the last post returned
					},
				},
				{
					description: "fetch own posts (user1) - limit 2, page 2",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit:         2,
						PaginationKey: strPtr("post-g1-03"),
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp hub.GetUserPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(2))
						Expect(
							resp.Posts[0].ID,
						).Should(Equal("post-g1-02"))
						// Compare with string
						Expect(resp.Posts[1].ID).Should(Equal("post-g1-01"))
						Expect(
							resp.PaginationKey,
						).Should(Equal("post-g1-01"))
						// Last post ID returned
					},
				},
				{
					description: "fetch own posts (user1) - limit 2, page 3 (empty)",
					token:       getUser1Token,
					request: hub.GetUserPostsRequest{
						Limit:         2,
						PaginationKey: strPtr("post-g1-01"),
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
						).Should(Equal("post-g1-04"))
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
