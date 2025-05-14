package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = FDescribe("Employer Posts", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0026-employerpost-up.pgsql")
	})

	AfterAll(func() {
		// Clean up the database using the down migration
		seedDatabase(pool, "0026-employerpost-down.pgsql")
		pool.Close()
	})

	// Helper function to create a test post
	createTestPost := func(token string, content string, tagIDs []common.VTagID, newTags []common.VTagName) string {
		request := employer.AddEmployerPostRequest{
			Content: content,
			TagIDs:  tagIDs,
			NewTags: newTags,
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/add-post",
			http.StatusOK,
		).([]byte)

		var addResp employer.AddEmployerPostResponse
		err := json.Unmarshal(resp, &addResp)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(addResp.PostID).ShouldNot(BeEmpty())

		return addResp.PostID
	}

	// Helper function to create multiple test posts
	createTestPosts := func(token string, count int) []string {
		postIDs := make([]string, count)
		for i := 0; i < count; i++ {
			postIDs[i] = createTestPost(
				token,
				fmt.Sprintf("Test post %d", i+1),
				nil,
				nil,
			)
		}
		return postIDs
	}

	Describe("Add Employer Post", func() {
		var (
			adminToken     string
			marketingToken string
			regularToken   string
		)

		BeforeEach(func() {
			// Login org users and get tokens
			var wg sync.WaitGroup
			wg.Add(3) // 3 org users to sign in

			employerSigninAsync(
				"0026-employerposts.example.com",
				"admin@0026-employerposts.example.com",
				"NewPassword123$",
				&adminToken,
				&wg,
			)

			employerSigninAsync(
				"0026-employerposts.example.com",
				"marketing@0026-employerposts.example.com",
				"NewPassword123$",
				&marketingToken,
				&wg,
			)

			employerSigninAsync(
				"0026-employerposts.example.com",
				"regular@0026-employerposts.example.com",
				"NewPassword123$",
				&regularToken,
				&wg,
			)

			wg.Wait()
		})

		type addEmployerPostTestCase struct {
			description string
			token       string
			request     employer.AddEmployerPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various add employer post scenarios", func() {
			testCases := []addEmployerPostTestCase{
				{
					description: "without authentication",
					token:       "",
					request: employer.AddEmployerPostRequest{
						Content: "Test post without auth",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: employer.AddEmployerPostRequest{
						Content: "Test post with invalid token",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "admin can add post",
					token:       adminToken,
					request: employer.AddEmployerPostRequest{
						Content: "New post by admin",
						NewTags: []common.VTagName{"admin-tag"},
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.AddEmployerPostResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.PostID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "marketing user can add post",
					token:       marketingToken,
					request: employer.AddEmployerPostRequest{
						Content: "New post by marketing",
						TagIDs: []common.VTagID{
							common.VTagID(
								"12345678-0026-0026-0026-000000050002",
							),
						}, // marketing tag
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.AddEmployerPostResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.PostID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "regular user cannot add post",
					token:       regularToken,
					request: employer.AddEmployerPostRequest{
						Content: "New post by regular user",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "empty content",
					token:       adminToken,
					request: employer.AddEmployerPostRequest{
						Content: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "too many tags",
					token:       adminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with too many tags",
						NewTags: []common.VTagName{
							"tag1",
							"tag2",
							"tag3",
							"tag4",
						},
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing AddEmployerPost: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/add-post",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Get Employer Post", func() {
		var (
			adminToken string
			postID     string
		)

		BeforeEach(func() {
			// Login admin user
			var wg sync.WaitGroup
			wg.Add(1)

			employerSigninAsync(
				"0026-employerposts2.example.com",
				"admin@0026-employerposts2.example.com",
				"NewPassword123$",
				&adminToken,
				&wg,
			)

			wg.Wait()

			// Create a fresh post for each test
			postID = createTestPost(
				adminToken,
				"Test post for get",
				[]common.VTagID{
					common.VTagID(
						"12345678-0026-0026-0026-000000050001",
					), // engineering
					common.VTagID(
						"12345678-0026-0026-0026-000000050003",
					), // golang
				},
				nil,
			)
		})

		type getEmployerPostTestCase struct {
			description string
			token       string
			request     employer.GetEmployerPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various get employer post scenarios", func() {
			// Login other users for this test
			var marketingToken, regularToken string
			var wg sync.WaitGroup
			wg.Add(2)

			employerSigninAsync(
				"0026-employerposts2.example.com",
				"marketing@0026-employerposts2.example.com",
				"NewPassword123$",
				&marketingToken,
				&wg,
			)

			employerSigninAsync(
				"0026-employerposts2.example.com",
				"regular@0026-employerposts2.example.com",
				"NewPassword123$",
				&regularToken,
				&wg,
			)

			wg.Wait()

			testCases := []getEmployerPostTestCase{
				{
					description: "without authentication",
					token:       "",
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "admin can get post",
					token:       adminToken,
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var post common.EmployerPost
						err := json.Unmarshal(respBody, &post)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(post.ID).Should(Equal(postID))
						Expect(post.Content).Should(Equal("Test post for get"))
						Expect(
							post.CompanyDomain,
						).Should(Equal("0026-employerposts2.example.com"))
						Expect(
							post.Tags,
						).Should(ContainElements("engineering", "golang"))
						Expect(post.CreatedAt).ShouldNot(BeZero())
						Expect(post.UpdatedAt).ShouldNot(BeZero())
					},
				},
				{
					description: "marketing user can get post",
					token:       marketingToken,
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var post common.EmployerPost
						err := json.Unmarshal(respBody, &post)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(post.ID).Should(Equal(postID))
						Expect(post.Content).Should(Equal("Test post for get"))
						Expect(
							post.Tags,
						).Should(ContainElements("engineering", "golang"))
					},
				},
				{
					description: "regular user cannot get post",
					token:       regularToken,
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "non-existent post",
					token:       adminToken,
					request: employer.GetEmployerPostRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing GetEmployerPost: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/get-post",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("List Employer Posts", func() {
		var (
			adminToken string
			postIDs    []string
		)

		BeforeEach(func() {
			// Login admin user
			var wg sync.WaitGroup
			wg.Add(1)

			employerSigninAsync(
				"0026-employerposts3.example.com",
				"admin@0026-employerposts3.example.com",
				"NewPassword123$",
				&adminToken,
				&wg,
			)

			wg.Wait()

			// Create fresh posts for each test
			postIDs = createTestPosts(adminToken, 4)
		})

		type listEmployerPostsTestCase struct {
			description string
			token       string
			request     employer.ListEmployerPostsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various list employer posts scenarios", func() {
			// Login other users for this test
			var marketingToken, regularToken string
			var wg sync.WaitGroup
			wg.Add(2)

			employerSigninAsync(
				"0026-employerposts3.example.com",
				"marketing@0026-employerposts3.example.com",
				"NewPassword123$",
				&marketingToken,
				&wg,
			)

			employerSigninAsync(
				"0026-employerposts3.example.com",
				"regular@0026-employerposts3.example.com",
				"NewPassword123$",
				&regularToken,
				&wg,
			)

			wg.Wait()

			testCases := []listEmployerPostsTestCase{
				{
					description: "without authentication",
					token:       "",
					request:     employer.ListEmployerPostsRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request:     employer.ListEmployerPostsRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "admin can list posts with default limit",
					token:       adminToken,
					request:     employer.ListEmployerPostsRequest{},
					wantStatus:  http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.ListEmployerPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(4))
						Expect(resp.PaginationKey).Should(BeEmpty())

						// Verify posts are ordered by updated_at DESC
						Expect(resp.Posts[0].ID).Should(Equal(postIDs[3]))
						Expect(resp.Posts[1].ID).Should(Equal(postIDs[2]))
						Expect(resp.Posts[2].ID).Should(Equal(postIDs[1]))
						Expect(resp.Posts[3].ID).Should(Equal(postIDs[0]))
					},
				},
				{
					description: "marketing user can list posts with custom limit",
					token:       marketingToken,
					request: employer.ListEmployerPostsRequest{
						Limit: 2,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.ListEmployerPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(2))
						Expect(resp.PaginationKey).ShouldNot(BeEmpty())

						// Verify posts are ordered by updated_at DESC
						Expect(resp.Posts[0].ID).Should(Equal(postIDs[3]))
						Expect(resp.Posts[1].ID).Should(Equal(postIDs[2]))
					},
				},
				{
					description: "pagination works correctly",
					token:       adminToken,
					request: employer.ListEmployerPostsRequest{
						Limit:         2,
						PaginationKey: postIDs[2], // Start after the third newest post
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.ListEmployerPostsResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Posts).Should(HaveLen(2))
						Expect(resp.Posts[0].ID).Should(Equal(postIDs[1]))
						Expect(resp.Posts[1].ID).Should(Equal(postIDs[0]))
					},
				},
				{
					description: "regular user cannot list posts",
					token:       regularToken,
					request:     employer.ListEmployerPostsRequest{},
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "invalid limit",
					token:       adminToken,
					request: employer.ListEmployerPostsRequest{
						Limit: 50, // Max is 40
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing ListEmployerPosts: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/list-posts",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Delete Employer Post", func() {
		var (
			adminToken string
			postID     string
		)

		BeforeEach(func() {
			// Login admin user
			var wg sync.WaitGroup
			wg.Add(1)

			employerSigninAsync(
				"0026-employerposts4.example.com",
				"admin@0026-employerposts4.example.com",
				"NewPassword123$",
				&adminToken,
				&wg,
			)

			wg.Wait()

			// Create a fresh post for each test
			postID = createTestPost(
				adminToken,
				"Test post for delete",
				nil,
				nil,
			)
		})

		type deleteEmployerPostTestCase struct {
			description string
			token       string
			request     employer.DeleteEmployerPostRequest
			wantStatus  int
		}

		It("should handle various delete employer post scenarios", func() {
			// Login other users for this test
			var marketingToken, regularToken string
			var wg sync.WaitGroup
			wg.Add(2)

			employerSigninAsync(
				"0026-employerposts4.example.com",
				"marketing@0026-employerposts4.example.com",
				"NewPassword123$",
				&marketingToken,
				&wg,
			)

			employerSigninAsync(
				"0026-employerposts4.example.com",
				"regular@0026-employerposts4.example.com",
				"NewPassword123$",
				&regularToken,
				&wg,
			)

			wg.Wait()

			testCases := []deleteEmployerPostTestCase{
				{
					description: "without authentication",
					token:       "",
					request: employer.DeleteEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: employer.DeleteEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "regular user cannot delete post",
					token:       regularToken,
					request: employer.DeleteEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "non-existent post",
					token:       adminToken,
					request: employer.DeleteEmployerPostRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "marketing user can delete post",
					token:       marketingToken,
					request: employer.DeleteEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing DeleteEmployerPost: %s\n",
					tc.description,
				)
				testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/delete-post",
					tc.wantStatus,
				)
			}
		})
	})
})
