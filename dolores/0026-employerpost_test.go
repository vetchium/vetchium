package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("Employer Posts", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool

		// Token variables for all tests
		// Employer 1 (0026-employerposts.example.com) - used by Add Employer Post tests
		employer1AdminToken     string
		employer1MarketingToken string
		employer1RegularToken   string

		// Employer 2 (0026-employerposts2.example.com) - used by Get Employer Post tests
		employer2AdminToken     string
		employer2MarketingToken string
		employer2RegularToken   string

		// Employer 3 (0026-employerposts3.example.com) - used by List Employer Posts tests
		employer3AdminToken     string
		employer3MarketingToken string
		employer3RegularToken   string

		// Employer 4 (0026-employerposts4.example.com) - used by Delete Employer Post tests
		employer4AdminToken     string
		employer4MarketingToken string
		employer4RegularToken   string

		// Org follow employer (0026-orgfollow.example.com) - used by Follow/Unfollow and Hub tests
		orgFollowAdminToken string

		// Different employer (0026-hubtest-different.example.com) - used by Hub tests for different employer
		differentAdminToken string

		// Hub users
		hubUserToken1 string
		hubUserToken2 string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0026-employerpost-up.pgsql")

		// Sign in all users once at the beginning
		var wg sync.WaitGroup
		wg.Add(16) // Total number of signin operations

		// Employer 1 (0026-employerposts.example.com) users
		employerSigninAsync(
			"0026-employerposts.example.com",
			"admin@0026-employerposts.example.com",
			"NewPassword123$",
			&employer1AdminToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts.example.com",
			"marketing@0026-employerposts.example.com",
			"NewPassword123$",
			&employer1MarketingToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts.example.com",
			"regular@0026-employerposts.example.com",
			"NewPassword123$",
			&employer1RegularToken,
			&wg,
		)

		// Employer 2 (0026-employerposts2.example.com) users
		employerSigninAsync(
			"0026-employerposts2.example.com",
			"admin@0026-employerposts2.example.com",
			"NewPassword123$",
			&employer2AdminToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts2.example.com",
			"marketing@0026-employerposts2.example.com",
			"NewPassword123$",
			&employer2MarketingToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts2.example.com",
			"regular@0026-employerposts2.example.com",
			"NewPassword123$",
			&employer2RegularToken,
			&wg,
		)

		// Employer 3 (0026-employerposts3.example.com) users
		employerSigninAsync(
			"0026-employerposts3.example.com",
			"admin@0026-employerposts3.example.com",
			"NewPassword123$",
			&employer3AdminToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts3.example.com",
			"marketing@0026-employerposts3.example.com",
			"NewPassword123$",
			&employer3MarketingToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts3.example.com",
			"regular@0026-employerposts3.example.com",
			"NewPassword123$",
			&employer3RegularToken,
			&wg,
		)

		// Employer 4 (0026-employerposts4.example.com) users
		employerSigninAsync(
			"0026-employerposts4.example.com",
			"admin@0026-employerposts4.example.com",
			"NewPassword123$",
			&employer4AdminToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts4.example.com",
			"marketing@0026-employerposts4.example.com",
			"NewPassword123$",
			&employer4MarketingToken,
			&wg,
		)
		employerSigninAsync(
			"0026-employerposts4.example.com",
			"regular@0026-employerposts4.example.com",
			"NewPassword123$",
			&employer4RegularToken,
			&wg,
		)

		// Org follow employer (0026-orgfollow.example.com)
		employerSigninAsync(
			"0026-orgfollow.example.com",
			"admin@0026-orgfollow.example.com",
			"NewPassword123$",
			&orgFollowAdminToken,
			&wg,
		)

		// Different employer (0026-hubtest-different.example.com)
		employerSigninAsync(
			"0026-hubtest-different.example.com",
			"admin@0026-hubtest-different.example.com",
			"NewPassword123$",
			&differentAdminToken,
			&wg,
		)

		// Hub users
		hubSigninAsync(
			"test1@0026-hubuser.example.com",
			"NewPassword123$",
			&hubUserToken1,
			&wg,
		)
		hubSigninAsync(
			"test2@0026-hubuser.example.com",
			"NewPassword123$",
			&hubUserToken2,
			&wg,
		)

		wg.Wait()
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
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "New post by admin",
						NewTags: []common.VTagName{"0026-admin-tag"},
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
					token:       employer1MarketingToken,
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
					token:       employer1RegularToken,
					request: employer.AddEmployerPostRequest{
						Content: "New post by regular user",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "empty content",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "too many tags",
					token:       employer1AdminToken,
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
				{
					description: "non-existent tag UUID",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with non-existent tag",
						TagIDs: []common.VTagID{
							common.VTagID(
								"12345678-0000-0000-0000-000000000000",
							),
						},
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "invalid UUID format for tag",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with invalid UUID format",
						TagIDs: []common.VTagID{
							common.VTagID("invalid-uuid-format"),
						},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "empty tag ID",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with empty tag ID",
						TagIDs: []common.VTagID{
							common.VTagID(""),
						},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "mix of valid and invalid tag IDs",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with mixed tag IDs",
						TagIDs: []common.VTagID{
							common.VTagID(
								"12345678-0026-0026-0026-000000050001",
							), // Valid tag
							common.VTagID(
								"12345678-0000-0000-0000-000000000000",
							), // Non-existent tag
						},
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "mix of existing and new tags",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with mixed tag types",
						TagIDs: []common.VTagID{
							common.VTagID(
								"12345678-0026-0026-0026-000000050001",
							), // Valid existing tag
						},
						NewTags: []common.VTagName{
							"0026-new-tag",
						},
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
					description: "create post with multiple new tags",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with multiple new tags",
						NewTags: []common.VTagName{
							"0026-multi-tag-1",
							"0026-multi-tag-2",
							"0026-multi-tag-3",
						},
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
					description: "create post with same tag name as existing (should reuse existing)",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with duplicate tag name",
						NewTags: []common.VTagName{
							"0026-engineering", // This tag already exists
						},
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
					description: "create post with mixed new tags (some existing, some new)",
					token:       employer1AdminToken,
					request: employer.AddEmployerPostRequest{
						Content: "Post with mixed new tag scenarios",
						NewTags: []common.VTagName{
							"0026-golang",    // Existing tag
							"0026-brand-new", // New tag
							"0026-rust",      // Existing tag
						},
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var resp employer.AddEmployerPostResponse
						err := json.Unmarshal(respBody, &resp)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.PostID).ShouldNot(BeEmpty())
					},
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

		It("should handle concurrent tag creation scenarios", func() {
			// Test concurrent creation of posts with the same tag names
			// This tests the tag creation logic under concurrent load
			var wg sync.WaitGroup
			numConcurrentPosts := 3 // Reduce concurrency to avoid overwhelming the test system
			postIDs := make([]string, numConcurrentPosts)
			errors := make([]error, numConcurrentPosts)
			resultMutex := sync.Mutex{} // Protect shared slices

			for i := 0; i < numConcurrentPosts; i++ {
				wg.Add(1)
				go func(index int) {
					defer GinkgoRecover() // Ensure Ginkgo can handle panics in goroutines
					defer wg.Done()

					// All goroutines try to create posts with the same tag name
					request := employer.AddEmployerPostRequest{
						Content: fmt.Sprintf("Concurrent post %d", index),
						NewTags: []common.VTagName{
							"0026-concurrent-tag", // Same tag name for all
							common.VTagName(fmt.Sprintf(
								"0026-unique-tag-%d",
								index,
							)), // Unique tag for each
						},
					}

					// Use a separate HTTP client for each goroutine to avoid connection sharing issues
					body, err := json.Marshal(request)
					if err != nil {
						resultMutex.Lock()
						errors[index] = err
						resultMutex.Unlock()
						return
					}

					req, err := http.NewRequest(
						http.MethodPost,
						serverURL+"/employer/add-post",
						bytes.NewBuffer(body),
					)
					if err != nil {
						resultMutex.Lock()
						errors[index] = err
						resultMutex.Unlock()
						return
					}

					req.Header.Set(
						"Authorization",
						"Bearer "+employer1AdminToken,
					)
					req.Header.Set("Content-Type", "application/json")

					// Create a separate client for this goroutine
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						resultMutex.Lock()
						errors[index] = err
						resultMutex.Unlock()
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						respBody, _ := io.ReadAll(resp.Body)
						resultMutex.Lock()
						errors[index] = fmt.Errorf(
							"unexpected status %d: %s",
							resp.StatusCode,
							string(respBody),
						)
						resultMutex.Unlock()
						return
					}

					respBody, err := io.ReadAll(resp.Body)
					if err != nil {
						resultMutex.Lock()
						errors[index] = err
						resultMutex.Unlock()
						return
					}

					var addResp employer.AddEmployerPostResponse
					err = json.Unmarshal(respBody, &addResp)
					if err != nil {
						resultMutex.Lock()
						errors[index] = err
						resultMutex.Unlock()
						return
					}

					resultMutex.Lock()
					postIDs[index] = addResp.PostID
					resultMutex.Unlock()
				}(i)
			}

			wg.Wait()

			// Verify all posts were created successfully
			for i, err := range errors {
				Expect(
					err,
				).ShouldNot(HaveOccurred(), "Post %d failed to create", i)
				Expect(
					postIDs[i],
				).ShouldNot(BeEmpty(), "Post %d has empty ID", i)
			}

			// Verify all posts can be retrieved
			for i, postID := range postIDs {
				getResp := testPOSTGetResp(
					employer1AdminToken,
					employer.GetEmployerPostRequest{PostID: postID},
					"/employer/get-post",
					http.StatusOK,
				)

				var post common.EmployerPost
				err := json.Unmarshal(getResp.([]byte), &post)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(
					post.Content,
				).Should(Equal(fmt.Sprintf("Concurrent post %d", i)))
				Expect(post.Tags).Should(ContainElement("0026-concurrent-tag"))
				Expect(
					post.Tags,
				).Should(ContainElement(fmt.Sprintf("0026-unique-tag-%d", i)))
			}
		})

		It("should handle rapid sequential tag creation", func() {
			// Test rapid sequential creation that might expose race conditions
			baseTags := []string{
				"0026-rapid-base-1",
				"0026-rapid-base-2",
				"0026-rapid-base-3",
			}

			var postIDs []string

			// Create posts rapidly in sequence
			for i := 0; i < 5; i++ { // Reduce from 10 to 5 to be less aggressive
				request := employer.AddEmployerPostRequest{
					Content: fmt.Sprintf("Rapid post %d", i),
					NewTags: []common.VTagName{
						// Mix of reused and new tags
						common.VTagName(
							baseTags[i%len(baseTags)],
						), // Reused tag
						common.VTagName(
							fmt.Sprintf("0026-rapid-unique-%d", i),
						), // New tag
					},
				}

				resp := testPOSTGetResp(
					employer1AdminToken,
					request,
					"/employer/add-post",
					http.StatusOK,
				)

				var addResp employer.AddEmployerPostResponse
				err := json.Unmarshal(resp.([]byte), &addResp)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(addResp.PostID).ShouldNot(BeEmpty())
				postIDs = append(postIDs, addResp.PostID)
			}

			// Verify all posts were created with correct tags
			for i, postID := range postIDs {
				getResp := testPOSTGetResp(
					employer1AdminToken,
					employer.GetEmployerPostRequest{PostID: postID},
					"/employer/get-post",
					http.StatusOK,
				)

				var post common.EmployerPost
				err := json.Unmarshal(getResp.([]byte), &post)
				Expect(err).ShouldNot(HaveOccurred())

				expectedBaseTag := baseTags[i%len(baseTags)]
				expectedUniqueTag := fmt.Sprintf("0026-rapid-unique-%d", i)

				Expect(post.Tags).Should(ContainElement(expectedBaseTag))
				Expect(post.Tags).Should(ContainElement(expectedUniqueTag))
			}
		})
	})

	Describe("Get Employer Post", func() {
		var (
			postID string
		)

		BeforeEach(func() {
			// Create a fresh post for each test
			postID = createTestPost(
				employer2AdminToken,
				"Test post for get",
				[]common.VTagID{
					common.VTagID(
						"12345678-0026-0026-0026-000000050001",
					), // 0026-engineering
					common.VTagID(
						"12345678-0026-0026-0026-000000050003",
					), // 0026-golang
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
					token:       employer2AdminToken,
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
							post.EmployerDomainName,
						).Should(Equal("0026-employerposts2.example.com"))
						Expect(
							post.Tags,
						).Should(ContainElements("0026-engineering", "0026-golang"))
						Expect(post.CreatedAt).ShouldNot(BeZero())
						Expect(post.UpdatedAt).ShouldNot(BeZero())
					},
				},
				{
					description: "marketing user can get post",
					token:       employer2MarketingToken,
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
						).Should(ContainElements("0026-engineering", "0026-golang"))
					},
				},
				{
					description: "regular user cannot get post",
					token:       employer2RegularToken,
					request: employer.GetEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "non-existent post",
					token:       employer2AdminToken,
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
			postIDs []string
		)

		BeforeEach(func() {
			// Create fresh posts for each test
			postIDs = createTestPosts(employer3AdminToken, 4)
		})

		type listEmployerPostsTestCase struct {
			description string
			token       string
			request     employer.ListEmployerPostsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various list employer posts scenarios", func() {

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
					token:       employer3AdminToken,
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
					token:       employer3MarketingToken,
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
					token:       employer3AdminToken,
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
					token:       employer3RegularToken,
					request:     employer.ListEmployerPostsRequest{},
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "invalid limit",
					token:       employer3AdminToken,
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
			postID string
		)

		BeforeEach(func() {
			// Create a fresh post for each test
			postID = createTestPost(
				employer4AdminToken,
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
					token:       employer4RegularToken,
					request: employer.DeleteEmployerPostRequest{
						PostID: postID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "non-existent post",
					token:       employer4AdminToken,
					request: employer.DeleteEmployerPostRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "marketing user can delete post",
					token:       employer4MarketingToken,
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

	Describe("Follow/Unfollow Organization", func() {

		It("should handle various follow org scenarios", func() {
			testCases := []struct {
				description string
				token       string
				request     hub.FollowOrgRequest
				wantStatus  int
			}{
				{
					description: "without authentication",
					token:       "",
					request: hub.FollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.FollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "follow non-existent org",
					token:       hubUserToken1,
					request: hub.FollowOrgRequest{
						Domain: "non-existent.example.com",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "follow org successfully",
					token:       hubUserToken1,
					request: hub.FollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "follow org again (should be idempotent)",
					token:       hubUserToken1,
					request: hub.FollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing FollowOrg: %s\n",
					tc.description,
				)
				testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/follow-org",
					tc.wantStatus,
				)
			}
		})

		It("should handle various unfollow org scenarios", func() {
			// First follow the org
			testPOSTGetResp(
				hubUserToken2,
				hub.FollowOrgRequest{
					Domain: "0026-orgfollow.example.com",
				},
				"/hub/follow-org",
				http.StatusOK,
			)

			testCases := []struct {
				description string
				token       string
				request     hub.UnfollowOrgRequest
				wantStatus  int
			}{
				{
					description: "without authentication",
					token:       "",
					request: hub.UnfollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.UnfollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "unfollow non-existent org",
					token:       hubUserToken2,
					request: hub.UnfollowOrgRequest{
						Domain: "non-existent.example.com",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "unfollow org successfully",
					token:       hubUserToken2,
					request: hub.UnfollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "unfollow org again (should be idempotent)",
					token:       hubUserToken2,
					request: hub.UnfollowOrgRequest{
						Domain: "0026-orgfollow.example.com",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing UnfollowOrg: %s\n",
					tc.description,
				)
				testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/unfollow-org",
					tc.wantStatus,
				)
			}
		})

		It("should show employer posts in timeline after following", func() {
			// First follow the org
			testPOSTGetResp(
				hubUserToken1,
				hub.FollowOrgRequest{
					Domain: "0026-orgfollow.example.com",
				},
				"/hub/follow-org",
				http.StatusOK,
			)

			// Create a post as org admin
			postContent := uuid.New().String()
			postID := createTestPost(orgFollowAdminToken, postContent, nil, nil)

			var timeline hub.MyHomeTimeline
			var found bool
			for i := 0; i < 5; i++ {
				// Wait for timeline to be refreshed
				<-time.After(30 * time.Second)
				resp := testPOSTGetResp(
					hubUserToken1,
					hub.GetMyHomeTimelineRequest{},
					"/hub/get-my-home-timeline",
					http.StatusOK,
				).([]byte)

				err := json.Unmarshal(resp, &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the post is in the timeline
				for _, post := range timeline.EmployerPosts {
					if post.ID == postID && post.Content == postContent {
						found = true
						break
					}
				}

				if found {
					break
				}
			}
			Expect(found).Should(BeTrue(), "not got post after following org")

			// Unfollow the org
			testPOSTGetResp(
				hubUserToken1,
				hub.UnfollowOrgRequest{
					Domain: "0026-orgfollow.example.com",
				},
				"/hub/unfollow-org",
				http.StatusOK,
			)

			// Create another post
			postID2 := createTestPost(
				orgFollowAdminToken,
				"after unfollow",
				nil,
				nil,
			)

			found = false
			for i := 0; i < 5; i++ {
				// Wait for timeline to be refreshed
				<-time.After(30 * time.Second)
				resp := testPOSTGetResp(
					hubUserToken1,
					hub.GetMyHomeTimelineRequest{},
					"/hub/get-my-home-timeline",
					http.StatusOK,
				).([]byte)

				err := json.Unmarshal(resp, &timeline)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the new post is in the timeline
				for _, post := range timeline.EmployerPosts {
					if post.ID == postID2 {
						found = true
						break
					}
				}

				if found {
					break
				}
			}

			Expect(found).Should(BeFalse(), "got post after unfollowing org")
		})
	})

	Describe("Get Employer Post Details (Hub)", func() {
		var (
			postID string
		)

		BeforeEach(func() {
			// Create a test post with tags
			postID = createTestPost(
				orgFollowAdminToken,
				"Test employer post for hub users to read",
				[]common.VTagID{
					common.VTagID(
						"12345678-0026-0026-0026-000000050001",
					), // 0026-engineering
				},
				[]common.VTagName{"0026-hub-test"},
			)
		})

		type getEmployerPostDetailsTestCase struct {
			description string
			token       string
			request     hub.GetEmployerPostDetailsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various get employer post details scenarios", func() {
			testCases := []getEmployerPostDetailsTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: postID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "hub user can get employer post details",
					token:       hubUserToken1,
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: postID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var post common.EmployerPost
						err := json.Unmarshal(respBody, &post)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(post.ID).Should(Equal(postID))
						Expect(
							post.Content,
						).Should(Equal("Test employer post for hub users to read"))
						Expect(post.EmployerName).ShouldNot(BeEmpty())
						Expect(
							post.EmployerDomainName,
						).Should(Equal("0026-orgfollow.example.com"))
						Expect(
							post.Tags,
						).Should(ContainElements("0026-engineering", "0026-hub-test"))
						Expect(post.CreatedAt).ShouldNot(BeZero())
						Expect(post.UpdatedAt).ShouldNot(BeZero())
					},
				},
				{
					description: "different hub user can also get employer post details",
					token:       hubUserToken2,
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: postID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						var post common.EmployerPost
						err := json.Unmarshal(respBody, &post)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(post.ID).Should(Equal(postID))
						Expect(
							post.Content,
						).Should(Equal("Test employer post for hub users to read"))
						Expect(post.EmployerName).ShouldNot(BeEmpty())
						Expect(
							post.EmployerDomainName,
						).Should(Equal("0026-orgfollow.example.com"))
						Expect(post.Tags).Should(HaveLen(2))
					},
				},
				{
					description: "non-existent post",
					token:       hubUserToken1,
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "empty post ID",
					token:       hubUserToken1,
					request: hub.GetEmployerPostDetailsRequest{
						EmployerPostID: "",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing GetEmployerPostDetails (Hub): %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/get-employer-post-details",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})

		It("should verify all response fields are populated correctly", func() {
			// Create a fresh post with no tags for this specific test
			postWithoutTags := createTestPost(
				orgFollowAdminToken,
				"Post without tags for field verification",
				nil,
				nil,
			)

			resp := testPOSTGetResp(
				hubUserToken1,
				hub.GetEmployerPostDetailsRequest{
					EmployerPostID: postWithoutTags,
				},
				"/hub/get-employer-post-details",
				http.StatusOK,
			).([]byte)

			var post common.EmployerPost
			err := json.Unmarshal(resp, &post)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify all required fields are present and have expected types
			Expect(post.ID).Should(Equal(postWithoutTags))
			Expect(
				post.Content,
			).Should(Equal("Post without tags for field verification"))
			Expect(post.EmployerName).ShouldNot(BeEmpty())
			Expect(
				post.EmployerDomainName,
			).Should(Equal("0026-orgfollow.example.com"))
			Expect(post.Tags).Should(BeEmpty()) // This post has no tags
			Expect(post.CreatedAt).ShouldNot(BeZero())
			Expect(post.UpdatedAt).ShouldNot(BeZero())

			// Verify timestamps are reasonable (within last minute)
			now := time.Now()
			Expect(post.CreatedAt).Should(BeTemporally("~", now, time.Minute))
			Expect(post.UpdatedAt).Should(BeTemporally("~", now, time.Minute))
			Expect(post.UpdatedAt).Should(BeTemporally(">=", post.CreatedAt))
		})

		It("should handle posts from different employers correctly", func() {
			// Create a post from different employer
			differentPostID := createTestPost(
				differentAdminToken,
				"Post from different employer",
				[]common.VTagID{
					common.VTagID(
						"12345678-0026-0026-0026-000000050002",
					), // 0026-marketing
				},
				nil,
			)

			// Hub user should be able to get posts from any employer
			resp := testPOSTGetResp(
				hubUserToken1,
				hub.GetEmployerPostDetailsRequest{
					EmployerPostID: differentPostID,
				},
				"/hub/get-employer-post-details",
				http.StatusOK,
			).([]byte)

			var post common.EmployerPost
			err := json.Unmarshal(resp, &post)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(post.ID).Should(Equal(differentPostID))
			Expect(post.Content).Should(Equal("Post from different employer"))
			Expect(
				post.EmployerDomainName,
			).Should(Equal("0026-hubtest-different.example.com"))
			Expect(post.Tags).Should(ContainElement("0026-marketing"))
		})
	})
})
