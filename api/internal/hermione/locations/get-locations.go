package locations

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetLocations(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetLocations")
		var getLocationsReq employer.GetLocationsRequest
		err := json.NewDecoder(r.Body).Decode(&getLocationsReq)
		if err != nil {
			h.Dbg("failed to decode getLocationsReq", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getLocationsReq) {
			h.Dbg("failed to validate getLocationsReq", "error", err)
			return
		}
		h.Dbg("Validated", "getLocationsReq", getLocationsReq)

		if getLocationsReq.Limit == 0 {
			getLocationsReq.Limit = 100
			h.Dbg("set default limit", "limit", getLocationsReq.Limit)
		}

		locations, err := h.DB().GetLocations(r.Context(), getLocationsReq)
		if err != nil {
			h.Dbg("failed to get locations", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Got locations", "locations", locations)
		err = json.NewEncoder(w).Encode(locations)
		if err != nil {
			h.Err("failed to encode locations", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
