package dolores

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Incognito Posts API", Ordered, func() {
	var (
		pool *pgxpool.Pool

		// User tokens for dedicated test users
		user0034Token1 string
		user0034Token2 string
		user0034Token3 string
		user0034Token4 string
		user0034Token5 string
		user0034Token6 string
		user0034Token7 string
		user0034Token8 string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0034-incognito-posts-up.pgsql")

		var wg sync.WaitGroup
		wg.Add(8)
		hubSigninAsync(
			"user0034-1@0034-test.com",
			"NewPassword123$",
			&user0034Token1,
			&wg,
		)
		hubSigninAsync(
			"user0034-2@0034-test.com",
			"NewPassword123$",
			&user0034Token2,
			&wg,
		)
		hubSigninAsync(
			"user0034-3@0034-test.com",
			"NewPassword123$",
			&user0034Token3,
			&wg,
		)
		hubSigninAsync(
			"user0034-4@0034-test.com",
			"NewPassword123$",
			&user0034Token4,
			&wg,
		)
		hubSigninAsync(
			"user0034-5@0034-test.com",
			"NewPassword123$",
			&user0034Token5,
			&wg,
		)
		hubSigninAsync(
			"user0034-6@0034-test.com",
			"NewPassword123$",
			&user0034Token6,
			&wg,
		)
		hubSigninAsync(
			"user0034-7@0034-test.com",
			"NewPassword123$",
			&user0034Token7,
			&wg,
		)
		hubSigninAsync(
			"user0034-8@0034-test.com",
			"NewPassword123$",
			&user0034Token8,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(pool, "0034-incognito-posts-down.pgsql")
		pool.Close()
	})

	Describe("AddIncognitoPost", func() {
		It("should create an incognito post successfully with valid data",
			func() {
				reqBody := hub.AddIncognitoPostRequest{
					Content: "This is a new incognito post about technology and career development.",
					TagIDs:  []common.VTagID{"technology"},
				}

				respData := testPOSTGetResp(
					user0034Token1,
					reqBody,
					"/hub/add-incognito-post",
					http.StatusOK,
				)

				var response hub.AddIncognitoPostResponse
				err := json.Unmarshal(respData.([]byte), &response)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.IncognitoPostID).ShouldNot(BeEmpty())

				getReq := hub.GetIncognitoPostRequest{
					IncognitoPostID: response.IncognitoPostID,
				}
				getResp := testPOSTGetResp(
					user0034Token1,
					getReq,
					"/hub/get-incognito-post",
					http.StatusOK,
				)

				var getResponse hub.IncognitoPost
				err = json.Unmarshal(getResp.([]byte), &getResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getResponse.Content).Should(Equal(reqBody.Content))
				Expect(getResponse.IsCreatedByMe).Should(BeTrue())
				Expect(getResponse.IsDeleted).Should(BeFalse())
				Expect(len(getResponse.Tags)).Should(Equal(1))
				Expect(
					getResponse.Tags[0].ID,
				).Should(Equal(common.VTagID("technology")))
			},
		)

		It("should create an incognito post with multiple tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "Post with multiple tags about personal development.",
				TagIDs: []common.VTagID{
					"personal-development",
					"careers",
					"mentorship",
				},
			}

			respData := testPOSTGetResp(
				user0034Token2,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var response hub.AddIncognitoPostResponse
			err := json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.IncognitoPostID).ShouldNot(BeEmpty())

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: response.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				user0034Token2,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Tags)).Should(Equal(3))
		})

		It("should fail without authentication", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This should fail without auth.",
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				"",
				reqBody,
				"/hub/add-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail with empty content", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "",
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				user0034Token3,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with content too long", func() {
			longContent := make([]byte, 1025)
			for i := range longContent {
				longContent[i] = 'a'
			}

			reqBody := hub.AddIncognitoPostRequest{
				Content: string(longContent),
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				user0034Token4,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with no tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has no tags which should fail.",
				TagIDs:  []common.VTagID{},
			}

			testPOST(
				user0034Token5,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with too many tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has too many tags.",
				TagIDs: []common.VTagID{
					"technology",
					"careers",
					"mentorship",
					"startups",
				},
			}

			testPOST(
				user0034Token6,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with invalid tag IDs", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has invalid tags.",
				TagIDs:  []common.VTagID{"nonexistent-tag"},
			}

			testPOST(
				user0034Token7,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})
	})

	Describe("GetIncognitoPost", func() {
		It("should retrieve an incognito post successfully", func() {
			createReq := hub.AddIncognitoPostRequest{
				Content: "Test post for retrieval.",
				TagIDs:  []common.VTagID{"technology"},
			}

			createResp := testPOSTGetResp(
				user0034Token1,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var createResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &createResponse)
			Expect(err).ShouldNot(HaveOccurred())

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			getResp := testPOSTGetResp(
				user0034Token2,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.Content).Should(Equal(createReq.Content))
			Expect(getResponse.IsCreatedByMe).Should(BeFalse())
			Expect(getResponse.CanUpvote).Should(BeTrue())
			Expect(getResponse.CanDownvote).Should(BeTrue())
		})

		It("should show ownership for creator", func() {
			createReq := hub.AddIncognitoPostRequest{
				Content: "Test post for ownership check.",
				TagIDs:  []common.VTagID{"technology"},
			}

			createResp := testPOSTGetResp(
				user0034Token3,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var createResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &createResponse)
			Expect(err).ShouldNot(HaveOccurred())

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			getResp := testPOSTGetResp(
				user0034Token3,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.IsCreatedByMe).Should(BeTrue())
			Expect(getResponse.CanUpvote).Should(BeFalse())
			Expect(getResponse.CanDownvote).Should(BeFalse())
		})

		It("should fail without authentication", func() {
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: "some-id",
			}

			testPOST(
				"",
				getReq,
				"/hub/get-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: "nonexistent-post-id",
			}

			testPOST(
				user0034Token4,
				getReq,
				"/hub/get-incognito-post",
				http.StatusNotFound,
			)
		})
	})

	Describe("DeleteIncognitoPost", func() {
		It("should delete own incognito post successfully", func() {
			createReq := hub.AddIncognitoPostRequest{
				Content: "Test post for deletion.",
				TagIDs:  []common.VTagID{"technology"},
			}

			createResp := testPOSTGetResp(
				user0034Token5,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var createResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &createResponse)
			Expect(err).ShouldNot(HaveOccurred())

			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			testPOST(
				user0034Token5,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusOK,
			)

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			testPOST(
				user0034Token5,
				getReq,
				"/hub/get-incognito-post",
				http.StatusNotFound,
			)
		})

		It("should fail to delete someone else's post", func() {
			createReq := hub.AddIncognitoPostRequest{
				Content: "Test post for forbidden deletion.",
				TagIDs:  []common.VTagID{"technology"},
			}

			createResp := testPOSTGetResp(
				user0034Token6,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var createResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &createResponse)
			Expect(err).ShouldNot(HaveOccurred())

			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			testPOST(
				user0034Token7,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusForbidden,
			)
		})

		It("should fail without authentication", func() {
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: "some-id",
			}

			testPOST(
				"",
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: "nonexistent-post-id",
			}

			testPOST(
				user0034Token8,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusNotFound,
			)
		})
	})
})
