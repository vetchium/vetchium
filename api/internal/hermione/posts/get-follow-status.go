package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetFollowStatus(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetFollowStatus")
		var req hub.GetFollowStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		h.Dbg("Validated", "req", req)

		status, err := h.DB().GetFollowStatus(r.Context(), string(req.Handle))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("Handle not found", "handle", req.Handle)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Follow status", "status", status)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			h.Err("Failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
