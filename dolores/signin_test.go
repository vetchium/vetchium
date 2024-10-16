package dolores

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Signin", func() {
	Describe("Employer Signin", func() {
		It("various invalid inputs", func() {
			signinReqBody, err := json.Marshal(
				vetchi.EmployerSignInRequest{
					ClientID: "test",
					Email:    "test",
					Password: "test",
				})
			Expect(err).ShouldNot(HaveOccurred())

			signinReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/signin",
				bytes.NewBuffer(signinReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			signinReq.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(signinReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			var validationErr vetchi.ValidationErrors
			err = json.NewDecoder(resp.Body).Decode(&validationErr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(validationErr.Errors).Should(HaveLen(3))
			Expect(validationErr.Errors).Should(ContainElements([]string{
				"client_id", "email", "password",
			}))
		})

		It("non-existent client_id", func() {
			signinReqBody, err := json.Marshal(
				vetchi.EmployerSignInRequest{
					ClientID: "bad-client-id.example",
					Email:    "admin@domain-onboarded.example",
					Password: "NewPassword123$",
				})
			Expect(err).ShouldNot(HaveOccurred())

			signinReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/signin",
				bytes.NewBuffer(signinReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			signinReq.Header.Set("Content-Type", "application/json")

			signinResp, err := http.DefaultClient.Do(signinReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signinResp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("good client_id, non-existent email, bad password", func() {
			signinReqBody, err := json.Marshal(
				vetchi.EmployerSignInRequest{
					ClientID: "domain-onboarded.example",
					Email:    "non-existent-email@domain-onboarded.example",
					Password: "BadPassword11234$",
				})
			Expect(err).ShouldNot(HaveOccurred())

			signinReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/signin",
				bytes.NewBuffer(signinReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			signinReq.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(signinReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("good client_id, good email, bad password", func() {
			signinReqBody, err := json.Marshal(
				vetchi.EmployerSignInRequest{
					ClientID: "domain-onboarded.example",
					Email:    "admin@domain-onboarded.example",
					Password: "BadPassword11234$",
				})
			Expect(err).ShouldNot(HaveOccurred())

			signinReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/signin",
				bytes.NewBuffer(signinReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			signinReq.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(signinReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("good client_id, good email, good password", func() {
			signinReqBody, err := json.Marshal(
				vetchi.EmployerSignInRequest{
					ClientID: "domain-onboarded.example",
					Email:    "admin@domain-onboarded.example",
					Password: "NewPassword123$",
				})
			Expect(err).ShouldNot(HaveOccurred())

			signinReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/signin",
				bytes.NewBuffer(signinReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())

			signinReq.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(signinReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
	})
})
