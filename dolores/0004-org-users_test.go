package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = Describe("Org Users", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOrgUsersToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0004-org-users-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@orgusers.example":        &adminToken,
			"crud@orgusers.example":         &crudToken,
			"viewer@orgusers.example":       &viewerToken,
			"non-orgusers@orgusers.example": &nonOrgUsersToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"orgusers.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0004-org-users-down.pgsql")
		db.Close()
	})

	Describe("OrgUsers Tests", func() {
		It("Add OrgUser", func() {
			type addOrgUserTestCase struct {
				description   string
				token         string
				request       employer.AddOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []addOrgUserTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "new1@orgusers.example",
						Name:  "New User One",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: employer.AddOrgUserRequest{
						Email: "new2@orgusers.example",
						Name:  "New User Two",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: employer.AddOrgUserRequest{
						Email: "new3@orgusers.example",
						Name:  "New User Three",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "with non-orgusers token",
					token:       nonOrgUsersToken,
					request: employer.AddOrgUserRequest{
						Email: "new4@orgusers.example",
						Name:  "New User Four",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "with duplicate email",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "new1@orgusers.example",
						Name:  "Duplicate User",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusConflict,
				},
				{
					description: "with invalid email",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "invalid-email",
						Name:  "Invalid Email User",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "with short name",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "short@orgusers.example",
						Name:  "ab",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with long name",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "long@orgusers.example",
						Name:  strings.Repeat("a", 257),
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with empty roles",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "noroles@orgusers.example",
						Name:  "No Roles User",
						Roles: []common.OrgUserRole{},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"roles"},
				},
				{
					description: "with invalid role",
					token:       adminToken,
					request: employer.AddOrgUserRequest{
						Email: "invalid@orgusers.example",
						Name:  "Invalid Role User",
						Roles: []common.OrgUserRole{"INVALID_ROLE"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"roles"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				if len(tc.wantErrFields) > 0 {
					validationErrors := testAddOrgUserGetResp(
						tc.token,
						tc.request,
						tc.wantStatus,
					)
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
				} else {
					testPOST(tc.token, tc.request, "/employer/add-org-user", tc.wantStatus)
				}
			}
		})

		It("Update OrgUser", func() {
			type updateOrgUserTestCase struct {
				description   string
				token         string
				request       employer.UpdateOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []updateOrgUserTestCase{
				{
					description: "with Admin token - update name",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Updated CRUD User",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Admin token - update roles",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Updated CRUD User",
						Roles: []common.OrgUserRole{
							"ORG_USERS_CRUD",
							"COST_CENTERS_CRUD",
						},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: employer.UpdateOrgUserRequest{
						Email: "viewer@orgusers.example",
						Name:  "Updated Viewer User",
						Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: employer.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Should Not Update",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "update non-existent user",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "nonexistent@orgusers.example",
						Name:  "Non-existent User",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "update last admin",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "admin@orgusers.example",
						Name:  "Last Admin",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid email format",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "invalid-email",
						Name:  "Invalid Email User",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "with short name",
					token:       adminToken,
					request: employer.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "ab",
						Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", tc.description)
				if len(tc.wantErrFields) > 0 {
					validationErrors := testUpdateOrgUserGetResp(
						tc.token,
						tc.request,
						tc.wantStatus,
					)
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
				} else {
					testPOST(tc.token, tc.request, "/employer/update-org-user", tc.wantStatus)
				}
			}
		})

		It("Disable OrgUser", func() {
			// First create some test users that we can disable
			testUsers := []employer.AddOrgUserRequest{
				{
					Email: "to-disable1@orgusers.example",
					Name:  "To Disable User 1",
					Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
				},
				{
					Email: "to-disable2@orgusers.example",
					Name:  "To Disable User 2",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
			}

			for _, user := range testUsers {
				testPOST(
					adminToken,
					user,
					"/employer/add-org-user",
					http.StatusOK,
				)
			}

			type disableOrgUserTestCase struct {
				description string
				token       string
				request     employer.DisableOrgUserRequest
				wantStatus  int
			}

			testCases := []disableOrgUserTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: employer.DisableOrgUserRequest{
						Email: "to-disable1@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: employer.DisableOrgUserRequest{
						Email: "to-disable2@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: employer.DisableOrgUserRequest{
						Email: "crud@orgusers.example",
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "disable last admin",
					token:       adminToken,
					request: employer.DisableOrgUserRequest{
						Email: "admin@orgusers.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "disable non-existent user",
					token:       adminToken,
					request: employer.DisableOrgUserRequest{
						Email: "nonexistent@orgusers.example",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "disable already disabled user",
					token:       adminToken,
					request: employer.DisableOrgUserRequest{
						Email: "to-disable1@orgusers.example",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "###### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/disable-org-user",
					tc.wantStatus,
				)
			}
		})

		It("Filter OrgUsers with Pagination", func() {
			// First create bulk test users
			bulkAddFilterOrgUsers(
				adminToken,
				"run-1",
				30,
				4,
			) // count not divisible by limit
			bulkAddFilterOrgUsers(
				adminToken,
				"run-2",
				32,
				4,
			) // count divisible by limit
			bulkAddFilterOrgUsers(
				adminToken,
				"run-3",
				2,
				4,
			) // count less than limit

			type filterOrgUsersTestCase struct {
				description   string
				token         string
				request       employer.FilterOrgUsersRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []filterOrgUsersTestCase{
				{
					description: "with Admin token - no filters",
					token:       adminToken,
					request:     employer.FilterOrgUsersRequest{},
					wantStatus:  http.StatusOK,
				},
				{
					description: "with CRUD token - prefix filter",
					token:       crudToken,
					request: employer.FilterOrgUsersRequest{
						Prefix: "admin",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - state filter",
					token:       viewerToken,
					request: employer.FilterOrgUsersRequest{
						State: []employer.OrgUserState{"DISABLED_ORG_USER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-orgusers token",
					token:       nonOrgUsersToken,
					request:     employer.FilterOrgUsersRequest{},
					wantStatus:  common.ErrEmployerRBAC,
				},
				{
					description: "with invalid limit",
					token:       adminToken,
					request: employer.FilterOrgUsersRequest{
						Limit: 41,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"limit"},
				},
				{
					description: "with invalid state",
					token:       adminToken,
					request: employer.FilterOrgUsersRequest{
						State: []employer.OrgUserState{"INVALID_STATE"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"state"},
				},
				{
					description: "with pagination",
					token:       adminToken,
					request: employer.FilterOrgUsersRequest{
						Limit:         2,
						PaginationKey: "crud@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", tc.description)
				if len(tc.wantErrFields) > 0 {
					validationErrors := testFilterOrgUsersGetResp(
						tc.token,
						tc.request,
						tc.wantStatus,
					)
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
				} else {
					testPOST(tc.token, tc.request, "/employer/filter-org-users", tc.wantStatus)
				}
			}
		})

		It("Filter OrgUsers with Prefix Search", func() {
			// First create test users with specific prefixes
			prefixUsers := []employer.AddOrgUserRequest{
				{
					Email: "alpha1@orgusers.example",
					Name:  "Beta User One",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "beta1@orgusers.example",
					Name:  "Alpha User One",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "gamma1@orgusers.example",
					Name:  "Alpha User Two",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
			}

			for _, user := range prefixUsers {
				testPOST(
					adminToken,
					user,
					"/employer/add-org-user",
					http.StatusOK,
				)
			}

			type prefixTestCase struct {
				description string
				prefix      string
				wantEmails  []string
			}

			testCases := []prefixTestCase{
				{
					description: "search by email prefix",
					prefix:      "alpha",
					wantEmails:  []string{"alpha1@orgusers.example"},
				},
				{
					description: "search by name prefix",
					prefix:      "Alpha",
					wantEmails: []string{
						"beta1@orgusers.example",
						"gamma1@orgusers.example",
					},
				},
				{
					description: "search with no matches",
					prefix:      "nonexistent",
					wantEmails:  []string{},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				resp := testPOSTGetResp(
					adminToken,
					employer.FilterOrgUsersRequest{Prefix: tc.prefix},
					"/employer/filter-org-users",
					http.StatusOK,
				).([]byte)

				var users []employer.OrgUser
				err := json.Unmarshal(resp, &users)
				Expect(err).ShouldNot(HaveOccurred())

				gotEmails := make([]string, len(users))
				for i, user := range users {
					gotEmails[i] = user.Email
				}

				// Sort both slices to ensure consistent comparison
				sort.Strings(gotEmails)
				sort.Strings(tc.wantEmails)

				Expect(gotEmails).Should(ContainElements(tc.wantEmails))
			}
		})

		It("Filter OrgUsers RBAC", func() {
			type rbacTestCase struct {
				description string
				token       string
				wantStatus  int
			}

			testCases := []rbacTestCase{
				{
					description: "with CRUD token",
					token:       crudToken,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with non-orgusers token",
					token:       nonOrgUsersToken,
					wantStatus:  common.ErrEmployerRBAC,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					employer.FilterOrgUsersRequest{Prefix: "alpha"},
					"/employer/filter-org-users",
					tc.wantStatus,
				)
			}
		})

		It("Enable OrgUser", func() {
			// First create and disable some test users that we can enable
			testUsers := []employer.AddOrgUserRequest{
				{
					Email: "to-enable1@orgusers.example",
					Name:  "To Enable User 1",
					Roles: []common.OrgUserRole{"ORG_USERS_CRUD"},
				},
				{
					Email: "to-enable2@orgusers.example",
					Name:  "To Enable User 2",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "active-user@orgusers.example",
					Name:  "Active User",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
			}

			// Delete any existing invite emails for these users
			for _, user := range testUsers {
				baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
				Expect(err).ShouldNot(HaveOccurred())
				query := url.Values{}
				query.Add(
					"query",
					fmt.Sprintf(
						"to:%s subject:Vetchium Employer Invitation",
						user.Email,
					),
				)
				baseURL.RawQuery = query.Encode()

				mailPitResp, err := http.Get(baseURL.String())
				Expect(err).ShouldNot(HaveOccurred())
				body, err := io.ReadAll(mailPitResp.Body)
				Expect(err).ShouldNot(HaveOccurred())
				mailPitResp.Body.Close()

				var mailPitRespObj MailPitResponse
				err = json.Unmarshal(body, &mailPitRespObj)
				Expect(err).ShouldNot(HaveOccurred())

				if len(mailPitRespObj.Messages) > 0 {
					messageIDs := make([]string, len(mailPitRespObj.Messages))
					for i, msg := range mailPitRespObj.Messages {
						messageIDs[i] = msg.ID
					}

					deleteReqBody, err := json.Marshal(
						MailPitDeleteRequest{IDs: messageIDs},
					)
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
				}
			}

			// Create the test users
			for _, user := range testUsers {
				testPOST(
					adminToken,
					user,
					"/employer/add-org-user",
					http.StatusOK,
				)
			}

			// Disable users that we'll enable in our tests
			disableUsers := []string{
				"to-enable1@orgusers.example",
				"to-enable2@orgusers.example",
			}

			for _, email := range disableUsers {
				testPOST(
					adminToken,
					employer.DisableOrgUserRequest{Email: email},
					"/employer/disable-org-user",
					http.StatusOK,
				)
			}

			type enableOrgUserTestCase struct {
				description string
				token       string
				request     employer.EnableOrgUserRequest
				wantStatus  int
			}

			testCases := []enableOrgUserTestCase{
				{
					description: "with Admin token - enable disabled user",
					token:       adminToken,
					request: employer.EnableOrgUserRequest{
						Email: "to-enable1@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token - enable disabled user",
					token:       crudToken,
					request: employer.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - should fail",
					token:       viewerToken,
					request: employer.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "with non-orgusers token - should fail",
					token:       nonOrgUsersToken,
					request: employer.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: common.ErrEmployerRBAC,
				},
				{
					description: "enable non-existent user",
					token:       adminToken,
					request: employer.EnableOrgUserRequest{
						Email: "nonexistent@orgusers.example",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "enable already active user",
					token:       adminToken,
					request: employer.EnableOrgUserRequest{
						Email: "active-user@orgusers.example",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "enable already enabled user",
					token:       adminToken,
					request: employer.EnableOrgUserRequest{
						Email: "to-enable1@orgusers.example",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "###### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/enable-org-user",
					tc.wantStatus,
				)
			}

			// Verify the enabled users are actually enabled by trying to filter them
			resp := testPOSTGetResp(
				adminToken,
				employer.FilterOrgUsersRequest{
					State:  []employer.OrgUserState{employer.AddedOrgUserState},
					Prefix: "to-enable",
				},
				"/employer/filter-org-users",
				http.StatusOK,
			).([]byte)

			var users []employer.OrgUser
			err := json.Unmarshal(resp, &users)
			Expect(err).ShouldNot(HaveOccurred())

			enabledEmails := []string{
				"to-enable1@orgusers.example",
				"to-enable2@orgusers.example",
			}
			for _, email := range enabledEmails {
				found := false
				for _, user := range users {
					if user.Email == email {
						found = true
						Expect(
							user.State,
						).Should(Equal(employer.AddedOrgUserState))
						break
					}
				}
				Expect(found).Should(BeTrue(), "Enabled but notfound %q", email)
			}

			// After running test cases, verify that invite emails were sent to enabled users

			// Wait for 10 seconds to ensure mailpit has processed the emails
			<-time.After(10 * time.Second)

			enabledUsers := []string{
				"to-enable1@orgusers.example",
				"to-enable2@orgusers.example",
			}

			for _, email := range enabledUsers {
				baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
				Expect(err).ShouldNot(HaveOccurred())
				query := url.Values{}
				query.Add(
					"query",
					fmt.Sprintf(
						"to:%s subject:Vetchium Employer Invitation",
						email,
					),
				)
				baseURL.RawQuery = query.Encode()

				var messageFound bool
				for i := 0; i < 3; i++ { // Retry up to 3 times
					<-time.After(
						15 * time.Second,
					) // Sleep at start of each iteration

					mailPitResp, err := http.Get(baseURL.String())
					Expect(err).ShouldNot(HaveOccurred())
					body, err := io.ReadAll(mailPitResp.Body)
					Expect(err).ShouldNot(HaveOccurred())
					mailPitResp.Body.Close()

					var mailPitRespObj MailPitResponse
					err = json.Unmarshal(body, &mailPitRespObj)
					Expect(err).ShouldNot(HaveOccurred())

					if len(mailPitRespObj.Messages) > 0 {
						messageFound = true
						// Clean up the message
						deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
							IDs: []string{mailPitRespObj.Messages[0].ID},
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
						Expect(
							deleteResp.StatusCode,
						).Should(Equal(http.StatusOK))
						break
					}
				}

				Expect(
					messageFound,
				).Should(BeTrue(), "No invite email found for %s", email)
			}
		})

		It("Test SignUp of OrgUsers", func() {
			type signupOrgUserTestCase struct {
				description   string
				request       employer.SignupOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			// First create an org user that we can use to get an invite token
			testPOST(
				adminToken,
				employer.AddOrgUserRequest{
					Email: "to-signup@orgusers.example",
					Name:  "To Signup User",
					Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				"/employer/add-org-user",
				http.StatusOK,
			)

			// Get the invite token from mailpit
			baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
			Expect(err).ShouldNot(HaveOccurred())
			query := url.Values{}
			// Subject comes from api/internal/hermione/orgusers/generate-invite.go
			// and potentially in future, from hedwig templates. Changes should
			// be synced in both places.
			query.Add(
				"query",
				"to:to-signup@orgusers.example subject:Vetchium Employer Invitation",
			)
			baseURL.RawQuery = query.Encode()

			finalURL := baseURL.String()
			fmt.Fprintf(GinkgoWriter, "finalURL: %s\n", finalURL)

			var messageID string
			// Retry a few times as email delivery might take time
			for i := 0; i < 3; i++ {
				<-time.After(15 * time.Second)
				fmt.Fprintf(GinkgoWriter, "Trying to get Invite mail\n")

				mailPitReq, err := http.NewRequest("GET", finalURL, nil)
				Expect(err).ShouldNot(HaveOccurred())
				mailPitReq.Header.Add("Content-Type", "application/json")

				mailPitResp, err := http.DefaultClient.Do(mailPitReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

				body, err := io.ReadAll(mailPitResp.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var mailPitRespObj MailPitResponse
				err = json.Unmarshal(body, &mailPitRespObj)
				Expect(err).ShouldNot(HaveOccurred())

				if len(mailPitRespObj.Messages) > 0 {
					Expect(len(mailPitRespObj.Messages)).Should(Equal(1))
					messageID = mailPitRespObj.Messages[0].ID
					break
				}
			}
			Expect(messageID).ShouldNot(BeEmpty())

			// Get the email content
			mailURL := fmt.Sprintf(
				"%s/api/v1/message/%s",
				mailPitURL,
				messageID,
			)
			mailPitReq, err := http.NewRequest("GET", mailURL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq.Header.Add("Content-Type", "application/json")

			mailPitResp, err := http.DefaultClient.Do(mailPitReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailPitResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			// Extract the invite token from the email body
			// The link format is: employer.example/signup-orguser/<token>
			re := regexp.MustCompile(`/signup-orguser/([a-zA-Z0-9]+)`)
			matches := re.FindStringSubmatch(string(body))
			Expect(len(matches)).Should(BeNumerically(">=", 2))
			inviteToken := matches[1]
			Expect(inviteToken).ShouldNot(BeEmpty())

			// Delete the email from mailpit
			mailPitDeleteReqBody, err := json.Marshal(MailPitDeleteRequest{
				IDs: []string{messageID},
			})
			Expect(err).ShouldNot(HaveOccurred())

			mailPitReq, err = http.NewRequest(
				"DELETE",
				mailPitURL+"/api/v1/messages",
				bytes.NewBuffer(mailPitDeleteReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq.Header.Set("Accept", "application/json")
			mailPitReq.Header.Add("Content-Type", "application/json")

			mailPitDeleteResp, err := http.DefaultClient.Do(mailPitReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitDeleteResp.StatusCode).Should(Equal(http.StatusOK))

			// Now continue with the test cases using the extracted invite token
			testCases := []signupOrgUserTestCase{
				{
					description: "valid signup",
					request: employer.SignupOrgUserRequest{
						InviteToken: inviteToken,
						Name:        "Signed Up User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-existent token",
					request: employer.SignupOrgUserRequest{
						InviteToken: "non-existent-token",
						Name:        "Invalid Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with already used token",
					request: employer.SignupOrgUserRequest{
						InviteToken: inviteToken, // Using the same token again
						Name:        "Duplicate Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with empty token",
					request: employer.SignupOrgUserRequest{
						InviteToken: "",
						Name:        "Empty Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"invite_token"},
				},
				{
					description: "with short name",
					request: employer.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        "ab", // Too short
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with long name",
					request: employer.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        strings.Repeat("a", 257), // Too long
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with empty password",
					request: employer.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        "Valid Name",
						Password:    "",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
				{
					description: "with short password",
					request: employer.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        "Valid Name",
						Password:    "short",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				if len(tc.wantErrFields) > 0 {
					validationErrors := testSignupOrgUserGetResp(
						tc.request,
						tc.wantStatus,
					)
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
				} else {
					// Note: Not using testPOST here since this endpoint doesn't require auth
					jsonBytes, err := json.Marshal(tc.request)
					Expect(err).ShouldNot(HaveOccurred())

					resp, err := http.Post(
						fmt.Sprintf("%s/employer/signup-orguser", serverURL),
						"application/json",
						bytes.NewBuffer(jsonBytes),
					)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).Should(Equal(tc.wantStatus))
				}
			}

			// Verify the successful signup by trying to signin
			sessionToken, err := employerSignin(
				"orgusers.example",
				"to-signup@orgusers.example",
				"ValidPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sessionToken).ShouldNot(BeEmpty())
		})
	})
})

func testAddOrgUserGetResp(
	token string,
	request employer.AddOrgUserRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/add-org-user", wantStatus).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func testUpdateOrgUserGetResp(
	token string,
	request employer.UpdateOrgUserRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/update-org-user", wantStatus).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func testFilterOrgUsersGetResp(
	token string,
	request employer.FilterOrgUsersRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/filter-org-users", wantStatus).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func bulkAddFilterOrgUsers(token string, runID string, count int, limit int) {
	wantUsers := []employer.OrgUser{}

	for i := 0; i < count; i++ {
		email := fmt.Sprintf("bulk-%s-%d@orgusers.example", runID, i)
		name := fmt.Sprintf("Bulk User %s-%d", runID, i)

		request := employer.AddOrgUserRequest{
			Email: email,
			Name:  name,
			Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
		}

		fmt.Fprintf(GinkgoWriter, "Adding user %q %q\n", email, name)
		testPOST(token, request, "/employer/add-org-user", http.StatusOK)

		wantUsers = append(wantUsers, employer.OrgUser{
			Email: email,
			Name:  name,
			Roles: []common.OrgUserRole{"ORG_USERS_VIEWER"},
			State: employer.AddedOrgUserState,
		})
	}

	paginationKey := ""
	gotUsers := []employer.OrgUser{}

	for {
		request := employer.FilterOrgUsersRequest{
			PaginationKey: paginationKey,
			Limit:         limit,
		}

		resp := testPOSTGetResp(token, request, "/employer/filter-org-users", http.StatusOK).([]byte)

		var users []employer.OrgUser
		err := json.Unmarshal(resp, &users)
		Expect(err).ShouldNot(HaveOccurred())

		if len(users) == 0 {
			break
		}

		gotUsers = append(gotUsers, users...)

		if len(users) < limit {
			break
		}

		paginationKey = users[len(users)-1].Email
	}

	// fmt.Fprintf(GinkgoWriter, "Got users: %v\n", gotUsers)

	// Verify the bulk created users are found
	for _, wantUser := range wantUsers {
		found := false
		for _, gotUser := range gotUsers {
			if gotUser.Email == wantUser.Email {
				Expect(gotUser.Name).Should(Equal(wantUser.Name))
				Expect(gotUser.Roles).Should(Equal(wantUser.Roles))
				Expect(gotUser.State).Should(Equal(wantUser.State))
				found = true
				break
			}
		}
		Expect(found).Should(BeTrue(), "User %s not found", wantUser.Email)
	}
}

// Helper function for signup validation error responses
func testSignupOrgUserGetResp(
	request employer.SignupOrgUserRequest,
	wantStatus int,
) common.ValidationErrors {
	jsonBytes, err := json.Marshal(request)
	Expect(err).ShouldNot(HaveOccurred())

	resp, err := http.Post(
		fmt.Sprintf("%s/employer/signup-orguser", serverURL),
		"application/json",
		bytes.NewBuffer(jsonBytes),
	)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(wantStatus))

	body, err := io.ReadAll(resp.Body)
	Expect(err).ShouldNot(HaveOccurred())
	defer resp.Body.Close()

	var validationErrors common.ValidationErrors
	err = json.Unmarshal(body, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}
