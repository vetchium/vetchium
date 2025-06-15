package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetIncognitoPostComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetIncognitoPostComments")
		var req hub.GetIncognitoPostCommentsRequest
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
			req.Limit = 25
		}
		if req.SortBy == "" {
			req.SortBy = hub.IncognitoPostCommentSortByTop
		}
		if req.DirectRepliesPerComment == 0 {
			req.DirectRepliesPerComment = 3
		}

		h.Dbg("Validated", "req", req)

		response, err := h.DB().GetIncognitoPostComments(r.Context(), req)
		if err != nil {
			h.Dbg("getting incognito post comments failed", "error", err)
			switch err {
			case db.ErrNoIncognitoPost:
				h.Dbg("incognito post not found",
					"incognito_post_id", req.IncognitoPostID)
				http.Error(w, "Incognito post not found", http.StatusNotFound)
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

		h.Dbg("successfully retrieved incognito post comments",
			"incognito_post_id", req.IncognitoPostID,
			"count", len(response.Comments),
			"total_top_level_count", response.TotalCommentsCount,
			"pagination_key", response.PaginationKey)
	}
}
