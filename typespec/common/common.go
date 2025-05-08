package common

//go:generate go run generate_countries.go
//go:generate go run generate_timezones.go

import (
	"regexp"
)

type ValidationErrors struct {
	Errors []string `json:"errors"`
}

type EmailAddress string
type Password string
type City string
type Handle string
type Domain string

func (h Handle) IsValid() bool {
	// The regex pattern allows hyphens and underscores but requires they not be at start/end
	// and not be consecutive
	matched, err := regexp.MatchString(
		"^[A-Za-z][A-Za-z0-9_-]*$", // Updated regex
		string(h),
	)
	if err != nil || !matched {
		return false
	}
	return len(h) > 1 && len(h) <= 64
}

type CountryCode string

var validCountryCodes map[string]struct{}

func (cc CountryCode) IsValid() bool {
	if len(cc) != 3 {
		return false
	}

	if cc == GlobalCountryCode {
		return true
	}

	_, ok := validCountryCodes[string(cc)]
	return ok
}

const GlobalCountryCode CountryCode = "ZZG"

type Currency string

type TimeZone string

var validTimezones map[string]struct{}

func (t TimeZone) IsValid() bool {
	_, ok := validTimezones[string(t)]
	return ok
}

type OrgUserRole string

type OrgUserRoles []OrgUserRole

func (roles OrgUserRoles) StringArray() []string {
	var rolesStr []string
	for _, role := range roles {
		rolesStr = append(rolesStr, string(role))
	}
	return rolesStr
}

func (r OrgUserRole) IsValid() bool {
	switch r {
	case Admin, AnyOrgUser,
		ApplicationsCRUD, ApplicationsViewer,
		CostCentersCRUD, CostCentersViewer,
		LocationsCRUD, LocationsViewer,
		OpeningsCRUD, OpeningsViewer,
		OrgUsersCRUD, OrgUsersViewer:
		return true
	default:
		return false
	}
}

func (roles OrgUserRoles) IsValid() bool {
	if len(roles) == 0 {
		return false
	}
	for _, role := range roles {
		if !role.IsValid() {
			return false
		}
	}
	return true
}

const (
	Admin OrgUserRole = "ADMIN"

	// This ANY role is not saved in database. If this role is the value for
	// the allowedRoles in the middleware, then any OrgUser in that Org can
	// access that route.
	AnyOrgUser OrgUserRole = "ANY_ORG_USER"

	ApplicationsCRUD   OrgUserRole = "APPLICATIONS_CRUD"
	ApplicationsViewer OrgUserRole = "APPLICATIONS_VIEWER"

	CostCentersCRUD   OrgUserRole = "COST_CENTERS_CRUD"
	CostCentersViewer OrgUserRole = "COST_CENTERS_VIEWER"

	LocationsCRUD   OrgUserRole = "LOCATIONS_CRUD"
	LocationsViewer OrgUserRole = "LOCATIONS_VIEWER"

	OpeningsCRUD   OrgUserRole = "OPENINGS_CRUD"
	OpeningsViewer OrgUserRole = "OPENINGS_VIEWER"

	OrgUsersCRUD   OrgUserRole = "ORG_USERS_CRUD"
	OrgUsersViewer OrgUserRole = "ORG_USERS_VIEWER"
)
