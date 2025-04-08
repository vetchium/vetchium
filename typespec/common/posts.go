package common

import (
	"time"
)

type Post struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	AuthorName   string    `json:"author_name"`
	AuthorHandle Handle    `json:"author_handle"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
