package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DownvoteUserPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DownvoteUserPost")

		var req hub.DownvoteUserPostRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Err("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Err("validation failed", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		err = h.DB().DownvoteUserPost(r.Context(), req)
		if err != nil {
			if errors.Is(err, db.ErrNonVoteableUserPost) {
				h.Dbg("Cannot vote for the post")
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Err("failed to downvote user post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
