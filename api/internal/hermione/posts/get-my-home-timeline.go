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
			h.Dbg("Failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Validation failed", "req", req)
			return
		}

		if req.Limit == 0 {
			h.Dbg("Limit is 0, setting to default of 25")
			req.Limit = 25
		}

		h.Dbg("Validated", "req", req)

		getMyHomeTimelineResp, err := h.DB().GetMyHomeTimeline(r.Context(), req)
		if err != nil {
			// Check for invalid pagination key
			if errors.Is(err, db.ErrInvalidPaginationKey) {
				h.Dbg("Invalid pagination key", "error", err)
				w.WriteHeader(http.StatusUnprocessableEntity) // 422
				return
			}

			h.Dbg("Failed to get home timeline", "error", err)
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
