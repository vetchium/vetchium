package sortinghat

// ApplicationSortRequest represents a request to score a single application's resume
type ApplicationSortRequest struct {
	// ApplicationID is the unique identifier for the application
	ApplicationID string `json:"application_id"`
	// ResumePath is the S3 path to the resume file in format s3://bucket/key
	ResumePath string `json:"resume_path"`
}

// SortingHatRequest represents a request to score multiple resumes against a job description in a batch
type SortingHatRequest struct {
	// JobDescription is the job description to score resumes against
	JobDescription string `json:"job_description"`
	// ApplicationSortRequests is the list of applications to score
	ApplicationSortRequests []ApplicationSortRequest `json:"application_sort_requests"`
}

// ModelScore represents a score from a specific model
type ModelScore struct {
	// ModelName is the name of the model that generated the score
	ModelName string `json:"model_name"`
	// Score is the score value from 0 to 100
	Score int32 `json:"score"`
}

// SortingHatScore represents scores for a single application from all models
type SortingHatScore struct {
	// ApplicationID is the application ID this score relates to
	ApplicationID string `json:"application_id"`
	// ModelScores are the scores from different models
	ModelScores []ModelScore `json:"model_scores"`
}

// SortingHatResponse represents the response containing scores for all applications in the batch
type SortingHatResponse struct {
	// Scores is the list of application scores
	Scores []SortingHatScore `json:"scores"`
}
