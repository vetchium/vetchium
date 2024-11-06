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

	"github.com/psankar/vetchi/api/pkg/vetchi"

	"github.com/jackc/pgx/v5/pgxpool"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				request       vetchi.AddOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []addOrgUserTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "new1@orgusers.example",
						Name:  "New User One",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: vetchi.AddOrgUserRequest{
						Email: "new2@orgusers.example",
						Name:  "New User Two",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: vetchi.AddOrgUserRequest{
						Email: "new3@orgusers.example",
						Name:  "New User Three",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with non-orgusers token",
					token:       nonOrgUsersToken,
					request: vetchi.AddOrgUserRequest{
						Email: "new4@orgusers.example",
						Name:  "New User Four",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with duplicate email",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "new1@orgusers.example",
						Name:  "Duplicate User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusConflict,
				},
				{
					description: "with invalid email",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "invalid-email",
						Name:  "Invalid Email User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "with short name",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "short@orgusers.example",
						Name:  "ab",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with long name",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "long@orgusers.example",
						Name:  strings.Repeat("a", 257),
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with empty roles",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "noroles@orgusers.example",
						Name:  "No Roles User",
						Roles: []vetchi.OrgUserRole{},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"roles"},
				},
				{
					description: "with invalid role",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "invalid@orgusers.example",
						Name:  "Invalid Role User",
						Roles: []vetchi.OrgUserRole{"INVALID_ROLE"},
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
				request       vetchi.UpdateOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []updateOrgUserTestCase{
				{
					description: "with Admin token - update name",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Updated CRUD User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Admin token - update roles",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Updated CRUD User",
						Roles: []vetchi.OrgUserRole{
							"ORG_USERS_CRUD",
							"COST_CENTERS_CRUD",
						},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "viewer@orgusers.example",
						Name:  "Updated Viewer User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "Should Not Update",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "update non-existent user",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "nonexistent@orgusers.example",
						Name:  "Non-existent User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "update last admin",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "admin@orgusers.example",
						Name:  "Last Admin",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid email format",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "invalid-email",
						Name:  "Invalid Email User",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "with short name",
					token:       adminToken,
					request: vetchi.UpdateOrgUserRequest{
						Email: "crud@orgusers.example",
						Name:  "ab",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
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
			testUsers := []vetchi.AddOrgUserRequest{
				{
					Email: "to-disable1@orgusers.example",
					Name:  "To Disable User 1",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
				},
				{
					Email: "to-disable2@orgusers.example",
					Name:  "To Disable User 2",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
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
				request     vetchi.DisableOrgUserRequest
				wantStatus  int
			}

			testCases := []disableOrgUserTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "to-disable1@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "to-disable2@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "crud@orgusers.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "disable last admin",
					token:       adminToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "admin@orgusers.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "disable non-existent user",
					token:       adminToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "nonexistent@orgusers.example",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "disable already disabled user",
					token:       adminToken,
					request: vetchi.DisableOrgUserRequest{
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
				request       vetchi.FilterOrgUsersRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []filterOrgUsersTestCase{
				{
					description: "with Admin token - no filters",
					token:       adminToken,
					request:     vetchi.FilterOrgUsersRequest{},
					wantStatus:  http.StatusOK,
				},
				{
					description: "with CRUD token - prefix filter",
					token:       crudToken,
					request: vetchi.FilterOrgUsersRequest{
						Prefix: "admin",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - state filter",
					token:       viewerToken,
					request: vetchi.FilterOrgUsersRequest{
						State: []vetchi.OrgUserState{"DISABLED_ORG_USER"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-orgusers token",
					token:       nonOrgUsersToken,
					request:     vetchi.FilterOrgUsersRequest{},
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with invalid limit",
					token:       adminToken,
					request: vetchi.FilterOrgUsersRequest{
						Limit: 41,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"limit"},
				},
				{
					description: "with invalid state",
					token:       adminToken,
					request: vetchi.FilterOrgUsersRequest{
						State: []vetchi.OrgUserState{"INVALID_STATE"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"state"},
				},
				{
					description: "with pagination",
					token:       adminToken,
					request: vetchi.FilterOrgUsersRequest{
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
			prefixUsers := []vetchi.AddOrgUserRequest{
				{
					Email: "alpha1@orgusers.example",
					Name:  "Beta User One",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "beta1@orgusers.example",
					Name:  "Alpha User One",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "gamma1@orgusers.example",
					Name:  "Alpha User Two",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
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
					vetchi.FilterOrgUsersRequest{Prefix: tc.prefix},
					"/employer/filter-org-users",
					http.StatusOK,
				).([]byte)

				var users []vetchi.OrgUser
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
					wantStatus:  http.StatusForbidden,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					vetchi.FilterOrgUsersRequest{Prefix: "alpha"},
					"/employer/filter-org-users",
					tc.wantStatus,
				)
			}
		})

		It("Enable OrgUser", func() {
			// First create and disable some test users that we can enable
			testUsers := []vetchi.AddOrgUserRequest{
				{
					Email: "to-enable1@orgusers.example",
					Name:  "To Enable User 1",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_CRUD"},
				},
				{
					Email: "to-enable2@orgusers.example",
					Name:  "To Enable User 2",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
				},
				{
					Email: "active-user@orgusers.example",
					Name:  "Active User",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
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

			// Disable users that we'll enable in our tests
			disableUsers := []string{
				"to-enable1@orgusers.example",
				"to-enable2@orgusers.example",
			}

			for _, email := range disableUsers {
				testPOST(
					adminToken,
					vetchi.DisableOrgUserRequest{Email: email},
					"/employer/disable-org-user",
					http.StatusOK,
				)
			}

			type enableOrgUserTestCase struct {
				description string
				token       string
				request     vetchi.EnableOrgUserRequest
				wantStatus  int
			}

			testCases := []enableOrgUserTestCase{
				{
					description: "with Admin token - enable disabled user",
					token:       adminToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "to-enable1@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token - enable disabled user",
					token:       crudToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - should fail",
					token:       viewerToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with non-orgusers token - should fail",
					token:       nonOrgUsersToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "to-enable2@orgusers.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "enable non-existent user",
					token:       adminToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "nonexistent@orgusers.example",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "enable already active user",
					token:       adminToken,
					request: vetchi.EnableOrgUserRequest{
						Email: "active-user@orgusers.example",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "enable already enabled user",
					token:       adminToken,
					request: vetchi.EnableOrgUserRequest{
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
				vetchi.FilterOrgUsersRequest{
					State: []vetchi.OrgUserState{vetchi.AddedOrgUserState},
				},
				"/employer/filter-org-users",
				http.StatusOK,
			).([]byte)

			var users []vetchi.OrgUser
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
						).Should(Equal(vetchi.AddedOrgUserState))
						break
					}
				}
				Expect(found).Should(BeTrue(), "Enabled but notfound %q", email)
			}

			// TODO Check that the invite email was sent to the enabled users by querying mailpit
		})

		FIt("Test SignUp of OrgUsers", func() {
			type signupOrgUserTestCase struct {
				description   string
				request       vetchi.SignupOrgUserRequest
				wantStatus    int
				wantErrFields []string
			}

			// First create an org user that we can use to get an invite token
			testPOST(
				adminToken,
				vetchi.AddOrgUserRequest{
					Email: "to-signup@orgusers.example",
					Name:  "To Signup User",
					Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
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
				"to:to-signup@orgusers.example subject:Vetchi Employer Invitation",
			)
			baseURL.RawQuery = query.Encode()

			finalURL := baseURL.String()
			fmt.Fprintf(GinkgoWriter, "finalURL: %s\n", finalURL)

			var messageID string
			// Retry a few times as email delivery might take time
			for i := 0; i < 3; i++ {
				<-time.After(10 * time.Second)
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
			// The link format is: vetchi.example/signup-orguser/<token>
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
					request: vetchi.SignupOrgUserRequest{
						InviteToken: inviteToken,
						Name:        "Signed Up User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-existent token",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: "non-existent-token",
						Name:        "Invalid Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with already used token",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: inviteToken, // Using the same token again
						Name:        "Duplicate Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with empty token",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: "",
						Name:        "Empty Token User",
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"invite_token"},
				},
				{
					description: "with short name",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        "ab", // Too short
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with long name",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        strings.Repeat("a", 257), // Too long
						Password:    "ValidPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"name"},
				},
				{
					description: "with empty password",
					request: vetchi.SignupOrgUserRequest{
						InviteToken: "some-token",
						Name:        "Valid Name",
						Password:    "",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
				{
					description: "with short password",
					request: vetchi.SignupOrgUserRequest{
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
						fmt.Sprintf("%s/employer/signup-org-user", serverURL),
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
	request vetchi.AddOrgUserRequest,
	wantStatus int,
) vetchi.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/add-org-user", wantStatus).([]byte)
	var validationErrors vetchi.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func testUpdateOrgUserGetResp(
	token string,
	request vetchi.UpdateOrgUserRequest,
	wantStatus int,
) vetchi.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/update-org-user", wantStatus).([]byte)
	var validationErrors vetchi.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func testFilterOrgUsersGetResp(
	token string,
	request vetchi.FilterOrgUsersRequest,
	wantStatus int,
) vetchi.ValidationErrors {
	resp := testPOSTGetResp(token, request, "/employer/filter-org-users", wantStatus).([]byte)
	var validationErrors vetchi.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func bulkAddFilterOrgUsers(token string, runID string, count int, limit int) {
	wantUsers := []vetchi.OrgUser{}

	for i := 0; i < count; i++ {
		email := fmt.Sprintf("bulk-%s-%d@orgusers.example", runID, i)
		name := fmt.Sprintf("Bulk User %s-%d", runID, i)

		request := vetchi.AddOrgUserRequest{
			Email: email,
			Name:  name,
			Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
		}

		fmt.Fprintf(GinkgoWriter, "Adding user %q %q\n", email, name)
		testPOST(token, request, "/employer/add-org-user", http.StatusOK)

		wantUsers = append(wantUsers, vetchi.OrgUser{
			Email: email,
			Name:  name,
			Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
			State: vetchi.AddedOrgUserState,
		})
	}

	paginationKey := ""
	gotUsers := []vetchi.OrgUser{}

	for {
		request := vetchi.FilterOrgUsersRequest{
			PaginationKey: paginationKey,
			Limit:         limit,
		}

		resp := testPOSTGetResp(token, request, "/employer/filter-org-users", http.StatusOK).([]byte)

		var users []vetchi.OrgUser
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
	request vetchi.SignupOrgUserRequest,
	wantStatus int,
) vetchi.ValidationErrors {
	jsonBytes, err := json.Marshal(request)
	Expect(err).ShouldNot(HaveOccurred())

	resp, err := http.Post(
		fmt.Sprintf("%s/employer/signup-org-user", serverURL),
		"application/json",
		bytes.NewBuffer(jsonBytes),
	)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(wantStatus))

	body, err := io.ReadAll(resp.Body)
	Expect(err).ShouldNot(HaveOccurred())
	defer resp.Body.Close()

	var validationErrors vetchi.ValidationErrors
	err = json.Unmarshal(body, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}
