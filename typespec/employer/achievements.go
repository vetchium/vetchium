package employer

import (
	"github.com/psankar/vetchi/typespec/common"
)

type ListHubUserAchievementsRequest struct {
	Handle common.Handle `json:"handle"`
}
