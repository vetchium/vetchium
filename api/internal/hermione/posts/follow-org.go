package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func FollowOrg(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FollowOrg")
		var req hub.FollowOrgRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("FollowOrgRequest Validation failed")
			return
		}
		h.Dbg("Validated", "req", req)

		err := h.DB().FollowOrg(r.Context(), req.Domain)
		if err != nil {
			if errors.Is(err, db.ErrNoDomain) {
				h.Dbg("Domain not found", "domain", req.Domain)
				http.Error(w, "Domain not found", http.StatusNotFound)
				return
			}

			h.Dbg("Failed to follow org", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Followed org", "domain", req.Domain)
		w.WriteHeader(http.StatusOK)
	}
}
