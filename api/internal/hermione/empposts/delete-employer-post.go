package empposts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func DeleteEmployerPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DeleteEmployerPost")
		var req employer.DeleteEmployerPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			return
		}
		h.Dbg("validated", "deleteEmployerPostRequest", req)

		err := h.DB().DeleteEmployerPost(r.Context(), req.PostID)
		if err != nil {
			if errors.Is(err, db.ErrNoPost) {
				h.Dbg("Post not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to delete employer post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("employer post deleted", "post_id", req.PostID)
	}
}
