package applications

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func WithdrawApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered WithdrawApplication")
		var withdrawApplicationReq hub.WithdrawApplicationRequest
		if err := json.NewDecoder(r.Body).Decode(&withdrawApplicationReq); err != nil {
			h.Err("failed to decode withdraw application request", "error", err)
			h.Dbg("exiting WithdrawApplication - decode error")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &withdrawApplicationReq) {
			h.Dbg("invalid request", "request", withdrawApplicationReq)
			h.Dbg("exiting WithdrawApplication - validation error")
			return
		}
		h.Dbg("validated request", "request", withdrawApplicationReq)

		err := h.DB().WithdrawApplication(
			r.Context(),
			withdrawApplicationReq.ApplicationID,
		)
		if err != nil {
			switch err {
			case db.ErrNoApplication:
				h.Dbg("exiting WithdrawApplication - application not found")
				http.Error(w, err.Error(), http.StatusNotFound)
			case db.ErrApplicationStateInCompatible:
				h.Dbg("exiting WithdrawApplication - incompatible state")
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			default:
				h.Dbg("exiting WithdrawApplication - internal error")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("withdrew application", "id", withdrawApplicationReq)
		w.WriteHeader(http.StatusOK)
	}
}
