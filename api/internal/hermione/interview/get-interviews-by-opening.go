package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/common"
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

		if getInterviewsReq.Limit == 0 {
			getInterviewsReq.Limit = 40
		}

		interviews, err := h.DB().
			GetInterviewsByOpening(r.Context(), getInterviewsReq)
		if err != nil {
			if errors.Is(err, db.ErrInvalidPaginationKey) {
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"pagination_key"},
				})
				if err != nil {
					h.Err("error encoding validation errors", "error", err)
					// This will cause suprefluous error because of multiple headers
					// being written. But this code is unlikely to execute.
					http.Error(w, "", http.StatusInternalServerError)
					return
				}
				return
			}

			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("no opening found", "opening_id", getInterviewsReq)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			h.Dbg("error getting interviews", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got interviews", "interviews", interviews)

		err = json.NewEncoder(w).Encode(interviews)
		if err != nil {
			h.Err("error encoding interviews", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
