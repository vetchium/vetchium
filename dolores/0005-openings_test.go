package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var bachelorEducation *vetchi.EducationLevel

var _ = XDescribe("Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOpeningsToken string
	var recruiterToken, hiringManagerToken string

	bachelor := vetchi.BachelorEducation
	bachelorEducation = &bachelor

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0005-openings-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@openings.example":          &adminToken,
			"crud@openings.example":           &crudToken,
			"viewer@openings.example":         &viewerToken,
			"non-openings@openings.example":   &nonOpeningsToken,
			"recruiter@openings.example":      &recruiterToken,
			"hiring-manager@openings.example": &hiringManagerToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"openings.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0005-openings-down.pgsql")
		db.Close()
	})

	Describe("Openings Tests", func() {
		It("Create Opening", func() {
			validOpening := vetchi.CreateOpeningRequest{
				Title:         "Software Engineer",
				Positions:     2,
				JD:            "Looking for talented software engineers",
				Recruiter:     "recruiter@openings.example",
				HiringManager: "hiring-manager@openings.example",
				HiringTeam: []vetchi.EmailAddress{
					"crud@openings.example",
					"viewer@openings.example",
				},
				CostCenterName: "Engineering",
				LocationTitles: []string{
					"Bangalore Office",
					"Chennai Office",
				},
				RemoteCountryCodes: []vetchi.CountryCode{
					"IND",
					"USA",
				},
				RemoteTimezones: []vetchi.TimeZone{
					"IST Indian Standard Time GMT+0530",
				},
				OpeningType:       vetchi.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: bachelorEducation,
				Salary: &vetchi.Salary{
					MinAmount: 50000,
					MaxAmount: 100000,
					Currency:  "USD",
				},
			}

			type createOpeningTestCase struct {
				description   string
				token         string
				request       vetchi.CreateOpeningRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []createOpeningTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request:     validOpening,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request:     validOpening,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request:     validOpening,
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request:     validOpening,
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with missing title",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.Title = ""
						return r
					}(),
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with invalid positions",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.Positions = 0
						return r
					}(),
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"positions"},
				},
				{
					description: "with invalid JD",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.JD = "short"
						return r
					}(),
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"jd"},
				},
				{
					description: "with non-existent recruiter",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.Recruiter = "nonexistent@openings.example"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent hiring manager",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.HiringManager = "nonexistent@openings.example"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent cost center",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.CostCenterName = "NonExistent"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent location",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.LocationTitles = []string{"NonExistent"}
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with invalid YOE range",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.YoeMin = 6
						r.YoeMax = 5
						return r
					}(),
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"yoe_min", "yoe_max"},
				},
				{
					description: "with invalid salary range",
					token:       adminToken,
					request: func() vetchi.CreateOpeningRequest {
						r := validOpening
						r.Salary.MinAmount = 200000
						r.Salary.MaxAmount = 100000
						return r
					}(),
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"salary"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				if len(tc.wantErrFields) > 0 {
					validationErrors := testCreateOpeningGetResp(
						tc.token,
						tc.request,
						tc.wantStatus,
					)
					Expect(validationErrors.Errors).Should(
						ContainElements(tc.wantErrFields),
					)
				} else {
					testPOST(
						tc.token,
						tc.request,
						"/employer/create-opening",
						tc.wantStatus,
					)
				}
			}
		})

		It("Filter Openings", func() {
			// First create some test openings
			createTestOpenings(adminToken)

			type filterOpeningsTestCase struct {
				description string
				token       string
				request     vetchi.FilterOpeningsRequest
				wantStatus  int
			}

			testCases := []filterOpeningsTestCase{
				{
					description: "with Admin token - no filters",
					token:       adminToken,
					request:     vetchi.FilterOpeningsRequest{},
					wantStatus:  http.StatusOK,
				},
				{
					description: "with CRUD token - state filter",
					token:       crudToken,
					request: vetchi.FilterOpeningsRequest{
						State: []vetchi.OpeningState{vetchi.DraftOpening},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - date filter",
					token:       viewerToken,
					request: vetchi.FilterOpeningsRequest{
						FromDate: &time.Time{},
						ToDate:   &time.Time{},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request:     vetchi.FilterOpeningsRequest{},
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with invalid date range",
					token:       adminToken,
					request: vetchi.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t := time.Now().AddDate(0, 0, 1)
							return &t
						}(),
						ToDate: func() *time.Time {
							t := time.Now()
							return &t
						}(),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid limit",
					token:       adminToken,
					request: vetchi.FilterOpeningsRequest{
						Limit: 41,
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/filter-openings",
					tc.wantStatus,
				)
			}
		})

		It("Get Opening", func() {
			// First create a test opening
			openingID := createTestOpening(adminToken)

			type getOpeningTestCase struct {
				description string
				token       string
				request     vetchi.GetOpeningRequest
				wantStatus  int
			}

			testCases := []getOpeningTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: vetchi.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: vetchi.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: vetchi.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request: vetchi.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with non-existent ID",
					token:       adminToken,
					request: vetchi.GetOpeningRequest{
						ID: "non-existent-id",
					},
					wantStatus: http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/get-opening",
					tc.wantStatus,
				)
			}
		})

		It("Test Opening Pagination", func() {
			// Create bulk openings and test pagination
			bulkCreateOpenings(
				adminToken,
				"run-1",
				30,
				4,
			) // count not divisible by limit
			bulkCreateOpenings(
				adminToken,
				"run-2",
				32,
				4,
			) // count divisible by limit
			bulkCreateOpenings(
				adminToken,
				"run-3",
				2,
				4,
			) // count less than limit
		})
	})
})

