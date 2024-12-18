package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetInterviewsByOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetInterviewsByOpening")
		var getInterviewsReq employer.GetInterviewsByOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&getInterviewsReq)
		if err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getInterviewsReq) {
			h.Dbg("validation failed", "getInterviewsReq", getInterviewsReq)
			return
		}
		h.Dbg("validated", "getInterviewsByOpeningReq", getInterviewsReq)

		interviews, err := h.DB().
			GetInterviewsByOpening(r.Context(), getInterviewsReq)
		if err != nil {
			h.Dbg("error getting interviews", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got interviews", "interviews", interviews)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(interviews)
	}
}
