package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddIncognitoPostComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddIncognitoPostComment")
		var req hub.AddIncognitoPostCommentRequest
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

		commentID := util.RandomUniqueID(vetchi.CommentIDLenBytes)

		dbReq := db.AddIncognitoPostCommentRequest{
			Context:                        r.Context(),
			CommentID:                      commentID,
			AddIncognitoPostCommentRequest: req,
		}

		response, err := h.DB().AddIncognitoPostComment(r.Context(), dbReq)
		if err != nil {
			h.Dbg("adding incognito post comment failed", "error", err)
			switch err {
			case db.ErrNoIncognitoPost:
				h.Dbg(
					"incognito post not found",
					"incognito_post_id",
					req.IncognitoPostID,
				)
				http.Error(w, "Incognito post not found", http.StatusNotFound)
			case db.ErrInvalidParentComment:
				h.Dbg("invalid parent comment", "in_reply_to", req.InReplyTo)
				http.Error(
					w,
					"Parent comment not found or has been deleted",
					http.StatusNotFound,
				)
			case db.ErrMaxCommentDepthReached:
				h.Dbg("max comment depth reached", "in_reply_to", req.InReplyTo)
				http.Error(
					w,
					"Maximum comment depth reached. Cannot reply to this comment.",
					http.StatusUnprocessableEntity,
				)
			case db.ErrNoIncognitoPostComment:
				h.Dbg("parent comment not found", "in_reply_to", req.InReplyTo)
				http.Error(w, "Parent comment not found", http.StatusNotFound)
			default:
				h.Err("internal unhandled error", "error", err)
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

		h.Dbg("successfully added incognito post comment",
			"comment_id", response.CommentID,
			"incognito_post_id", response.IncognitoPostID)
	}
}
