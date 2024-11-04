package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

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
