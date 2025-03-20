package vetchi

// Types here should match the ones in sortinghat/main.py

type ApplicationSortRequest struct {
	ApplicationID string `json:"application_id"`
	ResumePath    string `json:"resume_path"`
}

type SortingHatRequest struct {
	JobDescription          string                   `json:"job_description"`
	ApplicationSortRequests []ApplicationSortRequest `json:"application_sort_requests"`
}

type SortingHatResponse struct {
	Scores []SortingHatScore `json:"scores"`
}

type SortingHatScore struct {
	ApplicationID string       `json:"application_id"`
	ModelScores   []ModelScore `json:"model_scores"`
}

type ModelScore struct {
	ModelName string `json:"model_name"`
	Score     int    `json:"score"`
}
