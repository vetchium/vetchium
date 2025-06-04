package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DeleteIncognitoPostComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DeleteIncognitoPostComment")
		var req hub.DeleteIncognitoPostCommentRequest
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

		err := h.DB().DeleteIncognitoPostComment(r.Context(), req)
		if err != nil {
			h.Dbg("deleting incognito post comment failed", "error", err)
			switch err {
			case db.ErrNoIncognitoPost:
				h.Dbg(
					"incognito post not found",
					"incognito_post_id",
					req.IncognitoPostID,
				)
				http.Error(w, "Incognito post not found", http.StatusNotFound)
			case db.ErrNoIncognitoPostComment:
				h.Dbg(
					"incognito post comment not found",
					"comment_id",
					req.CommentID,
				)
				http.Error(w, "Comment not found", http.StatusNotFound)
			case db.ErrNotIncognitoPostCommentAuthor:
				h.Dbg("not comment author", "comment_id", req.CommentID)
				http.Error(
					w,
					"You are not the author of this comment",
					http.StatusForbidden,
				)
			default:
				h.Dbg("internal unhandled error")
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg(
			"successfully deleted incognito post comment",
			"comment_id",
			req.CommentID,
			"incognito_post_id",
			req.IncognitoPostID,
		)
		w.WriteHeader(http.StatusOK)
	}
}
