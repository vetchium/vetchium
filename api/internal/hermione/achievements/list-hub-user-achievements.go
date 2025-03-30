package achievements

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func ListHubUserAchievements(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listHubUserAchievementsReq employer.ListHubUserAchievementsRequest
		if err := json.NewDecoder(r.Body).Decode(&listHubUserAchievementsReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &listHubUserAchievementsReq) {
			h.Dbg("invalid request", "req", listHubUserAchievementsReq)
			return
		}

		h.Dbg("validated", "req", listHubUserAchievementsReq)

		achievements, err := h.DB().
			ListHubUserAchievements(r.Context(), listHubUserAchievementsReq)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				<-time.After(time.Duration(rand.Intn(2)) * time.Second)
				h.Dbg("not found", "error", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			h.Dbg("failed to list hub user achievements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("listed hub user achievements", "length", len(achievements))
		err = json.NewEncoder(w).Encode(achievements)
		if err != nil {
			h.Dbg("failed to encode achievements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
