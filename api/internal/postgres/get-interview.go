package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetInterview(
	ctx context.Context,
	interviewID string,
) (employer.EmployerInterview, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.EmployerInterview{}, db.ErrInternal
	}

	// TODO: Right now we are making the getInterview and getInterviews
	// calls for all the orgUsers. We should do some RBAC on this.
	query := `
		SELECT 
			i.id as interview_id,
			i.interview_state as interview_state,
			i.start_time,
			i.end_time,
			i.interview_type,
			i.description,
			i.interviewers_decision,
			i.positives,
			i.negatives,
			i.overall_assessment,
			i.feedback_to_candidate,
			i.feedback_submitted_by,
			i.feedback_submitted_at,
			i.created_at
		FROM interviews i
		WHERE i.id = $1
		AND i.employer_id = $2
	`

	var interview employer.EmployerInterview
	err := p.pool.QueryRow(ctx, query, interviewID, orgUser.EmployerID).Scan(
		&interview.InterviewID,
		&interview.InterviewState,
		&interview.StartTime,
		&interview.EndTime,
		&interview.InterviewType,
		&interview.Description,
		&interview.InterviewersDecision,
		&interview.Positives,
		&interview.Negatives,
		&interview.OverallAssessment,
		&interview.FeedbackToCandidate,
		&interview.FeedbackSubmittedBy,
		&interview.FeedbackSubmittedAt,
		&interview.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.log.Dbg("interview not found")
			return employer.EmployerInterview{}, db.ErrNoInterview
		}
		p.log.Err("failed to get interview", "error", err)
		return employer.EmployerInterview{}, db.ErrInternal
	}

	// Get interviewers
	interviewersQuery := `
		SELECT 
			ou.name,
			ou.email
		FROM interview_interviewers ii
		JOIN org_users ou ON ii.interviewer_id = ou.id
		WHERE ii.interview_id = $1
		AND ii.employer_id = $2
	`

	rows, err := p.pool.Query(
		ctx,
		interviewersQuery,
		interviewID,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to get interviewers", "error", err)
		return employer.EmployerInterview{}, db.ErrInternal
	}
	defer rows.Close()

	var interviewers []employer.OrgUserTiny
	for rows.Next() {
		var interviewer employer.OrgUserTiny
		err := rows.Scan(
			&interviewer.Name,
			&interviewer.Email,
		)
		if err != nil {
			p.log.Err("failed to scan interviewer", "error", err)
			return employer.EmployerInterview{}, db.ErrInternal
		}
		interviewers = append(interviewers, interviewer)
	}
	interview.Interviewers = interviewers

	return interview, nil
}
