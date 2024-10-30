package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type OrgUser struct {
	ID           uuid.UUID            `db:"id"`
	Email        string               `db:"email"`
	PasswordHash string               `db:"password_hash"`
	OrgUserRoles []vetchi.OrgUserRole `db:"org_user_roles"`
	OrgUserState OrgUserState         `db:"org_user_state"`
	EmployerID   uuid.UUID            `db:"employer_id"`
	CreatedAt    time.Time            `db:"created_at"`
}

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

type AddOrgUserReq struct {
	Email        string
	OrgUserRoles []vetchi.OrgUserRole
	OrgUserState OrgUserState
	EmployerID   uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	AddingUser uuid.UUID
}
