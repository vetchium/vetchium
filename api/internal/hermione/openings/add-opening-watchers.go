package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddOpeningWatchers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddOpeningWatchers")
		var addWatchersReq vetchi.AddWatchersRequest
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		emails := make([]string, len(addWatchersReq.Emails))
		for i, email := range addWatchersReq.Emails {
			emails[i] = string(email)
		}

		err = h.DB().AddOpeningWatchers(r.Context(), db.AddOpeningWatchersReq{
			ID:         addWatchersReq.ID,
			Emails:     emails,
			EmployerID: orgUser.EmployerID,
			AddedBy:    orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", addWatchersReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to add watchers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added watchers", "id", addWatchersReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
