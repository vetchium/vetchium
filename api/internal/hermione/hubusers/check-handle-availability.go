package hubusers

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func CheckHandleAvailability(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.CheckHandleAvailabilityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("request failed validation", "request", req)
			return
		}

		h.Dbg("checkHandleAvailabilityRequest validated", "request", req)

		res, err := h.DB().CheckHandleAvailability(r.Context(), req.Handle)
		if err != nil {
			h.Dbg("failed to check handle availability", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("checkHandleAvailabilityResponse", "response", res)

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
