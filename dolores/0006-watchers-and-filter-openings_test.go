package dolores

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var bachelorEducation_0006 *vetchi.EducationLevel

var _ = Describe("Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOpeningsToken string
	var recruiterToken, hiringManagerToken string

	bachelor := vetchi.BachelorEducation
	bachelorEducation_0006 = &bachelor

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0006-watchers-and-filter-openings-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@openings0006.example":          &adminToken,
			"crud@openings0006.example":           &crudToken,
			"viewer@openings0006.example":         &viewerToken,
			"non-openings@openings0006.example":   &nonOpeningsToken,
			"recruiter@openings0006.example":      &recruiterToken,
			"hiring-manager@openings0006.example": &hiringManagerToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"openings0006.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0006-watchers-and-filter-openings-down.pgsql")
		db.Close()
	})

	Describe("Filter Openings", func() {
		Expect(true).To(BeTrue())
	})

	Describe("Watchers", func() {

	})
})
