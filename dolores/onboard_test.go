package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

var _ = Describe("GetOnboardStatus", func() {

	type Message struct {
		ID string `json:"ID"`
	}

	type MailPitResponse struct {
		Messages []Message `json:"messages"`
	}

	var db *pgxpool.Pool
	var ctx context.Context
	var employerID, domainID, orgUserID int64

	BeforeEach(func() {
		ctx = context.Background()
		db = setupTestDB()

		err := db.QueryRow(
			ctx,
			`
INSERT INTO employers	(client_id_type, employer_state, 
						onboard_admin_email, onboard_secret_token)
VALUES ('DOMAIN', 'ONBOARDED', 'admin@domain-onboarded.example', 'token') 
RETURNING id
`,
		).Scan(&employerID)
		Expect(err).ShouldNot(HaveOccurred())

		err = db.QueryRow(
			ctx,
			`
INSERT INTO domains (domain_name, domain_state, employer_id) 
VALUES ('domain-onboarded.example', 'VERIFIED', $1) 
RETURNING id`,
			employerID,
		).Scan(&domainID)
		Expect(err).ShouldNot(HaveOccurred())

		err = db.QueryRow(
			ctx,
			`
INSERT INTO org_users (email, password_hash, org_user_role, employer_id)
VALUES ('admin@domain-onboarded.example', 'password_hash', 'ADMIN', $1)
RETURNING id`,
			employerID,
		).Scan(&orgUserID)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		_, err := db.Exec(ctx, `DELETE FROM org_users WHERE id = $1`, orgUserID)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(ctx, `DELETE FROM domains WHERE id = $1`, domainID)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(ctx, `DELETE FROM employers WHERE id = $1`, employerID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(err).ShouldNot(HaveOccurred())
		db.Close()
	})

	var _ = Describe("GetOnboardStatus", func() {
		It("returns the onboard status", func() {
			var tests = []struct {
				clientID string
				want     libvetchi.OnboardStatus
			}{
				{
					clientID: "domain-onboarded.example",
					want:     libvetchi.DomainOnboarded,
				},
				{
					clientID: "example.com",
					want:     libvetchi.DomainNotVerified,
				},
				{
					clientID: "secretsapp.com",
					want:     libvetchi.DomainVerifiedOnboardPending,
				},
			}

			for _, test := range tests {
				log.Println("Testing for domain", test.clientID)
				getOnboardStatusRequest := libvetchi.GetOnboardStatusRequest{
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

				var got libvetchi.GetOnboardStatusResponse
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

			// TODO: Once password validation is added, add a testcase with
			// an invalid password

			// Set password for the admin
			setOnboardPasswordBody, err := json.Marshal(
				libvetchi.SetOnboardPasswordRequest{
					ClientID: "domain-onboarded.example",
					Password: "NewPassword123$",
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
			).Should(Equal(http.StatusOK))

			// Get Onboard Status should now return DomainOnboarded
			getOnboardStatusRequest := libvetchi.GetOnboardStatusRequest{
				ClientID: "domain-onboarded.example",
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

			var got libvetchi.GetOnboardStatusResponse
			err = json.NewDecoder(getOnboardStatusResp.Body).Decode(&got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(got.Status).Should(Equal(libvetchi.DomainOnboarded))

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
			getOnboardStatusRequest := libvetchi.GetOnboardStatusRequest{
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

			var got libvetchi.GetOnboardStatusResponse
			err = json.NewDecoder(resp.Body).Decode(&got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				got.Status,
			).Should(Equal(libvetchi.DomainVerifiedOnboardPending))

			// Sleep for 3 minutes to allow granger to email the token
			<-time.After(3 * time.Minute)

			// Sleep for 2 minutes to allow the email to be sent by granger
			<-time.After(2 * time.Minute)

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

			setPasswordRequest := libvetchi.SetOnboardPasswordRequest{
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
