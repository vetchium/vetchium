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
		It("create a new cost center without a session token", func() {
			createCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			createCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(createCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(createCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("create a new cost center with a session token", func() {
			createCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "New Cost Center",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			createCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(createCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			createCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(createCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var createCostCenterResp vetchi.AddCostCenterResponse
			err = json.NewDecoder(resp.Body).Decode(&createCostCenterResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createCostCenterResp.CostCenterName).ShouldNot(BeEmpty())
		})

		It("create a new cost center with an invalid name", func() {
			createCostCenterRequestBody, err := json.Marshal(
				vetchi.AddCostCenterRequest{
					Name: "",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			createCostCenterReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/add-cost-center",
				bytes.NewBuffer(createCostCenterRequestBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			createCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(createCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
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
			Expect(resp3.StatusCode).Should(Equal(http.StatusNotImplemented))
		})
	})
})
