package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = Describe("Employer Signin", func() {
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

			var signinResp vetchi.EmployerSignInResponse
			err = json.NewDecoder(resp.Body).Decode(&signinResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signinResp.Token).ShouldNot(BeEmpty())

			// Get the tfa code from the email by querying mailpit
			fmt.Fprintf(GinkgoWriter, "Sleeping to allow granger to email\n")
			<-time.After(2 * time.Minute)
			fmt.Fprintf(GinkgoWriter, "Wokeup\n")

			url := "http://localhost:8025/api/v1/search?query=to%3Aadmin%40domain-onboarded.example%20subject%3AVetchi%20Two%20Factor%20Authentication"
			fmt.Fprintf(GinkgoWriter, "URL: %s\n", url)

			mailPitReq1, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq1.Header.Add("Content-Type", "application/json")

			mailPitResp1, err := http.DefaultClient.Do(mailPitReq1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp1.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailPitResp1.Body)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Fprintf(GinkgoWriter, "Body: %s\n", string(body))

			var mailPitResp1Obj MailPitResponse
			err = json.Unmarshal(body, &mailPitResp1Obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(mailPitResp1Obj.Messages)).Should(Equal(1))

			mailURL := "http://localhost:8025/api/v1/message/" + mailPitResp1Obj.Messages[0].ID
			fmt.Fprintf(GinkgoWriter, "Mail URL: %s\n", mailURL)

			mailPitReq2, err := http.NewRequest("GET", mailURL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq2.Header.Add("Content-Type", "application/json")

			mailPitResp2, err := http.DefaultClient.Do(mailPitReq2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp2.StatusCode).Should(Equal(http.StatusOK))

			body, err = io.ReadAll(mailPitResp2.Body)
			Expect(err).ShouldNot(HaveOccurred())

			// Extracting the token from the mail body
			re := regexp.MustCompile(`Token:\s*([a-zA-Z0-9]+)\s*`)

			tokens := re.FindAllStringSubmatch(string(body), -1)
			Expect(len(tokens)).Should(BeNumerically(">=", 1))

			tfaCode := tokens[0][1] // The token is captured in the first group
			fmt.Fprintf(GinkgoWriter, "TFACode: %s\n", tfaCode)
			fmt.Fprintf(GinkgoWriter, "TGT: %s\n", signinResp.Token)

			// TFA with the two tokens
			tfaReqBody, err := json.Marshal(
				vetchi.EmployerTFARequest{
					TGT:     signinResp.Token,
					TFACode: tfaCode,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			tfaReq, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/employer/tfa",
				bytes.NewBuffer(tfaReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			tfaReq.Header.Set("Content-Type", "application/json")

			tfaResp, err := http.DefaultClient.Do(tfaReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tfaResp.StatusCode).Should(Equal(http.StatusOK))

			var tfaRespObj vetchi.EmployerTFAResponse
			err = json.NewDecoder(tfaResp.Body).Decode(&tfaRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tfaRespObj.SessionToken).ShouldNot(BeEmpty())
		})
	})
})
