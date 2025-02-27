package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) EndorseApplication(
	ctx context.Context,
	endorseReq hub.EndorseApplicationRequest,
) error {
	p.log.Dbg("Endorsing application", "endorseReq", endorseReq)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	query := `
UPDATE application_endorsements
SET state = $1
WHERE application_id = $2
AND endorser_id = $3
AND state = $4
`

	args := []interface{}{
		hub.Endorsed,
		endorseReq.ApplicationID,
		hubUserID,
		hub.SoughtEndorsement,
	}

	result, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to execute query", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		p.log.Err("no endorsement state change", "endorseReq", endorseReq)
		return db.ErrNoApplication
	}

	p.log.Dbg("Endorsed application", "endorseReq", endorseReq)
	return nil
}

func (p *PG) RejectEndorsement(
	ctx context.Context,
	rejectReq hub.RejectEndorsementRequest,
) error {
	p.log.Dbg("Rejecting endorsement", "rejectReq", rejectReq)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	query := `
UPDATE application_endorsements
SET state = $1
WHERE application_id = $2
AND endorser_id = $3
AND state = $4
`

	args := []interface{}{
		hub.DeclinedEndorsement,
		rejectReq.ApplicationID,
		hubUserID,
		hub.SoughtEndorsement,
	}

	result, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to execute query", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		p.log.Err("no endorsement state change", "rejectReq", rejectReq)
		return db.ErrNoApplication
	}

	p.log.Dbg("Rejected endorsement", "rejectReq", rejectReq)
	return nil
}
