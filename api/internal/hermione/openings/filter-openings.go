package openings

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func FilterOpenings(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterOpenings")
		var filterOpeningsReq employer.FilterOpeningsRequest
		err := json.NewDecoder(r.Body).Decode(&filterOpeningsReq)
		if err != nil {
			h.Dbg("failed to decode filter openings request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterOpeningsReq) {
			h.Dbg("validation failed", "filterOpeningsReq", filterOpeningsReq)
			return
		}
		h.Dbg("validated", "filterOpeningsReq", filterOpeningsReq)

		defaultFromDate := time.Now().AddDate(0, 0, -30)
		if filterOpeningsReq.FromDate == nil {
			h.Dbg("setting default fromdate", "date", defaultFromDate)
			filterOpeningsReq.FromDate = &defaultFromDate
		}

		defaultToDate := time.Now().AddDate(0, 0, 1)
		if filterOpeningsReq.ToDate == nil {
			h.Dbg("setting default todate", "date", defaultToDate)
			filterOpeningsReq.ToDate = &defaultToDate
		}

		if filterOpeningsReq.FromDate.After(*filterOpeningsReq.ToDate) {
			h.Dbg("fromdate > todate", "filterOpeningsReq", filterOpeningsReq)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).
				Encode(common.ValidationErrors{Errors: []string{"fromdate"}})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		if filterOpeningsReq.Limit <= 0 {
			filterOpeningsReq.Limit = 40
			h.Dbg("set default limit", "limit", filterOpeningsReq.Limit)
		}

		if len(filterOpeningsReq.State) == 0 {
			filterOpeningsReq.State = []common.OpeningState{
				common.ActiveOpening,
				common.DraftOpening,
				common.SuspendedOpening,
			}
			h.Dbg("set default state", "state", filterOpeningsReq.State)
		}

		openingInfos, err := h.DB().
			FilterOpenings(r.Context(), filterOpeningsReq)
		if err != nil {
			h.Dbg("failed to filter openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Filtered Openings", "openingInfos", openingInfos)

		err = json.NewEncoder(w).Encode(openingInfos)
		if err != nil {
			h.Err("failed to encode openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
