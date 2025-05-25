package comments

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DeleteMyComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DeleteMyComment")
		var req hub.DeleteMyCommentRequest
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

		err := h.DB().DeleteMyComment(r.Context(), req)
		if err != nil {
			h.Dbg("deleting my comment failed", "error", err)
			// Per API spec, no error is returned if comment is not found
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg(
			"successfully deleted my comment",
			"comment_id",
			req.CommentID,
			"post_id",
			req.PostID,
		)
		w.WriteHeader(http.StatusOK)
	}
}
