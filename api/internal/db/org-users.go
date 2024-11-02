package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// OrgUserTO should be used within the API server and
// should not be exposed to the outside world
type OrgUserTO struct {
	ID           uuid.UUID            `db:"id"             json:"-"`
	Name         string               `db:"name"           json:"-"`
	Email        string               `db:"email"          json:"-"`
	PasswordHash string               `db:"password_hash"  json:"-"`
	OrgUserRoles []vetchi.OrgUserRole `db:"org_user_roles" json:"-"`
	OrgUserState vetchi.OrgUserState  `db:"org_user_state" json:"-"`
	EmployerID   uuid.UUID            `db:"employer_id"    json:"-"`
	CreatedAt    time.Time            `db:"created_at"     json:"-"`
}

type AddOrgUserReq struct {
	Name         string
	Email        string
	OrgUserRoles []vetchi.OrgUserRole
	OrgUserState vetchi.OrgUserState

	InviteMail Email

	InviteToken TokenReq

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	AddingUserID uuid.UUID
}

type DisableOrgUserReq struct {
	Email      string
	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	DisablingUserID uuid.UUID
}

type FilterOrgUsersReq struct {
	Prefix     string
	State      []vetchi.OrgUserState
	EmployerID uuid.UUID

	PaginationKey string
	Limit         int
}

type UpdateOrgUserReq struct {
	Email      string
	Name       string
	Roles      []vetchi.OrgUserRole
	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	UpdatingUserID uuid.UUID
}
