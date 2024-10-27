package db

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDupLocationName = errors.New("location name already exists")
	ErrNoLocation      = errors.New("location not found")
)

type AddLocationReq struct {
	Title            string
	CountryCode      string
	PostalAddress    string
	PostalCode       string
	OpenStreetMapURL string
	CityAka          []string

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type DefunctLocationReq struct {
	Title string

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type GetLocByNameReq struct {
	Title string

	EmployerID uuid.UUID
}

type GetLocationsReq struct {
	States []string

	PaginationKey string
	Limit         int

	EmployerID uuid.UUID
}

type RenameLocationReq struct {
	OldTitle string
	NewTitle string

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type UpdateLocationReq struct {
	Title            string
	CountryCode      string
	PostalAddress    string
	PostalCode       string
	OpenStreetMapURL string
	CityAka          []string

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}
