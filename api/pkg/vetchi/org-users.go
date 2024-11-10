package vetchi

type OrgUserRole string

type OrgUserRoles []OrgUserRole

func (roles OrgUserRoles) StringArray() []string {
	var rolesStr []string
	for _, role := range roles {
		rolesStr = append(rolesStr, string(role))
	}
	return rolesStr
}

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

	// The user has been added to the organization but not yet signed up
	AddedOrgUserState OrgUserState = "ADDED_ORG_USER"

	// The user is no longer active in the organization
	DisabledOrgUserState OrgUserState = "DISABLED_ORG_USER"

	// The user is replicated from a different directory service (e.g. LDAP, Google, Microsoft Active Directory, etc.)
	ReplicatedOrgUserState OrgUserState = "REPLICATED_ORG_USER"
)

type AddOrgUserRequest struct {
	Name  string       `json:"name"  validate:"required,min=3,max=255"`
	Email string       `json:"email" validate:"required,email,min=3,max=255"`
	Roles OrgUserRoles `json:"roles" validate:"required,validate_org_user_roles"`
}

type DisableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type EnableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Should we graphqlize this for auto-completion ?
type FilterOrgUsersRequest struct {
	Prefix string         `json:"prefix" validate:"omitempty,min=1,max=255"`
	State  []OrgUserState `json:"state"  validate:"omitempty,validate_org_user_state"`

	PaginationKey string `json:"pagination_key" validate:"omitempty,email"`
	Limit         int    `json:"limit"          validate:"omitempty,min=0,max=40"`
}

func (filterOrgUsersReq *FilterOrgUsersRequest) StatesAsStrings() []string {
	if len(filterOrgUsersReq.State) == 0 {
		return []string{string(ActiveOrgUserState), string(AddedOrgUserState)}
	}

	var states []string
	for _, state := range filterOrgUsersReq.State {
		// Already validated by Vator validate_org_user_state
		states = append(states, string(state))
	}
	return states
}

type OrgUser struct {
	Name  string       `json:"name"  db:"name"`
	Email string       `json:"email" db:"email"`
	Roles OrgUserRoles `json:"roles" db:"org_user_roles"`
	State OrgUserState `json:"state" db:"org_user_state"`
}

// OrgUserShort is intended to be used in rendering of OrgUsers
// within the Employer UI, Autocompletion on Employer UI, etc.
// Not to be used on the Hub UI or even exposed on Hub APIs.
type OrgUserShort struct {
	Name  string `json:"name"`
	Email string `json:"email"`

	// If there is a HubUser (hub_users table in the db) who has the above email
	// as one of the VERIFIED emails in hub_users_official_emails table,
	// then this field will contain the handle of that HubUser.
	VetchiHandle string `json:"vetchi_handle,omitempty"`
}

type SignupOrgUserRequest struct {
	Name        string `json:"name"         validate:"required,min=3,max=255"`
	Password    string `json:"password"     validate:"required,password"`
	InviteToken string `json:"invite_token" validate:"required,min=1,max=255"`
}

type UpdateOrgUserRequest struct {
	Email string       `json:"email" validate:"required,email,min=3,max=255"`
	Name  string       `json:"name"  validate:"required,min=3,max=255"`
	Roles OrgUserRoles `json:"roles" validate:"required,validate_org_user_roles"`
}
