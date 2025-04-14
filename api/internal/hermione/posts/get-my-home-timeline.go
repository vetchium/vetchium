package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetMyHomeTimeline(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetMyHomeTimeline")
		var req hub.GetMyHomeTimelineRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		h.Dbg("Validated", "req", req)

		getMyHomeTimelineResp, err := h.DB().GetMyHomeTimeline(r.Context(), req)
		if err != nil {
			// Check for invalid pagination key
			if errors.Is(err, db.ErrInvalidPaginationKey) {
				h.Err("Invalid pagination key", "error", err)
				w.WriteHeader(http.StatusUnprocessableEntity) // 422
				return
			}

			h.Err("Failed to get home timeline", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("GetMyHomeTimelineResponse", "resp", getMyHomeTimelineResp)
		if err := json.NewEncoder(w).Encode(getMyHomeTimelineResp); err != nil {
			h.Err("encode failure", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
