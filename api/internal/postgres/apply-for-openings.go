package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) CreateApplication(
	ctx context.Context,
	req db.ApplyOpeningReq,
) error {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		return db.ErrNoHubUser
	}

	query := `
INSERT INTO applications (id, employer_id, opening_id, cover_letter, original_filename, internal_filename, hub_user_id, application_state)
    VALUES ($1, (
            SELECT
                employer_id
            FROM
                domains
            WHERE
                DOMAIN = $2), $3, $4, $5, $6, $7, $8)
RETURNING
    id
`

	var applicationID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		req.ApplicationID,
		req.CompanyDomain,
		req.OpeningIDWithinCompany,
		req.CoverLetter,
		req.OriginalFilename,
		req.InternalFilename,
		hubUser.ID,
		vetchi.AppliedAppState,
	).Scan(&applicationID)
	if err != nil {
		p.log.Err("failed to create application", "error", err)
		return db.ErrInternal
	}
	p.log.Dbg("created application", "application_id", applicationID)

	return nil
}
