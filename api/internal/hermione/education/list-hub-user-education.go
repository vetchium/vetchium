package education

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func ListHubUserEducation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listHubUserEducationReq employer.ListHubUserEducationRequest
		err := json.NewDecoder(r.Body).Decode(&listHubUserEducationReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &listHubUserEducationReq) {
			h.Dbg(
				"invalid request",
				"listHubUserEducationReq",
				listHubUserEducationReq,
			)
			return
		}

		h.Dbg("validated", "listHubUserEducationReq", listHubUserEducationReq)

		educations, err := h.DB().
			ListHubUserEducation(r.Context(), listHubUserEducationReq)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("failed to list education", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to list education", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(educations)
		if err != nil {
			h.Dbg("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
