package postgres

import (
	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) AddOfficialEmail(req db.AddOfficialEmailReq) error {
	ctx := req.Context

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Add the official email to the hub_users_official_emails table
	officialEmailsQuery := ``

	err = tx.QueryRow(ctx, officialEmailsQuery).Scan()
	if err != nil {
		return err
	}

	// Send the email with the token to the added official email address
	tokenMailQuery := ``

	err = tx.QueryRow(ctx, tokenMailQuery).Scan()
	if err != nil {
		return err
	}

	return nil
}
