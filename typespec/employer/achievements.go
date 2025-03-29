package employer

import (
	"github.com/psankar/vetchi/typespec/common"
)

type ListAchievementsRequest struct {
	Handle common.Handle `json:"handle"`
}
