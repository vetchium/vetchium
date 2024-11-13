package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddOpeningWatchers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddOpeningWatchers")
		var addWatchersReq vetchi.AddOpeningWatchersRequest
		err := json.NewDecoder(r.Body).Decode(&addWatchersReq)
		if err != nil {
			h.Dbg("failed to decode add watchers request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addWatchersReq) {
			h.Dbg("validation failed", "addWatchersReq", addWatchersReq)
			return
		}
		h.Dbg("validated", "addWatchersReq", addWatchersReq)

		err = h.DB().AddOpeningWatchers(r.Context(), addWatchersReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", addWatchersReq.OpeningID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNoOrgUser) {
				h.Dbg("org user not found", "email", addWatchersReq.Emails)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(vetchi.ValidationErrors{
					Errors: []string{"emails"},
				})
				if err != nil {
					h.Err("failed to encode validation errors", "error", err)
				}
				return
			}

			h.Dbg("failed to add watchers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added watchers", "id", addWatchersReq.OpeningID)
		w.WriteHeader(http.StatusOK)
	}
}
