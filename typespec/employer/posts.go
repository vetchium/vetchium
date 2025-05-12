package employer

type EmployerPost struct {
	ID            string   `json:"id"`
	Content       string   `json:"content"`
	Tags          []string `json:"tags"`
	CompanyDomain string   `json:"company_domain"`
	CreatedAt     string   `json:"created_at"`
}

type AddEmployerPostRequest struct {
	Content string   `json:"content" validate:"required,min=1,max=4096"`
	Tags    []string `json:"tags"    validate:"max=3"`
}

type AddEmployerPostResponse struct {
	PostID string `json:"post_id"`
}

type UpdateEmployerPostRequest struct {
	PostID  string   `json:"post_id"`
	Content string   `json:"content" validate:"required,min=1,max=4096"`
	Tags    []string `json:"tags"    validate:"max=3"`
}

type DeleteEmployerPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}
