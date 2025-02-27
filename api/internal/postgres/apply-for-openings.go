package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) CreateApplication(
	ctx context.Context,
	req db.ApplyOpeningReq,
) error {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hub user", "error", db.ErrNoHubUser)
		return db.ErrNoHubUser
	}

	// Start a transaction
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	// Validate endorsers are colleagues (if any)
	if len(req.EndorserHandles) > 0 {
		// Check if all endorsers are colleagues
		for _, handle := range req.EndorserHandles {
			var isColleague bool
			err := tx.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM colleague_connections
					WHERE (
						(requester_id = $1 AND requested_id = (SELECT id FROM hub_users WHERE handle = $2))
						OR 
						(requester_id = (SELECT id FROM hub_users WHERE handle = $2) AND requested_id = $1)
					)
					AND state = $3
				)
			`, hubUser.ID, string(handle), db.ColleagueAccepted).Scan(&isColleague)
			if err != nil {
				p.log.Err("failed to check colleague status", "error", err)
				return db.ErrInternal
			}

			if !isColleague {
				p.log.Dbg("endorser is not a colleague", "handle", handle)
				return db.ErrNotColleague
			}
		}
	}

	// Create the application
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
	err = tx.QueryRow(
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

	// Create endorsement requests if any
	for _, handle := range req.EndorserHandles {
		// Get endorser ID
		var endorserID string
		err := tx.QueryRow(ctx, `
			SELECT id FROM hub_users WHERE handle = $1
		`, string(handle)).Scan(&endorserID)
		if err != nil {
			p.log.Err("failed to get endorser", "error", err, "handle", handle)
			return db.ErrInternal
		}

		// Create endorsement record with auto-generated UUID
		_, err = tx.Exec(ctx, `
			INSERT INTO application_endorsements (
				application_id, endorser_id, state
			) VALUES (
				$1, $2, $3
			)
		`, applicationID, endorserID, hub.SoughtEndorsement)
		if err != nil {
			p.log.Err("failed to create endorsement", "error", err)
			return db.ErrInternal
		}

		p.log.Dbg("endorsement requested", "endorser", handle)
	}

	// Send emails to endorsers
	for _, email := range req.EndorsementEmails {
		emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
VALUES ($1, $2, $3, $4, $5, $6)
`

		_, err = tx.Exec(ctx, emailQuery,
			email.EmailFrom,
			email.EmailTo,
			email.EmailSubject,
			email.EmailHTMLBody,
			email.EmailTextBody,
			email.EmailState,
		)
		if err != nil {
			p.log.Err("failed to send email to endorser", "error", err)
			return db.ErrInternal
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

// GetHubUsersByHandles gets the details of hub users by their handles
func (p *PG) GetHubUsersByHandles(
	ctx context.Context,
	handles []common.Handle,
) ([]db.HubUserContact, error) {
	if len(handles) == 0 {
		return []db.HubUserContact{}, nil
	}

	// Convert handles to a format suitable for SQL IN clause
	handleParams := make([]string, len(handles))
	for i, handle := range handles {
		handleParams[i] = string(handle)
	}

	query := `
SELECT handle, full_name, email
FROM hub_users
WHERE handle = ANY($1)
`

	rows, err := p.pool.Query(ctx, query, handleParams)
	if err != nil {
		p.log.Err("failed to query hub users", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	var users []db.HubUserContact
	for rows.Next() {
		var user db.HubUserContact
		err := rows.Scan(&user.Handle, &user.FullName, &user.Email)
		if err != nil {
			p.log.Err("failed to scan hub user", "error", err)
			return nil, db.ErrInternal
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("error iterating over hub users", "error", err)
		return nil, db.ErrInternal
	}

	return users, nil
}
