package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Hub Openings", Ordered, func() {
	var db *pgxpool.Pool
	var hubUserToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0008-hub-openings-up.pgsql")

		// Login as hub user
		loginReqBody, err := json.Marshal(vetchi.LoginRequest{
			Email:    "hubopening@hub.example",
			Password: "NewPassword123$",
		})
		Expect(err).ShouldNot(HaveOccurred())

		loginResp, err := http.Post(
			serverURL+"/hub/login",
			"application/json",
			bytes.NewBuffer(loginReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))

		var loginRespObj vetchi.LoginResponse
		err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
		Expect(err).ShouldNot(HaveOccurred())
		hubUserToken = loginRespObj.Token
	})

	AfterAll(func() {
		seedDatabase(db, "0008-hub-openings-down.pgsql")
		db.Close()
	})

	Describe("Find Hub Openings", func() {
		type findOpeningsTestCase struct {
			description string
			request     vetchi.FindHubOpeningsRequest
			wantStatus  int
			wantCount   int
			validate    func([]vetchi.HubOpening)
		}

		It("should find openings with various filters", func() {
			testCases := []findOpeningsTestCase{
				// Basic pagination and limit tests
				{
					description: "find all openings with default limit",
					request:     vetchi.FindHubOpeningsRequest{},
					wantStatus:  http.StatusOK,
					wantCount:   40, // Default limit
				},
				{
					description: "find openings with custom limit",
					request: vetchi.FindHubOpeningsRequest{
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					wantCount:  10,
				},
				{
					description: "find openings with invalid limit (too high)",
					request: vetchi.FindHubOpeningsRequest{
						Limit: 101,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "find openings with invalid limit (too low)",
					request: vetchi.FindHubOpeningsRequest{
						Limit: 0,
					},
					wantStatus: http.StatusBadRequest,
				},

				// Company domain filters
				{
					description: "find openings by single company domain",
					request: vetchi.FindHubOpeningsRequest{
						CompanyDomains: []string{"hubopening1.example"},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						for _, o := range openings {
							Expect(
								o.CompanyDomain,
							).Should(Equal("hubopening1.example"))
						}
					},
				},
				{
					description: "find openings by multiple company domains",
					request: vetchi.FindHubOpeningsRequest{
						CompanyDomains: []string{
							"hubopening1.example",
							"hubopening2.example",
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						for _, o := range openings {
							Expect(o.CompanyDomain).Should(Or(
								Equal("hubopening1.example"),
								Equal("hubopening2.example"),
							))
						}
					},
				},

				// Experience range filters
				{
					description: "find openings by experience range (entry level)",
					request: vetchi.FindHubOpeningsRequest{
						ExperienceRange: &vetchi.ExperienceRange{
							YoeMin: 0,
							YoeMax: 3,
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by experience range (mid level)",
					request: vetchi.FindHubOpeningsRequest{
						ExperienceRange: &vetchi.ExperienceRange{
							YoeMin: 3,
							YoeMax: 6,
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by experience range (senior level)",
					request: vetchi.FindHubOpeningsRequest{
						ExperienceRange: &vetchi.ExperienceRange{
							YoeMin: 6,
							YoeMax: 10,
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Salary range filters
				{
					description: "find openings by salary range (USD)",
					request: vetchi.FindHubOpeningsRequest{
						SalaryRange: &vetchi.SalaryRange{
							Currency: "USD",
							Min:      50000,
							Max:      100000,
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Location filters
				{
					description: "find openings by country",
					request: vetchi.FindHubOpeningsRequest{
						Countries: []vetchi.CountryCode{"IND"},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by specific location",
					request: vetchi.FindHubOpeningsRequest{
						Locations: []vetchi.LocationFilter{
							{
								CountryCode: "IND",
								City:        "Bangalore",
							},
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Education level filters
				{
					description: "find openings by minimum education level (Bachelor's)",
					request: vetchi.FindHubOpeningsRequest{
						MinEducationLevel: vetchi.BachelorEducation,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Remote work filters
				{
					description: "find openings by remote timezone",
					request: vetchi.FindHubOpeningsRequest{
						RemoteTimezones: []vetchi.TimeZone{
							"IST Indian Standard Time GMT+0530",
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by remote country",
					request: vetchi.FindHubOpeningsRequest{
						RemoteCountryCodes: []vetchi.CountryCode{"IND"},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Combined filters
				{
					description: "find openings with multiple filters",
					request: vetchi.FindHubOpeningsRequest{
						ExperienceRange: &vetchi.ExperienceRange{
							YoeMin: 3,
							YoeMax: 6,
						},
						SalaryRange: &vetchi.SalaryRange{
							Currency: "USD",
							Min:      80000,
							Max:      150000,
						},
						Countries:         []vetchi.CountryCode{"USA"},
						MinEducationLevel: vetchi.MasterEducation,
						RemoteTimezones: []vetchi.TimeZone{
							"PST Pacific Standard Time GMT-0800",
						},
					},
					wantStatus: http.StatusOK,
					validate: func(openings []vetchi.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				reqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/hub/find-openings",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubUserToken)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantStatus == http.StatusOK {
					var openings []vetchi.HubOpening
					err = json.NewDecoder(resp.Body).Decode(&openings)
					Expect(err).ShouldNot(HaveOccurred())

					if tc.wantCount > 0 {
						Expect(openings).Should(HaveLen(tc.wantCount))
					}

					if tc.validate != nil {
						tc.validate(openings)
					}
				}
			}
		})

		It("should handle pagination correctly", func() {
			// Get first page
			firstPageReq := vetchi.FindHubOpeningsRequest{
				Limit: 10,
			}
			firstPageBody, err := json.Marshal(firstPageReq)
			Expect(err).ShouldNot(HaveOccurred())

			firstPageResp, err := http.NewRequest(
				"POST",
				serverURL+"/hub/find-openings",
				bytes.NewBuffer(firstPageBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			firstPageResp.Header.Set("Authorization", "Bearer "+hubUserToken)

			resp, err := http.DefaultClient.Do(firstPageResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var firstPage []vetchi.HubOpening
			err = json.NewDecoder(resp.Body).Decode(&firstPage)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(firstPage).Should(HaveLen(10))

			// Get second page using pagination key
			secondPageReq := vetchi.FindHubOpeningsRequest{
				Limit:         10,
				PaginationKey: firstPage[len(firstPage)-1].PaginationKey,
			}
			secondPageBody, err := json.Marshal(secondPageReq)
			Expect(err).ShouldNot(HaveOccurred())

			secondPageResp, err := http.NewRequest(
				"POST",
				serverURL+"/hub/find-openings",
				bytes.NewBuffer(secondPageBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			secondPageResp.Header.Set("Authorization", "Bearer "+hubUserToken)

			resp, err = http.DefaultClient.Do(secondPageResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var secondPage []vetchi.HubOpening
			err = json.NewDecoder(resp.Body).Decode(&secondPage)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(secondPage).Should(HaveLen(10))

			// Verify no duplicates between pages
			firstPageIDs := make(map[string]bool)
			for _, o := range firstPage {
				firstPageIDs[o.OpeningIDWithinCompany] = true
			}
			for _, o := range secondPage {
				Expect(firstPageIDs[o.OpeningIDWithinCompany]).Should(BeFalse())
			}
		})
	})
})
