package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetUserPosts(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetUserPosts")
		var getUserPostsReq hub.GetUserPostsRequest
		if err := json.NewDecoder(r.Body).Decode(&getUserPostsReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getUserPostsReq) {
			h.Dbg("validation failed", "getUserPostsReq", getUserPostsReq)
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		if getUserPostsReq.Limit == 0 {
			getUserPostsReq.Limit = 10
		}

		h.Dbg("Validated", "getUserPostsReq", getUserPostsReq)

		resp, err := h.DB().GetUserPosts(r.Context(), getUserPostsReq)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("invalid handle passed")
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("GetUserPosts failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Dbg("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
