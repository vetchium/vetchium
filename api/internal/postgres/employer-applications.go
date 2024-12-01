package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) GetApplicationsForEmployer(
	c context.Context,
	req vetchi.GetApplicationsRequest,
) ([]vetchi.Application, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return "", db.ErrInternal
	}

	query := `
SELECT * FROM applications WHERE employer_id = $1
`
}
