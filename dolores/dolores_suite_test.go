package dolores

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDolores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dolores Suite")
}

const serverURL = "http://localhost:8081"

var db *pgxpool.Pool

var _ = BeforeSuite(func() {
	db = setupTestDB()

	/*
		seed, err := os.ReadFile("seed.pgsql")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = db.Exec(context.Background(), string(seed))
		Expect(err).ShouldNot(HaveOccurred())
	*/
})
