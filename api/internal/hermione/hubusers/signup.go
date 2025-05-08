package hubusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func SignupHubUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered SignupHubUser")
		var req hub.SignupHubUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "request", req)
			return
		}

		h.Dbg("signupHubUserRequest validated", "request", req)

		err := h.DB().SignupHubUser(r.Context(), string(req.Email))
		if err != nil {
			if errors.Is(err, db.ErrUnsupportedDomain) {
				h.Dbg("email domain is not supported for signup")
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to set handle", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
