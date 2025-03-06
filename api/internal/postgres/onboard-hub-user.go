package postgres

import (
	"context"
	"fmt"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) OnboardHubUser(
	ctx context.Context,
	onboardHubUserReq db.OnboardHubUserReq,
) (string, error) {
	p.log.Dbg("Onboarding hub user", "request", onboardHubUserReq)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("Failed to begin transaction", "error", err)
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create new user with generated handle, using email from invite
	userQuery := `
WITH invite AS (
	SELECT email FROM hub_user_invites WHERE token = $1
),
new_handle AS (
	SELECT generate_unique_handle($2) as handle
)
INSERT INTO hub_users (
	full_name,
	handle,
	email,
	password_hash,
	tier,
	resident_country_code,
	preferred_language,
	short_bio,
	long_bio,
	state
)
SELECT
	$2, handle, email, $3, $4, $5, $6, $7, $8, $9
FROM new_handle, invite
RETURNING handle`

	var handle string
	err = tx.QueryRow(
		ctx,
		userQuery,
		onboardHubUserReq.InviteToken,
		onboardHubUserReq.FullName,
		onboardHubUserReq.PasswordHash,
		onboardHubUserReq.Tier,
		onboardHubUserReq.ResidentCountryCode,
		onboardHubUserReq.PreferredLanguage,
		onboardHubUserReq.ShortBio,
		onboardHubUserReq.LongBio,
		hub.ActiveHubUserState,
	).Scan(&handle)
	if err != nil {
		return "", fmt.Errorf("failed to create hub user: %w", err)
	}

	// Delete the invite token
	_, err = tx.Exec(
		ctx,
		`DELETE FROM hub_user_invites WHERE token = $1`,
		onboardHubUserReq.InviteToken,
	)
	if err != nil {
		p.log.Err("Failed to delete invite token", "error", err)
		return "", fmt.Errorf("failed to delete invite token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		p.log.Err("Failed to commit transaction", "error", err)
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	p.log.Dbg("Hub user onboarded", "handle", handle)
	return handle, nil
}
