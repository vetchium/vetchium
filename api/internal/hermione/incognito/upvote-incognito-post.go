package incognito

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func UpvoteIncognitoPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpvoteIncognitoPost")

		var req hub.UpvoteIncognitoPostRequest
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
		h.Dbg("Validated", "UpvoteIncognitoPostRequest", req)

		err = h.DB().UpvoteIncognitoPost(r.Context(), req)
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

			if errors.Is(err, db.ErrIncognitoPostVoteConflict) {
				h.Dbg(
					"Vote conflict: user has already voted in opposite direction",
				)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to upvote incognito post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Upvoted incognito post", "UpvoteIncognitoPostRequest", req)
		w.WriteHeader(http.StatusOK)
	}
}
