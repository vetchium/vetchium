package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func UpdateLocation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateLocation")
		var updateLocationReq employer.UpdateLocationRequest
		err := json.NewDecoder(r.Body).Decode(&updateLocationReq)
		if err != nil {
			h.Dbg("failed to decode update location request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateLocationReq) {
			h.Dbg("validation failed", "updateLocationReq", updateLocationReq)
			return
		}
		h.Dbg("validated", "updateLocationReq", updateLocationReq)

		err = h.DB().UpdateLocation(r.Context(), updateLocationReq)
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				h.Dbg("not found", "title", updateLocationReq.Title)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("updated location", "updateLocationReq", updateLocationReq)
		w.WriteHeader(http.StatusOK)
	}
}
