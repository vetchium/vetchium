package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
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
			if errors.Is(err, db.ErrNoInterview) {
				h.Dbg("no interview found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNotAnInterviewer) {
				h.Dbg("not an interviewer", "error", err)
				http.Error(w, "", http.StatusForbidden)
				return
			}

			if errors.Is(err, db.ErrStateMismatch) {
				h.Dbg("interview state mismatch", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("error putting assessment", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("assessment put successfully", "assessment", assessment)
		w.WriteHeader(http.StatusOK)
	}
}
