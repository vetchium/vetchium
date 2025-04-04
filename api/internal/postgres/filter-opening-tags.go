package postgres

import (
	"context"
	"fmt"

	"github.com/vetchium/vetchium/typespec/common"
)

func (p *PG) FilterOpeningTags(
	ctx context.Context,
	req common.FilterOpeningTagsRequest,
) ([]common.OpeningTag, error) {
	query := `
SELECT id, name
FROM opening_tags
WHERE 1=1
`
	args := make([]interface{}, 0)
	argPos := 1

	if req.Prefix != nil {
		// TODO: Use Semantic matching search instead of just prefix
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, fmt.Sprintf("%s%%", *req.Prefix))
		argPos++
	}

	query += ` ORDER BY name ASC`

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to filter opening tags", "error", err)
		return nil, err
	}
	defer rows.Close()

	var tags []common.OpeningTag
	for rows.Next() {
		var tag common.OpeningTag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			p.log.Err("failed to scan opening tag", "error", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	p.log.Dbg("filtered opening tags", "tags", tags)

	return tags, nil
}
