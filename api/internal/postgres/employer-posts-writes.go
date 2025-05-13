package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) AddEmployerPost(req db.AddEmployerPostRequest) error {
	return nil
}

func (p *PG) UpdateEmployerPost(req db.UpdateEmployerPostRequest) error {
	return nil
}

func (p *PG) DeleteEmployerPost(ctx context.Context, postID string) error {
	return nil
}
