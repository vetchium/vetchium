package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = Describe("Employer Locations", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, viewerToken string
	var crud1Token, crud2Token string
	var nonLocationToken, multipleNonLocationRolesToken string

	dummyLocation := employer.AddLocationRequest{
		Title:         "Location-dummy",
		CountryCode:   "SCO",
		PostalCode:    "TN-1234",
		PostalAddress: "Hogwarts School of Witchcraft and Wizardry, Highlands, Scotland",
		CityAka:       []string{"Hogwarts", "School"},
	}

	location1 := employer.AddLocationRequest{
		Title:       "Location-1",
		CountryCode: "UAE",
		PostalCode:  "12345",
		PostalAddress: `
Number 6, Viveganandhar Theru,
Dubai Kurukkuchandhu,
Dubai Main Road,
Dubai
PIN: 12345`,
		CityAka: []string{"Saarja", "Beghireen", "Abidhaabi"},
	}

	location2 := employer.AddLocationRequest{
		Title:            "Location-2",
		CountryCode:      "IND",
		PostalCode:       "12345",
		PostalAddress:    "4 Privet Drive, Little Whinging, Surrey",
		OpenStreetMapURL: "https://www.openstreetmap.org/way/966341718",
		CityAka:          []string{"Dursleys"},
	}

	location3 := employer.AddLocationRequest{
		Title:         "Location-3",
		CountryCode:   "USA",
		PostalCode:    "12345",
		PostalAddress: "6, Murray Hills, New Jersey",
		CityAka:       []string{},
	}

	location4 := employer.AddLocationRequest{
		Title:         "Location-4",
		CountryCode:   "GBR",
		PostalCode:    "23456",
		PostalAddress: "Number 12, Grimmauld Place, London PIN: 23456",
		CityAka:       []string{"Order of the Phoenix"},
	}

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0003-employer-locations-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@location.example":                       &adminToken,
			"crud1@location.example":                       &crud1Token,
			"crud2@location.example":                       &crud2Token,
			"viewer@location.example":                      &viewerToken,
			"non-location@location.example":                &nonLocationToken,
			"multiple-non-location-roles@location.example": &multipleNonLocationRolesToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"location.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}

		// Wait until all the signin operations are complete
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0003-employer-locations-down.pgsql")
		db.Close()
	})

	Describe("Locations related Tests", func() {
		It("Add Location", func() {
			type locationTestCase struct {
				description string
				token       string
				location    employer.AddLocationRequest
				wantStatus  int
			}

			testCases := []locationTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					location:    location1,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Admin token second location",
					token:       adminToken,
					location:    location2,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Crud1 token",
					token:       crud1Token,
					location:    location3,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Crud2 token",
					token:       crud2Token,
					location:    location4,
					wantStatus:  http.StatusOK,
				},
				{
					description: "with Viewer token",
					token:       viewerToken,
					location:    dummyLocation,
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with Non-location role token",
					token:       nonLocationToken,
					location:    dummyLocation,
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with Multiple non-location roles token",
					token:       multipleNonLocationRolesToken,
					location:    dummyLocation,
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with Invalid token",
					token:       "invalid-token",
					location:    dummyLocation,
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with Empty token",
					token:       "",
					location:    dummyLocation,
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with Admin token with duplicate title",
					token:       adminToken,
					location:    location1,
					wantStatus:  http.StatusConflict,
				},
				{
					description: "with Crud1 token with duplicate title",
					token:       crud1Token,
					location:    location1,
					wantStatus:  http.StatusConflict,
				},
			}

			for _, testCase := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
				testAddLocation(
					testCase.token,
					testCase.location,
					testCase.wantStatus,
				)
			}
		})

		It("Add Location Validation", func() {
			type locationValidationTestCase struct {
				description   string
				token         string
				location      employer.AddLocationRequest
				wantStatus    int
				wantErrFields []string
			}
			testCases := []locationValidationTestCase{
				{
					description: "with missing title",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with small invalid title",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "a",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with long invalid title",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         strings.Repeat("a", 33),
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with missing country code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with small invalid country code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IN",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with long invalid country code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "INDIA",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with missing postal code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with small invalid postal code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "1",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with long invalid postal code",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    strings.Repeat("1", 17),
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with missing postal address",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with small invalid postal address",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "xy",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with long invalid postal address",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: strings.Repeat("x", 1025),
						CityAka:       []string{"Dursleys"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with invalid city aka",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys", "xy"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
				{
					description: "with long city aka",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka: []string{
							"Dursleys",
							strings.Repeat("x", 100),
						},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
				{
					description: "with invalid number of city aka",
					token:       adminToken,
					location: employer.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka: []string{
							"Harry Potter and the Sorcerer's Stone",
							"Harry Potter and the Chamber of Secrets",
							"Harry Potter and the Prisoner of Azkaban",
							"Harry Potter and the Goblet of Fire",
							"Harry Potter and the Order of the Phoenix",
							"Harry Potter and the Half-Blood Prince",
							"Harry Potter and the Deathly Hallows",
						},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
			}

			for _, testCase := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
				validationErrors := testAddLocationGetResp(
					testCase.token,
					testCase.location,
					testCase.wantStatus,
				)
				Expect(
					validationErrors.Errors,
				).Should(ContainElements(testCase.wantErrFields))
			}
		})

		It("Get Locations", func() {
			type testGetLocationsTestCase struct {
				description         string
				token               string
				getLocationsRequest employer.GetLocationsRequest
				wantStatus          int
			}

			testCases := []testGetLocationsTestCase{
				{
					description:         "with Admin token and no filters",
					token:               adminToken,
					getLocationsRequest: employer.GetLocationsRequest{},
					wantStatus:          http.StatusOK,
				},
				{
					description:         "with Viewer token and no filters",
					token:               viewerToken,
					getLocationsRequest: employer.GetLocationsRequest{},
					wantStatus:          http.StatusOK,
				},
				{
					description:         "with Crud1 token and no filters",
					token:               crud1Token,
					getLocationsRequest: employer.GetLocationsRequest{},
					wantStatus:          http.StatusOK,
				},
				{
					description: "with invalid state",
					token:       adminToken,
					getLocationsRequest: employer.GetLocationsRequest{
						States: []employer.LocationState{"invalid-state"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid low limit",
					token:       adminToken,
					getLocationsRequest: employer.GetLocationsRequest{
						Limit: -1,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid high limit",
					token:       adminToken,
					getLocationsRequest: employer.GetLocationsRequest{
						Limit: 101,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description:         "with invalid token",
					token:               "invalid-token",
					getLocationsRequest: employer.GetLocationsRequest{},
					wantStatus:          http.StatusUnauthorized,
				},
				{
					description:         "with empty token",
					token:               "",
					getLocationsRequest: employer.GetLocationsRequest{},
					wantStatus:          http.StatusUnauthorized,
				},
			}

			for _, testCase := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
				locations := testGetLocations(
					testCase.token,
					testCase.getLocationsRequest,
					testCase.wantStatus,
				)
				if testCase.wantStatus == http.StatusOK {
					Expect(locations).Should(HaveLen(4))
					Expect(locations).Should(ContainElements(
						makeLocation(location1, employer.ActiveLocation),
						makeLocation(location2, employer.ActiveLocation),
						makeLocation(location3, employer.ActiveLocation),
						makeLocation(location4, employer.ActiveLocation),
					))
				}
			}
		})

		It("Defunct Location", func() {
			type testDefunctLocationTestCase struct {
				description string
				token       string
				location    employer.DefunctLocationRequest
				wantStatus  int
			}

			testCases := []testDefunctLocationTestCase{
				{
					description: "with Admin token",
					token:       adminToken,
					location: employer.DefunctLocationRequest{
						Title: location1.Title,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Crud1 token",
					token:       crud1Token,
					location: employer.DefunctLocationRequest{
						Title: location2.Title,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					location: employer.DefunctLocationRequest{
						Title: location3.Title,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with empty token",
					token:       "",
					location: employer.DefunctLocationRequest{
						Title: location4.Title,
					},
					wantStatus: http.StatusUnauthorized,
				},
			}

			for _, testCase := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
				testDefunctLocation(
					testCase.token,
					testCase.location,
					testCase.wantStatus,
				)
			}

			// Verify that the defunct location is not returned in the list
			locations := testGetLocations(
				adminToken,
				employer.GetLocationsRequest{},
				http.StatusOK,
			)
			Expect(locations).Should(HaveLen(2))
			Expect(locations).Should(ContainElements(
				makeLocation(location3, employer.ActiveLocation),
				makeLocation(location4, employer.ActiveLocation),
			))

			// Verify that the defunct location is returned in the list
			locations = testGetLocations(
				adminToken,
				employer.GetLocationsRequest{
					States: []employer.LocationState{employer.DefunctLocation},
				},
				http.StatusOK,
			)
			Expect(locations).Should(HaveLen(2))
			Expect(locations).Should(ContainElements(
				makeLocation(location1, employer.DefunctLocation),
				makeLocation(location2, employer.DefunctLocation),
			))
		})

		It("Update and Get Location", func() {
			location3Updated := makeLocation(location1, employer.ActiveLocation)
			location3Updated.Title = location3.Title

			updateLocationRequest := employer.UpdateLocationRequest{
				Title:         location3.Title,
				CountryCode:   location3Updated.CountryCode,
				PostalCode:    location3Updated.PostalCode,
				PostalAddress: location3Updated.PostalAddress,
				CityAka:       location3Updated.CityAka,
			}

			_ = doPOST(
				adminToken,
				updateLocationRequest,
				"/employer/update-location",
				http.StatusOK,
				false,
			)

			getLocationRequest := employer.GetLocationRequest{
				Title: location3.Title,
			}
			locationRaw := doPOST(
				adminToken,
				getLocationRequest,
				"/employer/get-location",
				http.StatusOK,
				true,
			).([]byte)
			var location employer.Location
			err := json.Unmarshal(locationRaw, &location)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(location).Should(Equal(location3Updated))
		})

		It("Update Location Validation", func() {
			type testUpdateLocationValidationTestCase struct {
				description   string
				token         string
				location      employer.UpdateLocationRequest
				wantStatus    int
				wantErrFields []string
			}

			testCases := []testUpdateLocationValidationTestCase{
				{
					description: "with missing title",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title: "",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with small invalid title",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title: "a",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with long invalid title",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title: strings.Repeat("a", 33),
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"title"},
				},
				{
					description: "with missing country code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   "",
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with small invalid country code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   "IN",
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with long invalid country code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   "INDIA",
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"country_code"},
				},
				{
					description: "with missing postal code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    "",
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with small invalid postal code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    "1",
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with long invalid postal code",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    strings.Repeat("1", 17),
						PostalAddress: location3.PostalAddress,
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_code"},
				},
				{
					description: "with missing postal address",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: "",
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with small invalid postal address",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: "xy",
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with long invalid postal address",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: strings.Repeat("x", 1025),
						CityAka:       location3.CityAka,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"postal_address"},
				},
				{
					description: "with invalid city aka",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka:       []string{"Dursleys", "xy"},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
				{
					description: "with long city aka",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka: []string{
							"Dursleys",
							strings.Repeat("x", 100),
						},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
				{
					description: "with invalid number of city aka",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:         location3.Title,
						CountryCode:   location3.CountryCode,
						PostalCode:    location3.PostalCode,
						PostalAddress: location3.PostalAddress,
						CityAka: []string{
							"Harry Potter and the Sorcerer's Stone",
							"Harry Potter and the Chamber of Secrets",
							"Harry Potter and the Prisoner of Azkaban",
							"Harry Potter and the Goblet of Fire",
							"Harry Potter and the Order of the Phoenix",
							"Harry Potter and the Half-Blood Prince",
							"Harry Potter and the Deathly Hallows",
						},
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"city_aka"},
				},
				{
					description: "with invalid small openstreetmap url",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:            location3.Title,
						CountryCode:      location3.CountryCode,
						PostalCode:       location3.PostalCode,
						PostalAddress:    location3.PostalAddress,
						CityAka:          location3.CityAka,
						OpenStreetMapURL: "a",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"openstreetmap_url"},
				},
				{
					description: "with invalid long openstreetmap url",
					token:       adminToken,
					location: employer.UpdateLocationRequest{
						Title:            location3.Title,
						CountryCode:      location3.CountryCode,
						PostalCode:       location3.PostalCode,
						PostalAddress:    location3.PostalAddress,
						CityAka:          location3.CityAka,
						OpenStreetMapURL: strings.Repeat("x", 256),
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"openstreetmap_url"},
				},
			}

			for _, testCase := range testCases {
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
				validationErrors := testUpdateLocationGetResp(
					testCase.token,
					testCase.location,
					testCase.wantStatus,
				)
				Expect(
					validationErrors.Errors,
				).Should(ContainElements(testCase.wantErrFields))
			}
		})

		It("Create Get Defunct Locations in bulk; Test Pagination", func() {
			// runID parameter is needed because even if a CC is defunct,
			// its name cannot be used by another CC

			fmt.Fprintf(GinkgoWriter, "count is not divisible by limit\n")
			bulkAddDefunctLocation(adminToken, "run-1", 30, 4)

			fmt.Fprintf(GinkgoWriter, "count is divisible by limit\n")
			bulkAddDefunctLocation(adminToken, "run-2", 32, 4)

			fmt.Fprintf(GinkgoWriter, "limit is greater than count\n")
			bulkAddDefunctLocation(adminToken, "run-3", 30, 40)
		})
	})
})

func testAddLocation(
	token string,
	location employer.AddLocationRequest,
	wantStatus int,
) {
	fmt.Fprintf(
		GinkgoWriter,
		"testAddLocation: token=%s, location=%v, wantStatus=%d\n",
		token, location, wantStatus,
	)
	reqBody := employer.AddLocationRequest{
		Title:            location.Title,
		CountryCode:      location.CountryCode,
		PostalCode:       location.PostalCode,
		PostalAddress:    location.PostalAddress,
		CityAka:          location.CityAka,
		OpenStreetMapURL: location.OpenStreetMapURL,
	}
	testPOST(token, reqBody, "/employer/add-location", wantStatus)
}

func testAddLocationGetResp(
	token string,
	location employer.AddLocationRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(
		token,
		location,
		"/employer/add-location",
		wantStatus,
	).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func testGetLocations(
	token string,
	getLocationsRequest employer.GetLocationsRequest,
	wantStatus int,
) []employer.Location {
	resp := testPOSTGetResp(
		token,
		getLocationsRequest,
		"/employer/get-locations",
		wantStatus,
	).([]byte)

	var locations []employer.Location
	if wantStatus == http.StatusOK {
		err := json.Unmarshal(resp, &locations)
		Expect(err).ShouldNot(HaveOccurred())
	}
	return locations
}

func makeLocation(
	req employer.AddLocationRequest,
	state employer.LocationState,
) employer.Location {
	return employer.Location{
		Title:            req.Title,
		CountryCode:      req.CountryCode,
		PostalCode:       req.PostalCode,
		PostalAddress:    req.PostalAddress,
		CityAka:          req.CityAka,
		OpenStreetMapURL: req.OpenStreetMapURL,
		State:            state,
	}
}

func testDefunctLocation(
	token string,
	defunctLocationRequest employer.DefunctLocationRequest,
	wantStatus int,
) {
	fmt.Fprintf(
		GinkgoWriter,
		"testDefunctLocation: token=%s, defunctLocationRequest=%v, wantStatus=%d\n",
		token,
		defunctLocationRequest,
		wantStatus,
	)
	reqBody := employer.DefunctLocationRequest{
		Title: defunctLocationRequest.Title,
	}
	testPOST(token, reqBody, "/employer/defunct-location", wantStatus)
}

func testUpdateLocationGetResp(
	token string,
	location employer.UpdateLocationRequest,
	wantStatus int,
) common.ValidationErrors {
	resp := testPOSTGetResp(
		token,
		location,
		"/employer/update-location",
		wantStatus,
	).([]byte)
	var validationErrors common.ValidationErrors
	err := json.Unmarshal(resp, &validationErrors)
	Expect(err).ShouldNot(HaveOccurred())
	return validationErrors
}

func bulkAddDefunctLocation(
	token string,
	runID string,
	count int,
	limit int,
) {
	wantLocations := []employer.Location{}

	dummyLocation := employer.AddLocationRequest{
		Title:         "dummy-location",
		CountryCode:   "IND",
		PostalCode:    "123456",
		PostalAddress: "dummy-address",
		CityAka:       []string{"dummy-city"},
	}

	for i := 0; i < count; i++ {
		newLocation := makeLocation(dummyLocation, employer.ActiveLocation)
		newLocation.Title = fmt.Sprintf("LOC-%s-%d", runID, i)
		wantLocations = append(wantLocations, newLocation)

		newLocationReq := employer.AddLocationRequest{
			Title:         newLocation.Title,
			CountryCode:   newLocation.CountryCode,
			PostalCode:    newLocation.PostalCode,
			PostalAddress: newLocation.PostalAddress,
			CityAka:       newLocation.CityAka,
		}
		testAddLocation(token, newLocationReq, http.StatusOK)
	}

	paginationKey := ""
	gotLocations := []employer.Location{}
	for {
		getLocationsRequest := employer.GetLocationsRequest{
			PaginationKey: paginationKey,
			Limit:         limit,
		}
		locations := testGetLocations(token, getLocationsRequest, http.StatusOK)

		if len(locations) == 0 {
			break
		}

		gotLocations = append(gotLocations, locations...)
		if len(locations) < limit {
			break
		}
		paginationKey = locations[len(locations)-1].Title
	}

	Expect(gotLocations).Should(ContainElements(wantLocations))

	for _, location := range gotLocations {
		testDefunctLocation(token, employer.DefunctLocationRequest{
			Title: location.Title,
		}, http.StatusOK)
	}
}
