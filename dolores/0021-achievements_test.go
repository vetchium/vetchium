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
	"github.com/vetchium/vetchium/typespec/employer"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Achievements", Ordered, func() {
	var db *pgxpool.Pool
	var addUserToken, listUserToken, deleteUserToken, flowUserToken, secondUserToken string
	var adminOrgToken string
	var flowAchievementUserHandle, employerViewAchievementUserHandle, nonexistentUserHandle common.Handle

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0021-achievements-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(5) // 5 hub users to sign in
		hubSigninAsync(
			"add-user@achievement-hub.example",
			"NewPassword123$",
			&addUserToken,
			&wg,
		)
		hubSigninAsync(
			"list-user@achievement-hub.example",
			"NewPassword123$",
			&listUserToken,
			&wg,
		)
		hubSigninAsync(
			"delete-user@achievement-hub.example",
			"NewPassword123$",
			&deleteUserToken,
			&wg,
		)
		hubSigninAsync(
			"flow-user@achievement-hub.example",
			"NewPassword123$",
			&flowUserToken,
			&wg,
		)
		hubSigninAsync(
			"second-user@achievement-hub.example",
			"NewPassword123$",
			&secondUserToken,
			&wg,
		)
		wg.Wait()

		// Login org user and get token
		wg.Add(1)
		employerSigninAsync(
			"achievement-employer.example",
			"admin@achievement-employer.example",
			"NewPassword123$",
			&adminOrgToken,
			&wg,
		)
		wg.Wait()

		flowAchievementUserHandle = "flow-achievement-user"
		employerViewAchievementUserHandle = "employer-view-achievement-user"
		nonexistentUserHandle = "nonexistent-user"
	})

	AfterAll(func() {
		seedDatabase(db, "0021-achievements-down.pgsql")
		db.Close()
	})

	Describe("Add Achievement", func() {
		type addAchievementTestCase struct {
			description string
			token       string
			request     hub.AddAchievementRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []addAchievementTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "Test Patent",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "Test Patent",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid patent achievement",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "Innovation in Cloud Computing",
						Description: strptr(
							"A patent for innovative cloud architecture",
						),
						URL: strptr(
							"https://patent.example.com/cloud-innovation",
						),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddAchievementResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.ID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add valid publication achievement",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "PUBLICATION",
						Title: "Research on Distributed Systems",
						Description: strptr(
							"Publication about distributed system optimization",
						),
						URL: strptr(
							"https://journal.example.com/distributed-systems",
						),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddAchievementResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.ID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add valid certification achievement",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "CERTIFICATION",
						Title: "Google Cloud Professional Architect",
						Description: strptr(
							"Professional certification for Google Cloud architecture",
						),
						URL: strptr(
							"https://google.example.com/certification",
						),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddAchievementResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.ID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add achievement with missing required type",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Title: "Missing Type Achievement",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add achievement with missing required title",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type: "PATENT",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add achievement with title too short",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "AB", // Less than 3 characters
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add achievement with title too long",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type: "PATENT",
						Title: strings.Repeat(
							"x",
							129,
						), // More than 128 characters
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add achievement with description too long",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "Test Patent",
						Description: strptr(
							strings.Repeat("x", 1025),
						), // More than 1024 characters
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add achievement with URL too long",
					token:       addUserToken,
					request: hub.AddAchievementRequest{
						Type:  "PATENT",
						Title: "Test Patent",
						URL: strptr(
							strings.Repeat("x", 1025),
						), // More than 1024 characters
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/add-achievement",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("List Achievements", func() {
		type listAchievementsTestCase struct {
			description string
			token       string
			request     hub.ListAchievementsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle list achievements cases correctly", func() {
			testCases := []listAchievementsTestCase{
				{
					description: "list achievements without authentication",
					token:       "",
					request: hub.ListAchievementsRequest{
						Type: "PATENT",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "list achievements with invalid token",
					token:       "invalid-token",
					request: hub.ListAchievementsRequest{
						Type: "PATENT",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "list own patent achievements",
					token:       listUserToken,
					request: hub.ListAchievementsRequest{
						Type: "PATENT",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var achievements []common.Achievement
						err := json.Unmarshal(resp, &achievements)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(achievements).Should(HaveLen(1))

						achievement := achievements[0]
						Expect(
							achievement.ID,
						).Should(Equal("12345678-0021-0021-0021-000000000010"))
						Expect(string(achievement.Type)).Should(Equal("PATENT"))
						Expect(
							achievement.Title,
						).Should(Equal("Machine Learning Patent"))
						Expect(
							*achievement.Description,
						).Should(Equal("A patent for innovative ML algorithms"))
						Expect(
							*achievement.URL,
						).Should(Equal("https://patent.example.com/ml-innovation"))
						Expect(achievement.At).ShouldNot(BeNil())
					},
				},
				{
					description: "list own publication achievements",
					token:       listUserToken,
					request: hub.ListAchievementsRequest{
						Type: "PUBLICATION",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var achievements []common.Achievement
						err := json.Unmarshal(resp, &achievements)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(achievements).Should(HaveLen(1))

						achievement := achievements[0]
						Expect(
							achievement.ID,
						).Should(Equal("12345678-0021-0021-0021-000000000011"))
						Expect(
							string(achievement.Type),
						).Should(Equal("PUBLICATION"))
						Expect(
							achievement.Title,
						).Should(Equal("Research on AI Ethics"))
					},
				},
				{
					description: "list own certification achievements (empty)",
					token:       listUserToken,
					request: hub.ListAchievementsRequest{
						Type: "CERTIFICATION",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var achievements []common.Achievement
						err := json.Unmarshal(resp, &achievements)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(achievements).Should(BeEmpty())
					},
				},
				{
					description: "list other user's patent achievements (should not include IDs)",
					token:       secondUserToken,
					request: hub.ListAchievementsRequest{
						Type:   "PATENT",
						Handle: &employerViewAchievementUserHandle,
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var achievements []common.Achievement
						err := json.Unmarshal(resp, &achievements)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(achievements).Should(HaveLen(1))

						// When viewed by another user, ID should be empty
						achievement := achievements[0]
						Expect(achievement.ID).Should(BeEmpty())
						Expect(string(achievement.Type)).Should(Equal("PATENT"))
						Expect(
							achievement.Title,
						).Should(Equal("Blockchain Security Patent"))
					},
				},
				{
					description: "list achievements for non-existent user",
					token:       listUserToken,
					request: hub.ListAchievementsRequest{
						Type:   "PATENT",
						Handle: &nonexistentUserHandle,
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/list-achievements",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Delete Achievement", func() {
		type deleteAchievementTestCase struct {
			description string
			token       string
			request     hub.DeleteAchievementRequest
			wantStatus  int
		}

		It("should handle delete achievement cases correctly", func() {
			// Use the pre-loaded achievement for the delete test user
			achievementID := "12345678-0021-0021-0021-000000000012" // AWS certification for delete user

			// Verify the achievement exists before deletion
			listResp := testPOSTGetResp(
				deleteUserToken,
				hub.ListAchievementsRequest{
					Type: "CERTIFICATION",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var achievementsBefore []common.Achievement
			err := json.Unmarshal(listResp.([]byte), &achievementsBefore)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(achievementsBefore).Should(HaveLen(1))
			Expect(achievementsBefore[0].ID).Should(Equal(achievementID))

			testCases := []deleteAchievementTestCase{
				{
					description: "delete achievement without authentication",
					token:       "",
					request: hub.DeleteAchievementRequest{
						ID: achievementID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete achievement with invalid token",
					token:       "invalid-token",
					request: hub.DeleteAchievementRequest{
						ID: achievementID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete with empty achievement ID",
					token:       deleteUserToken,
					request: hub.DeleteAchievementRequest{
						ID: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "delete with non-existent achievement ID",
					token:       deleteUserToken,
					request: hub.DeleteAchievementRequest{
						ID: "87654321-0021-0021-0021-000000000000",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete another user's achievement (should fail)",
					token:       deleteUserToken,
					request: hub.DeleteAchievementRequest{
						ID: "12345678-0021-0021-0021-000000000010", // Machine Learning Patent of list user
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete achievement successfully",
					token:       deleteUserToken,
					request: hub.DeleteAchievementRequest{
						ID: achievementID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete already deleted achievement",
					token:       deleteUserToken,
					request: hub.DeleteAchievementRequest{
						ID: achievementID,
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/delete-achievement",
					tc.wantStatus,
				)
			}

			// Verify the achievement was actually deleted
			listRespAfter := testPOSTGetResp(
				deleteUserToken,
				hub.ListAchievementsRequest{
					Type: "CERTIFICATION",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var achievementsAfter []common.Achievement
			err = json.Unmarshal(listRespAfter.([]byte), &achievementsAfter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(achievementsAfter).Should(BeEmpty())
		})
	})

	Describe("Complete Achievement Flow", func() {
		It("should handle add, list, and delete operations correctly", func() {
			// Step 1: Verify no initial achievements for the flow user
			initialListResp := testPOSTGetResp(
				flowUserToken,
				hub.ListAchievementsRequest{
					Type: "PATENT",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var initialAchievements []common.Achievement
			err := json.Unmarshal(
				initialListResp.([]byte),
				&initialAchievements,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(initialAchievements).Should(BeEmpty())

			// Step 2: Add first achievement (patent)
			addResp1 := testPOSTGetResp(
				flowUserToken,
				hub.AddAchievementRequest{
					Type:  "PATENT",
					Title: "Quantum Computing Algorithm",
					Description: strptr(
						"A patent for quantum computing optimization",
					),
					URL: strptr(
						"https://patent.example.com/quantum-computing",
					),
				},
				"/hub/add-achievement",
				http.StatusOK,
			)

			var addResponse1 hub.AddAchievementResponse
			err = json.Unmarshal(addResp1.([]byte), &addResponse1)
			Expect(err).ShouldNot(HaveOccurred())
			firstAchievementID := addResponse1.ID

			// Step 3: Add second achievement (publication)
			addResp2 := testPOSTGetResp(
				flowUserToken,
				hub.AddAchievementRequest{
					Type:  "PUBLICATION",
					Title: "Machine Learning Research",
					Description: strptr(
						"Research paper on advanced ML techniques",
					),
					URL: strptr(
						"https://journal.example.com/ml-research",
					),
				},
				"/hub/add-achievement",
				http.StatusOK,
			)

			var addResponse2 hub.AddAchievementResponse
			err = json.Unmarshal(addResp2.([]byte), &addResponse2)
			Expect(err).ShouldNot(HaveOccurred())
			secondAchievementID := addResponse2.ID

			// Step 4: List patent achievements
			listPatentResp := testPOSTGetResp(
				flowUserToken,
				hub.ListAchievementsRequest{
					Type: "PATENT",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var patentAchievements []common.Achievement
			err = json.Unmarshal(listPatentResp.([]byte), &patentAchievements)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(patentAchievements).Should(HaveLen(1))
			Expect(patentAchievements[0].ID).Should(Equal(firstAchievementID))
			Expect(
				patentAchievements[0].Title,
			).Should(Equal("Quantum Computing Algorithm"))

			// Step 5: List publication achievements
			listPublicationResp := testPOSTGetResp(
				flowUserToken,
				hub.ListAchievementsRequest{
					Type: "PUBLICATION",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var publicationAchievements []common.Achievement
			err = json.Unmarshal(
				listPublicationResp.([]byte),
				&publicationAchievements,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(publicationAchievements).Should(HaveLen(1))
			Expect(
				publicationAchievements[0].ID,
			).Should(Equal(secondAchievementID))
			Expect(
				publicationAchievements[0].Title,
			).Should(Equal("Machine Learning Research"))

			// Step 6: View flow user's achievements as another user
			otherUserViewResp := testPOSTGetResp(
				secondUserToken,
				hub.ListAchievementsRequest{
					Type:   "PATENT",
					Handle: &flowAchievementUserHandle,
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var otherUserView []common.Achievement
			err = json.Unmarshal(otherUserViewResp.([]byte), &otherUserView)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(otherUserView).Should(HaveLen(1))
			Expect(
				otherUserView[0].ID,
			).Should(BeEmpty())
			// No ID when viewed by another user
			Expect(
				otherUserView[0].Title,
			).Should(Equal("Quantum Computing Algorithm"))

			// Step 7: Delete first achievement
			testPOSTGetResp(
				flowUserToken,
				hub.DeleteAchievementRequest{
					ID: firstAchievementID,
				},
				"/hub/delete-achievement",
				http.StatusOK,
			)

			// Step 8: Verify first achievement is deleted
			listPatentAfterResp := testPOSTGetResp(
				flowUserToken,
				hub.ListAchievementsRequest{
					Type: "PATENT",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var patentAchievementsAfter []common.Achievement
			err = json.Unmarshal(
				listPatentAfterResp.([]byte),
				&patentAchievementsAfter,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(patentAchievementsAfter).Should(BeEmpty())

			// Step 9: Delete second achievement
			testPOSTGetResp(
				flowUserToken,
				hub.DeleteAchievementRequest{
					ID: secondAchievementID,
				},
				"/hub/delete-achievement",
				http.StatusOK,
			)

			// Step 10: Verify second achievement is deleted
			listPublicationAfterResp := testPOSTGetResp(
				flowUserToken,
				hub.ListAchievementsRequest{
					Type: "PUBLICATION",
				},
				"/hub/list-achievements",
				http.StatusOK,
			)

			var publicationAchievementsAfter []common.Achievement
			err = json.Unmarshal(
				listPublicationAfterResp.([]byte),
				&publicationAchievementsAfter,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(publicationAchievementsAfter).Should(BeEmpty())
		})
	})

	Describe("Employer List Hub User Achievements", func() {
		type listHubUserAchievementsTestCase struct {
			description string
			token       string
			request     employer.ListHubUserAchievementsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It(
			"should handle employer listing hub user achievements correctly",
			func() {
				testCases := []listHubUserAchievementsTestCase{
					{
						description: "without authentication",
						token:       "",
						request: employer.ListHubUserAchievementsRequest{
							Handle: employerViewAchievementUserHandle,
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "with invalid token",
						token:       "invalid-token",
						request: employer.ListHubUserAchievementsRequest{
							Handle: employerViewAchievementUserHandle,
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "with empty handle",
						token:       adminOrgToken,
						request: employer.ListHubUserAchievementsRequest{
							Handle: "",
						},
						wantStatus: http.StatusBadRequest,
					},
					{
						description: "with non-existent user handle",
						token:       adminOrgToken,
						request: employer.ListHubUserAchievementsRequest{
							Handle: nonexistentUserHandle,
						},
						wantStatus: http.StatusNotFound,
					},
					{
						description: "list achievements for a user with achievements",
						token:       adminOrgToken,
						request: employer.ListHubUserAchievementsRequest{
							Handle: employerViewAchievementUserHandle,
						},
						wantStatus: http.StatusOK,
						validate: func(resp []byte) {
							var achievements []common.Achievement
							err := json.Unmarshal(resp, &achievements)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(achievements).Should(HaveLen(2))

							// Verify both entries are present without IDs
							for _, achievement := range achievements {
								Expect(achievement.ID).Should(BeEmpty())
								Expect(achievement.Title).ShouldNot(BeEmpty())
							}

							// Verify specific types of achievements
							patentFound := false
							publicationFound := false
							for _, achievement := range achievements {
								if achievement.Type == "PATENT" {
									patentFound = true
									Expect(
										achievement.Title,
									).Should(Equal("Blockchain Security Patent"))
								}
								if achievement.Type == "PUBLICATION" {
									publicationFound = true
									Expect(
										achievement.Title,
									).Should(Equal("Research on Quantum Computing"))
								}
							}
							Expect(patentFound).Should(BeTrue())
							Expect(publicationFound).Should(BeTrue())
						},
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(
						GinkgoWriter,
						"### Testing: %s\n",
						tc.description,
					)
					resp := testPOSTGetResp(
						tc.token,
						tc.request,
						"/employer/list-hub-user-achievements",
						tc.wantStatus,
					)
					if tc.validate != nil && tc.wantStatus == http.StatusOK {
						tc.validate(resp.([]byte))
					}
				}
			},
		)
	})
})
