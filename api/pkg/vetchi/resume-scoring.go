package vetchi

type ModelScore struct {
	ModelName string `json:"model_name"`
	Score     int    `json:"score"`
}

type SortingHatScore struct {
	ApplicationID string       `json:"application_id"`
	ModelScore    []ModelScore `json:"model_score"`
}

type SortingHatScores struct {
	Scores []SortingHatScore `json:"scores"`
}
