package common

import (
	"time"
)

type Post struct {
	ID             string    `json:"id"`
	Content        string    `json:"content"`
	Tags           []string  `json:"tags"`
	AuthorName     string    `json:"author_name"`
	AuthorHandle   Handle    `json:"author_handle"`
	CreatedAt      time.Time `json:"created_at"`
	UpvotesCount   int32     `json:"upvotes_count"`
	DownvotesCount int32     `json:"downvotes_count"`
	Score          int32     `json:"score"`
}
