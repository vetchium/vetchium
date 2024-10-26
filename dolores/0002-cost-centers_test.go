package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = Describe("Cost Centers", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, viewerToken string
	var nonCostCenterToken, multipleNonCostCenterRolesToken string
	var crud1Token, crud2Token string

	BeforeAll(func() {
		db = setupTestDB()

		seed, err := os.ReadFile("0002-cost-centers-up.pgsql")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(context.Background(), string(seed))
		Expect(err).ShouldNot(HaveOccurred())

		var wg sync.WaitGroup

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"admin@cost-center.example",
			"NewPassword123$",
			&adminToken,
			&wg,
		)

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"viewer@cost-center.example",
			"NewPassword123$",
			&viewerToken,
			&wg,
		)

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"non-cost-center@cost-center.example",
			"NewPassword123$",
			&nonCostCenterToken,
			&wg,
		)

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"multiple-non-cost-center-roles@cost-center.example",
			"NewPassword123$",
			&multipleNonCostCenterRolesToken,
			&wg,
		)

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"crud1@cost-center.example",
			"NewPassword123$",
			&crud1Token,
			&wg,
		)

		wg.Add(1)
		employerSigninAsync(
			"cost-center.example",
			"crud2@cost-center.example",
			"NewPassword123$",
			&crud2Token,
			&wg,
		)

		wg.Wait()
	})

	AfterAll(func() {
		seed, err := os.ReadFile("0002-cost-centers-down.pgsql")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(context.Background(), string(seed))
		Expect(err).ShouldNot(HaveOccurred())

		db.Close()
	})

	Context("Get Cost Centers", func() {
		It("add a new cost center without a session token", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("add a new cost center with an invalid session token", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", "blah blah blah")

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("add a new cost center with an invalid name", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", adminToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("add a new cost center with viewer role", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", viewerToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a new cost center with non-cost-center role", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", nonCostCenterToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a cost center with multiple non-cost-center roles", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set(
				"Authorization",
				multipleNonCostCenterRolesToken,
			)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a new cost center with an admin session token", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", adminToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var addCostCenterResp vetchi.AddCostCenterResponse
			err = json.NewDecoder(resp.Body).Decode(&addCostCenterResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(addCostCenterResp.Name).Should(Equal("CC1-Admin"))
		})

		It("add a cost center with a duplicate name as Admin", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", adminToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusConflict))
		})

		It("add a cost center with cost-centers-crud session token", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC2-Crud",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", crud1Token)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var addCostCenterResp vetchi.AddCostCenterResponse
			err = json.NewDecoder(resp.Body).Decode(&addCostCenterResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(addCostCenterResp.Name).Should(Equal("CC2-Crud"))
		})

		It("add a cost center with a duplicate name as Crud", func() {
			addCostCenterReqBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC2-Crud",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", crud2Token)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusConflict))
		})

		It("should return the cost centers for the employer", func() {
			getCostCentersReqBody, err := json.Marshal(
				vetchi.GetCostCentersRequest{},
			)
			Expect(err).ShouldNot(HaveOccurred())

			// First try without a session token
			getCostCentersReq1, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp1, err := http.DefaultClient.Do(getCostCentersReq1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp1.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Then try with an invalid session token
			getCostCentersReq2, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq2.Header.Set("Authorization", "blah blah blah")

			resp2, err := http.DefaultClient.Do(getCostCentersReq2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp2.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Then try with a valid session token
			getCostCentersReq3, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq3.Header.Set("Authorization", adminToken)

			resp3, err := http.DefaultClient.Do(getCostCentersReq3)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp3.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp3.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(2))
			Expect(
				[]string{costCenters[0].Name, costCenters[1].Name},
			).Should(ConsistOf("CC1-Admin", "CC2-Crud"))
		})

		It("defunct a cost center with no session token", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("defunct a cost center with an invalid session token", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", "blah blah blah")

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("defunct a cost center with a viewer session token", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", viewerToken)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("defunct a cost center with a non-cost-center session", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", nonCostCenterToken)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("defunct a cc with a multiple-non-cost-center session", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set(
				"Authorization",
				multipleNonCostCenterRolesToken,
			)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("defunct a cost center with an admin session token", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
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

			// Get the cost centers and verify that the cost center is defunct
			getCostCentersReqBody, err := json.Marshal(
				vetchi.GetCostCentersRequest{},
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq.Header.Set("Authorization", adminToken)

			resp, err = http.DefaultClient.Do(getCostCentersReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(1))
			Expect(costCenters[0].Name).Should(Equal("CC2-Crud"))
		})

		It("defunct cc with a cost-centers-crud session token", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC2-Crud",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", crud1Token)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get the cost centers and verify that the cost center is defunct
			getCostCentersReqBody, err := json.Marshal(
				vetchi.GetCostCentersRequest{},
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq.Header.Set("Authorization", adminToken)

			resp, err = http.DefaultClient.Do(getCostCentersReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(0))
		})

		It("defunct a cost center with an invalid name", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "",
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
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "respBody: %+v\n", respBody)
			Expect(respBody.Errors).Should(HaveLen(1))
			Expect(respBody.Errors[0]).Should(ContainSubstring("name"))
		})

		It("defunct a cost center with a name that doesn't exist", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "Non-existent Cost Center",
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
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
		})

		It("defunct a cost center with a name that is too long", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: strings.Repeat("A", 65),
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
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(respBody.Errors).Should(HaveLen(1))
			Expect(respBody.Errors[0]).Should(ContainSubstring("name"))
		})

		It("defunct a cost center with a name that is too short", func() {
			defunctCostCenterReqBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "A",
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
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(respBody.Errors).Should(HaveLen(1))
			Expect(respBody.Errors[0]).Should(ContainSubstring("name"))
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
				[]string{costCenters[0].Name, costCenters[1].Name},
			).Should(ConsistOf("CC1-Admin", "CC2-Crud"))
		})

		It("create/defunct cost centers in bulk; verify pagination", func() {
			// runID parameter is needed because even if a CC is defunct,
			// its name cannot be used by another CC

			fmt.Fprintf(GinkgoWriter, "count is not divisible by limit\n")
			bulkAddDefunctCC(adminToken, serverURL, "run-1", 30, 4)

			fmt.Fprintf(GinkgoWriter, "count is divisible by limit\n")
			bulkAddDefunctCC(adminToken, serverURL, "run-2", 32, 4)

			fmt.Fprintf(GinkgoWriter, "count is less than limit\n")
			bulkAddDefunctCC(adminToken, serverURL, "run-3", 2, 4)
		})
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

		ccb4Update, err := getCostCenter(
			adminToken,
			vetchi.GetCostCenterRequest{
				Name: "CC-update-test-1",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ccb4Update.Name).Should(Equal("CC-update-test-1"))
		Expect(ccb4Update.Notes).Should(Equal("This is a test cost center"))

		fmt.Fprintf(GinkgoWriter, "Update cost center with invalid name\n")
		statusCode, err = updateCostCenter(
			adminToken,
			vetchi.UpdateCostCenterRequest{
				Name: "Non-existent Cost Center",
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(statusCode).Should(Equal(http.StatusNotFound))

		fmt.Fprintf(GinkgoWriter, "Updating cost center with invalid notes\n")
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

		ccAfterUpdate, err := getCostCenter(
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
})

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
) (vetchi.CostCenter, error) {
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
	Expect(resp.StatusCode).Should(Equal(http.StatusOK))

	var costCenter vetchi.CostCenter
	err = json.NewDecoder(resp.Body).Decode(&costCenter)
	Expect(err).ShouldNot(HaveOccurred())

	return costCenter, nil
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
				Name: ccName,
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
				PaginationKey: paginationKey,
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
			gotCC = append(gotCC, costCenter.Name)
		}
		paginationKey = costCenters[len(costCenters)-1].Name

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
				Name: ccName,
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
