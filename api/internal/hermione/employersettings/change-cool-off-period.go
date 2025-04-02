package employersettings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func ChangeCoolOffPeriod(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangeCoolOffPeriod")
		var req employer.ChangeCoolOffPeriodRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("validation failed", "req", req)
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		err := h.DB().ChangeCoolOffPeriod(r.Context(), req.CoolOffPeriod)
		if err != nil {
			h.Dbg("failed to change cool off period", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