// Helper functions

func testCreateOpeningGetResp(
	token string,
	request vetchi.CreateOpeningRequest,
	wantStatus int,
) vetchi.ValidationErrors {
	resp := testPOSTGetResp(
		token,
		request,
		"/employer/create-opening",
		wantStatus,
	).([]byte)
	var validationErrors vetchi.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func createTestOpening(token string) string {
	request := vetchi.CreateOpeningRequest{
		Title:             "Test Opening",
		Positions:         1,
		JD:                "Test Job Description",
		Recruiter:         "recruiter@openings.example",
		HiringManager:     "hiring-manager@openings.example",
		CostCenterName:    "Engineering",
		OpeningType:       vetchi.FullTimeOpening,
		YoeMin:            0,
		YoeMax:            5,
		MinEducationLevel: bachelorEducation,
		Salary: &vetchi.Salary{
			MinAmount: 50000,
			MaxAmount: 100000,
			Currency:  "USD",
		},
	}

	resp := testPOSTGetResp(
		token,
		request,
		"/employer/create-opening",
		http.StatusOK,
	).([]byte)

	var opening vetchi.Opening
	err := json.Unmarshal(resp, &opening)
	Expect(err).ShouldNot(HaveOccurred())

	return opening.ID
}

func createTestOpenings(token string) {

	for i := 0; i < 3; i++ {
		request := vetchi.CreateOpeningRequest{
			Title:             fmt.Sprintf("Test Opening %d", i),
			Positions:         1,
			JD:                fmt.Sprintf("Test Job Description %d", i),
			Recruiter:         "recruiter@openings.example",
			HiringManager:     "hiring-manager@openings.example",
			CostCenterName:    "Engineering",
			OpeningType:       vetchi.FullTimeOpening,
			YoeMin:            0,
			YoeMax:            5,
			MinEducationLevel: bachelorEducation,
			Salary: &vetchi.Salary{
				MinAmount: 50000,
				MaxAmount: 100000,
				Currency:  "USD",
			},
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/create-opening",
			http.StatusOK,
		).([]byte)

		var opening vetchi.Opening
		err := json.Unmarshal(resp, &opening)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

func bulkCreateOpenings(token string, runID string, count int, limit int) {
	wantOpenings := []string{}

	for i := 0; i < count; i++ {
		request := vetchi.CreateOpeningRequest{
			Title:     fmt.Sprintf("Bulk Opening %s-%d", runID, i),
			Positions: 1,
			JD: fmt.Sprintf(
				"Bulk Job Description %s-%d",
				runID,
				i,
			),
			Recruiter:         "recruiter@openings.example",
			HiringManager:     "hiring-manager@openings.example",
			CostCenterName:    "Engineering",
			OpeningType:       vetchi.FullTimeOpening,
			YoeMin:            0,
			YoeMax:            5,
			MinEducationLevel: bachelorEducation,
			Salary: &vetchi.Salary{
				MinAmount: 50000,
				MaxAmount: 100000,
				Currency:  "USD",
			},
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/create-opening",
			http.StatusOK,
		).([]byte)

		var opening vetchi.Opening
		err := json.Unmarshal(resp, &opening)
		Expect(err).ShouldNot(HaveOccurred())
		wantOpenings = append(wantOpenings, opening.ID)
	}

	paginationKey := ""
	gotOpenings := []string{}

	for {
		request := vetchi.FilterOpeningsRequest{
			PaginationKey: paginationKey,
			Limit:         limit,
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/filter-openings",
			http.StatusOK,
		).([]byte)

		var openings []vetchi.OpeningInfo
		err := json.Unmarshal(resp, &openings)
		Expect(err).ShouldNot(HaveOccurred())

		if len(openings) == 0 {
			break
		}

		for _, opening := range openings {
			gotOpenings = append(gotOpenings, opening.ID)
		}

		if len(openings) < limit {
			break
		}

		paginationKey = openings[len(openings)-1].ID
	}

	Expect(gotOpenings).Should(ContainElements(wantOpenings))
}
