package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
)

func getHubUserID(ctx context.Context) (string, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		return "", db.ErrInternal
	}
	return hubUser.ID.String(), nil
}
