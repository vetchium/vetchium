package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type AddAchievementRequest struct {
	Type        common.AchievementType `json:"type"        validate:"required"`
	Title       string                 `json:"title"       validate:"required,max=128"`
	Description *string                `json:"description" validate:"omitempty,max=1024"`
	URL         *string                `json:"url"         validate:"omitempty,max=1024"`
	At          *time.Time             `json:"at"          validate:"omitempty"`
}

type AddAchievementResponse struct {
	ID string `json:"id"`
}
