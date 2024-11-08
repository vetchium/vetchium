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
	OrgUserRoles vetchi.OrgUserRoles
	OrgUserState vetchi.OrgUserState

	InviteMail Email

	InviteToken TokenReq

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	AddingUserID uuid.UUID
}

type EnableOrgUserReq struct {
	Email      string
	EmployerID uuid.UUID

	InviteMail  Email
	InviteToken TokenReq
	// Currently unused, but will be used in the future for audit logs
	EnablingUserID uuid.UUID
}

type SignupOrgUserReq struct {
	InviteToken  string
	Name         string
	PasswordHash string
}
