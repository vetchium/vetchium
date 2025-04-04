package achievements

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func ListAchievements(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listAchievementsReq hub.ListAchievementsRequest
		if err := json.NewDecoder(r.Body).Decode(&listAchievementsReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &listAchievementsReq) {
			h.Dbg("invalid request", "req", listAchievementsReq)
			return
		}

		h.Dbg("validated", "listAchievementsReq", listAchievementsReq)

		achievements, err := h.DB().
			ListAchievements(r.Context(), listAchievementsReq)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				<-time.After(time.Duration(rand.Intn(2)) * time.Second)
				h.Dbg("no hub user found", "error", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			h.Dbg("failed to list achievements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got achievements", "count", len(achievements))
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(achievements)
		if err != nil {
			h.Dbg("failed to encode achievements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
