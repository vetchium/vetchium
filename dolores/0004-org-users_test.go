package dolores

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Org Users", Ordered, func() {
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
				description string
				token       string
				request     vetchi.AddOrgUserRequest
				wantStatus  int
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
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with short name",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "short@orgusers.example",
						Name:  "ab",
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long name",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "long@orgusers.example",
						Name:  strings.Repeat("a", 257),
						Roles: []vetchi.OrgUserRole{"ORG_USERS_VIEWER"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with empty roles",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "noroles@orgusers.example",
						Name:  "No Roles User",
						Roles: []vetchi.OrgUserRole{},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid role",
					token:       adminToken,
					request: vetchi.AddOrgUserRequest{
						Email: "invalid@orgusers.example",
						Name:  "Invalid Role User",
						Roles: []vetchi.OrgUserRole{"INVALID_ROLE"},
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/add-org-user",
					tc.wantStatus,
				)
			}
		})

		It("Update OrgUser", func() {
			type updateOrgUserTestCase struct {
				description string
				token       string
				request     vetchi.UpdateOrgUserRequest
				wantStatus  int
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
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/update-org-user",
					tc.wantStatus,
				)
			}
		})

		It("Disable OrgUser", func() {
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
						Email: "crud@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "viewer@orgusers.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: vetchi.DisableOrgUserRequest{
						Email: "non-orgusers@orgusers.example",
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
						Email: "disabled@orgusers.example",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/disable-org-user",
					tc.wantStatus,
				)
			}
		})

		It("Filter OrgUsers", func() {
			type filterOrgUsersTestCase struct {
				description string
				token       string
				request     vetchi.FilterOrgUsersRequest
				wantStatus  int
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
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid state",
					token:       adminToken,
					request: vetchi.FilterOrgUsersRequest{
						State: []vetchi.OrgUserState{"INVALID_STATE"},
					},
					wantStatus: http.StatusBadRequest,
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
				testPOST(
					tc.token,
					tc.request,
					"/employer/filter-org-users",
					tc.wantStatus,
				)
			}
		})
	})
})
