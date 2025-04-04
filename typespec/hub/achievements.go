package hub

import (
	"time"

	"github.com/vetchium/vetchium/typespec/common"
)

type AddAchievementRequest struct {
	Type        common.AchievementType `json:"type"        validate:"required"`
	Title       string                 `json:"title"       validate:"required,min=3,max=128"`
	Description *string                `json:"description" validate:"omitempty,min=3,max=1024"`
	URL         *string                `json:"url"         validate:"omitempty,min=3,max=1024"`
	At          *time.Time             `json:"at"          validate:"omitempty,ltefield=now"`
}

type AddAchievementResponse struct {
	ID string `json:"id"`
}

type ListAchievementsRequest struct {
	Type   common.AchievementType `json:"type"   validate:"omitempty,validate_achievement_type"`
	Handle *common.Handle         `json:"handle" validate:"omitempty,validate_handle"`
}

type DeleteAchievementRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
