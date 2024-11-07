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
			h.Dbg(
				"validation failed",
				"getOpeningWatchersReq",
				getOpeningWatchersReq,
			)
			return
		}
		h.Dbg("validated", "getOpeningWatchersReq", getOpeningWatchersReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		watchers, err := h.DB().
			GetOpeningWatchers(r.Context(), db.GetOpeningWatchersReq{
				ID:         getOpeningWatchersReq.ID,
				EmployerID: orgUser.EmployerID,
			})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", getOpeningWatchersReq.ID)
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
