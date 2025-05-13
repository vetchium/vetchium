package empposts

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func ListEmployerPosts(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ListEmployerPosts")
		var req employer.ListEmployerPostsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			return
		}
		h.Dbg("validated", "listEmployerPostsRequest", req)

		resp, err := h.DB().
			ListEmployerPosts(db.ListEmployerPostsRequest{
				Context:                  r.Context(),
				ListEmployerPostsRequest: req,
			})
		if err != nil {
			h.Dbg("failed to list employer posts", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("employer posts fetched", "len", len(resp.Posts))

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Err("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
