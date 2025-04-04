package colleagues

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func ConnectColleague(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ConnectColleague")

		var req hub.ConnectColleagueRequest
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

		err := h.DB().ConnectColleague(r.Context(), string(req.Handle))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("Colleague not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNotColleaguable) {
				h.Dbg("handle is not colleaguable for the logged in user now")
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("Error connecting colleague", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
