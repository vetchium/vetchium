package incognito

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DownvoteIncognitoPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DownvoteIncognitoPost")

		var req hub.DownvoteIncognitoPostRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("validation failed", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("Validated", "DownvoteIncognitoPostRequest", req)

		err = h.DB().DownvoteIncognitoPost(r.Context(), req)
		if err != nil {
			if errors.Is(err, db.ErrNoIncognitoPost) {
				h.Dbg("Incognito post not found")
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNonVoteableIncognitoPost) {
				h.Dbg("Cannot vote for the incognito post")
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to downvote incognito post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Downvoted incognito post", "DownvoteIncognitoPostRequest", req)
		w.WriteHeader(http.StatusOK)
	}
}
