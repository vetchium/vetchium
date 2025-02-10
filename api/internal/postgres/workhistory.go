package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) AddWorkHistory(
	ctx context.Context,
	req hub.AddWorkHistoryRequest,
) (string, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return "", err
	}

	// Start a transaction since we might need to create a domain
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return "", db.ErrInternal
	}
	defer tx.Rollback(ctx)

	// Get or create domain
	var domainID string
	err = tx.QueryRow(ctx, `
		WITH domain_insert AS (
			INSERT INTO domains (domain_name, domain_state)
			VALUES ($1, $2)
			ON CONFLICT (domain_name) DO NOTHING
			RETURNING id
		)
		SELECT id FROM domain_insert
		UNION ALL
		SELECT id FROM domains WHERE domain_name = $1
		LIMIT 1
	`, req.EmployerDomain, db.UnverifiedDomainState).Scan(&domainID)
	if err != nil {
		return "", db.ErrInternal
	}

	// Insert work history
	var id string
	err = tx.QueryRow(ctx, `
		INSERT INTO work_history (
			hub_user_id,
			domain_id,
			title,
			start_date,
			end_date,
			description
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`, hubUserID, domainID, req.Title, req.StartDate, req.EndDate, req.Description).Scan(&id)

	if err != nil {
		return "", db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", db.ErrInternal
	}

	return id, nil
}

func (p *PG) DeleteWorkHistory(
	ctx context.Context,
	req hub.DeleteWorkHistoryRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	tag, err := p.pool.Exec(ctx, `
		DELETE FROM work_history
		WHERE id = $1 AND hub_user_id = $2
	`, req.ID, hubUserID)

	if err != nil {
		return db.ErrInternal
	}

	if tag.RowsAffected() == 0 {
		return db.ErrNoWorkHistory
	}

	return nil
}

func (p *PG) ListWorkHistory(
	ctx context.Context,
	req hub.ListWorkHistoryRequest,
) ([]hub.WorkHistory, error) {
	var userID string
	if req.UserHandle != nil && *req.UserHandle != "" {
		err := p.pool.QueryRow(ctx, `
			SELECT id FROM hub_users WHERE handle = $1
		`, *req.UserHandle).Scan(&userID)
		if err == pgx.ErrNoRows {
			return nil, db.ErrNoHubUser
		}
		if err != nil {
			return nil, db.ErrInternal
		}
	} else {
		var err error
		userID, err = getHubUserID(ctx)
		if err != nil {
			return nil, err
		}
	}

	rows, err := p.pool.Query(ctx, `
		SELECT 
			w.id,
			d.domain_name,
			e.company_name,
			w.title,
			w.start_date,
			w.end_date,
			w.description
		FROM work_history w
		JOIN domains d ON d.id = w.domain_id
		LEFT JOIN employers e ON e.id = d.employer_id AND e.employer_state = $2
		WHERE w.hub_user_id = $1
		ORDER BY w.start_date DESC
	`, userID, db.OnboardedEmployerState)

	if err != nil {
		return nil, db.ErrInternal
	}
	defer rows.Close()

	var workHistories []hub.WorkHistory
	for rows.Next() {
		var wh hub.WorkHistory
		var companyName, endDate, description pgtype.Text

		err := rows.Scan(
			&wh.ID,
			&wh.EmployerDomain,
			&companyName,
			&wh.Title,
			&wh.StartDate,
			&endDate,
			&description,
		)
		if err != nil {
			return nil, db.ErrInternal
		}

		if companyName.Valid {
			wh.EmployerName = &companyName.String
		}
		if endDate.Valid {
			wh.EndDate = &endDate.String
		}
		if description.Valid {
			wh.Description = &description.String
		}

		workHistories = append(workHistories, wh)
	}

	return workHistories, nil
}

func (p *PG) UpdateWorkHistory(
	ctx context.Context,
	req hub.UpdateWorkHistoryRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	tag, err := p.pool.Exec(ctx, `
		UPDATE work_history
		SET 
			title = $1,
			start_date = $2,
			end_date = $3,
			description = $4,
			updated_at = $5
		WHERE id = $6 AND hub_user_id = $7
	`, req.Title, req.StartDate, req.EndDate, req.Description, time.Now().UTC(), req.ID, hubUserID)

	if err != nil {
		return db.ErrInternal
	}

	if tag.RowsAffected() == 0 {
		return db.ErrNoWorkHistory
	}

	return nil
}
