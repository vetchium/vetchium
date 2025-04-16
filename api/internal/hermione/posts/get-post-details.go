package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetPostDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetPostDetails")
		var req hub.GetPostDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Failed to validate request body")
			return
		}
		h.Dbg("Validated", "getPostDetailsReq", req)

		post, err := h.DB().GetPost(db.GetPostRequest{
			Context: r.Context(),
			PostID:  req.PostID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoPost) {
				h.Dbg("Post not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("Failed to get post", "error", err)
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
