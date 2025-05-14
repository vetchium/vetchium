package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func (p *PG) GetEmployerPost(
	ctx context.Context,
	postID string,
) (common.EmployerPost, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return common.EmployerPost{}, db.ErrInternal
	}

	// Query to get post details and tags
	query := `
	SELECT 
		ep.id, 
		ep.content, 
		ep.created_at, 
		ep.updated_at,
		d.domain_name as company_domain,
		ARRAY_AGG(t.name) FILTER (WHERE t.name IS NOT NULL) as tags
	FROM employer_posts ep
	JOIN employers e ON ep.employer_id = e.id
	JOIN employer_primary_domains epd ON e.id = epd.employer_id
	JOIN domains d ON epd.domain_id = d.id
	LEFT JOIN employer_post_tags ept ON ep.id = ept.employer_post_id
	LEFT JOIN tags t ON ept.tag_id = t.id
	WHERE ep.id = $1 AND ep.employer_id = $2
	GROUP BY ep.id, ep.content, ep.created_at, ep.updated_at, d.domain_name
	`

	var post common.EmployerPost
	var tags []string

	err := p.pool.QueryRow(ctx, query, postID, orgUser.EmployerID).Scan(
		&post.ID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.EmployerDomainName,
		&tags,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("employer post not found",
				"post_id", postID,
				"employer_id", orgUser.EmployerID,
			)
			return common.EmployerPost{}, db.ErrNoEmployerPost
		}

		p.log.Err("failed to scan employer post", "error", err)
		return common.EmployerPost{}, db.ErrInternal
	}

	// Set tags if they exist
	if len(tags) > 0 && tags[0] != "" {
		post.Tags = tags
	} else {
		post.Tags = []string{}
	}

	p.log.Dbg("fetched employer post", "post_id", postID)
	return post, nil
}

func (p *PG) ListEmployerPosts(
	req db.ListEmployerPostsRequest,
) (employer.ListEmployerPostsResponse, error) {
	orgUser, ok := req.Context.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.ListEmployerPostsResponse{}, db.ErrInternal
	}

	// Use the limit from the request
	limit := req.Limit

	// Build the query with pagination
	query := `
	WITH ranked_posts AS (
		SELECT 
			ep.id, 
			ep.content, 
			ep.created_at, 
			ep.updated_at,
			d.domain_name as company_domain,
			ARRAY_AGG(t.name) FILTER (WHERE t.name IS NOT NULL) as tags
		FROM employer_posts ep
		JOIN employers e ON ep.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		LEFT JOIN employer_post_tags ept ON ep.id = ept.employer_post_id
		LEFT JOIN tags t ON ept.tag_id = t.id
		WHERE ep.employer_id = $1
`

	// Add pagination condition if pagination_key is provided
	args := []interface{}{orgUser.EmployerID}
	if req.PaginationKey != "" {
		query += `		AND (ep.updated_at < (SELECT updated_at FROM employer_posts WHERE id = $2) 
		      OR (ep.updated_at = (SELECT updated_at FROM employer_posts WHERE id = $2) AND ep.id < $2))
`
		args = append(args, req.PaginationKey)
	}

	// Complete the query with grouping, ordering and limit
	query += `		GROUP BY ep.id, ep.content, ep.created_at, ep.updated_at, d.domain_name
		ORDER BY ep.updated_at DESC, ep.id DESC
		LIMIT $` + strconv.Itoa(
		len(args)+1,
	) + `
	)
	SELECT * FROM ranked_posts
	`

	args = append(args, limit)

	// Execute the query
	rows, err := p.pool.Query(req.Context, query, args...)
	if err != nil {
		p.log.Err("failed to query employer posts", "error", err)
		return employer.ListEmployerPostsResponse{}, db.ErrInternal
	}
	defer rows.Close()

	// Process the results
	posts := []common.EmployerPost{}
	var lastPostID string

	for rows.Next() {
		var post common.EmployerPost
		var tags []string

		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.EmployerDomainName,
			&tags,
		)

		if err != nil {
			p.log.Err("failed to scan employer post", "error", err)
			return employer.ListEmployerPostsResponse{}, db.ErrInternal
		}

		// Set tags if they exist
		if len(tags) > 0 && tags[0] != "" {
			post.Tags = tags
		} else {
			post.Tags = []string{}
		}

		posts = append(posts, post)
		lastPostID = post.ID
	}

	if rows.Err() != nil {
		p.log.Err("error iterating employer post rows", "error", rows.Err())
		return employer.ListEmployerPostsResponse{}, rows.Err()
	}

	// Set the next pagination key if there are more posts
	var nextPaginationKey string
	if len(posts) == limit {
		nextPaginationKey = lastPostID
	}

	p.log.Dbg(
		"listed employer posts",
		"count",
		len(posts),
		"employer_id",
		orgUser.EmployerID,
	)
	return employer.ListEmployerPostsResponse{
		Posts:         posts,
		PaginationKey: nextPaginationKey,
	}, nil
}
