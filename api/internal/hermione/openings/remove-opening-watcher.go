package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func RemoveOpeningWatcher(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveOpeningWatcher")
		var removeWatcherReq employer.RemoveOpeningWatcherRequest
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
				h.Dbg("opening not found", "id", removeWatcherReq.OpeningID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to remove watcher", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("removed watcher", "id", removeWatcherReq.OpeningID)
		w.WriteHeader(http.StatusOK)
	}
}
