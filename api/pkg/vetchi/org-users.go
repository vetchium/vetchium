package vetchi

type OrgUserRole string

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
	Name  string        `json:"name"  validate:"required,min=1,max=255"`
	Email string        `json:"email" validate:"required,email,min=3,max=255"`
	Roles []OrgUserRole `json:"roles" validate:"required"`
}

type DisableOrgUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type FilterOrgUsersRequest struct {
	Prefix string         `json:"prefix" validate:"required,min=1,max=255"`
	State  []OrgUserState `json:"state"  validate:"required,validate_org_user_state"`

	PaginationKey string `json:"pagination_key" validate:"required,email"`
	Limit         int    `json:"limit"          validate:"required,min=1,max=40"`
}

type OrgUser struct {
	Name  string        `json:"name"  db:"name"`
	Email string        `json:"email" db:"email"`
	Roles []OrgUserRole `json:"roles" db:"org_user_roles"`
	State OrgUserState  `json:"state" db:"org_user_state"`
}
