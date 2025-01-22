package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetEmployerInterviewsByCandidacy(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetEmployerInterviewsByCandidacy")
		var getInterviewsReq employer.GetEmployerInterviewsByCandidacyRequest
		err := json.NewDecoder(r.Body).Decode(&getInterviewsReq)
		if err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getInterviewsReq) {
			h.Dbg("validation failed", "getInterviewsReq", getInterviewsReq)
			return
		}
		h.Dbg("validated", "getInterviewsByCandidacyReq", getInterviewsReq)

		interviews, err := h.DB().
			GetEmployerInterviewsByCandidacy(r.Context(), getInterviewsReq)
		if err != nil {
			h.Dbg("error getting interviews", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got interviews", "interviews", interviews)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(interviews)
		if err != nil {
			h.Dbg("error encoding interviews", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
