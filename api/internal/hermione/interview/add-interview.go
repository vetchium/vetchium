package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func AddInterview(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddInterview")
		var addInterviewReq employer.AddInterviewRequest
		if err := json.NewDecoder(r.Body).Decode(&addInterviewReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addInterviewReq) {
			h.Dbg("validation failed", "addInterviewReq", addInterviewReq)
			return
		}

		h.Dbg("validated", "addInterviewReq", addInterviewReq)
	})
}
