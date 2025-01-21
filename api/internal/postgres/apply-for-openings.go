package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
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
WITH employer AS (
    SELECT employer_id
    FROM domains
    WHERE domain_name = $2
),
valid_opening AS (
    SELECT 1
    FROM openings
    WHERE employer_id = (SELECT employer_id FROM employer)
      AND id = $3
)
INSERT INTO applications (
    id, employer_id, opening_id, cover_letter,
    resume_sha, hub_user_id, application_state
)
SELECT
    $1, (SELECT employer_id FROM employer), $3, $4, $5, $6, $7
WHERE EXISTS (SELECT 1 FROM valid_opening)
RETURNING id
`

	var applicationID string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.ApplicationID,
		req.CompanyDomain,
		req.OpeningIDWithinCompany,
		req.CoverLetter,
		req.ResumeSHA,
		hubUser.ID,
		common.AppliedAppState,
	).Scan(&applicationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("either domain or opening does not exist", "error", err)
			return db.ErrNoOpening
		}
		p.log.Err("failed to create application", "error", err)
		return db.ErrInternal
	}
	p.log.Dbg("created application", "application_id", applicationID)

	return nil
}
