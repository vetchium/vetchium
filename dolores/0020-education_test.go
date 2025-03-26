package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = FDescribe("Education", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2 string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0020-education-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(2)
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
						Degree: "This is an extremely long degree name that exceeds the maximum " +
							"allowed length of 64 characters according to the API specification",
						StartDate: strptr("2022-01-01"),
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
							"This is an extremely long description that exceeds the maximum " +
								"allowed length of 1024 characters according to the API specification. " +
								"It contains a lot of unnecessary text just to make it longer and longer " +
								"until it reaches and exceeds the 1024 character limit. This text is being " +
								"repeated multiple times to ensure that it exceeds the limit. " +
								"This is an extremely long description that exceeds the maximum " +
								"allowed length of 1024 characters according to the API specification. " +
								"It contains a lot of unnecessary text just to make it longer and longer " +
								"until it reaches and exceeds the 1024 character limit. This text is being " +
								"repeated multiple times to ensure that it exceeds the limit. " +
								"This is an extremely long description that exceeds the maximum " +
								"allowed length of 1024 characters according to the API specification. " +
								"It contains a lot of unnecessary text just to make it longer and longer " +
								"until it reaches and exceeds the 1024 character limit. This text is being " +
								"repeated multiple times to ensure that it exceeds the limit.",
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
})
