package hubusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func ChangeEmailAddress(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangeEmailAddress")
		var req hub.ChangeEmailAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "request", req)
			return
		}

		h.Dbg("changeEmailAddressRequest validated", "request", req)

		err := h.DB().ChangeEmailAddress(r.Context(), string(req.Email))
		if err != nil {
			if errors.Is(err, db.ErrDupEmail) {
				h.Dbg("email already in use", "error", err)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("failed to set handle", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
