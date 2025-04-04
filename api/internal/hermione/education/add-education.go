package education

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddEducation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addEducationReq hub.AddEducationRequest
		err := json.NewDecoder(r.Body).Decode(&addEducationReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addEducationReq) {
			h.Dbg("invalid request", "addEducationReq", addEducationReq)
			return
		}

		h.Dbg("validated", "addEducationReq", addEducationReq)

		educationID, err := h.DB().AddEducation(r.Context(), addEducationReq)
		if err != nil {
			h.Dbg("failed to add education", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("education added", "educationID", educationID)

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(hub.AddEducationResponse{
			EducationID: educationID,
		})
		if err != nil {
			h.Dbg("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
