package education

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func DeleteEducation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteEducationReq hub.DeleteEducationRequest
		err := json.NewDecoder(r.Body).Decode(&deleteEducationReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &deleteEducationReq) {
			h.Dbg("invalid request", "request", deleteEducationReq)
			return
		}

		err = h.DB().DeleteEducation(r.Context(), deleteEducationReq)
		if err != nil {
			if errors.Is(err, db.ErrNoEducation) {
				h.Dbg("education not found")
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to delete education", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
