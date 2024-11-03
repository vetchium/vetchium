package vetchi

type OrgUserRole string

const (
	Admin OrgUserRole = "ADMIN"

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
	// The user is active in the organization
	ActiveOrgUserState OrgUserState = "ACTIVE_ORG_USER"

	// The user has been invited to the organization but has not yet signed up
	InvitedOrgUserState OrgUserState = "INVITED_ORG_USER"

	// The user has been added to the organization but not yet sent an invitation email
	AddedOrgUserState OrgUserState = "ADDED_ORG_USER"

	// The user is no longer active in the organization
	DisabledOrgUserState OrgUserState = "DISABLED_ORG_USER"

	// The user is replicated from a different directory service (e.g. LDAP, Google, Microsoft Active Directory, etc.)
	ReplicatedOrgUserState OrgUserState = "REPLICATED_ORG_USER"
)

type AddOrgUserRequest struct {
	Name  string        `json:"name"  validate:"required,min=3,max=255"`
	Email string        `json:"email" validate:"required,email,min=3,max=255"`
	Roles []OrgUserRole `json:"roles" validate:"required,validate_org_user_roles"`
}

type DisableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Should we graphqlize this for auto-completion ?
type FilterOrgUsersRequest struct {
	Prefix string         `json:"prefix" validate:"required,min=1,max=255"`
	State  []OrgUserState `json:"state"  validate:"required,validate_org_user_state"`

	PaginationKey string `json:"pagination_key" validate:"email"`
	Limit         int    `json:"limit"          validate:"min=0,max=40"`
}

type OrgUser struct {
	Name  string        `json:"name"  db:"name"`
	Email string        `json:"email" db:"email"`
	Roles []OrgUserRole `json:"roles" db:"org_user_roles"`
	State OrgUserState  `json:"state" db:"org_user_state"`
}

type UpdateOrgUserRequest struct {
	Email string        `json:"email" validate:"required,email,min=3,max=255"`
	Name  string        `json:"name"  validate:"required,min=3,max=255"`
	Roles []OrgUserRole `json:"roles" validate:"required,validate_org_user_roles"`
}
