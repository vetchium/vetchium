package comments

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DeletePostComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DeletePostComment")
		var req hub.DeletePostCommentRequest
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

		err := h.DB().DeletePostComment(r.Context(), req)
		if err != nil {
			h.Dbg("deleting post comment failed", "error", err)
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
			"successfully deleted comment",
			"comment_id",
			req.CommentID,
			"post_id",
			req.PostID,
		)
		w.WriteHeader(http.StatusOK)
	}
}
