package posts

import (
	"encoding/json"
	"net/http"

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
			h.Dbg("Failed to follow org", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Followed org", "domain", req.Domain)
		w.WriteHeader(http.StatusOK)
	}
}
