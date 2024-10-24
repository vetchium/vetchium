package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = Describe("Cost Centers", Ordered, func() {
	var db *pgxpool.Pool
	var sessionToken string
	BeforeAll(func() {
		db = setupTestDB()

		seed, err := os.ReadFile("0002-cost-centers-up.pgsql")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(context.Background(), string(seed))
		Expect(err).ShouldNot(HaveOccurred())

		sessionToken, err = employerSignin(
			"cost-center.example",
			"admin@cost-center.example",
			"NewPassword123$",
		)
		Expect(err).ShouldNot(HaveOccurred())
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
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("add a new cost center with an invalid session token", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", "blah blah blah")

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("add a new cost center with an invalid name", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("add a new cost center with viewer role", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			sessionToken2, err := employerSignin(
				"cost-center.example",
				"viewer@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken2)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a new cost center with non-cost-center role", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			sessionToken3, err := employerSignin(
				"cost-center.example",
				"non-cost-center@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken3)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a cost center with multiple non-cost-center roles", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			sessionToken4, err := employerSignin(
				"cost-center.example",
				"multiple-non-cost-center-roles@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken4)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusForbidden))
		})

		It("add a new cost center with an admin session token", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var addCostCenterResp vetchi.AddCostCenterResponse
			err = json.NewDecoder(resp.Body).Decode(&addCostCenterResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(addCostCenterResp.Name).Should(Equal("CC1-Admin"))
		})

		It("add a cost center with a duplicate name as Admin", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusConflict))
		})

		It("add a cost center with cost-centers-crud session token", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC2-Crud",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			sessionToken5, err := employerSignin(
				"cost-center.example",
				"crud1@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken5)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var addCostCenterResp vetchi.AddCostCenterResponse
			err = json.NewDecoder(resp.Body).Decode(&addCostCenterResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(addCostCenterResp.Name).Should(Equal("CC2-Crud"))
		})

		It("add a cost center with a duplicate name as Crud", func() {
			addCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "CC2-Crud",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(addCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			sessionToken6, err := employerSignin(
				"cost-center.example",
				"crud2@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			addCostCenterReq.Header.Set("Authorization", sessionToken6)

			resp, err := http.DefaultClient.Do(addCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusConflict))
		})

		It("should return the cost centers for the employer", func() {
			getCostCentersRequestBody, err := json.Marshal(
				vetchi.GetCostCentersRequest{},
			)
			Expect(err).ShouldNot(HaveOccurred())

			// First try without a session token
			getCostCentersReq1, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp1, err := http.DefaultClient.Do(getCostCentersReq1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp1.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Then try with an invalid session token
			getCostCentersReq2, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/get-cost-centers",
				bytes.NewBuffer(getCostCentersRequestBody),
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
				bytes.NewBuffer(getCostCentersRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			getCostCentersReq3.Header.Set("Authorization", sessionToken)

			resp3, err := http.DefaultClient.Do(getCostCentersReq3)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp3.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp3.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(2))

			Expect(costCenters[0].Name).Should(Equal("CC1-Admin"))
			Expect(costCenters[1].Name).Should(Equal("CC2-Crud"))
		})

		It("defunct a cost center with no session token", func() {
			defunctCostCenterRequestBody, err := json.Marshal(
				vetchi.DefunctCostCenterRequest{
					Name: "CC1-Admin",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/defunct-cost-center",
				bytes.NewBuffer(defunctCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		// TODO: Add tests for pagination

	})
})
