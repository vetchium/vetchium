package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) AddWorkHistory(
	ctx context.Context,
	req hub.AddWorkHistoryRequest,
) (string, error) {
	return "", nil
}

func (p *PG) DeleteWorkHistory(
	ctx context.Context,
	req hub.DeleteWorkHistoryRequest,
) error {
	return nil
}

func (p *PG) ListWorkHistory(
	ctx context.Context,
	req hub.ListWorkHistoryRequest,
) ([]hub.WorkHistory, error) {
	return nil, nil
}

func (p *PG) UpdateWorkHistory(
	ctx context.Context,
	req hub.UpdateWorkHistoryRequest,
) error {
	return nil
}
