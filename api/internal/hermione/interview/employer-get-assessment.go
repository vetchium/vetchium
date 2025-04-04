package interview

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func EmployerGetAssessment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerGetAssessment")
		var getAssessmentReq employer.GetAssessmentRequest
		if err := json.NewDecoder(r.Body).Decode(&getAssessmentReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getAssessmentReq) {
			h.Dbg("validation failed", "getAssessmentReq", getAssessmentReq)
			return
		}
		h.Dbg("validated", "getAssessmentReq", getAssessmentReq)

		assessment, err := h.DB().
			GetAssessment(r.Context(), getAssessmentReq)
		if err != nil {
			h.Dbg("error getting assessment", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.Dbg("got assessment", "assessment", assessment)

		if err := json.NewEncoder(w).Encode(assessment); err != nil {
			h.Err("error encoding assessment", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
