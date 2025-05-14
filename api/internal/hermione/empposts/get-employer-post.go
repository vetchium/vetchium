package empposts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func GetEmployerPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetEmployerPost")
		var req employer.GetEmployerPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			return
		}
		h.Dbg("validated", "getEmployerPostRequest", req)

		post, err := h.DB().GetEmployerPost(r.Context(), req.PostID)
		if err != nil {
			if errors.Is(err, db.ErrNoEmployerPost) {
				h.Dbg("Post not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get employer post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("employer post fetched", "post_id", req.PostID)

		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			h.Err("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
