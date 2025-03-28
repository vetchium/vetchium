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
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = Describe("Education", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2, hubToken3, listToken, deleteToken, flowToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0020-education-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(6)
		hubSigninAsync(
			"user1@education-hub.example",
			"NewPassword123$",
			&hubToken1,
			&wg,
		)
		hubSigninAsync(
			"user2@education-hub.example",
			"NewPassword123$",
			&hubToken2,
			&wg,
		)
		hubSigninAsync(
			"user3@education-hub.example",
			"NewPassword123$",
			&hubToken3,
			&wg,
		)
		hubSigninAsync(
			"list-user@education-hub.example",
			"NewPassword123$",
			&listToken,
			&wg,
		)
		hubSigninAsync(
			"delete-user@education-hub.example",
			"NewPassword123$",
			&deleteToken,
			&wg,
		)
		hubSigninAsync(
			"flow-user@education-hub.example",
			"NewPassword123$",
			&flowToken,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0020-education-down.pgsql")
		db.Close()
	})

	Describe("Add Education", func() {
		type addEducationTestCase struct {
			description string
			token       string
			request     hub.AddEducationRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []addEducationTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "Bachelor of Science",
						StartDate:       strptr("2019-01-01"),
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "Bachelor of Science",
						StartDate:       strptr("2019-01-01"),
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid education entry",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "Bachelor of Science",
						StartDate:       strptr("2019-01-01"),
						EndDate:         strptr("2023-12-31"),
						Description:     strptr("Studied Computer Science"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddEducationResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.EducationID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add education with new institute domain",
					token:       hubToken2,
					request: hub.AddEducationRequest{
						InstituteDomain: "harvard.example",
						Degree:          "Master of Business Administration",
						StartDate:       strptr("2020-01-01"),
						EndDate:         strptr("2022-12-31"),
						Description:     strptr("Studied Business Management"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddEducationResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.EducationID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add education with invalid date format",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "PhD in Computer Science",
						StartDate:       strptr("invalid-date"),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add education with end date before start date",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "PhD in Computer Science",
						StartDate:       strptr("2022-01-01"),
						EndDate:         strptr("2021-12-31"),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add education with missing required institute domain",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						Degree:    "PhD in Computer Science",
						StartDate: strptr("2022-01-01"),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add education with degree exceeding max length",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          strings.Repeat("x", 65),
						StartDate:       strptr("2022-01-01"),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add education with description exceeding max length",
					token:       hubToken1,
					request: hub.AddEducationRequest{
						InstituteDomain: "stanford.example",
						Degree:          "PhD in Computer Science",
						StartDate:       strptr("2022-01-01"),
						Description: strptr(
							strings.Repeat("x", 1025),
						),
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/add-education",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("List Education", func() {
		type listEducationTestCase struct {
			description string
			token       string
			request     hub.ListEducationRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle list education cases correctly", func() {
			// These education entries are pre-loaded in the 0020-education-up.pgsql file
			// for the list-education-user user with ID 12345678-0020-0020-0020-000000000011
			user1EducId1 := "12345678-0020-0020-0020-000000000021" // MIT entry
			user1EducId2 := "12345678-0020-0020-0020-000000000022" // Caltech entry

			testCases := []listEducationTestCase{
				{
					description: "list education without authentication",
					token:       "",
					request:     hub.ListEducationRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "list education with invalid token",
					token:       "invalid-token",
					request:     hub.ListEducationRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "list own education (should include IDs)",
					token:       listToken,
					request:     hub.ListEducationRequest{},
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var educations []hub.Education
						err := json.Unmarshal(resp, &educations)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(educations).Should(HaveLen(2))

						// Verify both entries are present with IDs
						idFound1 := false
						idFound2 := false
						for _, edu := range educations {
							Expect(edu.ID).ShouldNot(BeEmpty())
							if edu.ID == user1EducId1 {
								idFound1 = true
								Expect(
									edu.InstituteDomain,
								).Should(Equal("mit.example"))
								Expect(
									edu.Degree,
								).Should(Equal("Bachelor of Engineering"))
								Expect(
									*edu.StartDate,
								).Should(Equal("2015-09-01"))
								Expect(*edu.EndDate).Should(Equal("2019-05-31"))
								Expect(
									*edu.Description,
								).Should(Equal("Electrical Engineering"))
							}
							if edu.ID == user1EducId2 {
								idFound2 = true
								Expect(
									edu.InstituteDomain,
								).Should(Equal("caltech.example"))
								Expect(
									edu.Degree,
								).Should(Equal("Master of Science"))
							}
						}
						Expect(idFound1).Should(BeTrue())
						Expect(idFound2).Should(BeTrue())
					},
				},
				{
					description: "list other user's education (should not include IDs)",
					token:       hubToken2,
					request: hub.ListEducationRequest{
						UserHandle: strptr("list-education-user"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var educations []hub.Education
						err := json.Unmarshal(resp, &educations)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(educations).Should(HaveLen(2))

						// Verify both entries are present without IDs
						for _, edu := range educations {
							Expect(edu.ID).Should(BeEmpty())
							Expect(edu.InstituteDomain).ShouldNot(BeEmpty())
							Expect(edu.Degree).ShouldNot(BeEmpty())
						}

						// Verify specific details
						foundMIT := false
						foundCaltech := false
						for _, edu := range educations {
							if edu.InstituteDomain == "mit.example" {
								foundMIT = true
							}
							if edu.InstituteDomain == "caltech.example" {
								foundCaltech = true
							}
						}
						Expect(foundMIT).Should(BeTrue())
						Expect(foundCaltech).Should(BeTrue())
					},
				},
				{
					description: "list education for non-existent user",
					token:       hubToken1,
					request: hub.ListEducationRequest{
						UserHandle: strptr("nonexistent-user"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var educations []hub.Education
						err := json.Unmarshal(resp, &educations)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(educations).Should(BeEmpty())
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/list-education",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Delete Education", func() {
		type deleteEducationTestCase struct {
			description string
			token       string
			request     hub.DeleteEducationRequest
			wantStatus  int
		}

		It("should handle delete education cases correctly", func() {
			// Use the pre-loaded education entry for delete-education-user
			educationID := "12345678-0020-0020-0020-000000000023" // Berkeley entry

			// Verify education was added
			listResp := testPOSTGetResp(
				deleteToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsBefore []hub.Education
			err := json.Unmarshal(listResp.([]byte), &educationsBefore)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify the added education exists
			found := false
			for _, edu := range educationsBefore {
				if edu.ID == educationID &&
					edu.InstituteDomain == "berkeley.example" {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			testCases := []deleteEducationTestCase{
				{
					description: "delete education without authentication",
					token:       "",
					request: hub.DeleteEducationRequest{
						EducationID: educationID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete education with invalid token",
					token:       "invalid-token",
					request: hub.DeleteEducationRequest{
						EducationID: educationID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete with improper UUID",
					token:       deleteToken,
					request: hub.DeleteEducationRequest{
						EducationID: "improper-uuid",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "delete with non-existent education ID",
					token:       deleteToken,
					request: hub.DeleteEducationRequest{
						EducationID: "87654321-0020-0020-0020-000000000000",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete education successfully",
					token:       deleteToken,
					request: hub.DeleteEducationRequest{
						EducationID: educationID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete already deleted education",
					token:       deleteToken,
					request: hub.DeleteEducationRequest{
						EducationID: educationID,
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete education of another user (should fail)",
					token:       hubToken2,
					request: hub.DeleteEducationRequest{
						EducationID: educationID,
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/delete-education",
					tc.wantStatus,
				)
			}

			// Verify the education was actually deleted by listing again
			listRespAfter := testPOSTGetResp(
				deleteToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsAfter []hub.Education
			err = json.Unmarshal(listRespAfter.([]byte), &educationsAfter)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify the deleted education no longer exists
			found = false
			for _, edu := range educationsAfter {
				if edu.ID == educationID {
					found = true
					break
				}
			}
			Expect(found).Should(BeFalse())
		})
	})

	Describe("Complete Education Flow", func() {
		It("should handle add, list, and delete operations correctly", func() {
			// Use the dedicated token for the flow test

			// Step 1: List education - should be empty initially
			listResp1 := testPOSTGetResp(
				flowToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var initialEducations []hub.Education
			err := json.Unmarshal(listResp1.([]byte), &initialEducations)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(initialEducations).Should(BeEmpty())

			// Step 2: Add first education
			addResp1 := testPOSTGetResp(
				flowToken,
				hub.AddEducationRequest{
					InstituteDomain: "princeton.example",
					Degree:          "Bachelor of Arts",
					StartDate:       strptr("2010-09-01"),
					EndDate:         strptr("2014-05-31"),
					Description:     strptr("English Literature"),
				},
				"/hub/add-education",
				http.StatusOK,
			)

			var addResponse1 hub.AddEducationResponse
			err = json.Unmarshal(addResp1.([]byte), &addResponse1)
			Expect(err).ShouldNot(HaveOccurred())
			firstEducationID := addResponse1.EducationID

			// Step 3: Add second education
			addResp2 := testPOSTGetResp(
				flowToken,
				hub.AddEducationRequest{
					InstituteDomain: "yale.example",
					Degree:          "Master of Fine Arts",
					StartDate:       strptr("2015-09-01"),
					EndDate:         strptr("2017-05-31"),
					Description:     strptr("Creative Writing"),
				},
				"/hub/add-education",
				http.StatusOK,
			)

			var addResponse2 hub.AddEducationResponse
			err = json.Unmarshal(addResp2.([]byte), &addResponse2)
			Expect(err).ShouldNot(HaveOccurred())
			secondEducationID := addResponse2.EducationID

			// Step 4: List education - should have two entries now
			listResp2 := testPOSTGetResp(
				flowToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsAfterAddition []hub.Education
			err = json.Unmarshal(listResp2.([]byte), &educationsAfterAddition)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(educationsAfterAddition).Should(HaveLen(2))

			// Step 5: Delete first education
			testPOSTGetResp(
				flowToken,
				hub.DeleteEducationRequest{
					EducationID: firstEducationID,
				},
				"/hub/delete-education",
				http.StatusOK,
			)

			// Step 6: List education - should have one entry now
			listResp3 := testPOSTGetResp(
				flowToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsAfterFirstDeletion []hub.Education
			err = json.Unmarshal(
				listResp3.([]byte),
				&educationsAfterFirstDeletion,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(educationsAfterFirstDeletion).Should(HaveLen(1))
			Expect(
				educationsAfterFirstDeletion[0].ID,
			).Should(Equal(secondEducationID))
			Expect(
				educationsAfterFirstDeletion[0].InstituteDomain,
			).Should(Equal("yale.example"))

			// Step 7: Delete second education
			testPOSTGetResp(
				flowToken,
				hub.DeleteEducationRequest{
					EducationID: secondEducationID,
				},
				"/hub/delete-education",
				http.StatusOK,
			)

			// Step 8: List education - should be empty again
			listResp4 := testPOSTGetResp(
				flowToken,
				hub.ListEducationRequest{},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsAfterSecondDeletion []hub.Education
			err = json.Unmarshal(
				listResp4.([]byte),
				&educationsAfterSecondDeletion,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(educationsAfterSecondDeletion).Should(BeEmpty())

			// Step 9: View education as another user
			// First, add new education for flow-education-user
			addResp3 := testPOSTGetResp(
				flowToken,
				hub.AddEducationRequest{
					InstituteDomain: "columbia.example",
					Degree:          "PhD in Psychology",
					StartDate:       strptr("2018-09-01"),
					EndDate:         strptr("2022-05-31"),
					Description:     strptr("Clinical Psychology"),
				},
				"/hub/add-education",
				http.StatusOK,
			)

			var addResponse3 hub.AddEducationResponse
			err = json.Unmarshal(addResp3.([]byte), &addResponse3)
			Expect(err).ShouldNot(HaveOccurred())

			// Check user3's education as user1 (should not include IDs)
			listRespByUser1 := testPOSTGetResp(
				hubToken1,
				hub.ListEducationRequest{
					UserHandle: strptr("flow-education-user"),
				},
				"/hub/list-education",
				http.StatusOK,
			)

			var educationsViewedByUser1 []hub.Education
			err = json.Unmarshal(
				listRespByUser1.([]byte),
				&educationsViewedByUser1,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(educationsViewedByUser1).Should(HaveLen(1))
			Expect(educationsViewedByUser1[0].ID).Should(BeEmpty())
			Expect(
				educationsViewedByUser1[0].InstituteDomain,
			).Should(Equal("columbia.example"))
			Expect(
				educationsViewedByUser1[0].Degree,
			).Should(Equal("PhD in Psychology"))
		})
	})

	Describe("Filter Institutes", func() {
		type filterInstitutesTestCase struct {
			description string
			token       string
			request     hub.FilterInstitutesRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle filter institutes cases correctly", func() {
			testCases := []filterInstitutesTestCase{
				{
					description: "filter institutes without authentication",
					token:       "",
					request: hub.FilterInstitutesRequest{
						Prefix: "stan",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "filter institutes with invalid token",
					token:       "invalid-token",
					request: hub.FilterInstitutesRequest{
						Prefix: "stan",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "filter with prefix too short",
					token:       hubToken1,
					request: hub.FilterInstitutesRequest{
						Prefix: "st",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "filter with prefix too long",
					token:       hubToken1,
					request: hub.FilterInstitutesRequest{
						Prefix: strings.Repeat("x", 65),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "filter with valid prefix (matching domains)",
					token:       hubToken1,
					request: hub.FilterInstitutesRequest{
						Prefix: "stan",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var institutes []hub.Institute
						err := json.Unmarshal(resp, &institutes)
						Expect(err).ShouldNot(HaveOccurred())
						foundStanford := false
						for _, inst := range institutes {
							if inst.Domain == "stanford.example" {
								foundStanford = true
								break
							}
						}
						Expect(foundStanford).Should(BeTrue())
					},
				},
				{
					description: "filter with valid prefix (matching institute names)",
					token:       hubToken2,
					request: hub.FilterInstitutesRequest{
						Prefix: "cal",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var institutes []hub.Institute
						err := json.Unmarshal(resp, &institutes)
						Expect(err).ShouldNot(HaveOccurred())
						foundCaltech := false
						for _, inst := range institutes {
							if inst.Domain == "caltech.example" {
								foundCaltech = true
								break
							}
						}
						Expect(foundCaltech).Should(BeTrue())
					},
				},
				{
					description: "filter with prefix that doesn't match any institute",
					token:       hubToken3,
					request: hub.FilterInstitutesRequest{
						Prefix: "xyz123",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var institutes []hub.Institute
						err := json.Unmarshal(resp, &institutes)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(institutes).Should(BeEmpty())
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/filter-institutes",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
