package hubusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func SetHandle(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.SetHandleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("request failed validation", "request", req)
			return
		}

		h.Dbg("setHandleRequest validated", "request", req)

		err := h.DB().SetHandle(r.Context(), req.Handle)
		if err != nil {
			if errors.Is(err, db.ErrUnpaidHubUser) {
				h.Dbg("user is not a paid hub user", "error", err)
				http.Error(w, "", http.StatusForbidden)
				return
			}

			if errors.Is(err, db.ErrDupHandle) {
				h.Dbg("handle already in use", "error", err)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("failed to set handle", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("handle set", "handle", req.Handle)
		w.WriteHeader(http.StatusOK)
	}
}
