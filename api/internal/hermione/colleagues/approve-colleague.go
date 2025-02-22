package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func ApproveColleague(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ApproveColleague")

		var req hub.ApproveColleagueRequest
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

		if err := h.DB().ApproveColleague(r.Context(), string(req.Handle)); err != nil {
			h.Dbg("failed to approve colleague", "error", err)
			switch err {
			case db.ErrNoHubUser:
				h.Dbg("no hub user found", "handle", req.Handle)
				http.Error(w, "", http.StatusNotFound)
			case db.ErrNoApplication:
				h.Dbg("no pending request found", "handle", req.Handle)
				http.Error(w, "", http.StatusNotFound)
			default:
				h.Dbg("internal server error", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
