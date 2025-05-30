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
			Expect(len(result)).To(BeNumerically(">", 0))
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

			Expect(got).To(ContainElement("PaaS (Platform as a Service)"))
			Expect(got).To(ContainElement("Product Management"))
			Expect(got).To(ContainElement("Python Programming"))
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
				TagIDs: []common.VTagID{"python", "golang"},
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

			Expect(got[0]).To(Equal("Go Programming Language"))
			Expect(got[1]).To(Equal("Python Programming"))
		})

		It("should create opening with mixed tags and verify them", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Mixed Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers with Python and PostgreSQL experience",
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
				TagIDs: []common.VTagID{
					"python",
					"sports",
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

			sortedTags := []string{
				string(opening.Tags[0].Name),
				string(opening.Tags[1].Name),
			}
			sort.Strings(sortedTags)
			Expect(sortedTags[0]).To(Equal("Python Programming"))
			Expect(sortedTags[1]).To(Equal("Sports"))
		})

		It(
			"should create opening with existing tag and verify it",
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
					TagIDs: []common.VTagID{"golang"},
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
				Expect(
					string(opening.Tags[0].Name),
				).To(Equal("Go Programming Language"))
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
				TagIDs: []common.VTagID{"python", "golang"},
			}

			testPOST(token, req, "/employer/create-opening", http.StatusOK)
		})

		It("should create opening with mixed tags", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Mixed Tags",
				Positions:         2,
				JD:                "Looking for talented software engineers with Python and PostgreSQL experience",
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
				TagIDs: []common.VTagID{
					"python",     // Python
					"technology", // Technology (replacing PostgreSQL which isn't in vetchium-tags.json)
				},
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
				TagIDs: []common.VTagID{
					"golang",     // Go
					"java",       // Java
					"python",     // Python
					"technology", // Technology (4th tag to exceed limit)
				},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})

		It("should fail when more than 3 existing tags are provided", func() {
			req := employer.CreateOpeningRequest{
				Title:             "Test Opening with Too Many Existing Tags",
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
				TagIDs: []common.VTagID{
					"golang",     // Go
					"java",       // Java
					"python",     // Python
					"technology", // Technology (4th tag to exceed limit)
				},
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
				TagIDs: []common.VTagID{
					"golang",     // Go
					"java",       // Java
					"python",     // Python
					"technology", // Technology (4th tag to exceed limit)
				},
			}

			testPOST(
				token,
				req,
				"/employer/create-opening",
				http.StatusBadRequest,
			)
		})

		It(
			"should return 400 with ValidationErrors for invalid tag IDs",
			func() {
				req := employer.CreateOpeningRequest{
					Title:             "Opening - Invalid TagID",
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
					TagIDs: []common.VTagID{
						"invalid-tag-id", // Invalid tag ID
					},
				}

				resp := doPOST(
					token,
					req,
					"/employer/create-opening",
					http.StatusBadRequest,
					true,
				)

				var validationErrors common.ValidationErrors
				err := json.Unmarshal(resp.([]byte), &validationErrors)
				Expect(err).NotTo(HaveOccurred())
				Expect(validationErrors.Errors).To(ContainElement("tags"))
			},
		)

		It(
			"should return 400 with ValidationErrors for mix of valid and invalid tag IDs",
			func() {
				req := employer.CreateOpeningRequest{
					Title:             "Test Opening with Mixed Tag IDs",
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
					TagIDs: []common.VTagID{
						"python",             // Valid tag ID
						"nonexistent-tag-id", // Invalid tag ID
					},
				}

				resp := doPOST(
					token,
					req,
					"/employer/create-opening",
					http.StatusBadRequest,
					true,
				)

				var validationErrors common.ValidationErrors
				err := json.Unmarshal(resp.([]byte), &validationErrors)
				Expect(err).NotTo(HaveOccurred())
				Expect(validationErrors.Errors).To(ContainElement("tags"))
			},
		)
	})
})
