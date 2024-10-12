package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

var _ = Describe("GetOnboardStatus", func() {
	var db *pgxpool.Pool
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		db = setupTestDB()
		_, err := db.Exec(
			ctx,
			`
INSERT INTO employers (client_id, onboard_status) VALUES
	('domain-onboarded.example', 'DOMAIN_ONBOARDED')`,
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		_, err := db.Exec(
			ctx,
			`
DELETE FROM employers WHERE client_id IN('domain-onboarded.example')`,
		)
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
	})
})
