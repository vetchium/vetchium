package postgres

import (
	"context"
	"fmt"

	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) GetMyEndorsementApprovals(
	ctx context.Context,
	req hub.MyEndorseApprovalsRequest,
) (hub.MyEndorseApprovalsResponse, error) {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return hub.MyEndorseApprovalsResponse{}, err
	}

	// Prepare arguments
	var args []interface{}
	args = append(args, loggedInUserID)
	args = append(args, req.State)

	// Build the query to fetch endorsement requests
	query := `
		SELECT
			ae.id,
			ae.application_id,
			ae.state,
			hu_applicant.handle as applicant_handle,
			hu_applicant.full_name as applicant_name,
			hu_applicant.short_bio as applicant_short_bio,
			e.company_name as employer_name,
			d.domain_name as employer_domain,
			o.title as opening_title,
			concat('https://', d.domain_name, '/openings/', o.id) as opening_url,
			a.application_state as application_status,
			a.created_at as application_created_at
		FROM application_endorsements ae
		JOIN applications a ON ae.application_id = a.id
		JOIN hub_users hu_applicant ON a.hub_user_id = hu_applicant.id
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		JOIN employers e ON a.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		WHERE ae.endorser_id = $1
		AND ae.state = ANY($2)
	`

	// Add pagination
	if req.PaginationKey != nil {
		// Get the timestamp and ID from the pagination key
		query += fmt.Sprintf(`
			AND (ae.created_at, ae.id) > (
				SELECT created_at, id
				FROM application_endorsements
				WHERE id = $%d
			)
		`, len(args)+1)
		args = append(args, *req.PaginationKey)
	}

	// Add ordering and limit
	query += fmt.Sprintf(`
		ORDER BY ae.created_at, ae.id
		LIMIT $%d
	`, len(args)+1)

	args = append(args, req.Limit)

	// Execute the query
	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to query endorsement approvals", "error", err)
		return hub.MyEndorseApprovalsResponse{}, err
	}
	defer rows.Close()

	// Process the results
	endorsements := make([]hub.MyEndorseApproval, 0)
	var lastID string
	for rows.Next() {
		var endorsement hub.MyEndorseApproval
		var id string

		err := rows.Scan(
			&id,
			&endorsement.ApplicationID,
			&endorsement.EndorsementStatus,
			&endorsement.ApplicantHandle,
			&endorsement.ApplicantName,
			&endorsement.ApplicantShortBio,
			&endorsement.EmployerName,
			&endorsement.EmployerDomain,
			&endorsement.OpeningTitle,
			&endorsement.OpeningURL,
			&endorsement.ApplicationStatus,
			&endorsement.ApplicationCreatedAt,
		)
		if err != nil {
			p.log.Err("failed to scan endorsement approval", "error", err)
			return hub.MyEndorseApprovalsResponse{}, err
		}

		lastID = id
		endorsements = append(endorsements, endorsement)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over endorsement approvals", "error", err)
		return hub.MyEndorseApprovalsResponse{}, err
	}

	// Return the response with pagination key if there are results
	response := hub.MyEndorseApprovalsResponse{
		Endorsements: endorsements,
	}

	if len(endorsements) > 0 {
		response.PaginationKey = lastID
	}

	return response, nil
}
