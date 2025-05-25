package comments

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DisablePostComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DisablePostComments")
		var req hub.DisablePostCommentsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("validation failed", "req", req)
			return
		}

		h.Dbg("Validated", "req", req)

		err := h.DB().DisablePostComments(r.Context(), req)
		if err != nil {
			h.Dbg("disabling post comments failed", "error", err)
			switch err {
			case db.ErrNoPost:
				h.Dbg("post not found", "post_id", req.PostID)
				http.Error(w, "Post not found", http.StatusNotFound)
			case db.ErrNotPostAuthor:
				h.Dbg("not post author", "post_id", req.PostID)
				http.Error(
					w,
					"You are not the author of this post",
					http.StatusForbidden,
				)
			default:
				h.Dbg("internal unhandled error")
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg(
			"successfully disabled comments",
			"post_id",
			req.PostID,
			"delete_existing",
			req.DeleteExistingComments,
		)
		w.WriteHeader(http.StatusOK)
	}
}
