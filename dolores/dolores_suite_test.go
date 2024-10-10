package dolores_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDolores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dolores Suite")
}

var serverURL = "http://localhost:8080"
