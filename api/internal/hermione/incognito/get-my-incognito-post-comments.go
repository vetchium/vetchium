package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetMyIncognitoPostComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetMyIncognitoPostComments")
		var req hub.GetMyIncognitoPostCommentsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Failed to validate request body")
			return
		}
		h.Dbg("Validated", "getMyIncognitoPostCommentsReq", req)

		resp, err := h.DB().GetMyIncognitoPostComments(r.Context(), req)
		if err != nil {
			h.Dbg("Failed to get my incognito post comments", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Err("Failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Successfully retrieved my incognito post comments",
			"count", len(resp.Comments))
	}
}
