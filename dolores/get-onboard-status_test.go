package dolores

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

var _ = Describe("GetOnboardStatus", func() {
	It("returns the onboard status", func() {
		var tests = []struct {
			clientID string
			want     libvetchi.OnboardStatus
		}{
			{clientID: "notfound.example", want: libvetchi.DomainNotVerified},
			{
				clientID: "domain-verified-email-not-sent.example",
				want:     libvetchi.DomainVerifiedEmailNotSent,
			},
			{
				clientID: "domain-verified-email-sent.example",
				want:     libvetchi.DomainVerifiedEmailSent,
			},
			{
				clientID: "domain-onboarded.example",
				want:     libvetchi.DomainOnboarded,
			},
		}

		for _, test := range tests {
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
