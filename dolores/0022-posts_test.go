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
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("Posts", Ordered, func() {
	var db *pgxpool.Pool
	var addUserToken, authTestUserToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0022-posts-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(2) // 2 hub users to sign in
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
		wg.Wait()
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		// We might need to explicitly delete posts if the down migration doesn't cover it
		// For now, assume the down migration handles users and cascades or manual cleanup exists.
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
						Tags:    []string{"test"},
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddPostRequest{
						Content: "Another post that should not be added.",
						Tags:    []string{"fail"},
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid post with content only",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: "This is my first post!",
						Tags:    []string{}, // Empty tags
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
						Tags:    []string{"golang", "testing", "backend"},
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
						Tags:    []string{"empty"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add post with content exceeding max length",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: strings.Repeat("x", 4097), // MaxLength is 4096
						Tags:    []string{"long"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add post with exactly max content length",
					token:       addUserToken,
					request: hub.AddPostRequest{
						Content: strings.Repeat("y", 4096), // Exactly MaxLength
						Tags:    []string{"maxlength"},
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
						Tags: []string{
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
						Tags: []string{
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
						Tags:    nil,
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

	// TODO: Add Describe blocks for GetTimeline tests later

})
