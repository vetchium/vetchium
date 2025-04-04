package workhistory

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddWorkHistory(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addWorkHistoryReq hub.AddWorkHistoryRequest
		err := json.NewDecoder(r.Body).Decode(&addWorkHistoryReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addWorkHistoryReq) {
			h.Dbg("invalid request", "addWorkHistoryReq", addWorkHistoryReq)
			return
		}
		h.Dbg("validated", "addWorkHistoryReq", addWorkHistoryReq)

		workHistoryID, err := h.DB().
			AddWorkHistory(r.Context(), addWorkHistoryReq)
		if err != nil {
			h.Dbg("failed to add work history", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("work history added", "work_history_id", workHistoryID)

		err = json.NewEncoder(w).Encode(hub.AddWorkHistoryResponse{
			WorkHistoryID: workHistoryID,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
