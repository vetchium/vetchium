package dolores

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type locationTestCase struct {
	description string
	token       string
	location    vetchi.AddLocationRequest
	wantStatus  int
}

var _ = FDescribe("Employer Locations", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, viewerToken string
	var crud1Token, crud2Token string
	var nonLocationToken, multipleNonLocationRolesToken string

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
			dummyLocation := vetchi.AddLocationRequest{
				Title:         "Location-dummy",
				CountryCode:   "SCO",
				PostalCode:    "TN-1234",
				PostalAddress: "Hogwarts School of Witchcraft and Wizardry, Highlands, Scotland",
				CityAka:       []string{"Hogwarts", "School"},
			}

			location1 := vetchi.AddLocationRequest{
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
					location: vetchi.AddLocationRequest{
						Title:            "Location-2",
						CountryCode:      "IND",
						PostalCode:       "12345",
						PostalAddress:    "4 Privet Drive, Little Whinging, Surrey",
						OpenStreetMapURL: "https://www.openstreetmap.org/way/966341718",
						CityAka:          []string{"Dursleys"},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Crud1 token",
					token:       crud1Token,
					location: vetchi.AddLocationRequest{
						Title:         "Location-3",
						CountryCode:   "USA",
						PostalCode:    "12345",
						PostalAddress: "6, Murray Hills, New Jersey",
						CityAka:       []string{},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with Crud2 token",
					token:       crud2Token,
					location: vetchi.AddLocationRequest{
						Title:         "Location-4",
						CountryCode:   "GBR",
						PostalCode:    "23456",
						PostalAddress: "Number 12, Grimmauld Place, London PIN: 23456",
						CityAka:       []string{"Order of the Phoenix"},
					},
					wantStatus: http.StatusOK,
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
				{
					description: "with missing title",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with small invalid title",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "a",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long invalid title",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         strings.Repeat("a", 33),
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with missing country code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with small invalid country code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IN",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long invalid country code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "INDIA",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with missing postal code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with small invalid postal code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "1",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long invalid postal code",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    strings.Repeat("1", 17),
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with missing postal address",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with small invalid postal address",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "xy",
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long invalid postal address",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: strings.Repeat("x", 1025),
						CityAka:       []string{"Dursleys"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid city aka",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka:       []string{"Dursleys", "xy"},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with long city aka",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
						Title:         "Location-6",
						CountryCode:   "IND",
						PostalCode:    "12345",
						PostalAddress: "4 Privet Drive, Little Whinging, Surrey",
						CityAka: []string{
							"Dursleys",
							strings.Repeat("x", 100),
						},
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid number of city aka",
					token:       adminToken,
					location: vetchi.AddLocationRequest{
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
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, testCase := range testCases {
				testAddLocation(
					testCase.token,
					testCase.location,
					testCase.wantStatus,
				)
				fmt.Fprintf(GinkgoWriter, "%s\n", testCase.description)
			}
		})
	})
})

func testAddLocation(
	token string,
	location vetchi.AddLocationRequest,
	wantStatus int,
) {
	fmt.Fprintf(
		GinkgoWriter,
		"testAddLocation: token=%s, location=%v, wantStatus=%d\n",
		token, location, wantStatus,
	)
	reqBody := vetchi.AddLocationRequest{
		Title:            location.Title,
		CountryCode:      location.CountryCode,
		PostalCode:       location.PostalCode,
		PostalAddress:    location.PostalAddress,
		CityAka:          location.CityAka,
		OpenStreetMapURL: location.OpenStreetMapURL,
	}
	testPOST(token, reqBody, "/employer/add-location", wantStatus)
}
