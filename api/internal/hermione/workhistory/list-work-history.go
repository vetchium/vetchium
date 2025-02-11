package workhistory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func ListWorkHistory(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listWorkHistoryReq hub.ListWorkHistoryRequest
		err := json.NewDecoder(r.Body).Decode(&listWorkHistoryReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &listWorkHistoryReq) {
			h.Dbg("invalid request", "error", err)
			return
		}
		h.Dbg("validated", "listWorkHistoryReq", listWorkHistoryReq)

		workHistory, err := h.DB().
			ListWorkHistory(r.Context(), listWorkHistoryReq)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("failed to list work history", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Err("failed to list work history", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("work history", "work_history", workHistory)

		err = json.NewEncoder(w).Encode(workHistory)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
