package postgres

import (
	"context"
	"errors"
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
		p.log.Err("failed to get hub user ID", "error", err)
		return "", err
	}

	// Start a transaction since we might need to create a domain
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return "", db.ErrInternal
	}
	defer tx.Rollback(ctx)

	// Get or create employer for the domain
	var employerID string
	err = tx.QueryRow(ctx, `
		SELECT get_or_create_dummy_employer($1)
	`, req.EmployerDomain).Scan(&employerID)
	if err != nil {
		p.log.Err(
			"failed to get or create dummy employer",
			"error",
			err,
			"domain",
			req.EmployerDomain,
		)
		return "", db.ErrInternal
	}

	// Insert work history
	var id string
	err = tx.QueryRow(ctx, `
		INSERT INTO work_history (
			hub_user_id,
			employer_id,
			title,
			start_date,
			end_date,
			description
		) VALUES (
			$1, $2, $3, $4::DATE, $5::DATE, $6
		) RETURNING id
	`, hubUserID, employerID, req.Title, req.StartDate, req.EndDate, req.Description).Scan(&id)

	if err != nil {
		p.log.Err("failed to insert work history", "error", err)
		return "", db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
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
		p.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tag, err := p.pool.Exec(ctx, `
		DELETE FROM work_history
		WHERE id = $1 AND hub_user_id = $2
	`, req.ID, hubUserID)

	if err != nil {
		p.log.Err(
			"failed to delete work history",
			"error",
			err,
			"work_history_id",
			req.ID,
		)
		return db.ErrInternal
	}

	if tag.RowsAffected() == 0 {
		p.log.Dbg(
			"work history not found or not owned by user",
			"work_history_id",
			req.ID,
			"hub_user_id",
			hubUserID,
		)
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
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("hub user not found", "handle", *req.UserHandle)
			return nil, db.ErrNoHubUser
		}
		if err != nil {
			p.log.Err(
				"failed to get hub user by handle",
				"error",
				err,
				"handle",
				*req.UserHandle,
			)
			return nil, db.ErrInternal
		}
	} else {
		var err error
		userID, err = getHubUserID(ctx)
		if err != nil {
			p.log.Err("failed to get hub user ID", "error", err)
			return nil, err
		}
	}

	rows, err := p.pool.Query(ctx, `
		SELECT 
			w.id,
			d.domain_name,
			e.company_name,
			w.title,
			w.start_date::TEXT,
			w.end_date::TEXT,
			w.description
		FROM work_history w
		JOIN employers e ON e.id = w.employer_id
		JOIN employer_primary_domains epd ON epd.employer_id = e.id
		JOIN domains d ON d.id = epd.domain_id
		WHERE w.hub_user_id = $1
		ORDER BY w.start_date DESC
	`, userID)

	if err != nil {
		p.log.Err(
			"failed to query work history",
			"error",
			err,
			"hub_user_id",
			userID,
		)
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
			p.log.Err("failed to scan work history row", "error", err)
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

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over work history rows", "error", err)
		return nil, db.ErrInternal
	}

	return workHistories, nil
}

func (p *PG) UpdateWorkHistory(
	ctx context.Context,
	req hub.UpdateWorkHistoryRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tag, err := p.pool.Exec(ctx, `
		UPDATE work_history
		SET 
			title = $1,
			start_date = $2::DATE,
			end_date = $3::DATE,
			description = $4,
			updated_at = $5
		WHERE id = $6 AND hub_user_id = $7
	`, req.Title, req.StartDate, req.EndDate, req.Description, time.Now().UTC(), req.ID, hubUserID)

	if err != nil {
		p.log.Err(
			"failed to update work history",
			"error",
			err,
			"work_history_id",
			req.ID,
		)
		return db.ErrInternal
	}

	if tag.RowsAffected() == 0 {
		p.log.Dbg(
			"work history not found or not owned by user",
			"work_history_id",
			req.ID,
			"hub_user_id",
			hubUserID,
		)
		return db.ErrNoWorkHistory
	}

	return nil
}
