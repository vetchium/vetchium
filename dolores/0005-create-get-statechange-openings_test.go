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
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

var _ = Describe("Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOpeningsToken string
	var recruiterToken, hiringManagerToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0005-create-get-statechange-openings-up.pgsql")

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
		seedDatabase(db, "0005-create-get-statechange-openings-down.pgsql")
		db.Close()
	})

	Describe("Openings Tests", func() {
		It("Create Opening", func() {
			validOpening := employer.CreateOpeningRequest{
				Title:         "Software Engineer",
				Positions:     2,
				JD:            "Looking for talented software engineers",
				Recruiter:     "recruiter@openings.example",
				HiringManager: "hiring-manager@openings.example",
				HiringTeam: []common.EmailAddress{
					"crud@openings.example",
					"viewer@openings.example",
				},
				CostCenterName: "Engineering",
				LocationTitles: []string{
					"Bangalore Office",
					"Chennai Office",
				},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				RemoteTimezones: []common.TimeZone{
					"IST Indian Standard Time GMT+0530",
				},
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.BachelorEducation,
				Salary: &common.Salary{
					MinAmount: 50000,
					MaxAmount: 100000,
					Currency:  "USD",
				},
				NewTags: []string{"DevOps"},
			}

			type createOpeningTestCase struct {
				description   string
				token         string
				request       employer.CreateOpeningRequest
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
					request: func() employer.CreateOpeningRequest {
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
					request: func() employer.CreateOpeningRequest {
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
					request: func() employer.CreateOpeningRequest {
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
					request: func() employer.CreateOpeningRequest {
						r := validOpening
						r.Recruiter = "nonexistent@openings.example"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent hiring manager",
					token:       adminToken,
					request: func() employer.CreateOpeningRequest {
						r := validOpening
						r.HiringManager = "nonexistent@openings.example"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent cost center",
					token:       adminToken,
					request: func() employer.CreateOpeningRequest {
						r := validOpening
						r.CostCenterName = "NonExistent"
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-existent location",
					token:       adminToken,
					request: func() employer.CreateOpeningRequest {
						r := validOpening
						r.LocationTitles = []string{"NonExistent"}
						return r
					}(),
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with invalid YOE range",
					token:       adminToken,
					request: func() employer.CreateOpeningRequest {
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
					request: func() employer.CreateOpeningRequest {

						return employer.CreateOpeningRequest{
							Title:         "Software Engineer",
							Positions:     2,
							JD:            "Looking for talented software engineers",
							Recruiter:     "recruiter@openings.example",
							HiringManager: "hiring-manager@openings.example",
							HiringTeam: []common.EmailAddress{
								"crud@openings.example",
								"viewer@openings.example",
							},
							CostCenterName: "Engineering",
							LocationTitles: []string{
								"Bangalore Office",
								"Chennai Office",
							},
							RemoteCountryCodes: []common.CountryCode{
								"IND",
								"USA",
							},
							RemoteTimezones: []common.TimeZone{
								"IST Indian Standard Time GMT+0530",
							},
							OpeningType:       common.FullTimeOpening,
							YoeMin:            2,
							YoeMax:            5,
							MinEducationLevel: common.BachelorEducation,
							Salary: &common.Salary{
								MinAmount: 200000,
								MaxAmount: 100000,
								Currency:  "USD",
							},
							NewTags: []string{"DevOps"},
						}
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

		It("Filter Openings on openings state", func() {
			fmt.Fprintf(GinkgoWriter, "#### Creating test openings\n")
			createTestOpenings(adminToken)
			fmt.Fprintf(GinkgoWriter, "#### Test openings created\n")

			type filterOpeningsTestCase struct {
				description string
				token       string
				request     employer.FilterOpeningsRequest
				wantStatus  int
			}

			testCases := []filterOpeningsTestCase{
				{
					description: "with Admin token - no filters",
					token:       adminToken,
					request:     employer.FilterOpeningsRequest{},
					wantStatus:  http.StatusOK,
				},
				{
					description: "with CRUD token - state filter",
					token:       crudToken,
					request: employer.FilterOpeningsRequest{
						State: []common.OpeningState{common.DraftOpening},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token - date filter",
					token:       viewerToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t := time.Now().AddDate(0, 0, -30)
							return &t
						}(),
						ToDate: func() *time.Time {
							t := time.Now()
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request:     employer.FilterOpeningsRequest{},
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with invalid date range",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
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
					description: "with todate < fromdate",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t := time.Now()
							return &t
						}(),
						ToDate: func() *time.Time {
							t := time.Now().AddDate(0, 0, -1)
							return &t
						}(),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid limit",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
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
				request     employer.GetOpeningRequest
				wantStatus  int
			}

			testCases := []getOpeningTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: employer.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: employer.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					request: employer.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request: employer.GetOpeningRequest{
						ID: openingID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with non-existent ID",
					token:       adminToken,
					request: employer.GetOpeningRequest{
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

		It("Change Opening State", func() {
			openingID := createTestOpening(adminToken)

			type changeOpeningStateTestCase struct {
				description string
				token       string
				request     employer.ChangeOpeningStateRequest
				wantStatus  int
			}

			testCases := []changeOpeningStateTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.DraftOpening,
						ToState:   common.ActiveOpening,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with CRUD token",
					token:       crudToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.ActiveOpening,
						ToState:   common.SuspendedOpening,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings token",
					token:       nonOpeningsToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.SuspendedOpening,
						ToState:   common.ClosedOpening,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with viewer token",
					token:       viewerToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.SuspendedOpening,
						ToState:   common.ClosedOpening,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with non-existent opening",
					token:       adminToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: "non-existent-id",
						FromState: common.DraftOpening,
						ToState:   common.ActiveOpening,
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "with invalid transition",
					token:       adminToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.SuspendedOpening,
						ToState:   common.DraftOpening,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with invalid from state",
					token:       adminToken,
					request: employer.ChangeOpeningStateRequest{
						OpeningID: openingID,
						FromState: common.DraftOpening,
						ToState:   common.ActiveOpening,
					},
					wantStatus: http.StatusConflict,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/change-opening-state",
					tc.wantStatus,
				)
			}
		})
	})
})

// Helper functions

func testCreateOpeningGetResp(
	token string,
	request employer.CreateOpeningRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(
		token,
		request,
		"/employer/create-opening",
		wantStatus,
	).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func createTestOpening(token string) string {
	request := employer.CreateOpeningRequest{
		Title:             "Test Opening",
		Positions:         1,
		JD:                "Test Job Description",
		Recruiter:         "recruiter@openings.example",
		HiringManager:     "hiring-manager@openings.example",
		CostCenterName:    "Engineering",
		OpeningType:       common.FullTimeOpening,
		YoeMin:            0,
		YoeMax:            5,
		MinEducationLevel: common.BachelorEducation,
		Salary: &common.Salary{
			MinAmount: 50000,
			MaxAmount: 100000,
			Currency:  "USD",
		},
		RemoteCountryCodes: []common.CountryCode{
			"IND",
			"USA",
		},
		NewTags: []string{"DevOps"},
	}

	resp := testPOSTGetResp(
		token,
		request,
		"/employer/create-opening",
		http.StatusOK,
	).([]byte)

	var opening employer.CreateOpeningResponse
	err := json.Unmarshal(resp, &opening)
	Expect(err).ShouldNot(HaveOccurred())

	return opening.OpeningID
}

func createTestOpenings(token string) {
	for i := 0; i < 3; i++ {
		request := employer.CreateOpeningRequest{
			Title:             fmt.Sprintf("Test Opening %d", i),
			Positions:         1,
			JD:                fmt.Sprintf("Test Job Description %d", i),
			Recruiter:         "recruiter@openings.example",
			HiringManager:     "hiring-manager@openings.example",
			CostCenterName:    "Engineering",
			OpeningType:       common.FullTimeOpening,
			YoeMin:            0,
			YoeMax:            5,
			MinEducationLevel: common.BachelorEducation,
			Salary: &common.Salary{
				MinAmount: 50000,
				MaxAmount: 100000,
				Currency:  "USD",
			},
			RemoteCountryCodes: []common.CountryCode{
				"IND",
				"USA",
			},
			NewTags: []string{"DevOps"},
		}

		_ = testPOSTGetResp(
			token,
			request,
			"/employer/create-opening",
			http.StatusOK,
		).([]byte)
	}
}

func bulkCreateOpenings(token string, runID string, count int, limit int) {
	wantOpenings := []string{}

	for i := 0; i < count; i++ {
		request := employer.CreateOpeningRequest{
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
			OpeningType:       common.FullTimeOpening,
			YoeMin:            0,
			YoeMax:            5,
			MinEducationLevel: common.BachelorEducation,
			Salary: &common.Salary{
				MinAmount: 50000,
				MaxAmount: 100000,
				Currency:  "USD",
			},
			RemoteCountryCodes: []common.CountryCode{
				"IND",
				"USA",
			},
			NewTags: []string{"DevOps"},
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/create-opening",
			http.StatusOK,
		).([]byte)

		var opening employer.CreateOpeningResponse
		err := json.Unmarshal(resp, &opening)
		Expect(err).ShouldNot(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "Appending opening: %+v\n", opening)
		wantOpenings = append(wantOpenings, opening.OpeningID)
	}

	paginationKey := ""
	gotOpenings := []string{}

	for {
		request := employer.FilterOpeningsRequest{
			PaginationKey: paginationKey,
			Limit:         limit,
		}

		resp := testPOSTGetResp(
			token,
			request,
			"/employer/filter-openings",
			http.StatusOK,
		).([]byte)

		var openings []employer.OpeningInfo
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
