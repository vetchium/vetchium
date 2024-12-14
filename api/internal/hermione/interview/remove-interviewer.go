package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func RemoveInterviewer(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveInterviewer")
		var removeInterviewerReq employer.RemoveInterviewerRequest
		err := json.NewDecoder(r.Body).Decode(&removeInterviewerReq)
		if err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &removeInterviewerReq) {
			h.Dbg("Validation failed", "req", removeInterviewerReq)
			return
		}
		h.Dbg("validated", "removeInterviewerReq", removeInterviewerReq)

		err = h.DB().RemoveInterviewer(r.Context(), removeInterviewerReq)
		if err != nil {
			if errors.Is(err, db.ErrInvalidInterviewState) {
				h.Dbg("Invalid interview state", "err", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("Error removing interviewer", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
