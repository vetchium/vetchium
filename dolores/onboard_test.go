package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = Describe("GetOnboardStatus", func() {
	type Message struct {
		ID string `json:"ID"`
	}

	type MailPitResponse struct {
		Messages []Message `json:"messages"`
	}

	var _ = Describe("GetOnboardStatus", func() {
		It("returns the onboard status", func() {
			var tests = []struct {
				clientID string
				want     vetchi.OnboardStatus
			}{
				{
					clientID: "domain-onboarded.example",
					want:     vetchi.DomainOnboarded,
				},
				{
					clientID: "example.com",
					want:     vetchi.DomainNotVerified,
				},
				{
					clientID: "secretsapp.com",
					want:     vetchi.DomainVerifiedOnboardPending,
				},
			}

			for _, test := range tests {
				fmt.Fprintf(
					GinkgoWriter,
					"Testing for domain %s\n",
					test.clientID,
				)
				getOnboardStatusRequest := vetchi.GetOnboardStatusRequest{
					ClientID: test.clientID,
				}

				req, err := json.Marshal(getOnboardStatusRequest)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/employer/get-onboard-status",
					"application/json",
					bytes.NewBuffer(req),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				var got vetchi.GetOnboardStatusResponse
				err = json.NewDecoder(resp.Body).Decode(&got)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(got.Status).Should(Equal(test.want))
			}
		})

		It("Check if mailpit got the email and set the admin password", func() {
			// Sleep for 2 minutes to allow the email to be sent by granger
			<-time.After(2 * time.Minute)

			url := "http://localhost:8025/api/v1/search?query=to%3Asecretsapp%40example.com%20subject%3AWelcome%20to%20Vetchi%20!"
			log.Println("URL:", url)

			listMailsReq, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())
			listMailsReq.Header.Add("Content-Type", "application/json")

			listMailsResp, err := http.DefaultClient.Do(listMailsReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listMailsResp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(listMailsResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			log.Println("Body:", string(body))

			var listMailsRespObj MailPitResponse
			err = json.Unmarshal(body, &listMailsRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				len(listMailsRespObj.Messages),
			).Should(BeNumerically(">=", 1))

			mailURL := "http://localhost:8025/api/v1/message/" + listMailsRespObj.Messages[0].ID
			log.Println("Mail URL:", mailURL)

			getMailReq, err := http.NewRequest("GET", mailURL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			getMailReq.Header.Add("Content-Type", "application/json")

			getMailResp, err := http.DefaultClient.Do(getMailReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getMailResp.StatusCode).Should(Equal(http.StatusOK))

			body, err = io.ReadAll(getMailResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			log.Println("Mail Body:", string(body))

			// Extracting the token from the URL
			re := regexp.MustCompile(
				`https://employer.vetchi.org/onboard/([^\\\s]+)`,
			)
			tokens := re.FindAllStringSubmatch(string(body), -1)
			Expect(len(tokens)).Should(BeNumerically(">=", 1))

			token := tokens[0][1] // The token is captured in the first group
			log.Println("Token:", token)

			// Test with an invalid password
			setOnboardPasswordBody, err := json.Marshal(
				vetchi.SetOnboardPasswordRequest{
					ClientID: "secretsapp.com",
					Password: "pass",
					Token:    token,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			setOnboardPasswordResp, err := http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setOnboardPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				setOnboardPasswordResp.StatusCode,
			).Should(Equal(http.StatusBadRequest))

			var v vetchi.ValidationErrors
			err = json.NewDecoder(setOnboardPasswordResp.Body).Decode(&v)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v.Errors).ShouldNot(BeEmpty())
			Expect(v.Errors).Should(ContainElement("password"))

			// Set password for the admin
			setOnboardPasswordBody, err = json.Marshal(
				vetchi.SetOnboardPasswordRequest{
					ClientID: "secretsapp.com",
					Password: "NewPassword123$",
					Token:    token,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			setOnboardPasswordResp, err = http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setOnboardPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				setOnboardPasswordResp.StatusCode,
			).Should(Equal(http.StatusOK))

			// Get Onboard Status should now return DomainOnboarded
			getOnboardStatusRequest := vetchi.GetOnboardStatusRequest{
				ClientID: "secretsapp.com",
			}
			getOnboardStatusBody, err := json.Marshal(getOnboardStatusRequest)
			Expect(err).ShouldNot(HaveOccurred())

			getOnboardStatusResp, err := http.Post(
				serverURL+"/employer/get-onboard-status",
				"application/json",
				bytes.NewBuffer(getOnboardStatusBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getOnboardStatusResp.StatusCode).Should(Equal(http.StatusOK))

			var got vetchi.GetOnboardStatusResponse
			err = json.NewDecoder(getOnboardStatusResp.Body).Decode(&got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(got.Status).Should(Equal(vetchi.DomainOnboarded))

			log.Println("Test if the same token can be used again")

			// Retry the set-password with the same token
			setOnboardPasswordResp2, err := http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setOnboardPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				setOnboardPasswordResp2.StatusCode,
			).Should(Equal(http.StatusUnprocessableEntity))
		})

		It("test if invite token can be used after validity", func() {
			getOnboardStatusRequest := vetchi.GetOnboardStatusRequest{
				ClientID: "aadal.in",
			}

			req, err := json.Marshal(getOnboardStatusRequest)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/employer/get-onboard-status",
				"application/json",
				bytes.NewBuffer(req),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var got vetchi.GetOnboardStatusResponse
			err = json.NewDecoder(resp.Body).Decode(&got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				got.Status,
			).Should(Equal(vetchi.DomainVerifiedOnboardPending))

			fmt.Fprintf(GinkgoWriter, "Sleeping to allow granger to email\n")
			<-time.After(3 * time.Minute)
			fmt.Fprintf(GinkgoWriter, "Wokeup\n")

			url := "http://localhost:8025/api/v1/search?query=to%3Aaadal%40example.com%20subject%3AWelcome%20to%20Vetchi%20!"
			log.Println("URL:", url)

			mailPitReq1, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq1.Header.Add("Content-Type", "application/json")

			mailPitResp1, err := http.DefaultClient.Do(mailPitReq1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp1.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailPitResp1.Body)
			Expect(err).ShouldNot(HaveOccurred())

			log.Println("Body:", string(body))

			var mailPitResp1Obj MailPitResponse
			err = json.Unmarshal(body, &mailPitResp1Obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(mailPitResp1Obj.Messages)).Should(BeNumerically(">=", 1))

			mailURL := "http://localhost:8025/api/v1/message/" + mailPitResp1Obj.Messages[0].ID
			log.Println("Mail URL:", mailURL)

			mailPitReq2, err := http.NewRequest("GET", mailURL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			mailPitReq2.Header.Add("Content-Type", "application/json")

			mailPitResp2, err := http.DefaultClient.Do(mailPitReq2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailPitResp2.StatusCode).Should(Equal(http.StatusOK))

			body, err = io.ReadAll(mailPitResp2.Body)
			Expect(err).ShouldNot(HaveOccurred())

			log.Println("Mail Body:", string(body))

			// Extracting the token from the URL
			re := regexp.MustCompile(
				`https://employer.vetchi.org/onboard/([^\\\s]+)`,
			)
			tokens := re.FindAllStringSubmatch(string(body), -1)
			Expect(len(tokens)).Should(BeNumerically(">=", 1))

			token := tokens[0][1] // The token is captured in the first group
			log.Println("Token:", token)

			// Sleep to allow the token to expire
			<-time.After(4 * time.Minute)

			setPasswordRequest := vetchi.SetOnboardPasswordRequest{
				ClientID: "aadal.in",
				Password: "NewPassword123$",
				Token:    token,
			}

			setPasswordBody, err := json.Marshal(setPasswordRequest)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err = http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				resp.StatusCode,
			).Should(Equal(http.StatusUnprocessableEntity))
		})
	})
})
