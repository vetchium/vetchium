package employer

import (
	"github.com/psankar/vetchi/typespec/common"
)

type ListHubUserAchievementsRequest struct {
	Handle common.Handle          `json:"handle" validate:"required,validatate_handle"`
	Type   common.AchievementType `json:"type"   validate:"omitempty,validatate_achievement_type"`
}
