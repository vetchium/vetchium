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

	var id string
	err = p.pool.QueryRow(ctx, `
		INSERT INTO work_history (
			hub_user_id,
			employer_domain,
			title,
			start_date,
			end_date,
			description
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`, hubUserID, req.EmployerDomain, req.Title, req.StartDate, req.EndDate, req.Description).Scan(&id)

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
			w.employer_domain,
			e.company_name,
			w.title,
			w.start_date,
			w.end_date,
			w.description
		FROM work_history w
		LEFT JOIN employers e ON e.client_id_type = 'DOMAIN' AND e.employer_state = 'ONBOARDED'
		LEFT JOIN domains d ON d.employer_id = e.id AND d.domain_name = w.employer_domain
		WHERE w.hub_user_id = $1
		ORDER BY w.start_date DESC
	`, userID)

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
