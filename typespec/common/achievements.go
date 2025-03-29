package common

import "time"

type AchievementType string

const (
	Patent        AchievementType = "PATENT"
	Publication   AchievementType = "PUBLICATION"
	Certification AchievementType = "CERTIFICATION"
)

type Achievement struct {
	ID          string          `json:"id"`
	Type        AchievementType `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	At          time.Time       `json:"at"`
}
