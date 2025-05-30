package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
)

func GetEmployerCandidacyInfo(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetEmployerCandidacyInfo")

		var getCandidacyInfoReq common.GetCandidacyInfoRequest
		err := json.NewDecoder(r.Body).Decode(&getCandidacyInfoReq)
		if err != nil {
			h.Dbg("Error decoding request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		candidacyInfo, err := h.DB().
			GetEmployerCandidacyInfo(r.Context(), getCandidacyInfoReq)
		if err != nil {
			h.Dbg("Error getting candidacy info", "error", err)
			if err == db.ErrNoCandidacy {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("got candidacy info", "candidacyInfo", candidacyInfo)
		err = json.NewEncoder(w).Encode(candidacyInfo)
		if err != nil {
			h.Dbg("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
