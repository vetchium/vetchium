package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
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

		It("check mailpit if it got the email", func() {
			// Sleep for 6 minutes to allow the email to be sent by granger
			<-time.After(3 * time.Minute)

			queryParams := url.Values{}
			queryParams.Add("to", "secretsapp@example.com")
			queryParams.Add("subject", "Welcome to Vetchi !")
			qpStr := queryParams.Encode()

			url := serverURL + "/employer/get-onboard-status?" + qpStr

			req, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			// Extracting the URL from the body string
			re := regexp.MustCompile(
				`https://employer.vetchi.org/onboard/[^"]+`,
			)
			urls := re.FindAllString(string(body), -1)
			Expect(len(urls)).Should(BeNumerically(">=", 1))

			for _, url := range urls {
				log.Println("URL:", url)
			}
		})
	})
})
