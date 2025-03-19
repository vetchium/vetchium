package vetchi

// Types here should match the ones in sortinghat/main.py

type ModelScore struct {
	ModelName string `json:"model_name"`
	Score     int    `json:"score"`
}

type SortingHatScore struct {
	ApplicationID string       `json:"application_id"`
	ModelScores   []ModelScore `json:"model_scores"`
}

type SortingHatRequest struct {
	JobDescription string   `json:"job_description"`
	ResumePaths    []string `json:"resume_paths"`
}

type SortingHatResponse struct {
	Scores []SortingHatScore `json:"scores"`
}
