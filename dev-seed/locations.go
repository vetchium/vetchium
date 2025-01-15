package main

import (
	"log"

	"github.com/psankar/vetchi/typespec/employer"
)

func addLocation(token string, req employer.AddLocationRequest) {
	makeRequest("POST", "/employer/add-location", token, req, nil)
}

func createLocations() {
	// Get tokens from the global map
	gryffindorVal, ok := sessionTokens.Load("admin@gryffindor.example")
	if !ok {
		log.Fatal("failed to get gryffindor token")
	}
	gryffindorToken := gryffindorVal.(string)

	hufflepuffVal, ok := sessionTokens.Load("admin@hufflepuff.example")
	if !ok {
		log.Fatal("failed to get hufflepuff token")
	}
	hufflepuffToken := hufflepuffVal.(string)

	ravenclawVal, ok := sessionTokens.Load("admin@ravenclaw.example")
	if !ok {
		log.Fatal("failed to get ravenclaw token")
	}
	ravenclawToken := ravenclawVal.(string)

	slytherinVal, ok := sessionTokens.Load("admin@slytherin.example")
	if !ok {
		log.Fatal("failed to get slytherin token")
	}
	slytherinToken := slytherinVal.(string)

	locations := []struct {
		token string
		req   employer.AddLocationRequest
	}{
		// Gryffindor locations
		{
			gryffindorToken,
			employer.AddLocationRequest{
				Title:            "London HQ",
				CountryCode:      "GBR",
				PostalAddress:    "1 Tower Bridge, London",
				PostalCode:       "SE1 2UP",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/1234567",
				CityAka:          []string{"Greater London", "The City"},
			},
		},
		{
			gryffindorToken,
			employer.AddLocationRequest{
				Title:            "Edinburgh Office",
				CountryCode:      "GBR",
				PostalAddress:    "1 Castle Terrace, Edinburgh",
				PostalCode:       "EH1 2EF",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/2345678",
				CityAka:          []string{"Auld Reekie"},
			},
		},
		{
			gryffindorToken,
			employer.AddLocationRequest{
				Title:            "Dublin Hub",
				CountryCode:      "IRL",
				PostalAddress:    "Grand Canal Dock, Dublin",
				PostalCode:       "D02 XR80",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/3456789",
				CityAka:          []string{"Baile Átha Cliath"},
			},
		},
		{
			gryffindorToken,
			employer.AddLocationRequest{
				Title:            "Toronto Base",
				CountryCode:      "CAN",
				PostalAddress:    "100 King Street West, Toronto",
				PostalCode:       "M5X 1E1",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/4567890",
				CityAka:          []string{"GTA", "The 6ix"},
			},
		},
		{
			gryffindorToken,
			employer.AddLocationRequest{
				Title:            "Sydney Office",
				CountryCode:      "AUS",
				PostalAddress:    "1 Macquarie Place, Sydney",
				PostalCode:       "NSW 2000",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/5678901",
				CityAka:          []string{"Harbour City"},
			},
		},

		// Hufflepuff locations
		{
			hufflepuffToken,
			employer.AddLocationRequest{
				Title:            "Amsterdam HQ",
				CountryCode:      "NLD",
				PostalAddress:    "Dam Square 1, Amsterdam",
				PostalCode:       "1012 JL",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/6789012",
				CityAka:          []string{"Mokum"},
			},
		},
		{
			hufflepuffToken,
			employer.AddLocationRequest{
				Title:            "Copenhagen Hub",
				CountryCode:      "DNK",
				PostalAddress:    "Rådhuspladsen 1, Copenhagen",
				PostalCode:       "1550",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/7890123",
				CityAka:          []string{"København"},
			},
		},
		{
			hufflepuffToken,
			employer.AddLocationRequest{
				Title:            "Oslo Office",
				CountryCode:      "NOR",
				PostalAddress:    "Karl Johans gate 1, Oslo",
				PostalCode:       "0154",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/8901234",
			},
		},
		{
			hufflepuffToken,
			employer.AddLocationRequest{
				Title:            "Stockholm Base",
				CountryCode:      "SWE",
				PostalAddress:    "Sergels torg 1, Stockholm",
				PostalCode:       "111 57",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/9012345",
				CityAka:          []string{"Eken"},
			},
		},
		{
			hufflepuffToken,
			employer.AddLocationRequest{
				Title:            "Helsinki Hub",
				CountryCode:      "FIN",
				PostalAddress:    "Mannerheimintie 1, Helsinki",
				PostalCode:       "00100",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/0123456",
				CityAka:          []string{"Stadi"},
			},
		},

		// Ravenclaw locations
		{
			ravenclawToken,
			employer.AddLocationRequest{
				Title:            "Singapore HQ",
				CountryCode:      "SGP",
				PostalAddress:    "1 Raffles Place",
				PostalCode:       "048616",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/1234567",
				CityAka:          []string{"Lion City"},
			},
		},
		{
			ravenclawToken,
			employer.AddLocationRequest{
				Title:            "Tokyo Office",
				CountryCode:      "JPN",
				PostalAddress:    "1-1 Marunouchi, Chiyoda",
				PostalCode:       "100-0005",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/2345678",
				CityAka:          []string{"東京"},
			},
		},
		{
			ravenclawToken,
			employer.AddLocationRequest{
				Title:            "Seoul Hub",
				CountryCode:      "KOR",
				PostalAddress:    "Jung-gu, Seoul",
				PostalCode:       "04533",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/3456789",
				CityAka:          []string{"서울"},
			},
		},
		{
			ravenclawToken,
			employer.AddLocationRequest{
				Title:            "Mumbai Base",
				CountryCode:      "IND",
				PostalAddress:    "Nariman Point, Mumbai",
				PostalCode:       "400021",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/4567890",
				CityAka:          []string{"Bombay"},
			},
		},
		{
			ravenclawToken,
			employer.AddLocationRequest{
				Title:            "Dubai Office",
				CountryCode:      "ARE",
				PostalAddress:    "Downtown Dubai",
				PostalCode:       "00000",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/5678901",
				CityAka:          []string{"دبي"},
			},
		},

		// Slytherin locations
		{
			slytherinToken,
			employer.AddLocationRequest{
				Title:            "Berlin HQ",
				CountryCode:      "DEU",
				PostalAddress:    "Unter den Linden 1, Berlin",
				PostalCode:       "10117",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/6789012",
			},
		},
		{
			slytherinToken,
			employer.AddLocationRequest{
				Title:            "Paris Office",
				CountryCode:      "FRA",
				PostalAddress:    "1 Avenue des Champs-Élysées, Paris",
				PostalCode:       "75008",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/7890123",
				CityAka:          []string{"Paname"},
			},
		},
		{
			slytherinToken,
			employer.AddLocationRequest{
				Title:            "Madrid Hub",
				CountryCode:      "ESP",
				PostalAddress:    "Puerta del Sol, Madrid",
				PostalCode:       "28013",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/8901234",
			},
		},
		{
			slytherinToken,
			employer.AddLocationRequest{
				Title:            "Rome Base",
				CountryCode:      "ITA",
				PostalAddress:    "Via del Corso 1, Rome",
				PostalCode:       "00186",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/9012345",
				CityAka:          []string{"Roma"},
			},
		},
		{
			slytherinToken,
			employer.AddLocationRequest{
				Title:            "Vienna Hub",
				CountryCode:      "AUT",
				PostalAddress:    "Stephansplatz 1, Vienna",
				PostalCode:       "1010",
				OpenStreetMapURL: "https://www.openstreetmap.org/way/0123456",
				CityAka:          []string{"Wien"},
			},
		},
	}

	for _, loc := range locations {
		addLocation(loc.token, loc.req)
	}
}
