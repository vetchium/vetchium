package achievements

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func AddAchievement(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addAchievementReq hub.AddAchievementRequest
		if err := json.NewDecoder(r.Body).Decode(&addAchievementReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addAchievementReq) {
			h.Dbg("invalid request", "req", addAchievementReq)
			return
		}

		h.Dbg("validated", "addAchievementReq", addAchievementReq)

		achievementID, err := h.DB().
			AddAchievement(r.Context(), addAchievementReq)
		if err != nil {
			h.Dbg("failed to add achievement", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("achievement added", "achievementID", achievementID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(hub.AddAchievementResponse{ID: achievementID})
	}
}
