package comments

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddPostComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddPostComment")
		var req hub.AddPostCommentRequest
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

		// Generate comment ID
		commentID := util.RandomUniqueID(vetchi.CommentIDLenBytes)

		// Create db request with generated comment ID
		dbReq := db.AddPostCommentRequest{
			CommentID:             commentID,
			AddPostCommentRequest: req,
		}

		response, err := h.DB().AddPostComment(r.Context(), dbReq)
		if err != nil {
			h.Dbg("adding post comment failed", "error", err)
			switch err {
			case db.ErrNoPost:
				h.Dbg("post not found", "post_id", req.PostID)
				http.Error(w, "Post not found", http.StatusNotFound)
			case db.ErrCommentsDisabled:
				h.Dbg("comments disabled", "post_id", req.PostID)
				http.Error(
					w,
					"Comments are disabled for this post",
					http.StatusForbidden,
				)
			default:
				h.Dbg("internal unhandled error")
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			h.Dbg("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("successfully added comment",
			"comment_id", response.CommentID,
			"post_id", response.PostID,
		)
	}
}
