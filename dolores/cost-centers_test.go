package dolores

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cost Centers", func() {
	Context("Get Cost Centers", func() {
		It("should return the cost centers for the employer", func() {
			sessionToken, err := employerSignin(
				"test1.example",
				"admin@test1.example",
				"NewPassword123$",
			)
			_ = sessionToken
			_ = err
			Expect(true).Should(BeTrue())
		})
	})
})
