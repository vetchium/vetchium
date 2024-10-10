package dolores

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

var _ = Describe("GetOnboardStatus", func() {
	Context("when the client ID is valid", func() {
		req := libvetchi.GetOnboardStatusRequest{
			ClientID: "123",
		}

		It("returns the onboard status", func() {
			resp, err := http.Post(
				serverURL+"/employer/get-onboard-status",
				"application/json",
				bytes.NewBuffer(req),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
	})
})
