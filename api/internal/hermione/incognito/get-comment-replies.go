package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetCommentReplies(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCommentReplies")
		var req hub.GetCommentRepliesRequest
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
		if req.Limit == 0 {
			req.Limit = 50
		}
		if req.MaxDepth == 0 {
			req.MaxDepth = 2
		}

		h.Dbg("Validated", "req", req)

		response, err := h.DB().GetCommentReplies(r.Context(), req)
		if err != nil {
			h.Dbg("getting comment replies failed", "error", err)
			switch err {
			case db.ErrNoIncognitoPost:
				h.Dbg("incognito post not found",
					"incognito_post_id", req.IncognitoPostID)
				http.Error(w, "Incognito post not found", http.StatusNotFound)
			case db.ErrNoIncognitoPostComment:
				h.Dbg("parent comment not found",
					"parent_comment_id", req.ParentCommentID)
				http.Error(w, "Parent comment not found", http.StatusNotFound)
			case db.ErrMaxCommentDepthReached:
				h.Dbg("max comment depth would be exceeded",
					"parent_comment_id", req.ParentCommentID)
				http.Error(w,
					"Maximum comment depth limit would be exceeded",
					http.StatusBadRequest)
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

		h.Dbg("successfully retrieved comment replies",
			"incognito_post_id", req.IncognitoPostID,
			"parent_comment_id", req.ParentCommentID,
			"count", len(response.Replies),
			"total_replies_count", response.TotalRepliesCount,
			"pagination_key", response.PaginationKey)
	}
}
