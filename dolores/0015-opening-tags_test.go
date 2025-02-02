package dolores

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const defaultPassword = "NewPassword123$"

var _ = Describe("OpeningsTags", Ordered, func() {
	var (
		token string
		db    *pgxpool.Pool
	)

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0015-opening-tags-up.pgsql")
		var err error
		token, err = employerSignin(
			"openingtags.example",
			"tags.test@openingtags.example",
			defaultPassword,
		)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		seedDatabase(db, "0015-opening-tags-down.pgsql")
		db.Close()
	})

	Describe("Filter Opening Tags", func() {
		It("should return all tags when no prefix is provided", func() {
			req := common.FilterOpeningTagsRequest{}
			resp := doPOST(
				token,
				req,
				"/employer/filter-opening-tags",
				http.StatusOK,
				true,
			)
			var result []common.OpeningTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(25))
			Expect(result[0].Name).To(Equal("Backend Developer"))
			Expect(result[9].Name).To(Equal("Go"))
			Expect(result[11].Name).To(Equal("Java"))
			Expect(result[14].Name).To(Equal("PostgreSQL"))
			Expect(result[16].Name).To(Equal("Python"))
			Expect(result[18].Name).To(Equal("React"))
		})

		It("should filter tags by prefix", func() {
			prefix := "P"
			req := common.FilterOpeningTagsRequest{Prefix: &prefix}
			resp := doPOST(
				token,
				req,
				"/employer/filter-opening-tags",
				http.StatusOK,
				true,
			)
			var result []common.OpeningTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(3))
			Expect(result[0].Name).To(Equal("PostgreSQL"))
			Expect(result[1].Name).To(Equal("Product Manager"))
			Expect(result[2].Name).To(Equal("Python"))
		})

		It("should return empty list for non-existent prefix", func() {
			prefix := "NonExistent"
			req := common.FilterOpeningTagsRequest{Prefix: &prefix}
			resp := doPOST(
				token,
				req,
				"/employer/filter-opening-tags",
				http.StatusOK,
				true,
			)
			var result []common.OpeningTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("Create Opening with Tags", func() {
		It("should create opening with existing tags", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Existing Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers with Python and Go experience",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				Tags: []common.OpeningTagID{
					"12345678-0015-0015-0015-000000070003", // Python
					"12345678-0015-0015-0015-000000070001", // Go
				},
			}

			testPOST(token, req, "/employer/create-opening", http.StatusOK)
		})

		It("should create opening with new tags", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with New Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers with Rust and TypeScript experience",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				NewTags: []string{"Scala", "Haskell"},
			}

			testPOST(token, req, "/employer/create-opening", http.StatusOK)
		})

		It("should create opening with both existing and new tags", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Mixed Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers with Python and Swift experience",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				Tags: []common.OpeningTagID{
					"12345678-0015-0015-0015-000000070003", // Python
				},
				NewTags: []string{"Swift"},
			}

			testPOST(token, req, "/employer/create-opening", http.StatusOK)
		})

		It("should fail when no tags are provided", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with No Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})

		It("should fail when more than 3 tags are provided", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Too Many Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				Tags: []common.OpeningTagID{
					"12345678-0015-0015-0015-000000070001", // Go
					"12345678-0015-0015-0015-000000070002", // Java
					"12345678-0015-0015-0015-000000070003", // Python
					"12345678-0015-0015-0015-000000070004", // PostgreSQL
				},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})

		It("should fail when more than 3 new tags are provided", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Too Many New Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				NewTags: []string{"Tag1", "Tag2", "Tag3", "Tag4"},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})

		It("should fail when combined tags exceed 3", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Too Many Combined Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers",
				Recruiter:         "tags.test@openingtags.example",
				HiringManager:     "tags.test@openingtags.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            5,
				MinEducationLevel: common.NotMattersEducation,
				LocationTitles:    []string{"Main Office"},
				RemoteCountryCodes: []common.CountryCode{
					"IND",
					"USA",
				},
				Tags: []common.OpeningTagID{
					"12345678-0015-0015-0015-000000070001", // Go
					"12345678-0015-0015-0015-000000070002", // Java
				},
				NewTags: []string{"Tag1", "Tag2"},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})
	})
})
