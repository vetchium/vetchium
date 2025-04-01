package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hub Openings", Ordered, func() {
	var db *pgxpool.Pool
	var hubUserToken string

	var bachelorEducation = common.BachelorEducation
	var masterEducation = common.MasterEducation

	var usaCountryCode = common.CountryCode("USA")

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0008-hub-openings-up.pgsql")

		// Login as hub user
		loginReqBody, err := json.Marshal(hub.LoginRequest{
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

		var loginRespObj hub.LoginResponse
		err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
		Expect(err).ShouldNot(HaveOccurred())
		tfaToken := loginRespObj.Token

		// Get TFA code from email
		baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
		Expect(err).ShouldNot(HaveOccurred())
		query := url.Values{}
		query.Add(
			"query",
			"to:hubopening@hub.example subject:Vetchi Two Factor Authentication",
		)
		baseURL.RawQuery = query.Encode()

		var messageID string
		for i := 0; i < 3; i++ {
			<-time.After(10 * time.Second)
			mailPitResp, err := http.Get(baseURL.String())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailPitResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			var mailPitRespObj MailPitResponse
			err = json.Unmarshal(body, &mailPitRespObj)
			Expect(err).ShouldNot(HaveOccurred())

			if len(mailPitRespObj.Messages) > 0 {
				messageID = mailPitRespObj.Messages[0].ID
				break
			}
		}
		Expect(messageID).ShouldNot(BeEmpty())

		// Get the email content
		mailResp, err := http.Get(mailPitURL + "/api/v1/message/" + messageID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

		body, err := io.ReadAll(mailResp.Body)
		Expect(err).ShouldNot(HaveOccurred())

		re := regexp.MustCompile(
			`Your Two Factor authentication code is:\s*([0-9]+)`,
		)
		matches := re.FindStringSubmatch(string(body))
		Expect(len(matches)).Should(BeNumerically(">=", 2))
		tfaCode := matches[1]

		// Clean up the email
		deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
			IDs: []string{messageID},
		})
		Expect(err).ShouldNot(HaveOccurred())

		req, err := http.NewRequest(
			"DELETE",
			mailPitURL+"/api/v1/messages",
			bytes.NewBuffer(deleteReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		req.Header.Set("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		deleteResp, err := http.DefaultClient.Do(req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(deleteResp.StatusCode).Should(Equal(http.StatusOK))

		// Complete TFA flow
		tfaReqBody, err := json.Marshal(hub.HubTFARequest{
			TFAToken:   tfaToken,
			TFACode:    tfaCode,
			RememberMe: false,
		})
		Expect(err).ShouldNot(HaveOccurred())

		tfaResp, err := http.Post(
			serverURL+"/hub/tfa",
			"application/json",
			bytes.NewBuffer(tfaReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tfaResp.StatusCode).Should(Equal(http.StatusOK))

		var tfaRespObj hub.HubTFAResponse
		err = json.NewDecoder(tfaResp.Body).Decode(&tfaRespObj)
		Expect(err).ShouldNot(HaveOccurred())
		hubUserToken = tfaRespObj.SessionToken
	})

	AfterAll(func() {
		seedDatabase(db, "0008-hub-openings-down.pgsql")
		db.Close()
	})

	Describe("Find Hub Openings", func() {
		type findOpeningsTestCase struct {
			description string
			request     hub.FindHubOpeningsRequest
			wantStatus  int
			wantCount   int
			validate    func([]hub.HubOpening)
		}

		It("should find openings with various filters", func() {
			testCases := []findOpeningsTestCase{
				// Basic pagination and limit tests
				{
					description: "find all openings with default limit",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					wantCount:  30,
				},
				{
					description: "find openings with custom limit",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
						Limit:       10,
					},
					wantStatus: http.StatusOK,
					wantCount:  10,
				},
				{
					description: "find openings with invalid limit (too high)",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
						Limit:       101,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "find openings with invalid limit (too low)",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
						Limit:       -1,
					},
					wantStatus: http.StatusBadRequest,
				},

				// Company domain filters
				{
					description: "find openings by single company domain",
					request: hub.FindHubOpeningsRequest{
						CompanyDomains: []string{"hubopening1.example"},
						CountryCode:    usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						for _, o := range openings {
							Expect(
								o.CompanyDomain,
							).Should(Equal("hubopening1.example"))
						}
					},
				},
				{
					description: "find openings by multiple company domains",
					request: hub.FindHubOpeningsRequest{
						CompanyDomains: []string{
							"hubopening1.example",
							"hubopening2.example",
						},
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
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
					request: hub.FindHubOpeningsRequest{
						ExperienceRange: &hub.ExperienceRange{
							YoeMin: 0,
							YoeMax: 3,
						},
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by experience range (mid level)",
					request: hub.FindHubOpeningsRequest{
						ExperienceRange: &hub.ExperienceRange{
							YoeMin: 3,
							YoeMax: 6,
						},
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "find openings by experience range (senior level)",
					request: hub.FindHubOpeningsRequest{
						ExperienceRange: &hub.ExperienceRange{
							YoeMin: 6,
							YoeMax: 10,
						},
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Salary range filters
				{
					description: "find openings by salary range (USD)",
					request: hub.FindHubOpeningsRequest{
						SalaryRange: &common.Salary{
							Currency:  "USD",
							MinAmount: 50000,
							MaxAmount: 100000,
						},
						CountryCode: usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Education level filters
				{
					description: "find openings by minimum education level (Bachelor's)",
					request: hub.FindHubOpeningsRequest{
						MinEducationLevel: &bachelorEducation,
						CountryCode:       usaCountryCode,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},

				// Combined filters
				{
					description: "find openings with multiple filters",
					request: hub.FindHubOpeningsRequest{
						ExperienceRange: &hub.ExperienceRange{
							YoeMin: 3,
							YoeMax: 6,
						},
						SalaryRange: &common.Salary{
							Currency:  "USD",
							MinAmount: 80000,
							MaxAmount: 150000,
						},
						CountryCode:       usaCountryCode,
						MinEducationLevel: &masterEducation,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
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

				fmt.Fprintf(
					GinkgoWriter,
					"Status, want: %d, got: %d\n",
					tc.wantStatus,
					resp.StatusCode,
				)
				if resp.StatusCode != tc.wantStatus &&
					resp.StatusCode == http.StatusBadRequest {
					body, err := io.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Fprintf(
						GinkgoWriter,
						"Validation Errors:\n%s\n",
						string(body),
					)
				}
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantStatus == http.StatusOK {
					var openings []hub.HubOpening
					err = json.NewDecoder(resp.Body).Decode(&openings)
					Expect(err).ShouldNot(HaveOccurred())

					if tc.wantCount > 0 {
						Expect(len(openings)).Should(Equal(tc.wantCount))
					}

					if tc.validate != nil {
						tc.validate(openings)
					}
				}
			}
		})

		It("should handle pagination correctly", func() {
			// Get first page with USA country filter
			firstPageReq := hub.FindHubOpeningsRequest{
				CountryCode: usaCountryCode,
				Limit:       10,
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

			var firstPage []hub.HubOpening
			err = json.NewDecoder(resp.Body).Decode(&firstPage)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(firstPage).Should(HaveLen(10))

			// Get second page using pagination key AND country code
			secondPageReq := hub.FindHubOpeningsRequest{
				CountryCode:   usaCountryCode,
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

			var secondPage []hub.HubOpening
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

		It("should handle location-based searches correctly", func() {
			testCases := []findOpeningsTestCase{
				{
					description: "find openings in different country",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
						Limit:       10,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "search for non-existent city in wrong country",
					request: hub.FindHubOpeningsRequest{
						CountryCode: usaCountryCode,
						Cities:      []string{"New York"},
						Limit:       10,
					},
					wantStatus: http.StatusOK,
					validate: func(openings []hub.HubOpening) {
						Expect(len(openings)).Should(Equal(0))
					},
				},
				// Add more test cases here
			}

			// Execute test cases
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

				fmt.Fprintf(
					GinkgoWriter,
					"Status, want: %d, got: %d\n",
					tc.wantStatus,
					resp.StatusCode,
				)
				if resp.StatusCode != tc.wantStatus &&
					resp.StatusCode == http.StatusBadRequest {
					body, err := io.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Fprintf(
						GinkgoWriter,
						"Validation Errors:\n%s\n",
						string(body),
					)
				}
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantStatus == http.StatusOK {
					var openings []hub.HubOpening
					err = json.NewDecoder(resp.Body).Decode(&openings)
					Expect(err).ShouldNot(HaveOccurred())

					if tc.wantCount > 0 {
						Expect(len(openings)).Should(Equal(tc.wantCount))
					}

					if tc.validate != nil {
						tc.validate(openings)
					}
				}
			}
		})
	})
})
