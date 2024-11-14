package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func GetOpeningWatchers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetOpeningWatchers")
		var getOpeningWatchersReq vetchi.GetOpeningWatchersRequest
		err := json.NewDecoder(r.Body).Decode(&getOpeningWatchersReq)
		if err != nil {
			h.Dbg("failed to decode get opening watchers request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getOpeningWatchersReq) {
			h.Dbg("validation failed", "req", getOpeningWatchersReq)
			return
		}
		h.Dbg("validated", "getOpeningWatchersReq", getOpeningWatchersReq)

		watchers, err := h.DB().
			GetOpeningWatchers(r.Context(), getOpeningWatchersReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("not found", "openingID", getOpeningWatchersReq.OpeningID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get opening watchers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(watchers)
		if err != nil {
			h.Err("failed to encode watchers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
