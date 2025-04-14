package dolores

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"

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
			req := common.FilterVTagsRequest{}
			resp := doPOST(
				token,
				req,
				"/employer/filter-vtags",
				http.StatusOK,
				true,
			)
			var result []common.VTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())

			// Check for existence of specific tags without enforcing positions
			tagNames := make([]string, len(result))
			for i, tag := range result {
				tagNames[i] = string(tag.Name)
			}
			Expect(tagNames).To(ContainElement("Backend Developer"))
			Expect(tagNames).To(ContainElement("Go"))
			Expect(tagNames).To(ContainElement("Java"))
			Expect(tagNames).To(ContainElement("PostgreSQL"))
			Expect(tagNames).To(ContainElement("Python"))
			Expect(tagNames).To(ContainElement("React"))
		})

		It("should filter tags by prefix", func() {
			prefix := "P"
			req := common.FilterVTagsRequest{Prefix: &prefix}
			resp := doPOST(
				token,
				req,
				"/employer/filter-vtags",
				http.StatusOK,
				true,
			)
			var result []common.VTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())

			got := make([]string, len(result))
			for i, tag := range result {
				got[i] = string(tag.Name)
			}

			Expect(got).To(ContainElement("PostgreSQL"))
			Expect(got).To(ContainElement("Product Manager"))
			Expect(got).To(ContainElement("Python"))
		})

		It("should return empty list for non-existent prefix", func() {
			prefix := "NonExistent"
			req := common.FilterVTagsRequest{Prefix: &prefix}
			resp := doPOST(
				token,
				req,
				"/employer/filter-vtags",
				http.StatusOK,
				true,
			)
			var result []common.VTag
			err := json.Unmarshal(resp.([]byte), &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("Create and Get Opening with Tags", func() {
		var openingID string

		It("should create opening with existing tags and verify them", func() {
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
				Tags: []common.VTagID{
					"12345678-0015-0015-0015-000000070003", // Python
					"12345678-0015-0015-0015-000000070001", // Go
				},
			}

			resp := doPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusOK,
				true,
			)
			var createResp employer.CreateOpeningResponse
			err := json.Unmarshal(resp.([]byte), &createResp)
			Expect(err).NotTo(HaveOccurred())
			openingID = createResp.OpeningID

			// Now get the opening and verify tags
			getReq := employer.GetOpeningRequest{ID: openingID}
			resp = doPOST(
				token,
				getReq,
				"/employer/get-opening",
				http.StatusOK,
				true,
			)
			var opening employer.Opening
			err = json.Unmarshal(resp.([]byte), &opening)
			Expect(err).NotTo(HaveOccurred())

			Expect(opening.Tags).To(HaveLen(2))

			got := make([]string, len(opening.Tags))
			for i, tag := range opening.Tags {
				got[i] = string(tag.Name)
			}
			sort.Strings(got)

			Expect(got[0]).To(Equal("Go"))
			Expect(got[1]).To(Equal("Python"))
		})

		It("should create opening with new tags and verify them", func() {
			newTags := []string{"Scala", "Haskell"}
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
				NewTags: newTags,
			}

			resp := doPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusOK,
				true,
			)
			var createResp employer.CreateOpeningResponse
			err := json.Unmarshal(resp.([]byte), &createResp)
			Expect(err).NotTo(HaveOccurred())
			openingID = createResp.OpeningID

			// Now get the opening and verify tags
			getReq := employer.GetOpeningRequest{ID: openingID}
			resp = doPOST(
				token,
				getReq,
				"/employer/get-opening",
				http.StatusOK,
				true,
			)
			var opening employer.Opening
			err = json.Unmarshal(resp.([]byte), &opening)
			Expect(err).NotTo(HaveOccurred())
			Expect(opening.Tags).To(HaveLen(2))

			got := make([]string, len(opening.Tags))
			for i, tag := range opening.Tags {
				got[i] = string(tag.Name)
			}
			sort.Strings(got)

			Expect(got[0]).To(BeElementOf(newTags))
			Expect(got[1]).To(BeElementOf(newTags))
		})

		It("should create opening with mixed tags and verify them", func() {
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
				Tags: []common.VTagID{
					"12345678-0015-0015-0015-000000070003", // Python
				},
				NewTags: []string{"Swift"},
			}

			resp := doPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusOK,
				true,
			)
			var createResp employer.CreateOpeningResponse
			err := json.Unmarshal(resp.([]byte), &createResp)
			Expect(err).NotTo(HaveOccurred())
			openingID = createResp.OpeningID

			// Now get the opening and verify tags
			getReq := employer.GetOpeningRequest{ID: openingID}
			resp = doPOST(
				token,
				getReq,
				"/employer/get-opening",
				http.StatusOK,
				true,
			)
			var opening employer.Opening
			err = json.Unmarshal(resp.([]byte), &opening)
			Expect(err).NotTo(HaveOccurred())

			Expect(opening.Tags).To(HaveLen(2))

			sortedTags := []string{
				string(opening.Tags[0].Name),
				string(opening.Tags[1].Name),
			}
			sort.Strings(sortedTags)
			Expect(sortedTags[0]).To(Equal("Python"))
			Expect(sortedTags[1]).To(Equal("Swift"))
		})

		It(
			"should create opening with existing tag passed as new tag and verify it",
			func() {
				req := employer.CreateOpeningRequest{
					Title:             "Go Developer",
					Positions:         2,
					JD:                "Looking for talented software engineers with Go experience",
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
					NewTags: []string{
						"Go",
					}, // Go is an existing tag but passed as new
				}

				resp := doPOST(
					token,
					req,
					"/employer/create-opening",
					http.StatusOK,
					true,
				)
				var createResp employer.CreateOpeningResponse
				err := json.Unmarshal(resp.([]byte), &createResp)
				Expect(err).NotTo(HaveOccurred())
				openingID = createResp.OpeningID

				// Now get the opening and verify tags
				getReq := employer.GetOpeningRequest{ID: openingID}
				resp = doPOST(
					token,
					getReq,
					"/employer/get-opening",
					http.StatusOK,
					true,
				)
				var opening employer.Opening
				err = json.Unmarshal(resp.([]byte), &opening)
				Expect(err).NotTo(HaveOccurred())

				Expect(opening.Tags).To(HaveLen(1))
				Expect(string(opening.Tags[0].Name)).To(Equal("Go"))
				// Verify that the ID matches the existing Go tag ID
				Expect(
					opening.Tags[0].ID,
				).To(Equal(common.VTagID("12345678-0015-0015-0015-000000070001")))
			},
		)
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
				Tags: []common.VTagID{
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
				Tags: []common.VTagID{
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
				Tags: []common.VTagID{
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
				Tags: []common.VTagID{
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
