package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func VerifyOfficialEmail(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.VerifyOfficialEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "req", req)
			return
		}
		h.Dbg("validated", "req", req)

		err = h.DB().
			VerifyOfficialEmail(r.Context(), string(req.Email), req.Code)
		if err != nil {
			if err == db.ErrOfficialEmailNotFound {
				h.Dbg("email not found", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
			if err == db.ErrInvalidVerificationCode {
				h.Dbg("invalid verification code", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
			h.Dbg("failed to verify official email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("verified official email", "email", req.Email)
		w.WriteHeader(http.StatusOK)
	}
}
