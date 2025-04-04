package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Work History", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2 string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0017-workhistory-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(2)
		hubSigninAsync(
			"user1@workhistory-hub.example",
			"NewPassword123$",
			&hubToken1,
			&wg,
		)
		hubSigninAsync(
			"user2@workhistory-hub.example",
			"NewPassword123$",
			&hubToken2,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0017-workhistory-down.pgsql")
		db.Close()
	})

	Describe("Add Work History", func() {
		type addWorkHistoryTestCase struct {
			description string
			token       string
			request     hub.AddWorkHistoryRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []addWorkHistoryTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "workhistory-employer1.example",
						Title:          "Software Developer",
						StartDate:      "2019-01-01",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "workhistory-employer1.example",
						Title:          "Software Developer",
						StartDate:      "2019-01-01",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid work history for onboarded employer",
					token:       hubToken2,
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "workhistory-employer1.example",
						Title:          "Software Developer",
						StartDate:      "2019-01-01",
						EndDate:        strptr("2019-12-31"),
						Description:    strptr("Worked on backend systems"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddWorkHistoryResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.WorkHistoryID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add valid work history for non-onboarded employer",
					token:       hubToken2,
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "non-onboarded-employer.example",
						Title:          "Software Developer",
						StartDate:      "2018-01-01",
						EndDate:        strptr("2018-12-31"),
						Description:    strptr("Worked at a startup"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddWorkHistoryResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.WorkHistoryID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add work history with invalid date format",
					token:       hubToken2,
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "workhistory-employer1.example",
						Title:          "Software Developer",
						StartDate:      "invalid-date",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add work history with end date before start date",
					token:       hubToken2,
					request: hub.AddWorkHistoryRequest{
						EmployerDomain: "workhistory-employer1.example",
						Title:          "Software Developer",
						StartDate:      "2019-01-01",
						EndDate:        strptr("2018-12-31"),
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/add-work-history",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("List Work History", func() {
		type listWorkHistoryTestCase struct {
			description string
			token       string
			request     hub.ListWorkHistoryRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []listWorkHistoryTestCase{
				{
					description: "without authentication",
					token:       "",
					request:     hub.ListWorkHistoryRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "list own work history",
					token:       hubToken1,
					request:     hub.ListWorkHistoryRequest{},
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var workHistory []hub.WorkHistory
						err := json.Unmarshal(resp, &workHistory)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(workHistory)).Should(Equal(2))
						Expect(
							workHistory[0].Title,
						).Should(Equal("Senior Engineer"))
						Expect(
							workHistory[1].Title,
						).Should(Equal("Software Engineer"))
					},
				},
				{
					description: "list another user's work history",
					token:       hubToken1,
					request: hub.ListWorkHistoryRequest{
						UserHandle: strptr("workhistory-user2"),
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var workHistory []hub.WorkHistory
						err := json.Unmarshal(resp, &workHistory)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							len(workHistory),
						).Should(Equal(2)) // From the previous add tests
					},
				},
				{
					description: "list non-existent user's work history",
					token:       hubToken1,
					request: hub.ListWorkHistoryRequest{
						UserHandle: strptr("nonexistent-user"),
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/list-work-history",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Update Work History", func() {
		type updateWorkHistoryTestCase struct {
			description string
			token       string
			request     hub.UpdateWorkHistoryRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []updateWorkHistoryTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.UpdateWorkHistoryRequest{
						ID:        "12345678-0017-0017-0017-000000000007",
						Title:     "Updated Title",
						StartDate: "2020-01-01",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "update own work history",
					token:       hubToken1,
					request: hub.UpdateWorkHistoryRequest{
						ID:          "12345678-0017-0017-0017-000000000007",
						Title:       "Updated Software Engineer",
						StartDate:   "2020-01-01",
						EndDate:     strptr("2021-12-31"),
						Description: strptr("Updated description"),
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "update another user's work history",
					token:       hubToken2,
					request: hub.UpdateWorkHistoryRequest{
						ID:        "12345678-0017-0017-0017-000000000007",
						Title:     "Malicious Update",
						StartDate: "2020-01-01",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "update non-existent work history",
					token:       hubToken1,
					request: hub.UpdateWorkHistoryRequest{
						ID:        "12345678-0017-0017-0017-000000000099",
						Title:     "Update Non-existent",
						StartDate: "2020-01-01",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "update with invalid date format",
					token:       hubToken1,
					request: hub.UpdateWorkHistoryRequest{
						ID:        "12345678-0017-0017-0017-000000000007",
						Title:     "Invalid Date",
						StartDate: "invalid-date",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "update with end date before start date",
					token:       hubToken1,
					request: hub.UpdateWorkHistoryRequest{
						ID:        "12345678-0017-0017-0017-000000000007",
						Title:     "Invalid Dates",
						StartDate: "2020-01-01",
						EndDate:   strptr("2019-12-31"),
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				reqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					http.MethodPut,
					serverURL+"/hub/update-work-history",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				if tc.token != "" {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))
			}
		})
	})

	Describe("Delete Work History", func() {
		type deleteWorkHistoryTestCase struct {
			description string
			token       string
			request     hub.DeleteWorkHistoryRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []deleteWorkHistoryTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.DeleteWorkHistoryRequest{
						ID: "12345678-0017-0017-0017-000000000008",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete another user's work history",
					token:       hubToken2,
					request: hub.DeleteWorkHistoryRequest{
						ID: "12345678-0017-0017-0017-000000000008",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete own work history",
					token:       hubToken1,
					request: hub.DeleteWorkHistoryRequest{
						ID: "12345678-0017-0017-0017-000000000008",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete non-existent work history",
					token:       hubToken1,
					request: hub.DeleteWorkHistoryRequest{
						ID: "12345678-0017-0017-0017-000000000099",
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/delete-work-history",
					tc.wantStatus,
				)
			}
		})
	})
})
