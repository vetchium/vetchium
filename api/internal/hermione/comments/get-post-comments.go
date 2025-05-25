package comments

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetPostComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetPostComments")
		var req hub.GetPostCommentsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("validation failed", "req", req)
			return
		}

		// Set default limit if not provided
		if req.Limit == 0 {
			req.Limit = 10
		}

		h.Dbg("Validated", "req", req)

		comments, err := h.DB().GetPostComments(r.Context(), req)
		if err != nil {
			h.Dbg("getting post comments failed", "error", err)
			switch err {
			case db.ErrNoPost:
				h.Dbg("post not found", "post_id", req.PostID)
				http.Error(w, "Post not found", http.StatusNotFound)
			default:
				h.Dbg("internal unhandled error")
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(comments)
		if err != nil {
			h.Dbg("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg(
			"successfully retrieved comments",
			"post_id",
			req.PostID,
			"count",
			len(comments),
		)
	}
}
