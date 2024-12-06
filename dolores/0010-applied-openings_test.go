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
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Applied Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, viewerToken, nonAppToken string
	var hubUser1Token, hubUser2Token string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0010-applied-openings-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@applied1.example":   &adminToken,
			"viewer@applied1.example":  &viewerToken,
			"non-app@applied1.example": &nonAppToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"applied1.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()

		// Login hub users
		hubUser1Token = loginHubUser("hub1@applied1.example", "NewPassword123$")
		hubUser2Token = loginHubUser("hub2@applied1.example", "NewPassword123$")
	})

	AfterAll(func() {
		seedDatabase(db, "0010-applied-openings-down.pgsql")
		db.Close()
	})

	Describe("Employer Application Management", func() {
		Context("Color Tag Management", func() {
			It("should handle color tag operations", func() {
				type colorTagTestCase struct {
					description string
					token       string
					request     vetchi.SetApplicationColorTagRequest
					wantStatus  int
				}

				testCases := []colorTagTestCase{
					{
						description: "set color tag with admin token",
						token:       adminToken,
						request: vetchi.SetApplicationColorTagRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-1",
							ColorTag:      vetchi.GreenApplicationColorTag,
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "set color tag with viewer token",
						token:       viewerToken,
						request: vetchi.SetApplicationColorTagRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-2",
							ColorTag:      vetchi.YellowApplicationColorTag,
						},
						wantStatus: http.StatusForbidden,
					},
					{
						description: "set color tag for non-existent application",
						token:       adminToken,
						request: vetchi.SetApplicationColorTagRequest{
							ApplicationID: "non-existent",
							ColorTag:      vetchi.RedApplicationColorTag,
						},
						wantStatus: http.StatusNotFound,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					testPOST(
						tc.token,
						tc.request,
						"/employer/set-application-color-tag",
						tc.wantStatus,
					)
				}

				// Test remove color tag
				type removeColorTagTestCase struct {
					description string
					token       string
					request     vetchi.RemoveApplicationColorTagRequest
					wantStatus  int
				}

				removeTestCases := []removeColorTagTestCase{
					{
						description: "remove color tag with admin token",
						token:       adminToken,
						request: vetchi.RemoveApplicationColorTagRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-1",
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "remove color tag with viewer token",
						token:       viewerToken,
						request: vetchi.RemoveApplicationColorTagRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-2",
						},
						wantStatus: http.StatusForbidden,
					},
					{
						description: "remove color tag for non-existent application",
						token:       adminToken,
						request: vetchi.RemoveApplicationColorTagRequest{
							ApplicationID: "non-existent",
						},
						wantStatus: http.StatusNotFound,
					},
				}

				for _, tc := range removeTestCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					testPOST(
						tc.token,
						tc.request,
						"/employer/remove-application-color-tag",
						tc.wantStatus,
					)
				}
			})
		})

		Context("Application State Changes", func() {
			It("should handle shortlisting applications", func() {
				type shortlistTestCase struct {
					description string
					token       string
					request     vetchi.ShortlistApplicationRequest
					wantStatus  int
				}

				testCases := []shortlistTestCase{
					{
						description: "shortlist with admin token",
						token:       adminToken,
						request: vetchi.ShortlistApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-3",
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "shortlist with viewer token",
						token:       viewerToken,
						request: vetchi.ShortlistApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-4",
						},
						wantStatus: http.StatusForbidden,
					},
					{
						description: "shortlist non-existent application",
						token:       adminToken,
						request: vetchi.ShortlistApplicationRequest{
							ApplicationID: "non-existent",
						},
						wantStatus: http.StatusNotFound,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					testPOST(
						tc.token,
						tc.request,
						"/employer/shortlist-application",
						tc.wantStatus,
					)
				}
			})

			It("should handle rejecting applications", func() {
				type rejectTestCase struct {
					description string
					token       string
					request     vetchi.RejectApplicationRequest
					wantStatus  int
				}

				testCases := []rejectTestCase{
					{
						description: "reject with admin token",
						token:       adminToken,
						request: vetchi.RejectApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-5",
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "reject with viewer token",
						token:       viewerToken,
						request: vetchi.RejectApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-6",
						},
						wantStatus: http.StatusForbidden,
					},
					{
						description: "reject non-existent application",
						token:       adminToken,
						request: vetchi.RejectApplicationRequest{
							ApplicationID: "non-existent",
						},
						wantStatus: http.StatusNotFound,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					testPOST(
						tc.token,
						tc.request,
						"/employer/reject-application",
						tc.wantStatus,
					)
				}
			})
		})
	})

	Describe("Hub User Application Management", func() {
		Context("View Applications", func() {
			It("should handle my applications requests", func() {
				type myAppsTestCase struct {
					description string
					token       string
					request     vetchi.MyApplicationsRequest
					wantStatus  int
					validate    func([]vetchi.HubApplication)
				}

				testCases := []myAppsTestCase{
					{
						description: "get all applications",
						token:       hubUser1Token,
						request: vetchi.MyApplicationsRequest{
							Limit: 10,
						},
						wantStatus: http.StatusOK,
						validate: func(apps []vetchi.HubApplication) {
							Expect(len(apps)).Should(BeNumerically(">", 0))
						},
					},
					{
						description: "get applications with pagination",
						token:       hubUser1Token,
						request: vetchi.MyApplicationsRequest{
							Limit: 5,
						},
						wantStatus: http.StatusOK,
						validate: func(apps []vetchi.HubApplication) {
							Expect(len(apps)).Should(Equal(5))
						},
					},
					{
						description: "get applications with state filter",
						token:       hubUser1Token,
						request: vetchi.MyApplicationsRequest{
							State: vetchi.AppliedAppState,
							Limit: 10,
						},
						wantStatus: http.StatusOK,
						validate: func(apps []vetchi.HubApplication) {
							for _, app := range apps {
								Expect(
									app.State,
								).Should(Equal(vetchi.AppliedAppState))
							}
						},
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					resp := testPOSTGetResp(
						tc.token,
						tc.request,
						"/hub/my-applications",
						tc.wantStatus,
					).([]byte)

					if tc.wantStatus == http.StatusOK {
						var apps []vetchi.HubApplication
						err := json.Unmarshal(resp, &apps)
						Expect(err).ShouldNot(HaveOccurred())
						tc.validate(apps)
					}
				}
			})
		})

		Context("Withdraw Applications", func() {
			It("should handle withdrawing applications", func() {
				type withdrawTestCase struct {
					description string
					token       string
					request     vetchi.WithdrawApplicationRequest
					wantStatus  int
				}

				testCases := []withdrawTestCase{
					{
						description: "withdraw own application",
						token:       hubUser1Token,
						request: vetchi.WithdrawApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-7",
						},
						wantStatus: http.StatusOK,
					},
					{
						description: "withdraw another user's application",
						token:       hubUser2Token,
						request: vetchi.WithdrawApplicationRequest{
							ApplicationID: "APP-12345678-0010-0010-0010-000000000201-8",
						},
						wantStatus: http.StatusNotFound,
					},
					{
						description: "withdraw non-existent application",
						token:       hubUser1Token,
						request: vetchi.WithdrawApplicationRequest{
							ApplicationID: "non-existent",
						},
						wantStatus: http.StatusNotFound,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
					testPOST(
						tc.token,
						tc.request,
						"/hub/withdraw-application",
						tc.wantStatus,
					)
				}
			})
		})
	})
})

// Helper function to login a hub user and get their token
func loginHubUser(email, password string) string {
	loginReqBody, err := json.Marshal(vetchi.LoginRequest{
		Email:    vetchi.EmailAddress(email),
		Password: vetchi.Password(password),
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

	tfaCode, messageID := getTFACode(email)
	sessionToken := getSessionToken(loginRespObj.Token, tfaCode, false)
	cleanupEmail(messageID)

	return sessionToken
}
