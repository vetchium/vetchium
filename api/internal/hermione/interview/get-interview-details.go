package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func GetInterviewDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetInterviewDetails")
		var getInterviewDetailsReq employer.GetInterviewDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&getInterviewDetailsReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getInterviewDetailsReq) {
			h.Dbg("validation failed")
			return
		}
		h.Dbg("validated", "request", getInterviewDetailsReq)

		interview, err := h.DB().
			GetInterview(r.Context(), getInterviewDetailsReq.InterviewID)
		if err != nil {
			if errors.Is(err, db.ErrNoInterview) {
				h.Dbg("interview not found")
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("error getting interview", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("interview found", "interview", interview)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(interview)
		if err != nil {
			h.Dbg("error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
