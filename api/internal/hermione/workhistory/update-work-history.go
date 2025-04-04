package workhistory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func UpdateWorkHistory(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateWorkHistoryReq hub.UpdateWorkHistoryRequest
		err := json.NewDecoder(r.Body).Decode(&updateWorkHistoryReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateWorkHistoryReq) {
			h.Dbg("invalid", "updateWorkHistoryReq", updateWorkHistoryReq)
			return
		}
		h.Dbg("validated", "updateWorkHistoryReq", updateWorkHistoryReq)

		err = h.DB().UpdateWorkHistory(r.Context(), updateWorkHistoryReq)
		if err != nil {
			if errors.Is(err, db.ErrNoWorkHistory) {
				h.Dbg("work history not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update work history", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
