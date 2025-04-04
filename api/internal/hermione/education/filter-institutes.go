package education

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func FilterInstitutes(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filterInstitutesReq hub.FilterInstitutesRequest
		err := json.NewDecoder(r.Body).Decode(&filterInstitutesReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterInstitutesReq) {
			h.Dbg("invalid request", "request", filterInstitutesReq)
			return
		}

		institutes, err := h.DB().
			FilterInstitutes(r.Context(), filterInstitutesReq)
		if err != nil {
			h.Dbg("failed to filter institutes", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(institutes)
		if err != nil {
			h.Dbg("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
