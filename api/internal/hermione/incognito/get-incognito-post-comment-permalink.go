package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetIncognitoPostCommentPermalink(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetIncognitoPostCommentPermalink")
		var req hub.GetIncognitoPostCommentPermalinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("validation failed", "req", req)
			return
		}

		// Set default values as per TypeSpec specification
		if req.ContextSiblingsCount == 0 {
			req.ContextSiblingsCount = 3
		}
		if req.ContextRepliesCount == 0 {
			req.ContextRepliesCount = 10
		}

		h.Dbg("Validated", "req", req)

		response, err := h.DB().
			GetIncognitoPostCommentPermalink(r.Context(), req)
		if err != nil {
			h.Dbg("getting comment permalink failed", "error", err)
			switch err {
			case db.ErrNoIncognitoPost:
				h.Dbg("incognito post not found",
					"incognito_post_id", req.IncognitoPostID)
				http.Error(w, "Incognito post not found", http.StatusNotFound)
			case db.ErrNoIncognitoPostComment:
				h.Dbg("comment not found", "comment_id", req.CommentID)
				http.Error(w, "Comment not found", http.StatusNotFound)
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

		h.Dbg("successfully retrieved comment permalink",
			"incognito_post_id", req.IncognitoPostID,
			"comment_id", req.CommentID,
			"context_comments_count", len(response.Comments),
			"target_comment_id", response.TargetCommentID,
			"breadcrumb_depth", len(response.BreadcrumbPath))
	}
}
