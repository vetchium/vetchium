package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

// OrgUserTO should be used within the API server and
// should not be exposed to the outside world
type OrgUserTO struct {
	ID           uuid.UUID             `db:"id"             json:"-"`
	Name         string                `db:"name"           json:"-"`
	Email        string                `db:"email"          json:"-"`
	PasswordHash string                `db:"password_hash"  json:"-"`
	OrgUserRoles []common.OrgUserRole  `db:"org_user_roles" json:"-"`
	OrgUserState employer.OrgUserState `db:"org_user_state" json:"-"`
	EmployerID   uuid.UUID             `db:"employer_id"    json:"-"`
	CreatedAt    time.Time             `db:"created_at"     json:"-"`
}

type AddOrgUserReq struct {
	Name         string
	Email        string
	OrgUserRoles common.OrgUserRoles
	OrgUserState employer.OrgUserState

	InviteMail Email

	InviteToken OrgUserInviteReq

	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	AddingUserID uuid.UUID
}

type EnableOrgUserReq struct {
	Email      string
	EmployerID uuid.UUID

	InviteMail  Email
	InviteToken OrgUserInviteReq

	// Currently unused, but will be used in the future for audit logs
	EnablingUserID uuid.UUID
}

type SignupOrgUserReq struct {
	InviteToken  string
	Name         string
	PasswordHash string
}
