package employer

import (
	"github.com/psankar/vetchi/typespec/common"
)

type OrgUserRole string

type OrgUserRoles []OrgUserRole

const (
	Admin OrgUserRole = "ADMIN"

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

type OrgUserState string

const (
	ActiveOrgUserState     OrgUserState = "ACTIVE_ORG_USER"
	AddedOrgUserState      OrgUserState = "ADDED_ORG_USER"
	DisabledOrgUserState   OrgUserState = "DISABLED_ORG_USER"
	ReplicatedOrgUserState OrgUserState = "REPLICATED_ORG_USER"
)

type OrgUserShort struct {
	Name         string              `json:"name"`
	Email        common.EmailAddress `json:"email"`
	VetchiHandle string              `json:"vetchi_handle,omitempty"`
}

type OrgUser struct {
	Name  string              `json:"name"`
	Email common.EmailAddress `json:"email"`
	Roles OrgUserRoles        `json:"roles"`
	State OrgUserState        `json:"state"`
}

type AddOrgUserRequest struct {
	Name  string       `json:"name"  validate:"required,min=3,max=255"`
	Email string       `json:"email" validate:"required,email,min=3,max=255"`
	Roles OrgUserRoles `json:"roles" validate:"required,validate_org_user_roles"`
}

type DisableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email,min=3,max=255"`
}

type EnableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email,min=3,max=255"`
}

type FilterOrgUsersRequest struct {
	Prefix        *string        `json:"prefix"         validate:"omitempty,min=3,max=255"`
	PaginationKey *string        `json:"pagination_key" validate:"omitempty"`
	Limit         *int           `json:"limit"          validate:"omitempty,max=100"`
	State         []OrgUserState `json:"state"          validate:"omitempty"`
}

type SignupOrgUserRequest struct {
	Name        string          `json:"name"         validate:"required,min=3,max=255"`
	Password    common.Password `json:"password"     validate:"required,password"`
	InviteToken string          `json:"invite_token" validate:"required,min=1,max=255"`
}

type UpdateOrgUserRequest struct {
	Email string       `json:"email" validate:"required,email,min=3,max=255"`
	Name  string       `json:"name"  validate:"required,min=3,max=255"`
	Roles OrgUserRoles `json:"roles" validate:"required,validate_org_user_roles"`
}

// ... rest of the file with similar updates to use common package types
