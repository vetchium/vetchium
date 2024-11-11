package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type costCenterTestCase struct {
	description    string
	token          string
	costCenterName string
	wantStatus     int
}

var _ = Describe("Cost Centers", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, viewerToken string
	var nonCostCenterToken, multipleNonCostCenterRolesToken string
	var crud1Token, crud2Token string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0002-cost-centers-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@cost-center.example":                          &adminToken,
			"viewer@cost-center.example":                         &viewerToken,
			"non-cost-center@cost-center.example":                &nonCostCenterToken,
			"multiple-non-cost-center-roles@cost-center.example": &multipleNonCostCenterRolesToken,
			"crud1@cost-center.example":                          &crud1Token,
			"crud2@cost-center.example":                          &crud2Token,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"cost-center.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}

		// Wait until all the signin operations are complete
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0002-cost-centers-down.pgsql")
		db.Close()
	})

	Describe("Cost Centers related Tests", func() {
		It("Add Cost Center", func() {
			testCases := []costCenterTestCase{
				{
					description:    "without a session token",
					costCenterName: "New Cost Center",
					token:          "",
					wantStatus:     http.StatusUnauthorized,
				},
				{
					description:    "with an invalid session token",
					costCenterName: "New Cost Center",
					token:          "blah blah blah",
					wantStatus:     http.StatusUnauthorized,
				},
				{
					description:    "with an invalid name",
					costCenterName: "",
					token:          adminToken,
					wantStatus:     http.StatusBadRequest,
				},
				{
					description:    "with viewer role",
					costCenterName: "New Cost Center",
					token:          viewerToken,
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with non-cost-center role",
					costCenterName: "New Cost Center",
					token:          nonCostCenterToken,
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with multiple non-cost-center roles",
					costCenterName: "New Cost Center",
					token:          multipleNonCostCenterRolesToken,
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with an admin session token",
					costCenterName: "CC1-Admin",
					token:          adminToken,
					wantStatus:     http.StatusOK,
				},
				{
					description:    "with a duplicate name as Admin",
					costCenterName: "CC1-Admin",
					token:          adminToken,
					wantStatus:     http.StatusConflict,
				},
				{
					description:    "with cost-centers-crud session token",
					costCenterName: "CC2-Crud",
					token:          crud1Token,
					wantStatus:     http.StatusOK,
				},
				{
					description:    "with a duplicate name as Crud",
					costCenterName: "CC2-Crud",
					token:          crud2Token,
					wantStatus:     http.StatusConflict,
				},
			}

			for _, i := range testCases {
				testAddCostCenter(i.token, i.costCenterName, i.wantStatus)
				fmt.Fprintf(GinkgoWriter, "%s\n", i.description)
			}

		})

		It("should return the cost centers for the employer", func() {
			testGetCostCenters(
				adminToken,
				2,
				[]string{"CC1-Admin", "CC2-Crud"},
			)
		})

		It("Defunct Cost Center", func() {
			testCases := []costCenterTestCase{
				{
					description:    "with no session token",
					token:          "",
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusUnauthorized,
				},
				{
					description:    "with an invalid session token",
					token:          "blah blah blah",
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusUnauthorized,
				},
				{
					description:    "with a viewer session token",
					token:          viewerToken,
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with a non-cost-center session",
					token:          nonCostCenterToken,
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with a multiple-non-cost-center session",
					token:          multipleNonCostCenterRolesToken,
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusForbidden,
				},
				{
					description:    "with an admin session token",
					token:          adminToken,
					costCenterName: "CC1-Admin",
					wantStatus:     http.StatusOK,
				},
				{
					description:    "with cost-centers-crud session token",
					token:          crud1Token,
					costCenterName: "CC2-Crud",
					wantStatus:     http.StatusOK,
				},
				{
					description:    "with an invalid name",
					token:          adminToken,
					costCenterName: "",
					wantStatus:     http.StatusBadRequest,
				},
				{
					description:    "with a name that doesn't exist",
					token:          adminToken,
					costCenterName: "Non-existent Cost Center",
					wantStatus:     http.StatusNotFound,
				},
				{
					description:    "with a name that is too long",
					token:          adminToken,
					costCenterName: strings.Repeat("A", 65),
					wantStatus:     http.StatusBadRequest,
				},
				{
					description:    "with a name that is too short",
					token:          adminToken,
					costCenterName: "A",
					wantStatus:     http.StatusBadRequest,
				},
			}

			for _, testCase := range testCases {
				testDefunctCostCenter(
					testCase.token,
					testCase.costCenterName,
					testCase.wantStatus,
				)
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
			}
		})

		It("get the list of defunct cost centers", func() {
			getDefunctCostCentersReqBody, err := json.Marshal(
				vetchi.GetCostCentersRequest{
					States: []vetchi.CostCenterState{
						vetchi.DefunctCC,
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			getDefunctCostCentersReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getDefunctCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getDefunctCostCentersReq.Header.Set("Authorization", adminToken)

			resp, err := http.DefaultClient.Do(getDefunctCostCentersReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(2))
			Expect(
				[]string{
					string(costCenters[0].Name),
					string(costCenters[1].Name),
				},
			).Should(ConsistOf("CC1-Admin", "CC2-Crud"))
		})

		It(
			"create/defunct cost centers in bulk; verify pagination",
			func() {
				// runID parameter is needed because even if a CC is defunct,
				// its name cannot be used by another CC

				fmt.Fprintf(
					GinkgoWriter,
					"count is not divisible by limit\n",
				)
				bulkAddDefunctCC(adminToken, serverURL, "run-1", 30, 4)

				fmt.Fprintf(GinkgoWriter, "count is divisible by limit\n")
				bulkAddDefunctCC(adminToken, serverURL, "run-2", 32, 4)

				fmt.Fprintf(GinkgoWriter, "count is less than limit\n")
				bulkAddDefunctCC(adminToken, serverURL, "run-3", 2, 4)
			},
		)
	})

	It("update a cost center", func() {
		statusCode, err := addCostCenter(
			adminToken,
			vetchi.AddCostCenterRequest{
				Name:  "CC-update-test-1",
				Notes: "This is a test cost center",
			})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusOK))

		ccb4Update, _, err := getCostCenter(
			adminToken,
			vetchi.GetCostCenterRequest{
				Name: "CC-update-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ccb4Update.Name).Should(Equal("CC-update-test-1"))
		Expect(ccb4Update.Notes).Should(Equal("This is a test cost center"))

		fmt.Fprintf(GinkgoWriter, "Update Cost Center without Notes\n")
		statusCode, err = updateCostCenter(
			adminToken,
			vetchi.UpdateCostCenterRequest{
				Name: "CC-update-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusBadRequest))

		fmt.Fprintf(GinkgoWriter, "Update Cost Center with invalid notes\n")
		statusCode, err = updateCostCenter(
			adminToken,
			vetchi.UpdateCostCenterRequest{
				Name:  "CC-update-test-1",
				Notes: strings.Repeat("A", 1025),
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusBadRequest))

		fmt.Fprintf(GinkgoWriter, "Updating cost center with new notes\n")
		statusCode, err = updateCostCenter(
			adminToken,
			vetchi.UpdateCostCenterRequest{
				Name:  "CC-update-test-1",
				Notes: "This is an updated test cost center",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusOK))

		ccAfterUpdate, _, err := getCostCenter(
			adminToken,
			vetchi.GetCostCenterRequest{
				Name: "CC-update-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ccAfterUpdate.Notes).Should(Equal(
			"This is an updated test cost center",
		))
	})

	It("update a cost center with a viewer session token", func() {
		statusCode, err := updateCostCenter(
			viewerToken,
			vetchi.UpdateCostCenterRequest{
				Name:  "CC-update-test-1",
				Notes: "This is an updated test cost center",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusForbidden))
	})

	It("update a cost center with a non-cost-center session token", func() {
		statusCode, err := updateCostCenter(
			nonCostCenterToken,
			vetchi.UpdateCostCenterRequest{
				Name:  "CC-update-test-1",
				Notes: "This is an updated test cost center",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusForbidden))
	})

	It(
		"update a cost center with a multiple-non-cost-center session token",
		func() {
			statusCode, err := updateCostCenter(
				multipleNonCostCenterRolesToken,
				vetchi.UpdateCostCenterRequest{
					Name:  "CC-update-test-1",
					Notes: "This is an updated test cost center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(statusCode).Should(Equal(http.StatusForbidden))
		},
	)

	It("rename a cost center", func() {
		statusCode, err := addCostCenter(
			adminToken,
			vetchi.AddCostCenterRequest{
				Name: "CC-rename-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusOK))

		renameCostCenterReq := vetchi.RenameCostCenterRequest{
			OldName: "CC-rename-test-1",
			NewName: "CC-rename-test-2",
		}

		renameCostCenterReqBody, err := json.Marshal(renameCostCenterReq)
		Expect(err).ShouldNot(HaveOccurred())

		req, err := http.NewRequest(
			http.MethodPost,
			serverURL+"/employer/rename-cost-center",
			bytes.NewBuffer(renameCostCenterReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())

		req.Header.Set("Authorization", adminToken)

		resp, err := http.DefaultClient.Do(req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		ccAfterRename, _, err := getCostCenter(
			adminToken,
			vetchi.GetCostCenterRequest{
				Name: "CC-rename-test-2",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ccAfterRename.Name).Should(Equal("CC-rename-test-2"))

		_, statusCode, err = getCostCenter(
			adminToken,
			vetchi.GetCostCenterRequest{
				Name: "CC-rename-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusNotFound))
	})

})

func testAddCostCenter(token, name string, expectedStatus int) {
	fmt.Fprintf(
		GinkgoWriter,
		"testAddCostCenter: token=%s, name=%s, expectedStatus=%d\n",
		token,
		name,
		expectedStatus,
	)
	reqBody := vetchi.AddCostCenterRequest{Name: vetchi.CostCenterName(name)}
	testPOST(token, reqBody, "/employer/add-cost-center", expectedStatus)
}

func testDefunctCostCenter(token, name string, expectedStatus int) {
	fmt.Fprintf(
		GinkgoWriter,
		"testDefunctCostCenter: token=%s, name=%s, expectedStatus=%d\n",
		token,
		name,
		expectedStatus,
	)
	reqBody := vetchi.DefunctCostCenterRequest{
		Name: vetchi.CostCenterName(name),
	}
	testPOST(
		token,
		reqBody,
		"/employer/defunct-cost-center",
		expectedStatus,
	)
}

func testGetCostCenters(token string, expectedLen int, expectedNames []string) {
	getCostCentersReqBody, err := json.Marshal(vetchi.GetCostCentersRequest{})
	Expect(err).ShouldNot(HaveOccurred())

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/get-cost-centers",
		bytes.NewBuffer(getCostCentersReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(http.StatusOK))

	var costCenters []vetchi.CostCenter
	err = json.NewDecoder(resp.Body).Decode(&costCenters)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(costCenters).Should(HaveLen(expectedLen))
	Expect(
		[]string{
			string(costCenters[0].Name),
			string(costCenters[1].Name),
		},
	).Should(ConsistOf(expectedNames))
}

// addCostCenter adds a cost center and returns the http status code. The status
// code should be used only if err is nil.
func addCostCenter(
	token string,
	addCostCenterReq vetchi.AddCostCenterRequest,
) (int, error) {
	addCostCenterReqBody, err := json.Marshal(addCostCenterReq)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/add-cost-center",
		bytes.NewBuffer(addCostCenterReqBody),
	)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}

	return resp.StatusCode, nil
}

func updateCostCenter(
	token string,
	updateCostCenterReq vetchi.UpdateCostCenterRequest,
) (int, error) {
	updateCostCenterReqBody, err := json.Marshal(updateCostCenterReq)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/update-cost-center",
		bytes.NewBuffer(updateCostCenterReqBody),
	)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}

	return resp.StatusCode, nil
}

func getCostCenter(
	token string,
	getCostCenterReq vetchi.GetCostCenterRequest,
) (vetchi.CostCenter, int, error) {
	getCostCenterReqBody, err := json.Marshal(getCostCenterReq)
	Expect(err).ShouldNot(HaveOccurred())

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/get-cost-center",
		bytes.NewBuffer(getCostCenterReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())
	if resp.StatusCode != http.StatusOK {
		return vetchi.CostCenter{}, resp.StatusCode, nil
	}

	Expect(resp.StatusCode).Should(Equal(http.StatusOK))

	var costCenter vetchi.CostCenter
	err = json.NewDecoder(resp.Body).Decode(&costCenter)
	Expect(err).ShouldNot(HaveOccurred())

	return costCenter, resp.StatusCode, nil
}

// bulkAddDefunctCC adds a number of defunct cost centers and verifies that
// they are paginated correctly. Not intended to be used in other test cases
// because it defuncts the CCs it creates.
func bulkAddDefunctCC(
	adminToken string,
	serverURL string,
	runID string,
	count, limit int,
) {
	wantCC := []string{}

	for i := 0; i < count; i++ {
		ccName := fmt.Sprintf("CC-%s-%d", runID, i)
		wantCC = append(wantCC, ccName)

		statusCode, err := addCostCenter(
			adminToken,
			vetchi.AddCostCenterRequest{
				Name: vetchi.CostCenterName(ccName),
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusOK))
	}

	paginationKey := ""
	gotCC := []string{}

	for {
		getCostCentersReqBody, err := json.Marshal(
			vetchi.GetCostCentersRequest{
				PaginationKey: vetchi.CostCenterName(paginationKey),
				Limit:         limit,
			},
		)
		Expect(err).ShouldNot(HaveOccurred())

		getCostCentersReq, err := http.NewRequest(
			http.MethodPost,
			serverURL+"/employer/get-cost-centers",
			bytes.NewBuffer(getCostCentersReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())

		getCostCentersReq.Header.Set("Authorization", adminToken)

		resp, err := http.DefaultClient.Do(getCostCentersReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		var costCenters []vetchi.CostCenter
		err = json.NewDecoder(resp.Body).Decode(&costCenters)
		Expect(err).ShouldNot(HaveOccurred())

		if len(costCenters) == 0 {
			break
		}

		for _, costCenter := range costCenters {
			gotCC = append(gotCC, string(costCenter.Name))
		}
		paginationKey = string(costCenters[len(costCenters)-1].Name)

		if len(costCenters) < limit {
			break
		}
	}

	Expect(gotCC).Should(HaveLen(count))
	Expect(gotCC).Should(ContainElements(wantCC))

	for i := 0; i < count; i++ {
		ccName := fmt.Sprintf("CC-%s-%d", runID, i)

		defunctCostCenterReqBody, err := json.Marshal(
			vetchi.DefunctCostCenterRequest{
				Name: vetchi.CostCenterName(ccName),
			},
		)
		Expect(err).ShouldNot(HaveOccurred())

		defunctCostCenterReq, err := http.NewRequest(
			http.MethodPost,
			serverURL+"/employer/defunct-cost-center",
			bytes.NewBuffer(defunctCostCenterReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())

		defunctCostCenterReq.Header.Set("Authorization", adminToken)

		resp, err := http.DefaultClient.Do(defunctCostCenterReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))
	}
}
