package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

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

			addCostCenterReq.Header.Set("Authorization", sessionToken)

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

			addCostCenterReq.Header.Set("Authorization", sessionToken)

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

			sessionToken7, err := employerSignin(
				"cost-center.example",
				"viewer@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", sessionToken7)

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

			sessionToken8, err := employerSignin(
				"cost-center.example",
				"non-cost-center@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", sessionToken8)

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

			sessionToken9, err := employerSignin(
				"cost-center.example",
				"multiple-non-cost-center-roles@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", sessionToken9)

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

			defunctCostCenterReq.Header.Set("Authorization", sessionToken)

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

			getCostCentersReq.Header.Set("Authorization", sessionToken)

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

			sessionToken10, err := employerSignin(
				"cost-center.example",
				"crud3@cost-center.example",
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			defunctCostCenterReq.Header.Set("Authorization", sessionToken10)

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

			getCostCentersReq.Header.Set("Authorization", sessionToken)

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

			defunctCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(respBody).Should(HaveLen(1))
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

			defunctCostCenterReq.Header.Set("Authorization", sessionToken)

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

			defunctCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(respBody).Should(HaveLen(1))
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

			defunctCostCenterReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(defunctCostCenterReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var respBody vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(respBody).Should(HaveLen(1))
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

			getDefunctCostCentersReq.Header.Set("Authorization", sessionToken)

			resp, err := http.DefaultClient.Do(getDefunctCostCentersReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var costCenters []vetchi.CostCenter
			err = json.NewDecoder(resp.Body).Decode(&costCenters)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(costCenters).Should(HaveLen(2))
			Expect(costCenters).Should(ContainElements("CC1-Admin", "CC2-Crud"))
		})

		// TODO: Add tests for pagination
	})
})
