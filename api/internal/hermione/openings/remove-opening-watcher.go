package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func RemoveOpeningWatcher(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveOpeningWatcher")
		var removeWatcherReq vetchi.RemoveOpeningWatcherRequest
		err := json.NewDecoder(r.Body).Decode(&removeWatcherReq)
		if err != nil {
			h.Dbg("failed to decode remove watcher request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &removeWatcherReq) {
			h.Dbg("validation failed", "removeWatcherReq", removeWatcherReq)
			return
		}
		h.Dbg("validated", "removeWatcherReq", removeWatcherReq)

		err = h.DB().RemoveOpeningWatcher(r.Context(), removeWatcherReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", removeWatcherReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to remove watcher", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("removed watcher", "id", removeWatcherReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
