package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) GetOldestUnsentEmails(ctx context.Context) ([]db.Email, error) {
	query := `
SELECT
    email_key,
    email_from,
    email_to,
    email_subject,
    email_html_body,
    email_text_body,
    email_state
FROM
    emails
WHERE
    email_state = $1
ORDER BY
    created_at ASC
LIMIT 10
`
	rows, err := p.pool.Query(ctx, query, db.EmailStatePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []db.Email
	for rows.Next() {
		var email db.Email
		err := rows.Scan(
			&email.EmailKey,
			&email.EmailFrom,
			&email.EmailTo,
			&email.EmailSubject,
			&email.EmailHTMLBody,
			&email.EmailTextBody,
			&email.EmailState,
		)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (p *PG) UpdateEmailState(
	ctx context.Context,
	emailStateChange db.EmailStateChange,
) error {
	query := `
UPDATE
    emails
SET
    email_state = $1,
    processed_at = NOW()
WHERE
    email_key = $2
`
	_, err := p.pool.Exec(
		ctx,
		query,
		emailStateChange.EmailState,
		emailStateChange.EmailDBKey,
	)
	if err != nil {
		p.log.Err("failed to update email state", "error", err)
		return err
	}

	return nil
}
