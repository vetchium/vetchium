package main

import (
	"log"

	"github.com/psankar/vetchi/typespec/employer"
)

// addCostCenter creates a new cost center for an employer
func addCostCenter(token, name, notes string) {
	req := employer.AddCostCenterRequest{
		Name:  employer.CostCenterName(name),
		Notes: notes,
	}

	makeRequest("POST", "/employer/add-cost-center", token, req, nil)
}

func createCostCenters() {
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

	costCenters := []struct {
		token string
		name  string
		notes string
	}{
		// Gryffindor cost centers
		{gryffindorToken, "UK Operations", "All UK based operations and staff"},
		{gryffindorToken, "Ireland Division", "Irish operations and expansion"},
		{gryffindorToken, "Canada Business", "North American presence"},
		{gryffindorToken, "APAC Operations", "Asia Pacific operations"},
		{
			gryffindorToken,
			"Global Marketing",
			"Marketing activities across all regions",
		},

		// Hufflepuff cost centers
		{
			hufflepuffToken,
			"Benelux Operations",
			"Netherlands, Belgium, Luxembourg operations",
		},
		{
			hufflepuffToken,
			"Nordic Division",
			"Scandinavian countries operations",
		},
		{hufflepuffToken, "EU Marketing", "European marketing initiatives"},
		{hufflepuffToken, "R&D Labs", "Research and development centers"},
		{hufflepuffToken, "EU Admin", "European administrative operations"},

		// Ravenclaw cost centers
		{
			ravenclawToken,
			"APAC Headquarters",
			"Singapore and surrounding operations",
		},
		{ravenclawToken, "Japan Operations", "Japanese market operations"},
		{ravenclawToken, "Korea Division", "Korean market presence"},
		{ravenclawToken, "India Operations", "Indian subcontinent operations"},
		{
			ravenclawToken,
			"Middle East Division",
			"Middle East expansion and operations",
		},

		// Slytherin cost centers
		{
			slytherinToken,
			"DACH Operations",
			"Germany, Austria, Switzerland operations",
		},
		{slytherinToken, "France Division", "French market operations"},
		{
			slytherinToken,
			"Southern Europe",
			"Spain, Italy, and Mediterranean operations",
		},
		{slytherinToken, "EU Projects", "Special European projects"},
		{
			slytherinToken,
			"Continental R&D",
			"European research and development",
		},
	}

	for _, cc := range costCenters {
		addCostCenter(cc.token, cc.name, cc.notes)
	}
}
