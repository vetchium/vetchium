package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func FollowUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FollowUser")
		var req hub.FollowUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("FollowUserRequest Validation failed")
			return
		}
		h.Dbg("Validated", "req", req)

		err := h.DB().FollowUser(r.Context(), string(req.Handle))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("not found or not active", "handle", req.Handle)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("Failed to follow user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Followed user", "handle", req.Handle)
		w.WriteHeader(http.StatusOK)
	}
}
