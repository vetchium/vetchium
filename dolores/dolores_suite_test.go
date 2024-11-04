package dolores

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDolores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dolores Suite")
}

var _ = BeforeSuite(func() {
	// Delete all mails on mailpit to clean up from past runs
	resp, err := http.DefaultClient.Get("http://localhost:8025/api/v1/messages")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(http.StatusOK))
})
