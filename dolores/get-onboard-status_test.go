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

			req, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			log.Println("Body:", string(body))

			type Message struct {
				ID string `json:"ID"`
			}

			type MailPitResponse struct {
				Messages []Message `json:"messages"`
			}

			var response MailPitResponse
			err = json.Unmarshal(body, &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(response.Messages)).Should(BeNumerically(">=", 1))

			mailURL := "http://localhost:8025/api/v1/message/" + response.Messages[0].ID
			log.Println("Mail URL:", mailURL)

			req, err = http.NewRequest("GET", mailURL, nil)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Add("Content-Type", "application/json")

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err = io.ReadAll(resp.Body)
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

			resp, err = http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setOnboardPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get Onboard Status should now return DomainOnboarded
			getOnboardStatusRequest := libvetchi.GetOnboardStatusRequest{
				ClientID: "domain-onboarded.example",
			}
			getOnboardStatusBody, err := json.Marshal(getOnboardStatusRequest)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err = http.Post(
				serverURL+"/employer/get-onboard-status",
				"application/json",
				bytes.NewBuffer(getOnboardStatusBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var got libvetchi.GetOnboardStatusResponse
			err = json.NewDecoder(resp.Body).Decode(&got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(got.Status).Should(Equal(libvetchi.DomainOnboarded))

			// Retry the set-password with the same token
			resp, err = http.Post(
				serverURL+"/employer/set-onboard-password",
				"application/json",
				bytes.NewBuffer(setOnboardPasswordBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				resp.StatusCode,
			).Should(Equal(http.StatusUnprocessableEntity))
		})
	})
})
