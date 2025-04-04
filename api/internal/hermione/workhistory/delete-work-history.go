package workhistory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func DeleteWorkHistory(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteWorkHistoryReq hub.DeleteWorkHistoryRequest
		err := json.NewDecoder(r.Body).Decode(&deleteWorkHistoryReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &deleteWorkHistoryReq) {
			h.Dbg("invalid", "deleteWorkHistoryReq", deleteWorkHistoryReq)
			return
		}
		h.Dbg("validated", "deleteWorkHistoryReq", deleteWorkHistoryReq)

		err = h.DB().DeleteWorkHistory(r.Context(), deleteWorkHistoryReq)
		if err != nil {
			if errors.Is(err, db.ErrNoWorkHistory) {
				h.Dbg("work history not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to delete work history", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
