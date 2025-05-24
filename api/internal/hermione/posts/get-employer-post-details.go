package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetEmployerPostDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetEmployerPostDetails")
		var req hub.GetEmployerPostDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Failed to validate request body")
			return
		}
		h.Dbg("Validated", "getEmployerPostDetailsReq", req)

		post, err := h.DB().
			GetEmployerPostForHub(r.Context(), req.EmployerPostID)
		if err != nil {
			if errors.Is(err, db.ErrNoEmployerPost) {
				h.Dbg("Employer post not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("Failed to get employer post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			h.Err("Failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
