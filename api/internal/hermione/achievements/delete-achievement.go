package achievements

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func DeleteAchievement(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteAchievementReq hub.DeleteAchievementRequest
		if err := json.NewDecoder(r.Body).Decode(&deleteAchievementReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &deleteAchievementReq) {
			h.Dbg("invalid request", "req", deleteAchievementReq)
			return
		}

		h.Dbg("validated", "deleteAchievementReq", deleteAchievementReq)

		err := h.DB().DeleteAchievement(r.Context(), deleteAchievementReq)
		if err != nil {
			if errors.Is(err, db.ErrNoAchievement) {
				<-time.After(time.Duration(rand.Intn(2)) * time.Second)
				h.Dbg("no achievement found", "error", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			h.Dbg("failed to delete achievement", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("deleted achievement", "achievementID", deleteAchievementReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
