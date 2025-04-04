package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func RenameLocation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RenameLocation")
		var renameLocationReq employer.RenameLocationRequest
		err := json.NewDecoder(r.Body).Decode(&renameLocationReq)
		if err != nil {
			h.Dbg("failed to decode rename location request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &renameLocationReq) {
			h.Dbg("validation failed", "renameLocationReq", renameLocationReq)
			return
		}
		h.Dbg("validated", "renameLocationReq", renameLocationReq)

		err = h.DB().RenameLocation(r.Context(), renameLocationReq)
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				h.Dbg("not found", "title", renameLocationReq.OldTitle)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrDupLocationName) {
				h.Dbg("location exists", "renameLocationReq", renameLocationReq)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("failed to rename location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("renamed location", "renameLocReq", renameLocationReq)
		w.WriteHeader(http.StatusOK)
	}
}
