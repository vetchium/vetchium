package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func UnlinkColleague(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UnlinkColleague")

		var req hub.UnlinkColleagueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated request", "req", req)

		if err := h.DB().UnlinkColleague(r.Context(), string(req.Handle)); err != nil {
			h.Dbg("failed to unlink colleague", "error", err)
			switch err {
			case db.ErrNoHubUser:
				h.Dbg("no hub user found", "handle", req.Handle)
				http.Error(w, "", http.StatusNotFound)
			case db.ErrNoConnection:
				h.Dbg(
					"no active colleague connection found",
					"handle",
					req.Handle,
				)
				http.Error(w, "", http.StatusNotFound)
			default:
				h.Dbg("internal server error", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("colleague unlinked", "handle", req.Handle)
		w.WriteHeader(http.StatusOK)
	}
}
