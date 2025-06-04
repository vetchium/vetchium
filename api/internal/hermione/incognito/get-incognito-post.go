package incognito

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetIncognitoPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetIncognitoPost")
		var req hub.GetIncognitoPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Failed to validate request body")
			return
		}
		h.Dbg("Validated", "getIncognitoPostReq", req)

		post, err := h.DB().GetIncognitoPost(r.Context(), req)
		if err != nil {
			if errors.Is(err, db.ErrNoIncognitoPost) {
				h.Dbg("Incognito post not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("Failed to get incognito post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			h.Err("Failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
