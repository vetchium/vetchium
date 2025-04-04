package employer

import (
	"github.com/vetchium/vetchium/typespec/common"
)

type ListHubUserAchievementsRequest struct {
	Handle common.Handle          `json:"handle" validate:"required,validate_handle"`
	Type   common.AchievementType `json:"type"   validate:"omitempty,validate_achievement_type"`
}
