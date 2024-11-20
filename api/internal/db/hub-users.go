package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type HubUserTFA struct {
	TFACode  string
	TFAToken HubTokenReq
	Email    Email
}

type HubUserTO struct {
	ID           uuid.UUID           `db:"id"`
	FullName     string              `db:"full_name"`
	Handle       string              `db:"handle"`
	Email        vetchi.EmailAddress `db:"email"`
	PasswordHash string              `db:"password_hash"`
	State        vetchi.HubUserState `db:"state"`
	CreatedAt    time.Time           `db:"created_at"`
	UpdatedAt    time.Time           `db:"updated_at"`
}

type HubUserPasswordResetReq struct {
	Email Email
	HubTokenReq
}
