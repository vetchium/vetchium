package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func EmployerPutAssessment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerPutAssessment")

		var assessment employer.Assessment
		if err := json.NewDecoder(r.Body).Decode(&assessment); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &assessment) {
			h.Dbg("validation failed", "assessment", assessment)
			return
		}

		if err := h.DB().PutAssessment(r.Context(), assessment); err != nil {
			h.Dbg("error putting assessment", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.Dbg("assessment put successfully", "assessment", assessment)
		w.WriteHeader(http.StatusOK)
	}
}
